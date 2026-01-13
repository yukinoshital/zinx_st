package apis

import (
	"fmt"
	"zinx/mmo_game/core"
	"zinx/mmo_game/pb"
	"zinx/ziface"
	"zinx/znet"

	"github.com/golang/protobuf/proto"
)

// 世界聊天的路由业务
type WorldChatApi struct {
	znet.Route
}

func (wc WorldChatApi) Handle(request ziface.Irequest) {
	//1、解析客户端传递的proto协议
	proto_msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("tale unmarshal err:", err)
		return
	}

	//2、当前聊天数据是那个玩家发送的
	pid, err := request.GetConnection().GetProperty("pid")

	//3、根据pid得到对应的player对象
	player := core.WorldManagerObj.GetPlayPid(pid.(int32))

	player.Talk(proto_msg.Content)

}
