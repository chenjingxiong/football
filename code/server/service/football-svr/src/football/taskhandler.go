package football

import (
// "fmt"
//"time"
)

type TaskQueryDayTaskResultMsg struct { ///查询天任务列表结果消息
	MsgHead             `json:"head"` ///"daytask", "querydaytask"
	HasAcceptAwardCount int           `json:"hasacceptawardcount"` ///已领取奖励的次数
	DayTaskList         TaskInfoList  `json:"tasklist"`            ///天任务列表
}

func (self *TaskQueryDayTaskResultMsg) GetTypeAndAction() (string, string) {
	return "daytask", "querydaytaskresult"
}

type TaskQueryDayTaskMsg struct { ///查询天任务列表
	MsgHead `json:"head"` ///"task", "querydaytask"
}

func (self *TaskQueryDayTaskMsg) GetTypeAndAction() (string, string) {
	return "daytask", "querydaytask"
}

func SendTaskQueryDayTaskResultMsg(client IClient) { ///发送查询天任务列表结果包	team := client.GetTeam()
	team := client.GetTeam()
	taskMgr := team.GetTaskMgr()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.QueryResetAttrib(ResetAttribTypeDayTask)
	taskQueryDayTaskResultMsg := new(TaskQueryDayTaskResultMsg)
	taskQueryDayTaskResultMsg.HasAcceptAwardCount = resetAttrib.Value1
	taskQueryDayTaskResultMsg.DayTaskList = taskMgr.GetDayTaskInfoList()
	client.SendMsg(taskQueryDayTaskResultMsg)
}

func (self *TaskQueryDayTaskMsg) processAction(client IClient) (result bool) { ///查询天任务列表
	team := client.GetTeam()
	taskMgr := team.GetTaskMgr()
	resetAttribMgr := team.GetResetAttribMgr()
	taskMgr.RefreshDayTask() ///刷新日常任务
	resetAttrib := resetAttribMgr.QueryResetAttrib(ResetAttribTypeDayTask)
	resetAttrib.ResetDayTask()
	SendTaskQueryDayTaskResultMsg(client)
	return true
}

type TaskAccpetDayTaskAwardResultMsg struct { ///领取天任务奖励结果消息
	MsgHead `json:"head"` ///"daytask", "accpetdaytaskawardresult"
	TaskID  int           `json:"taskid"` ///领取奖励的任务id
	Result  int           `json:"result"` ///0失败,1成功
}

func (self *TaskAccpetDayTaskAwardResultMsg) GetTypeAndAction() (string, string) {
	return "daytask", "querydaytaskresult"
}

type TaskAccpetDayTaskAwardMsg struct { ///领取天任务奖励
	MsgHead `json:"head"` ///"daytask", "accpetdaytaskaward"
	TaskID  int           `json:"taskid"` ///希望领取奖励的任务id
}

func (self *TaskAccpetDayTaskAwardMsg) GetTypeAndAction() (string, string) {
	return "daytask", "accpetdaytaskaward"
}

func (self *TaskAccpetDayTaskAwardMsg) checkAction(client IClient) (result bool) { ///领取天任务奖励
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	taskMgr := team.GetTaskMgr()
	resetAttribMgr := team.GetResetAttribMgr()
	task := taskMgr.GetTask(self.TaskID)
	if loger.CheckFail("task!=nil", task != nil, task, nil) {
		return false ///无效的任务
	}
	taskType := task.GetTypeInfo()
	if loger.CheckFail("taskType!=nil", taskType != nil, taskType, nil) {
		return false ///无效的任务类型
	}
	isDayTask := taskType.IsDayTask()
	if loger.CheckFail("taskType!=nil", isDayTask == true, isDayTask, true) {
		return false ///不是日常任务
	}
	if loger.CheckFail("task.IsDone>0", task.IsDone > 0, task.IsDone, 0) {
		return false ///日常任务未完成
	}
	if loger.CheckFail("task.IsAward<=0", task.IsAward <= 0, task.IsAward, 0) {
		return false ///日常任务已领过奖励了
	}
	isExpire := task.IsExpire()
	if loger.CheckFail(" isExpire==false", isExpire == false, isExpire, false) {
		return false ///日常任务已过期不能领取奖励了
	}
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeDayTask)
	if loger.CheckFail("resetAttrib!=nil", resetAttrib != nil, resetAttrib, nil) {
		return false ///日常任务信息生成与获得失败
	}
	acceptDayTaskAwardMaxCount := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTask,
		configTaskDayTaskRefresh, 6) ///取得配置表中日常任务每日最大领取奖励数
	if loger.CheckFail("resetAttrib!=nil", resetAttrib.Value1 < acceptDayTaskAwardMaxCount, resetAttrib.Value1,
		acceptDayTaskAwardMaxCount) {
		return false ///日常任务信息生成与获得失败
	}
	return true
}

func (self *TaskAccpetDayTaskAwardMsg) payAction(client IClient) (result bool) { ///领取天任务奖励
	team := client.GetTeam()
	taskMgr := team.GetTaskMgr()
	resetAttribMgr := team.GetResetAttribMgr()
	task := taskMgr.GetTask(self.TaskID)
	task.IsAward = 1 ///打上已领取标识
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeDayTask)
	resetAttrib.Value1++ ///已领取奖励次数加1
	return true
}

func (self *TaskAccpetDayTaskAwardMsg) doAction(client IClient) (result bool) { ///领取天任务奖励
	team := client.GetTeam()
	taskMgr := team.GetTaskMgr()
	task := taskMgr.GetTask(self.TaskID)
	taskType := task.GetTypeInfo()
	if taskType.Item1 > 0 && taskType.Num1 > 0 {
		team.AwardObject(taskType.Item1, taskType.Num1, 0, 0)
	}
	if taskType.Item2 > 0 && taskType.Num2 > 0 {
		team.AwardObject(taskType.Item2, taskType.Num2, 0, 0)
	}
	SendTaskQueryDayTaskResultMsg(client)
	return true
}

func (self *TaskAccpetDayTaskAwardMsg) processAction(client IClient) (result bool) { ///领取天任务奖励
	if self.checkAction(client) == false { ///检测
		return false
	}
	if self.payAction(client) == false { ///支付
		return false
	}
	if self.doAction(client) == false { ///发货
		return false
	}
	return true
}

type TaskTenDrawStarsMsg struct { ///请求球星十连抽
	MsgHead `json:"head"` ///"daytask", "tendrawstars"
}

func (self *TaskTenDrawStarsMsg) GetTypeAndAction() (string, string) {
	return "daytask", "tendrawstars"
}

func (self *TaskTenDrawStarsMsg) checkAction(client IClient) (result bool) {
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	starCenter := team.GetStarCenter()
	taskMgr := team.GetTaskMgr()
	staticDataMgr := GetServer().GetStaticDataMgr()
	taskTypeIDList := taskMgr.GetTaskTypeIDList(taskTypeDayTenDrawStars)
	taskTypeIDListLen := taskTypeIDList.Len()
	if loger.CheckFail("taskTypeIDListLen>0", taskTypeIDListLen > 0, taskTypeIDListLen, 0) {
		return false ///没有找到十连抽相关的日常任务
	}
	taskTypeID := taskTypeIDList[0] ///得到任务类型id
	taskType := staticDataMgr.GetTaskType(taskTypeID)
	needPayDiamonds := taskType.Part1 ///得到需要支付的钻石数
	currentDiamonds := team.GetTicket()
	if loger.CheckFail("needPayDiamonds<=currentDiamonds",
		needPayDiamonds <= currentDiamonds, needPayDiamonds, currentDiamonds) {
		return false ///余额不足
	}
	remainTenDrawStarCount := starCenter.GetRemainTenDrawStarCount()
	if loger.CheckFail("remainTenDrawStarCount<=0",
		remainTenDrawStarCount <= 0, remainTenDrawStarCount, 0) {
		return false ///还有未取得的球员存在,不能刷新,否则玩家利益损失
	}
	return true
}

func (self *TaskTenDrawStarsMsg) payAction(client IClient) (result bool) {
	team := client.GetTeam()
	taskMgr := team.GetTaskMgr()
	syncMgr := client.GetSyncMgr()
	staticDataMgr := GetServer().GetStaticDataMgr()
	taskTypeIDList := taskMgr.GetTaskTypeIDList(taskTypeDayTenDrawStars)
	taskTypeID := taskTypeIDList[0] ///得到任务类型id
	taskType := staticDataMgr.GetTaskType(taskTypeID)
	needPayDiamonds := taskType.Part1 ///得到需要支付的钻石数
	team.PayTicket(needPayDiamonds)   ///支付钻石
	syncMgr.SyncObject("BuyActionPointMsg", team)
	return true
}

func SendTaskQueryTenDrawStarsResultMsg(client IClient, resultType int) {
	team := client.GetTeam()
	starCenter := team.GetStarCenter()
	starCenterMemberInfoList := starCenter.GetStarCenterMemberInfoList(starCenterTypeStarTenDraw)
	taskQueryTenDrawStarsResultMsg := new(TaskQueryTenDrawStarsResultMsg)
	taskQueryTenDrawStarsResultMsg.ResultType = resultType
	taskQueryTenDrawStarsResultMsg.StarList = starCenterMemberInfoList
	//fmt.Println(starCenterMemberInfoList)
	client.SendMsg(taskQueryTenDrawStarsResultMsg)
}

func (self *TaskTenDrawStarsMsg) doAction(client IClient) (result bool) {
	team := client.GetTeam()
	taskMgr := team.GetTaskMgr()
	starCenter := team.GetStarCenter()
	starCenter.UpdateStarTenDraw()                ///更新球星十连抽数据
	SendTaskQueryTenDrawStarsResultMsg(client, 2) ///同步客户端十连抽状态
	///更新球星十连抽的日常任务
	taskMgr.UpdateDayTaskFuntion(client.GetElement(), dayTaskFunctionTenDrawStars)
	//! 更新十连抽开服活动
	team.OSActivityValue.Refresh(team, 1)
	return true
}

func (self *TaskTenDrawStarsMsg) processAction(client IClient) (result bool) {
	if self.checkAction(client) == false { ///检测
		return false
	}
	if self.payAction(client) == false { ///支付
		return false
	}
	if self.doAction(client) == false { ///发货
		return false
	}
	return true
}

//type TaskTenDrawStarsResultMsg struct { ///请求球星十连抽
//	MsgHead  `json:"head"`            ///"task", "tendrawstarsresult"
//	StarList StarCenterMemberInfoList `json:"starlist"` ///十连抽球员类型列表
//}

//func (self *TaskTenDrawStarsResultMsg) GetTypeAndAction() (string, string) {
//	return "task", "tendrawstarsresult"
//}

type TaskQueryTenDrawStarsMsg struct { ///请求查询十连抽已有球星
	MsgHead `json:"head"` ///"daytask", "querytendrawstars"
}

func (self *TaskQueryTenDrawStarsMsg) GetTypeAndAction() (string, string) {
	return "daytask", "querytendrawstars"
}

func (self *TaskQueryTenDrawStarsMsg) processAction(client IClient) bool {
	SendTaskQueryTenDrawStarsResultMsg(client, 1)
	return true
}

type TaskQueryTenDrawStarsResultMsg struct { ///请求查询十连抽已有球星
	MsgHead    `json:"head"`            ///"daytask", "querytendrawstarsresult"
	ResultType int                      `json:"resulttype"` ///结果类型 1表示查询结果 2表示十连抽结果 3表示翻卡后结果
	StarList   StarCenterMemberInfoList `json:"starlist"`   ///十连抽球员类型列表
}

func (self *TaskQueryTenDrawStarsResultMsg) GetTypeAndAction() (string, string) {
	return "daytask", "querytendrawstarsresult"
}

const (
	TakeTenDrawOperateTypeAwardStar  = 1 ///1.招至球队/提升本队相同球员星级
	TakeTenDrawOperateTypeConverStar = 2 ///2.转化经验
	TakeTenDrawOperateTypeGivePiece  = 3 ///3.给予碎片
)

type TaskTakeTenDrawStars struct { ///请求取得十连抽产生的球员
	MsgHead     `json:"head"` ///"daytask" "taketendrawstars"
	MemberID    int           `json:"memberid"`    ///请求取得的成员id,此成员存放球员类型
	OperateType int           `json:"operatetype"` ///发掘后操作类型: 1.招至球队/提升本队相同球员星级  2.转化经验
}

func (self *TaskTakeTenDrawStars) GetTypeAndAction() (string, string) {
	return "daytask", "taketendrawstars"
}

func (self *TaskTakeTenDrawStars) checkAction(client IClient) (result bool) {
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	starCenter := team.GetStarCenter()
	//fmt.Println(self.MemberID)
	starCenterMember := starCenter.GetStarCenterMember(starCenterTypeStarTenDraw, self.MemberID)
	if loger.CheckFail(" starCenterMember!=nil", starCenterMember != nil, starCenterMember, nil) {
		return false ///无效的球员
	}
	if loger.CheckFail("starCenterMember.ExpireTime<=0", starCenterMember.ExpireTime <= 0,
		starCenterMember.ExpireTime, 0) {
		return false ///球员已经被抽取过了
	}
	if loger.CheckFail("self.OperateType >=1 && self.OperateType <=3", self.OperateType >= 1 && self.OperateType <= 3,
		self.OperateType, 12) {
		return false ///无效的操作类型
	}
	return true
}

func (self *TaskTakeTenDrawStars) payAction(client IClient) (result bool) {
	team := client.GetTeam()
	starCenter := team.GetStarCenter()
	starCenterMember := starCenter.GetStarCenterMember(starCenterTypeStarTenDraw, self.MemberID)
	takeFlag := Now()
	if TakeTenDrawOperateTypeConverStar == self.OperateType || TakeTenDrawOperateTypeGivePiece == self.OperateType {
		takeFlag *= -1
	}
	starCenterMember.ExpireTime = takeFlag ///已抽选则意味马上过期
	return true
}

func (self *TaskTakeTenDrawStars) doAction(client IClient) (result bool) {
	team := client.GetTeam()
	starCenter := team.GetStarCenter()
	starCenterMember := starCenter.GetStarCenterMember(starCenterTypeStarTenDraw, self.MemberID)

	// if starCenterMember.StarType >= 110000 { //若抽取为道具
	// 	//!直接给予道具
	// 	team.AwardObject(starCenterMember.StarType, starCenterMember.EvolveCount, 0, 0)
	// 	SendTaskQueryTenDrawStarsResultMsg(client, 3)
	// 	return true
	// }

	starType := GetServer().GetStaticDataMgr().Unsafe().GetStarType(starCenterMember.StarType)
	switch self.OperateType {
	case TakeTenDrawOperateTypeAwardStar:
		if starCenterMember.StarType >= 110000 {
			return true
		}
		team.AwardObject(0, 0, starCenterMember.EvolveCount, starCenterMember.StarType)
	case TakeTenDrawOperateTypeConverStar:
		if starCenterMember.StarType >= 110000 {
			return true
		}
		starCardCount := team.GetStarCardCount(starType, starCenterMember.EvolveCount)
		team.AwardObject(ItemStarCard, starCardCount, 0, 0) ///赠星卡
	case TakeTenDrawOperateTypeGivePiece:
		team.AwardObject(starCenterMember.StarType, starCenterMember.EvolveCount, 0, 0)
	}
	SendTaskQueryTenDrawStarsResultMsg(client, 3)
	return true
}

func (self *TaskTakeTenDrawStars) processAction(client IClient) (result bool) {
	if self.checkAction(client) == false { ///检测
		return false
	}
	if self.payAction(client) == false { ///支付
		return false
	}
	if self.doAction(client) == false { ///发货
		return false
	}
	return true
}
