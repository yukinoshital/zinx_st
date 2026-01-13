package ziface

type IconnManager interface {
	Add(conn Iconnection)
	Remove(conn Iconnection)
	Get(uint32) (Iconnection,error)
	Len() int
	ClearConn()
}