package football

import (
	"fmt"
	"reflect"
	//"strconv"
)

///可重置属性管理器
const (
	ResetAttribTypeTeamVolunteerInfo     = 1  ///球星来投系统信息,value1今天剩于签约数 value2刷新累计计数
	ResetAttribTypeVolunteerStarType     = 2  ///球星来投系统刷新球星的类型
	ResetAttribTypeSkipLevelCount        = 3  ///征战八方跳过关卡战斗次数
	ResetAttribTypeTeamTrainMatch        = 4  ///维护球队训练赛信息数据
	ResetAttribTypeTeamShopTimes         = 5  ///玩家购买次数  Value1购买体力次数, Value2购买点金手次数 Value3点金手下次是否爆击
	ResetAttribTypeVipPrivilege          = 6  ///玩家VIP特权,Value1每日礼包已领取次数
	ResetAttribTypeMoneyBuyCount         = 7  ///玩家已购买钻石的套餐次数,v1对应id为1的money套餐信息
	ResetAttribTypeMonthCard             = 8  ///玩家vip月卡相关信息 v1对应月卡截止日期utc时间(0为非月卡) v2领取月卡礼包的次数 v3重置领取次数
	ResetAttribTypeDayTask               = 9  ///玩家日常任务相关信息 v1日常任务每日已领取奖励次数
	ResetAttribTypeTriStar               = 10 ///挑战三星球队的次数的前缀，注意是前缀，每个球队都有3次不得不分开
	ResetAttribTypeChallangeMatchPerfect = 11 ///挑战赛完美艺术
	ResetAttribTypeChallangeMatchCrazy   = 12 ///挑战赛疯狂轰炸
	ResetAttribTypeChallangeMatchDefend  = 13 ///挑战赛无懈可击
	ResetAttribTypeSkillStudyInfo        = 14 ///球星学习技能  value1 = starid  value2 = skilltype  value3 = 状态 1为学习 0为停止
)

type ResetAttribInfo struct { ///数据表记录 对应表dy_resetattrib
	ID        int `json:"id"`       ///id
	OwnID     int `json:"ownid"`    ///拥有者id
	Type      int `json:"type"`     ///属性类型
	ResetTime int `json:"resetime"` ///下次重置时间,utc秒
	Value1    int `json:"value1"`   ///数据1
	Value2    int `json:"value2"`   ///数据2
	Value3    int `json:"value3"`   ///数据3
	Value4    int `json:"value4"`   ///数据4VolunteerInfo
	Value5    int `json:"value5"`   ///数据5
	Value6    int `json:"value6"`   ///数据6
	Value7    int `json:"value7"`   ///数据7
	Value8    int `json:"value8"`   ///数据8
}

type ResetAttrib struct {
	ResetAttribInfo ///数据表记录
	DataUpdater     ///更新数据组件
}

func (self *ResetAttrib) ResetMatchLevelSkipCount(client IClient) { ///比赛跳过次数进行重置
	matchLevelSkipCountMax := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTeam, configTeamLevelMatchCommon, 2) ///取得可跳过次数
	matchLevelSkipCountResetTime := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTeam,
		configTeamLevelMatchCommon, 3) ///取得可跳过次数重置时间
	self.ResetTime = GetHourUTC(matchLevelSkipCountResetTime)
	self.Value1 = matchLevelSkipCountMax
	if client != nil {
		skipLevelMsgResult := new(SkipLevelMsgResult)
		skipLevelMsgResult.RemainSkipCount = self.Value1
		client.SendMsg(skipLevelMsgResult) ///同步最新的剩余次数给客户端
	}
}

func (self *ResetAttrib) ResetShopTimes(client IClient) { ///重置购买次数
	refreshShopTimesHours := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTeam,
		configTeamRestore, 2) ///取得配置表中每日刷新小时数
	self.ResetTime = GetHourUTC(refreshShopTimesHours)
	self.Value1 = 0 ///默认购买次数为1
	self.Value2 = 0
	self.Value3 = 0
	if client != nil {
		client.GetSyncMgr().SyncObject("ResetShopTimes", client.GetTeam())
	}
}

func (self *ResetAttrib) ResetDayTask() { ///重置每日日常任务信息
	now := Now()
	if self.Type != ResetAttribTypeDayTask {
		return ///非vip权限信息
	}
	if now < self.ResetTime {
		return ///未到更新重置时间
	}
	refreshPrivilegeTime := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTask,
		configTaskDayTaskRefresh, 5) ///取得配置表中每日刷新小时数
	self.ResetTime = GetHourUTC(refreshPrivilegeTime)
	self.Value1 = 0 ///日常任务每日已领取奖励次数
}

func (self *ResetAttrib) ResetVipMonthCard() { ///重置VIP特权
	now := Now()
	if self.Type != ResetAttribTypeMonthCard {
		return ///非vip权限信息
	}
	if now < self.ResetTime {
		return ///未到更新重置时间
	}
	refreshPrivilegeTime := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configVip,
		configVipData, 3) ///取得配置表中每日刷新小时数
	self.ResetTime = GetHourUTC(refreshPrivilegeTime)
	self.Value2 = 0 ///重置每日领取礼包次数
}

func (self *ResetAttrib) ResetVipPrivilege() { ///重置VIP特权
	now := Now()
	if self.Type != ResetAttribTypeVipPrivilege {
		return ///非vip权限信息
	}
	if now < self.ResetTime {
		return ///未到更新重置时间
	}
	refreshPrivilegeTime := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configVip,
		configVipData, 3) ///取得配置表中每日刷新小时数
	self.ResetTime = GetHourUTC(refreshPrivilegeTime)
	self.Value1 = 0 ///当前已领vip礼包的vip等级
}

func (self *ResetAttrib) ResetTrainMatch(client IClient) { ///重置购买次数
	refreshTrainMatchHours := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTrainMatch,
		configTrainMatchCommon, 1) ///取得配置表中每日刷新小时数
	self.ResetTime = GetHourUTC(refreshTrainMatchHours)
	self.Value1 = 0               ///总积分清空
	self.Value2 = 0               ///已领奖积分清空
	RefreshTrainMatchTarget(self) ///刷新训练赛信息中的目标信息
	RefreshTrainMatchType(self)   ///刷新训练赛信息中的可选训练项目信息 Value5 为新手指引使用训练次数
	if client != nil {
		SendQueryTrainMatchResultMsg(client)
	}
}

func (self *ResetAttrib) ResetSkillStudyInfo(client IClient) { ///技能学习
	self.ResetTime = -1

	team := client.GetTeam()
	skillMgr := team.GetSkillMgr()
	if self.Value1 != 0 && self.Value2 != 0 {
		skillMgr.AddSkill(self.Value1, self.Value2)
	}

	self.Value1 = 0 //! starID
	self.Value2 = 0 //! skillType
	self.Value3 = 0 //! 1 or 0
}

//func (self *ResetAttrib) ResetTeamVolunteerInfo(client IClient) { ///球星来投系统信息进行重置
//	staticDataMgr := GetServer().GetStaticDataMgr()
//	configStarVolunteer := staticDataMgr.GetConfigStaticData(configStarCenter, configItemStarVolunteer)
//	defaultSignCount, _ := strconv.Atoi(configStarVolunteer.Param1)      ///默认每日签约次数
//	defaultUpdateCount, _ := strconv.Atoi(configStarVolunteer.Param2)    ///默认每日补充刷新次数
//	defaultUpdateMaxCount, _ := strconv.Atoi(configStarVolunteer.Param3) ///刷新次数上限
//	finalUpdateCount := defaultUpdateCount + self.Value2                 ///累加刷新次数
//	if finalUpdateCount > defaultUpdateMaxCount {
//		finalUpdateCount = defaultUpdateMaxCount
//	}
//	newExpireTime, _ := strconv.Atoi(configStarVolunteer.Param4) ///默认过期时间,需要计算
//	newExpireTime = GetHourUTC(newExpireTime)                    ///计算得到过期时间
//	self.ResetTime = newExpireTime
//	self.Value1 = defaultSignCount
//	self.Value2 = finalUpdateCount
//	queryVolunteerResultMsg := new(QueryVolunteerResultMsg)
//	queryVolunteerResultMsg.RemainSignCount = resetAttribTeamVolunteerInfo.Value1
//	queryVolunteerResultMsg.RemainUpdateCount = resetAttribTeamVolunteerInfo.Value2
//	queryVolunteerResultMsg.StarTypeList = IntList{resetAttribTypeVolunteerStarType.Value1, resetAttribTypeVolunteerStarType.Value2,
//		resetAttribTypeVolunteerStarType.Value3, resetAttribTypeVolunteerStarType.Value4}
//	client.SendMsg(queryVolunteerResultMsg)
//}

func (self *ResetAttrib) Update(client IClient, team *Team) { ///时间到了更新自身状态
	switch self.Type {
	case ResetAttribTypeTeamVolunteerInfo: ///球星来投系统信息,value1今天剩于签约数 value2刷新累计计数
		//self.ResetTeamVolunteerInfo(client)
	case ResetAttribTypeVolunteerStarType: ///球星来投系统刷新球星的类型
	case ResetAttribTypeSkipLevelCount: ///征战八方跳过关卡战斗次数
		self.ResetMatchLevelSkipCount(client) ///重置跳过次数并同步客户端
	case ResetAttribTypeTeamShopTimes:
		self.ResetShopTimes(client) ///重置玩家购买次数
	case ResetAttribTypeTeamTrainMatch:
		self.ResetTrainMatch(client) ///重置训练赛信息
	case ResetAttribTypeSkillStudyInfo:
		self.ResetSkillStudyInfo(client) ///重置技能训练
		//	case ResetAttribTypeVipPrivilege:
		//		self.ResetVipPrivilege(client, team)
	}
}

func NewResetAttrib(resetAttribInfo *ResetAttribInfo) *ResetAttrib {
	resetAttrib := new(ResetAttrib)
	resetAttrib.ResetAttribInfo = *resetAttribInfo
	resetAttrib.InitDataUpdater(tableResetAttrib, &resetAttrib.ResetAttribInfo)
	return resetAttrib
}

//type IResetAttribMgr interface { ///被球队拥有处理管理器,用于时间推移给效果的系统
//	IGameMgr
//	AddResetAttrib(attribType int, resettime int, valueList IntList) int   ///添加一个可重置属性
//	GetResetAttrib(attribType int) *ResetAttrib                            ///得到指定类型的属性
//	Save()                                                                 ///保存数据
//	UpdateResetAttrib(attribType int, resettime int, valueList []int) bool ///更新一个可重置属性
//	Update(now int, client IClient)                                        ///可重置值管理器自身更新状态
//}

type ResetAttribList map[int]*ResetAttrib ///可重置属性列表
type ResetAttribMgr struct {              ///可重置属性管理器
	GameMgr
	resetAttribList ResetAttribList ///可重置属性列表
	//ownID           int             ///所属id
}

func NewResetAttribMgr(ownID int) *ResetAttribMgr { ///加载拥有者的可变属性数据
	resetAttribMgr := new(ResetAttribMgr)
	if resetAttribMgr.Init(ownID) == false {
		return nil
	}
	return resetAttribMgr
}

func (self *ResetAttribMgr) GetType() int { ///得到管理器类型
	return mgrTypeResetAttribMgr ///关卡管理器
}

func (self *ResetAttribMgr) SaveInfo() { ///保存数据
	for _, v := range self.resetAttribList {
		v.Save()
	}
}

func (self *ResetAttribMgr) Init(ownID int) bool { ///加载拥有者的可变属性数据
	self.resetAttribList = make(ResetAttribList)
	//	self.teamID = ownID ///存放所有者id

	resetAttribListQuery := fmt.Sprintf("select * from %s where ownid=%d limit 3000", tableResetAttrib, ownID)
	resetAttribInfo := new(ResetAttribInfo)
	resetAttribList := GetServer().GetDynamicDB().fetchAllRows(resetAttribListQuery, resetAttribInfo)
	if nil == resetAttribList {
		return false
	}
	for i := range resetAttribList {
		resetAttribInfo = resetAttribList[i].(*ResetAttribInfo)
		resetAttrib := NewResetAttrib(resetAttribInfo)
		self.resetAttribList[resetAttribInfo.Type] = resetAttrib
	}
	return true
}

func (self *ResetAttribMgr) UpdateResetAttrib(attribType int, resettime int, valueList []int) bool { ///添加一个可重置属性
	///更新内存
	resetAttrib := self.resetAttribList[attribType]
	if nil == resetAttrib {
		return false
	}
	resetAttrib.ResetTime = resettime
	v := reflect.ValueOf(resetAttrib).Elem()
	for i := range valueList {
		intValue := int64(valueList[i])
		fieldName := fmt.Sprintf("Value%d", i+1)
		fieldObj := v.FieldByName(fieldName)
		if fieldObj.IsValid() == false || fieldObj.Kind() != reflect.Int {
			break
		}
		fieldObj.SetInt(intValue)
	}
	return true
}

func (self *ResetAttribMgr) AddResetAttrib(attribType int, resettime int, valueList IntList) *ResetAttrib { ///添加一个可重置属性
	///创建RestAttrib对象
	resetAttrib := new(ResetAttrib)
	resetAttrib.Type = attribType
	resetAttrib.OwnID = self.team.GetID()
	resetAttrib.ResetTime = resettime
	v := reflect.ValueOf(resetAttrib).Elem()
	for i := range valueList {
		intValue := int64(valueList[i])
		fieldName := fmt.Sprintf("Value%d", i+1)
		fieldObj := v.FieldByName(fieldName)
		if fieldObj.IsValid() == false || fieldObj.Kind() != reflect.Int {
			break
		}
		fieldObj.SetInt(intValue)
	}

	insertNewRestAttribQuery := fmt.Sprintf("Insert %s set type=%d,ownid=%d,resettime=%d,value1=%d,value2=%d,value3=%d,value4=%d,value5=%d,value6=%d,value7=%d,value8=%d",
		tableResetAttrib, attribType, self.team.GetID(), resettime, resetAttrib.Value1, resetAttrib.Value2, resetAttrib.Value3, resetAttrib.Value4,
		resetAttrib.Value5, resetAttrib.Value6, resetAttrib.Value7, resetAttrib.Value8)
	///执行插入语句
	lastInsertRestAttribID, _ := GetServer().GetDynamicDB().Exec(insertNewRestAttribQuery)
	if lastInsertRestAttribID <= 0 {
		GetServer().GetLoger().Warn("RestAttribMgr AddRestAttrib insertNewRestAttribQuery fail! centerType")
		return nil
	}
	///更新内存
	resetAttrib.ID = lastInsertRestAttribID
	resetAttrib.InitDataUpdater(tableResetAttrib, &resetAttrib.ResetAttribInfo)
	self.resetAttribList[resetAttrib.Type] = resetAttrib
	return resetAttrib
}

func (self *ResetAttribMgr) QueryResetAttrib(attribType int) *ResetAttrib {
	resetAttrib := self.GetResetAttrib(attribType)
	if nil == resetAttrib {
		resetAttrib = self.AddResetAttrib(attribType, 0, nil)
	}
	return resetAttrib
}

func (self *ResetAttribMgr) GetResetAttrib(attribType int) *ResetAttrib {
	if nil == self.resetAttribList[attribType] {
		return nil
	}
	return self.resetAttribList[attribType]
}

func (self *ResetAttribMgr) Save() {
	for _, v := range self.resetAttribList {
		v.Save()
	}
}

func (self *ResetAttribMgr) Update(now int, client IClient) { ///可重置值管理器自身更新状态
	for _, v := range self.resetAttribList {
		if v.ResetTime > 0 && IsExpireTime(v.ResetTime) == true {
			v.ResetTime = 0             ///关闭时间,防止重复调用
			v.Update(client, self.team) ///要求对象马上更新自身状态
		}
	}
}

func GetTriStarAttribType(npcteamid int) int { ///取得npcteamid对应的可重置记录类型
	return ResetAttribTypeTriStar*10000000 + npcteamid
}

func (self *ResetAttribMgr) ResetTriStar(npcteamid int) { ///对一个三星球队挑战次数进行重置
	attribType := GetTriStarAttribType(npcteamid)  ///取得球队对应的类型
	resetAttrib := self.GetResetAttrib(attribType) ///这里就只找存在的数据而不重新创建
	if resetAttrib != nil {                        ///如果数据存在就更新
		challangeNumber := 3     ///可挑战次数，暂时硬编码
		challangeResetClock := 4 ///重置钟点暂，时硬编码
		resetAttrib.ResetTime = GetHourUTC(challangeResetClock)
		resetAttrib.Value1 = challangeNumber
		resetAttrib.Save() ///保存
	}
}

func (self *ResetAttribMgr) ResetChallangeMatch(challangeMatchType int, client IClient) { ///重置挑战赛数据
	challangeNumber := GetChallangeNumber(client)
	attribType := GetChallangeMatchAttribType(challangeMatchType)
	resetAttrib := self.GetResetAttrib(attribType) ///这里就只找存在的数据而不重新创建
	if resetAttrib != nil {                        ///如果数据存在就更新
		resetAttrib.ResetTime = GetHourUTC(ChallangeMatchResetClock)
		resetAttrib.Value1 = challangeNumber
		resetAttrib.Save() ///保存
	}
}

func (self *ResetAttribMgr) QueryChallangeMatchResetAttrib(challangeMatchType int, client IClient) *ResetAttrib { ///取得可重置数据，没有就新增
	attribType := GetChallangeMatchAttribType(challangeMatchType) ///取得对应的可重置数据类型
	resetAttrib := self.GetResetAttrib(attribType)                ///取得可重置数据对象
	if resetAttrib != nil {                                       ///如果数据存在就返回存在的数据
		if IsExpireTime(resetAttrib.ResetTime) { ///如果已经过期就重置
			self.ResetChallangeMatch(challangeMatchType, client) ///重置数据
		}
	} else { ///数据不存在就创建新的
		challangeNumber := GetChallangeNumber(client) ///取得最大可挑战的次数
		values := []int{challangeNumber}
		resetAttrib = self.AddResetAttrib(attribType, GetHourUTC(ChallangeMatchResetClock), values)
	}
	return resetAttrib
}

func (self *ResetAttribMgr) GetChallangeMatchResetAttrib(challangeMatchType int, client IClient) *ResetAttrib { ///仅仅取得可重置数据，没有就返回nil
	attribType := GetChallangeMatchAttribType(challangeMatchType) ///取得对应的可重置数据类型
	resetAttrib := self.GetResetAttrib(attribType)                ///取得可重置数据对象
	return resetAttrib
}
