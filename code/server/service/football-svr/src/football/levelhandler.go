package football

///关卡消息处理器
import (
	"math/rand"
)

const (
	levelFunctionStarEducation    = 1  ///球员培养
	levelFunctionItemEquip        = 2  ///装备开启
	levelFunctionAddTrainCell     = 3  ///训练功能
	levelFunctionAddFormation     = 4  ///送阵形
	levelFunctionStarSpy          = 5  ///开启球员挖掘
	levelFunctionSkill            = 6  ///开启技能系统并奖励一个技能
	levelFunctionVolunteer        = 7  ///开启球星来投
	levelFunctionStarConvince     = 8  ///开启球员游说
	levelFunctionStarEvolve       = 9  ///开启球员突破
	levelFunctionChangeStar       = 10 ///更换球员
	levelFunctionFindStar1        = 11 ///初级寻找球员
	levelFunctionFindStar2        = 12 ///中级寻找球员
	levelFunctionFindStar3        = 13 ///高级寻找球员
	levelFunctionTrainMatch       = 14 ///训练赛
	levelFunctionFormationUplevel = 15 ///阵型升级
	levelFunctionArenaMatch       = 16 ///天天联赛
	levelFunctionItemMerge        = 17 ///装备融合
	levelFunctionStarTrain        = 18 ///球员训练
	levelFunctionDayTask          = 19 ///日常任务功能开启
	levelFunctionChallangeMatch   = 20 ///挑战赛
	levelFunctionStarSkill        = 21 //!球员技能
	levelFunctionMannaStar        = 22 //!天赐球员
)

func (self *LevelHandler) getName() string { ///返回可处理的消息类型
	return "level"
}

func (self *LevelHandler) initHandler() { ///初始化消息处理器
	self.addActionToList(new(QueryLeagueInfoMsg))
	self.addActionToList(new(QueryLevelListMsg))
	self.addActionToList(new(PassLevelMsg))
	self.addActionToList(new(SkipLevelMsg))
}

type LevelHandler struct {
	MsgHandler
}

type QueryLeagueInfoResultMsg struct { ///查询球队所拥有联赛消息,所有地图信息
	MsgHead    `json:"head"`  ///"level", "queryleagueinforesult"
	LeagueList LeagueInfoList `json:"leaguelist"`
}

func (self *QueryLeagueInfoResultMsg) GetTypeAndAction() (string, string) {
	return "level", "queryleagueinforesult"
}

type QueryLeagueInfoMsg struct { ///查询球队所拥有联赛消息,所有地图信息
	MsgHead `json:"head"` ///"level", "queryleagueinfo"
}

func (self *QueryLeagueInfoMsg) GetTypeAndAction() (string, string) {
	return "level", "queryleagueinfo"
}

func (self *QueryLeagueInfoMsg) processAction(client IClient) bool {
	levelMgr := client.GetTeam().GetLevelMgr()
	queryLeagueInfoResultMsg := new(QueryLeagueInfoResultMsg)
	queryLeagueInfoResultMsg.LeagueList = levelMgr.GetLeagueInfoList()
	client.SendMsg(queryLeagueInfoResultMsg)
	return true
}

type QueryLevelListResultMsg struct { ///球队所拥有指定联赛所有关卡消息,指定地图信息,查询结果
	MsgHead         `json:"head"` ///"level", "querylevellistresult"
	LeagueType      int           `json:"leaguetype"`      ///请求指定联赛类型
	LevelList       LevelInfoList `json:"levelinfolist"`   ///关卡信息列表,可能为空
	TaskList        TaskInfoList  `json:"taskinfolist"`    ///关卡相关任务列表,可能为空
	LevelMatchCDUTC int           `json:"levelmatchcdutc"` ///关卡比赛cd的utc时间秒,0表示没有cd
	RemainSkipCount int           `json:"remainskipcount"` ///关卡比赛剩余跳过次数
}

func (self *QueryLevelListResultMsg) GetTypeAndAction() (string, string) {
	return "level", "querylevellistresult"
}

type QueryLevelListMsg struct { ///查询球队所拥有指定联赛所有关卡消息,指定地图信息
	MsgHead    `json:"head"` ///"level", "querylevellist"
	LeagueType int           `json:"leaguetype"` ///请求指定联赛类型
}

func (self *QueryLevelListMsg) GetTypeAndAction() (string, string) {
	return "level", "querylevellist"
}

func (self *QueryLevelListMsg) processAction(client IClient) bool { ///发送联赛中详细关卡与关卡任务信息给客户端
	levelMgr := client.GetTeam().GetLevelMgr()
	taskMgr := client.GetTeam().GetTaskMgr()
	processMgr := client.GetTeam().GetProcessMgr()
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeSkipLevelCount)
	if resetAttrib == nil {
		//matchLevelSkipCountMax := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTeam, configTeamLevelMatchCommon, 2) ///取得可跳过次数
		//matchLevelSkipCountResetTime := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTeam,
		//	configTeamLevelMatchCommon, 3) ///取得可跳过次数重置时间
		//resetTimeUTC := GetHourUTC(matchLevelSkipCountResetTime)
		resetAttrib = resetAttribMgr.AddResetAttrib(ResetAttribTypeSkipLevelCount,
			0, IntList{0})
		//resetAttrib = resetAttribMgr.GetResetAttrib(resetAttribID)
		resetAttrib.ResetMatchLevelSkipCount(nil)
	}
	queryLevelListResultMsg := new(QueryLevelListResultMsg)
	queryLevelListResultMsg.LeagueType = self.LeagueType
	queryLevelListResultMsg.LevelList = levelMgr.GetLevelInfoList(self.LeagueType)          ///得到关卡信息列表
	queryLevelListResultMsg.TaskList = taskMgr.GetTaskInfoList(self.LeagueType)             ///得到任务信息列表
	queryLevelListResultMsg.LevelMatchCDUTC = processMgr.GetTeamCDUTC(teamLevelMatchCDType) ///得到球队比赛cd时间
	queryLevelListResultMsg.RemainSkipCount = resetAttrib.Value1                            ///同步剩余跳过次数
	client.SendMsg(queryLevelListResultMsg)
	return true
}

type SkipLevelMsgResult struct { ///客户端请求跳过比赛结果
	MsgHead         `json:"head"` ///"level", "skiplevelresult"
	RemainSkipCount int           `json:"remainskipcount"` ///剩余跳过战斗次数
}

func (self *SkipLevelMsgResult) GetTypeAndAction() (string, string) {
	return "level", "skiplevelresult"
}

type SkipLevelMsg struct { ///客户端请求跳过比赛
	MsgHead `json:"head"` ///"level", "skiplevel"
}

func (self *SkipLevelMsg) GetTypeAndAction() (string, string) {
	return "level", "skiplevel"
}

func (self *SkipLevelMsg) checkAction(client IClient) bool { ///检测条件
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeSkipLevelCount)
	if loger.CheckFail("resetAttrib!=nil", resetAttrib != nil, resetAttrib, nil) {
		return false ///可跳过次数必须有效
	}
	remainSkipCount := resetAttrib.Value1
	if loger.CheckFail("remainSkipCount>0", remainSkipCount > 0, remainSkipCount, 0) {
		return false ///可跳过次数必须有效
	}
	return true
}

func (self *SkipLevelMsg) processAction(client IClient) bool { ///客户端请求通过指定关卡处理流程
	if self.checkAction(client) == false { ///跳关时的通用检测
		return false
	}
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeSkipLevelCount)
	resetAttrib.Value1-- ///跳过次数减１
	skipLevelMsgResult := new(SkipLevelMsgResult)
	skipLevelMsgResult.RemainSkipCount = resetAttrib.Value1
	client.SendMsg(skipLevelMsgResult)
	return true
}

//func (self *PassLevelMsg) checkAction(client IClient) bool { ///检测条件
//	return true
//}

//func (self *PassLevelMsg) payAction(client IClient) bool { ///付款

//	return true
//}

//func (self *PassLevelMsg) doAction(client IClient) bool { ///发货

//	return true
//}

type PassLevelMsg struct { ///客户端请求通过指定关卡
	MsgHead   `json:"head"` ///"level", "passlevel"
	LevelSort int           `json:"levelsort"` ///请求通关的关卡子类型
	LevelType int           `json:"leveltype"` ///请求通关的关卡类型
	LevelID   int           `json:"levelid"`   ///请求通关的关卡id
	Param1    int           `json:"param1"`    ///通关参数1,根据不同的关卡类型,param1的内容也不同
}

func (self *PassLevelMsg) GetTypeAndAction() (string, string) {
	return "level", "passlevel"
}

func (self *PassLevelMsg) processActionCardCheck(client IClient) bool { ///处理过关时的卡类关卡检测
	team := client.GetTeam()
	levelType := GetServer().GetStaticDataMgr().Unsafe().GetLevelType(self.LevelType)
	if nil == levelType {
		return false ///无效的关卡类型
	}
	currentTicket := team.GetTicket() ///得到当前球票
	if self.Param1 > 0 {
		if levelType.Sid2 <= 0 || levelType.Sid3 <= 0 {
			return false ///无效的特殊奖励配置数据
		}
		if currentTicket < levelType.Sid3 {
			return false ///余额不足
		}
	}
	return true
}

func (self *PassLevelMsg) processActionSkillCard(client IClient) bool { ///处理技能卡请求
	if self.processActionCardCheck(client) == false {
		return false
	}
	levelMgr := client.GetTeam().GetLevelMgr()
	team := client.GetTeam()
	//	skillMgr := team.GetSkillMgr()
	syncMgr := client.GetSyncMgr()
	levelType := GetServer().GetStaticDataMgr().Unsafe().GetLevelType(self.LevelType)
	awardSkillType := levelType.Sid1
	if self.Param1 > 0 {
		awardSkillType = levelType.Sid2
	}
	if self.Param1 > 0 { ///拿特殊奖励
		team.PayTicket(levelType.Sid3) ///扣特殊奖励所需球票
		client.SetMoneyRecord(PlayerCostMoney, Pay_SpecialAward, levelType.Sid3, team.GetTicket())
		syncMgr.SyncObject("PassLevelMsg", team)
	}
	//	skillID := skillMgr.AwardSkill(awardSkillType) ///生成新技能
	//	syncMgr.syncAddSkill(IntList{skillID}) ///同步新拿到的道具信息列表
	GetServer().GetLoger().CYDebug("%d", awardSkillType)
	levelID := levelMgr.AddLevel(self.LevelType, 0)
	client.GetSyncMgr().syncAddLevel(IntList{levelID})
	return true
}

func (self *PassLevelMsg) processActionItemCard(client IClient) bool { ///处理道具卡请求
	if self.processActionCardCheck(client) == false {
		return false
	}
	levelMgr := client.GetTeam().GetLevelMgr()
	syncMgr := client.GetSyncMgr()
	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	levelType := GetServer().GetStaticDataMgr().Unsafe().GetLevelType(self.LevelType)
	awardItemType := levelType.Sid1
	if self.Param1 > 0 {
		awardItemType = levelType.Sid2
	}
	if awardItemType <= 0 {
		return false ///无效的奖励道具类型
	}
	if self.Param1 > 0 { ///拿特殊奖励
		team.PayTicket(levelType.Sid3) ///扣特殊奖励所需球票
		client.SetMoneyRecord(PlayerCostMoney, Pay_SpecialAward, levelType.Sid3, team.GetTicket())
		syncMgr.SyncObject("PassLevelMsg", team)
	}
	itemIDList := itemMgr.AwardItem(awardItemType, 1) ///生成新道具
	if itemIDList != nil {
		for i := range itemIDList {
			item := itemMgr.GetItem(itemIDList[i])
			itemInfo := item.GetInfo()
			if self.Param1 > 0 {
				itemInfo.Color = levelType.Sid5
			} else {
				itemInfo.Color = levelType.Sid4
			}
		}

	}

	syncMgr.syncAddItem(itemIDList) ///同步新拿到的道具信息列表
	levelID := levelMgr.AddLevel(self.LevelType, 0)
	client.GetSyncMgr().syncAddLevel(IntList{levelID})
	return true
}

func (self *PassLevelMsg) processActionStarCard(client IClient) bool { ///处理球员卡请求
	levelMgr := client.GetTeam().GetLevelMgr()
	syncMgr := client.GetSyncMgr()
	team := client.GetTeam()
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	levelType := GetServer().GetStaticDataMgr().Unsafe().GetLevelType(self.LevelType)
	currentTicket := team.GetTicket() ///得到当前球票
	awardStarType := levelType.Sid1
	awardStarLevel := levelType.Sid4
	if self.Param1 > 0 {
		awardStarType = levelType.Sid2
		awardStarLevel = levelType.Sid5
		if levelType.Sid2 <= 0 || levelType.Sid3 <= 0 {
			return false ///无效的特殊奖励配置数据
		}
		if currentTicket < levelType.Sid3 {
			return false ///余额不足
		}
	}
	if awardStarType <= 0 {
		return false ///无效的球员类型
	}
	if staticDataMgr.GetStarType(awardStarType) == nil {
		return false ///在数据库静态表中找不到此球员类型
	}

	if self.Param1 > 0 { ///拿特殊奖励
		team.PayTicket(levelType.Sid3) ///扣特殊奖励所需球票
		client.SetMoneyRecord(PlayerCostMoney, Pay_SpecialAward, levelType.Sid3, team.GetTicket())
		syncMgr.SyncObject("PassLevelMsg", team)
	}
	//	starID := team.AwardStar(awardStarType) ///给特殊球员
	//	syncMgr.syncAddStar(IntList{starID})    ///同步新拿到的球员信息列表

	//	isHasStar := true
	star := team.GetStarFromType(awardStarType)
	//	if star == nil {
	//		isHasStar = false ///检查球星是否存在玩家队伍,根据存在状态,选择是否改变该球星属性
	//	}
	team.AwardObject(0, 0, awardStarLevel, awardStarType)

	//if isHasStar == false { ///球员奖励赠送的球员属性修改为当前星级的上一星级满属性（满等级，满培养）的状态  && self.Param1 > 0
	star = team.GetStarFromType(awardStarType)
	star.ChangeStarAttribute(awardStarLevel - 1)
	syncMgr.SyncObject("PassLevelMsg", star)
	//}

	levelID := levelMgr.AddLevel(self.LevelType, 0)
	client.GetSyncMgr().syncAddLevel(IntList{levelID})
	return true
}

func (self *PassLevelMsg) processActionLocker(client IClient) bool { ///处理关卡锁请求
	team := client.GetTeam()
	levelMgr := client.GetTeam().GetLevelMgr()
	levelType := GetServer().GetStaticDataMgr().Unsafe().GetLevelType(self.LevelType)
	taskMgr := client.GetTeam().GetTaskMgr()
	taskTypeList := IntList{levelType.Sid1, levelType.Sid2, levelType.Sid3}
	for i := range taskTypeList {
		taskType := taskTypeList[i]
		if taskType <= 0 {
			continue ///忽略
		}
		task := taskMgr.FindTask(taskType)
		if nil == task {
			taskID := taskMgr.AddTask(taskType)
			task = taskMgr.GetTask(taskID)
		}
		if nil == task {
			continue ///再次失败忽略
		}
		if taskMgr.IsTaskTypeDone(taskType) == false {
			return false
		}
	}
	addItemIDList := IntList{}
	for k := range taskTypeList {
		taskType := taskTypeList[k]
		if taskType <= 0 {
			continue
		}
		itemIDList := taskMgr.TakeTaskAward(team, taskType)
		if itemIDList != nil {
			addItemIDList = append(addItemIDList, itemIDList...)
		}
	}
	if len(addItemIDList) > 0 {
		client.GetSyncMgr().syncAddItem(addItemIDList)
	}
	levelID := levelMgr.AddLevel(self.LevelType, 0)
	client.GetSyncMgr().syncAddLevel(IntList{levelID})
	return true

	//if levelType.Sid1 > 0 && taskMgr.IsTaskTypeDone(levelType.Sid1) == false {
	//	return false ///任务1未完成
	//}
	//if levelType.Sid2 > 0 && taskMgr.IsTaskTypeDone(levelType.Sid2) == false {
	//	return false ///任务2未完成
	//}
	//if levelType.Sid3 > 0 && taskMgr.IsTaskTypeDone(levelType.Sid3) == false {
	//	return false ///任务3未完成
	//}

	//if levelType.Sid1 > 0 { ///任务1发奖
	//	itemID := taskMgr.TakeTaskAward(team, levelType.Sid1)
	//	if itemID > 0 {
	//		addItemIDList = append(addItemIDList, itemID)
	//	}
	//}
	//if levelType.Sid2 > 0 { ///任务2发奖
	//	itemID := taskMgr.TakeTaskAward(team, levelType.Sid2)
	//	if itemID > 0 {
	//		addItemIDList = append(addItemIDList, itemID)
	//	}
	//}
	//if levelType.Sid3 > 0 { ///任务3发奖
	//	itemID := taskMgr.TakeTaskAward(team, levelType.Sid3)
	//	if itemID > 0 {
	//		addItemIDList = append(addItemIDList, itemID)
	//	}
	//}

}

func (self *PassLevelMsg) processActionStartFinish(client IClient) bool { ///处理起点与终点请求
	levelMgr := client.GetTeam().GetLevelMgr()
	levelID := levelMgr.AddLevel(self.LevelType, 0)
	client.GetSyncMgr().syncAddLevel(IntList{levelID})
	return levelID > 0
}

func (self *PassLevelMsg) processActionFunction(client IClient) bool { ///处理游戏功能开启卡
	levelTypeStaticData := GetServer().GetStaticDataMgr().Unsafe().GetLevelType(self.LevelType)
	team := client.GetTeam()
	//	loger := GetServer().GetLoger()
	itemMgr := team.GetItemMgr()
	syncMgr := client.GetSyncMgr()
	//	skillMgr := team.GetSkillMgr()
	formationMgr := team.GetFormationMgr()
	switch levelTypeStaticData.Sid1 {
	case levelFunctionStarEducation: ///球员培养开启
		team.SetFunctionMask(functionMaskStarEducation, 1) ///开启指定功能
		team.AwardObject(levelTypeStaticData.Sid2, levelTypeStaticData.Sid3, 0, 0)
	case levelFunctionAddTrainCell: ///开训练系统并训练位加1
		team.SetFunctionMask(functionMaskStarTrain, 1)
		team.AwardExpPool(levelTypeStaticData.Sid3)
		syncMgr.SyncObject("levelFunctionStarEvolve", team)
		//team.AwardTrainCell(levelTypeStaticData.Sid2) ///加训练位
	case levelFunctionAddFormation: ///送新阵形
		formationID := formationMgr.AwardFormation(levelTypeStaticData.Sid2)
		syncMgr.syncAddFormation(IntList{formationID}) ///同步新加的阵型到客户端
	case levelFunctionItemEquip: ///道具装备开启
		team.SetFunctionMask(functionMaskItemEquip, 1)               ///开启指定功能
		itemIDList := itemMgr.AwardItem(levelTypeStaticData.Sid2, 1) ///奖励道具
		if itemIDList != nil {
			syncMgr.syncAddItem(itemIDList)
		}

	//case levelFunctionStarSpy: ///球员发掘开启
	//	team.SetFunctionMask(functionMaskStarSpy, 1) ///开启指定功能
	//case levelFunctionSkill: ///技能系统开启
	//	team.SetFunctionMask(functionMaskSkill, 1) ///开启指定功能
	//	skillID := skillMgr.AwardSkill(levelTypeStaticData.Sid2)
	//	syncMgr.syncAddSkill(IntList{skillID})
	case levelFunctionVolunteer: ///球星来投开启
		team.SetFunctionMask(functionMaskVolunteer, 1)               ///开启指定功能
		itemIDList := itemMgr.AwardItem(levelTypeStaticData.Sid2, 1) ///奖励道具
		if itemIDList != nil {
			syncMgr.syncAddItem(itemIDList)
		}
	case levelFunctionStarConvince: ///球星游说开启
		team.SetFunctionMask(functionMaskStarConvince, 1) ///开启指定功能
	case levelFunctionStarEvolve: ///球员进化开启
		team.SetFunctionMask(functionMaskStarEvolve, 1) ///开启指定功能
		team.AwardObject(levelTypeStaticData.Sid2, levelTypeStaticData.Sid3, 0, 0)
	case levelFunctionChangeStar: ///更换球员
		//staticDataMgr := GetServer().GetStaticDataMgr()
		//starType := staticDataMgr.GetStarType(levelTypeStaticData.Sid2)
		//if starType == nil {
		//	loger.Warn("processActionFunction starType == nil teamID: %d", team.GetID()) ///不合法球员
		//	return false
		//}

		//isHas := team.HasStar(levelTypeStaticData.Sid2)
		//if isHas == true {
		//	loger.Warn("processActionFunction Has this star teamID: %d starType: %d", team.GetID(), levelTypeStaticData.Sid2) ///队伍已有此球员
		//	return false
		//}
		//		team.SetFunctionMask(functionMaskChangeStar, 1)     ///开启指定功能
		//		team.AwardObject(0, 0, 0, levelTypeStaticData.Sid2) ///给予球星
	case levelFunctionFindStar1: ///初级寻找球员
		team.SetFunctionMask(functionMaskFindStar1, 1)  ///开启指定功能
		team.SetFunctionMask(functionMaskChangeStar, 1) ///开启指定功能
		team.SetFunctionMask(functionMaskFindStar2, 1)  ///开启指定功能
		team.SetFunctionMask(functionMaskFindStar3, 1)  ///开启指定功能  开启初级寻找球员同时开启中级高级
		//		team.SetFunctionMask(functionMaskChallangeMatch, 1) ///for debug
	case levelFunctionFindStar2: ///中级寻找球员
		team.SetFunctionMask(functionMaskFindStar2, 1) ///开启指定功能
	case levelFunctionFindStar3: ///高级寻找球员
		team.SetFunctionMask(functionMaskFindStar3, 1) ///开启指定功能
	case levelFunctionTrainMatch: ///训练赛
		team.SetFunctionMask(functionMaskTrainMatch, 1) ///开启指定功能

		//训练赛开启时重置新手引导刷新
		resetAttribMgr := team.GetResetAttribMgr()
		resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamTrainMatch)
		if nil != resetAttrib {
			resetAttrib.Value5 = 0 ///重置新手引导
			RefreshTrainMatchType(resetAttrib)
		}

	case levelFunctionFormationUplevel: ///阵型升级
		team.SetFunctionMask(functionMaskFormationUplevel, 1) ///开启指定功能
		team.AwardObject(levelTypeStaticData.Sid2, levelTypeStaticData.Sid3, 0, 0)
	case levelFunctionArenaMatch: ///天天联赛
		team.SetFunctionMask(functionMaskArenaMatch, 1) ///开启指定功能
	case levelFunctionItemMerge: ///装备融合
		team.SetFunctionMask(functionMaskItemMerge, 1) ///开启指定功能
	case levelFunctionStarTrain: ///球员训练
		team.SetFunctionMask(functionMaskTrainStar, 1) ///开启指定功能
	case levelFunctionDayTask: ///日常任务
		team.SetFunctionMask(functionMaskDayTask, 1) ///开启指定功能

		//! 开启功能时主动推送数据一次
		taskMgr := team.GetTaskMgr()
		resetAttribMgr := team.GetResetAttribMgr()
		taskMgr.RefreshDayTask() ///刷新日常任务
		resetAttrib := resetAttribMgr.QueryResetAttrib(ResetAttribTypeDayTask)
		resetAttrib.ResetDayTask()
		SendTaskQueryDayTaskResultMsg(client)

	case levelFunctionChallangeMatch: ///挑战赛
		team.SetFunctionMask(functionMaskChallangeMatch, 1) ///开启挑战赛功能

	case levelFunctionStarSkill: //! 技能
		team.SetFunctionMask(functionMaskStarSkill, 1)

	case levelFunctionMannaStar: //! 天赐球员
		team.SetFunctionMask(functionMaskMannaStar, 1)
	}
	client.GetSyncMgr().SyncObject("PassLevelMsg", team) ///同步最新的功能掩码给客户端
	levelMgr := client.GetTeam().GetLevelMgr()
	levelID := levelMgr.AddLevel(self.LevelType, 0)
	client.GetSyncMgr().syncAddLevel(IntList{levelID})
	return levelID > 0
}

//func (self *PassLevelMsg) sendPassLevelMatchResultMsg(client IClient, levelID int, usergoalCount int, npcGoalCount int) { ///处理比赛请求
//	passLevelMatchResultMsg := new(PassLevelMatchResultMsg)
//	passLevelMatchResultMsg.LevelID = levelID
//	passLevelMatchResultMsg.LevelType = self.LevelType
//	passLevelMatchResultMsg.NpcTeamType = self.Param1
//	passLevelMatchResultMsg.HomeGoal = usergoalCount
//	passLevelMatchResultMsg.GuestGoal = npcGoalCount
//	client.SendMsg(passLevelMatchResultMsg)
//}

//func (self *PassLevelMsg) processActionMatchPay(client IClient) bool { ///处理比赛请求扣代价
//	//	team := client.GetTeam()

//	return true
//}

//func (self *PassLevelMsg) processActionMatchCheck(client IClient) bool { ///处理比赛请求检测

//	return true
//}

func (self *PassLevelMsg) calcMatchStarCount(homeGoalCount int, guestGoalCount int) int { ///通过比分计算星级
	matchStarCount := 0 ///平局或战败为0星
	diffGoalCount := homeGoalCount - guestGoalCount
	if diffGoalCount >= 3 {
		matchStarCount = 3 //净胜球三球以上且不丢球3星
	} else if diffGoalCount >= 2 {
		matchStarCount = 2 //净胜球(进球数减去失球数)2球以上2星；
	} else if diffGoalCount >= 1 {
		matchStarCount = 1 ///战胜球队1星；
	}
	return matchStarCount
}

func (self *PassLevelMsg) getNpcTeamBaseStarCount(npcTeamType int) int { ///得到比赛的npcteam基础星级
	npcTeamBaseStarCount := 0
	levelTypeStaticData := GetServer().GetStaticDataMgr().Unsafe().GetLevelType(self.LevelType)
	if levelTypeStaticData.Sid1 == npcTeamType {
		npcTeamBaseStarCount = 0 ///普通级别0星
	} else if levelTypeStaticData.Sid2 == npcTeamType {
		npcTeamBaseStarCount = 3 ///经典级别3星
	} else if levelTypeStaticData.Sid3 == npcTeamType {
		npcTeamBaseStarCount = 6 ///经典级别6星
	}
	return npcTeamBaseStarCount
}

func (self *PassLevelMsg) processActionAwardItem(client IClient, finalstarcount int, starcount int) IntList { ///处理奖励道具
	itemMgr := client.GetTeam().GetItemMgr()
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	syncMgr := client.GetSyncMgr()
	npcTeamType := self.Param1 ///得到准备攻击的球队id
	staticData := staticDataMgr.GetNpcTeamType(npcTeamType)
	awardItemIDList, awardItemCountList := IntList{}, IntList{}
	///规则修改:2星掉落addteanmexp。 3星掉落additemtype1。其跟随字段（additemrate1，additemmincount1，additemmaxcount1）功能不变。
	//if staticData.AddItemType1 > 0 && finalstarcount >= 2 && starcount < 2 { ///奖励道具1
	if finalstarcount >= 2 && starcount < 2 { ///奖励道具1
		// randRate1 := rand.Intn(100) ///随机获得几率
		// ///随机获得个数, 可以获得0个道具,如果个数为0则跳过
		// randItemCount1 := Random(staticData.AddItemMinCount1, staticData.AddItemMaxCount1)
		// if randRate1 <= staticData.AddItemRate1 && randItemCount1 > 0 {
		// 	itemIDList := itemMgr.AwardItem(staticData.AddItemType1, randItemCount1)
		// 	if itemIDList != nil {
		// 		awardItemIDList = append(awardItemIDList, itemIDList...)
		// 	}

		// 	awardItemCountList = append(awardItemCountList, randItemCount1)
		// 	awardItemCountList = append(awardItemCountList, 0)
		// }

		team := client.GetTeam()
		team.AwardExpPool(staticData.AddTeamExp)
		syncMgr.SyncObject("PassLevelMsg", team)
	}
	if staticData.AddItemType1 > 0 && finalstarcount >= 3 && starcount < 3 { ///奖励道具2

		randRate1 := rand.Intn(100) ///随机获得几率
		///随机获得个数, 可以获得0个道具,如果个数为0则跳过
		randItemCount1 := Random(staticData.AddItemMinCount1, staticData.AddItemMaxCount1)
		if randRate1 <= staticData.AddItemRate1 && randItemCount1 > 0 {
			itemIDList := itemMgr.AwardItem(staticData.AddItemType1, randItemCount1)
			if itemIDList != nil {
				awardItemIDList = append(awardItemIDList, itemIDList...)
			}

			awardItemCountList = append(awardItemCountList, randItemCount1)
			awardItemCountList = append(awardItemCountList, 0)
		}

		// randRate2 := rand.Intn(100) ///随机获得几率
		// ///随机获得个数, 可以获得0个道具,如果个数为0则跳过
		// randItemCount2 := Random(staticData.AddItemMinCount2, staticData.AddItemMaxCount2)
		// if randRate2 <= staticData.AddItemRate2 && randItemCount2 > 0 {
		// 	itemIDList := itemMgr.AwardItem(staticData.AddItemType2, randItemCount2)
		// 	if itemIDList != nil {
		// 		awardItemIDList = append(awardItemIDList, itemIDList...)
		// 	}
		// 	awardItemCountList = append(awardItemCountList, 0)
		// 	awardItemCountList = append(awardItemCountList, randItemCount2)
		// }
	}
	if len(awardItemIDList) > 0 {
		syncMgr.syncAddItem(awardItemIDList) ///同步客户端已获得道具信息列表
	}
	return awardItemCountList
}

func (self *PassLevelMsg) processActionMatch(client IClient) bool { ///处理比赛请求发货
	loger := GetServer().GetLoger()
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	team := client.GetTeam()
	levelMgr := client.GetTeam().GetLevelMgr()
	syncMgr := client.GetSyncMgr()
	taskMgr := team.GetTaskMgr()
	//	processMgr := team.GetProcessMgr()
	awardCoin, awardExp, talentPoint := 0, 0, 0
	awardItemCountList, awardStarIDList := IntList{}, IntList{}
	//matchCD := 0
	///检测比赛cd是否合法
	//isExpireTeamCD := processMgr.isExpireTeamCD(teamLevelMatchCDType)
	//if loger.CheckFail(" isExpireTeamCD==true", isExpireTeamCD == true, isExpireTeamCD, true) {
	//	return false ///比赛cd有cd且未到期
	//}
	///检查客户端请求攻击球队的合法性
	levelTypeStaticData := staticDataMgr.GetLevelType(self.LevelType) ///关卡类型对象
	if levelTypeStaticData.Sid1 != self.Param1 && levelTypeStaticData.Sid2 != self.Param1 &&
		levelTypeStaticData.Sid3 != self.Param1 {
		return false ///非法的npc球队类型
	}
	npcTeamType := self.Param1 ///得到准备攻击的球队id
	npcTeamTypeStaticData := staticDataMgr.GetNpcTeamType(npcTeamType)
	if nil == npcTeamTypeStaticData {
		return false ///无效的npc球队类型,配置数据不存在
	}
	if npcTeamTypeStaticData.AddStarType > 0 && team.IsStarFull() == true {
		loger.Warn("npcTeamTypeStaticData.AddStarType > 0 && team.IsStarFull() == true")
		return false ///球队已满同时有奖励玩家入队就不让通关
	}
	isStarsContractPointEnough := team.IsStarsContractPointEnough()
	if loger.CheckFail("isStarsContractPointEnough==true", isStarsContractPointEnough == true, isStarsContractPointEnough, true) {
		return false ///所有上场球员的契约值均要求大于0
	}

	level := levelMgr.FindLevel(self.LevelType) ///移到前面，因为三星掉落需要用

	///判断是否是三星掉落比赛
	isTriStarMatch := false             ///是否是3星掉落比赛的标记
	matchType := npcTeamType / 100000   ///通过npc球队id取模的方式取得比赛的类型
	attribType := 0                     ///可重置数据的类型
	if matchType == 1 && level != nil { ///比赛类型必须是冠军之路而且必须有关卡对象才可能是三星掉落比赛
		teamIndex := levelTypeStaticData.GetNpcIndexInLevel(npcTeamType) ///球队在关卡中的位置
		if level.StarCount >= teamIndex*3 {                              ///判断星数是否满足要求
			isTriStarMatch = true
		} else {
			isTriStarMatch = false
		}
		///判断挑战次数
		attribType = GetTriStarAttribType(npcTeamType) ///取得可重置数据的类型
		resetAttribMgr := team.GetResetAttribMgr()     ///取得队伍的可重置管理器
		resetAttrib := resetAttribMgr.GetResetAttrib(attribType)
		if resetAttrib != nil && resetAttrib.Value1 > 0 { ///这个值可能是没有的这里判断一下，剩余挑战
			isTriStarMatch = true
		} else { ///如果挑战次数不够是要放过的，只是最后不给奖励，所以这里要设false
			isTriStarMatch = false
		}
	}
	///体力限制
	if isTriStarMatch { ///三星挑战需要消费体力
		if team.PayActionPoint(npcTeamTypeStaticData.AddItemRate2) { ///扣除体力成功后就要减去挑战次数
			resetAttribMgr := team.GetResetAttribMgr() ///取得队伍的可重置管理器
			resetAttrib := resetAttribMgr.GetResetAttrib(attribType)
			resetAttrib.Value1--
			resetAttrib.Save()
		} else {
			return false ///活力不足
		}
	}

	userGoalCount, npcGoalCount := npcTeamTypeStaticData.CalcMatchResult(team) ///计算比赛结果
	matchStarCount := self.calcMatchStarCount(userGoalCount, npcGoalCount)     ///得到比赛结果星级

	///是否已通过一次
	///根据结果给奖励
	levelID := 0

	if nil == level { ///还未打过个关卡
		levelID = levelMgr.AddLevel(self.LevelType, 0)
		level = levelMgr.GetLevel(levelID)
		client.GetSyncMgr().syncAddLevel(IntList{levelID})
	}

	levelInfoPtr := level.GetInfoPtr()
	npcTeamBaseStarCount := self.getNpcTeamBaseStarCount(npcTeamType) ///得到npcteam基础星级
	finalMatchStarCount := npcTeamBaseStarCount                       ///得到最终比赛产生的星级
	if userGoalCount > npcGoalCount && matchStarCount > (levelInfoPtr.StarCount-npcTeamBaseStarCount) {
		finalMatchStarCount = levelInfoPtr.StarCount + 1
		finalMatchStarCount = Min(npcTeamBaseStarCount+3, finalMatchStarCount)
	}

	levelID = levelInfoPtr.ID
	if levelInfoPtr.StarCount < finalMatchStarCount { ///有经验拿的情况
		if levelInfoPtr.StarCount == npcTeamBaseStarCount &&
			npcTeamTypeStaticData.AddStarType > 0 { ///第一次过此队有球星奖励
			starID := team.AwardStar(npcTeamTypeStaticData.AddStarType)
			awardStarIDList = append(awardStarIDList, starID)
			syncMgr.syncAddStar(awardStarIDList)
		}
		///计算奖励经验值
		awardTeamExp := (finalMatchStarCount - levelInfoPtr.StarCount) * npcTeamTypeStaticData.AddBaseExp
		team.AwardExp(client, awardTeamExp) ///给球队奖经验
		awardExp = awardTeamExp

		if (finalMatchStarCount >= 1 && levelInfoPtr.StarCount < 1) || (finalMatchStarCount >= 4 && levelInfoPtr.StarCount < 4) || (finalMatchStarCount >= 7 && levelInfoPtr.StarCount < 7) {
			///1星给球币与潜力点
			team.AwardCoin(npcTeamTypeStaticData.AddCoin)     ///给球队奖球币
			team.AwardTicket(npcTeamTypeStaticData.AddTicket) ///给球队奖钻石
			//	team.AwardTalentPoint(npcTeamTypeStaticData.AddTalentPoint)
			talentPoint = npcTeamTypeStaticData.AddTeamExp
			awardCoin = npcTeamTypeStaticData.AddCoin
		}

		///2星/3星 给予道具
		awardItemCountList = self.processActionAwardItem(client, finalMatchStarCount-npcTeamBaseStarCount, levelInfoPtr.StarCount-npcTeamBaseStarCount)

		levelInfoPtr.StarCount = finalMatchStarCount ///更新最新的星级
		syncMgr.SyncObject("PassLevelMsg", level)    ///同步关卡信息到客户端
	}

	///发放三星掉落奖励
	triStarAwardID := 0     ///三星掉落奖励道具ID
	triStarAwardNumber := 0 ///三星掉落奖励道具数量
	if isTriStarMatch {
		triStarAwardID = npcTeamTypeStaticData.AddItemType2
		triStarAwardNumberMax := npcTeamTypeStaticData.AddItemMaxCount2           ///掉落的最大数量
		triStarAwardNumberMin := npcTeamTypeStaticData.AddItemMinCount2           ///掉落的最小数量
		triStarAwardNumber = Random(triStarAwardNumberMin, triStarAwardNumberMax) ///随机取得道具数量
		team.AwardObject(triStarAwardID, triStarAwardNumber, 0, 0)                ///获取道具
		//if awardResult {                                                          ///获取成功就要更新到道具数量数组
		awardItemCountList = append(awardItemCountList, 0) ///客户端要求把数量放第二位
		awardItemCountList = append(awardItemCountList, triStarAwardNumber)
		//}
	}

	///	userGoalCount = 0
	//if userGoalCount <= npcGoalCount { ///打胜了奖励道具
	///平局或输了设置CD
	//		levelMatchCD := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTeam, ///读取cd间隔
	//			configTeamLevelMatchCommon, 1)
	///VIP整顿时间减少
	//vipLevel := team.GetVipLevel()
	// if vipLevel != 0 {
	// 	vipType := GetServer().GetStaticDataMgr().GetVipInfo(vipLevel)
	// 	levelMatchCD = vipType.Param14
	// }

	//		processMgr.SetTeamCD(teamLevelMatchCDType, levelMatchCD)
	//	matchCD = 0
	//}
	syncMgr.SyncObject("PassLevelMsg", team)                                    ///同步球队属性信息
	taskMgr.UpdateNpcTeamTask(client, npcTeamType, userGoalCount, npcGoalCount) ///更新任务信息

	passLevelMatchResultMsg := new(PassLevelMatchResultMsg)
	passLevelMatchResultMsg.LevelID = levelID
	passLevelMatchResultMsg.LevelType = self.LevelType
	passLevelMatchResultMsg.NpcTeamType = self.Param1
	passLevelMatchResultMsg.HomeGoal = userGoalCount
	passLevelMatchResultMsg.GuestGoal = npcGoalCount
	passLevelMatchResultMsg.AddCoin = awardCoin
	passLevelMatchResultMsg.AddExp = awardExp
	passLevelMatchResultMsg.AddTalentPoint = talentPoint
	passLevelMatchResultMsg.AddItemCountList = awardItemCountList
	passLevelMatchResultMsg.AddStarIDList = awardStarIDList
	//passLevelMatchResultMsg.LevelMatchCDUTC = Now() + matchCD ///同步比赛cd时间秒
	client.SendMsg(passLevelMatchResultMsg)
	team.SpendStarsContractPoint(1) ///打一场比赛,球队上阵球员均扣一点契约点数
	return true
}

func (self *PassLevelMsg) processActionCommonCheck(client IClient) bool { ///处理过关时的通用检测
	levelMgr := client.GetTeam().GetLevelMgr()
	level := levelMgr.FindLevel(self.LevelType)
	if level != nil && self.LevelSort != levelTypeMatch {
		return false ///非比赛关卡禁止重复过关
	}
	levelType := GetServer().GetStaticDataMgr().Unsafe().GetLevelType(self.LevelType)
	if nil == levelType {
		return false ///无效的关卡类型
	}
	if levelType.Untie > 0 { ///有前置条件要求
		levelUntie := levelMgr.FindLevel(levelType.Untie)
		if levelUntie == nil {
			loger := GetServer().GetLoger()
			loger.Warn("Untie level not pass  LevelType: %d  LevelID: %d Untie:%d", levelType.ID, self.LevelID, levelType.Untie)
			return false ///有前置条件但没满足
		}
		if levelUntie.IsPass() == false {
			return false ///战斗关卡未赢得1星以上
		}
	}
	if self.LevelSort != levelType.Type {
		return false ///检查客户端参数正确性
	}
	return true
}

func (self *PassLevelMsg) processAction(client IClient) bool { ///客户端请求通过指定关卡处理流程
	result := false
	if self.processActionCommonCheck(client) == false { ///处理过关时的通用检测
		return false
	}
	switch self.LevelSort {
	case levelTypeMatch:
		result = self.processActionMatch(client) ///处理比赛关卡逻辑
	case levelTypeLocker:
		result = self.processActionLocker(client) ///处理关卡锁处理逻辑
	case levelTypeStart:
		result = self.processActionStartFinish(client) ///起点与终点处理逻辑
	case levelTypeFinish:
		result = self.processActionStartFinish(client) ///起点与终点处理逻辑
	case levelTypeStarCard:
		result = self.processActionStarCard(client) ///处理球员卡请求
	case levelTypeItemCard:
		result = self.processActionItemCard(client) ///处理道具卡请求
	case levelTypeSkillCard:
		result = self.processActionSkillCard(client) ///处理技能卡请求
	case levelTypeFunctionCard:
		result = self.processActionFunction(client) ///处理功能卡请求
	}
	return result
}

type PassLevelMatchResultMsg struct { ///通过关卡比赛结果
	MsgHead          `json:"head"` ///"level", "passlevelresult"
	LevelType        int           `json:"leveltype"`        ///请求通关的关卡类型
	LevelID          int           `json:"levelid"`          ///请求通关的关卡id
	NpcTeamType      int           `json:"npcteamtype"`      ///请求比赛的npc球队类型
	HomeGoal         int           `json:"homegoal"`         ///自己队进球数
	GuestGoal        int           `json:"guestgoal"`        ///目标队进球数
	AddCoin          int           `json:"addcoin"`          ///奖励游戏币
	AddExp           int           `json:"addexp"`           ///奖励球队经验
	AddTalentPoint   int           `json:"addtalentpoint"`   ///奖励球队培养点数
	AddItemCountList IntList       `json:"additemcountlist"` ///奖励球队道具数量列表
	AddStarIDList    IntList       `json:"addstaridlist"`    ///奖励球队球员id列表
	LevelMatchCDUTC  int           `json:"levelmatchcdutc"`  ///球队比赛cd秒,0表示没有cd
}

func (self *PassLevelMatchResultMsg) GetTypeAndAction() (string, string) {
	return "level", "passlevelresult"
}

type GetPassLevelAwardMsg struct {
	MsgHead    `json:"head"` //"level", "passlevelaward"
	LeagueType int           `json:"leaguetype"` //地图类型
	AwardType  int           `json:"awardtype"`  //领奖类型: 1 关卡锁  2 星级  3 总进度
}

func (self *GetPassLevelAwardMsg) GetTypeAndAction() (string, string) {
	return "level", "passlevelaward"
}

//func (self *GetPassLevelAwardMsg) CountCheck(levelMgr *LevelMgr, leagueAwardInfo *LeagueAwardInfo) (bool, IntList, IntList) {
//	hasRecv := false
//	awardList, countList := IntList{}, IntList{}
//	loger := GetServer().GetLoger()
//	staticDataMgr := GetServer().GetStaticDataMgr()
//	leagueAwardType := staticDataMgr.GetLeagueAwardType(self.LeagueType)
//	if loger.CheckFail("leagueAwardType != nil", leagueAwardType != nil, leagueAwardType, nil) {
//		return false, awardList, countList //静态库不存在
//	}

//	begin, end, needCount, currentCount := 0, 0, 0, 0
//	itemID, itemNum := 0, 0
//	switch self.AwardType {
//	case leagueStarAward:
//		begin = leagueAwardStar1
//		end = leagueAwardStar3
//		currentCount = levelMgr.GetLevelStarCount(self.LeagueType)
//	case leagueLockAward:
//		begin = leagueAwardLock1
//		end = leagueAwardLock3
//		currentCount = levelMgr.GetLockNum(self.LeagueType)
//	case leagueScheduleAward:
//		begin = leagueAwardChe1
//		end = leagueAwardChe3
//		currentCount = levelMgr.GetSchedule(self.LeagueType)
//	}

//	for i := begin; i <= end; i++ {
//		isRecv := levelMgr.AwardIsRecv(*leagueAwardInfo, i)
//		if isRecv == true { //奖励已领取过
//			continue
//		}

//		switch self.AwardType {
//		case leagueStarAward:
//			switch i {
//			case leagueAwardStar1:
//				needCount = leagueAwardType.Stars1
//				itemID = leagueAwardType.Item1
//				itemNum = leagueAwardType.Num1
//			case leagueAwardStar2:
//				needCount = leagueAwardType.Stars2
//				itemID = leagueAwardType.Item2
//				itemNum = leagueAwardType.Num2
//			case leagueAwardStar3:
//				needCount = leagueAwardType.Stars3
//				itemID = leagueAwardType.Item3
//				itemNum = leagueAwardType.Num3
//			}
//		case leagueLockAward:
//			switch i {
//			case leagueAwardLock1:
//				needCount = leagueAwardType.Lock1
//				itemID = leagueAwardType.Item4
//				itemNum = leagueAwardType.Num4
//			case leagueAwardLock2:
//				needCount = leagueAwardType.Lock2
//				itemID = leagueAwardType.Item5
//				itemNum = leagueAwardType.Num5
//			case leagueAwardLock3:
//				needCount = leagueAwardType.Lock3
//				itemID = leagueAwardType.Item6
//				itemNum = leagueAwardType.Num6
//			}
//		case leagueScheduleAward:
//			switch i {
//			case leagueAwardChe1:
//				needCount = leagueAwardType.Che1
//				itemID = leagueAwardType.Item7
//				itemNum = leagueAwardType.Num7
//			case leagueAwardChe2:
//				needCount = leagueAwardType.Che2
//				itemID = leagueAwardType.Item8
//				itemNum = leagueAwardType.Num8
//			case leagueAwardChe3:
//				needCount = leagueAwardType.Che3
//				itemID = leagueAwardType.Item9
//				itemNum = leagueAwardType.Num9
//			}
//		}

//		if currentCount < needCount {
//			break ///未到领取条件
//		}

//		awardList = append(awardList, itemID)
//		countList = append(countList, itemNum)
//		hasRecv = true
//	}

//	return hasRecv, awardList, countList
//}

///判断当前条件下是否可以领奖
func (self *GetPassLevelAwardMsg) CanAcceptAward(client IClient) (bool, int, int) {
	topAwardCount := 3 ///领奖上限次数
	team := client.GetTeam()
	levelMgr := team.GetLevelMgr()
	staticDataMgr := GetServer().GetStaticDataMgr()
	leagueAward := levelMgr.FindLeagueAward(self.LeagueType)
	leagueAwardType := staticDataMgr.GetLeagueAwardType(self.LeagueType)
	processTargetCountList := IntList{leagueAwardType.Che1, leagueAwardType.Che2, leagueAwardType.Che3}
	starsTargetCountList := IntList{leagueAwardType.Stars1, leagueAwardType.Stars2, leagueAwardType.Stars3}
	locksTargetCountList := IntList{leagueAwardType.Lock1, leagueAwardType.Lock2, leagueAwardType.Lock3}
	starsAwardItemTypeList := IntList{leagueAwardType.Item1, leagueAwardType.Item2, leagueAwardType.Item3}
	starsAwardCountTypeList := IntList{leagueAwardType.Num1, leagueAwardType.Num2, leagueAwardType.Num3}
	locksAwardItemTypeList := IntList{leagueAwardType.Item4, leagueAwardType.Item5, leagueAwardType.Item6}
	locksAwardCountTypeList := IntList{leagueAwardType.Num4, leagueAwardType.Num5, leagueAwardType.Num6}
	processAwardItemTypeList := IntList{leagueAwardType.Item7, leagueAwardType.Item8, leagueAwardType.Item9}
	processAwardCountTypeList := IntList{leagueAwardType.Num7, leagueAwardType.Num8, leagueAwardType.Num9}
	currentCount, targetCount, index := 0, 0, 0
	awardItemType, awardItemCount := 0, 0
	canAcceptAward := false
	switch self.AwardType {
	case leagueStarAward:
		currentCount = levelMgr.GetLevelStarCount(self.LeagueType)
		index = Min(leagueAward.StarsAwardNum, 2)
		targetCount = starsTargetCountList[index]
		awardItemType = starsAwardItemTypeList[index]
		awardItemCount = starsAwardCountTypeList[index]
		canAcceptAward = (currentCount >= targetCount) && (leagueAward.StarsAwardNum < topAwardCount)
	case leagueLockAward:
		currentCount = levelMgr.GetLockNum(self.LeagueType)
		index = Min(leagueAward.LocksAwardNum, 2)
		targetCount = locksTargetCountList[index]
		awardItemType = locksAwardItemTypeList[index]
		awardItemCount = locksAwardCountTypeList[index]
		canAcceptAward = (currentCount >= targetCount) && (leagueAward.LocksAwardNum < topAwardCount)
	case leagueScheduleAward:
		index = Min(leagueAward.ProcessAwardNum, 2)
		currentCount = levelMgr.GetSchedule(self.LeagueType)
		targetCount = processTargetCountList[index]
		awardItemType = processAwardItemTypeList[index]
		awardItemCount = processAwardCountTypeList[index]
		canAcceptAward = (currentCount >= targetCount) && (leagueAward.ProcessAwardNum < topAwardCount)
	}
	return canAcceptAward, awardItemType, awardItemCount
}

func (self *GetPassLevelAwardMsg) checkActikon(client IClient) bool {
	team := client.GetTeam()
	loger := GetServer().GetLoger()
	levelMgr := team.GetLevelMgr()
	staticDataMgr := GetServer().GetStaticDataMgr()
	leagueAwardType := staticDataMgr.GetLeagueAwardType(self.LeagueType)
	if loger.CheckFail("leagueAwardType != nil", leagueAwardType != nil, leagueAwardType, nil) {
		return false //联赛类型不存在
	}
	leagueAward := levelMgr.FindLeagueAward(self.LeagueType)
	if leagueAward == nil {
		leagueAward = levelMgr.AddLeagueAwardInfo(self.LeagueType)
	}
	canAcceptAward, _, _ := self.CanAcceptAward(client)
	if loger.CheckFail("canAcceptAward==true", canAcceptAward == true, canAcceptAward, true) {
		return false //联赛类型不存在
	}
	//isCanRecv, _, _ := self.CountCheck(levelMgr, &leagueAward.LeagueAwardInfo)
	//if loger().CheckFail("isCanRecv == true", isCanRecv == true, isCanRecv, true) {
	//	return false ///不满足领取条件
	//}

	return true
}

//func (self *GetPassLevelAwardMsg) payAction(client IClient) bool {
//	//	syncMgr := client.GetSyncMgr()
//	team := client.GetTeam()
//	levelMgr := team.GetLevelMgr() //发奖
//	leagueAward := levelMgr.FindLeagueAward(self.LeagueType)
//	switch self.AwardType {
//	case leagueStarAward:
//		leagueAward.StarsAwardNum++
//	case leagueLockAward:
//		leagueAward.LocksAwardNum++
//	case leagueScheduleAward:
//		leagueAward.ProcessAwardNum++
//	}
//	return true
//}

type PassLevelAwardMsgResult struct {
	MsgHead    `json:"head"` //"level", "passlevelaward"
	LeagueType int           `json:"leaguetype"` //地图类型
	AwardType  int           `json:"awardtype"`  //领奖类型: 1 关卡锁  2 星级  3 总进度
	Result     int           `json:"result"`     ///1是成功,0是失败
}

func (self *PassLevelAwardMsgResult) GetTypeAndAction() (string, string) {
	return "level", "passlevelawardmsgresult"
}

func (self *GetPassLevelAwardMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	_, awardItemType, awardItemCount := self.CanAcceptAward(client)
	levelMgr := team.GetLevelMgr() //发奖
	leagueAward := levelMgr.FindLeagueAward(self.LeagueType)
	switch self.AwardType {
	case leagueStarAward:
		leagueAward.StarsAwardNum++
	case leagueLockAward:
		leagueAward.LocksAwardNum++
	case leagueScheduleAward:
		leagueAward.ProcessAwardNum++
	}
	team.AwardObject(awardItemType, awardItemCount, 0, 0)
	if awardItemType == awardTypeTicket && awardItemCount > 0 {
		client.RechargeRecord(Get_PassLevelAward, awardItemCount)
	}
	return true
}

func (self *GetPassLevelAwardMsg) processAction(client IClient) (result bool) {
	defer func() {
		passLevelAwardMsgResult := new(PassLevelAwardMsgResult)
		passLevelAwardMsgResult.LeagueType = self.LeagueType
		passLevelAwardMsgResult.AwardType = self.AwardType
		passLevelAwardMsgResult.Result = 1
		if false == result {
			passLevelAwardMsgResult.Result = 0
		}
		client.SendMsg(passLevelAwardMsgResult)
	}()
	if self.checkActikon(client) == false {
		return false
	}
	//if self.payAction(client) == false {
	//	return false
	//}
	if self.doAction(client) == false {
		return false
	}
	return true
}

type QueryMaxStartTimesMsg struct { ///请求冠军之路三星后剩余可挑战次数
	MsgHead   `json:"head"` ///"level", "querytimes"
	LevelSort int           `json:"levelsort"` ///请求通关的关卡子类型
	LevelType int           `json:"leveltype"` ///请求通关的关卡类型
	LevelID   int           `json:"levelid"`   ///请求通关的关卡id
	NpcID     int           `json:"npcid"`     ///挑战球队ID
}

func (self *QueryMaxStartTimesMsg) GetTypeAndAction() (string, string) { ///实现消息处理接口的返回类型方法
	return "level", "querytimes"
}

func (self *QueryMaxStartTimesMsg) processAction(client IClient) bool { ///实现消息处理接口的处理消息方法
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

///检测
func (self *QueryMaxStartTimesMsg) checkAction(client IClient) bool {
	///基本检查
	loger := loger()                            ///记录对象
	levelMgr := client.GetTeam().GetLevelMgr()  ///取得关卡管理器
	level := levelMgr.FindLevel(self.LevelType) ///取得关卡对象
	if loger.CheckFail("level!=nil and LevelType == 1 ", level != nil && self.LevelSort == levelTypeMatch, level, nil) {
		return false ///非比赛关卡不能查询星级
	}
	levelType := GetServer().GetStaticDataMgr().Unsafe().GetLevelType(self.LevelType) ///取得关卡类型对象
	if loger.CheckFail("levelType!=nil", levelType != nil, levelType, nil) {
		return false ///无效的关卡类型
	}
	///先要找到球队
	teamIndex := levelType.GetNpcIndexInLevel(self.NpcID) ///球队在关卡中的位置
	if loger.CheckFail("npcteam can't find", teamIndex != 0, self, nil) {
		return false ///找不到对应的npc队伍
	}
	///必须3星通关
	if level.StarCount < teamIndex*3 {
		return false ///必须达到指定的星级
	}
	return true
}

///支付
func (self *QueryMaxStartTimesMsg) payAction(client IClient) bool {
	return true
}

///发货
func (self *QueryMaxStartTimesMsg) doAction(client IClient) bool {
	///取得重置数据对象
	attribType := GetTriStarAttribType(self.NpcID) ///取得可重置数据的类型
	team := client.GetTeam()                       ///取得玩家队伍
	resetAttribMgr := team.GetResetAttribMgr()     ///取得队伍的可重置管理器
	resetAttrib := resetAttribMgr.GetResetAttrib(attribType)

	challangeNumber := 3     ///可挑战次数，暂时硬编码
	challangeResetClock := 4 ///重置钟点暂，时硬编码
	if resetAttrib != nil {  ///如果数据存在就检测是否到期
		if IsExpireTime(resetAttrib.ResetTime) { ///如果已经过期就重置
			resetAttribMgr.ResetTriStar(self.NpcID) ///重置数据
		}
	} else { ///数据不存在就创建新的
		values := []int{challangeNumber}
		resetAttrib = resetAttribMgr.AddResetAttrib(attribType, GetHourUTC(challangeResetClock), values)
	}

	queryMaxStartTimesResultMsg := new(QueryMaxStartTimesResultMsg)
	queryMaxStartTimesResultMsg.NpcID = self.NpcID
	queryMaxStartTimesResultMsg.Remain = resetAttrib.Value1
	client.SendMsg(queryMaxStartTimesResultMsg) ///发送给客户端返回信息
	return true
}

type QueryMaxStartTimesResultMsg struct { ///  返回冠军之路三星后剩余可挑战次数
	MsgHead `json:"head"` /// " level ", "querytimesresult"
	NpcID   int           `json:"npcid"`  ///挑战球队ID
	Remain  int           `json:"remain"` ///剩余挑战次数
}

func (self *QueryMaxStartTimesResultMsg) GetTypeAndAction() (string, string) { ///实现消息处理接口的返回类型方法
	return "level", "querytimesresult"
}
