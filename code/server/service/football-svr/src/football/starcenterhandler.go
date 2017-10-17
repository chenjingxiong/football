package football

import (
	"fmt"
	//	"math/rand"
	"reflect"
	"strconv"
	"time"
)

type StarCenterHandler struct {
	MsgHandler
}

func (self *StarCenterHandler) getName() string { ///返回可处理的消息类型
	return "starcenter"
}

type StarCenterTransferMsg struct { ///请求球员中心转会球员消息,由客户端发起
	MsgHead          `json:"head"`
	StarCenterType   int     `json:"starcentertype"`   ///请求转会的球员所在球员中心的类型
	MemberID         int     `json:"memberid"`         ///请求转会的球员会员id
	ExchangeStarList IntList `json:"exchangestarlist"` ///用于抵偿花费的球员id列表
}

func (self *StarCenterTransferMsg) GetTypeAndAction() (string, string) {
	return "starcenter", "starcentertransfer"
}

type StarCenterTransferResultMsg struct { ///球员中心转会球员结果消息,由服务器回应
	MsgHead `json:"head"`
	Result  string `json:"result"` ///结果 ok为成功 or fail为失败
}

func NewStarCenterTransferResultMsg(result string) *StarCenterTransferResultMsg {
	msg := new(StarCenterTransferResultMsg)
	msg.Result = result
	return msg
}

func (self *StarCenterTransferResultMsg) GetTypeAndAction() (string, string) {
	return "starcenter", "starcentertransferresult"
}

func (self *StarCenterTransferMsg) getDeductPrice(client IClient) int { ///得到抵扣价格
	totalDeductPrice := 0
	team := client.GetTeam()
	staticDataMgr := GetServer().GetStaticDataMgr()
	for i := range self.ExchangeStarList {
		starID := self.ExchangeStarList[i] ///得到抵扣球员id
		if team.IsStarInCurrentFormation(starID) == true {
			continue ///首先阵形中的球员不准抵偿
		}
		star := team.GetStar(starID)
		if nil == star {
			continue ///球员不存在
		}

		starType := team.GetStarType(starID)
		if starType <= 0 {

		}
		starBasePrice := staticDataMgr.GetStarTypeBasePrice(starType) ////得到请求基础身价
		if starBasePrice <= 0 {
			continue ///不允许球员身价为0
		}
		totalDeductPrice += starBasePrice ///抵偿价格
	}
	return totalDeductPrice
}

func (self *StarCenterTransferMsg) checkExchangeLegal(targetClass int, exchangeClass int) bool {
	if targetClass == classE {
		/// 最低等级不用比较
		return true
	}
	return (targetClass - exchangeClass) < 50
}

func (self *StarCenterTransferMsg) checkAction(client IClient) bool { ///检测
	const exchangeStarListLenMax = 3 ///抵偿球员最多为3名
	team := client.GetTeam()
	now := int(time.Now().Unix())
	starCenter := client.GetTeam().GetStarCenter()
	staticDataMgr := GetServer().GetStaticDataMgr()
	exchangeStarListLen := len(self.ExchangeStarList)
	loger := GetServer().GetLoger()
	if loger.CheckFail("exchangeStarListLen<=exchangeStarListLenMax", ///禁止抵偿球员大于3名
		exchangeStarListLen <= exchangeStarListLenMax, exchangeStarListLen, exchangeStarListLenMax) {
		return false
	}
	starCenterMember := starCenter.GetStarCenterMember(self.StarCenterType, self.MemberID)
	if loger.CheckFail("starCenterMember!=nil", ///转会中心并不存在此球员
		starCenterMember != nil, starCenterMember, nil) {
		return false
	}
	if loger.CheckFail("now<starCenterMember.ExpireTime", ///球员已过期
		now < starCenterMember.ExpireTime, now, starCenterMember.ExpireTime) {
		return false
	}
	isTeamHasStar := team.HasStar(starCenterMember.StarType)
	if loger.CheckFail("isTeamHasStar==false", ///已拥有指定类型的球员,不准重复拥有
		isTeamHasStar == false, isTeamHasStar, false) {
		return false
	}
	needPayPrice := staticDataMgr.GetStarTypeBasePrice(starCenterMember.StarType) ///得到请求基础身价
	if loger.CheckFail("needPayPrice>0", needPayPrice > 0, needPayPrice, 0) {     ///不允许球员身价为0
		return false
	}

	targetStarClass := staticDataMgr.GetStarTypeClass(starCenterMember.StarType)
	for _, v := range self.ExchangeStarList {
		exchangeStar := team.GetStar(v)
		if loger.CheckFail("exchangeStar!=nil", exchangeStar != nil, exchangeStar, nil) { ///不允许球员身价为0
			return false
		}
		exchangeStarInfo := exchangeStar.GetInfo()
		exchangeStarClass := staticDataMgr.GetStarTypeClass(exchangeStarInfo.Type)
		isEvolveCountLegal := self.checkExchangeLegal(targetStarClass, exchangeStarClass)
		if loger.CheckFail("isEvolveCountLegal == true", isEvolveCountLegal == true,
			isEvolveCountLegal, true) {
			return false ///替换球员成长评价不能小于目标球员评价-1级以上
		}
	}

	deductPrice := self.getDeductPrice(client)      ///得到抵扣价格
	needPayPrice = Max(0, needPayPrice-deductPrice) ///不得减为负数
	currentCoin := team.GetCoin()                   ///扣钱
	if loger.CheckFail("currentCoin>=needPayPrice", ///钱不够扣
		currentCoin >= needPayPrice, currentCoin, needPayPrice) {
		return false
	}
	teamMaxStarCount := GetServer().GetStaticDataMgr().GetTeamMaxStarCount()
	teamMaxStarCount += team.AddStarPos
	exchangeStarCount := self.ExchangeStarList.Len()
	teamTotalStarCount := team.GetTotalStarCount()
	teamRemainStarCount := teamTotalStarCount + 1 - exchangeStarCount
	if loger.CheckFail("teamRemainStarCount<=teamMaxStarCount",
		teamRemainStarCount <= teamMaxStarCount, teamRemainStarCount, teamMaxStarCount) {
		return false
	}
	//isTeamStarFull := team.IsStarFull() ///判断球队是否已满
	//if loger.CheckFail("isTeamStarFull==false", isTeamStarFull == false, isTeamStarFull, false) {
	//	return false
	//}
	return true
}

func (self *StarCenterTransferMsg) payAction(client IClient) bool { ///支付
	team := client.GetTeam()
	starCenter := client.GetTeam().GetStarCenter()
	syncMgr := client.GetSyncMgr()
	staticDataMgr := GetServer().GetStaticDataMgr()
	starCenterMember := starCenter.GetStarCenterMember(self.StarCenterType, self.MemberID)
	needPayPrice := staticDataMgr.GetStarTypeBasePrice(starCenterMember.StarType) ////得到请求基础身价
	if self.ExchangeStarList.Len() > 0 {
		team.RemoveStar(self.ExchangeStarList)        ///扣球员
		syncMgr.syncRemoveStar(self.ExchangeStarList) ///同步球员删除
	}

	//每个球员减少50%  两个或以上则全免
	deductPrice := float32(self.ExchangeStarList.Len()) * 0.5

	if deductPrice < 1.0 {
		needPayPrice = int(float32(needPayPrice) * deductPrice)
	} else {
		needPayPrice = 0
	}

	//deductPrice := self.getDeductPrice(client)      ///得到抵扣价格
	//needPayPrice = Max(0, needPayPrice-deductPrice) ///不得减为负数
	if needPayPrice > 0 {
		team.PayCoin(needPayPrice)                        ///扣钱
		syncMgr.SyncObject("StarCenterTransferMsg", team) ///同步最新的球队信息给客户端
	}
	return true
}

func (self *StarCenterTransferMsg) doAction(client IClient) bool { ///发货
	starCenter := client.GetTeam().GetStarCenter()
	syncMgr := client.GetSyncMgr()
	loger := GetServer().GetLoger()
	starCenterMember := starCenter.GetStarCenterMember(self.StarCenterType, self.MemberID)
	starCenterMemberStarType := starCenterMember.StarType
	///扣球会中心中的球员会员
	starCenter.RemoveMember(self.StarCenterType, IntList{self.MemberID}) ///先从球员中心删除此成员
	starCenterMemberRemoveMsg := NewStarCenterMemberRemoveMsg(self.StarCenterType, self.MemberID)
	client.SendMsg(starCenterMemberRemoveMsg) ///告诉客户端从转会中心删除此会员
	///奖球员
	team := client.GetTeam()
	starID := team.AwardStar(starCenterMemberStarType)      ///给球队加此球员
	if loger.CheckFail("starID>0", starID > 0, starID, 0) { ///已拥有指定类型的球员,不准重复拥有
		return false
	}

	///给予球员默认星级
	star := team.GetStar(starID)
	starInfo := star.GetInfo()
	starInfo.EvolveCount = starCenter.GetStarEvolveCount(starCenterTypeDiscover, starCenterMemberStarType)
	syncMgr.syncAddStar(IntList{starID}) ///同步新加的球员到客户端
	return true
}

func (self *StarCenterTransferMsg) processAction(client IClient) (result bool) {
	defer func() {
		if false == result {
			self.sendStarCenterTransferResultMsg(client, msgResultFail) ///发失败结果消息
		} else {
			self.sendStarCenterTransferResultMsg(client, msgResultOK) ///发失败结果消息
		}
	}()
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

func (self *StarCenterTransferMsg) sendStarCenterTransferResultMsg(client IClient, result string) {
	msg := NewStarCenterTransferResultMsg(result)
	client.SendMsg(msg)
}

func (self *StarCenterHandler) initHandler() { ///初始化消息处理器
	self.addActionToList(new(QueryStarCenterMemberListMsg))
	self.addActionToList(new(StarCenterTransferMsg))
	self.addActionToList(new(QueryVolunteerInfoMsg))
	self.addActionToList(new(VolunteerSignMsg))
	self.addActionToList(new(GetStarLobbyAwardMsg))
}

///球员中心删除成员消息,有可能是客户端发送,也可能是服务器发送
type StarCenterMemberRemoveMsg struct {
	MsgHead        `json:"head"`
	StarCenterType int     `json:"starcentertype"` ///球员中心类型
	MemberIDList   IntList `json:"memberidlist"`   ///球员中心会员列表
}

func (self *StarCenterMemberRemoveMsg) GetTypeAndAction() (string, string) {
	return "starcenter", "starcentermemberremove"
}

func NewStarCenterMemberRemoveMsg(starCenterType int, memberID int) *StarCenterMemberRemoveMsg {
	msg := new(StarCenterMemberRemoveMsg)
	msg.StarCenterType = starCenterType
	msg.MemberIDList = append(msg.MemberIDList, memberID)
	return msg
}

type StarCenterMemberAddMsg struct { ///球员中心添加成员消息
	MsgHead              `json:"head"`
	StarCenterMemberList []StarCenterMember `json:"starcentermemberlist"` ///球员中心会员列表
}

type QueryStarCenterMemberListMsg struct { ///请求球员中心成员列表消息
	MsgHead        `json:"head"`
	StarCenterType int `json:"starcentertype"` ///球员中心类型
}

func (self *QueryStarCenterMemberListMsg) GetTypeAndAction() (string, string) {
	return "starcenter", "querystarcentermemberlist"
}

func (self *QueryStarCenterMemberListMsg) processAction(client IClient) bool {
	msg := NewStarCenterAddMemberMsg()
	starCenter := client.GetTeam().GetStarCenter()
	memberList := starCenter.GetStarCenterMemberList(self.StarCenterType)
	for i := range memberList {
		memberID := memberList[i]
		starCenterMember := starCenter.GetStarCenterMember(self.StarCenterType, memberID)
		if starCenterMember != nil {
			msg.AddMember(starCenterMember)
		}
	}
	client.SendMsg(msg) ///有可能是空消息
	return true
}

func (self *StarCenterMemberAddMsg) GetTypeAndAction() (string, string) {
	return "starcenter", "starcentermemberadd"
}

func NewStarCenterAddMemberMsg() *StarCenterMemberAddMsg {
	msg := new(StarCenterMemberAddMsg)
	msg.MsgType = "starcenter"
	msg.Action = "starcentermemberadd"
	return msg
}

func (self *StarCenterMemberAddMsg) AddMember(starCenterMember *StarCenterMember) {
	self.StarCenterMemberList = append(self.StarCenterMemberList, *starCenterMember)
}

///球员中心查询球员来投信息请求
type QueryVolunteerInfoMsg struct {
	MsgHead           `json:"head"` ///"starcenter", "queryvolunteerinfo"
	IsRefreshStarList bool          `json:"isrefreshstarlist"` ///是否用元宝刷新一次球员列表(修改: 招募卡)
}

func (self *QueryVolunteerInfoMsg) GetTypeAndAction() (string, string) {
	return "starcenter", "queryvolunteerinfo"
}

func (self *QueryVolunteerInfoMsg) createDefaultVolunteer(client IClient) bool { ///创建默认球星来投信息
	team := client.GetTeam()
	//loger := GetServer().GetLoger()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttribTeamVolunteerInfo := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamVolunteerInfo)
	if resetAttribTeamVolunteerInfo != nil {
		return true ///如果已存在球星来投信息则直接返回
	}
	//if loger.CheckFail("resetAttribTeamVolunteerInfo==nil", resetAttribTeamVolunteerInfo == nil,
	//	resetAttribTeamVolunteerInfo, nil) { ///创建默认球星来投信息默认数据必须为nil
	//	return false
	//}
	resetAttribMgr.AddResetAttrib(ResetAttribTypeTeamVolunteerInfo, 0, nil)
	resetAttribMgr.AddResetAttrib(ResetAttribTypeVolunteerStarType, 0, nil)
	return true
}

func (self *QueryVolunteerInfoMsg) updateVolunteerSystemInfo(client IClient) bool { ///更新球星来投系统信息
	team := client.GetTeam()
	loger := GetServer().GetLoger()
	staticDataMgr := GetServer().GetStaticDataMgr()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttribTeamVolunteerInfo := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamVolunteerInfo)
	if loger.CheckFail("resetAttribTeamVolunteerInfo!=nil", resetAttribTeamVolunteerInfo != nil,
		resetAttribTeamVolunteerInfo, nil) { ///更新球星来投系统信息时信息不能为空
		return false
	}
	isExpire := IsExpireTime(resetAttribTeamVolunteerInfo.ResetTime)
	if false == isExpire {
		return true ///未到更新时间不做处理
	}
	//expireTime := resetAttribTeamVolunteerInfo.ResetTime ///得到过期时间
	//now := int(time.Now().Unix())                        ///得到当前时间
	//if now <= expireTime {
	//	return true ///未到更新时间不做处理
	//}
	///处理系统信息刷新
	configStarVolunteer := staticDataMgr.GetConfigStaticData(configStarCenter, configItemStarVolunteer)
	if loger.CheckFail("configStarVolunteer!=nil", configStarVolunteer != nil,
		configStarVolunteer, nil) { ///判断球星来投配置数据有效性
		return false
	}
	//if nil == configStarVolunteer {
	//	return false
	//}
	defaultSignCount := 0 /// strconv.Atoi(configStarVolunteer.Param1) ///默认每日签约次数
	//	if loger.CheckFail("defaultSignCount>0", defaultSignCount > 0,
	//		defaultSignCount, 0) {
	//		return false
	//	}

	// vipLevel := team.GetVipLevel()
	// vipInfo := GetServer().GetStaticDataMgr().GetVipInfo(vipLevel)
	// if loger.CheckFail("vipInfo != nil", vipInfo != nil, vipInfo, nil) { ///默认每日补充刷新次数
	// 	return false
	// }

	//	defaultUpdateCount := vipInfo.Param11
	//	if loger.CheckFail("defaultUpdateCount>0", defaultUpdateCount > 0,
	//		defaultUpdateCount, 0) {
	//		return false
	//	}

	hasUpdateCount := 0 ///已刷新次数
	//remainUpdateCount := team.GetVipFreeVolunteerUpdateCount()
	//strconv.Atoi(configStarVolunteer.Param3) ///重置剩余刷新次数

	newExpireTime, _ := strconv.Atoi(configStarVolunteer.Param4) ///默认过期时间,需要计算
	if loger.CheckFail("newExpireTime>0", newExpireTime > 0, newExpireTime, 0) {
		return false
	}
	newExpireTime = GetHourUTC(newExpireTime)                                                 //int(time.Now().Unix()) + 3600*newExpireTime
	valueList := []int{defaultSignCount, hasUpdateCount, resetAttribTeamVolunteerInfo.Value3} ///生成属性列表
	ok := resetAttribMgr.UpdateResetAttrib(ResetAttribTypeTeamVolunteerInfo, newExpireTime, valueList)
	return ok
}

//得到随机星级
func GetRandomEvolove() int {

	///简化公式3^(7-m)/40  m = 星级 星级上限为7 m >= 4
	///概率 0.675 0.225 0.075 0.025
	iProbabilityList := []int{678, 225, 75, 25}
	nTemp := 0
	randNum := Random(0, 1000)
	randEvolove := 4
	for i := 4; i <= 7; i++ {
		if randNum >= nTemp && randNum < nTemp+iProbabilityList[i-4] {
			randEvolove = i
			break
		}
		nTemp += iProbabilityList[i-4]
	}
	return randEvolove
}

func (self *QueryVolunteerInfoMsg) updateVolunteerStarList(client IClient) bool { ///更新球星来投信息
	team := client.GetTeam()
	syncMgr := client.GetSyncMgr()
	starCenter := team.GetStarCenter()
	staticDataMgr := GetServer().GetStaticDataMgr()
	resetAttribMgr := team.GetResetAttribMgr()
	loger := GetServer().GetLoger()
	resetAttribTypeVolunteerStarType := resetAttribMgr.GetResetAttrib(ResetAttribTypeVolunteerStarType)
	//if loger.CheckFail("resetAttribTeamVolunteerInfo!=nil", resetAttribTeamVolunteerInfo != nil,
	//	resetAttribTeamVolunteerInfo, nil) { ///更新球星来投系统信息时信息不能为空
	//	return false
	//}
	if nil == resetAttribTypeVolunteerStarType { ///球星来投系统信息
		return false
	}
	resetAttribTeamVolunteerInfo := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamVolunteerInfo)
	if nil == resetAttribTeamVolunteerInfo { ///球星来投系统信息
		return false
	}
	expireTime := resetAttribTypeVolunteerStarType.ResetTime  ///得到过期时间
	now := int(time.Now().Unix())                             ///得到当前时间
	if now <= expireTime && false == self.IsRefreshStarList { ///时间未到或是非球票刷新均跳过
		return true ///未到更新时间不做处理
	}
	//	if true == self.IsRefreshStarList && resetAttribTeamVolunteerInfo.Value2 <= 0 {
	//		return false
	//	}
	configStarVolunteer := staticDataMgr.GetConfigStaticData(configStarCenter, configItemStarVolunteer)
	if nil == configStarVolunteer {
		return false
	}
	if true == self.IsRefreshStarList {
		///如果是球票刷球员则需要扣代价
		configVolunteerUpdateStarListPay, _ := strconv.Atoi(configStarVolunteer.Param5) ///每次球票刷新球员所支付的球票数
		if configVolunteerUpdateStarListPay <= 0 {
			return false ///不允许代价为0
		}
		//currentTicket := team.GetTicket()
		//if currentTicket < configVolunteerUpdateStarListPay {
		//	return false ///余额不足
		//}
		//currentTicket = team.PayTicket(configVolunteerUpdateStarListPay) ///支付球票
		//actionAttribChangeMsg := NewActionAttribChangeMsg(systemTypeStarVolunteer, teamTicketAttribType, 0, currentTicket)
		//client.SendMsg(actionAttribChangeMsg) ///同步客户端球票信息
		maxFreeUpdateCount := team.GetVipFreeVolunteerUpdateCount()
		if resetAttribTeamVolunteerInfo.Value2 < maxFreeUpdateCount { ///如果有免费次数则优先扣免费次数
			resetAttribTeamVolunteerInfo.Value2++ ///增加已刷新次数
		} else {
			itemMgr := team.GetItemMgr()
			itemIsEnough := itemMgr.HasEnoughItem(ItemCanvassCard, configVolunteerUpdateStarListPay)
			if loger.CheckFail("itemIsEnough == true", itemIsEnough == true, itemIsEnough, true) {
				return false ///道具数额不足
			}

			removeList, influenceItemID, _ := itemMgr.PayItemType(ItemCanvassCard, configVolunteerUpdateStarListPay)
			if removeList.Len() > 0 {
				syncMgr.SyncRemoveItem(removeList) ///同步道具删除
			}

			item := itemMgr.GetItem(influenceItemID)
			if item != nil {
				syncMgr.SyncObject("QueryVolunteerInfoMsg", item)
			}
		}
	}
	newExpireTime, _ := strconv.Atoi(configStarVolunteer.Param4) ///默认过期时间,需要计算
	if loger.CheckFail("newExpireTime>0", newExpireTime > 0, newExpireTime, 0) {
		return false
	}
	newExpireTime = GetHourUTC(newExpireTime) //int(time.Now().Unix()) + 3600*newExpireTime

	///处理球星列表刷新
	resultList := starCenter.RollVolunteerStarTypeList(client, 3, nil)
	if nil == resultList {
		return false
	}
	teamVolunteerInfo := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamVolunteerInfo)

	resultList = append(resultList, 0) ///扩展
	for v := 0; v < 3; v++ {
		randEvolove := GetRandomEvolove()

		//!潜规则
		if teamVolunteerInfo.Value3 == 0 && false == self.IsRefreshStarList { //! 第一次不要Roll出七星
			if randEvolove == 7 {
				randEvolove--
				teamVolunteerInfo.Value3++
				resetAttribMgr.UpdateResetAttrib(ResetAttribTypeTeamVolunteerInfo, newExpireTime, IntList{teamVolunteerInfo.Value1, teamVolunteerInfo.Value2, teamVolunteerInfo.Value3})
			}
		} else if teamVolunteerInfo.Value3 <= 1 && true == self.IsRefreshStarList { //!用道具刷新第一次必出七星
			randEvolove = 7
			teamVolunteerInfo.Value3 = 2
			resetAttribMgr.UpdateResetAttrib(ResetAttribTypeTeamVolunteerInfo, newExpireTime, IntList{teamVolunteerInfo.Value1, teamVolunteerInfo.Value2, teamVolunteerInfo.Value3})
		}

		resultList = append(resultList, randEvolove)
	}
	resultList = append(resultList, 0) ///扩展
	resetAttribMgr.UpdateResetAttrib(ResetAttribTypeVolunteerStarType, newExpireTime, resultList)
	return true
}

func (self *QueryVolunteerInfoMsg) processAction(client IClient) bool { ///查询球星来投信息
	team := client.GetTeam()
	staticDataMgr := GetServer().GetStaticDataMgr()
	resetAttribMgr := team.GetResetAttribMgr()
	configStarVolunteer := staticDataMgr.GetConfigStaticData(configStarCenter, configItemStarVolunteer)
	if nil == configStarVolunteer {
		return false
	}
	self.createDefaultVolunteer(client)    ///如果没有球星来投系统信息则创建默认的
	self.updateVolunteerSystemInfo(client) ///更新球星来投系统信息
	self.updateVolunteerStarList(client)   ///更新球星来投球星列表信息
	resetAttribTeamVolunteerInfo := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamVolunteerInfo)
	if nil == resetAttribTeamVolunteerInfo { ///球星来投系统信息
		return false
	}
	resetAttribTypeVolunteerStarType := resetAttribMgr.GetResetAttrib(ResetAttribTypeVolunteerStarType)
	if nil == resetAttribTypeVolunteerStarType { ///球星来投球星类型列表
		return false
	}
	///发送查询结果包
	queryVolunteerResultMsg := new(QueryVolunteerResultMsg)
	queryVolunteerResultMsg.RemainSignCount = resetAttribTeamVolunteerInfo.Value1
	queryVolunteerResultMsg.RemainUpdateCount = resetAttribTeamVolunteerInfo.Value2
	queryVolunteerResultMsg.StarTypeList = IntList{resetAttribTypeVolunteerStarType.Value1, resetAttribTypeVolunteerStarType.Value2,
		resetAttribTypeVolunteerStarType.Value3, resetAttribTypeVolunteerStarType.Value4}
	queryVolunteerResultMsg.StarEvolveList = IntList{resetAttribTypeVolunteerStarType.Value5, resetAttribTypeVolunteerStarType.Value6,
		resetAttribTypeVolunteerStarType.Value7, resetAttribTypeVolunteerStarType.Value8}
	client.SendMsg(queryVolunteerResultMsg)
	return true
}

type StarFateSignMsg struct { ///球星缘分系统签约请求消息
	MsgHead  `json:"head"` ///"starcenter", "starfatesign"
	StarType int           `json:"startype"` ///签约球星type
}

func (self *StarFateSignMsg) GetTypeAndAction() (string, string) {
	return "starcenter", "starfatesign"
}

func (self *StarFateSignMsg) checkAction(client IClient) bool {
	staticDataMgr := GetServer().GetStaticDataMgr()
	loger := GetServer().GetLoger()

	starTypeInfo := staticDataMgr.GetStarType(self.StarType)
	if loger.CheckFail("starTypeInfo != nil", starTypeInfo != nil, starTypeInfo, nil) {
		return false ///球星类型非法
	}

	team := client.GetTeam()
	itemMgr := team.GetItemMgr()

	isItemEnoughItem := itemMgr.HasEnoughItem(ItemStarCard, starTypeInfo.Ticket) ///判断是否有足够的球星卡
	if loger.CheckFail("isItemEnoughItem == true", isItemEnoughItem == true, isItemEnoughItem, true) {
		return false ///道具不足够或价格非法
	}

	return true
}

func (self *StarFateSignMsg) payAction(client IClient) bool {
	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	syncMgr := client.GetSyncMgr()
	// resetAttribMgr := team.GetResetAttribMgr()
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	starTypeInfo := staticDataMgr.GetStarType(self.StarType)

	removeList, influenceItemID, _ := itemMgr.PayItemType(ItemStarCard, starTypeInfo.Ticket)
	if removeList.Len() > 0 {
		syncMgr.SyncRemoveItem(removeList) ///同步道具删除
	}

	item := itemMgr.GetItem(influenceItemID)
	if item != nil {
		syncMgr.SyncObject("StarFateSignMsg", item)
	}

	return true
}

func (self *StarFateSignMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	//evolveCount := calcVolunteerStarCount(client)
	team.AwardObject(0, 0, 1, self.StarType)
	return true
}

func (self *StarFateSignMsg) processAction(client IClient) bool {

	msgResult := msgResultFail
	defer self.sendResultMsg(client, &msgResult)
	if self.checkAction(client) == false {
		return false
	}

	if self.payAction(client) == false {
		return false
	}

	if self.doAction(client) == false {
		return false
	}

	msgResult = msgResultOK

	return true
}

type FateSignResultMsg struct { ///球星来投签约球员请求结果消息
	MsgHead `json:"head"` ///"starcenter", "getstaryuanfen"
	Result  string        `json:"result"` ///签约结果 ok or fail
}

func (self *FateSignResultMsg) GetTypeAndAction() (string, string) {
	return "starcenter", "getstaryuanfen"
}

func (self *StarFateSignMsg) sendResultMsg(client IClient, result *string) {
	msg := new(FateSignResultMsg)
	msg.Result = *result
	client.SendMsg(msg)
}

type VolunteerSignMsg struct { ///球星来投签约球员请求消息
	MsgHead  `json:"head"` ///"starcenter", "volunteersign"
	StarType int           `json:"startype"` ///签约球星type 服务器要验证type是否拥有
}

func (self *VolunteerSignMsg) GetTypeAndAction() (string, string) {
	return "starcenter", "volunteersign"
}

type QueryVolunteerResultMsg struct { ///查询球星来投信息结果
	MsgHead           `json:"head"`
	RemainSignCount   int     `json:"remainsigncount"`   ///剩余签约次数
	RemainUpdateCount int     `json:"remainupdatecount"` ///剩余刷新次数
	StarTypeList      IntList `json:"startypelist"`      ///球星类型列表
	StarEvolveList    IntList `json:"starevolvelist"`    //球星星级列表
}

func (self *QueryVolunteerResultMsg) GetTypeAndAction() (string, string) {
	return "starcenter", "queryvolunteerresult"
}

func (self *VolunteerSignMsg) checkAction(client IClient) bool { ///检测
	///判断签约剩余数是否够足
	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	loger := GetServer().GetLoger()
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	resetAttribMgr := team.GetResetAttribMgr()
	if team.IsStarFull() == true {
		return false ///球队球员数已满
	}
	//if team.GetStarFromType(self.StarType) != nil {
	//	return false ///球队里已有此球员了
	//}
	resetAttribTeamVolunteerInfo := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamVolunteerInfo)
	if nil == resetAttribTeamVolunteerInfo {
		return false
	}
	// if resetAttribTeamVolunteerInfo.Value1 <= 0 {
	// 	return false ///剩余签约数为0
	// }
	starTypeInfo := staticDataMgr.GetStarType(self.StarType)
	if loger.CheckFail("starTypeInfo != nil", starTypeInfo != nil, starTypeInfo, nil) {
		return false ///球星类型非法
	}
	//currentTicket := team.GetTicket()
	//if currentTicket < starTypeInfo.Ticket || starTypeInfo.Ticket <= 0 {
	//	return false ///球票不足够扣或价格非法
	//}
	//configStarVolunteer := GetServer().GetStaticDataMgr().GetConfigStaticData(configStarCenter, configItemStarVolunteer)
	//if loger.CheckFail("configStarVolunteer != nil", configStarVolunteer != nil, configStarVolunteer, nil) {
	//	return false
	//}

	//configVolunteerUpdateStarListPay, _ := strconv.Atoi(configStarVolunteer.Param5) ///每次球票刷新球员所支付的球票数
	//if configVolunteerUpdateStarListPay <= 0 {
	//	return false ///不允许代价为0
	//}

	isItemEnoughItem := itemMgr.HasEnoughItem(ItemStarCard, starTypeInfo.Ticket) ///判断是否有足够的球星卡
	if loger.CheckFail("isItemEnoughItem == true", isItemEnoughItem == true, isItemEnoughItem, true) {
		return false ///道具不足够或价格非法
	}

	///验证客户端转来starType合法性
	foundStarType := 0
	resetAttribTypeVolunteerStarType := resetAttribMgr.GetResetAttrib(ResetAttribTypeVolunteerStarType)
	if nil == resetAttribTeamVolunteerInfo {
		return false
	}
	valueType := reflect.ValueOf(resetAttribTypeVolunteerStarType).Elem()
	for i := 1; i < 10; i++ {
		fieldName := fmt.Sprintf("Value%d", i)
		fieldValue := valueType.FieldByName(fieldName)
		if fieldValue.IsValid() == false {
			break
		}
		foundStarType = int(fieldValue.Int())
		if foundStarType == self.StarType {
			break
		}
	}
	if foundStarType != self.StarType {
		return false ///非法的球员类型
	}
	return true
}

func (self *VolunteerSignMsg) replaceStarType(client IClient) { ///支付
	team := client.GetTeam()
	starCenter := team.GetStarCenter()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttribTypeVolunteerStarType := resetAttribMgr.GetResetAttrib(ResetAttribTypeVolunteerStarType)
	valueType := reflect.ValueOf(resetAttribTypeVolunteerStarType).Elem()
	newStarType := int64(0)
	for i := 1; i < 10; i++ {
		fieldName := fmt.Sprintf("Value%d", i)
		fieldValue := valueType.FieldByName(fieldName)
		if fieldValue.IsValid() == false {
			break
		}
		starType := int(fieldValue.Int())
		if starType == self.StarType {
			excludeList := IntList{resetAttribTypeVolunteerStarType.Value1, resetAttribTypeVolunteerStarType.Value2,
				resetAttribTypeVolunteerStarType.Value3, resetAttribTypeVolunteerStarType.Value4}
			newStarTypeList := starCenter.RollVolunteerStarTypeList(client, 1, excludeList)
			if len(newStarTypeList) > 0 {
				///随机品质
				gradefieldName := fmt.Sprintf("Value%d", i+4)
				gradefieldValue := valueType.FieldByName(gradefieldName)
				newGrade := calcVolunteerStarCount(client)
				gradefieldValue.SetInt(int64(newGrade))

				newStarType = int64(newStarTypeList[0])
			}
			fieldValue.SetInt(newStarType)
			break
		}
	}
}

func (self *VolunteerSignMsg) payAction(client IClient) bool { ///支付
	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	syncMgr := client.GetSyncMgr()
	// resetAttribMgr := team.GetResetAttribMgr()
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	starTypeInfo := staticDataMgr.GetStarType(self.StarType)
	//currentTicket := team.PayTicket(starTypeInfo.Ticket) ///支付球员球票价格
	//syncMgr.syncAttribChangeItem(systemTypeStarVolunteer, teamTicketAttribType, 0, currentTicket)
	removeList, influenceItemID, _ := itemMgr.PayItemType(ItemStarCard, starTypeInfo.Ticket)
	if removeList.Len() > 0 {
		syncMgr.SyncRemoveItem(removeList) ///同步道具删除
	}

	item := itemMgr.GetItem(influenceItemID)
	if item != nil {
		syncMgr.SyncObject("VolunteerSignMsg", item)
	}

	///扣签约次数 (删除)
	//resetAttribTeamVolunteerInfo := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamVolunteerInfo)
	//resetAttribTeamVolunteerInfo.Value1-- ///减一次签约次数
	//resetAttribTeamVolunteerInfo.Save()
	//resetAttribTeamVolunteerInfo.Save()
	return true
}

func calcVolunteerStarCount(client IClient) int { ///计算球星来投星级
	///简化公式3^(7-m)/40  m = 星级 星级上限为7 m >= 4
	///概率 0.675 0.225 0.75 0.25
	iProbabilityList := []int{678, 225, 75, 25}
	randNum := Random(0, 1000)
	nTemp := 0
	evolveCount := 4
	for i := 4; i <= 7; i++ {
		if randNum >= nTemp && randNum < nTemp+iProbabilityList[i-4] {
			evolveCount = i
			break
		}
	}
	return evolveCount
}

func (self *VolunteerSignMsg) doAction(client IClient) bool { ///发货
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	//awardStarID := team.AwardStar(self.StarType)
	//if awardStarID <= 0 {
	//	return false
	//}
	volunteerStarType := resetAttribMgr.GetResetAttrib(ResetAttribTypeVolunteerStarType)
	if nil == volunteerStarType {
		return false
	}
	//star := team.GetStar(awardStarID)
	//starInfo := star.GetInfo()
	evolveCount := 0
	starList := IntList{volunteerStarType.Value1, volunteerStarType.Value2, volunteerStarType.Value3, volunteerStarType.Value4}
	StarEvolveList := IntList{volunteerStarType.Value5, volunteerStarType.Value6, volunteerStarType.Value7, volunteerStarType.Value8}
	for i, v := range starList {
		if v == self.StarType {
			evolveCount = StarEvolveList[i]
		}
	}

	team.AwardObject(0, 0, evolveCount, self.StarType)
	///替换已签约的starType为一个新starType
	self.replaceStarType(client)
	//star.SetPermanentContract() ///设置球员为永久球员
	//syncMgr := client.GetSyncMgr()
	//syncMgr.syncAddStar(IntList{awardStarID}) ///通知客户端球队新进了一名球员
	///发送查询结果包
	//	resetAttribMgr := team.GetResetAttribMgr()
	resetAttribTeamVolunteerInfo := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamVolunteerInfo)
	resetAttribTypeVolunteerStarType := resetAttribMgr.GetResetAttrib(ResetAttribTypeVolunteerStarType)
	queryVolunteerResultMsg := new(QueryVolunteerResultMsg)
	queryVolunteerResultMsg.RemainSignCount = resetAttribTeamVolunteerInfo.Value1
	queryVolunteerResultMsg.RemainUpdateCount = resetAttribTeamVolunteerInfo.Value2
	queryVolunteerResultMsg.StarTypeList = IntList{resetAttribTypeVolunteerStarType.Value1, resetAttribTypeVolunteerStarType.Value2,
		resetAttribTypeVolunteerStarType.Value3, resetAttribTypeVolunteerStarType.Value4}
	queryVolunteerResultMsg.StarEvolveList = IntList{resetAttribTypeVolunteerStarType.Value5, resetAttribTypeVolunteerStarType.Value6,
		resetAttribTypeVolunteerStarType.Value7, resetAttribTypeVolunteerStarType.Value8}
	client.SendMsg(queryVolunteerResultMsg)
	return true
}

func (self *VolunteerSignMsg) processAction(client IClient) bool {
	msgResult := msgResultFail
	defer self.sendResultMsg(client, &msgResult)
	if self.checkAction(client) == false { ///检测
		return false
	}
	if self.payAction(client) == false { ///支付
		return false
	}
	if self.doAction(client) == false { ///发货
		return false
	}
	msgResult = msgResultOK
	return true
}

func (self *VolunteerSignMsg) sendResultMsg(client IClient, result *string) {
	msg := new(VolunteerSignResultMsg)
	msg.Result = *result
	client.SendMsg(msg)
}

type VolunteerSignResultMsg struct { ///球星来投签约球员请求结果消息
	MsgHead `json:"head"` ///"starcenter", "volunteersignresult"
	Result  string        `json:"result"` ///签约结果 ok or fail
}

func (self *VolunteerSignResultMsg) GetTypeAndAction() (string, string) {
	return "starcenter", "volunteersignresult"
}

type GetStarLobbyAwardMsg struct { ///激活游说球星 来自客户端请求
	MsgHead     `json:"head"` ///"starcenter", "getstarlobby"
	StarLobbyID int           `json:"starlobbyid"`
}

func (self *GetStarLobbyAwardMsg) GetTypeAndAction() (string, string) {
	return "starcenter", "getstarlobby"
}

func (self *GetStarLobbyAwardMsg) IsStarLobbyCondition(client IClient) bool {
	///获取达成奖励球员所需的条件球员信息
	lobbyType := GetServer().GetStaticDataMgr().Unsafe().GetStarLobbyType(self.StarLobbyID)
	if nil == lobbyType {
		return false
	}

	needStarList := IntList{lobbyType.Needstartype1, lobbyType.Needstartype2, lobbyType.Needstartype3,
		lobbyType.Needstartype4, lobbyType.Needstartype5}

	///判断条件是否达成
	team := client.GetTeam()
	for i := range needStarList {
		if needStarList[i] == 0 {
			continue
		}

		if team.GetStarFromType(needStarList[i]) == nil {
			return false
		}
	}
	return true
}

func (self *GetStarLobbyAwardMsg) IsLobbyStarExist(client IClient) bool {
	///奖励球星是否存在于玩家队伍
	team := client.GetTeam()
	lobbyType := GetServer().GetStaticDataMgr().Unsafe().GetStarLobbyType(self.StarLobbyID)
	if lobbyType == nil {
		return false
	}
	if team.GetStarFromType(lobbyType.Awardstartype) == nil {
		return false
	}
	return true
}

func (self *GetStarLobbyAwardMsg) checkAction(client IClient) bool {

	loger := GetServer().GetLoger()

	///验证游说ID合法性
	IsloddyIDLegal := (self.StarLobbyID > starLobbyLimitMax || self.StarLobbyID <= 0)
	if loger.CheckFail("IsloddyIDLegal == true", IsloddyIDLegal == false,
		IsloddyIDLegal, false) {
		return false
	}

	///判断奖励球星所需条件是否具备
	lobbyType := GetServer().GetStaticDataMgr().Unsafe().GetStarLobbyType(self.StarLobbyID)
	if loger.CheckFail("lobbyType != nil", lobbyType != nil,
		lobbyType, nil) {
		return false
	}

	needStarExist := self.IsStarLobbyCondition(client)
	if loger.CheckFail("needStarExist == true", needStarExist == true,
		needStarExist, true) {
		return false
	}

	///判断奖励球星是否已存在于玩家队伍
	awardStarExist := self.IsLobbyStarExist(client)
	if loger.CheckFail("awardStarExist == false", awardStarExist == false,
		awardStarExist, false) {
		return false
	}

	///判断奖励球星是否已存在于转会中心
	// starCenter := client.GetTeam().GetStarCenter()
	// starCenterExist := starCenter.IsTypeExistStarCenter(starCenterTypeLobby, lobbyType.Awardstartype)
	// if loger.CheckFail("starCenterExist == false", starCenterExist == false,
	// 	starCenterExist, false) {
	// 	return false
	// }

	///检测当前发掘球员中心的会员数是否已到达上限
	// payConfig := GetServer().GetStaticDataMgr().GetConfigStaticData(configStarSpy, configItemDiscoverConfig)
	// maxMemberCount, _ := strconv.Atoi(payConfig.Param1) ///得到发掘球员中心的会员数上限数

	///得到当前发掘球员中心的会员数
	// currentMemberCount := client.GetTeam().GetStarCenter().GetStarCenterMemberCount(starCenterTypeDiscover)
	// if loger.CheckFail("currentMemberCount < maxMemberCount", currentMemberCount < maxMemberCount,
	// 	currentMemberCount, maxMemberCount) {
	// 	return false
	// }

	return true
}

func (self *GetStarLobbyAwardMsg) doAction(client IClient) bool {
	///设置球星游说Mask
	lobbyType := GetServer().GetStaticDataMgr().Unsafe().GetStarLobbyType(self.StarLobbyID)

	client.GetTeam().SetLobbyMask(self.StarLobbyID - 1)

	/// 将奖励球员放入转会中心 改: 球员直接入队
	team := client.GetTeam()
	// team.AwardStar(lobbyType.Awardstartype)

	team.AwardObject(0, 0, 1, lobbyType.Awardstartype)
	return true
}

type GetStarLobbyAwardResultMsg struct { ///球星游说的结果
	MsgHead `json:"head"` ///"starcenter", "getstarlobbyresult"
	Result  string        `json:"result"` ///签订结果 ok or fail
}

func (self *GetStarLobbyAwardResultMsg) GetTypeAndAction() (string, string) {
	return "starcenter", "getstarlobbyresult"
}

func (self *GetStarLobbyAwardMsg) sendResultMsg(client IClient, result *string) {
	msg := new(GetStarLobbyAwardResultMsg)
	msg.Result = *result
	client.SendMsg(msg)
}

func (self *GetStarLobbyAwardMsg) processAction(client IClient) bool {
	msgResult := msgResultFail
	defer self.sendResultMsg(client, &msgResult)
	if self.checkAction(client) == false { ///检测
		return false
	}
	//if self.payAction(client) == false { ///支付
	//	return false
	//}
	if self.doAction(client) == false { ///发货
		return false
	}
	msgResult = msgResultOK
	return true
}
