package football

import (
	"math"
)

type StarHandler struct {
	MsgHandler
}

///球员消息处理器
func (self *StarHandler) getName() string { ///返回可处理的消息类型
	return "star"
}

func (self *StarHandler) initHandler() { ///初始化消息处理器
	self.addActionToList(new(StarEducationMsg))
	self.addActionToList(new(StarEvolveMsg))
	//self.addActionToList(new(StarRenewalMsg))
	self.addActionToList(new(StarSackMsg))
}

type StarEvolveMsg struct { ///请求球员突破消息
	MsgHead    `json:"head"` ///"star", "starevolve"
	StarID     int           `json:"starid"`   ///请求突破球员id
	PayItemID1 int           `json:"payitem1"` ///扣除道具ID1
	PayItemID2 int           `json:"payitem2"` ///扣除道具ID2
}

func (self *StarEvolveMsg) GetTypeAndAction() (string, string) {
	return "star", "starevolve"
}

func (self *StarEvolveMsg) checkMaxLimitAndNeedItem(client IClient) bool {
	team := client.GetTeam()
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	loger := GetServer().GetLoger()
	itemMgr := team.GetItemMgr()
	star := team.GetStar(self.StarID)
	starInfo := star.GetInfo()
	levelExpTypeStaticData := staticDataMgr.GetLevelExpType(levelExpTypeStarEvolve, starInfo.EvolveCount)

	evolveCountMax := star.GetEvolveCountMax()                                                             ///升星上限
	if loger.CheckFail("EvolveCount < evolveCountMax", starInfo.EvolveCount < evolveCountMax, star, nil) { ///升星等级必须小于上限
		return false
	}

	needItem, insteadItem := star.GetEvolveNeedItem()  ///取得升星所需道具和替代道具
	needCount := levelExpTypeStaticData.NeedItemCount1 ///需求的升星道具数量

	if needCount == 0 { ///如果需求的道具数量为0就放过
		return true
	}

	needItemCount := itemMgr.GetItemCountByType(needItem) ///取背包中需求道具的数量
	if needItemCount >= needCount {                       ///如果需求的道具足够就通过检验
		return true
	}

	if star.IsMannaStar == 1 {
		return false //!万能碎片不适用于天赐球员
	}

	remainNeedCount := needCount - needItemCount                ///仍然需要的数量
	insteadItemcount := itemMgr.GetItemCountByType(insteadItem) ///取背包中替代道具的数量

	if loger.CheckFail("Item Enough", insteadItemcount >= remainNeedCount, star, nil) { ///判断替代道具的数量是否足够
		return false
	}

	return true
}

func (self *StarEvolveMsg) checkAction(client IClient) bool { ///检测
	team := client.GetTeam()
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	loger := GetServer().GetLoger()
	//itemMgr := team.GetItemMgr()
	star := team.GetStar(self.StarID)
	//element := GetServer().GetStaticDataMgr().GetStaticData(tableTaskType, 1)
	//itemType := element.(*ItemTypeStaticData)
	//itemType.Body = 0
	if loger.CheckFail("star!=nil", star != nil, star, nil) { ///请求突破的球员必须存在
		return false
	}
	// isReachEvolveLimit := star.IsReachEvolveLimit() //修改: 没有突破颜色限制
	// if loger.CheckFail("isReachEvolveLimit==false", isReachEvolveLimit == false, isReachEvolveLimit, false) { ///请求突破的球员必须存在
	// 	return false
	// }
	starInfo := star.GetInfo()
	levelExpTypeStaticData := staticDataMgr.GetLevelExpType(levelExpTypeStarEvolve, starInfo.EvolveCount)
	if loger.CheckFail("levelTypeStaticData!=nil", levelExpTypeStaticData != nil, levelExpTypeStaticData, nil) {
		return false ///请求突破信息无法在配置数据表中找到
	}
	if loger.CheckFail("starInfo.Level >= levelExpTypeStaticData.NeedLevel",
		starInfo.Level >= levelExpTypeStaticData.NeedLevel,
		starInfo.Level, levelExpTypeStaticData.NeedLevel) {
		return false ///未到达可突破的等级
	}

	managerLevel := team.GetLevel()
	if loger.CheckFail("managerLever >= levelExpTypeStaticData.NeedLevel",
		managerLevel >= levelExpTypeStaticData.NeedTeamLevel,
		managerLevel, levelExpTypeStaticData.NeedTeamLevel) {
		return false ///经理等级未达到可突破等级
	}

	currentCoin := team.GetCoin()
	if loger.CheckFail("currentCoin>=levelTypeStaticData.PayCoin", currentCoin >= levelExpTypeStaticData.PayCoin,
		currentCoin, levelExpTypeStaticData.PayCoin) {
		return false ///球币不够扣,余额不足
	}

	///检测升星上限和道具需求
	checkLimitAndItemRes := self.checkMaxLimitAndNeedItem(client)
	if !checkLimitAndItemRes {
		return false
	}

	//if levelExpTypeStaticData.NeedItemType1 > 0 {
	//	payItem1 := itemMgr.GetItem(self.PayItemID1)
	//	if loger.CheckFail("payItem1!=nil", payItem1 != nil, payItem1, nil) {
	//		return false ///需扣的道具1不存在
	//	}
	//	payItemInfo1 := payItem1.GetInfo()
	//	if loger.CheckFail("payItemInfo1.Type == levelExpTypeStaticData.NeedItemType1",
	//		payItemInfo1.Type == levelExpTypeStaticData.NeedItemType1,
	//		payItemInfo1.Type, levelExpTypeStaticData.NeedItemType1) {
	//		return false ///需扣的道具1类型与静态配置表中的类型不匹配
	//	}
	//	if loger.CheckFail("payItemInfo1.Count>= levelExpTypeStaticData.NeedItemCount1",
	//		payItemInfo1.Count >= levelExpTypeStaticData.NeedItemCount1,
	//		payItemInfo1.Count, levelExpTypeStaticData.NeedItemCount1) {
	//		return false ///需扣的道具1库存数量不够扣除
	//	}
	//	//hasEnoughItem1 := itemMgr.HasEnoughItem(levelExpTypeStaticData.NeedItemType1,
	//	//	levelExpTypeStaticData.NeedItemCount1)
	//	//if loger.CheckFail("hasEnoughItem1==true", hasEnoughItem1 == true, hasEnoughItem1, true) {
	//	//	return false ///需扣的道具1数量不够扣
	//	//}
	//}
	//if levelExpTypeStaticData.NeedItemType2 > 0 {
	//	payItem2 := itemMgr.GetItem(self.PayItemID2)
	//	if loger.CheckFail("payItem2!=nil", payItem2 != nil, payItem2, nil) {
	//		return false ///需扣的道具2不存在
	//	}
	//	payItemInfo2 := payItem2.GetInfo()
	//	if loger.CheckFail("payItemInfo2.Type == levelExpTypeStaticData.NeedItemType2",
	//		payItemInfo2.Type == levelExpTypeStaticData.NeedItemType2,
	//		payItemInfo2.Type, levelExpTypeStaticData.NeedItemType2) {
	//		return false ///需扣的道具2类型与静态配置表中的类型不匹配
	//	}
	//	if loger.CheckFail("payItemInfo2.Count>= levelExpTypeStaticData.NeedItemCount2",
	//		payItemInfo2.Count >= levelExpTypeStaticData.NeedItemCount2,
	//		payItemInfo2.Count, levelExpTypeStaticData.NeedItemCount2) {
	//		return false ///需扣的道具2库存数量不够扣除
	//	}
	//	//hasEnoughItem2 := itemMgr.HasEnoughItem(levelExpTypeStaticData.NeedItemType2,
	//	//	levelExpTypeStaticData.NeedItemCount2)
	//	//if loger.CheckFail("hasEnoughItem1==true", hasEnoughItem2 == true, hasEnoughItem2, true) {
	//	//	return false ///需扣的道具2数量不够扣
	//	//}
	//}
	return true
}

func (self *StarEvolveMsg) payAction(client IClient) bool { ///支付
	team := client.GetTeam()
	star := team.GetStar(self.StarID)
	syncMgr := client.GetSyncMgr()
	starInfo := star.GetInfo()
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	itemMgr := team.GetItemMgr()
	removeItemIDList := IntList{}    ///同步删除道具id列表
	syncItemList := SyncObjectList{} ///同步道具属性变更列表
	levelTypeStaticData := staticDataMgr.GetLevelExpType(levelExpTypeStarEvolve, starInfo.EvolveCount)
	if levelTypeStaticData.PayCoin > 0 {
		team.PayCoin(levelTypeStaticData.PayCoin) ///扣突破所需要的球币
		syncMgr.SyncObject("StarEvolveMsg", team) ///同步最新的球队属性到客户端
	}
	//if levelTypeStaticData.NeedItemType1 > 0 { ///扣道具1
	//	itemMgr.PayItem(self.PayItemID1, levelTypeStaticData.NeedItemCount1) ///扣除道具1
	//	item1 := itemMgr.GetItem(self.PayItemID1)
	//	if item1 != nil {
	//		syncItemList = append(syncItemList, item1)
	//	} else {
	//		removeItemIDList = append(removeItemIDList, self.PayItemID1)
	//	}
	//	//removeItemIDList1, _ := itemMgr.PayItemType(levelTypeStaticData.NeedItemType1, levelTypeStaticData.NeedItemCount1)
	//	//removeTotalItemIDList = append(removeTotalItemIDList, removeItemIDList1...)
	//}
	//if levelTypeStaticData.NeedItemType2 > 0 { ///扣道具2
	//	itemMgr.PayItem(self.PayItemID2, levelTypeStaticData.NeedItemCount2) ///扣除道具2
	//	item2 := itemMgr.GetItem(self.PayItemID2)
	//	//removeItemIDList2, _ := itemMgr.PayItemType(levelTypeStaticData.NeedItemType2, levelTypeStaticData.NeedItemCount2)
	//	//removeTotalItemIDList = append(removeTotalItemIDList, removeItemIDList2...)
	//	if item2 != nil {
	//		syncItemList = append(syncItemList, item2)
	//	} else {
	//		removeItemIDList = append(removeItemIDList, self.PayItemID2)
	//	}
	//}

	///扣除球员碎片和万能碎片
	needCount := levelTypeStaticData.NeedItemCount1 ///需求的升星道具数量
	if needCount > 0 {                              ///如果需求的道具数量不为0就扣除道具
		needItem, insteadItem := star.GetEvolveNeedItem()     ///取得升星所需道具和替代道具
		needItemCount := itemMgr.GetItemCountByType(needItem) ///取背包中需求道具的数量
		insteadItemcount := 0                                 ///替代道具数量
		if needItemCount < needCount {                        ///如果需求的道具不够就用替代道具来顶
			insteadItemcount = needCount - needItemCount ///需要的替代道具数量，前面检测过，这里数量肯定是够的
		} else {
			needItemCount = needCount
		}
		///扣除需求道具
		if needItemCount > 0 {
			removeList, updatItemID, _ := itemMgr.PayItemType(needItem, needItemCount) ///扣除道具
			removeItemIDList = append(removeItemIDList, removeList...)                 ///需要删除的道具
			updatItem := itemMgr.GetItem(updatItemID)
			if updatItem != nil {
				syncItemList = append(syncItemList, updatItem) ///需要更新的道具
			}
		}
		///扣除替代道具
		if insteadItemcount > 0 {
			removeList, updatItemID, _ := itemMgr.PayItemType(insteadItem, insteadItemcount) ///扣除道具
			removeItemIDList = append(removeItemIDList, removeList...)                       ///需要删除的道具
			updatItem := itemMgr.GetItem(updatItemID)
			if updatItem != nil {
				syncItemList = append(syncItemList, updatItem) ///需要更新的道具
			}
		}
	}

	if len(syncItemList) > 0 { ///有道具数量属性变更的需要同步客户端
		syncMgr.SyncObjectArray("StarEvolveMsg", syncItemList)
	}
	if removeItemIDList.Len() > 0 { ///有道具删除消息需要同步客户端
		syncMgr.SyncRemoveItem(removeItemIDList)
	}
	return true
}

func (self *StarEvolveMsg) doAction(client IClient) bool { ///发货
	team := client.GetTeam()
	star := team.GetStar(self.StarID)
	mannaStarMgr := team.GetMannaStarMgr()
	starInfo := star.GetInfo()
	starInfo.EvolveCount++ ///突破等级提升一级

	///根据突破等级决定球员品质
	switch starInfo.EvolveCount {
	case 3:
		starInfo.Grade = starGradeBlue
	case 5:
		starInfo.Grade = starGradePurple
	case 7:
		starInfo.Grade = starGradeOrange
	case 9:
		starInfo.Grade = starGradeRed
	}
	star.CalcScore()
	team.CalcScore()
	client.GetSyncMgr().SyncObject("StarEvolveMsg", team)
	///同步最新的球员信息到客户端
	client.GetSyncMgr().SyncObject("StarEvolveMsg", star)

	//!更新自创球员课题位置
	if star.IsMannaStar == 1 {
		mannaStarMgr.UpdateMannaStarSeat()
	}

	return true
}

func (self *StarEvolveMsg) processAction(client IClient) bool {
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

type StarEducationMsg struct { ///球员培养消息
	MsgHead        `json:"head"` ///"star", "stareducation"
	StarID         int           `json:"starid"`         ///培训球星id
	PassAdd        int           `json:"passadd"`        ///传球潜力加点
	StealsAdd      int           `json:"stealsadd"`      ///抢断潜力加点
	DribblingAdd   int           `json:"dribblingadd"`   ///盘带潜力加点
	SlidingAdd     int           `json:"slidingadd"`     ///铲球潜力加点
	ShootingAdd    int           `json:"shootingadd"`    ///射门潜力加点
	GoalKeepingAdd int           `json:"goalkeepingadd"` ///守门潜力加点
}

func (self *StarEducationMsg) GetTypeAndAction() (string, string) {
	return "star", "stareducation"
}

func (self *StarEducationMsg) GetAddPointMax(addPointType int, star *Star) (int, int) { ///得到球员某一属性的加点上限
	//	starInfo := star.GetInfo()
	if star.IsMannaStar == 1 {

		mannaStarMgr := star.team.GetMannaStarMgr()
		starType := mannaStarMgr.GetMannaStar(star.Type)
		//	evolveCount := float32(starInfo.EvolveCount)
		currentValue, growValue := float32(0), float32(0)
		switch addPointType {
		case starAttribPass:
			currentValue, growValue = float32(starType.Pass), float32(starType.PassGrow)
		case starAttribSteals:
			currentValue, growValue = float32(starType.Steals), float32(starType.StealsGrow)
		case starAttribDribbling:
			currentValue, growValue = float32(starType.Dribbling), float32(starType.DribblingGrow)
		case starAttribSliding:
			currentValue, growValue = float32(starType.Sliding), float32(starType.SlidingGrow)
		case starAttribShooting:
			currentValue, growValue = float32(starType.Shooting), float32(starType.ShootingGrow)
		case starAttribGoalKeeping:
			currentValue, growValue = float32(starType.GoalKeeping), float32(starType.GoalKeepingGrow)
		}
		//addPointMax := currentValue + (currentValue+growValue)*(((1+evolveCount)*evolveCount/2)*0.08+0.2*evolveCount) + growValue
		evolveCount := float32(star.GetInfo().EvolveCount)
		//addPointMax := float32(1+evolveCount*(evolveCount+1)/2) * ((growValue*growValue)/100 + growValue - 20)
		//INT((星级*(星级+1)/2)*(成长值^2/90+0.5*成长值-30))
		//addPointMax := (evolveCount * (evolveCount + 1.0) / 2.0) * (growValue*growValue/90.0 + 0.5*growValue - 30.0)

		addPointMax := star.GetAddPointMax(evolveCount, growValue)
		if addPointMax < 0.0 {
			addPointMax = 0.0
		}

		return int(addPointMax + currentValue), int(currentValue)
	}
	starType := star.GetTypeInfo()
	//	evolveCount := float32(starInfo.EvolveCount)
	currentValue, growValue := float32(0), float32(0)
	switch addPointType {
	case starAttribPass:
		currentValue, growValue = float32(starType.Pass), float32(starType.PassGrow)
	case starAttribSteals:
		currentValue, growValue = float32(starType.Steals), float32(starType.StealsGrow)
	case starAttribDribbling:
		currentValue, growValue = float32(starType.Dribbling), float32(starType.DribblingGrow)
	case starAttribSliding:
		currentValue, growValue = float32(starType.Sliding), float32(starType.SlidingGrow)
	case starAttribShooting:
		currentValue, growValue = float32(starType.Shooting), float32(starType.ShootingGrow)
	case starAttribGoalKeeping:
		currentValue, growValue = float32(starType.GoalKeeping), float32(starType.GoalKeepingGrow)
	}
	//addPointMax := currentValue + (currentValue+growValue)*(((1+evolveCount)*evolveCount/2)*0.08+0.2*evolveCount) + growValue
	evolveCount := float32(star.GetInfo().EvolveCount)
	//addPointMax := float32(1+evolveCount*(evolveCount+1)/2) * ((growValue*growValue)/100 + growValue - 20)
	//INT((星级*(星级+1)/2)*(成长值^2/90+0.5*成长值-30))
	//addPointMax := (evolveCount * (evolveCount + 1.0) / 2.0) * (growValue*growValue/90.0 + 0.5*growValue - 30.0)

	addPointMax := star.GetAddPointMax(evolveCount, growValue)
	if addPointMax < 0.0 {
		addPointMax = 0.0
	}

	return int(addPointMax + currentValue), int(currentValue)
}

///计算总的支付培养点数
func (self *StarEducationMsg) calcTotalPayTalentPoint(client IClient, currentPointList IntList, addPointList IntList) int {
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	star := team.GetStar(self.StarID)
	//	starInfo := star.GetInfo()
	totalPayPoint := 0
	//staticDataMgr := GetServer().GetStaticDataMgr()
	//levelExpCount := staticDataMgr.GetLevelExpCount(levelExpTypeStarEducationLevel)
	//evolveLevelExpType := staticDataMgr.Unsafe().GetLevelExpType(levelExpTypeStarEvolve, starInfo.EvolveCount)
	//if loger.CheckFail("evolveLevelExpType!=nil", evolveLevelExpType != nil,
	//	evolveLevelExpType, nil) { ///配置表中必须有突破等级数据
	//	return 0
	//}
	//if loger.CheckFail("levelExpCount>0", levelExpCount > 0,
	//	levelExpCount, 0) { ///配置表中必须有加点经验配置表
	//	return 0
	//}
	if len(currentPointList) != len(addPointList) {
		return 0 ///输入的两个数组必须长度一致
	}
	for k := range addPointList {
		addPoint := addPointList[k]
		if addPoint <= 0 {
			continue
		}
		currentPoint := currentPointList[k]
		goalPoint := currentPoint + addPoint                               // 培养后属性值
		currentAddPointMax, staticAttrib := self.GetAddPointMax(k+1, star) ///得到此属性的最大加点值和静态属性值
		///fmt.Println(k+1, currentAddPointMax)
		//loger.Print("上限为: %d\r\n", currentAddPointMax)
		if loger.CheckFail("(staticAttrib+goalPoint) <= currentAddPointMax", (staticAttrib+goalPoint) <= currentAddPointMax,
			staticAttrib+goalPoint, currentAddPointMax) {
			return 0 ///客户端请求加点值必须在突破等级限制范围内
		}
		var value float64 = 0
		for i := currentPoint + 1; i <= goalPoint; i++ {
			value = float64(staticAttrib+i) / 5.0
			needExp := 10 + int(math.Ceil(value))

			totalPayPoint += needExp
		}

	}

	return totalPayPoint
}

func (self *StarEducationMsg) checkAction(client IClient) bool { ///检测条件
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	starInfoCopy := team.GetStarInfoCopy(self.StarID)
	currentPointList := IntList{starInfoCopy.PassTalentAdd, starInfoCopy.StealsTalentAdd, starInfoCopy.DribblingTalentAdd,
		starInfoCopy.SlidingTalentAdd, starInfoCopy.ShootingTalentAdd, starInfoCopy.GoalKeepingTalentAdd}
	addPointList := IntList{self.PassAdd, self.StealsAdd, self.DribblingAdd, self.SlidingAdd, self.ShootingAdd, self.GoalKeepingAdd}
	totalPayTalentPoint := self.calcTotalPayTalentPoint(client, currentPointList, addPointList)
	if loger.CheckFail("totalPayTalentPoint>0", totalPayTalentPoint > 0, totalPayTalentPoint, 0) {
		return false ///加培养点数时必须要有代价
	}
	//	teamTalentPoint := team.GetInfo().TalentPoint ///验证余额要足够扣
	teamTalentPoint := team.GetInfo().Coin
	if loger.CheckFail("totalPayTalentPoint<teamTalentPoint", totalPayTalentPoint < teamTalentPoint,
		totalPayTalentPoint, teamTalentPoint) {
		return false
	}
	return true
}

func (self *StarEducationMsg) payAction(client IClient) bool { ///支付代价
	team := client.GetTeam()
	starInfoCopy := team.GetStarInfoCopy(self.StarID)
	currentPointList := IntList{starInfoCopy.PassTalentAdd, starInfoCopy.StealsTalentAdd, starInfoCopy.DribblingTalentAdd,
		starInfoCopy.SlidingTalentAdd, starInfoCopy.ShootingTalentAdd, starInfoCopy.GoalKeepingTalentAdd}
	addPointList := IntList{self.PassAdd, self.StealsAdd, self.DribblingAdd, self.SlidingAdd, self.ShootingAdd, self.GoalKeepingAdd}
	totalPayTalentPoint := self.calcTotalPayTalentPoint(client, currentPointList, addPointList)
	team.PayTalentPoint(totalPayTalentPoint)
	star := team.GetStar(self.StarID)
	star.AddTotalPayTalentPoint(totalPayTalentPoint) ///记录新的消费点,便于以后洗点返还
	_, msgAction := self.GetTypeAndAction()
	client.GetSyncMgr().SyncObject(msgAction, team) ///通知客户端潜力点变更了
	return true
}

func (self *StarEducationMsg) doAction(client IClient) bool { ///发货
	team := client.GetTeam()
	star := team.GetStar(self.StarID)
	star.AddTalentPoint(self.PassAdd, self.StealsAdd, self.DribblingAdd, self.SlidingAdd,
		self.ShootingAdd, self.GoalKeepingAdd) ///增加球员培养加成点数
	_, msgAction := self.GetTypeAndAction()
	star.CalcScore()
	team.CalcScore()
	client.GetSyncMgr().SyncObject(msgAction, team)
	client.GetSyncMgr().SyncObject(msgAction, star) ///同步最新的球员属性给客户端

	///更新天天联赛中的日常任务
	team.GetTaskMgr().UpdateDayTaskFuntion(client.GetElement(), dayTaskFunctionStarEducation)
	return true
}

func (self *StarEducationMsg) processAction(client IClient) bool {
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

//type StarRenewalMsg struct { ///请求球员续约消息
//	MsgHead `json:"head"` ///"star", "renewal"
//	StarID  int           `json:"starid"` ///续约球员id
//	//	AddContractPoint int           `json:"addcontractpoint"` ///申请续约点数
//	UseTicket bool `json:"useticket"` ///是否使用球票续约
//}

//func (self *StarRenewalMsg) GetTypeAndAction() (string, string) {
//	return "star", "renewal"
//}

//func (self *StarRenewalMsg) checkAction(client IClient) bool { ///验货
//	staticDataMgr := GetServer().GetStaticDataMgr()
//	loger := GetServer().GetLoger()
//	team := client.GetTeam()
//	teamInfo := team.GetInfo()
//	star := team.GetStar(self.StarID)
//	starInfo := star.GetInfo()
//	if loger.CheckFail(" star!=nil", star != nil, star, nil) {
//		return false ///请求解雇的球员不存在
//	}
//	///得到球票配置价格
//	configTicketPrice := staticDataMgr.GetConfigStaticDataInt(configStar, configItemStarCommonConfig, 4)
//	if loger.CheckFail("configTicketPrice>0", configTicketPrice > 0, configTicketPrice, 0) {
//		return false ///球票配置价格必须大于０
//	}
//	///得到球币配置价格
//	configCoinPrice := staticDataMgr.GetConfigStaticDataInt(configStar, configItemStarCommonConfig, 3)
//	if loger.CheckFail("configCoinPrice>0", configCoinPrice > 0, configCoinPrice, 0) {
//		return false ///球币配置价格必须大于０
//	}
//	if true == self.UseTicket {
//		needPayTicket := starInfo.Grade * configTicketPrice
//		if loger.CheckFail("teamInfo.Ticket>=needPayTicket", teamInfo.Ticket >= needPayTicket,
//			teamInfo.Ticket, needPayTicket) {
//			return false ///球队球票余额要足够扣
//		}
//	} else {
//		needAddContractPoint := ContractPointMax - starInfo.ContractPoint ///得到补满所需合约点
//		if loger.CheckFail("needAddContractPoint", needAddContractPoint > 0, needAddContractPoint, 0) {
//			return false ///补满所需合约点必须大于0
//		}
//		needPayCoin := needAddContractPoint * configCoinPrice
//		if loger.CheckFail("teamInfo.Coin>= needPayCoin", teamInfo.Coin >= needPayCoin,
//			teamInfo.Coin, needPayCoin) {
//			return false ///球队球币余额要足够扣
//		}
//	}
//	return true
//}

//func (self *StarRenewalMsg) payAction(client IClient) bool { ///付款
//	staticDataMgr := GetServer().GetStaticDataMgr()
//	team := client.GetTeam()
//	syncMgr := client.GetSyncMgr()
//	star := team.GetStar(self.StarID)
//	starInfo := star.GetInfo()
//	///得到球票配置价格
//	configTicketPrice := staticDataMgr.GetConfigStaticDataInt(configStar, configItemStarCommonConfig, 4)
//	///得到球币配置价格
//	configCoinPrice := staticDataMgr.GetConfigStaticDataInt(configStar, configItemStarCommonConfig, 3)
//	if true == self.UseTicket {
//		///球票续约永久球员
//		renewalTicketPrice := configTicketPrice * starInfo.Grade
//		team.SpendTicket(renewalTicketPrice, "StarRenewalMsg")
//		syncMgr.SyncObject("StarRenewalMsg", team) ///同步最新的球队球币信息
//	} else {
//		needAddContractPoint := ContractPointMax - starInfo.ContractPoint ///得到补满所需合约点
//		renewalPriceCoin := configCoinPrice * needAddContractPoint
//		team.PayCoin(renewalPriceCoin)
//		syncMgr.SyncObject("StarRenewalMsg", team) ///同步最新的球队球币信息
//	}
//	return true
//}

//func (self *StarRenewalMsg) doAction(client IClient) bool { ///发货
//	team := client.GetTeam()
//	star := team.GetStar(self.StarID)
//	starInfo := star.GetInfo()
//	if true == self.UseTicket {
//		starInfo.ContractPoint = LifeTimeContractPoint ///升为永久球员
//	} else {
//		///加契约点
//		needAddContractPoint := ContractPointMax - starInfo.ContractPoint ///得到补满所需合约点
//		starInfo.ContractPoint = Min(starInfo.ContractPoint+needAddContractPoint, ContractPointMax)
//	}
//	client.GetSyncMgr().SyncObject("StarRenewalMsg", star) ///同步最新的球员信息给客户端
//	return true
//}

//func (self *StarRenewalMsg) processAction(client IClient) bool {
//	if self.checkAction(client) == false { ///检测
//		return false
//	}
//	if self.payAction(client) == false { ///支付
//		return false
//	}
//	if self.doAction(client) == false { ///发货
//		return false
//	}
//	return true
//}

type StarSackMsg struct { ///请求解雇球员消息
	MsgHead   `json:"head"` ///"star", "sack"
	StarID    int           `json:"starid"`    ///解雇球员id
	UseTicket bool          `json:"useticket"` ///是否使用球票解雇
}

func (self *StarSackMsg) GetTypeAndAction() (string, string) {
	return "star", "sack"
}

func (self *StarSackMsg) checkAction(client IClient) bool { ///验货
	team := client.GetTeam()
	//vipLevel := team.GetInfo().VipLevel
	loger := GetServer().GetLoger()
	staticDataMgr := GetServer().GetStaticDataMgr()
	payStarSackTicket := staticDataMgr.GetConfigStaticDataInt(configStarCenter,
		configItemStarCenterCommonConfig, 2)

	///根据VIP等级得到训练点与经验返还百分比
	// vipLevel := team.GetVipLevel()
	// vipPrivilege := GetServer().GetStaticDataMgr().GetVipInfo(vipLevel)
	// if loger.CheckFail("vipPrivilege != nil", vipPrivilege != nil, vipPrivilege, nil) {
	// 	return false ///VIP数据为空
	// }

	//	trainPointSackRate := vipPrivilege.Param12
	//	expSackRate := vipPrivilege.Param13

	// if loger.CheckFail("trainPointSackRate > 0", trainPointSackRate > 0, trainPointSackRate, 0) {
	// 	return false ///解雇球员训练点返回率必须大于0
	// }

	// if loger.CheckFail("expSackRate > 0", expSackRate > 0, expSackRate, 0) {
	// 	return false ///解雇球员经验返回率必须大于0
	// }
	NormalSackRate := staticDataMgr.GetConfigStaticDataInt(configStar,
		configItemStarCommonConfig, 5) ///普通解雇球员返还率
	TicketSackRate := staticDataMgr.GetConfigStaticDataInt(configStar,
		configItemStarCommonConfig, 6) ///球票解雇球员返还率
	if loger.CheckFail("NormalSackRate>0", NormalSackRate > 0, NormalSackRate, 0) {
		return false ///普通解雇球员训练点返回率必须大于0
	}
	if loger.CheckFail(" TicketSackRate>0", TicketSackRate > 0, TicketSackRate, 0) {
		return false ///球票解雇球员训练点返回率必须大于0
	}
	currentTicket := team.GetTicket()
	isStarInCurrentFormation := team.IsStarInCurrentFormation(self.StarID)
	star := team.GetStar(self.StarID)
	if loger.CheckFail(" star!=nil", star != nil, star, nil) {
		return false ///请求解雇的球员不存在
	}

	if loger.CheckFail("ismannastar == 0", star.IsMannaStar == 0, star.IsMannaStar, 0) {
		return false //! 自创球员无法解雇
	}

	if loger.CheckFail("isStarInCurrentFormation==false", isStarInCurrentFormation == false,
		isStarInCurrentFormation, false) {
		return false ///上阵球员不得解雇
	}
	if true == self.UseTicket {
		if loger.CheckFail("payStarSackTicket>0", payStarSackTicket > 0, payStarSackTicket, 0) {
			return false ///配置球票余额必须大于０
		}
		if loger.CheckFail("currentTicket>=payStarSackTicket", currentTicket >= payStarSackTicket,
			currentTicket, payStarSackTicket) {
			return false ///球票余额不足
		}
	}

	processMgr := client.GetTeam().GetProcessMgr()
	isStarTrain := processMgr.FindProcessByObjID(ProcessTypeStarTrain, self.StarID) > 0
	if loger.CheckFail("isStarTrain == true", isStarTrain == false, isStarTrain, false) {
		return false // 球员训练中,禁止解雇
	}

	return true
}

func (self *StarSackMsg) payAction(client IClient) bool { ///付款
	team := client.GetTeam()
	staticDataMgr := GetServer().GetStaticDataMgr()
	syncMgr := client.GetSyncMgr()
	payStarSackTicket := staticDataMgr.GetConfigStaticDataInt(configStarCenter,
		configItemStarCenterCommonConfig, 2)
	if true == self.UseTicket {
		team.PayTicket(payStarSackTicket)
		client.SetMoneyRecord(PlayerCostMoney, Pay_StarSack, payStarSackTicket, team.GetTicket())
		syncMgr.SyncObject("StarSackMsg", team)
	}
	return true
}

func (self *StarSackMsg) doAction(client IClient) bool { ///发货
	team := client.GetTeam()
	syncMgr := client.GetSyncMgr()
	itemMgr := team.GetItemMgr()
	//	staticDataMgr := GetServer().GetStaticDataMgr()
	star := team.GetStar(self.StarID)
	starInfo := star.GetInfo()

	if star.IsMannaStar == 1 {
		mannaStarMgr := team.GetMannaStarMgr()
		starType := mannaStarMgr.GetMannaStar(star.Type)
		sackTalentRate, sackExpRate := team.GetVipStarSackRepayRate()

		returnPayTalentPoint := starInfo.TotalPayTalentPoint * sackTalentRate / 100
		returnExp := starInfo.Exp * sackExpRate / 100

		starCardCount := team.GetMannaStarCardCount(starType, star.EvolveCount)

		team.AwardTalentPoint(returnPayTalentPoint)         ///返回训练点数
		team.AwardExpPool(returnExp)                        ///返还解雇经验
		team.AwardObject(ItemStarCard, starCardCount, 0, 0) ///赠星卡
	} else {
		starType := star.GetTypeInfo()
		sackTalentRate, sackExpRate := team.GetVipStarSackRepayRate()

		returnPayTalentPoint := starInfo.TotalPayTalentPoint * sackTalentRate / 100
		returnExp := starInfo.Exp * sackExpRate / 100

		starCardCount := team.GetStarCardCount(starType, star.EvolveCount)

		team.AwardTalentPoint(returnPayTalentPoint)         ///返回训练点数
		team.AwardExpPool(returnExp)                        ///返还解雇经验
		team.AwardObject(ItemStarCard, starCardCount, 0, 0) ///赠星卡
	}

	///根据VIP等级得到训练点与经验返还百分比
	//	vipLevel := team.GetVipLevel()
	// vipPrivilege := staticDataMgr.GetVipInfo(vipLevel)

	// trainPointSackRate := vipPrivilege.Param12
	// expSackRate := vipPrivilege.Param13

	//	NormalSackRate := staticDataMgr.GetConfigStaticDataInt(configStar,
	//		configItemStarCommonConfig, 5) ///普通解雇球员返还率
	// TicketSackRate := staticDataMgr.GetConfigStaticDataInt(configStar,
	// 	configItemStarCommonConfig, 6) ///球票解雇球员返还率

	///计算取得返回训练点数
	// returnPayTalentPoint := starInfo.TotalPayTalentPoint * trainPointSackRate / 100
	// returnPayExp := starInfo.Exp * expSackRate / 100

	itemMgr.RemoveEquipment(client, self.StarID) ///卸下该球员装备
	team.RemoveStar(IntList{self.StarID})        ///删除此球员
	syncMgr.syncRemoveStar(IntList{self.StarID}) ///同步客户端此球员被删除

	team.CalcScore() ///更新球队评分
	syncMgr.SyncObject("StarSackMsg", team)
	return true
}

func (self *StarSackMsg) processAction(client IClient) bool {
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
