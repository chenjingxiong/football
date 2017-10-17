package football

import (
	//	"math/rand"
	"strconv"

	//	"time"
)

const (
	drawGroupStarSpyDiscoverPrimary = 1 ///初级球探发掘球员所有抽卡分组号
	drawGroupStarSpyDiscoverMiddle  = 2 ///中级球探发掘球员所有抽卡分组号
	drawGroupStarSpyDiscoverExpert  = 3 ///高级球探发掘球员所有抽卡分组号
)

type StarSpyHandler struct {
	MsgHandler
}

func (self *StarSpyHandler) getName() string { ///返回可处理的消息类型
	return "starspy"
}

func (self *StarSpyHandler) initHandler() { ///初始化消息处理器
	self.addActionToList(new(StarSpyDiscoverMsg))
	self.addActionToList(new(StarSpyOperateMsg))
}

const (
	discoverResultOK   = "ok"   ///普通成功
	discoverResultLuck = "luck" ///幸运成功
	discoverResultFail = "fail" ///失败
)

type StarSpyDiscoverResultMsg struct { ///球探发掘球员结果消息
	MsgHead                       `json:"head"`
	Result                        string  `json:"result"`                        ///成功为ok,失败为fail
	DiscoverResultStarType        int     `json:"discoverresultstartype"`        ///成功发掘球员类型
	DiscoverStarList              IntList `json:"discoverstarlist"`              ///发掘候选球员列表
	DiscoverResultStarEvolveCount int     `json:"discoverresultstarevolvecount"` ///成功发掘球员星级
}

func (self *StarSpyDiscoverResultMsg) GetTypeAndAction() (string, string) {
	return "starspy", "discoverresult"
}

func newStarSpyDiscoverResultMsg(result string) *StarSpyDiscoverResultMsg {
	msg := new(StarSpyDiscoverResultMsg)
	msg.MsgType = "starspy"
	msg.Action = "discoverresult"
	msg.Result = result
	return msg
}

type StarSpyDiscoverMsg struct { ///请求球探发掘球员消息
	MsgHead     `json:"head"`
	StarSpyType int `json:"starspytype"` ///球探类型 1初级 2中级  3高级
}

func (self *StarSpyDiscoverMsg) GetTypeAndAction() (string, string) {
	return "starspy", "starspydiscover"
}

///根据不同的球探类型返回相应的发掘球员后幸运值加成
func (self *StarSpyDiscoverMsg) getConfigStarSpyDiscoverLuckAward() int {
	result := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configStarSpy, configItemDiscoverAddLuck, self.StarSpyType)
	// paramIntMap := GetServer().GetStaticDataMgr().getConfigStaticDataParamIntMap(configStarSpy, configItemDiscoverAddLuck)
	// result, ok := (*paramIntMap)[self.StarSpyType]
	// if false == ok {
	// 	return 0
	// }
	return result
}

///根据不同的球探类型返回相应的发掘球员配置数据
func (self *StarSpyDiscoverMsg) getConfigStarSpyDiscoverParam(starSpyType int, subType string) int {
	paramIntMap := GetServer().GetStaticDataMgr().getConfigStaticDataParamIntMap(configStarSpy, subType)
	result, ok := (*paramIntMap)[starSpyType]
	if false == ok {
		return 0
	}
	return result
}

func (self *StarSpyDiscoverMsg) getConfigMaxDiscoverStarNum(starSpyType int) int { ///根据不同的球探类型返回相应的发掘球员代价
	//result := 0
	payConfig := GetServer().GetStaticDataMgr().GetConfigStaticData(configStarSpy, configItemDiscoverExtraPay)
	if nil == payConfig {
		return 0
	}
	result, err := strconv.Atoi(payConfig.Param1)
	if err != nil {
		return 0
	}
	return result
}

func (self *StarSpyDiscoverMsg) getConfigDiscoverRate(starSpyType int) []int { ///根据不同的球探类型返回相应的发掘球员几率
	result := []int{}
	configDiscoverRateDic := map[int]string{primerStarSpy: configItemDiscoverPrimaryRate,
		middleStarSpy: configItemDiscoverMiddleRate,
		expertStarSpy: configItemDiscoverExpertRate} ///发掘机率对应表
	configItemDiscoverParam, ok := configDiscoverRateDic[starSpyType]
	if true == ok {
		result = GetServer().GetStaticDataMgr().getConfigStaticDataParamIntList(configStarSpy, configItemDiscoverParam)
	}
	return result
}

func (self *StarSpyDiscoverMsg) getConfigDiscoverPay() int { ///根据不同的球探类型返回相应的发掘球员代价
	result := 0
	payConfig := GetServer().GetStaticDataMgr().GetConfigStaticData(configStarSpy, configItemDiscoverExtraPay)
	if nil == payConfig {
		return 0
	}
	switch self.StarSpyType {
	case primerStarSpy:
		result, _ = strconv.Atoi(payConfig.Param1)
	case middleStarSpy:
		result, _ = strconv.Atoi(payConfig.Param2)
	case expertStarSpy:
		result, _ = strconv.Atoi(payConfig.Param3)
	}
	return result
}

func (self *StarSpyDiscoverMsg) checkAction(client IClient) bool { ///检测消息合法性
	if nil == client.GetTeam() { ///检测客户端是是否已经创建过球队了
		GetServer().GetLoger().Warn("StarSpyDiscoverMsg CheckAction fail! client's team is nil!")
		client.SendErrorMsg(failStarSpyDiscover, failInvalidMsg)
		return false
	}
	if self.StarSpyType < primerStarSpy || self.StarSpyType > expertStarSpy { ///客户端发来无效的球探类型
		GetServer().GetLoger().Warn("StarSpyDiscoverMsg CheckAction fail! StarSpyType is Invalid!")
		client.SendErrorMsg(failStarSpyDiscover, failInvalidParam)
		return false
	}
	///检测抽选组内是否满足5张
	starCenter := client.GetTeam().GetStarCenter()
	staticDataMgr := GetServer().GetStaticDataMgr()
	loger := GetServer().GetLoger()
	drawGroupList := staticDataMgr.GetDrawGroupIndexList(1) ///取得抽卡索引列表复本(修改: 废除Group2-6, 球探共用Group1)
	if loger.CheckFail("drawGroupList!=nil", drawGroupList != nil, drawGroupList, nil) {
		return false
	}
	///过滤掉已获得员球列表
	const maxDrawCount = 5 ///最多抽5张卡
	drawGroupList = starCenter.GetDrawGroupFilterIndexList(starCenterTypeDiscover, drawGroupList)
	drawGroupListLen := len(drawGroupList)
	if loger.CheckFail("drawGroupListLen >= maxDrawCount", drawGroupListLen >= maxDrawCount, drawGroupListLen, maxDrawCount) {
		return false
	}
	///检测当前发掘球员中心的会员数是否已到达上限
	payConfig := GetServer().GetStaticDataMgr().GetConfigStaticData(configStarSpy, configItemDiscoverConfig)
	maxMemberCount, _ := strconv.Atoi(payConfig.Param1) ///得到发掘球员中心的会员数上限数
	///得到当前发掘球员中心的会员数
	currentMemberCount := client.GetTeam().GetStarCenter().GetStarCenterMemberCount(starCenterTypeDiscover)
	if currentMemberCount >= maxMemberCount { ///判断会员中心的成员数是否超上限
		GetServer().GetLoger().Warn("StarSpyDiscoverMsg CheckAction fail! currentMemberCount>=maxMemberCount!")
		client.SendErrorMsg(failStarSpyDiscover, failReachLimit)
		return false
	}
	starSpy := client.GetTeam().GetStarSpy()                ///取得球探接口
	spyNeedItem := starSpy.GetSpyNeedItem(self.StarSpyType) ///取得球探所需道具
	if loger.CheckFail("spyNeedItem != 0", spyNeedItem != 0, spyNeedItem, 0) {
		return false
	}

	//configDiscoverPay := self.getConfigDiscoverPay() ///取得球探在cd状态和没有剩余次数时发掘球员所支付的球票代价
	//if loger.CheckFail("configDiscoverPay >0", configDiscoverPay > 0, configDiscoverPay, 0) {
	//	return false
	//}

	//discoverCD := starSpy.GetDiscoverCD(self.StarSpyType) ///得到当前球探的cd时间
	//currentTeamTicket := client.GetTeam().GetInfo().Ticket       ///得到当前球队球票的余额
	discoverRemainCount := starSpy.GetDiscoverRemainCount(self.StarSpyType)
	if discoverRemainCount <= 0 {
		itemMgr := client.GetTeam().GetItemMgr()
		isEnough := itemMgr.HasEnoughItem(spyNeedItem, 1) ///有cd时判断道具是否足够
		if loger.CheckFail("isEnough == true", isEnough == true, isEnough, 1) {
			return false
		}
	}

	//if discoverCD > 0 && currentTeamTicket < configDiscoverPay { ///有cd时判断余额是否足够
	//	GetServer().GetLoger().Warn("StarSpyDiscoverMsg CheckAction fail! insufficient ticket!")
	//	client.SendErrorMsg(failStarSpyDiscover, failInsufficientTicket)
	//	return false
	//}

	//	if discoverCD <= 0 && discoverRemainCount <= 0 { ///没cd时需要检测剩于次数要大于1,可扣
	//		GetServer().GetLoger().Warn("StarSpyDiscoverMsg CheckAction fail! insufficient discover count!")
	//		client.SendErrorMsg(failStarSpyDiscover, failInsufficientDiscoverCount)
	//		return false
	//	}
	//	configDiscoverCD := self.getConfigStarSpyDiscoverParam(self.StarSpyType, configItemDiscoverCD) ///得到配置中的球探发掘cd时间
	//	if loger.CheckFail("configDiscoverCD >0", configDiscoverCD > 0, configDiscoverCD, 0) {         ///验证配置中的球探发掘cd时间必须合法
	//		return false
	//	}
	return true
}

func (self *StarSpyDiscoverMsg) payAction(client IClient) bool { ///支付代价逻辑
	starSpy := client.GetTeam().GetStarSpy() ///取得球探接口
	//	configDiscoverPay := self.getConfigDiscoverPay()             ///取得球探在cd状态和没有剩余次数时发掘球员所支付的球票代价
	//currentDiscoverCD := starSpy.GetDiscoverCD(self.StarSpyType) ///得到当前球探的CD时间
	discoverRemainCount := starSpy.GetDiscoverRemainCount(self.StarSpyType)
	if discoverRemainCount <= 0 {
		/////处于有cd状态时需要扣球队的球票
		//ok := client.GetTeam().SpendTicket(configDiscoverPay, "StarSpyDiscover")
		//currentTicket := client.GetTeam().GetInfo().Ticket
		//msgTicketChange := NewActionAttribChangeMsg(systemTypeStarSpyDiscover, teamTicketAttribType, 0, currentTicket) ///生成一个属性变更消息
		//client.SendMsg(msgTicketChange)

		///处于有cd状态时需要扣除玩家道具
		team := client.GetTeam()
		itemMgr := team.GetItemMgr()
		syncMgr := client.GetSyncMgr()
		spyNeedItem := starSpy.GetSpyNeedItem(self.StarSpyType) ///取得球探所需道具
		removeList, influenceItemID, _ := itemMgr.PayItemType(spyNeedItem, 1)
		if removeList.Len() > 0 {
			syncMgr.SyncRemoveItem(removeList) ///同步道具删除
		}

		item := itemMgr.GetItem(influenceItemID)
		if item != nil {
			syncMgr.SyncObject("StarSpyDiscoverMsg", item)
		}

		return true
	}
	payConfig := GetServer().GetStaticDataMgr().GetConfigStaticData(configStarSpy, configItemDiscoverResetCount)
	///初始化为vip0的默认次数
	maxRemainFreeCountList := IntList{0, 0, 0}
	maxRemainFreeCountList[0], _ = strconv.Atoi(payConfig.Param2)
	maxRemainFreeCountList[1], _ = strconv.Atoi(payConfig.Param4)
	maxRemainFreeCountList[2], _ = strconv.Atoi(payConfig.Param6)

	freeCountResetTimeList := IntList{0, 0, 0}
	freeCountResetTimeList[0], _ = strconv.Atoi(payConfig.Param1)
	freeCountResetTimeList[1], _ = strconv.Atoi(payConfig.Param3)
	freeCountResetTimeList[2], _ = strconv.Atoi(payConfig.Param5)

	team := client.GetTeam()
	vipLevel := team.GetVipLevel()
	vipPrivilege := GetServer().GetStaticDataMgr().GetVipInfo(vipLevel)
	if vipPrivilege != nil {
		maxRemainFreeCountList[0] = vipPrivilege.Param6
		maxRemainFreeCountList[1] = vipPrivilege.Param7
		maxRemainFreeCountList[2] = vipPrivilege.Param8
	}
	maxRemainFreeCount := maxRemainFreeCountList[self.StarSpyType-1]
	if discoverRemainCount >= maxRemainFreeCount { ///在次数上限消费需要将cd置0
		freeCountResetTime := freeCountResetTimeList[self.StarSpyType-1]
		starSpy.SetResetRemainCD(self.StarSpyType, Now()+freeCountResetTime)
	}
	discoverRemainCountFinal := discoverRemainCount - 1                        ///扣一次
	starSpy.SetDiscoverRemainCount(self.StarSpyType, discoverRemainCountFinal) ///更新
	client.GetSyncMgr().SyncObject("StarSpyDiscoverMsg", starSpy)

	//msg := NewActionAttribChangeMsg(systemTypeStarSpyDiscover, "", 0, 0) ///生成一个属性变更消息
	/////无cd时需要扣次数
	//starSpyLabel := starSpy.GetStarSpyLabel(self.StarSpyType)               ///得到球探类型标签
	//starSpyDiscoverRemainCountLable := starSpyLabel + "DiscoverRemainCount" ///组合球探属性剩余次数标签
	//discoverRemainCount := starSpy.GetDiscoverRemainCount(self.StarSpyType)
	//discoverRemainCountFinal := discoverRemainCount - 1                               ///扣一次
	//starSpy.SetDiscoverRemainCount(self.StarSpyType, discoverRemainCountFinal)        ///更新
	//msg.AddAttribChangeInt(starSpyDiscoverRemainCountLable, discoverRemainCountFinal) ///通知剩余次数变更
	/////设置cd
	//starSpyCDLable := starSpyLabel + "CD"                                                          ///组合球探属性cd标签
	//configDiscoverCD := self.getConfigStarSpyDiscoverParam(self.StarSpyType, configItemDiscoverCD) ///得到配置中的球探发掘cd时间
	//discoverCDFinal := int(time.Now().Unix()) + configDiscoverCD                                   ///将分钟数转成秒
	//starSpy.SetDiscoverCD(self.StarSpyType, discoverCDFinal)                                       ///更新最新的球探CD时间
	//msg.AddAttribChangeInt(starSpyCDLable, discoverCDFinal)                                        ///通知CD变更
	//client.SendMsg(msg)                                                                            ///有存在的属性更新项才发送
	return true
}

///发送发掘球员结果消息return true
func (self *StarSpyDiscoverMsg) sendStarSpyDiscoverResultMsg(client IClient, result string, DiscoverResultStarType int, DiscoverStarList []int, DiscoverResultStarEvolveCount int) {
	msg := newStarSpyDiscoverResultMsg(result)
	msg.DiscoverResultStarType = DiscoverResultStarType
	msg.DiscoverStarList = DiscoverStarList
	msg.DiscoverResultStarEvolveCount = DiscoverResultStarEvolveCount
	client.SendMsg(msg)
}

///当幸运值到达100%时根据大列表生成稀有列表
func (self *StarSpyDiscoverMsg) discoverDrawRare(drawGroupList []int) []int {
	staticDataMgr := GetServer().GetStaticDataMgr()
	drawGrade := starGradePurple ///默认抽取品质为紫色
	if self.StarSpyType > primerStarSpy {
		drawGrade = starGradeOrange ///中级和高级球探取品质为橙色
	}
	drawList := []int{} ///候选组
	for i := range drawGroupList {
		///生成指定颜色的列表
		drawGroupIndex := drawGroupList[i]
		drawGroupStaticData := staticDataMgr.GetDrawGroupStaticData(drawGroupIndex)
		if nil == drawGroupStaticData {
			GetServer().GetLoger().Warn("StarSpyDiscoverMsg discoverDrawRare fail! drawGroupIndex %d is Invalid!", drawGroupIndex)
			break
		}
		element := staticDataMgr.GetStaticData(tableStarType, drawGroupStaticData.AwardType)
		if nil == element {
			continue ///查配置表无法找到指定记录时忽略
		}
		starTypeStaticData := element.(*StarTypeStaticData)
		if starTypeStaticData.Grade == drawGrade {
			drawList = append(drawList, drawGroupIndex) ///生成新的稀有表
		}
	}
	return drawList
}

func (self *StarSpyDiscoverMsg) GetDrawGroupList(client IClient) IntList { ///球探发掘球员逻辑
	staticDataMgr := GetServer().GetStaticDataMgr()
	starSpy := client.GetTeam().GetStarSpy() ///取得球探接口
	isFullDiscoverLuck := starSpy.IsFullDiscoverLuck(self.StarSpyType)
	groupNumber := drawGroupTypeStarSpy1
	switch self.StarSpyType {
	case primerStarSpy:
		if false == isFullDiscoverLuck {
			groupNumber = drawGroupTypeStarSpy1
		} else {
			groupNumber = drawGroupTypeStarSpyFull1
		}
		break
	case middleStarSpy:
		if false == isFullDiscoverLuck {
			groupNumber = drawGroupTypeStarSpy2
		} else {
			groupNumber = drawGroupTypeStarSpyFull2
		}
		break
	case expertStarSpy:
		if false == isFullDiscoverLuck {
			groupNumber = drawGroupTypeStarSpy3
		} else {
			groupNumber = drawGroupTypeStarSpyFull3
		}
		break
	}
	drawGroupList := staticDataMgr.GetDrawGroupIndexList(groupNumber) ///取得抽卡索引列表
	return drawGroupList
}

func (self *StarSpyDiscoverMsg) doAction(client IClient) bool { ///球探发掘球员逻辑
	///判断此类球探是否已经满幸运了
	const maxDrawCount = 5 ///最多抽5张卡
	msgResult := discoverResultOK
	starSpy := client.GetTeam().GetStarSpy() ///取得球探接口
	starCenter := client.GetTeam().GetStarCenter()
	team := client.GetTeam()
	msgStarSpyLuckChange := NewActionAttribChangeMsg(systemTypeStarSpyDiscover, "", 0, 0) ///生成一个属性变更消息
	starSpyLabel := starSpy.GetStarSpyLabel(self.StarSpyType)                             ///得到球探类型标签
	starSpyLuckLable := starSpyLabel + "Luck"                                             ///组合球探属性幸运标签
	isFullDiscoverLuck := starSpy.IsFullDiscoverLuck(self.StarSpyType)
	currentDiscoverLuck := starSpy.GetDiscoverLuck(self.StarSpyType)
	staticDataMgr := GetServer().GetStaticDataMgr()
	//drawGroupList := staticDataMgr.GetDrawGroupIndexList(drawGroupTypeStarSpy1) ///取得抽卡索引列表
	drawGroupList := self.GetDrawGroupList(client)
	///过滤掉已在球员中心中的球员列表
	drawGroupList = starCenter.GetDrawGroupFilterIndexList(starCenterTypeDiscover, drawGroupList)
	///过滤掉已在球队中的球员列表(修改: 不过滤)
	drawGroupList = team.GetDrawGroupFullIndexList(drawGroupList)
	totalTakeWeight, totalShowWeight := discoverGetDrawWeightTotal(drawGroupList) ///得到权重总和
	drawResultList := []int{}                                                     ///结果列表
	drawTakeOne, drawTakeRare, drawShowOne := 0, 0, 0                             ///选中id 展示id
	if true == isFullDiscoverLuck {
		starSpy.SetDiscoverLuck(self.StarSpyType, 0)                                         ///满幸运后需要重置
		msgStarSpyLuckChange.AddAttribChangeInt(starSpyLuckLable, 0)                         ///通知幸运变更
		drawRareList := IntList{}                                                            ///稀有列表
		drawRareList = self.discoverDrawRare(drawGroupList)                                  ///满幸运时生成稀有列表
		totalTakeWeightRare, totalShowWeightRare := discoverGetDrawWeightTotal(drawRareList) ///得到权重总和
		if totalTakeWeightRare > 0 {                                                         ///有稀有卡可抽时才去抽
			drawTakeRare, _ = discoverDrawOne(&drawRareList, &totalTakeWeightRare, &totalShowWeightRare, 0, true) ///抽取稀有卡
			msgResult = discoverResultLuck
		}
	}
	///不管是否幸运均需抽取一张普通的,避免因为权重值造成抽不到球员
	drawTakeOne, _ = discoverDrawOne(&drawGroupList, &totalTakeWeight, &totalShowWeight, drawTakeRare, true) ///第一张抽卡,后面四张为展示
	if drawTakeRare > 0 {
		drawTakeOne = drawTakeRare ///更新成抽取稀有球员type
	}

	starEvolveCount := 1 //starCenter.GetProbabilityEvolveCount(self.StarSpyType)

	if true == isFullDiscoverLuck {
		switch self.StarSpyType {
		case primerStarSpy:
			starEvolveCount = 1
		case middleStarSpy:
			starEvolveCount = 1
		case expertStarSpy:
			starEvolveCount = 1
		}
	}

	//初级球探寻找次数

	switch self.StarSpyType {
	case primerStarSpy:
		if starSpy.GetPrimerDiscoverCount(primerStarSpy) == false {
			drawTakeOne = staticDataMgr.GetConfigStaticDataInt(configStarSpy, configStarSpyDiscover, 1)
			starSpy.SetPrimerDiscoverCount(primerStarSpy)
		}
	case middleStarSpy:
		if starSpy.GetPrimerDiscoverCount(middleStarSpy) == false {
			drawTakeOne = staticDataMgr.GetConfigStaticDataInt(configStarSpy, configStarSpyDiscover, 2)
			starSpy.SetPrimerDiscoverCount(middleStarSpy)
		}
	case expertStarSpy:
		if starSpy.GetPrimerDiscoverCount(expertStarSpy) == false {
				drawTakeOne = staticDataMgr.GetConfigStaticDataInt(configStarSpy, configStarSpyDiscover, 3)
			//drawTakeOne = 4201
			starSpy.SetPrimerDiscoverCount(expertStarSpy)
		}
	}

	// discoverCount := starSpy.GetPrimerDiscoverCount(primerStarSpy)
	// if discoverCount < 2 {
	// 	// 初级寻找球员前3次 必定抽到2星球员  且球员固定 ID是4344，4297,4289

	// 	switch discoverCount {
	// 	case 0:
	// 		drawTakeOne = staticDataMgr.GetConfigStaticDataInt(configStarSpy, configStarSpyDiscover, 1)
	// 	case 1:
	// 		drawTakeOne = staticDataMgr.GetConfigStaticDataInt(configStarSpy, configStarSpyDiscover, 2)
	// 		//case 2:
	// 		//	drawTakeOne = staticDataMgr.GetConfigStaticDataInt(configStarSpy, configStarSpyDiscover, 3)
	// 	}
	// 	starEvolveCount = 1
	// 	starSpy.SetPrimerDiscoverCount()
	// }

	///存储奖励球员类型
	starSpy.SetAwardStarType(drawTakeOne, starEvolveCount)

	drawResultList = append(drawResultList, drawTakeOne) ///放入抽取普通球员type
	for i := 2; i <= maxDrawCount; i++ {                 ///一共抽4张展示卡
		drawShowOne, _ = discoverDrawOne(&drawGroupList, &totalTakeWeight, &totalShowWeight, 0, false) ///后面四张为展示
		if 0 == drawShowOne {
			GetServer().GetLoger().Warn("StarSpyDiscoverMsg doAction discoverDrawOne show get a zero draw!")
			return false
		}
		drawResultList = append(drawResultList, drawShowOne) ///取得结果
	}
	///发结果消息
	self.sendStarSpyDiscoverResultMsg(client, msgResult, drawTakeOne, drawResultList, starEvolveCount)
	///处理完成一次发掘给球探加经验逻辑
	if false == isFullDiscoverLuck {
		configDiscoverLuck := self.getConfigStarSpyDiscoverLuckAward()
		currentDiscoverLuck += configDiscoverLuck                                      ///累加
		starSpy.SetDiscoverLuck(self.StarSpyType, currentDiscoverLuck)                 ///满幸运后需要重置
		msgStarSpyLuckChange.AddAttribChangeInt(starSpyLuckLable, currentDiscoverLuck) ///通知幸运变更
	}
	client.SendMsg(msgStarSpyLuckChange)

	///更新球员训练中的日常任务
	dayTaskFunctionStarSpyType := dayTaskFunctionStarSpyLow + self.StarSpyType - 1
	team.GetTaskMgr().UpdateDayTaskFuntion(client.GetElement(), dayTaskFunctionStarSpyType)
	return true
}

func (self *StarSpyDiscoverMsg) processAction(client IClient) (result bool) {
	defer func() {
		if false == result {
			self.sendStarSpyDiscoverResultMsg(client, discoverResultFail, 0, nil, 0) ///发失败结果消息
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

type StarSpyOperateMsg struct { ///请求处理发掘球员消息
	MsgHead     `json:"head"` ///"starspy" "starspyoperatemsg"
	OperateType int           `json:"operatetype"` ///发掘后操作类型: 1.招至转会中心  2.转化经验  3.提升本队相同球员星级
}

func (self *StarSpyOperateMsg) GetTypeAndAction() (string, string) {
	return "starspy", "starspyoperatemsg"
}

func (self *StarSpyOperateMsg) checkStarCenter(client IClient) bool {
	loger := GetServer().GetLoger()
	payConfig := GetServer().GetStaticDataMgr().GetConfigStaticData(configStarSpy, configItemDiscoverConfig)
	maxMemberCount, _ := strconv.Atoi(payConfig.Param1) ///得到发掘球员中心的会员数上限数

	///得到当前发掘球员中心的会员数
	currentMemberCount := client.GetTeam().GetStarCenter().GetStarCenterMemberCount(starCenterTypeDiscover)
	if loger.CheckFail("currentMemberCount < maxMemberCount",
		currentMemberCount < maxMemberCount, currentMemberCount, maxMemberCount) {
		return false
	}

	return true
}

func (self *StarSpyOperateMsg) checkAction(client IClient) bool {
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	teamInfo := team.GetInfo()
	starSpy := team.GetStarSpy()
	isOperateLegal := (self.OperateType > OperateBegin) && (self.OperateType < OperateEnd)
	if loger.CheckFail("isOperateLegal == true", isOperateLegal == true, isOperateLegal, true) {
		return false ///检测操作类型合法性
	}

	awardStarType, awardStarEvolveCount := starSpy.GetAwardStarType()
	if loger.CheckFail("awardStarType != 0", awardStarType != 0, awardStarType, 0) {
		return false ///检测奖励球员类型是否存在
	}

	starType := GetServer().GetStaticDataMgr().Unsafe().GetStarType(awardStarType)
	if loger.CheckFail("starType != nil", starType != nil, starType, nil) {
		return false ///检测奖励球员类型是否存在
	}

	switch self.OperateType {
	case OperateCanvass:
		///检测球员中心人员上限
		// if self.checkStarCenter(client) == false {
		// 	return false
		// }

		///检测队伍是否存在球星
		isTeamHas := team.HasStar(awardStarType)
		if loger.CheckFail("isTeamHas == false", isTeamHas == false, isTeamHas, false) {
			return false
		}

	case OperateAddExp:
		///经验池检查
		if loger.CheckFail("StarExpPool <= ExpPoolLimit", teamInfo.StarExpPool <= ExpPoolLimit,
			teamInfo.StarExpPool, ExpPoolLimit) {
			return false
		}

	case OperateLevelup:
		///检测队伍是否存在奖励球星
		isTeamHas := team.HasStar(awardStarType)
		if loger.CheckFail("isTeamHas == true", isTeamHas == true, isTeamHas, true) {
			return false
		}

		///检测队伍中目标球星星级是否低于奖励类型
		teamStar := team.GetStarFromType(awardStarType)
		teamStarInfo := teamStar.GetInfo()
		if loger.CheckFail("awardStarEvolveCount > teamStarInfo.EvolveCount",
			awardStarEvolveCount > teamStarInfo.EvolveCount, awardStarEvolveCount, teamStarInfo.EvolveCount) {
			return false
		}
	}

	return true
}

func (self *StarSpyOperateMsg) GetFirstAndGrowValue(starTypeInfo *StarTypeStaticData) (int, int) { //得到初始值与成长值

	firstValue := starTypeInfo.Pass + starTypeInfo.Steals +
		starTypeInfo.Dribbling + starTypeInfo.Sliding +
		starTypeInfo.Shooting + starTypeInfo.GoalKeeping +
		starTypeInfo.Body + starTypeInfo.Speed
	growValue := starTypeInfo.PassGrow + starTypeInfo.StealsGrow +
		starTypeInfo.DribblingGrow + starTypeInfo.SlidingGrow +
		starTypeInfo.ShootingGrow + starTypeInfo.GoalKeepingGrow

	return firstValue, growValue
}

func (self *StarSpyOperateMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	starSpy := client.GetTeam().GetStarSpy()
	//	starCenter := client.GetTeam().GetStarCenter()
	sync := client.GetSyncMgr()
	awardStarType, awardStarEvolveCount := starSpy.GetAwardStarType()
	switch self.OperateType {
	case OperateCanvass:
		// awardMemberID := starCenter.AwardSpecialMember(starCenterTypeDiscover, awardStarType, awardStarEvolveCount, 0) ///发奖

		// ///同步球员中心新加成员信息给客户端
		// starCenterAddMemberMsg := NewStarCenterAddMemberMsg()
		// starCenterMember := starCenter.GetStarCenterMember(starCenterTypeDiscover, awardMemberID)
		// starCenterAddMemberMsg.AddMember(starCenterMember)
		// client.SendMsg(starCenterAddMemberMsg)
		// starSpy.SetAwardStarType(0, 0) ///处理完毕,清空类型

		///奖励球员,并赋予球员初始星级
		starID := team.AwardStar(awardStarType)
		star := team.GetStar(starID)
		star.SetStarCount(awardStarEvolveCount)
		//star.GetInfo().EvolveCount = awardStarEvolveCount
		sync.syncAddStar(IntList{starID})

	case OperateAddExp:
		///（初始值之和+3*成长值之和）*星级^2
		starType := GetServer().GetStaticDataMgr().Unsafe().GetStarType(awardStarType)
		//		firstValue, growValue := self.GetFirstAndGrowValue(starType)
		//		teamInfo := team.GetInfo()
		//		exp := (firstValue + growValue*3) * (awardStarEvolveCount * awardStarEvolveCount)

		//		if teamInfo.StarExpPool+exp > ExpPoolLimit {
		//			teamInfo.StarExpPool = ExpPoolLimit
		//		} else {
		//			teamInfo.StarExpPool += exp
		//		}
		//		fmt.Printf("FirstValue = %d  \r\n   GrowValue = %d  \r\n  StarLevel = %d \r\n  EXP = %d \r\n",
		//			firstValue, growValue, awardStarEvolveCount, exp)
		///同步到客户端
		//		sync.SyncObject("StarSpyOperateMsg", team)
		starSpy.SetAwardStarType(0, 0) ///处理完毕,清空类型
		///给星卡
		starCardCount := team.GetStarCardCount(starType, awardStarEvolveCount)
		team.AwardObject(ItemStarCard, starCardCount, 0, 0) ///赠星卡
	case OperateLevelup:
		teamStar := team.GetStarFromType(awardStarType)
		starType := teamStar.GetTypeInfo()
		///给星卡
		starCardCount := team.GetStarCardCount(starType, teamStar.EvolveCount)
		team.AwardObject(ItemStarCard, starCardCount, 0, 0) ///赠星卡

		teamStar.SetStarCount(awardStarEvolveCount)
		//teamStarInfo := teamStar.GetInfo()
		//teamStarInfo.EvolveCount = awardStarEvolveCount
		/////根据突破次数决定品质
		//if awardStarEvolveCount <= 2 {
		//	teamStarInfo.Grade = starGradeGreen
		//} else if awardStarEvolveCount > 2 && awardStarEvolveCount <= 4 {
		//	teamStarInfo.Grade = starGradeBlue
		//} else if awardStarEvolveCount > 4 && awardStarEvolveCount <= 6 {
		//	teamStarInfo.Grade = starGradePurple
		//} else {
		//	teamStarInfo.Grade = starGradeOrange
		//}

		sync.syncStarCalcInfo(teamStar)
		sync.SyncObject("StarSpyOperateMsg", teamStar)
		starSpy.SetAwardStarType(0, 0) ///处理完毕,清空类型
	}
	return true
}

func (self *StarSpyOperateMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false {
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
