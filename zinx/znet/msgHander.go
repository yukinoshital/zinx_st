package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHander struct {
	Apis map[uint32] ziface.Iroute
	WorkPoolSize uint32
	TaskQueue []chan ziface.Irequest
}

func NewMsgHander() *MsgHander {
	return &MsgHander{
		Apis: make(map[uint32] ziface.Iroute),
		WorkPoolSize: utils.GlobalObject.WorkPoolSize,
		TaskQueue: make([]chan ziface.Irequest, utils.GlobalObject.WorkPoolSize),
	}
}

func (mh *MsgHander) DoMsgHandler(request ziface.Irequest) {
	hander,ok := mh.Apis[request.GetDataId()]
	if !ok {
		fmt.Println("api msgid:",request.GetDataId()," is not found")
		return
	}

	hander.PreHandle(request)
	hander.Handle(request)
	hander.AfterHandle(request)
}


func (mh *MsgHander) AddRouter(id uint32,route ziface.Iroute) {
	if _,ok := mh.Apis[id];ok {
		fmt.Println("api msgid:",id," is exist")
		return
	}

	mh.Apis[id] = route
	fmt.Println("api msgid:",id," add success")
}

func (mh *MsgHander) StartWorkPool() {
	for i:=0; i<int(mh.WorkPoolSize);i++ {
		mh.TaskQueue[i] = make(chan ziface.Irequest,utils.GlobalObject.MaxWorkTaskLen)
		go mh.StartOneWorker(i,mh.TaskQueue[i])
	}
}

func (mh *MsgHander) StartOneWorker(workid int,taskqueue chan ziface.Irequest) {
	fmt.Println("workid:",workid," is start")
	for {
		select {
		case request := <- taskqueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHander) SendMsgToTaskQueue(request ziface.Irequest) {
	workerID := request.GetConnection().GetConnId() % mh.WorkPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnId(),
	" request msgID=", request.GetDataId(), "to workerID=", workerID)
	mh.TaskQueue[workerID] <- request
}