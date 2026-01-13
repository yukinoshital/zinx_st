package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Println("[client start...]")
	conn,err := net.Dial("tcp","127.0.0.1:8887")
	if err != nil {
		fmt.Println("dial err:",err)
		return
	}
	for {
		dp := znet.NewDataPack()
		binarydata,err := dp.Pack(znet.NewMessage(0,[]byte("client test message0")))
		if err != nil {
			fmt.Println("pack err:",err)
			return
		}
		if _,err := conn.Write(binarydata);err != nil {
			fmt.Println("write err:",err)
			return
		}

		buf := make([]byte,dp.GetMsgHead())
		if _,err := io.ReadFull(conn,buf);err != nil {
			fmt.Println("read full err:",err)
			return
		}

		Msg,err := dp.UnPack(buf)
		if err != nil {
			fmt.Println("unpack err:",err)
			return
		}

		if Msg.GetMsgLen() > 0 {
			msg := Msg.(*znet.Message)
			msg.MsgInfo = make([]byte, msg.MsgLen)
			if _,err := io.ReadFull(conn,msg.MsgInfo);err != nil {
				fmt.Println("read full err:",err)
				return
			}

			fmt.Println("recv from msgid=",msg.MsgId," msglen=",msg.MsgLen," msginfo=",string(msg.MsgInfo))
		}

		time.Sleep(1*time.Second)
	}

}