package core

import (
	"fmt"
	"math/rand"
	"sync"
	"zinx/mmo_game/pb"
	"zinx/ziface"

	"github.com/golang/protobuf/proto"
)

// 玩家对象
type Player struct {
	Pid  uint32             //玩家id
	Conn ziface.Iconnection //当前玩家的链接
	X    float32            //平面x坐标
	Y    float32            //高度
	Z    float32            //平面Y坐标
	V    float32            //旋转的0·360角度
}

// player id 生成器
var PidGen uint32 = 1
var IdLock sync.Mutex

// 创建玩家的方法
func NewPlayer(conn ziface.Iconnection) *Player {
	//生成一个玩家id
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	//创建一个玩家对象
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), //随机在160坐标点 基于x轴若干偏移
		Y:    0,
		Z:    float32(140 + rand.Intn(20)),
		V:    0,
	}

	return p
}

// 提供一个发送给客户端消息的方法 主要是将pb的protobuf数据序列化之后，在调用zinx框架的SendMsg方法
func (p *Player) SendMsg(msgid uint32, data proto.Message) {
	//将proto的消息data序列化
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal err:", err)
		return
	}

	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}
	if err := p.Conn.Send(msgid, msg); err != nil {
		fmt.Println("player send msg err")
		return
	}

	return

}

// 告知客户端玩家pid 同步已生成的玩家id给客户端
func (p *Player) SyncPid() {
	//组建msgid：0的proto数据
	data := &pb.SyncPid{
		Pid: int32(p.Pid),
	}

	p.SendMsg(1, data)
}

// 广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	//组建msgid：200的proto数据
	data := &pb.BroadCast{
		Pid: int32(p.Pid),
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	p.SendMsg(200, data)
}

// 玩家广播世界聊天信息
func (p *Player) Talk(content string) {
	//组建msgid：200的proto数据
	proto_msg := &pb.BroadCast{
		Pid: int32(p.Pid),
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	//得到当前世界所有玩家的公告
	players := WorldManagerObj.GetAllPlayers()

	//广播给所有在线玩家
	for _, p := range players {
		p.SendMsg(200, proto_msg)
	}
}

// 同步玩家上线的位置消息
func (p *Player) SyncSurrounding() {
	//1、获取当前玩家周围的玩家有哪些
	pids := WorldManagerObj.AoiMgr.GetPIDsByPos(p.X, p.Z)
	player := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		player = append(player, WorldManagerObj.GetPlayPid(int32(pid)))
	}
	//2、当前玩家的位置信息通过msgid 200 发送给周边玩家（让其他玩家看到自己）
	//2.1 组建msgid：200的proto数据
	proto_msg := &pb.BroadCast{
		Pid: int32(p.Pid),
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//2.2 分别给周边的全部玩家发送当前玩家的位置信息
	for _, player := range player {
		player.SendMsg(200, proto_msg)
	}

	//3、将周围的全部玩家的位置信息发送给当前的玩家客户端（让自己看到周边玩家）
	//3.1 制作msgid 202 proto数据
	players_proto_msg := make([]*pb.Player, 0, len(player))
	for _, player := range player {
		//3.2 制作一个message
		p := &pb.Player{
			Pid: int32(player.Pid),
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}

		players_proto_msg = append(players_proto_msg, p)
	}
	//封装syncplayer protobuf数据
	SyncPlayers_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:],
	}
	//3.2 将周边玩家的位置发送给当前玩家
	p.SendMsg(202, SyncPlayers_proto_msg)
}

// 广播当前玩家的位置信息
func (p *Player) UnpdatePos(x, y, z, v float32) {
	//更新当前player的玩家坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v
	//组建proto广播协议 msgid 200 tp 4
	proto_msg := &pb.BroadCast{
		Pid: int32(p.Pid),
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	//获取当前玩家周边玩家
	players := p.GetSurroundingPlayer()
	//一次给每个玩家对应的客户端发送当前玩家更新的信息
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

// 获取当前玩家的周边玩家aoi九宫格之内的玩家
func (p *Player) GetSurroundingPlayer() []*Player {
	//得到当前九宫格内所有玩家的pid
	pids := WorldManagerObj.AoiMgr.GetPIDsByPos(p.X, p.Z)

	//将所有的pid对应的player放到player切片中
	players := make([]*Player, 0, len(pids))

	for _, pid := range pids {
		players = append(players, WorldManagerObj.GetPlayPid(int32(pid)))
	}

	return players
}

func (p *Player) Offline() {
	//得到当前玩家周边九宫格内有哪些玩家
	players := p.GetSurroundingPlayer()
	//给周围玩家广播msgid201消息
	proto_msg := &pb.SyncPid{
		Pid: int32(p.Pid),
	}

	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}

	WorldManagerObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldManagerObj.RemovePlayerByPid(int32(p.Pid))
}
