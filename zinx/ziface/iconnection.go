package ziface

import (
	"net"
)

type Iconnection interface {
	Start() 
	Stop() 
	GetTcpConnection() *net.TCPConn
	GetConnId() uint32
	GetAddr() net.Addr
	Send(uint32,[]byte) error


	SetProperty(string,interface{})
	GetProperty(string) (interface{},error)
	RemoveProperty(string)
}

type HandleApi func(*net.TCPConn,[]byte,int) error