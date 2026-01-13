package znet

import "zinx/ziface"

type Request struct {
	conn ziface.Iconnection
	data ziface.Imessage
}

func(r *Request) GetConnection() ziface.Iconnection {
	return  r.conn
}

func(r *Request) GetData() []byte {
	return r.data.GetMsgInfo()
}

func (r *Request) GetDataId() uint32 {
	return r.data.GetMsgId()
}