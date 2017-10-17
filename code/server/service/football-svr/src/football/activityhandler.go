package football

import (
//"fmt"
//"time"
)

type QueryActivityResultMsg struct { ///查询活动结果信息
	MsgHead      `json:"head"`         ///"activity", "queryresult"
	ActivityInfo `json:"activityinfo"` ///活动信息
}

func (self *QueryActivityResultMsg) GetTypeAndAction() (string, string) {
	return "activity", "queryresult"
}

type QueryActivityInfoMsg struct { ///查询活动信息
	MsgHead      `json:"head"` ///"activity", "queryinfo"
	ActivityType int           `json:"activitytype"` ///活动类型
}

func (self *QueryActivityInfoMsg) GetTypeAndAction() (string, string) {
	return "activity", "queryinfo"
}

func SendQueryActivityResultMsg(activityType int, client IClient) {
	team := client.GetTeam()
	activityMgr := team.GetActivityMgr()
	activity := activityMgr.QueryActivity(activityType)
	queryActivityResultMsg := new(QueryActivityResultMsg)
	queryActivityResultMsg.ActivityInfo = activity.ActivityInfo
	client.SendMsg(queryActivityResultMsg) ///发送给客户端此用户所在竞技场信息
}

func (self *QueryActivityInfoMsg) processAction(client IClient) bool { ///查询活动详细信息
	activityType := GetActivityType(self.ActivityType)
	loger := loger() ///定义日志对象
	if loger.CheckFail("activityType!=nil", activityType != nil, activityType, nil) {
		return false ///被查询的活动必须是有效的
	}
	SendQueryActivityResultMsg(self.ActivityType, client)
	return true
}

type AwardActivityMsg struct { ///客户端请求领奖
	MsgHead      `json:"head"` ///"activity", "accpetaward"
	ActivityType int           `json:"activitytype"` ///活动类型
	AwardIndex   int           `json:"awardindex"`   ///领取第几个奖项
}

func (self *AwardActivityMsg) GetTypeAndAction() (string, string) {
	return "activity", "accpetaward"
}

func (self *AwardActivityMsg) checkAction(client IClient) bool {
	loger := loger() ///定义日志对象
	team := client.GetTeam()
	activityMgr := team.GetActivityMgr()
	activity := activityMgr.GetActivity(self.ActivityType)
	if loger.CheckFail("activity!=nil", activity != nil, activity, nil) {
		return false ///被查询的活动必须是有效的
	}
	if loger.CheckFail("self.AwardIndex>0 && self.AwardIndex<=ActivityAwardMaxCount",
		self.AwardIndex > 0 && self.AwardIndex <= ActivityAwardMaxCount,
		self.AwardIndex, ActivityAwardMaxCount) {
		return false ///被查询的活动必须是有效的
	}
	curAwardCount, maxAwardCount := activity.GetCurrentAwardInfo(self.AwardIndex)
	if loger.CheckFail("curAwardCount<maxAwardCount", curAwardCount < maxAwardCount,
		curAwardCount, maxAwardCount) {
		return false ///已领奖次数必须小于可领奖次数
	}
	return true
}

func (self *AwardActivityMsg) payAction(client IClient) bool {
	team := client.GetTeam()
	activityMgr := team.GetActivityMgr()
	activity := activityMgr.GetActivity(self.ActivityType)
	curAwardCount, maxAwardCount := activity.GetCurrentAwardInfo(self.AwardIndex)
	curAwardCount++
	///更新已领奖次数
	activity.SetCurrentAwardInfo(self.AwardIndex, curAwardCount, maxAwardCount)
	return true
}

func (self *AwardActivityMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	activityType := GetActivityType(self.ActivityType)
	activityAwardType := activityType.GetAwardType(self.AwardIndex)
	awardType := activityType.GetActivityAwardType(activityAwardType)
	addItemIDList := IntList{}
	///得到对应vip等级的奖励
	awardTypeList, awardCountList := awardType.GetAwardItemList(team.VipLevel)
	for i := range awardTypeList {
		awardItemType := awardTypeList[i]
		awardItemCount := awardCountList[i]
		itemIDList := itemMgr.AwardItem(awardItemType, awardItemCount)
		if itemIDList != nil {
			addItemIDList = append(addItemIDList, itemIDList...)
		}
	}
	if addItemIDList.Len() > 0 {
		client.GetSyncMgr().syncAddItem(addItemIDList)
	}
	SendQueryActivityResultMsg(self.ActivityType, client) ///领完奖后需要将新的活动信息发给客户端
	return true
}

func (self *AwardActivityMsg) processAction(client IClient) bool {
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
