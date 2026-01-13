package ziface

type Irequest interface {
	GetConnection() Iconnection
	GetData() []byte
	GetDataId() uint32
}