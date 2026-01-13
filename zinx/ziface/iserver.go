package ziface

type Iserver  interface {
	Start()
	Stop()
	Run()
	AddRoute(uint32,Iroute)
	GetConnMgr() IconnManager

	SetOnConnStart(func(Iconnection))
	SetOnConnStop(func(Iconnection))
	CallOnConnStart(Iconnection)
	CallOnConnStop(Iconnection)
}