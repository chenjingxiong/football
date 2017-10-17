package football

///挑战赛消息处理器
import ()

type QueryChallengeTimesMsg struct { ///客户端请求挑战赛日剩余挑战次数
	MsgHead `json:"head"` ///"challenge", "querytimes"
	Type    int           `json:"type"` ///挑战赛类型
}

func (self *QueryChallengeTimesMsg) GetTypeAndAction() (string, string) { ///实现消息处理接口的返回类型方法
	return "challenge", "querytimes"
}

func (self *QueryChallengeTimesMsg) processAction(client IClient) bool { ///实现消息处理接口的处理消息方法
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
func (self *QueryChallengeTimesMsg) checkAction(client IClient) bool {
	loger := loger() ///记录对象
	/////需要开启功能
	team := client.GetTeam() ///取得玩家队伍
	//isEnable := team.TestFunctionMask(functionMaskChallangeMatch) ///验证掩码
	//if loger.CheckFail("ChallangeMatch is Enable", isEnable, team, nil) {
	//	return false
	//}
	challangeMatchType := GetChallangeMatchAttribType(self.Type)
	if loger.CheckFail("ChallangeMatchType is vailed", challangeMatchType == ResetAttribTypeChallangeMatchPerfect ||
		challangeMatchType == ResetAttribTypeChallangeMatchCrazy || challangeMatchType == ResetAttribTypeChallangeMatchDefend, team, nil) {
		return false
	}
	return true
}

///支付
func (self *QueryChallengeTimesMsg) payAction(client IClient) bool {
	return true
}

///发货
func (self *QueryChallengeTimesMsg) doAction(client IClient) bool {
	///取得重置数据对象
	team := client.GetTeam()                   ///取得玩家队伍
	resetAttribMgr := team.GetResetAttribMgr() ///取得队伍的可重置管理器
	resetAttrib := resetAttribMgr.QueryChallangeMatchResetAttrib(self.Type, client)

	queryChallengeTimesResultMsg := new(QueryChallengeTimesResultMsg)
	queryChallengeTimesResultMsg.Type = self.Type
	queryChallengeTimesResultMsg.Remain = resetAttrib.Value1
	client.SendMsg(queryChallengeTimesResultMsg) ///发送给客户端返回信息
	return true
}

type QueryChallengeTimesResultMsg struct { ///  返回挑战赛日剩余挑战次数
	MsgHead `json:"head"` /// "challenge", "querytimesresult"
	Type    int           `json:"type"`   ///挑战赛类型
	Remain  int           `json:"remain"` ///剩余挑战次数
}

func (self *QueryChallengeTimesResultMsg) GetTypeAndAction() (string, string) { ///实现消息处理接口的返回类型方法
	return "challenge", "querytimesresult"
}

type QueryChallengeFightMsg struct { //客户端请求挑战赛挑战
	MsgHead `json:"head"` ///"challenge", "queryfight"
	Type    int           `json:"type"`    ///挑战赛类型
	MatchID int           `json:"matchid"` ///挑战赛id
}

func (self *QueryChallengeFightMsg) GetTypeAndAction() (string, string) { ///实现消息处理接口的返回类型方法
	return "challenge", "queryfight"
}

func (self *QueryChallengeFightMsg) processAction(client IClient) bool { ///实现消息处理接口的处理消息方法
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
func (self *QueryChallengeFightMsg) checkAction(client IClient) bool {
	loger := loger() ///记录对象
	///检查功能开启
	team := client.GetTeam() ///取得玩家队伍
	//isEnable := team.TestFunctionMask(functionMaskChallangeMatch) ///验证掩码
	//if loger.CheckFail("ChallangeMatch is Enable", isEnable, team, nil) {
	//	return false
	//}
	challangeMatchStaticData := GetServer().GetStaticDataMgr().Unsafe().GetChallangeMatchType(self.MatchID) ///取挑战赛类型对象
	if loger.CheckFail("challangeMatchStaticData != nil", challangeMatchStaticData != nil, self.MatchID, nil) {
		return false
	}
	///检查申请类型和挑战赛类型是否一致
	if loger.CheckFail("self.type == challangMatchtype", self.Type == challangeMatchStaticData.Type, team, nil) {
		return false
	}
	///检查等级
	if loger.CheckFail("teamLV >= challangMatchLV", team.GetLevel() >= challangeMatchStaticData.LV, team, nil) {
		return false
	}
	///检查剩余次数
	resetAttribMgr := team.GetResetAttribMgr() ///取得队伍的可重置管理器
	resetAttrib := resetAttribMgr.GetChallangeMatchResetAttrib(challangeMatchStaticData.Type, client)
	if loger.CheckFail("resetAttrib != nil", resetAttrib != nil, team, nil) {
		return false
	}
	if loger.CheckFail("challangeNum > 0", resetAttrib.Value1 > 0, team, nil) { ///检查可用挑战次数
		return false
	}
	return true
}

///支付
func (self *QueryChallengeFightMsg) payAction(client IClient) bool {
	team := client.GetTeam()                   ///取得玩家队伍
	resetAttribMgr := team.GetResetAttribMgr() ///取得队伍的可重置管理器
	resetAttrib := resetAttribMgr.GetChallangeMatchResetAttrib(self.Type, client)
	resetAttrib.Value1--
	resetAttrib.Save()
	return true
}

///发货
func (self *QueryChallengeFightMsg) doAction(client IClient) bool {
	team := client.GetTeam() ///取得玩家队伍
	///取得挑战赛npc队伍
	challangeMatchStaticData := GetServer().GetStaticDataMgr().Unsafe().GetChallangeMatchType(self.MatchID) ///取挑战赛类型对象
	npcTeamTypeStaticData := GetServer().GetStaticDataMgr().Unsafe().GetNpcTeamType(challangeMatchStaticData.TeamID)
	if nil == npcTeamTypeStaticData {
		return false ///无效的npc球队类型,配置数据不存在
	}
	///计算比赛结果
	userGoalCount, npcGoalCount := npcTeamTypeStaticData.CalcMatchResult(team) ///计算比赛结果
	///计算收益
	awardIDList, awardNumList := CountChallangeMatchAward(userGoalCount, npcGoalCount, challangeMatchStaticData)
	///发放收益
	//awardResultIDList := IntList{}
	//awardResultNumList := IntList{}
	for i := 0; i < len(awardIDList); i++ {
		//if team.AwardObject(awardIDList[i], awardNumList[i], 0, 0) { ///获得道具，成功了就加到成功道具返回列表
		team.AwardObject(awardIDList[i], awardNumList[i], 0, 0) ///获得道具，无法判断成功与否，就直接返回数量把
		//awardResultIDList = append(awardResultIDList, awardIDList[i])
		//awardResultNumList = append(awardResultNumList, awardNumList[i])
		//} else { ///失败了要写0
		//	awardResultIDList = append(awardResultIDList, 0)
		//	awardResultNumList = append(awardResultNumList, 0)
		//}
	}
	///返回结果给客户端
	queryChallengeFightResultMsg := new(QueryChallengeFightResultMsg)
	queryChallengeFightResultMsg.NpcTeamID = challangeMatchStaticData.TeamID
	queryChallengeFightResultMsg.HomeGoal = userGoalCount
	queryChallengeFightResultMsg.GuestGoal = npcGoalCount
	//queryChallengeFightResultMsg.FixAddType = awardResultIDList[0]
	//queryChallengeFightResultMsg.AddList = awardResultIDList
	//queryChallengeFightResultMsg.AddNumberList = awardResultNumList
	queryChallengeFightResultMsg.FixAddType = awardIDList[0]
	queryChallengeFightResultMsg.AddList = awardIDList
	queryChallengeFightResultMsg.AddNumberList = awardNumList
	client.SendMsg(queryChallengeFightResultMsg) ///发送给客户端返回信息
	return true
}

type QueryChallengeFightResultMsg struct { ///返回挑战赛比赛结果
	MsgHead       `json:"head"` ///"challenge", "queryfightresult"
	NpcTeamID     int           `json:"npcteamid"`     ///挑战的npc球队id
	HomeGoal      int           `json:"homegoal"`      ///自己队进球数
	GuestGoal     int           `json:"guestgoal"`     ///目标队进球数
	FixAddType    int           `json:"fixaddtype"`    ///固定奖励类型
	AddList       IntList       `json:"addlist"`       ///奖励列表
	AddNumberList IntList       `json:"addnumberlist"` ///奖励数量列表
}

func (self *QueryChallengeFightResultMsg) GetTypeAndAction() (string, string) { ///实现消息处理接口的返回类型方法
	return "challenge", "queryfightresult"
}
