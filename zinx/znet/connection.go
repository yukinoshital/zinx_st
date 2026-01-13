package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	TcpServer ziface.Iserver
	Conn *net.TCPConn
	Connid uint32
	isclose bool
	Handleapi ziface.HandleApi
	exitchan chan bool
	MsgHander ziface.ImsgHander
	msgChan chan []byte
	property map[string]interface{}
	propertylock sync.RWMutex
}

func NewConnection(server ziface.Iserver,conn *net.TCPConn,connid uint32,msghander ziface.ImsgHander) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn: conn,
		Connid: connid,
		isclose: false,
		MsgHander: msghander,
		exitchan: make(chan bool,1),
		msgChan: make(chan []byte),	
		property: make(map[string]interface{}),	
	}

	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func(c *Connection) StartRead() {
	fmt.Println("conn gorotinue is running")
	defer fmt.Printf("connid=%d,addr=%s",c.Connid,c.GetAddr().String())
	defer c.Stop()
	for {
		dp := NewDataPack()
		buf := make([]byte,dp.GetMsgHead())
		if _,err := io.ReadFull(c.Conn,buf);err != nil {
			fmt.Println("read full err:",err)
			return
		}

		Msg,err := dp.UnPack(buf)
		if err != nil {
			fmt.Println("unpack err:",err)
			return
		}

		var data []byte
		if Msg.GetMsgLen() > 0 {
			data = make([]byte,Msg.GetMsgLen())
			if _,err := io.ReadFull(c.Conn,data);err != nil {
				fmt.Println("read full err:",err)
				return
			}
		}

		Msg.SetMsgInfo(data)

		req := Request{
			conn: c,
			data: Msg,
		}

		if utils.GlobalObject.WorkPoolSize > 0 {
			c.MsgHander.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHander.DoMsgHandler(&req)
		}
		
	}
}

func (c *Connection) StartWrite() {
	fmt.Println("[write gorotinue is running]")
	defer fmt.Println(c.GetAddr().String(),"[conn writer stop]")
	for {
		select {
		case data := <- c.msgChan:
			if _,err := c.Conn.Write(data);err != nil {
				fmt.Println("weiter err:",err)
				return
			}
		case <- c.exitchan:
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("conn start")
	go c.StartRead()
	go c.StartWrite()
	c.TcpServer.CallOnConnStart(c)
}


func (c *Connection) Stop() {
	fmt.Println("conn stop")
	if c.isclose == true {
		return
	}

	c.isclose = true
	c.TcpServer.CallOnConnStop(c)
	c.Conn.Close()
	c.TcpServer.GetConnMgr().Remove(c)
	c.exitchan <- true
	close(c.exitchan)
	close(c.msgChan)
}


func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnId() uint32 {
	return c.Connid
}

func (c *Connection) GetAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(id uint32,data []byte) error {
	dp := NewDataPack()
	databinery,err := dp.Pack(NewMessage(id,data))
	if err != nil {
		fmt.Println("pack err:",err)
		return err
	}
	c.msgChan <- databinery
	return nil
}


func (c *Connection) SetProperty(key string,value interface{}) {
	c.propertylock.Lock()
	defer c.propertylock.Unlock()

	c.property[key] = value

}

func (c *Connection) GetProperty(key string) (interface{},error) {
	c.propertylock.RLock()
	defer c.propertylock.RUnlock()

	if value,ok := c.property[key];ok {
		return value,nil
	} else {
		return nil,errors.New("property not found")
	}
}

func(c *Connection) RemoveProperty(key string) {
	c.propertylock.Lock()
	defer c.propertylock.Unlock()

	delete(c.property,key)
}