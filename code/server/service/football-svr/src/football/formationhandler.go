package football

import (
//	"fmt"
//"sort"
)

///阵型消息处理器
type FormationHandler struct {
	MsgHandler
}

///道具消息处理器
func (self *FormationHandler) getName() string { ///返回可处理的消息类型
	return "formation"
}

func (self *FormationHandler) initHandler() { ///初始化消息处理器
	self.addActionToList(new(FormationSetCurrentMsg))
	self.addActionToList(new(FormationUplevelMsg))
	self.addActionToList(new(FormationChangeStarMsg))
}

type FormationChangeStarMsg struct { ///请求换人消息
	MsgHead     `json:"head"` ///"formation", "changestar"
	FormationID int           `json:"formationid"` ///阵型id
	StarList    IntList       `json:"starlist"`    ///新的阵形球员列表pos1pos2.....pos11
}

func (self *FormationChangeStarMsg) GetTypeAndAction() (string, string) {
	return "formation", "changestar"
}

func (self *FormationChangeStarMsg) processAction(client IClient) bool {
	team := client.GetTeam()
	loger := GetServer().GetLoger()
	syncMgr := client.GetSyncMgr()
	starListLen := len(self.StarList)
	if starListLen != formationMinStarCount {
		return false ///阵形人数错误
	}
	formation := team.GetFormationMgr().GetFormation(self.FormationID)
	if nil == formation {
		return false ///球队没有此阵形
	}
	formationInfo := formation.GetInfo()
	posPtrList := IntPtrList{&formationInfo.Pos1, &formationInfo.Pos2, &formationInfo.Pos3, &formationInfo.Pos4,
		&formationInfo.Pos5, &formationInfo.Pos6, &formationInfo.Pos7, &formationInfo.Pos8,
		&formationInfo.Pos9, &formationInfo.Pos10, &formationInfo.Pos11} ///位置数据数组化
	if len(posPtrList) != starListLen {
		return false
	}

	isUnique := self.StarList.IsUnique()

	if loger.CheckFail("isUnique == true", isUnique == true,
		isUnique, true) {
		return false ///队伍存在一个以上相同的球星
	}

	//	fmt.Println(posPtrList)
	isChange := false
	for i := range self.StarList {
		starID := self.StarList[i]
		if *posPtrList[i] != starID {
			*posPtrList[i] = starID
			isChange = true
		}
	}
	//fmt.Println(posPtrList)
	if true == isChange {
		syncMgr.SyncObject("FormationChangeStarMsg", formation)
		team.CalcScore() ///更新球队角色卡
		client.GetSyncMgr().SyncObject("FormationSetCurrentMsg", team)
		syncMgr.syncAllStars()
	}
	return true
}

type FormationSetCurrentMsg struct { ///设置球队当前阵型
	MsgHead     `json:"head"` ///"formation", "setcurrentformation"
	FormationID int           `json:"formationid"` ///阵型id
	TacticType  int           `json:"tactictype"`  ///战术类型,非空为设置阵形战术,0表示设置当前阵形
	BuyTactic   bool          `json:"buytactic"`   ///true是花球票开启未达条件的战术,false不希望花球票
}

func (self *FormationSetCurrentMsg) GetTypeAndAction() (string, string) {
	return "formation", "setcurrentformation"
}

func (self *FormationSetCurrentMsg) buyTactic(client IClient) bool { ///通知球票购买战术
	team := client.GetTeam()
	staticDataMgr := GetServer().GetStaticDataMgr()
	formation := team.GetFormationMgr().GetFormation(self.FormationID)
	formationInfo := formation.GetInfo()
	staticData := staticDataMgr.GetStaticData(tableTacticType, self.TacticType)
	if nil == staticData {
		return false ///无此战术类型
	}
	tacticTypeData := staticData.(*TacticTypeStaticData)
	if formationInfo.Type != tacticTypeData.OpenFormType {
		return false ///战术与阵形类型不匹配
	}
	currentTicket := team.GetTicket()
	if currentTicket < tacticTypeData.OpenTicketPay {
		return false ///余额不足
	}
	tacticTypeStorePtrList := IntPtrList{&formationInfo.Tactic1, &formationInfo.Tactic2, &formationInfo.Tactic3}
	isChange := false
	for i := range tacticTypeStorePtrList {
		tacticTypeStorePtr := tacticTypeStorePtrList[i]
		//if *tacticTypeStorePtr > self.TacticType {
		//	break ///必须从低到高来买战术,根据战术类型id由小到大
		//}
		if *tacticTypeStorePtr == 0 { ///寻找空位填入战术值
			*tacticTypeStorePtr = self.TacticType
			isChange = true
			break
		}
	}
	if true == isChange {
		team.PayTicket(tacticTypeData.OpenTicketPay) ///消费
		client.SetMoneyRecord(PlayerCostMoney, Pay_BuyTactic, tacticTypeData.OpenTicketPay, team.GetTicket())
		client.GetSyncMgr().SyncObject("FormationSetCurrentMsg", team)
	}
	return true
}

func (self *FormationSetCurrentMsg) processAction(client IClient) bool {
	team := client.GetTeam()
	formation := team.GetFormationMgr().GetFormation(self.FormationID)
	if nil == formation {
		return false ///球队没有此阵形
	}
	if self.TacticType > 0 {
		if formation.HasTactic(self.TacticType) == false {
			if false == self.BuyTactic {
				return false ///非法设置并未开启的阵形战术
			}
			///希望用球票买战术
			if self.buyTactic(client) == false {
				return false ///购买战术失败
			}
		}
		///设置战术
		//		formationInfo := formation.GetInfo()
		//		formationInfo.CurrentTactic = self.TacticType
		formation.CurrentTactic = self.TacticType
		client.GetSyncMgr().SyncObject("FormationSetCurrentMsg", formation) ///同步阵形新属性到客户端
	}
	if self.FormationID > 0 && team.GetCurrentFormation() != self.FormationID {
		///设置当前阵形
		team.SetCurrentFormation(self.FormationID) ///设置球队当前阵型
		//		client.GetSyncMgr().SyncObject("FormationSetCurrentMsg", team) ///同步球队新属性到客户端
	}
	team.CalcScore() ///更新球队角色卡
	client.GetSyncMgr().SyncObject("FormationSetCurrentMsg", team)
	client.GetSyncMgr().syncMainStars()
	return true
}

type FormationUplevelMsg struct { ///请求球队阵型升级
	MsgHead `json:"head"` ///"formation", "formationuplevel"
}

func (self *FormationUplevelMsg) GetTypeAndAction() (string, string) {
	return "formation", "formationuplevel"
}

func (self *FormationUplevelMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	teamInfo := team.GetInfo()
	currentFormLevel := teamInfo.FormationLevel ///得到当前球队阵形等级
	staticDataMgr := GetServer().GetStaticDataMgr()
	needExp := staticDataMgr.GetLevelExpNeedExp(levelExpTypeFormation, currentFormLevel) ///得到升级所需战术点数
	if 0 == needExp {
		return false ///没有此等级经验配置
	}
	if teamInfo.TacticalPoint < needExp {
		return false ///战术点不足
	}
	levelExpCount := staticDataMgr.GetLevelExpCount(levelExpTypeFormation) ///得到球队阵形升级等级上限
	if currentFormLevel >= levelExpCount {
		client.SendErrorMsg(failFormationUplevel, failInreachmaxlevel)
		return false ///超等级上限或是无效的经验配置
	}
	return true
}

func (self *FormationUplevelMsg) payAction(client IClient) bool {
	///扣战术点
	team := client.GetTeam()
	teamInfo := team.GetInfo()
	currentFormLevel := teamInfo.FormationLevel ///得到当前球队阵形等级
	staticDataMgr := GetServer().GetStaticDataMgr()
	needExp := staticDataMgr.GetLevelExpNeedExp(levelExpTypeFormation, currentFormLevel) ///得到升级所需战术点数
	team.PayTacticalPoint(needExp)                                                       ///扣战术点
	client.GetSyncMgr().SyncObject("FormationUplevelMsg", team)                          ///同步战术点变化
	return true
}

func (self *FormationUplevelMsg) UpdateFormTactic(client IClient) bool { ///阵形等级提升后更新所有阵形所拥有的战术
	team := client.GetTeam()
	syncMgr := client.GetSyncMgr()
	formationList := team.GetFormationMgr().GetFormationList()
	teamFormLevel := team.FormationLevel ///得到球队阵形等级
	for i := range formationList {
		formation := formationList[i]
		formation.UpdateFormTactic(teamFormLevel)
		syncMgr.SyncObject("FormationUplevelMsg", formation)
		//formationInfo := formation.GetInfo()
		//		tacticTypeStorePtrList := []*int{&formationInfo.Tactic1, &formationInfo.Tactic2, &formationInfo.Tactic3}
		//tacticTypeList := staticDataMgr.GetTacticTypeList(formationInfo.Type, teamFormLevel)
		//if formationInfo.Tactic1 > 0 {
		//	tacticTypeList = append(tacticTypeList, formationInfo.Tactic1)
		//}
		//if formationInfo.Tactic2 > 0 {
		//	tacticTypeList = append(tacticTypeList, formationInfo.Tactic2)
		//}
		//if formationInfo.Tactic3 > 0 {
		//	tacticTypeList = append(tacticTypeList, formationInfo.Tactic3)
		//}
		//tacticTypeList = tacticTypeList.Unique()
		//if
		// if len(tacticTypeStorePtrList) != len(tacticTypeList) {
		// 	continue ///长度必须一致
		// }
		//sort.Ints(tacticTypeList) ///从小到大排序,低级的战术id小
		//isChange := false
		//for i := range tacticTypeList {
		//	tacticType := tacticTypeList[i]
		//	tacticTypeStorePtr := tacticTypeStorePtrList[i]
		//	if *tacticTypeStorePtr > 0 { ///判断是否有空位
		//		continue ///非零不更新,没空位了
		//	}
		//	*tacticTypeStorePtr = tacticType ///得到新的战术
		//	isChange = true
		//}
		//		if true == isChange { ///有值改变就同步给客户端
		//syncMgr.SyncObject("FormationUplevelMsg", formation)
		//		}
	}
	return true
}

func (self *FormationUplevelMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	team.AddFormationLevel(1) ///升一级
	///////////////////////////////////////////////////////////
	//DEBUG:  给予所有阵型
	// syncMgr := client.GetSyncMgr()
	// formationList := IntList{}
	// formationMgr := team.GetFormationMgr()
	// for i := 1; i < 9; i++ {
	// 	isHas := formationMgr.HasFormation(i)
	// 	if isHas == true {
	// 		continue
	// 	}
	// 	formationid := formationMgr.AwardFormation(i)
	// 	formationList = append(formationList, formationid)
	// }
	// syncMgr.syncAddFormation(formationList)

	////////////////////////////////////////////////////////////
	client.GetSyncMgr().SyncObject("FormationUplevelMsg", team)
	self.UpdateFormTactic(client)
	team.CalcScore() ///更新球队角色卡
	client.GetSyncMgr().SyncObject("FormationSetCurrentMsg", team)
	client.GetSyncMgr().syncMainStars()
	return true
}

func (self *FormationUplevelMsg) processAction(client IClient) bool {
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
