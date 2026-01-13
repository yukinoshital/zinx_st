package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	Gid int //格子id
	MinX int //格子左边界坐标
	MaxX int //格子右边界坐标
	MinY int //格子上边界坐标
	MaxY int //格子下边界坐标
	playerIds map[int]bool //当前格子内的玩家或者物体成员ID
	pIdLock sync.RWMutex //playerIDs的保护map的锁
}

//初始化一个格子
func NewGrid(gID,MinX,MaxX,MinY,MaxY int) *Grid {
	return &Grid{
		Gid: gID,
		MinX: MinX,
		MaxX: MaxX,
		MinY: MinY,
		MaxY: MaxY,
		playerIds: make(map[int]bool),
	}
}


//向当前格子中添加一个玩家
func (g *Grid) AddPlayer(pid int) {
	g.pIdLock.Lock()
	defer g.pIdLock.Unlock()

	g.playerIds[pid] = true
}

//从格子中删除一个玩家
func (g *Grid) RemovePlayer(pid int) {
	g.pIdLock.Lock()
	defer g.pIdLock.Unlock()
	delete(g.playerIds,pid)
}

//得到当前格子中所有的玩家
func (g *Grid) GetAllPlayerIds() (playerIDs []int) {
	g.pIdLock.RLock()
	defer g.pIdLock.RUnlock()

	for k,_ := range g.playerIds {
		playerIDs = append(playerIDs,k)
	}

	return
}

//打印信息方法
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",
		g.Gid, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIds)
}