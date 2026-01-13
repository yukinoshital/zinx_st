package ziface

type Imessage interface {
	GetMsgId() uint32
	GetMsgLen() uint32
	GetMsgInfo() []byte


	SetMsgId(uint32)
	SetMsgLen(uint32)
	SetMsgInfo([]byte)
}