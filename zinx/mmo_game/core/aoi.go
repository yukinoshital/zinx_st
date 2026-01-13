package core

import "fmt"

/*
   AOI管理模块
*/

//定义aoi的边界值
const (
	AOI_MIN_X int = 85
	AOI_MAX_X int = 410
	AOI_CNT_X int = 10
	AOI_MIN_Y int = 75
	AOI_MAX_Y int = 400
	AOI_CNT_Y int = 20
)

type AOIManager struct {
	MinX int //区域左边界坐标
	MaxX int //区域右边界坐标
	CntX int //x方向格子的数量
	MinY int //区域上边界坐标
	MaxY int //区域下边界坐标
	CntY int //y方向的格子数量
	grids map[int]*Grid //当前区域中都有哪些格子，key=格子ID， value=格子对象
}

func NewAOIManager(minX,maxX,cntX,minY,maxY,cntY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX: minX,
		MaxX: maxX,
		CntX: cntX,
		MinY: minY,
		MaxY: maxY,
		CntY: cntY,
		grids: make(map[int]*Grid),
	}

	//给AOI初始化区域中所有的格子
	for y:=0; y<cntY;y++ {
		for x:=0; x<cntX;x++ {
			//格子编号：id = idy *nx + idx  (利用格子坐标得到格子编号)
			gid := y*cntX + x
			//初始化一个格子放在AOI中的map里，key是当前格子的ID
			aoiMgr.grids[gid] = NewGrid(
				gid,
				aoiMgr.MinX + x*aoiMgr.GetGridWidth(),
				aoiMgr.MinX + (x+1)*aoiMgr.GetGridWidth(),
				aoiMgr.MinY + y*aoiMgr.GetGridLength(),
				aoiMgr.MinY + (y+1)*aoiMgr.GetGridLength(),
			)

		}
	}

	return aoiMgr
}

//得到每个格子在x轴方向的宽度
func (m *AOIManager) GetGridWidth() int {
	return (m.MaxX - m.MinX) / m.CntX
}

//得到每个格子在x轴方向的长度
func (m *AOIManager) GetGridLength() int {
	return (m.MaxY - m.MinY) / m.CntY
}

//打印信息方法
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManagr:\nminX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d\n Grids in AOI Manager:\n",
		m.MinX, m.MaxX, m.CntX, m.MinY, m.MaxY, m.CntY)
	
	for _,grid := range m.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}


//根据格子的gID得到当前周边的九宫格信息
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grid []*Grid) {
	//判断gID是否存在
	if _,ok := m.grids[gID]; !ok {
		return
	}

	//将当前gid添加到九宫格中
	grid = append(grid, m.grids[gID])

	//根据gid得到当前格子所在的X轴编号
	idx := gID % m.CntX

	if idx > 0 {
		grid = append(grid, m.grids[gID-1])
	}

	if idx < m.CntX - 1 {
		grid = append(grid, m.grids[gID+1])
	}

	//将x轴当前的格子都取出，进行遍历，再分别得到每个格子的上下是否有格子
	//得到当前x轴的格子id集合
	gidX := make([]int,0,len(grid))
	for _,v := range grid {
		gidX = append(gidX, v.Gid)
	}

	for _,v := range gidX {
		idy := v/m.CntX
		if idy > 0 {
			grid = append(grid, m.grids[v-m.CntX])
		}

		if idy < m.CntY - 1 {
			grid = append(grid, m.grids[v+m.CntX])
		}
	}

	return
}


//通过横纵坐标获取对应的格子ID
func (m *AOIManager) GetGidByPos(x,y float32) int {
	gx := (int(x) - m.MinX) / m.GetGridWidth()
	gy := (int(x) - m.MinY) / m.GetGridLength()

	return gy * m.CntX + gx

}


//通过横纵坐标得到周边九宫格内的全部PlayerIDs
func (m *AOIManager) GetPIDsByPos(x, y float32) (playerIDs []int) {
	//根据横纵坐标得到当前坐标属于哪个格子ID
	gID := m.GetGidByPos(x, y)

	//根据格子ID得到周边九宫格的信息
	grids := m.GetSurroundGridsByGid(gID)
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetAllPlayerIds()...)
		fmt.Printf("===> grid ID : %d, pids : %v  ====", v.Gid, v.GetAllPlayerIds())
	}

	return
}

//通过GID获取当前格子的全部playerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.grids[gID].GetAllPlayerIds()
	return
}


//移除一个格子中的PlayerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.grids[gID].RemovePlayer(pID)
}

//添加一个PlayerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].AddPlayer(pID)
}


//通过横纵坐标添加一个Player到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.AddPlayer(pID)
}



//通过横纵坐标把一个Player从对应的格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.RemovePlayer(pID)
}