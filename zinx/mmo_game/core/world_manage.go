package core

import "sync"

/*
当前游戏世界的管理模块
*/

type WorldManager struct {
	//AOIManager 当前世界地图的aoi管理模块
	AoiMgr *AOIManager
	//当前全部在线的player集合
	Players map[int32] *Player
	//保护player集合的锁
	pLock sync.RWMutex
}

var WorldManagerObj *WorldManager
//初始化方法
func init() {
	WorldManagerObj = &WorldManager {
		//创建世界地图规划
		AoiMgr:NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNT_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNT_Y),
		//初始化pkayer集合
		Players:make(map[int32]*Player),
	}

}
//添加一个玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[int32(player.Pid)] = player
	wm.pLock.Unlock()

	//将玩家添加到aoimanager中
	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

//删除一个玩家
func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	player := wm.Players[pid]
	wm.AoiMgr.RemoveFromGridByPos(int(pid), player.X, player.Z)

	wm.pLock.Lock()
	delete(wm.Players,pid)
	wm.pLock.Unlock()

}
//通过玩家id查询player对象
func (wm *WorldManager) GetPlayPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	return wm.Players[pid]	
}

//获取全部在线的player玩家
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player,0)
	
	for _,p := range wm.Players {
		players = append(players, p)
	}

	return players
}