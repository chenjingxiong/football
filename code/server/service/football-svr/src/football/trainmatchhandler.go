package football

import (
	"fmt"
	"strconv"
	"strings"
)

const ( ///训练项目类别
	trainTypePass        = 1 ///传球
	trainTypeSteals      = 2 ///抢断
	trainTypebDribbling  = 3 ///盘带
	trainTypeSliding     = 4 ///铲球
	trainTypeShooting    = 5 ///射门
	trainTypeGoalKeeping = 6 ///守门
	trainTypeBody        = 7 ///身体值
	trainTypeSpeed       = 8 ///速度
)

const (
	trainGradeBegin  = 1 ///开始符
	trainGradeGreen  = 1 ///绿色品质训练
	trainGradeBlue   = 2 ///蓝色品质训练
	trainGradePurple = 3 ///紫色品质训练
	trainGradeOrange = 4 ///橙色品质训练
	trainGradeRed    = 5 ///红色品质训练
	trainGradeEnd        ///结束符
)

type PlayTrainMatchMsg struct { ///客户端请求打训练赛
	MsgHead    `json:"head"` ///"trainmatch", "queryinfo"
	TrainIndex int           `json:"trainindex"` ///训练类型索引,从1开始
}

func (self *PlayTrainMatchMsg) GetTypeAndAction() (string, string) {
	return "trainmatch", "playmatch" ///打比赛
}

func (self *PlayTrainMatchMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	loger := GetServer().GetLoger()
	staticDataMgr := GetServer().GetStaticDataMgr()
	awardTrainMatchPoint := staticDataMgr.GetConfigStaticDataInt(configTrainMatch,
		configTrainMatchCommon, 4)
	if loger.CheckFail("awardTrainMatchPoint>0", awardTrainMatchPoint > 0, awardTrainMatchPoint, 0) {
		return false ///训练赛全部目标达成时的奖励点数必须大于0
	}

	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamTrainMatch)
	if loger.CheckFail("resetAttrib!=nil", resetAttrib != nil, resetAttrib, nil) {
		return false ///训练赛信息不存在
	}
	_, trainItemList := SeparateIntList(resetAttrib.Value4)
	isTrainIndexLegal := self.TrainIndex >= 1 && self.TrainIndex <= trainItemList.Len()
	if loger.CheckFail(" isTrainIndexLegal==true", isTrainIndexLegal == true, isTrainIndexLegal, false) {
		return false ///操作给的训练项目索引必须合法
	}
	trainType := trainItemList[self.TrainIndex-1] ///从索引转换成训练类型
	isTrainTypeLegal := trainItemList.Search(trainType) >= 0
	if loger.CheckFail("isTrainTypeLegal==true", isTrainTypeLegal == true, isTrainTypeLegal, true) {
		return false ///客户端请求的训练项目是非法的,训练类型不在候选列表中
	}
	//trainTargetList, _ := SeparateIntList(resetAttrib.Value3)
	//isTrainTypeTargetLegal := trainTargetList.Search(trainType) >= 0
	//if loger.CheckFail("isTrainTypeTargetLegal==true", isTrainTypeTargetLegal == true, isTrainTypeTargetLegal, true) {
	//	return false ///客户端请求的训练项目是非法的,训练类型不是训练目标列表中
	//}
	payActionPoint := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTrainMatch, configTrainMatchCommon, 2)
	currentActionPoint := team.GetInfo().ActionPoint
	if loger.CheckFail("currentActionPoint>=payActionPoint", currentActionPoint >= payActionPoint, currentActionPoint, payActionPoint) {
		return false ///球队行动点不足
	}
	return true
}

func (self *PlayTrainMatchMsg) payAction(client IClient) bool {
	payActionPoint := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTrainMatch,
		configTrainMatchCommon, 2)
	team := client.GetTeam()
	team.PayActionPoint(payActionPoint)
	client.GetSyncMgr().SyncObject("PlayTrainMatchMsg", team) ///同步最新球队属性到客户端
	return true
}

func (self *PlayTrainMatchMsg) isTrainTargetComplete(client IClient) bool { ///判断当前的训练目标是否已完成
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamTrainMatch)
	_, doneFlagList := SeparateIntList(resetAttrib.Value3)
	isAllComplete := (doneFlagList.Search(0) < 0) ///找不到0表示全部完成
	return isAllComplete
}

func (self *PlayTrainMatchMsg) finishTrainTarget(client IClient) { ///处理训练完成时的逻辑处理
	team := client.GetTeam()
	awardTrainMatchPoint := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTrainMatch,
		configTrainMatchCommon, 4)
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamTrainMatch)
	resetAttrib.Value1 += awardTrainMatchPoint ///加奖励点数
	RefreshTrainMatchTarget(resetAttrib)       ///刷新下一批目标
}

func (self *PlayTrainMatchMsg) finishTrainItem(client IClient) { ///处理完成单次训练项目时奖励逻辑
	team := client.GetTeam()
	teamLevel := team.GetInfo().Level
	resetAttribMgr := team.GetResetAttribMgr()
	trainAwardType := FindTrainAwardByLevel(teamLevel)
	awardTrainPointList := IntList{trainAwardType.GreenAward, trainAwardType.BlueAward,
		trainAwardType.PurpleAward, trainAwardType.OrangeAward} ///组品质奖励列表
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamTrainMatch)
	trainGradeList, _ := SeparateIntList(resetAttrib.Value4)
	indexTrainItem := self.TrainIndex - 1                // trainTypeList.Search(self.TrainType)
	trainGrade := trainGradeList[indexTrainItem]         ///得到对应训练赛项目品质
	awardTrainPoint := awardTrainPointList[trainGrade-1] ///通过品质取得对应训练点数
	team.AwardTalentPoint(awardTrainPoint)
	client.GetSyncMgr().SyncObject("PlayTrainMatchMsg", team) ///同步最新的球员训练点数
}

func (self *PlayTrainMatchMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamTrainMatch)
	trainTargetString := fmt.Sprintf("%d", resetAttrib.Value3) ///转成串处理
	_, trainTypeList := SeparateIntList(resetAttrib.Value4)
	trainType := trainTypeList[self.TrainIndex-1]                                   ///得到训练类型
	oldString := fmt.Sprintf("%d0", trainType)                                      ///组合旧数据
	newString := fmt.Sprintf("%d1", trainType)                                      ///组合新数据
	trainTargetString = strings.Replace(trainTargetString, oldString, newString, 1) ///只替换一次
	resetAttrib.Value3, _ = strconv.Atoi(trainTargetString)
	self.finishTrainItem(client)                                ///给单次训练奖励
	isTrainTargetComplete := self.isTrainTargetComplete(client) ///判断是否所有训练目标均达成
	if true == isTrainTargetComplete {
		self.finishTrainTarget(client) ///训练赛目标全部达成去处理逻辑
	}
	RefreshTrainMatchType(resetAttrib) ///训练一次也会刷新下一批候选训练类型
	SendQueryTrainMatchResultMsg(client)

	team.GetTaskMgr().UpdateDayTaskFuntion(client.GetElement(), trainType) ///更新训练赛中的日常任务
	return true
}

func (self *PlayTrainMatchMsg) processAction(client IClient) bool {
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

type QueryTrainMatchResultMsg struct { ///查询训练赛信息结果消息
	MsgHead     `json:"head"` ///"trainmatch", "queryresult"
	TotalScore  int           `json:"totalscore"`  ///训练赛已获得总积分,用于领奖
	AcceptScore int           `json:"acceptscore"` ///训练赛中已领过奖的积分
	TrainTarget int           `json:"traintarget"` ///目标信息 x(训练类型)x(是否已完成) 例如203041,可重复
	TrainItem   int           `json:"trainitem"`   ///训练项目列表 x(品质)x(训练类型)12345612 共八位
}

func (self *QueryTrainMatchResultMsg) GetTypeAndAction() (string, string) {
	return "trainmatch", "queryresult"
}

type QueryTrainMatchMsg struct { ///查询训练赛信息消息
	MsgHead `json:"head"` ///"trainmatch", "queryinfo"
}

func (self *QueryTrainMatchMsg) GetTypeAndAction() (string, string) {
	return "trainmatch", "queryinfo"
}

func RefreshTrainMatchTarget(resetAttrib *ResetAttrib) { ///刷新训练赛信息中的目标信息
	targetStr := "" ///目标串
	trainTypeRand := 0
	for i := 1; i <= 3; i++ {
		trainTypeRand = Random(trainTypePass, trainTypeSpeed) ///随机训练目标类型
		targetStr += fmt.Sprintf("%d0", trainTypeRand)
	}
	resetAttrib.Value3, _ = strconv.Atoi(targetStr) ///生成新的目标信息
}

func RefreshTrainMatchType(resetAttrib *ResetAttrib) { ///刷新训练赛信息中的可训练项目类型
	targetStr := "" ///目标串
	trainGradeRand := 0
	trainTypeRand := 0
	for i := 1; i <= 4; i++ {
		trainGradeRand = Random(trainGradeGreen, trainGradeOrange) ///随机训练品质
		trainTypeRand = Random(trainTypePass, trainTypeSpeed)      ///随机训练类型
		targetStr += fmt.Sprintf("%d%d", trainGradeRand, trainTypeRand)
	}

	if resetAttrib.Value5 < 3 {
		targetStr = GuideRefresh(resetAttrib.Value3, resetAttrib.Value5)
		resetAttrib.Value5++
	}
	resetAttrib.Value4, _ = strconv.Atoi(targetStr) ///生成新的目标信息
}

func GuideRefresh(target int, currentTimes int) string { ///新手指引刷新

	//得到训练目标
	targetList := IntList{}
	for i := 0; i < 3; i++ {
		target /= 10
		targetList = append(targetList, target%10)
		target /= 10
	}

	//训练项目为四项,则随机多加一个训练
	targetList = append(targetList, Random(trainTypePass, trainTypeSpeed))
	newtrainStr := ""
	for i := 0; i <= 3; i++ {
		trainGradeRand := Random(trainGradeGreen, trainGradeOrange) ///随机训练品质

		if i == currentTimes {
			trainGradeRand = trainGradeOrange
		}
		newtrainStr += fmt.Sprintf("%d%d", trainGradeRand, targetList[i])
	}
	return newtrainStr
}

///创建默认训练赛信息
func (self *QueryTrainMatchMsg) createDefaultTrainMatchInfo(client IClient) *ResetAttrib {
	team := client.GetTeam()
	refreshTrainMatchHours := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTrainMatch,
		configTrainMatchCommon, 1) ///取得配置表中每日刷新小时数
	resetAttribMgr := team.GetResetAttribMgr()
	resettime := GetHourUTC(refreshTrainMatchHours)
	resetAttribMgr.AddResetAttrib(ResetAttribTypeTeamTrainMatch, resettime, IntList{0})
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamTrainMatch)
	if resetAttrib != nil {
		RefreshTrainMatchTarget(resetAttrib) ///刷新训练赛信息中的目标信息
		RefreshTrainMatchType(resetAttrib)   ///刷新训练赛信息中的可选训练项目信息
	}
	return resetAttrib
}

///发送训练赛查询结果
func SendQueryTrainMatchResultMsg(client IClient) {
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamTrainMatch)
	queryTrainMatchResultMsg := new(QueryTrainMatchResultMsg)
	queryTrainMatchResultMsg.TotalScore = resetAttrib.Value1
	queryTrainMatchResultMsg.AcceptScore = resetAttrib.Value2
	queryTrainMatchResultMsg.TrainTarget = resetAttrib.Value3
	queryTrainMatchResultMsg.TrainItem = resetAttrib.Value4
	client.SendMsg(queryTrainMatchResultMsg) ///通知客户端训练赛信息
}

func (self *QueryTrainMatchMsg) processAction(client IClient) bool { ///通用检测
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamTrainMatch)
	if nil == resetAttrib { ///是否第一次查询训练赛信息
		resetAttrib = self.createDefaultTrainMatchInfo(client)
	}
	if loger.CheckFail("resetAttrib!=nil", resetAttrib != nil, resetAttrib, nil) {
		return false ///创建默认训练赛信息失败
	}
	SendQueryTrainMatchResultMsg(client)
	return true
}

type RefeshTrainMatchMsg struct { ///客户端请求刷新训练赛训练项目列表
	MsgHead `json:"head"` ///"trainmatch", "refreshmatch"
}

func (self *RefeshTrainMatchMsg) GetTypeAndAction() (string, string) {
	return "trainmatch", "refreshmatch" ///刷新训练赛训练项目列表
}

func (self *RefeshTrainMatchMsg) processAction(client IClient) bool {
	team := client.GetTeam()
	loger := GetServer().GetLoger()
	refreshTrainMatchCoinPay := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTrainMatch,
		configTrainMatchCommon, 3) ///刷新训练赛训练项目列表时需要支付的球币代价
	if loger.CheckFail("refreshTrainMatchCoinPay>0", refreshTrainMatchCoinPay > 0, refreshTrainMatchCoinPay, 0) {
		return false ///球币代价未配置或配置为非法值
	}
	currentDiamond := team.GetTicket()
	if loger.CheckFail("refreshTrainMatchCoinPay<=currentCoin", refreshTrainMatchCoinPay <= currentDiamond,
		refreshTrainMatchCoinPay, currentDiamond) {
		return false ///钻石不足支付
	}
	client.SetMoneyRecord(PlayerCostMoney, Pay_RefreshTrainMatch, refreshTrainMatchCoinPay, team.GetTicket())
	team.PayTicket(refreshTrainMatchCoinPay)                    ///支持刷新代价
	client.GetSyncMgr().SyncObject("RefeshTrainMatchMsg", team) ///同步最新的球币信息给客户端
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamTrainMatch)
	if loger.CheckFail("resetAttrib!=nil", resetAttrib != nil, resetAttrib, nil) {
		return false ///训练赛信息有效性验证
	}
	RefreshTrainMatchType(resetAttrib)                                                       ///刷新训练赛可选训练项目列表
	SendQueryTrainMatchResultMsg(client)                                                     ///同步最新的训练赛信息给客户端
	team.GetTaskMgr().UpdateDayTaskFuntion(client.GetElement(), dayTaskFunctionTrainRefresh) ///更新训练赛中的日常任务
	return true
}

type AwardTrainMatchMsg struct { ///客户端请求领取训练赛积分奖励
	MsgHead `json:"head"` ///"trainmatch", "awardmatch"
}

func (self *AwardTrainMatchMsg) GetTypeAndAction() (string, string) {
	return "trainmatch", "awardmatch" ///领取训练赛积分奖励
}

func (self *AwardTrainMatchMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	loger := GetServer().GetLoger()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamTrainMatch)
	if loger.CheckFail("resetAttrib!=nil", resetAttrib != nil, resetAttrib, nil) {
		return false ///训练赛信息有效性验证
	}
	awardNeedScore, awardType, awardObj, awardNum := FindTrainMatchAward(client)
	if loger.CheckFail("awardNeedScore>0", awardNeedScore > 0, awardNeedScore, 0) {
		return false ///奖励需求分数非法
	}
	if loger.CheckFail("awardType>0", awardType > 0, awardType, 0) {
		return false ///奖励类型非法
	}
	if loger.CheckFail("awardObj>0", awardObj > 0, awardObj, 0) {
		return false ///奖励对象非法
	}
	if loger.CheckFail("awardNum>0", awardNum > 0, awardNum, 0) {
		return false ///奖励数量非法
	}

	trainMatchScore := resetAttrib.Value1 ///总训练赛积分
	awardScore := resetAttrib.Value2      ///奖励积分
	if loger.CheckFail("awardScore<trainMatchScore", awardScore < trainMatchScore,
		awardScore, trainMatchScore) {
		return false ///无可领奖项目
	}
	if loger.CheckFail("awardScore<trainMatchScore", awardScore < trainMatchScore,
		awardScore, trainMatchScore) {
		return false ///无可领奖项目
	}
	if TrainMatchAwardTypeItem == awardType {
		///如果是奖励道具还需要验收背包是否已满
		IsStoreFull := team.IsStoreFull(awardObj, awardNum)
		if loger.CheckFail(" IsStoreFull==false", IsStoreFull == false,
			IsStoreFull, false) {
			return false ///无可领奖项目
		}
	}
	return true
}

func FindTrainMatchAward(client IClient) (int, int, int, int) {
	awardNeedScore, awardType, awardObj, awardNum := 0, 0, 0, 0
	team := client.GetTeam()
	teamLevel := team.GetLevel()
	trainAwardType := FindTrainAwardByLevel(teamLevel)
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamTrainMatch)
	awardScore := resetAttrib.Value2
	needScoreList := IntList{trainAwardType.Score1, trainAwardType.Score2, trainAwardType.Score3,
		trainAwardType.Score4, trainAwardType.Score5, trainAwardType.Score6,
		trainAwardType.Score7, trainAwardType.Score8, trainAwardType.Score9}
	awardTypeList := IntList{trainAwardType.AwardType1, trainAwardType.AwardType2, trainAwardType.AwardType3,
		trainAwardType.AwardType4, trainAwardType.AwardType5, trainAwardType.AwardType6,
		trainAwardType.AwardType7, trainAwardType.AwardType8, trainAwardType.AwardType9}
	awardObjList := IntList{trainAwardType.AwardObj1, trainAwardType.AwardObj2, trainAwardType.AwardObj3,
		trainAwardType.AwardObj4, trainAwardType.AwardObj5, trainAwardType.AwardObj6,
		trainAwardType.AwardObj7, trainAwardType.AwardObj8, trainAwardType.AwardObj9}
	awardNumList := IntList{trainAwardType.AwardNum1, trainAwardType.AwardNum2, trainAwardType.AwardNum3,
		trainAwardType.AwardNum4, trainAwardType.AwardNum5, trainAwardType.AwardNum6,
		trainAwardType.AwardNum7, trainAwardType.AwardNum8, trainAwardType.AwardNum9}
	for i := range needScoreList {
		needScore := needScoreList[i]
		if needScore > awardScore {
			awardNeedScore, awardType, awardObj, awardNum = needScore, awardTypeList[i], awardObjList[i], awardNumList[i]
			break
		}
	}
	return awardNeedScore, awardType, awardObj, awardNum
}

const (
	TrainMatchAwardTypeItem = 1 ///奖励道具
	TrainMatchAwardAttrib   = 2 ///奖励货币
	TrainMatchAwardTypeStar = 3 ///奖励球员
)

func (self *AwardTrainMatchMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	//	itemMgr := team.GetItemMgr()
	//	syncMgr := client.GetSyncMgr()
	awardNeedScore, _, awardObj, awardNum := FindTrainMatchAward(client)
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamTrainMatch)
	resetAttrib.Value2 = awardNeedScore ///保存已领奖的档次分数,支付过程
	team.AwardObject(awardObj, awardNum, 0, 0)
	//switch awardType {
	//case TrainMatchAwardTypeItem: ///奖励道具类或道具
	//	newItemID := itemMgr.AwardItem(awardObj, awardNum)
	//	if newItemID > 0 {
	//		client.GetSyncMgr().syncAddItem(IntList{newItemID})
	//	}
	//case TrainMatchAwardAttrib: ///奖励货币
	//	team.AwardAttrib(awardObj, awardNum)
	//	syncMgr.SyncObject("AwardTrainMatchMsg", team)
	//case TrainMatchAwardTypeStar: ///奖励球员
	//	///保留
	//}
	SendQueryTrainMatchResultMsg(client) ///同步最新的训练赛信息给客户端
	return true
}

func (self *AwardTrainMatchMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false { ///检测
		return false
	}
	if self.doAction(client) == false { ///发货
		return false
	}
	return true
}
