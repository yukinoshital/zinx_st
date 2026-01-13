package ziface

type IdataPack interface {
	GetMsgHead() uint32
	Pack(Imessage) ([]byte,error)
	UnPack([]byte) (Imessage,error)
}