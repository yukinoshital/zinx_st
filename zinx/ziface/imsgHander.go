package ziface

type ImsgHander interface {
	DoMsgHandler(Irequest)
	AddRouter(uint32,Iroute)
	StartWorkPool()
	SendMsgToTaskQueue(Irequest)
}