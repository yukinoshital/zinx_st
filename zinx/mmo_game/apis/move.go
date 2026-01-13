package apis

import (
	"fmt"
	"zinx/mmo_game/core"
	"zinx/mmo_game/pb"
	"zinx/ziface"
	"zinx/znet"

	"github.com/golang/protobuf/proto"
)

type MoveApi struct {
	znet.Route
}

func (m *MoveApi) Handle(request ziface.Irequest) {
	//解析客户端传递过来的proto协议
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move : Position Unmarshal error", err)
		return
	}
	//得到当前发送位置的是哪一个玩家
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProerty pid err:", err)
		return
	}

	fmt.Printf("Player id=%d, move(%f,%f,%f,%f)", pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
	//给其他玩家进行当前玩家位置的广播
	player := core.WorldManagerObj.GetPlayPid(pid.(int32))
	//广播并更新当前玩家的坐标
	player.UnpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
