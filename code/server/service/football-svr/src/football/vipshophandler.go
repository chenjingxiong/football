package football

//商城上架/下架商品消息
// type VipShopCommodityChangeMsg struct {
// 	MsgHead         `json:"head"` // "vipshop", "vipshopcommoditychange"
// 	CommodityIDList IntList       `json:"commdityidlist"` //商品列表
// 	ChangeType      bool          `json:"changetype"`     //true 为上架   false 为下架
// }

// func (self *VipShopCommodityChangeMsg) GetTypeAndAction() (string, string) {
// 	return "vipshop", "vipshopcommoditychange"
// }

//查询商城商品信息
type VipShopCommodityQueryMsg struct {
	MsgHead  `json:"head"` // "vipshop", "vipshopcommodityquery"
	SaleType int           `json:"saletype"` // 1为促销  2为日常
}

func (self *VipShopCommodityQueryMsg) GetTypeAndAction() (string, string) {
	return "vipshop", "vipshopcommodityquery"
}

// func (self *VipShopCommodityQueryMsg) checkVersion(client IClient) (bool, int) {
// 	//比对版本号
// 	vipShopMgr := client.GetTeam().GetVipShopMgr()
// 	currentVer := vipShopMgr.GetVersion()
// 	if currentVer == self.Version {
// 		return true, currentVer //版本号一致
// 	}

// 	return false, currentVer
// }

func GetCommodityAllMsg(client IClient) ([]*VipShopStaticData, []*VipShopStaticData) {
	staticDataMgr := GetServer().GetStaticDataMgr()
	vipShopMgr := client.GetTeam().GetVipShopMgr()
	//	currentVer := vipShopMgr.GetVersion()

	//发送最新版本的商品信息
	commodityList := staticDataMgr.GetVipShopStaticDataMap()
	dailyList := []*VipShopStaticData{}      ///日常商品
	promotionsList := []*VipShopStaticData{} ///促销商品
	for _, v := range commodityList {
		commodityType := v

		//判断商品类型
		if commodityType.Type2 == TimeInfinite {
			//无时限的日常商品直接发送信息
			node := new(VipShopStaticData)
			node.ID = commodityType.ID
			node.ItemType = commodityType.ItemType
			node.Num = commodityType.Num
			node.Type2 = commodityType.Type2
			node.Class = commodityType.Class
			node.MoneyID = commodityType.MoneyID
			node.MoneyNum = commodityType.MoneyNum
			node.StartTime = commodityType.StartTime
			node.Duration = commodityType.Duration
			node.Limit = 0
			node.Class = commodityType.Class
			node.Discount = commodityType.Discount
			node.CD = commodityType.CD

			dailyList = append(dailyList, node)

		} else if commodityType.Type2 == TimeLimit {
			//有时限的促销商品需判断

			isPut, times := vipShopMgr.CheckShoppingTime(v.ID)
			commodityInfo := vipShopMgr.GetCommodityInfo(v.ID)
			if isPut == false {
				continue
			}

			nowTime := Now()                                                             //当前时间
			startTime := int(vipShopMgr.GetStampFromStaticData(commodityType.StartTime)) //开始时间
			cycle := commodityType.CD + commodityType.Duration
			if 0 == cycle {
				GetServer().GetLoger().Warn("vipshop cycle == 0  commodityID:%d", commodityType.ID)
				continue
			}
			curTime := (nowTime - startTime) % (cycle) //当前轮回中时间
			startTime = nowTime - curTime

			if times >= 1 && commodityInfo.Cycle != times {
				// 轮回次数超过一次,重新给予限购次数
				commodityInfo.Limittimes = commodityType.Limit
				commodityInfo.Cycle = times
			}

			//推送促销商品信息
			node := new(VipShopStaticData)
			node.ID = commodityType.ID
			node.ItemType = commodityType.ItemType
			node.Num = commodityType.Num
			node.Type2 = commodityType.Type2
			node.Class = commodityType.Class
			node.MoneyID = commodityType.MoneyID
			node.MoneyNum = commodityType.MoneyNum
			node.StartTime = startTime
			node.Duration = commodityType.Duration
			node.Limit = commodityInfo.Limittimes
			node.Class = commodityType.Class
			node.Discount = commodityType.Discount
			node.CD = commodityType.CD

			promotionsList = append(promotionsList, node)
		}
	}
	return dailyList, promotionsList
}

func (self *VipShopCommodityQueryMsg) processAction(client IClient) bool {

	//	isVerSame, _ := self.checkVersion(client)

	//得到所有商品信息
	dailList, promotionsList := GetCommodityAllMsg(client)

	if self.SaleType == TimeLimit {
		SendVipShopCommodityQueryResultMsg(client, promotionsList) //促销
	} else if self.SaleType == TimeInfinite {
		SendVipShopCommodityQueryResultMsg(client, dailList) //日常
	}

	return true
}

//返回查询结果信息
// type VipShopCommodityQueryResultMsg struct {
// 	MsgHead    `json:"head"` // "vipshop", "vipshopcommodityqueryresult"
// 	ItemID     int           `json:"itemid"`     //记录ID
// 	ItemType   int           `json:"itemtype"`   //物品类型
// 	ItemNum    int           `json:"itemnum"`    //物品数量
// 	Version    int           `json:"version"`    //版本号
// 	SaleType   int           `json:"saletype"`   //1.限时促销 2.日常商品
// 	MoneyID    int           `json:"moneyid"`    //使用货币id
// 	MoneyNum   int           `json:"moneynum"`   //货币支付数量
// 	PastTime   int           `json:"pastTime"`   //倒计时间
// 	Duration   int           `json:"duration"`   //持续时间
// 	LimitTimes int           `json:"limittimes"` //可购买次数
// 	Class      int           `json:"class"`      //类别
// 	Discount   int           `json:"discount"`   //折扣
// }

type VipShopCommodityQueryResultMsg struct {
	MsgHead       `json:"head"`        // "vipshop", "vipshopcommodityqueryresult"
	GoodsTypeList []*VipShopStaticData `json:"goodstype"`
}

func SendVipShopCommodityQueryResultMsg(client IClient, GoodsTypeList []*VipShopStaticData) {
	msg := new(VipShopCommodityQueryResultMsg)
	msg.GoodsTypeList = GoodsTypeList
	client.SendMsg(msg)
}

// func NewVipShopCommodityQueryResultMsg(id int, itemType int, itemNum int, version int, saleType int,
// 	moneyID int, moneyNum int, pastTime int, limitTimes int, duration int, class int, discount int) *VipShopCommodityQueryResultMsg {
// 	msg := new(VipShopCommodityQueryResultMsg)
// 	msg.ItemID = id
// 	msg.ItemType = itemType
// 	msg.ItemNum = itemNum
// 	msg.Version = version
// 	msg.SaleType = saleType
// 	msg.MoneyID = moneyID
// 	msg.MoneyNum = moneyNum
// 	msg.PastTime = pastTime
// 	msg.LimitTimes = limitTimes
// 	msg.Duration = duration
// 	msg.Class = class
// 	msg.Discount = discount
// 	return msg
// }

func (self *VipShopCommodityQueryResultMsg) GetTypeAndAction() (string, string) {
	return "vipshop", "vipshopcommodityqueryresult"
}

//购买物品消息
type VipShopCommodityBuyMsg struct {
	MsgHead `json:"head"` //"vipshop", "vipshopcommoditybuy"
	ItemID  int           `json:"itemid"`  //记录ID
	ItemNum int           `json:"itemnum"` //数量
	Version int           `json:"version"` //版本号
}

func (self *VipShopCommodityBuyMsg) GetTypeAndAction() (string, string) {
	return "vipshop", "vipshopcommoditybuy"
}

func (self *VipShopCommodityBuyMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	vipShopMgr := team.GetVipShopMgr()
	staticDataMgr := GetServer().GetStaticDataMgr()
	loger := GetServer().GetLoger()

	//取得购买商品的静态信息
	commodityType := staticDataMgr.Unsafe().GetVipShopItemInfo(self.ItemID)
	if loger.CheckFail("commodityType != nil", commodityType != nil, commodityType, nil) {
		return false //商品不存在
	}

	// serverVer := vipShopMgr.GetVersion()
	// if loger.CheckFail("serverVer == self.Version", serverVer == self.Version, serverVer, self.Version) {
	// 	return false //版本号错误
	// }

	//检测玩家货币是否足够
	isMoneyEnough := false
	singleNeedMoney := float32(commodityType.MoneyNum) * float32(commodityType.Discount) / 10.0 //得到折扣
	needMoney := int(singleNeedMoney) * self.ItemNum
	switch commodityType.MoneyID {
	case awardTypeCoin: //球币
		isMoneyEnough = needMoney < team.GetCoin()
	case awardTypeTicket: //球票
		isMoneyEnough = needMoney < team.GetTicket()
	}

	if loger.CheckFail("isMoneyEnough == true", isMoneyEnough == true, isMoneyEnough, true) {
		return false //货币不足
	}

	//检测玩家背包是否已满
	isStoreFull := team.IsStoreFull(commodityType.ItemType, self.ItemNum)
	if loger.CheckFail("isStoreFull == false", isStoreFull == false, isStoreFull, false) {
		return false
	}

	//判断商品的销售类型
	if commodityType.Type2 == TimeLimit {
		isShoppingTime, _ := vipShopMgr.CheckShoppingTime(self.ItemID)
		if loger.CheckFail("isShoppingTime == true", isShoppingTime == true, isShoppingTime, true) {
			//得到所有商品信息
			dailList, promotionsList := GetCommodityAllMsg(client)
			SendVipShopCommodityQueryResultMsg(client, dailList)       //日常
			SendVipShopCommodityQueryResultMsg(client, promotionsList) //促销
			return false                                               //非抢购时间
		}

		commodityInfo := vipShopMgr.GetCommodityInfo(self.ItemID)
		if loger.CheckFail("commodityInfo != nil", commodityInfo != nil, commodityInfo, nil) {
			return false //无该信息
		}

		if loger.CheckFail("LimitTimes - self.ItemNum >= 0", commodityInfo.Limittimes-self.ItemNum >= 0, commodityInfo.Limittimes, self.ItemNum) {
			return false //限购次数已尽
		}
	}
	//若为限购物品,则扣除次数
	if commodityType.Type1 == DayTaskBuyGoods {
		if loger.CheckFail("self.ItemNum==1", self.ItemNum == 1, self.ItemNum, 1) {
			return false //日常任务购买商品的指定数量只能是1
		}
	}
	return true
}

func (self *VipShopCommodityBuyMsg) payAction(client IClient) bool {
	syncMgr := client.GetSyncMgr()
	staticDataMgr := GetServer().GetStaticDataMgr()
	commodityType := staticDataMgr.Unsafe().GetVipShopItemInfo(self.ItemID)
	team := client.GetTeam()
	vipShopMgr := team.GetVipShopMgr()

	singleNeedMoney := float32(commodityType.MoneyNum) //* float32(commodityType.Discount) / 10.0得到折扣
	needMoney := int(singleNeedMoney) * self.ItemNum
	switch commodityType.MoneyID {
	case awardTypeCoin: //扣除球币
		team.PayCoin(needMoney)
	case awardTypeTicket: //扣除球票
		team.PayTicket(needMoney)
		client.SetMoneyRecord(PlayerCostMoney, Pay_VipShop, needMoney, team.GetTicket())
	}

	syncMgr.SyncObject("VipShopCommodityBuyMsg", team)

	//若为限购物品,则扣除次数
	if commodityType.Type2 == TimeLimit {
		commdityInfo := vipShopMgr.GetCommodityInfo(self.ItemID)
		commdityInfo.Limittimes -= self.ItemNum
	}

	return true
}

func (self *VipShopCommodityBuyMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	staticDataMgr := GetServer().GetStaticDataMgr()
	commodityType := staticDataMgr.Unsafe().GetVipShopItemInfo(self.ItemID)
	buyCount := commodityType.Num * self.ItemNum
	team.AwardObject(commodityType.ItemType, buyCount, 0, 0)
	// itemID := itemMgr.AwardItem(commodityType.ItemType, commodityType.Num*self.ItemNum)
	// item := itemMgr.GetItem(itemID)

	// syncMgr := client.GetSyncMgr()
	// syncMgr.SyncObject("VipShopCommodityBuyMsg", item)

	///更新日常任务中的商城购买
	team.GetTaskMgr().UpdateDayTaskVipBuy(client.GetElement(), self.ItemID)
	return true
}

func (self *VipShopCommodityBuyMsg) processAction(client IClient) (result bool) {
	defer func() {
		if false == result {
			self.sendVipShopResultMsg(client, msgResultFail) ///发失败结果消息
		} else {
			self.sendVipShopResultMsg(client, msgResultOK) ///发失败结果消息
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

type VipShopCommodityBuyResultMsg struct {
	MsgHead `json:"head"` //"vipshop", "vipshopcommoditybuyresult"
	Result  string        `json:"result"` ///购买结果 ok or fail
}

func (self *VipShopCommodityBuyResultMsg) GetTypeAndAction() (string, string) {
	return "vipshop", "vipshopcommoditybuyresult"
}

func NewVipShopCommodityBuyResultMsg(client IClient, result string) *VipShopCommodityBuyResultMsg {
	msg := new(VipShopCommodityBuyResultMsg)
	msg.Result = result
	return msg
}

func (self *VipShopCommodityBuyMsg) sendVipShopResultMsg(client IClient, result string) {
	msg := NewVipShopCommodityBuyResultMsg(client, result)
	client.SendMsg(msg)
}

type VipShopQueryBuyDiamondCountMsg struct { ///客户端请求查询购钻石套餐次数
	MsgHead `json:"head"` // "vipshop", "querybuydiamondcount"
}

func (self *VipShopQueryBuyDiamondCountMsg) GetTypeAndAction() (string, string) {
	return "vipshop", "querybuydiamondcount"
}

type VipShopQueryBuyDiamondCountResultMsg struct { ///客户端请求查询购钻石套餐次数
	MsgHead `json:"head"` // "vipshop", "querybuydiamondcountresult"
	//MoneyID           int           `json:"moneyid"`           ///完成购买套餐id,0表示刷新全部,非0表示刷新指定套餐
	MoneyBuyCountList IntList `json:"moneybuycountlist"` ///已购买次数列表
}

func (self *VipShopQueryBuyDiamondCountResultMsg) GetTypeAndAction() (string, string) {
	return "vipshop", "querybuydiamondcountresult"
}

func SendQueryBuyDiamondCountResultMsg(client IClient) {
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeMoneyBuyCount)

	queryBuyDiamondCountResultMsg := new(VipShopQueryBuyDiamondCountResultMsg)
	//queryBuyDiamondCountResultMsg.MoneyID = moneyID
	queryBuyDiamondCountResultMsg.MoneyBuyCountList = IntList{resetAttrib.Value1,
		resetAttrib.Value2, resetAttrib.Value3, resetAttrib.Value4, resetAttrib.Value5,
		resetAttrib.Value6, resetAttrib.Value7, resetAttrib.Value8}
	client.SendMsg(queryBuyDiamondCountResultMsg)
}

func (self *VipShopQueryBuyDiamondCountMsg) processAction(client IClient) (result bool) {
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeMoneyBuyCount)
	if nil == resetAttrib {
		resetAttrib = resetAttribMgr.AddResetAttrib(ResetAttribTypeMoneyBuyCount, 0, IntList{0})
		//resetAttrib = resetAttribMgr.GetResetAttrib(resetAttribID)
	}
	SendQueryBuyDiamondCountResultMsg(client)
	//	vipShopQueryBuyDiamondCountResultMsg := new(VipShopQueryBuyDiamondCountResultMsg)
	//	vipShopQueryBuyDiamondCountResultMsg.MoneyBuyCountList = IntList{resetAttrib.Value1,
	//		resetAttrib.Value2, resetAttrib.Value3, resetAttrib.Value4, resetAttrib.Value5,
	//		resetAttrib.Value6, resetAttrib.Value7, resetAttrib.Value8}
	//	client.SendMsg(vipShopQueryBuyDiamondCountResultMsg)
	return true
}

//type VipShopCreateOrderMsg struct { ///客户端请求创建订单
//	MsgHead `json:"head"` // "vipshop", "createorder"
//	MoneyID int           `json:"moneyid"` ///请求购买套餐id
//}

//func (self *VipShopCreateOrderMsg) GetTypeAndAction() (string, string) {
//	return "vipshop", "CreateOrder"
//}

type VipShopBuyDiamondMsg struct { ///客户端请求购钻石与月卡
	MsgHead    `json:"head"` // "vipshop", "buydiamond"
	MoneyID    int           `json:"moneyid"`    ///请求购买套餐id
	TeamID     int           `json:"teamid"`     ///球队id
	PayOrderID int           `json:"payorderid"` ///支付订单编号
	PayMoney   int           `json:"paymoney"`   ///支付多少(一级货币)人民币,单位分
}

func (self *VipShopBuyDiamondMsg) GetTypeAndAction() (string, string) {
	return "vipshop", "buydiamond"
}

func (self *VipShopBuyDiamondMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	//	vipShopMgr := team.GetVipShopMgr()
	loger := GetServer().GetLoger()
	if loger.CheckFail("self.MoneyID > 0", self.MoneyID > 0, self.MoneyID, 0) {
		return false //套餐编号不正确
	}
	moneyType := GetServer().GetStaticDataMgr().GetMoneyType(self.MoneyID)
	if loger.CheckFail("moneyType!=nil", moneyType != nil, moneyType, nil) {
		return false //套餐类型不存在
	}
	if team.IsGM() == false { ///如果是GM则不做严格验证,方便测试
		moneyPrice := moneyType.Money * 100 ///转换将人民币"元"价格转换成"分"价格
		if loger.CheckFail("self.PayMoney", self.PayMoney == moneyPrice, self.PayMoney, moneyPrice) {
			return false //商品价格与平台扣费不符,客户端被篡改?回调消息被破解?
		}
	}
	if moneyType.IsMonthCard() == false {
		resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeMoneyBuyCount)
		if nil == resetAttrib {
			resetAttrib = resetAttribMgr.AddResetAttrib(ResetAttribTypeMoneyBuyCount, 0, IntList{0})
		}
		if loger.CheckFail("resetAttrib!=nil", resetAttrib != nil, resetAttrib, nil) {
			return false //已购次信息不存在
		}
	}
	if loger.CheckFail("self.TeamID>0", self.TeamID > 0, self.TeamID, 0) {
		return false //球队编号必须大于0
	}
	if loger.CheckFail("self.PayOrderID>0", self.PayOrderID > 0, self.PayOrderID, 0) {
		return false //订单号不存在
	}
	return true
}

func (self *VipShopBuyDiamondMsg) payAction(client IClient) bool {
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeMoneyBuyCount)
	clientObj := client.GetElement()
	moneyType := GetServer().GetStaticDataMgr().GetMoneyType(self.MoneyID)
	MoneyBuyCountList := IntList{resetAttrib.Value1, resetAttrib.Value2,
		resetAttrib.Value3, resetAttrib.Value4, resetAttrib.Value5,
		resetAttrib.Value6, resetAttrib.Value7, resetAttrib.Value8}
	moneyIndex := self.MoneyID - 1
	moneyBuyCount := MoneyBuyCountList[moneyIndex]
	totalAwardNum := moneyType.Diamonds
	if moneyBuyCount >= moneyType.Num {
		totalAwardNum += moneyType.Give ///限购次数已用完
	} else {
		totalAwardNum += moneyType.NumGive
	}
	if moneyType.IsMonthCard() == true { ///买月卡特殊处理
		totalAwardNum = 0
	}
	balance := team.Ticket + totalAwardNum
	donePayOrder := GetServer().GetSDKMgr().DonePayOrder(self.PayOrderID, team.Name, clientObj.userID, balance)
	if loger.CheckFail("donePayOrder==true", donePayOrder == true, donePayOrder, true) {
		return false //订单号不存在
	}
	if moneyType.IsMonthCard() == true { ///买月卡特殊处理
		return true ///月卡无限购次数
	}
	//resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeMoneyBuyCount)
	MoneyBuyCountPtrList := IntPtrList{&resetAttrib.Value1, &resetAttrib.Value2,
		&resetAttrib.Value3, &resetAttrib.Value4, &resetAttrib.Value5,
		&resetAttrib.Value6, &resetAttrib.Value7, &resetAttrib.Value8}
	//moneyIndex := self.MoneyID - 1
	currentBuyCountPtr := MoneyBuyCountPtrList[moneyIndex]
	*currentBuyCountPtr = (*currentBuyCountPtr) + 1 ///增加一次购买次数
	return true
}

func (self *VipShopBuyDiamondMsg) doAction(client IClient) bool {
	//	loger := GetServer().GetLoger()
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeMoneyBuyCount)
	//	vipShopMgr := team.GetVipShopMgr()
	moneyType := GetServer().GetStaticDataMgr().GetMoneyType(self.MoneyID)
	if moneyType.IsMonthCard() == true { ///买月卡特殊处理
		vipBuyMonthCardMsg := new(VipBuyMonthCardMsg)
		vipBuyMonthCardMsg.processAction(client)
		return true
	}
	MoneyBuyCountList := IntList{resetAttrib.Value1, resetAttrib.Value2,
		resetAttrib.Value3, resetAttrib.Value4, resetAttrib.Value5,
		resetAttrib.Value6, resetAttrib.Value7, resetAttrib.Value8}
	moneyIndex := self.MoneyID - 1
	moneyBuyCount := MoneyBuyCountList[moneyIndex] - 1
	totalAwardNum := moneyType.Diamonds
	if moneyBuyCount >= moneyType.Num {
		totalAwardNum += moneyType.Give ///限购次数已用完
	} else {
		totalAwardNum += moneyType.NumGive
	}
	team.AwardObject(awardTypeTicket, totalAwardNum, 0, 0)      ///给钻石
	team.AwardObject(awardTypeVipExp, moneyType.Diamonds, 0, 0) ///给vip经验
	clientObj := client.GetElement()
	clientObj.RechargeRecord(Get_VIPShopMoney, totalAwardNum)
	SendQueryBuyDiamondCountResultMsg(client)
	resetAttrib.Save()
	team.CheckPayAward(self.MoneyID)

	//! 开服活动
	team.OSActivityValue.Refresh(team, 0)

	return true
}

func (self *VipShopBuyDiamondMsg) broacastMsg(userMgr *UserMgr) bool {
	client := userMgr.GetClientByTeamID(self.TeamID)
	if client != nil {
		self.processAction(client)
	}
	return true
}

func (self *VipShopBuyDiamondMsg) processAction(client IClient) (result bool) {
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

type VipQueryDayAwardResultMsg struct { ///客户端查询vip每日礼包状态
	MsgHead           `json:"head"` // "vipshop", "querydayawardresult"
	HasAccpetVipLevel int           `json:"hasaccpetviplevel"` ///已领奖的viplevel
	NextResetDelay    int           `json:"nextresetdelay"`    ///离下次重置的时间间隔,单位是秒
}

func (self *VipQueryDayAwardResultMsg) GetTypeAndAction() (string, string) {
	return "vipshop", "querydayawardresult"
}

type VipQueryDayAwardMsg struct { ///客户端查询vip每日礼包状态
	MsgHead `json:"head"` // "vipshop", "querydayaward"
}

func (self *VipQueryDayAwardMsg) GetTypeAndAction() (string, string) {
	return "vipshop", "querydayaward"
}

func (self *VipQueryDayAwardMsg) processAction(client IClient) (result bool) { ///客户端查询vip每日礼包状态
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeVipPrivilege)
	if nil == resetAttrib {
		resetAttribMgr.AddResetAttrib(ResetAttribTypeVipPrivilege, 0, nil)
		resetAttrib = resetAttribMgr.GetResetAttrib(ResetAttribTypeVipPrivilege)
	}
	resetAttrib.ResetVipPrivilege() ///尝试重置vip信息
	vipQueryDayAwardResultMsg := new(VipQueryDayAwardResultMsg)
	vipQueryDayAwardResultMsg.HasAccpetVipLevel = resetAttrib.Value1               ///已可领奖次数
	vipQueryDayAwardResultMsg.NextResetDelay = Max(0, resetAttrib.ResetTime-Now()) ///得到下次重置时间间隔,秒单位
	client.SendMsg(vipQueryDayAwardResultMsg)
	return true
}

type VipAccpetDayAwardMsg struct { ///客户端请求领取vip每日礼包
	MsgHead `json:"head"` // "vipshop", "accpetdayaward"
}

func (self *VipAccpetDayAwardMsg) GetTypeAndAction() (string, string) {
	return "vipshop", "accpetdayaward"
}

func (self *VipAccpetDayAwardMsg) checkAction(client IClient) (result bool) { ///验货
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	vipLevel := team.GetVipLevel()
	resetAttribMgr := team.GetResetAttribMgr()
	if loger.CheckFail("vipLevel > 0", vipLevel > 0, vipLevel, 0) {
		return false //有vip的人才能领取
	}
	vipInfo := GetServer().GetStaticDataMgr().GetVipInfo(vipLevel)
	if loger.CheckFail("vipInfo!=nil", vipInfo != nil, vipInfo, nil) {
		return false //找不到vip配置信息
	}
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeVipPrivilege)
	if loger.CheckFail("resetAttrib!=nil", resetAttrib != nil, resetAttrib, nil) {
		return false //找不到vip礼包领取信息
	}
	hasAccpetVipLevel := resetAttrib.Value1 ///已领取次数
	if loger.CheckFail("hasAccpetCount==0", hasAccpetVipLevel < vipLevel, hasAccpetVipLevel, vipLevel) {
		return false //此vip等级的vip礼包已领取了
	}
	return true
}

func (self *VipAccpetDayAwardMsg) payAction(client IClient) (result bool) { ///付款
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeVipPrivilege)
	resetAttrib.Value1 = team.VipLevel ///更新已领取vip每日礼包的vip等级
	return true
}

func (self *VipAccpetDayAwardMsg) doAction(client IClient) (result bool) { ///发货
	team := client.GetTeam()
	vipLevel := team.GetVipLevel()
	vipInfo := GetServer().GetStaticDataMgr().GetVipInfo(vipLevel)
	if vipInfo.Param4 > 0 && vipInfo.Param5 > 0 {
		team.AwardObject(vipInfo.Param4, vipInfo.Param5, 0, 0) ///赏道具1
	}
	if vipInfo.Param9 > 0 && vipInfo.Param10 > 0 {
		team.AwardObject(vipInfo.Param9, vipInfo.Param10, 0, 0) ///赏道具1
	}
	return true
}

func (self *VipAccpetDayAwardMsg) processAction(client IClient) (result bool) { ///领取vip每日礼包
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

type VipQueryMonthCardResultMsg struct { ///客户端查询月卡结果信息
	MsgHead        `json:"head"` // "vipshop", "querymonthcardresult"
	DeadLine       int           `json:"deadline"`       ///最后期限,0表示未买月卡
	HasAccpetCount int           `json:"hasaccpetcount"` ///已领取的礼包次数
}

func (self *VipQueryMonthCardResultMsg) GetTypeAndAction() (string, string) {
	return "vipshop", "querymonthcardresult"
}

type VipQueryMonthCardMsg struct { ///客户端请求查询月卡信息
	MsgHead `json:"head"` // "vipshop", "querymonthcard"
}

func (self *VipQueryMonthCardMsg) GetTypeAndAction() (string, string) {
	return "vipshop", "querymonthcard"
}

func (self *VipQueryMonthCardMsg) processAction(client IClient) (result bool) { ///客户端请求查询月卡信息
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeMonthCard)
	if nil == resetAttrib {
		resetAttribMgr.AddResetAttrib(ResetAttribTypeMonthCard, 0, nil)
		resetAttrib = resetAttribMgr.GetResetAttrib(ResetAttribTypeMonthCard)
	}
	resetAttrib.ResetVipMonthCard() ///尝试重置vip月卡信息
	vipQueryMonthCardResultMsg := new(VipQueryMonthCardResultMsg)
	vipQueryMonthCardResultMsg.DeadLine = resetAttrib.Value1
	vipQueryMonthCardResultMsg.HasAccpetCount = resetAttrib.Value2
	client.SendMsg(vipQueryMonthCardResultMsg)
	team.RedressPay() ///发可能产生的离线充值补偿
	return true
}

type VipBuyMonthCardMsg struct { ///客户端请求购买月卡
	MsgHead `json:"head"` // "vipshop", "buymonthcard"
}

func (self *VipBuyMonthCardMsg) GetTypeAndAction() (string, string) {
	return "vipshop", "buymonthcard"
}

func (self *VipBuyMonthCardMsg) processAction(client IClient) (result bool) { ///客户端请求购买月卡
	//	loger := GetServer().GetLoger()
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeMonthCard)
	if nil == resetAttrib {
		resetAttribMgr.AddResetAttrib(ResetAttribTypeMonthCard, 0, nil)
		resetAttrib = resetAttribMgr.GetResetAttrib(ResetAttribTypeMonthCard)
	}
	//	vipShopMgr := team.GetVipShopMgr()
	now := Now()
	//	if resetAttrib.Value1 <= now {
	resetAttrib.Value1 = Max(resetAttrib.Value1, now) ///如果以前没买过月卡需要将日期设置为当前时间
	resetAttrib.Value1 += 30 * 24 * 60 * 60           ///向后推30天
	//	}
	moneyType := GetServer().GetStaticDataMgr().GetMoneyType(MonthCardMoneyID)
	team.AwardObject(awardTypeVipExp, moneyType.Money*10, 0, 0)
	vipQueryMonthCardResultMsg := new(VipQueryMonthCardResultMsg)
	vipQueryMonthCardResultMsg.DeadLine = resetAttrib.Value1
	vipQueryMonthCardResultMsg.HasAccpetCount = resetAttrib.Value2
	client.SendMsg(vipQueryMonthCardResultMsg)
	//clientObj := client.GetElement()
	//donePayOrder := GetServer().GetSDKMgr().DonePayOrder(payOrderID, team.Name, clientObj.userID, team.Ticket)
	//if loger.CheckFail("donePayOrder==true", donePayOrder == true, donePayOrder, true) {
	//	return false //订单号不存在
	//}
	resetAttrib.Save()
	team.CheckPayAward(0)
	return true
}

type VipAccpetMonthCardAwardMsg struct { ///客户端请求领取vip每日月卡礼包
	MsgHead `json:"head"` // "vipshop", "accpetmonthcardaward"
}

func (self *VipAccpetMonthCardAwardMsg) GetTypeAndAction() (string, string) {
	return "vipshop", "accpetmonthcardaward"
}

func (self *VipAccpetMonthCardAwardMsg) checkAction(client IClient) (result bool) { ///验货
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeMonthCard)
	if loger.CheckFail("resetAttrib!=nil", resetAttrib != nil, resetAttrib, nil) {
		return false //找不到vip礼包领取信息
	}
	isDeadLine := resetAttrib.Value1 <= Now()
	if loger.CheckFail("isDeadLine!=true", isDeadLine != true, isDeadLine, true) {
		return false //过期月卡
	}
	if loger.CheckFail("resetAttrib.Value1<=0", resetAttrib.Value2 <= 0, resetAttrib.Value2, 0) {
		return false //已领取过了不能再取
	}
	return true
}

func (self *VipAccpetMonthCardAwardMsg) payAction(client IClient) (result bool) { ///收钱
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeMonthCard)
	resetAttrib.Value2 += 1 ///记录领取次数
	return true
}

func (self *VipAccpetMonthCardAwardMsg) doAction(client IClient) (result bool) { ///发货
	team := client.GetTeam()
	//	vipShopMgr := team.GetVipShopMgr()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeMonthCard)
	moneyType := GetServer().GetStaticDataMgr().GetMoneyType(MonthCardMoneyID)
	team.AwardObject(awardTypeTicket, moneyType.Diamonds, 0, 0)
	vipQueryMonthCardResultMsg := new(VipQueryMonthCardResultMsg)
	vipQueryMonthCardResultMsg.DeadLine = resetAttrib.Value1
	vipQueryMonthCardResultMsg.HasAccpetCount = resetAttrib.Value2
	client.SendMsg(vipQueryMonthCardResultMsg)
	client.RechargeRecord(Get_MonthCard, moneyType.Diamonds)
	return true
}

func (self *VipAccpetMonthCardAwardMsg) processAction(client IClient) (result bool) { ///领取vip每日礼包
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
