package main

import (
	"fmt"
	"zinx/mmo_game/apis"
	"zinx/mmo_game/core"
	"zinx/ziface"
	"zinx/znet"
)

// 当前客户端建立连接之后的hook函数
func OnConnectionAdd(conn ziface.Iconnection) {
	//创建一个玩家对象
	player := core.NewPlayer(conn)

	//给客户端发送msgid：1的消息 同步playerid给客户端
	player.SyncPid()

	//给客户端发送msgid：200的消息 同步player的初始位置给客户端
	player.BroadCastStartPosition()

	//将新上线的玩家添加到world
	core.WorldManagerObj.AddPlayer(player)

	//将该连接绑定一个pid 玩家id的属性
	conn.SetProperty("pid", player.Pid)

	//同步周边玩家 告知当前玩家已经上线 广播当前玩家的位置
	player.SyncSurrounding()

	fmt.Println("=======>player id =", player.Pid, "<========")

}

// 给当前连接断开之前出发的hook函数
func OnConnectionLost(conn ziface.Iconnection) {

	pid, _ := conn.GetProperty("pid")
	player := core.WorldManagerObj.GetPlayPid(pid.(int32))

	//触发玩家下线的业务
	player.Offline()

	fmt.Println("======>player id =", pid, " offline...<=======")

}

func main() {
	//创建zinx server句柄
	s := znet.NewServer("Mmo Game Zinx")

	//连接创建和销毁的hook函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)
	//注册一些路由服务
	s.AddRoute(2, &apis.WorldChatApi{})

	s.AddRoute(3, &apis.MoveApi{})
	//启动服务
	s.Run()
}
