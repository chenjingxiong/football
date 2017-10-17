package football

import (
	"fmt"
)

///球队任务管理器
//type ITaskMgr interface {
//	IGameMgr
//	GetTaskInfoList(leagueType int) TaskInfoList                                                 ///得到任务信息列表
//	IsTaskTypeDone(taskType int) bool                                                            ///判断指定类型任务是否完成
//	FindTask(taskType int) ITask                                                                 ///得到任务信息
//	TakeTaskAward(team *Team, taskType int) IntList                                              ///领取任务奖励
//	UpdateNpcTeamTask(client IClient, npcTeamType int, userGoalCount int, npcGoalCount int) bool ///更新NpcTeam相关任务信息
//	//AccpetTask(taskType int) ITask                                               ///领取此任务
//	GetTask(taskID int) ITask ///得到任务对象
//	AddTask(taskType int) int ///领取此任务
//}

type TaskInfoList []TaskInfo
type TaskList map[int]*Task

type TaskMgr struct { ///任务管理器
	GameMgr
	taskList TaskList
}

//func (self *TaskMgr) Save() { ///保存数据
//	for _, v := range self.taskList {
//		v.Save()
//	}
//}

func (self *TaskMgr) GetType() int { ///得到管理器类型
	return mgrTypeTaskMgr ///任务管理器
}

func (self *TaskMgr) SaveInfo() { ///保存数据
	for _, v := range self.taskList {
		v.Save()
	}
}

func (self *TaskMgr) IsTaskTypeDone(taskType int) bool { ///判断指定类型任务是否完成
	task := self.FindTask(taskType)
	if nil == task {
		return false ///不存在的任务对象为false
	}
	if task.IsComplete() == false {
		task.UpdateDone(self.team)
	}
	taskInfoPtr := task.GetInfoPtr()
	return taskInfoPtr.IsDone > 0
}

///更新日常任务功能相关
func (self *TaskMgr) UpdateDayTaskVipBuy(client *Client, vipGoodsType int) {
	syncMgr := client.GetSyncMgr()
	team := client.GetTeam()
	syncTaskList := SyncObjectList{}
	dayTaskList := self.GetDayTaskList()
	for i := range dayTaskList {
		dayTaskID := dayTaskList[i]
		dayTask := self.GetTask(dayTaskID)
		dayTaskType := dayTask.GetTypeInfo()
		if nil == dayTask {
			continue ///无效的任务忽略
		}
		if dayTaskType.IsDayTaskVipBuy() == false {
			continue ///非日常功能商城购买任务跳过
		}
		if dayTask.IsExpire() == true {
			continue ///过期任务跳过
		}
		if dayTask.IsComplete() == true {
			continue ///已完成任务跳过
		}
		isChange := dayTask.UpdateDayTaskVipBuy(team, vipGoodsType)
		if true == isChange {
			syncTaskList = append(syncTaskList, dayTask)
		}
	}
	syncMgr.SyncObjectArray("UpdateDayTaskVipBuy", syncTaskList) ///更新任务信息
}

func (self *TaskMgr) UpdateDayTaskFuntion(client *Client, typeDoFunction int) { ///更新日常任务功能相关
	syncMgr := client.GetSyncMgr()
	team := client.GetTeam()
	syncTaskList := SyncObjectList{}
	dayTaskList := self.GetDayTaskList()
	for i := range dayTaskList {
		dayTaskID := dayTaskList[i]
		dayTask := self.GetTask(dayTaskID)
		dayTaskType := dayTask.GetTypeInfo()
		if nil == dayTask {
			continue ///无效的任务忽略
		}
		if dayTaskType.IsDayTaskFunction() == false {
			continue ///非日常功能任务跳过
		}
		if dayTask.IsExpire() == true {
			continue ///过期任务跳过
		}
		if dayTask.IsComplete() == true {
			continue ///已完成任务跳过
		}
		isChange := dayTask.UpdateDayTaskFuntion(team, typeDoFunction)
		if true == isChange {
			syncTaskList = append(syncTaskList, dayTask)
		}
	}
	syncMgr.SyncObjectArray("UpdateDayTaskFuntion", syncTaskList) ///更新任务信息
}

func (self *TaskMgr) UpdateDayTask(client IClient, npcTeamType int, userGoalCount int, npcGoalCount int) bool {
	syncMgr := client.GetSyncMgr()
	team := client.GetTeam()
	syncTaskList := SyncObjectList{}
	dayTaskList := self.GetDayTaskList()
	for i := range dayTaskList {
		dayTaskID := dayTaskList[i]
		daytask := self.GetTask(dayTaskID)
		if nil == daytask {
			continue ///无效的任务忽略
		}
		if daytask.IsExpire() == true {
			continue
		}
		if daytask.IsComplete() == true {
			continue
		}
		if daytask.Need1 == npcTeamType {
			if daytask.UpdateNpcTeamTask(team, userGoalCount, npcGoalCount) == true {
				syncTaskList = append(syncTaskList, daytask)
			}
		}
	}
	syncMgr.SyncObjectArray("UpdateDayTask", syncTaskList) ///更新任务信息
	return true
}

///更新NpcTeam相关任务信息
func (self *TaskMgr) UpdateNpcTeamTask(client IClient, npcTeamType int, userGoalCount int, npcGoalCount int) bool {
	staticDataMgrUnsafe := GetServer().GetStaticDataMgr().Unsafe()
	syncMgr := client.GetSyncMgr()
	team := client.GetTeam()
	taskTypeList := staticDataMgrUnsafe.GetTaskTypeList()
	addTaskIDList := IntList{}
	syncTaskList := SyncObjectList{}
	for i := range taskTypeList {
		taskType := taskTypeList[i]
		if taskType.IsNpcTeamTask() == false {
			continue ///非npc球队任务忽略
		}
		if taskType.Part1 != npcTeamType {
			continue ///非指定npc球队类型任务忽略
		}
		task := self.FindTask(taskType.ID)
		if nil == task { ///天任务不得自动创建任务
			taskID := self.AddTask(taskType.ID)
			task = self.GetTask(taskID)
			//syncMgr.syncAddTask(task)
			addTaskIDList = append(addTaskIDList, taskID)
		}
		if nil == task {
			continue ///无效的任务忽略
		}
		if task.IsComplete() == true {
			continue
		}
		if task.UpdateNpcTeamTask(team, userGoalCount, npcGoalCount) == true {
			//syncMgr.SyncObject("UpdateNpcTeamTask", task) ///更新任务信息
			syncTaskList = append(syncTaskList, task)
		}
	}
	syncMgr.SyncAddTask(addTaskIDList)
	syncMgr.SyncObjectArray("UpdateNpcTeamTask", syncTaskList)           ///更新任务信息
	self.UpdateDayTask(client, npcTeamType, userGoalCount, npcGoalCount) ///更新日常任务
	return true
}

func (self *TaskMgr) GetTask(taskID int) *Task { ///得到任务对象
	return self.taskList[taskID]
}

func (self *TaskMgr) GetTaskTypeIDList(taskType int) IntList { ///得到任务信息
	taskTypeResultList := IntList{}
	taskTypeList := GetServer().GetStaticDataMgr().GetTaskTypeList()
	for i := range taskTypeList {
		taskTypeInfo := taskTypeList[i]
		if taskTypeInfo.Type == taskType {
			taskTypeResultList = append(taskTypeResultList, taskTypeInfo.ID)
		}
	}
	return taskTypeResultList
}

func (self *TaskMgr) FindTask(taskType int) *Task { ///得到任务信息
	for _, v := range self.taskList {
		taskInfoPtr := v.GetInfoPtr()
		if taskInfoPtr.Type == taskType {
			return v
		}
	}
	return nil
}

func (self *TaskMgr) AddTask(taskType int) int { ///领取此任务
	if taskType <= 0 {
		return 0
	}
	accpetTaskQuery := fmt.Sprintf("insert %s set teamid=%d,type=%d", tableTask, self.team.GetID(), taskType) ///组插入记录SQL
	lastInsertTaskID, _ := GetServer().GetDynamicDB().Exec(accpetTaskQuery)
	if lastInsertTaskID <= 0 {
		GetServer().GetLoger().Warn("TaskMgr AccpetTask fail! taskType:%d teamid:%d", taskType, self.team.GetID())
		return 0
	}
	taskInfo := new(TaskInfo)
	taskInfo.ID = lastInsertTaskID
	taskInfo.TeamID = self.team.GetID()
	taskInfo.Type = taskType
	self.taskList[lastInsertTaskID] = NewTask(taskInfo) ///生成关卡对象
	return lastInsertTaskID
}

func (self *TaskMgr) TakeTaskAward(team *Team, taskType int) IntList { ///领取任务奖励,得到道具对象id
	task := self.FindTask(taskType)
	if nil == task {
		return nil ///无此任务
	}
	taskInfoPtr := task.GetInfoPtr()
	if taskInfoPtr.IsDone <= 0 {
		return nil ///未完成不能领奖
	}
	if taskInfoPtr.IsAward > 0 {
		return nil ///已领奖
	}
	staticDataMgrUnsafe := GetServer().GetStaticDataMgr().Unsafe()
	taskTypeStaticData := staticDataMgrUnsafe.GetTaskType(taskType)
	if taskTypeStaticData.Item1 <= 0 || taskTypeStaticData.Num1 <= 0 {
		return nil ///无效的奖品配置数据
	}
	itemIDList := team.GetItemMgr().AwardItem(taskTypeStaticData.Item1, taskTypeStaticData.Num1)
	if itemIDList != nil {
		///领奖成功更新打上已领奖标识
		taskInfoPtr.IsAward = 1
	}
	return itemIDList
}

func (self *TaskMgr) GetTaskInfoList(leagueType int) TaskInfoList { ///得到任务信息列表
	taskInfoList := TaskInfoList{}
	staticDataMgrUnsafe := GetServer().GetStaticDataMgr().Unsafe()
	for _, v := range self.taskList {
		taskInfoPtr := v.GetInfoPtr()
		taskType := staticDataMgrUnsafe.GetTaskType(taskInfoPtr.Type)
		if nil == taskType {
			continue
		}
		if taskType.League == leagueType {
			taskInfoList = append(taskInfoList, *taskInfoPtr)
		}
	}
	return taskInfoList
}

func NewTaskMgr(teamID int) IGameMgr {
	taskMgr := new(TaskMgr)
	taskMgr.taskList = make(TaskList)
	//	taskMgr.teamID = teamID ///存放自己的球队id
	taskListQuery := fmt.Sprintf("select * from %s where teamid=%d limit 3000", tableTask, teamID)
	taskInfo := new(TaskInfo)
	taskInfoList := GetServer().GetDynamicDB().fetchAllRows(taskListQuery, taskInfo)
	for i := range taskInfoList {
		taskInfo = taskInfoList[i].(*TaskInfo)
		taskMgr.taskList[taskInfo.ID] = NewTask(taskInfo)
	}
	return taskMgr
}

///得到日常任务列表
func (self *TaskMgr) GetDayTaskList() IntList {
	dayTaskList := IntList{}
	for _, v := range self.taskList {
		taskType := v.GetTypeInfo()
		if taskType.IsDayTask() == true {
			dayTaskList = append(dayTaskList, v.ID)
		}
	}
	return dayTaskList
}

///清空日常任务列表
func (self *TaskMgr) CleanDayTaskList() {
	for _, v := range self.taskList {
		taskType := v.GetTypeInfo()
		if taskType.IsDayTask() == true {
			delete(self.taskList, v.ID)
		}
	}
}

///创建球队所有日常任务
func (self *TaskMgr) CreateAllDayTask() {
	dataParamIntList := GetServer().GetStaticDataMgr().getConfigStaticDataParamIntList(configTask, configTaskDayTaskInit)
	for i := range dataParamIntList {
		taskType := dataParamIntList[i]
		if taskType <= 0 {
			continue
		}
		self.AddTask(taskType) ///添加总进球数日常任务
	}
}

///刷新球队日常任务
func (self *TaskMgr) RefreshDayTask() {
	dayTaskList := self.GetDayTaskList()
	if dayTaskList.Len() <= 0 {
		self.CreateAllDayTask() ///如果没有日常任务则新创建几个
		dayTaskList = self.GetDayTaskList()
	}
	for i := range dayTaskList {
		taskID := dayTaskList[i]
		task := self.GetTask(taskID)
		if task.IsExpire() == false {
			break ///没过期任务不处理
		}
		task.RefreshDayTask(self.team) ///刷新自己
	}
}

func (self *TaskMgr) GetDayTaskInfoList() TaskInfoList { ///得到日常任务信息列表
	taskInfoList := TaskInfoList{}
	for _, v := range self.taskList {
		taskType := v.GetTypeInfo()
		if nil == taskType {
			continue
		}
		if taskType.IsDayTask() == true {
			taskInfoList = append(taskInfoList, v.TaskInfo)
		}
	}
	return taskInfoList
}

func (self *TaskMgr) onInit() {
	//	self.RefreshDayTask()
}
