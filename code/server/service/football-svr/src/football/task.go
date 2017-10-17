package football

import (
	"reflect"
	"time"
)

///任务信息
type TaskArray []*Task

//type ITask interface {
//	IDataUpdater
//	GetInfoPtr() *TaskInfo                                                  ///得到信息接口
//	UpdateNpcTeamTask(team *Team, userGoalCount int, npcGoalCount int) bool ///更新npc球队相关任务信息
//	GetID() int                                                             ///得到ID编号
//	GetReflectValue() reflect.Value                                         ///得到球队反射对象
//	IsComplete() bool                                                       ///判断任务是否已完成
//	GetTypeInfo() *TaskTypeStaticData                                       ///取得任务类型静态数据信息
//	UpdateDone(team *Team) bool                                             ///有修改返回true
//}

type TaskTypePtrList []*TaskTypeStaticData

type TaskInfo struct { ///任务信息,对应动态表中的dy_task
	ID         int `json:"id"`         ///任务id
	TeamID     int `json:"teamid"`     ///拥有任务信息的球队id
	Type       int `json:"type"`       ///任务类型',
	Progress   int `json:"progress"`   ///当前任务进度信息
	IsDone     int `json:"isdone"`     ///当前任务是否已完成 0未完成 1已完成
	IsAward    int `json:"isaward"`    ///当前任务是否已领取奖励
	Need1      int `json:"need1"`      ///完成条件,只有特殊任务的完成条件是动态的
	Need2      int `json:"need2"`      ///完成条件,只有特殊任务的完成条件是动态的
	ExpireTime int `json:"expiretime"` ///任务过期时间utc
}

//	taskTypeTeamStarColor              = 10 ///球队拥有指定位置上数量达到品质要求的球员
const (
	taskTypeNpcTeamFormWin             = 1  ///对指定球队用指定阵型获胜
	taskTypeNpcTeamTacticWin           = 2  ///对指定球队用指定战术获胜
	taskTypeNpcTeamTotalGoalCount      = 3  ///对指定球队进球总数达到指定数量
	taskTypeNpcTeamGoalDifference      = 4  ///对指定球队达到指定数量净胜球
	taskTypeNpcTeamSingleGoalCount     = 5  ///对指定球队单场进球达到指定数量
	taskTypeNpcTeamSingleLostGoalCount = 6  ///对指定球队单场失球不超过指定数量,往少了记
	taskTypeNpcTeamTotalWin            = 7  ///对指定球队获胜达到指定数量
	taskTypeTeamStarScore              = 8  ///球队拥有指定位置上数量达到评分要求的球员
	taskTypeTeamStarStarCount          = 9  ///球队拥有指定位置上数量达到星级要求的球员
	taskTypeDayNpcTeamTotalGoalCount   = 10 //日常任务,对指定球队进球
	taskTypeDayNpcTeamGoalDifference   = 11 //日常任务,对指定球队进净胜球
	taskTypeDayFunctionDoCount         = 12 //日常任务,系统功能执行次数
	taskTypeDayVipShopBuyItem          = 13 //日常任务,促销购买
	taskTypeDayTenDrawStars            = 14 //日常任务,十连抽球星
)

const (
	dayTaskFunctionTrainPass        = 1  ///日常任务功能条件-传球训练
	dayTaskFunctionTrainSteals      = 2  ///日常任务功能条件-抢断训练
	dayTaskFunctionTrainDribbling   = 3  ///日常任务功能条件-盘断训练
	dayTaskFunctionTrainSliding     = 4  ///日常任务功能条件-铲球训练
	dayTaskFunctionTrainShooting    = 5  ///日常任务功能条件-射门训练
	dayTaskFunctionTrainGoalKeeping = 6  ///日常任务功能条件-守门训练
	dayTaskFunctionTrainBody        = 7  ///日常任务功能条件-身体训练
	dayTaskFunctionTrainSpeed       = 8  ///日常任务功能条件-传球训练
	dayTaskFunctionTrainRefresh     = 9  ///日常任务功能条件-训练赛刷新
	dayTaskFunctionArenaMatch       = 10 ///日常任务功能条件-天天联赛
	dayTaskFunctionBuyActionPoint   = 11 ///日常任务功能条件-购买体力
	dayTaskFunctionBuyCoin          = 12 ///日常任务功能条件-使用点金手
	dayTaskFunctionStarEducation    = 13 ///日常任务功能条件-球员培养
	dayTaskFunctionStarTrain        = 14 ///日常任务功能条件-球员训练
	dayTaskFunctionStarSpyLow       = 15 ///日常任务功能条件-初级球探
	dayTaskFunctionStarSpyMid       = 16 ///日常任务功能条件-中级球探
	dayTaskFunctionStarSpyHigh      = 17 ///日常任务功能条件-高级球探
	dayTaskFunctionTenDrawStars     = 18 ///日常任务功能条件-十连抽
)

type TaskTypeStaticData struct { ///关卡类型,对应动态表中的dy_leveltype
	ID     int ///任务类型id
	Type   int ///任务类型
	League int ///任务关联的联赛类型
	Part1  int ///参数1
	Part2  int ///参数1
	Part3  int ///参数1
	Item1  int ///奖励道具类型1
	Num1   int ///奖励道具数量1
	Item2  int ///奖励道具类型2
	Num2   int ///奖励道具数量2
}

func (self *TaskTypeStaticData) IsNpcTeamTask() bool {
	return self.Type >= taskTypeNpcTeamFormWin && self.Type <= taskTypeNpcTeamTotalWin
}

func (self *TaskTypeStaticData) IsDayTask() bool {
	isDayTask := self.Type >= taskTypeDayNpcTeamTotalGoalCount && self.Type <= taskTypeDayTenDrawStars
	return isDayTask
}

func (self *TaskTypeStaticData) IsDayTaskFunction() bool {
	isDayTaskFunction := (taskTypeDayFunctionDoCount == self.Type || taskTypeDayTenDrawStars == self.Type)
	return isDayTaskFunction
}

func (self *TaskTypeStaticData) IsDayTaskVipBuy() bool {
	isDayTaskVipBuy := (taskTypeDayVipShopBuyItem == self.Type)
	return isDayTaskVipBuy
}

type Task struct { ///任务对象
	TaskInfo ///任务信息
	DataUpdater
}

func (self *Task) GetInfoPtr() *TaskInfo {
	return &self.TaskInfo
}

func (self *Task) GetID() int {
	return self.ID
}

///判断此任务是否已过期
func (self *Task) IsExpire() bool {
	//	if 0 == self.ExpireTime {
	//		return false ///永不过期任务
	//	}
	now := Now()
	isExpire := now >= self.ExpireTime
	return isExpire
}

func NewTask(taskInfo *TaskInfo) *Task {
	task := new(Task)
	task.TaskInfo = *taskInfo
	task.InitDataUpdater(tableTask, &task.TaskInfo)
	return task
}

func (self *Task) GetTypeInfo() *TaskTypeStaticData { ///取得任务类型静态数据信息
	taskTypeStaticData := GetServer().GetStaticDataMgr().Unsafe().GetTaskType(self.Type)
	return taskTypeStaticData
}

func (self *Task) UpdateDone(team *Team) bool { ///有修改返回true
	if self.IsDone > 0 {
		return false
	}
	isTaskDone := false
	oldProgress := self.Progress ///保存旧值
	taskTypeStaticData := self.GetTypeInfo()
	needFieldType := taskTypeStaticData.Part1 ///得到需求场类型
	if taskTypeStaticData.Part1 >= fieldTypeAll {
		needFieldType = fieldTypeNone ///指定全场类型
	}
	switch taskTypeStaticData.Type {
	//	case taskTypeTeamStarColor: ///球队拥有指定位置上数量达到品质要求的球员
	//		self.Progress = team.FindStarCount(needFieldType, taskTypeStaticData.Part3, 0, 0)
	//		isTaskDone = (self.Progress >= taskTypeStaticData.Part2)
	case taskTypeTeamStarStarCount: ///球队拥有指定位置上数量达到星级要求的球员
		self.Progress = team.FindStarCount(needFieldType, 0, taskTypeStaticData.Part3, 0)
		isTaskDone = (self.Progress >= taskTypeStaticData.Part2)
	case taskTypeTeamStarScore: ///球队拥有指定位置上数量达到评分要求的球员
		self.Progress = team.FindStarCount(needFieldType, 0, 0, taskTypeStaticData.Part3)
		isTaskDone = (self.Progress >= taskTypeStaticData.Part2)
	}
	if true == isTaskDone { ///判断任务是否完成
		self.IsDone = 1 ///任务完成
	}
	return isTaskDone || self.Progress != oldProgress ///返回外界任务是否已更新了
}

func (self *Task) UpdateDayTaskVipBuy(team *Team, vipGoodsType int) bool { ///更新日常任务功能相关
	if self.Need1 != vipGoodsType {
		return false ///非购买任务指定商品
	}
	//self.Save()     ///先保存一次
	self.Progress++ ///增加一次执行次数
	if self.Progress >= self.Need2 {
		self.IsDone = 1 ///更新任务完成状态
	}
	return true
}

func (self *Task) UpdateDayTaskFuntion(team *Team, typeDoFunction int) bool { ///更新日常任务功能相关
	taskType := self.GetTypeInfo()
	if taskType.Part3 != typeDoFunction {
		return false
	}
	self.Progress++ ///增加一次执行次数
	if self.Progress >= self.Need1 {
		self.IsDone = 1 ///更新任务完成状态
	}
	return true
}

func (self *Task) UpdateNpcTeamTask(team *Team, userGoalCount int, npcGoalCount int) bool { ///有修改返回true
	if userGoalCount <= npcGoalCount {
		return false ///只有战胜npc球队时才更新此npc球队相关的任务信息
	}
	isTaskDone := false
	oldProgress := self.Progress ///保存旧值
	staticDataMgrUnsafe := GetServer().GetStaticDataMgr().Unsafe()
	taskTypeStaticData := staticDataMgrUnsafe.GetTaskType(self.Type)
	formation := team.GetCurrentFormObject() ///得到当前阵形
	formationInfo := formation.GetInfo()
	switch taskTypeStaticData.Type {
	case taskTypeNpcTeamFormWin: ///对指定球队用指定阵型获胜
		isTaskDone = formationInfo.Type == taskTypeStaticData.Part2
	case taskTypeNpcTeamTacticWin: ///对指定球队用指定战术获胜
		isTaskDone = (formationInfo.CurrentTactic == taskTypeStaticData.Part2)
	case taskTypeNpcTeamTotalGoalCount: ///对指定球队进球总数达到指定数量
		self.Progress += userGoalCount ///累计进球数
		isTaskDone = (self.Progress >= taskTypeStaticData.Part2)
	case taskTypeNpcTeamGoalDifference: ///对指定球队达到指定数量净胜球
		self.Progress = Max(self.Progress, userGoalCount-npcGoalCount) ///更新净胜球数
		isTaskDone = (self.Progress >= taskTypeStaticData.Part2)
	case taskTypeNpcTeamSingleGoalCount: ///对指定球队单场进球达到指定数量
		self.Progress = Max(self.Progress, userGoalCount)
		isTaskDone = (self.Progress >= taskTypeStaticData.Part2)
	case taskTypeNpcTeamSingleLostGoalCount: ///对指定球队单场失球不超过指定数量
		//if userGoalCount > npcGoalCount {
		if self.Progress <= 0 {
			self.Progress = npcGoalCount
		}
		self.Progress = Min(self.Progress, npcGoalCount)
		//}
		isTaskDone = (self.Progress <= taskTypeStaticData.Part2)
	case taskTypeNpcTeamTotalWin: ///对指定球队获胜达到指定数量
		//if userGoalCount > npcGoalCount {
		self.Progress++
		//}
		isTaskDone = (self.Progress >= taskTypeStaticData.Part2)
	case taskTypeDayNpcTeamTotalGoalCount:
		self.Progress += userGoalCount ///累计进球数
		isTaskDone = (self.Progress >= self.Need2)
	case taskTypeDayNpcTeamGoalDifference: ///更新累加净胜球数
		self.Progress += Max(0, userGoalCount-npcGoalCount)
		isTaskDone = (self.Progress >= self.Need2)
	}
	if true == isTaskDone { ///判断任务是否完成
		self.IsDone = 1 ///任务完成
	}
	return isTaskDone || self.Progress != oldProgress ///返回外界任务是否已更新了
}

func (self *Task) IsComplete() bool { ///判断任务是否已完成
	return self.IsDone > 0
}

func (self *Task) GetReflectValue() reflect.Value { ///得到球队反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

func (self *Task) Reset() {
	self.Progress = 0
	self.IsDone = 0
	self.IsAward = 0
	self.Need1 = 0
	self.Need2 = 0
	self.ExpireTime = 0
	//	self.Save()
}

func (self *Task) GetDayTaskNextRefreshTime() int { ///得到天任务下次更新时间utc
	currentHour := time.Now().Hour()
	paramIntList := GetServer().GetStaticDataMgr().getConfigStaticDataParamIntList(configTask,
		configTaskDayTaskRefresh)
	nextRefreshTime := paramIntList[0] ///初始化为第一个值
	for i := range paramIntList {
		refreshHour := paramIntList[i]
		if refreshHour <= 0 {
			continue
		}
		if refreshHour > currentHour {
			nextRefreshTime = refreshHour
			break
		}
	}
	nextRefreshTime = GetHourUTC(nextRefreshTime)
	return nextRefreshTime
}

func (self *Task) RefreshDayTaskNpcTeamTotalGoalCount(team *Team, hasPassNpcTeamList IntList) {
	listLen := hasPassNpcTeamList.Len()
	rndNpcTeamIndex := Random(0, listLen-1) ///随机球队索引
	npcTeamType := 101011
	if listLen > 0 {
		npcTeamType = hasPassNpcTeamList[rndNpcTeamIndex]
	}
	taskType := self.GetTypeInfo()
	rndGoalCount := Random(taskType.Part1, taskType.Part2)
	self.Need1 = npcTeamType  ///要求打的npc球队类型
	self.Need2 = rndGoalCount ///要求总进球数
}

func (self *Task) RefreshDayTaskNpcTeamGoalDifference(team *Team, hasPassNpcTeamList IntList) {
	listLen := hasPassNpcTeamList.Len()
	rndNpcTeamIndex := Random(0, listLen-1) ///随机球队索引
	npcTeamType := 101011
	if listLen > 0 {
		npcTeamType = hasPassNpcTeamList[rndNpcTeamIndex]
	}
	taskType := self.GetTypeInfo()
	rndGoalCount := Random(taskType.Part1, taskType.Part2)
	self.Need1 = npcTeamType  ///要求打的npc球队类型
	self.Need2 = rndGoalCount ///要求净剩球数
}

func (self *Task) RefreshDayTaskFunctionDoCount(team *Team) {
	taskMgr := team.GetTaskMgr()
	taskType := self.GetTypeInfo()
	taskTypeIDList := taskMgr.GetTaskTypeIDList(taskType.Type)
	taskCount := taskTypeIDList.Len()
	rndTaskIndex := Random(0, taskCount-1)
	rndTaskID := taskTypeIDList[rndTaskIndex]
	self.Type = rndTaskID
	rndDoCount := Random(taskType.Part1, taskType.Part2)
	self.Need1 = rndDoCount ///要求完成功能次数
}

func (self *Task) RefreshDayTaskVipShopBuyItem(team *Team) {
	dayTaskVipBuyGroup := 2
	vipShopStaticDataList := GetServer().GetStaticDataMgr().GetVipShopStaticDataList(dayTaskVipBuyGroup)
	goodCount := vipShopStaticDataList.Len()
	rndGoodIndex := Random(0, goodCount-1)
	rndGoodID := vipShopStaticDataList[rndGoodIndex]
	self.Need1 = rndGoodID ///要求完成购买商品中种功能次数
	self.Need2 = 1         ///要求完成购买行为的次数
}

func (self *Task) RefreshDayTenDrawStars(team *Team) {
	self.Need1 = 1 ///执行一次十连抽
}

func (self *Task) RefreshDayTask(team *Team) {
	levelMgr := team.GetLevelMgr()
	hasPassNpcTeamList := levelMgr.GetHasPassNpcTeamList()
	dayTaskNextRefreshTime := self.GetDayTaskNextRefreshTime()
	if team.FunctionMask&(1<<(functionMaskDayTask-1)) == 0 { ///没开启日常任务功能时,任务不变
		self.ExpireTime = 0 ///更新过期时间
		return
	}
	self.Reset()                             ///重置任务信息
	self.ExpireTime = dayTaskNextRefreshTime ///更新过期时间

	taskType := self.GetTypeInfo()
	switch taskType.Type {
	case taskTypeDayNpcTeamTotalGoalCount:
		self.RefreshDayTaskNpcTeamTotalGoalCount(team, hasPassNpcTeamList)
	case taskTypeDayNpcTeamGoalDifference:
		self.RefreshDayTaskNpcTeamGoalDifference(team, hasPassNpcTeamList)
	case taskTypeDayFunctionDoCount:
		self.RefreshDayTaskFunctionDoCount(team)
	case taskTypeDayVipShopBuyItem:
		self.RefreshDayTaskVipShopBuyItem(team)
	case taskTypeDayTenDrawStars:
		self.RefreshDayTenDrawStars(team)
	}
}
