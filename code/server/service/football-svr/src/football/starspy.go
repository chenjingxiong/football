package football

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

const (
	primerStarSpy = 1 ///初级球探
	middleStarSpy = 2 ///中级球探
	expertStarSpy = 3 ///高级球探
)

const (
	primerSpyNeedItem = 200001 ///初级球探所需道具
	middleSpyNeedItem = 200002 ///中级球探所需道具
	expertSpyNeedItem = 200003 ///高级球探所需道具
)

const (
	ExpPoolLimit = 0x7FFFFFFF ///经验池上限
	LevelLimit   = 100        ///等级上限
)

const (
	OperateBegin   = 0 //开始
	OperateCanvass = 1 //招至工会
	OperateAddExp  = 2 //转至经验
	OperateLevelup = 3 //提升本队相同球员星级
	OperateEnd     = 4 //结束
)

type StarSpyInfo struct { ///球探信息,和dy_starspy一一对应
	ID                   int ///记录id
	Teamid               int ///所属球队id
	Discoverluck1        int ///初级球探发掘球员幸运值
	Discoverluck2        int ///中级球探发掘球员幸运值
	Discoverluck3        int ///高级球探发掘球员幸运值
	Discovercd1          int ///初级球探发掘球员cd到期utc时间
	Discovercd2          int ///中级球探发掘球员cd到期utc时间
	Discovercd3          int ///高级球探发掘球员cd到期utc时间
	Discoverremaincount1 int ///初级球探发掘球员剩余次数
	Discoverremaincount2 int ///中级球探发掘球员剩余次数
	Discoverremaincount3 int ///高级球探发掘球员剩余次数
	DiscoverResetTime1   int ///初级球探发掘球员次数下次重置时间
	DiscoverResetTime2   int ///中级球探发掘球员次数下次重置时间
	DiscoverResetTime3   int ///高级球探发掘球员次数下次重置时间
	Primerdiscovercount  int ///初级球探发掘次数
}

//type IStarSpy interface {
//	IDataUpdater
//	IsFullDiscoverLuck(starSpyType int) bool                         ///判断球探发掘球员幸运值是否已满
//	Init(teamID int) bool                                            ///加载球队所属球探信息
//	SetDiscoverLuck(starSpyType int, discoverLuck int)               ///设置球探发掘球员幸运值
//	GetDiscoverLuck(starSpyType int) int                             ///设置球探发掘球员幸运值
//	SetDiscoverCD(starSpyType int, discoverCD int)                   ///设置球探发掘球员幸运值
//	GetDiscoverCD(starSpyType int) int                               ///设置球探发掘球员幸运值
//	GetStarSpyLabel(starSpyType int) string                          ///得到球探类型标签
//	SetDiscoverRemainCount(starSpyType int, discoverRemainCount int) ///设置球探发掘球员剩余次值
//	GetDiscoverRemainCount(starSpyType int) int                      ///得到球探发掘球员剩余次值
//	Update(now int, client IClient)                                  ///球探自身更新状态
//	GetInfoCopy() *StarSpyInfo                                       ///得到球探信息副本
//	GetSpyNeedItem(spyType int) int                                  ///得到球探发掘所需道具
//	SetAwardStarType(starType int, evolveCount int)                  ///设置奖励球员类型
//	GetAwardStarType() (int, int)                                    ///得到奖励球员类型
//	GetPrimerDiscoverCount() int                                     ///得到使用初级球探次数
//	SetPrimerDiscoverCount()                                         ///设置使用初级球探次数
//}

type StarSpy struct {
	StarSpyInfo
	DataUpdater
	awardStarType        int ///奖励球员类型
	awardStarEvolveCount int ///奖励球员星级
	///starSpyInfoUpdate  ///球队信息更新组件
}

func (self *StarSpy) GetReflectValue() reflect.Value { ///得到球队反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

func (self *StarSpy) GetInfoCopy() *StarSpyInfo { ///得到球探信息副本
	infoCopy := self.StarSpyInfo
	return &infoCopy
}

//func (self *StarSpy) CalcAddCount(offsetTime int) { ///计算增加免费次数
//}

func (self *StarSpy) CalcOfflineFreeCount(team *Team) { ///计算离线免费次数
	payConfig := GetServer().GetStaticDataMgr().GetConfigStaticData(configStarSpy, configItemDiscoverResetCount)
	///初始化为vip0的默认次数
	maxRemainFreeCount1, maxRemainFreeCount2, maxRemainFreeCount3 := 0, 0, 0
	maxRemainFreeCount1, _ = strconv.Atoi(payConfig.Param2)
	maxRemainFreeCount2, _ = strconv.Atoi(payConfig.Param4)
	maxRemainFreeCount3, _ = strconv.Atoi(payConfig.Param6)

	freeCountResetTime1, freeCountResetTime2, freeCountResetTime3 := 0, 0, 0
	freeCountResetTime1, _ = strconv.Atoi(payConfig.Param1)
	freeCountResetTime2, _ = strconv.Atoi(payConfig.Param3)
	freeCountResetTime3, _ = strconv.Atoi(payConfig.Param5)

	vipLevel := team.GetVipLevel()
	vipPrivilege := GetServer().GetStaticDataMgr().GetVipInfo(vipLevel)
	if vipPrivilege != nil {
		maxRemainFreeCount1 = vipPrivilege.Param6
		maxRemainFreeCount2 = vipPrivilege.Param7
		maxRemainFreeCount3 = vipPrivilege.Param8
	}
	now := Now()                                 ///得到当前时间
	offsetTime1 := now - self.DiscoverResetTime1 ///取得离线后经过时间
	remainTime1 := offsetTime1 % freeCountResetTime1
	needAddCount1 := offsetTime1 / freeCountResetTime1
	if self.DiscoverResetTime1 > 0 && needAddCount1 > 0 && self.Discoverremaincount1 < maxRemainFreeCount1 {
		self.Discoverremaincount1 += needAddCount1 + 1
		self.Discoverremaincount1 = Min(self.Discoverremaincount1, maxRemainFreeCount1) ///防止超上限
		self.DiscoverResetTime1 = now + freeCountResetTime1 - remainTime1
	}

	offsetTime2 := now - self.DiscoverResetTime2
	remainTime2 := offsetTime2 % freeCountResetTime2
	needAddCount2 := offsetTime2 / freeCountResetTime2
	if self.DiscoverResetTime2 > 0 && needAddCount2 > 0 && self.Discoverremaincount2 < maxRemainFreeCount2 {
		self.Discoverremaincount2 += needAddCount2 + 1
		self.Discoverremaincount2 = Min(self.Discoverremaincount2, maxRemainFreeCount2) ///防止超上限
		self.DiscoverResetTime2 = now + freeCountResetTime2 - remainTime2
	}

	offsetTime3 := now - self.DiscoverResetTime3
	remainTime3 := offsetTime3 % freeCountResetTime3
	needAddCount3 := offsetTime3 / freeCountResetTime3
	if self.DiscoverResetTime3 > 0 && needAddCount3 > 0 && self.Discoverremaincount3 < maxRemainFreeCount3 {
		self.Discoverremaincount3 += needAddCount3 + 1
		self.Discoverremaincount3 = Min(self.Discoverremaincount3, maxRemainFreeCount3) ///防止超上限
		self.DiscoverResetTime3 = now + freeCountResetTime3 - remainTime3
	}
}

func (self *StarSpy) Init(teamID int, team *Team) bool { ///加载球队所属球探信息
	starSpyInfoQuery := fmt.Sprintf("select * from %s where teamid=%d limit 1", tableStarSpy, teamID)
	GetServer().GetDynamicDB().fetchOneRow(starSpyInfoQuery, &self.StarSpyInfo)
	self.CalcOfflineFreeCount(team)
	self.InitDataUpdater(tableStarSpy, &self.StarSpyInfo) ///球队信息更新组件
	return self.ID > 0
}

func (self *StarSpy) SetDiscoverRemainCount(starSpyType int, discoverRemainCount int) { ///设置球探发掘球员剩余次值
	if discoverRemainCount < 0 {
		return ///拒绝负值
	}
	switch starSpyType {
	case primerStarSpy:
		self.Discoverremaincount1 = discoverRemainCount
	case middleStarSpy:
		self.Discoverremaincount2 = discoverRemainCount
	case expertStarSpy:
		self.Discoverremaincount3 = discoverRemainCount
	}
}

func (self *StarSpy) GetDiscoverRemainCount(starSpyType int) int { ///设置球探发掘球员剩余值
	result := 0
	switch starSpyType {
	case primerStarSpy:
		result = self.Discoverremaincount1
	case middleStarSpy:
		result = self.Discoverremaincount2
	case expertStarSpy:
		result = self.Discoverremaincount3
	}
	return result
}

func (self *StarSpy) SetDiscoverCD(starSpyType int, discoverCD int) { ///设置球探发掘球员CD时间
	if discoverCD < 0 {
		return ///拒绝负值
	}
	switch starSpyType {
	case primerStarSpy:
		self.Discovercd1 = discoverCD
	case middleStarSpy:
		self.Discovercd2 = discoverCD
	case expertStarSpy:
		self.Discovercd3 = discoverCD
	}
}

func (self *StarSpy) SetResetRemainCD(starSpyType int, resetTime int) { ///设置球探发掘球员幸运值
	switch starSpyType {
	case primerStarSpy:
		self.DiscoverResetTime1 = resetTime
	case middleStarSpy:
		self.DiscoverResetTime2 = resetTime
	case expertStarSpy:
		self.DiscoverResetTime3 = resetTime
	}
}

func (self *StarSpy) GetResetRemainCD(starSpyType int) int { ///设置球探发掘球员幸运值
	result := 0
	switch starSpyType {
	case primerStarSpy:
		result = self.DiscoverResetTime1
	case middleStarSpy:
		result = self.DiscoverResetTime2
	case expertStarSpy:
		result = self.DiscoverResetTime3
	}
	return result
}

func (self *StarSpy) GetDiscoverCD(starSpyType int) int { ///设置球探发掘球员幸运值
	result := 0
	now := int(time.Now().Unix())
	switch starSpyType {
	case primerStarSpy:
		result = self.Discovercd1
	case middleStarSpy:
		result = self.Discovercd2
	case expertStarSpy:
		result = self.Discovercd3
	}
	if now >= result {
		result = 0 ///过期自动重置时间
	}
	return result
}

func (self *StarSpy) SetAwardStarType(starType int, evolveCount int) { ///设置奖励球员类型
	self.awardStarType = starType
	self.awardStarEvolveCount = evolveCount
}

func (self *StarSpy) GetAwardStarType() (int, int) { ///得到奖励球员类型
	return self.awardStarType, self.awardStarEvolveCount
}

func (self *StarSpy) SetDiscoverLuck(starSpyType int, discoverLuck int) { ///设置球探发掘球员幸运值
	if discoverLuck < 0 {
		return ///拒绝负值
	}
	if discoverLuck > 100 {
		discoverLuck = 100
	}
	switch starSpyType {
	case primerStarSpy:
		self.Discoverluck1 = discoverLuck
	case middleStarSpy:
		self.Discoverluck2 = discoverLuck
	case expertStarSpy:
		self.Discoverluck3 = discoverLuck
	}
}

func (self *StarSpy) GetStarSpyLabel(starSpyType int) string { ///得到球探类型标签
	StarSpyLabelList := map[int]string{1: "PrimerStarSpy", 2: "MiddleStarSpy", 3: "ExpertStarSpy"}
	starSpyLabel, ok := StarSpyLabelList[starSpyType]
	if false == ok {
		return "Unkown"
	}
	return starSpyLabel
}

func (self *StarSpy) GetDiscoverLuck(starSpyType int) int { ///设置球探发掘球员幸运值
	result := 0
	switch starSpyType {
	case primerStarSpy:
		result = self.Discoverluck1
	case middleStarSpy:
		result = self.Discoverluck2
	case expertStarSpy:
		result = self.Discoverluck3
	}
	return result
}

func (self *StarSpy) resetDiscoverRemainCount(now int, client IClient) int { ///重置发掘剩余次数
	result := 0
	payConfig := GetServer().GetStaticDataMgr().GetConfigStaticData(configStarSpy, configItemDiscoverResetCount)
	if nil == payConfig {
		return 0
	}
	///初始化为vip0的默认次数
	maxRemainFreeCount1, maxRemainFreeCount2, maxRemainFreeCount3 := 0, 0, 0
	maxRemainFreeCount1, _ = strconv.Atoi(payConfig.Param2)
	maxRemainFreeCount2, _ = strconv.Atoi(payConfig.Param4)
	maxRemainFreeCount3, _ = strconv.Atoi(payConfig.Param6)

	team := client.GetTeam()
	vipLevel := team.GetVipLevel()
	vipPrivilege := GetServer().GetStaticDataMgr().GetVipInfo(vipLevel)
	if vipPrivilege != nil {
		maxRemainFreeCount1 = vipPrivilege.Param6
		maxRemainFreeCount2 = vipPrivilege.Param7
		maxRemainFreeCount3 = vipPrivilege.Param8
	}

	if now >= self.DiscoverResetTime1 {
		if self.Discoverremaincount1 < maxRemainFreeCount1 {
			//discoverCount, _ := strconv.Atoi(payConfig.Param2) ///重置次数
			self.Discoverremaincount1 += 1 ///减已使用一次,则加一次可用次数
			result += 1
			self.DiscoverResetTime1, _ = strconv.Atoi(payConfig.Param1)
			self.DiscoverResetTime1 = self.DiscoverResetTime1 + now ///写入下次重置时间秒
		}
		//		} else {
		//			self.DiscoverResetTime1 = 0
		//		}
	}
	if now >= self.DiscoverResetTime2 {
		if self.Discoverremaincount2 < maxRemainFreeCount2 {
			//discoverCount, _ := strconv.Atoi(payConfig.Param4) ///重置次数
			self.Discoverremaincount2 += 1
			result += 1
			self.DiscoverResetTime2, _ = strconv.Atoi(payConfig.Param3)
			self.DiscoverResetTime2 = self.DiscoverResetTime2 + now
		}
		//		} else {
		//			self.DiscoverResetTime2 = 0
		//		}
	}
	if now >= self.DiscoverResetTime3 {
		if self.Discoverremaincount3 < maxRemainFreeCount3 {
			//discoverCount, _ := strconv.Atoi(payConfig.Param6)
			self.Discoverremaincount3 += 1 ///重置次数
			result += 1
			self.DiscoverResetTime3, _ = strconv.Atoi(payConfig.Param5)
			self.DiscoverResetTime3 = self.DiscoverResetTime3 + now
		}
		//		} else {
		//			self.DiscoverResetTime3 = 0
		//		}
	}
	return result
}

func (self *StarSpy) resetDiscoverCD(now int) int { ///重置发掘CD
	result := 0
	if self.Discovercd1 > 0 && now >= self.Discovercd1 {
		self.SetDiscoverCD(primerStarSpy, 0) ///重置cd1
		result += 1
	}
	if self.Discovercd2 > 0 && now >= self.Discovercd2 {
		self.SetDiscoverCD(middleStarSpy, 0) ///重置cd2
		result += 1
	}
	if self.Discovercd3 > 0 && now >= self.Discovercd3 {
		self.SetDiscoverCD(expertStarSpy, 0) ///重置cd3
		result += 1
	}
	return result
}

func (self *StarSpy) Update(now int, client IClient) { ///球队自身更新状态
	result := 0
	//	result += self.resetDiscoverCD(now)                  ///重置发掘CD
	result += self.resetDiscoverRemainCount(now, client) ///重置发掘剩余次数
	if result > 0 {
		client.GetSyncMgr().SyncObject("StarSpyUpdate", self)
	}
}

func (self *StarSpy) GetID() int { ///得到球队反射对象
	return self.ID
}

func (self *StarSpy) IsFullDiscoverLuck(starSpyType int) bool { ///判断球探发掘球员幸运值是否已满
	isFullDiscoverLuck := false
	switch starSpyType {
	case primerStarSpy:
		isFullDiscoverLuck = self.Discoverluck1 >= 100
	case middleStarSpy:
		isFullDiscoverLuck = self.Discoverluck2 >= 100
	case expertStarSpy:
		isFullDiscoverLuck = self.Discoverluck3 >= 100
	}
	return isFullDiscoverLuck
}

func (self *StarSpy) GetSpyNeedItem(spyType int) int { ///得到球探所需道具
	if spyType < primerStarSpy || spyType > expertStarSpy {
		return 0
	}

	spyNeedItemList := IntList{primerSpyNeedItem, middleSpyNeedItem, expertSpyNeedItem}
	spyTypeIndex := spyType - 1
	return spyNeedItemList[spyTypeIndex]

}

func (self *StarSpy) GetPrimerDiscoverCount(awardMask int) bool {
	return TestMask(self.Primerdiscovercount, awardMask)
}

func (self *StarSpy) SetPrimerDiscoverCount(awardMask int) {
	self.Primerdiscovercount = SetMask(self.Primerdiscovercount, awardMask, 1)
	//self.Primerdiscovercount = self.Primerdiscovercount | int((1 << (awardMask - 1)))
}

//! 退出时需要判断是否抽取球员
func (self *StarSpy) OnLogout(team *Team) {
	teamInfo := team.GetInfo()

	awardStarType, awardStarEvolveCount := self.GetAwardStarType()
	if awardStarType == 0 {
		return
	}

	starType := GetServer().GetStaticDataMgr().Unsafe().GetStarType(awardStarType)
	if starType == nil {
		return
	}

	isTeamHas := team.HasStar(awardStarType)

	if isTeamHas == false { //! 若不在则招募
		///奖励球员,并赋予球员初始星级
		starID := team.AwardStar(awardStarType)
		star := team.GetStar(starID)
		star.SetStarCount(awardStarEvolveCount)
		//star.GetInfo().EvolveCount = awardStarEvolveCount
		//sync.syncAddStar(IntList{starID})

		GetServer().GetLoger().CYDebug("get a player")
	} else { //! 若存在则转化
		if teamInfo.StarExpPool >= ExpPoolLimit {
			return
		}

		///（初始值之和+3*成长值之和）*星级^2
		starType := GetServer().GetStaticDataMgr().Unsafe().GetStarType(awardStarType)
		self.SetAwardStarType(0, 0) ///处理完毕,清空类型
		starCardCount := team.GetStarCardCount(starType, awardStarEvolveCount)
		team.AwardObject(ItemStarCard, starCardCount, 0, 0) ///赠星卡

		GetServer().GetLoger().CYDebug("trans a player")
	}

	//sync := client.GetSyncMgr()
	//awardStarType, awardStarEvolveCount := starSpy.GetAwardStarType()
	//switch self.OperateType {
	//case OperateCanvass:
	//	// awardMemberID := starCenter.AwardSpecialMember(starCenterTypeDiscover, awardStarType, awardStarEvolveCount, 0) ///发奖

	//	// ///同步球员中心新加成员信息给客户端
	//	// starCenterAddMemberMsg := NewStarCenterAddMemberMsg()
	//	// starCenterMember := starCenter.GetStarCenterMember(starCenterTypeDiscover, awardMemberID)
	//	// starCenterAddMemberMsg.AddMember(starCenterMember)
	//	// client.SendMsg(starCenterAddMemberMsg)
	//	// starSpy.SetAwardStarType(0, 0) ///处理完毕,清空类型

	//	///奖励球员,并赋予球员初始星级
	//	starID := team.AwardStar(awardStarType)
	//	star := team.GetStar(starID)
	//	star.SetStarCount(awardStarEvolveCount)
	//	//star.GetInfo().EvolveCount = awardStarEvolveCount
	//	sync.syncAddStar(IntList{starID})

	//case OperateAddExp:
	//	///（初始值之和+3*成长值之和）*星级^2
	//	starType := GetServer().GetStaticDataMgr().Unsafe().GetStarType(awardStarType)
	//	//		firstValue, growValue := self.GetFirstAndGrowValue(starType)
	//	//		teamInfo := team.GetInfo()
	//	//		exp := (firstValue + growValue*3) * (awardStarEvolveCount * awardStarEvolveCount)

	//	//		if teamInfo.StarExpPool+exp > ExpPoolLimit {
	//	//			teamInfo.StarExpPool = ExpPoolLimit
	//	//		} else {
	//	//			teamInfo.StarExpPool += exp
	//	//		}
	//	//		fmt.Printf("FirstValue = %d  \r\n   GrowValue = %d  \r\n  StarLevel = %d \r\n  EXP = %d \r\n",
	//	//			firstValue, growValue, awardStarEvolveCount, exp)
	//	///同步到客户端
	//	//		sync.SyncObject("StarSpyOperateMsg", team)
	//	starSpy.SetAwardStarType(0, 0) ///处理完毕,清空类型
	//	///给星卡
	//	starCardCount := team.GetStarCardCount(starType, awardStarEvolveCount)
	//	team.AwardObject(ItemStarCard, starCardCount, 0, 0) ///赠星卡
	//}
}
