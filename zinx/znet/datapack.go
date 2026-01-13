package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct {

}

func NewDataPack() *DataPack {
	return &DataPack{}
}


func (dp *DataPack)	GetMsgHead() uint32 {
	return 8
}

func (dp *DataPack) Pack(msg ziface.Imessage) ([]byte,error) {
	databuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(databuff,binary.LittleEndian,msg.GetMsgId());err != nil {
		return nil,err
	}

	if err := binary.Write(databuff,binary.LittleEndian,msg.GetMsgLen());err != nil {
		return nil,err
	}

	if err := binary.Write(databuff,binary.LittleEndian,msg.GetMsgInfo());err != nil {
		return nil,err
	}

	return databuff.Bytes(),nil
}

func (dp *DataPack) UnPack(data []byte) (ziface.Imessage,error) {
	databuff := bytes.NewReader(data)
	msg := &Message{}
	if err := binary.Read(databuff,binary.LittleEndian,&msg.MsgId);err != nil {
		return nil,err
	}

	if err := binary.Read(databuff,binary.LittleEndian,&msg.MsgLen);err != nil {
		return nil,err
	}

	if (msg.MsgLen > 0 && msg.MsgLen > utils.GlobalObject.MaxPackageSize) {
		return nil,errors.New("pack too large")
	}

	return msg,nil
}