package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRoute struct {
	znet.Route
}



func (this *PingRoute) Handle(request ziface.Irequest) {
	fmt.Println("call route handle")
	fmt.Println("recv from client : msgId=", request.GetDataId(), ", data=", string(request.GetData()))
	if err := request.GetConnection().Send(0,[]byte("ping ping ping"));err != nil {
		fmt.Println("write err:",err)
		return
	}
}


type HelloRoute struct {
	znet.Route
}



func (this *HelloRoute) Handle(request ziface.Irequest) {
	fmt.Println("call route handle")
	fmt.Println("recv from client : msgId=", request.GetDataId(), ", data=", string(request.GetData()))
	if err := request.GetConnection().Send(1,[]byte("hello zinx"));err != nil {
		fmt.Println("write err:",err)
		return
	}
}

//创建连接的时候执行
func DoConnectionBegin(conn ziface.Iconnection) {
	fmt.Println("DoConnectionBegin is called...")
	fmt.Println("Set conn Name, Home done!")
	conn.SetProperty("Name", "Aceld")
	conn.SetProperty("Home", "https://www.jianshu.com/u/35261429b7f1")
	err := conn.Send(2,[]byte("DoConnection Begin"))
	if err != nil {
		fmt.Println("send err:",err)
		return
	}
}

//连接断开的时候执行
func DoConnectionLost(conn ziface.Iconnection) {
	if name, err:= conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home = ", home)
	fmt.Println("DoConneciotnLost is Called ... ")
	}
}



func main() {
	s := znet.NewServer("zinxserver")
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	s.AddRoute(0,&PingRoute{})
	s.AddRoute(1,&HelloRoute{})
	s.Run()
	select {}
}