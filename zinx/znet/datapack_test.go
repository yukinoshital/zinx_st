package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)


func TestDataPack(t *testing.T) {

	listen,err := net.Listen("tcp","127.0.0.1:8888")
	if err != nil {
		fmt.Println("listen err:",err)
		return
	}

	go func(){
		for {
			conn,err := listen.Accept()
			if err != nil {
				fmt.Println("accept err:",err)
				return
			}

			go func(conn net.Conn){
				dp := NewDataPack()
				for {
					headData := make([]byte,dp.GetMsgHead())
					if _,err := io.ReadFull(conn,headData);err != nil {
						fmt.Println("read full err:",err)
						return
					}

					Msg,err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("unpack err:",err)
						return
					}

					if Msg.GetMsgLen() > 0 {
						msg := Msg.(*Message)
						msg.MsgInfo = make([]byte,msg.MsgLen)
						if _,err := io.ReadFull(conn,msg.MsgInfo); err != nil {
							fmt.Println("read full err:",err)
							return
						}

						fmt.Println("recv msgid:",msg.MsgId," msglen:",msg.MsgLen," msginfo:",string(msg.MsgInfo))
					}
				}
			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("dial err", err)
		return
	}

	dp := NewDataPack()

	msg1 := &Message{
		MsgId:   1,
		MsgLen: 4,
		MsgInfo: []byte{'z','i','n','x'},
	}

	senddata1,err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack err:",err)
		return
	}

	msg2 := &Message{
		MsgId: 2,
		MsgLen: 5,
		MsgInfo: []byte{'h','e','l','l','o'},
	}

	sendata2,err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("pack err:",err)
		return
	}

	senddata1 = append(senddata1,sendata2...)

	conn.Write(senddata1)
}