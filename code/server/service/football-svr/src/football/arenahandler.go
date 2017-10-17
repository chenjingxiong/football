package football

import (
//	"fmt"
)

type QueryArenaMatchResultMsg struct { ///查询联赛信息结果消息
	MsgHead        `json:"head"`          ///"arenamatch", "queryresult"
	ArenaInfoSlice `json:"arenainfolist"` ///联赛信息列表
	CurrentDay     int                    `json:"currentday"` ///联赛当前天数
}

func (self *QueryArenaMatchResultMsg) GetTypeAndAction() (string, string) {
	return "arenamatch", "queryresult"
}

type QueryArenaMatchMsg struct { ///查询联赛信息消息
	MsgHead `json:"head"` ///"arenamatch", "queryinfo"
}

func (self *QueryArenaMatchMsg) GetTypeAndAction() (string, string) {
	return "arenamatch", "queryinfo"
}

//func (self *QueryArenaMatchResultMsg) broacastMsg(userMgr *UserMgr) bool {
//	loger := loger()
//	senderClient := userMgr.GetClientByTeamID(self.ArenaInfo.TeamID) ///得到发送者
//	if nil == senderClient {
//		return false
//	}
//	if loger.CheckFail("self.ArenaInfo.Group<=0", self.ArenaInfo.Group <= 0, self.ArenaInfo.Group, 0) {
//		return false ///已有分组号的消息不应该到广播函数中来
//	}
//	///分配新的分组号

//	senderClient.SendMsg(self) ///转发给接受收
//	return true
//}
func SendQueryArenaMatchResultMsg(client IClient) {
	team := client.GetTeam()
	arenaMgr := team.GetArenaMgr()
	arenaMgr.RefreshArenaDate() ///尝试更新联赛数据
	queryArenaMatchResultMsg := new(QueryArenaMatchResultMsg)

	currentDay := GetOpenServerSamsaraDay_Plus(ArenaSettleAccounts)

	if currentDay < 0 {
		currentDay = 0
	}

	currentDay = 3 - currentDay
	//if currentDay <= 0 {
	//	currentDay = ArenaSettleAccounts
	//}

	queryArenaMatchResultMsg.CurrentDay = currentDay
	queryArenaMatchResultMsg.ArenaInfoSlice = arenaMgr.GetArenaInfoSlice()
	client.SendMsg(queryArenaMatchResultMsg) ///发送给客户端此用户所在竞技场信息
}

func (self *QueryArenaMatchMsg) checkAction(client IClient) bool {

	team := client.GetTeam()
	loger := GetServer().GetLoger()
	isOpenFunction := TestMask(team.FunctionMask, functionMaskArenaMatch)

	if loger.CheckFail("isOpenFunction == true", isOpenFunction == true, isOpenFunction, true) {
		return false //!新增条件限制: 未开启天天联赛,则名单不加入联赛数据
	}

	return true
}

func (self *QueryArenaMatchMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false {
		return false
	}

	SendQueryArenaMatchResultMsg(client)
	return true
}

type QueryTeamInfoResultMsg struct { ///查询球队信息结果
	MsgHead       `json:"head"`          ///"team", "queryteaminforesult"
	TeamInfo      `json:"teaminfo"`      ///球队的teaminfo
	StarInfoList  `json:"starlist"`      ///球队首发球员列表
	FormationInfo `json:"formationinfo"` ///球队阵型信息
	ItemInfoList  `json:"equipmentlist"` ///首发球员装备道具信息
	StarSkill     []*SkillState          `json:"starskill"`     //!球员技能信息
	MannaStarList MannaStarSlice         `json:"mannastarlist"` //!自创球员信息
}

func (self *QueryTeamInfoResultMsg) GetTypeAndAction() (string, string) {
	return "team", "queryteaminforesult"
}

func (self *QueryTeamInfoMsg) skillIsStudy(team *Team, skillType int, starID int) bool {
	skillMgr := team.GetSkillMgr()
	isStudy := false
	starSkillLst := skillMgr.GetStarSkillSlice(starID)
	for i := 0; i < len(starSkillLst); i++ {
		v := starSkillLst[i]
		if v.Type == skillType {
			isStudy = true
			break
		}
	}
	return isStudy
}

func (self *QueryTeamInfoMsg) checkSkillState(team *Team, skillType int, starID int) (int, int) {
	attrMgr := team.GetResetAttribMgr()
	skillAttrMgr := attrMgr.GetResetAttrib(ResetAttribTypeSkillStudyInfo)

	state := 0
	time := 0

	//!判断技能是否已学习
	isStudy := self.skillIsStudy(team, skillType, starID)
	if isStudy == true {
		state = 2
		time = 0

		return state, time
	}

	//!判断技能是否在学习中
	if skillAttrMgr.Value1 == starID && skillAttrMgr.Value2 == skillType && skillAttrMgr.Value3 != 0 {
		state = 1
		time = skillAttrMgr.ResetTime - Now()

		return state, time
	}

	//!否则为未学习技能
	return state, time
}

func (self *QueryTeamInfoMsg) GetSkillInfoList(team *Team) []*SkillState {
	attrMgr := team.GetResetAttribMgr()
	skillAttrMgr := attrMgr.GetResetAttrib(ResetAttribTypeSkillStudyInfo)
	if skillAttrMgr == nil {
		//! 不存在则创建默认
		skillAttrMgr = attrMgr.AddResetAttrib(ResetAttribTypeSkillStudyInfo, -1, IntList{0, 0, 0})
	}

	starSkill := []*SkillState{}
	staticDataMgr := GetServer().GetStaticDataMgr()
	allStarList := team.GetAllStarList()
	for i := 0; i < allStarList.Len(); i++ {
		star := team.GetStar(allStarList[i])

		if star.IsMannaStar == 1 {

			mannaStarMgr := team.GetMannaStarMgr()
			starType := mannaStarMgr.GetMannaStar(star.Type)

			if starType.Skill1 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill1
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill1, star.ID)
				starSkill = append(starSkill, node)
			}

			if starType.Skill2 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill2
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill2, star.ID)
				starSkill = append(starSkill, node)
			}

			if starType.Skill3 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill3
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill3, star.ID)

				starSkill = append(starSkill, node)
			}

			if starType.Skill4 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill4
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill4, star.ID)

				starSkill = append(starSkill, node)
			}

		} else {
			starType := staticDataMgr.GetStarType(star.Type)

			if starType.Skill1 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill1
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill1, star.ID)
				starSkill = append(starSkill, node)
			}

			if starType.Skill2 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill2
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill2, star.ID)
				starSkill = append(starSkill, node)
			}

			if starType.Skill3 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill3
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill3, star.ID)

				starSkill = append(starSkill, node)
			}

			if starType.Skill4 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill4
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill4, star.ID)

				starSkill = append(starSkill, node)
			}
		}
	}

	//!调试输出
	// for i := 0; i < len(starSkill); i++ {
	// 	v := starSkill[i]
	// 	fmt.Println("target team starID: ", v.StarID)
	// 	fmt.Println("target team skillType: ", v.SkillType)
	// 	fmt.Println("target teams killState: ", v.SkillState)
	// 	fmt.Println("target team time: ", v.Time)
	// }

	return starSkill
}

type QueryTeamInfoMsg struct { ///查询球队信息
	MsgHead `json:"head"` ///"team", "queryteaminfo"
	TeamID  int           `json:"teamid"` ///查询球队id
}

func (self *QueryTeamInfoMsg) GetTypeAndAction() (string, string) {
	return "team", "queryteaminfo"
}

func (self *QueryTeamInfoMsg) GetMannaStarList(team *Team) MannaStarSlice {
	starList := MannaStarSlice{}
	mannaStarMgr := team.GetMannaStarMgr()
	for i := MannaStarSeatOne; i <= MannaStarSeatThree; i++ {
		star := mannaStarMgr.GetMannaStarFromSeat(i)
		if star == nil {
			continue
		}

		//! 获取自创球员信息
		starList = append(starList, &star.MannaStarType)
	}
	return starList
}

func (self *QueryTeamInfoMsg) processAction(client IClient) bool {
	syncMgr := client.GetSyncMgr()

	teamTarget := new(Team)
	teamTarget.Create(0, self.TeamID, syncMgr) ///加载球队信息
	queryTeamInfoResultMsg := new(QueryTeamInfoResultMsg)
	queryTeamInfoResultMsg.TeamInfo = teamTarget.TeamInfo                               ///放入teamInfo
	queryTeamInfoResultMsg.FormationInfo = *teamTarget.GetCurrentFormObject().GetInfo() ///得到目标球队首发阵形信息
	queryTeamInfoResultMsg.ItemInfoList = teamTarget.GetCurrentFormStarItemInfo()       ///得到球队首发球员装备信息列表
	queryTeamInfoResultMsg.StarInfoList = teamTarget.GetStartersList()                  ///得到首发球员列表
	queryTeamInfoResultMsg.StarSkill = self.GetSkillInfoList(teamTarget)                //!得到对手球队技能信息

	starList := MannaStarSlice{}

	mannaStarList := self.GetMannaStarList(teamTarget)
	for i := 0; i < len(mannaStarList); i++ {
		mannaStarList[i].MannaSeat = 0
		starList = append(starList, mannaStarList[i])
	}

	mannaStarList = self.GetMannaStarList(client.GetTeam())
	for i := 0; i < len(mannaStarList); i++ {
		starList = append(starList, mannaStarList[i])
	}

	queryTeamInfoResultMsg.MannaStarList = self.GetMannaStarList(teamTarget) //!获取自创球员信息
	client.SendMsg(queryTeamInfoResultMsg)
	return true
}

type PlayArenaMatchResultMsg struct { ///请求比赛
	MsgHead     `json:"head"` ///"arenamatch", "matchresult"
	NpcTeamType int           `json:"npcteamtype"` ///请求比赛的npc球队类型
	TeamID      int           `json:"teamid"`      ///请求比赛的玩家球队
	HomeGoal    int           `json:"homegoal"`    ///自己队进球数
	GuestGoal   int           `json:"guestgoal"`   ///目标队进球数
}

func (self *PlayArenaMatchResultMsg) GetTypeAndAction() (string, string) {
	return "arenamatch", "matchresult"
}

type PlayArenaMatchMsg struct { ///请求比赛
	MsgHead     `json:"head"` ///"arenamatch", "playmatch"
	TeamID      int           `json:"teamid"`      ///对手球队id,0表示打npc球队
	NpcTeamType int           `json:"npcteamtype"` ///请求比赛的npc球队类型
	TeamNumber  int           `json:"indexteam"`   ///请求比赛的球队索引号,从1开始
}

func (self *PlayArenaMatchMsg) GetTypeAndAction() (string, string) {
	return "arenamatch", "playmatch"
}

func (self *PlayArenaMatchMsg) checkAction(client IClient, targetTeam *Team) bool {
	loger := loger()
	team := client.GetTeam()
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	arenaMgr := team.GetArenaMgr()
	arenaSelfInfo := arenaMgr.GetInfo()
	teamCount := arenaMgr.arenaIDList.Len()
	if self.TeamID > 0 { ///球家球队需要验证球队索引合法性
		if loger.CheckFail("self.TeamNumber>0 && self.TeamNumber<=teamCount",
			self.TeamNumber > 0 && self.TeamNumber <= teamCount,
			self.TeamNumber, teamCount) {
			return false ///客户端发来的球队序列号必须在有效范围内
		}
		arenaTargetID := arenaMgr.arenaIDList[self.TeamNumber-1] ///得到竞技场信息编号
		//fmt.Println(arenaMgr.arenaInfoList)
		arenaTargetInfo := arenaMgr.arenaInfoList[arenaTargetID]
		if loger.CheckFail("arenaInfo!=nil", arenaTargetInfo != nil, arenaTargetInfo, nil) {
			return false ///竞技场信息对象必须有效
		}
		if loger.CheckFail(" arenaTargetInfo.TeamID==self.TeamID", arenaTargetInfo.TeamID == self.TeamID,
			arenaTargetInfo.TeamID, self.TeamID) {
			return false ///客户端请求攻击的球队id必须与实际相符
		}
	} else {
		npcTeamType := staticDataMgr.GetNpcTeamType(self.NpcTeamType)
		if loger.CheckFail("npcTeamType!=nil", npcTeamType != nil, npcTeamType, nil) {
			return false ///客户端请求攻击的球队id必须与实际相符
		}
	}

	//根据联赛当前天数,获取可比赛次数
	currentDay := GetOpenServerSamsaraDay(ArenaSettleAccounts)
	totalMatchCount := staticDataMgr.GetConfigStaticDataInt(configArenaMatch, configArenaMatchCommon, 1)     ///总次数
	firstDayMatchCount := staticDataMgr.GetConfigStaticDataInt(configArenaMatch, configArenaMatchCommon, 4)  ///第一天比赛次数
	secondDayMatchCount := staticDataMgr.GetConfigStaticDataInt(configArenaMatch, configArenaMatchCommon, 5) ///第二天比赛次数
	LastDatMatchCount := staticDataMgr.GetConfigStaticDataInt(configArenaMatch, configArenaMatchCommon, 6)   ///第三天比赛次数
	remainMatchCount := 0                                                                                    ///当天可比赛次数
	switch currentDay {
	case 0:
		remainMatchCount = totalMatchCount - firstDayMatchCount
	case 1:
		remainMatchCount = totalMatchCount - firstDayMatchCount - secondDayMatchCount
	case 2:
		remainMatchCount = totalMatchCount - firstDayMatchCount - secondDayMatchCount - LastDatMatchCount
	}

	if loger.CheckFail("arenaSelfInfo.RemainMatchCoun > remainMatchCount", arenaSelfInfo.RemainMatchCount > remainMatchCount,
		arenaSelfInfo.RemainMatchCount, remainMatchCount) {
		return false ///可比赛次数必须大于
	}

	if loger.CheckFail("arenaSelfInfo.RemainMatchCoun > 0", arenaSelfInfo.RemainMatchCount > 0,
		arenaSelfInfo.RemainMatchCount, 0) {
		return false ///可比赛次数必须大于
	}

	isPlayedMatch := TestMask(arenaSelfInfo.PlayMask, self.TeamNumber) ///判断是否已打过比赛了
	if loger.CheckFail("isPlayedMatch", isPlayedMatch == false, isPlayedMatch, false) {
		return false ///必须未跟此球队打过比赛
	}
	return true
}

func (self *PlayArenaMatchMsg) payAction(client IClient, targetTeam *Team) bool {
	team := client.GetTeam()
	arenaMgr := team.GetArenaMgr()
	arenaInfo := arenaMgr.GetInfo()
	arenaInfo.RemainMatchCount--                                         ///消费一次剩余比赛次数
	arenaInfo.PlayMask = SetMask(arenaInfo.PlayMask, self.TeamNumber, 1) ///打上此序列球队已比赛的掩码
	return true
}

///根据比赛结果更新联赛信息
func (self *PlayArenaMatchMsg) updateArenaMatchInfo(client IClient, userGoalCount int, targetGoalCount int) {
	team := client.GetTeam()
	arenaMgr := team.GetArenaMgr()
	arenaInfo := arenaMgr.GetInfo()
	oldScore := arenaInfo.Score          ///保存历史数据
	if userGoalCount > targetGoalCount { ///胜
		arenaInfo.Score += 3
		arenaInfo.WinCount++
	} else if userGoalCount == targetGoalCount { ///平
		arenaInfo.Score += 1
		arenaInfo.DrawCount++
	} else { ///负
		arenaInfo.LostCount++
	}
	if arenaInfo.Score > oldScore {
		arenaInfo.UpdateUTC = Now()
	}
}

func (self *PlayArenaMatchMsg) doAction(client IClient, teamTarget *Team) bool {
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	teamSelf := client.GetTeam()
	userGoalCount, targetGoalCount := 0, 0
	//	mailMgr := teamSelf.GetMailMgr()
	arenaMgr := teamSelf.GetArenaMgr()
	playArenaMatchResultMsg := new(PlayArenaMatchResultMsg)
	if teamTarget != nil { ///判断是否攻击玩家球队
		userGoalCount, targetGoalCount = teamSelf.CalcMatchResult(teamTarget)
		playArenaMatchResultMsg.TeamID = self.TeamID
		//		targetFormationMgr := teamTarget.GetFormationMgr()
		//		formationID := teamTarget.GetFormation()
		//		targetFormation := targetFormationMgr.GetFormation(formationID)
		//		currentTactic := targetFormation.CurrentTactic
		//		currentType := targetFormation.Type
		//mailMgr.SendMatchReport(ReportMail, ArenaSend, teamTarget.Name, teamTarget.Icon, 0, currentType,
		//	currentTactic, userGoalCount, targetGoalCount, teamSelf.Score, teamTarget.Score)

		//playArenaMatchResultMsg.TeamInfo = teamTarget.TeamInfo                               ///放入teamInfo
		//playArenaMatchResultMsg.FormationInfo = *teamTarget.GetCurrentFormObject().GetInfo() ///得到目标球队首发阵形信息
		//playArenaMatchResultMsg.EquipmentList = teamTarget.GetCurrentFormStarItemInfo()      ///得到球队首发球员装备信息列表
		//playArenaMatchResultMsg.TargetStarList = teamTarget.GetStartersList()                ///得到首发球员列表
	} else {
		npcTeamTypeStaticData := staticDataMgr.GetNpcTeamType(self.NpcTeamType)
		userGoalCount, targetGoalCount = npcTeamTypeStaticData.CalcMatchResult(teamSelf) ///计算比赛结果
		playArenaMatchResultMsg.NpcTeamType = self.NpcTeamType                           ///放入npc球队类型

		//		npcScore := npcTeamTypeStaticData.AttackScore + npcTeamTypeStaticData.DefenseScore + npcTeamTypeStaticData.OrganizeScore
		//mailMgr.SendMatchReport(ReportMail, ArenaSend, npcTeamTypeStaticData.Name, npcTeamTypeStaticData.Icon, npcTeamTypeStaticData.ID, npcTeamTypeStaticData.Formation,
		//	npcTeamTypeStaticData.Tactical, userGoalCount, targetGoalCount, teamSelf.Score, npcScore)
	}
	playArenaMatchResultMsg.HomeGoal = userGoalCount    ///攻击球队进球数
	playArenaMatchResultMsg.GuestGoal = targetGoalCount ///防守球队进球数
	client.SendMsg(playArenaMatchResultMsg)
	self.updateArenaMatchInfo(client, userGoalCount, targetGoalCount) ///根据比赛成绩刷新联赛信息
	SendQueryArenaMatchResultMsg(client)                              ///刷新联赛信息给客户端
	arenaMgr.Save()                                                   ///打完一场比赛即时保存数据

	//新增奖励随机宝箱
	arenaInfo := arenaMgr.GetInfo()
	arenaType := GetArenaType(arenaInfo.ArenaType)
	teamSelf.AwardObject(arenaType.AwardItem, 1, 0, 0)
	///更新天天联赛中的日常任务
	teamSelf.GetTaskMgr().UpdateDayTaskFuntion(client.GetElement(), dayTaskFunctionArenaMatch)
	return true
}

func (self *PlayArenaMatchMsg) processAction(client IClient) bool {
	syncMgr := client.GetSyncMgr()
	var teamTarget *Team = nil
	if self.TeamID > 0 {
		teamTarget = new(Team)
		teamTarget.Create(0, self.TeamID, syncMgr) ///加载球队信息
	} else {
		teamTarget = nil //不存球家球队置空清除
	}
	//self.TeamNumber = 1                                /// for test
	if self.checkAction(client, teamTarget) == false { ///检测
		return false
	}
	if self.payAction(client, teamTarget) == false { ///支付
		return false
	}
	if self.doAction(client, teamTarget) == false { ///发货
		return false
	}
	return true
}

type AcceptArenaAwardMsg struct { ///请求领取竞技场奖励
	MsgHead `json:"head"` ///"arenamatch", "acceptaward"
}

func (self *AcceptArenaAwardMsg) GetTypeAndAction() (string, string) {
	return "arenamatch", "acceptaward"
}

func (self *AcceptArenaAwardMsg) checkAction(client IClient) bool {
	return true
}

const (
	AwardTicketNone     = 0 ///无奖励
	AwardTicketGave     = 1 ///已领奖
	AwardTicketPromote  = 2 ///晋级奖
	AwardTicketKeep     = 3 ///保级奖
	AwardTicketDemotion = 4 ///降级奖
)

func (self *AcceptArenaAwardMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	arenaMgr := team.GetArenaMgr()
	mailMgr := team.GetMailMgr()
	arenaInfo := arenaMgr.GetInfo()
	staticDataMgr := GetServer().GetStaticDataMgr()
	arenaType := GetArenaType(arenaInfo.ArenaType)
	awardCoin, awardTactic := 0, 0
	awardTicket := arenaInfo.AwardTicket
	arenaInfo.AwardTicket = AwardTicketGave ///打上已领奖标识
	switch awardTicket {
	case AwardTicketPromote:
		awardCoin = arenaType.UpAwardCoin
		awardTactic = arenaType.Upawardtactic

		//晋级奖励
		//if arenaInfo.ArenaType <= 5 && arenaInfo.ArenaType > 0 && TestMask(arenaInfo.AwardMask, arenaInfo.ArenaType) == false {

		//	rankAwaard := staticDataMgr.GetConfigStaticDataInt(configArenaMatch, configArenaPromotAward, 5-arenaInfo.ArenaType+1)
		//	mailMgr.SendSysAwardMail(SystemMail, ArenaRankUp, IntList{awardTypeTicket}, IntList{0}, IntList{rankAwaard})
		//	SetMask(arenaInfo.AwardMask, arenaInfo.ArenaType, 1) ///设置掩码
		//}
		//晋级奖励
		if 0 == arenaInfo.AwardMask || arenaInfo.AwardMask > arenaInfo.ArenaType { ///首次晋级
			rankAwaard := staticDataMgr.GetConfigStaticDataInt(configArenaMatch, configArenaPromotAward, 7-arenaInfo.ArenaType)
			if rankAwaard > 0 {
				mailMgr.SendSysAwardMail(ArenaFirstUpMail, ArenaRankUp, IntList{awardTypeTicket}, IntList{arenaInfo.ArenaType}, IntList{rankAwaard}, "", "")
				arenaInfo.AwardMask = arenaInfo.ArenaType
			}
		}
	case AwardTicketKeep:
		awardCoin = arenaType.KeepAwardCoin
		awardTactic = arenaType.KeepAwardTactic
	case AwardTicketDemotion:
		awardCoin = arenaType.DownAwardCoin
		awardTactic = arenaType.DownAwardTactic
	}
	team.AwardCoin(awardCoin)
	team.AwardTactic(awardTactic)
	client.GetSyncMgr().SyncObject("AcceptArenaAwardMsg", team)
	SendQueryArenaMatchResultMsg(client)
	return true
}

func (self *AcceptArenaAwardMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false { ///检测
		return false
	}
	if self.doAction(client) == false { ///发货
		return false
	}
	return true
}
