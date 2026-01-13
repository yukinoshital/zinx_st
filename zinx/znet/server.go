package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	IpVersion string
	Ip string
	Port int 
	Name string
	Route ziface.Iroute
	MsgHander ziface.ImsgHander
	ConMgr ziface.IconnManager
	OnConnStart func(conn ziface.Iconnection)
	OnConnStop func(conn ziface.Iconnection)
}

func NewServer(name string) ziface.Iserver {
	s := &Server{
		IpVersion: "tcp4",
		Ip: utils.GlobalObject.Host,
		Port: utils.GlobalObject.TcpPort,
		Name: utils.GlobalObject.Name,
		Route: nil,
		MsgHander: NewMsgHander(),
		ConMgr: NewConnManager(),
	}

	return s
}

func (s *Server) GetConnMgr() ziface.IconnManager {
	return s.ConMgr
}

func (s *Server) AddRoute(id uint32,route ziface.Iroute) {
	s.MsgHander.AddRouter(id,route)
	fmt.Println("add route success")
}


func (s *Server) Start() {
	fmt.Println("[server start]")
	fmt.Printf("srevername:%s, listen ip:%s,port:%d\n",
	utils.GlobalObject.Name,utils.GlobalObject.Host,utils.GlobalObject.TcpPort)
	fmt.Printf("zinxversion:%s,maxconn:%d,maxpackagesize:%d\n",
		utils.GlobalObject.Version,utils.GlobalObject.MaxConn,utils.GlobalObject.MaxPackageSize)	
	go func(){
		s.MsgHander.StartWorkPool()
		userAddr,err := net.ResolveTCPAddr(s.IpVersion,fmt.Sprintf("%s:%d",s.Ip,s.Port))
		if err != nil {
			fmt.Println("reslover addr err ",err)
			return
		}

		listener,err := net.ListenTCP(s.IpVersion,userAddr)
		if err != nil {
			fmt.Println("listen err:",err)
			return
		}

		var cid uint32
		cid = 0 

		for {
			conn,err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept err:",err)
				return
			}

			dealconn := NewConnection(s,conn,cid,s.MsgHander)
			cid++
			dealconn.Start()
		}
	}()

}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name " , s.Name)
	s.ConMgr.ClearConn()
}

func (s *Server) Run() {
	s.Start()
	select {}
}

func (s *Server) SetOnConnStart(hookFunc func(ziface.Iconnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(ziface.Iconnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.Iconnection) {
	if s.OnConnStart != nil {
		fmt.Println("----->CallOnConnStart<-----")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.Iconnection) {
	if s.OnConnStop != nil {
		fmt.Println("---->CallOnConnStop<-----")
		s.OnConnStop(conn)
	}
}