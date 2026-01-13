package znet

type Message struct {
	MsgId uint32
	MsgLen uint32
	MsgInfo []byte
}

func NewMessage(id uint32,data []byte) *Message {
	return &Message{
		MsgId: id,
		MsgLen: uint32(len(data)),
		MsgInfo: data,
	}
}
 

func (m *Message) GetMsgId() uint32 {
	return m.MsgId
}

func (m *Message) GetMsgLen() uint32 {
	return m.MsgLen
}

func (m *Message) GetMsgInfo() []byte {
	return m.MsgInfo
}


func (m *Message) SetMsgId(id uint32) {
	m.MsgId = id
}

func (m *Message) SetMsgLen(msglen uint32) {
	m.MsgLen = msglen
}

func (m *Message) SetMsgInfo(data []byte) {
	m.MsgInfo = data
}