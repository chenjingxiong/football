package football

import (
	"fmt"

//	"time"
)

const (
	teamLevelMatchCDType = 1 ///征战八方平局或输了添加的CD类型
)

type IProcessMgr interface {
	IGameMgr
	//	Init(teamID int) bool                                             ///加载指定球队id
	AddProcess(centerType int, processLevel int) int                  ///在指定处理中心中添加一个处理流程
	GetProcessInfoList(centerType int, processID int) ProcessInfoList ///得到指定类型处理中心中所有处理内容
	GetProcessMaxIndex(centerType int) int                            ///得到指定处理中心最大索引号,索引号从1
	//UpdateProcess(centerType int, processID int, objID int, expireTime int) bool ///更新处理对象
	Update(now int, client IClient)                    ///处理中心更新自己状态
	GetProcess(centerType int, processID int) IProcess ///得到处理对象
	FindProcessByObjID(centerType int, objID int) int  ///根据objID查找处理进程对象id
	GetProcessList(centerType int) ProcessList         ///得到处理对象列表
	SetTeamCD(cdType int, expireSecs int) IProcess     ///设置球队指定类型CD
	isExpireTeamCD(cdType int) bool                    ///判断球队cd是否已到期
	GetTeamCDUTC(cdType int) int                       ///得到球队指定类型CD的utc时间
}

type ProcessInfoList []ProcessInfo         ///任务处理对象列表
type ProcessList map[int]IProcess          ///任务处理对象图表
type ProcessCenterList map[int]ProcessList ///处理中心列表
type ProcessMgr struct {                   ///被球队拥有处理管理器,用于时间推移给效果的系统
	GameMgr
	processCenterList ProcessCenterList ///处理列表
}

func NewProcessMgr(teamID int) IProcessMgr {
	processMgr := new(ProcessMgr)
	if processMgr.Init(teamID) == false {
		return nil
	}
	return processMgr
}

func (self *ProcessMgr) GetType() int { ///得到管理器类型
	return mgrTypeProcessMgr ///关卡管理器
}

func (self *ProcessMgr) GetProcessMaxIndex(centerType int) int { ///得到指定处理中心最大索引号,索引号从1
	maxIndex := 0
	if nil == self.processCenterList[centerType] {
		return 0
	}
	for _, v := range self.processCenterList[centerType] {
		processInfo := v.GetInfo()
		if processInfo.Pos > maxIndex {
			maxIndex = processInfo.Pos
		}
	}
	return maxIndex
}

func (self *ProcessMgr) starTrainUpdate(now int, client IClient) { ///处理中心更新自己状态
	//const starTrainAwardExpInterval = 60 ///每60秒奖励球员经验
	//if now%starTrainAwardExpInterval > 0 {
	//	return ///每60秒判断一次
	//}
	if nil == self.processCenterList[ProcessTypeStarTrain] {
		return ///没数据不处理
	}
	syncObjectList := SyncObjectList{}
	team := client.GetTeam()
	syncMgr := client.GetSyncMgr()
	teamLevel := team.GetLevel()
	for _, v := range self.processCenterList[ProcessTypeStarTrain] {
		processInfo := v.GetInfo()
		if processInfo.ObjID <= 0 {
			continue ///空槽跳过
		}
		if IsExpireTime(processInfo.ExpireTime) == true {
			continue ///忽略已过期的训练项目
		}
		//remainSec:=processInfo.ExpireTime-now
		//nowTime := time.Now()
		//expireTime := time.Unix(int64(processInfo.ExpireTime), 0)
		//duration := expireTime.Sub(nowTime)
		//GetServer().GetLoger().Info("%v", duration)
		//duration := time.Now() - processInfo.ExpireTime
		if v.CanProcess() == false {
			continue ///未到处理时间
		}
		awardExp := v.CalcTrainAwardExp(teamLevel, false)
		if awardExp <= 0 {
			continue ///无效的经验不处理
		}
		star := team.GetStar(processInfo.ObjID)
		star.AwardExp(awardExp)
		if star.isMaxLevel() == true {
			processInfo.ExpireTime = Now() ///球员满级后直接处于完成状态
		}
		processInfo.NextProcessTime = now + trainAwardExpInterval ///记录发放经验的时间
		processInfo.Param1 += awardExp                            ///累加获得经验数
		syncObjectList = append(syncObjectList, star)
	}
	if len(syncObjectList) > 0 { ///同步玩员属性变更
		syncMgr.SyncObjectArray("starTrainUpdate", syncObjectList)
	}
}

func (self *ProcessMgr) Update(now int, client IClient) { ///处理中心更新自己状态
	///处理球员训练更新自身状态
	self.starTrainUpdate(now, client) ///处理中心更新自己状态
}

func (self *ProcessMgr) GetProcessList(centerType int) ProcessList { ///得到处理对象列表
	if nil == self.processCenterList[centerType] {
		return nil
	}
	return self.processCenterList[centerType]
}

func (self *ProcessMgr) GetProcess(centerType int, processID int) IProcess { ///得到处理对象
	if nil == self.processCenterList[centerType] {
		return nil
	}
	return self.processCenterList[centerType][processID]
}

///根据objID查找处理进程对象id
func (self *ProcessMgr) FindProcessByObjID(centerType int, objID int) int {
	processID := 0
	if nil == self.processCenterList[centerType] {
		return 0
	}
	for _, v := range self.processCenterList[centerType] {
		processInfo := v.GetInfo()
		if processInfo.ObjID == objID {
			processID = processInfo.ID
			break
		}
	}
	return processID
}

//func (self *ProcessMgr) UpdateProcess(centerType int, processID int, objID int, expireTime int) bool { ///更新处理对象
//	if nil == self.processCenterList[centerType] {
//		return false
//	}
//	processInfo, ok := self.processCenterList[centerType][processID]
//	if false == ok {
//		return false ///不存在的ProcessID
//	}
//	if processInfo.ObjID == objID && processInfo.ExpireTime == expireTime {
//		return true
//	}
//	updateProcessQuery := fmt.Sprintf("update %s set objid=%d,expiretime=%d where id=%d", tableProcessCenter, objID, expireTime, processID) ///生成更新SQL语句
//	///执行更新语句
//	_, rowsProcessAffected := GetServer().GetDynamicDB().Exec(updateProcessQuery)
//	if rowsProcessAffected <= 0 {
//		GetServer().GetLoger().Warn("ProcessMgr UpdateProcess updateProcessQuery fail!")
//		return false
//	}
//	///更新内存
//	processInfo.ObjID = objID
//	processInfo.ExpireTime = expireTime
//	return rowsProcessAffected > 0
//}

func (self *ProcessMgr) AddProcess(centerType int, processLevel int) int { ///加载指定球队id
	if nil == self.processCenterList[centerType] {
		self.processCenterList[centerType] = make(ProcessList)
	}
	maxIndex := self.GetProcessMaxIndex(centerType) ///得到最大索引号
	insertNewProcessQuery := fmt.Sprintf("Insert %s set teamid=%d,type=%d,pos=%d,level=%d",
		tableProcessCenter, self.team.GetID(), centerType, maxIndex+1, processLevel)
	///执行插入语句
	lastInsertProcessID, _ := GetServer().GetDynamicDB().Exec(insertNewProcessQuery)
	if lastInsertProcessID <= 0 {
		GetServer().GetLoger().Warn("TeamCreateMsg AddPorcess insertNewProcessQuery fail! centerType:%d", centerType)
		return 0
	}
	///创建Process对象
	loadProcessQuery := fmt.Sprintf("select * from %s where id=%d limit 1", tableProcessCenter, lastInsertProcessID)
	processInfo := new(ProcessInfo)
	ok := GetServer().GetDynamicDB().fetchOneRow(loadProcessQuery, processInfo)
	if false == ok {
		return 0
	}
	self.processCenterList[centerType][lastInsertProcessID] = NewProcess(processInfo)
	return lastInsertProcessID
}

func (self *ProcessMgr) SaveInfo() { ///保存数据
	for _, v := range self.processCenterList {
		for _, t := range v {
			t.Save()
		}
	}
}

func (self *ProcessMgr) GetProcessInfoList(centerType int, processID int) ProcessInfoList { ///得到指定类型处理中心中所有处理内容
	if nil == self.processCenterList[centerType] {
		return nil
	}
	processInfoList := ProcessInfoList{}
	for _, v := range self.processCenterList[centerType] {
		processInfo := v.GetInfo()
		if processID > 0 && processID != processInfo.ID {
			continue ///processID非零时要求匹配ID
		}
		processInfoList = append(processInfoList, *processInfo)
	}
	return processInfoList
}

func (self *ProcessMgr) Init(teamID int) bool { ///加载指定球队id
	self.processCenterList = make(ProcessCenterList)
	//	self.teamID = teamID ///存放自己的球队id
	processCenterQuery := fmt.Sprintf("select * from %s where teamid=%d limit 3000", tableProcessCenter, teamID)
	processInfo := new(ProcessInfo)
	processInfoList := GetServer().GetDynamicDB().fetchAllRows(processCenterQuery, processInfo)
	if nil == processInfoList {
		return false
	}
	for i := range processInfoList {
		processInfo = processInfoList[i].(*ProcessInfo)
		centerType := processInfo.Type
		processID := processInfo.ID
		if nil == self.processCenterList[centerType] {
			self.processCenterList[centerType] = make(ProcessList)
		}
		self.processCenterList[centerType][processID] = NewProcess(processInfo)
	}
	return true
}

func (self *ProcessMgr) GetTeamCDUTC(cdType int) int { ///得到球队指定类型CD的utc时间
	result := 0
	if nil == self.processCenterList[ProcessTypeTeamCD] {
		return 0
	}
	for _, v := range self.processCenterList[ProcessTypeTeamCD] {
		processInfo := v.GetInfo()
		if processInfo.ObjID == cdType {
			result = processInfo.ExpireTime
			break
		}
	}
	return result
}

func (self *ProcessMgr) SetTeamCD(cdType int, expireSecs int) IProcess { ///设置球队指定类型CD
	if nil == self.processCenterList[ProcessTypeTeamCD] {
		self.processCenterList[ProcessTypeTeamCD] = make(ProcessList)
	}
	process := IProcess(nil)
	for _, v := range self.processCenterList[ProcessTypeTeamCD] {
		processInfo := v.GetInfo()
		if processInfo.ObjID == cdType {
			process = v
			break
		}
	}
	if nil == process {
		processID := self.AddProcess(ProcessTypeTeamCD, 0)
		process = self.GetProcess(ProcessTypeTeamCD, processID)
	}
	processInfo := process.GetInfo()
	processInfo.ObjID = cdType
	processInfo.ExpireTime = Now() + expireSecs ///设置过期时间
	return process
}

func (self *ProcessMgr) isExpireTeamCD(cdType int) bool { ///判断球队cd是否已到期
	result := true
	if nil == self.processCenterList[ProcessTypeTeamCD] {
		return true
	}
	for _, v := range self.processCenterList[ProcessTypeTeamCD] {
		processInfo := v.GetInfo()
		if processInfo.ObjID == cdType {
			result = IsExpireTime(processInfo.ExpireTime)
			break
		}
	}
	return result
}
