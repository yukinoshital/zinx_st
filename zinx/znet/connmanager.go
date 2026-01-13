package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	connections map[uint32] ziface.Iconnection
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32] ziface.Iconnection),
	}
}

func (cm *ConnManager) Add(conn ziface.Iconnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[conn.GetConnId()] = conn
	fmt.Println("connection add to ConnManager successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Remove(conn ziface.Iconnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, conn.GetConnId())
	fmt.Println("connection remove to ConnManager successfully: conn num = ", cm.Len())

}

func (cm *ConnManager) Get(id uint32) (ziface.Iconnection,error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn,ok := cm.connections[id];ok {
		return conn,nil
	} else {
		return nil,errors.New("connection is not found")
	}

}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	for connid,conn := range cm.connections {
		conn.Stop()
		delete(cm.connections,connid)
	}
		
	fmt.Println("Clear All Connections successfully: conn num = ", cm.Len())
		

}