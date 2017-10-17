package football

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	drawGroupTypeStarSpy1          = 1  ///初级球探
	drawGroupTypeStarTenDrawSuper  = 2  ///十连抽高级组
	drawGroupTypeStarTenDrawNormal = 3  ///十连抽普通组
	drawGroupTypeTenDrawItem       = 4  ///十连抽道具组
	drawGroupTypeStarSpy2          = 5  ///组5为中级寻找球员
	drawGroupTypeStarSpy3          = 6  ///组6为高级寻找球员
	drawGroupTypeVolunteer         = 7  ///球星来投组号
	drawGroupTypeDefaultMasterStar = 8  ///初始当家球星分组
	drawGroupTypeDefaultStarGK     = 9  ///初始门将分组
	drawGroupTypeDefaultStarLB     = 10 ///初始左后卫分组
	drawGroupTypeDefaultStarCB     = 11 ///初始中后卫分组
	drawGroupTypeDefaultStarRB     = 12 ///初始右后卫分组
	drawGroupTypeDefaultStarLMF    = 13 ///初始左前卫分组
	drawGroupTypeDefaultStarCMF    = 14 ///初始中前卫分组
	drawGroupTypeDefaultStarRMF    = 15 ///初始右前卫分组
	drawGroupTypeDefaultStarCF     = 16 ///初始中锋分组
	drawGroupTypeStarSpyFull1      = 17 ///组17为初级满幸运
	drawGroupTypeStarSpyFull2      = 18 ///组18为中级满幸运
	drawGroupTypeStarSpyFull3      = 19 ///组19为高级满幸运
	drawGroupTypeMannaStar         = 20 ///组20为天赐球员碎片组
	drawGroupTypeMergeStar1        = 21 //!组21为融合球员抽取组1
	drawGroupTypeMergeStar2        = 22 //!组22为融合球员抽取组2
)

const (
	configStarSpy    = "starspy"    ///球探配置信息
	configStarCenter = "starcenter" ///球员中心配置信息
	configTeam       = "team"       ///球员中心配置信息
	configStar       = "star"       ///球员中心配置信息
	configItem       = "item"       ///道具装备配置信息
	configTrainMatch = "trainmatch" ///训练赛配置信息
	configArenaMatch = "arenamatch" ///联赛配置信息
	configServer     = "server"     ///服务器配置信息
	configVip        = "vip"        ///VIP配置信息
	configGoldFinger = "goldfinger" ///金手指配置信息
	configTask       = "task"       ///任务配置信息
)
const (
	configItemDiscoverExtraPay       = "discoverextrapay"      ///发掘球员代价
	configItemDiscoverPrimaryRate    = "primarydiscoverrate"   ///初级球探发掘球员几率
	configItemDiscoverMiddleRate     = "middlediscoverrate"    ///中级球探发掘球员几率
	configItemDiscoverExpertRate     = "expertdiscoverate"     ///高级球探发掘球员几率
	configItemDiscoverAddLuck        = "discoveraddluck"       ///球探每次发掘后幸运提升值,顺着标记初级,中级,高级球探
	configItemDiscoverCD             = "discovercd"            ///球探每次发掘后CD时间
	configItemDiscoverResetCount     = "resetdiscovercount"    ///球探发掘重置次数
	configItemDiscoverConfig         = "discoverconfig"        ///球探发掘通用配置 param1球探发掘的球员上限数
	configItemStarCenterCommonConfig = "commonconfig"          ///球员中心通用配置,param1为球员中心成员最大停留时间
	configItemStarTrain              = "startrain"             ///球员训练配置参数 param1训练位上训 param2开一格价格 param3开训练格封顶价格
	configItemDefaultTeamParam       = "defaultteamparam"      ///创建球队时默认配置参数
	configItemStarVolunteer          = "startvolunteer"        ///球星来投配置 1每天签约次数 2每天新加刷新次数 3刷新次数上限 4每天重置球星来投信息时间(4点钟)
	configItemStarCommonConfig       = "commonconfig"          ///球员通用配置 1球员等级上限
	configStarGradeEffect            = "gradeeffect"           ///品质会影响球员升级时提升主属性,千分单位
	configItemColorAdditionParam     = "coloradditionparam"    ///道具装备颜色品质加成系数 param1对应color为1的装备,以此类推
	configItemSumColorAdditionParam  = "colorsumadditionparam" ///道具装备颜色累计加成系统 param1为color为1的系数,以此类推
	configTeamLevelMatchCommon       = "levelmatchcommon"      ///球队在征战八方系统中的通用参数配置 param1 征战八方输产生CD秒
	configStarEvolveLevelLimit       = "evolvelevellimit"      ///球员突破次数限制,p1为绿色球员最多只能突破2次,后面以此类推
	configTrainMatchCommon           = "commonconfig"          ///训练赛一般配置信息
	configArenaMatchCommon           = "commonconfig"          ///联赛一般配置信息
	configServerCommon               = "commonconfig"          ///服务器一般配置信息
	configVipData                    = "vipconfig"             ///VIP一般配置信息
	configTeamRestore                = "restoretime"           ///球队恢复时间周期配置信息
	configGoldFingerCommonConfig     = "goldfingerconfig"      ///金手指基本配置信息
	configArenaPromotAward           = "promotaward"           ///联赛晋级奖励配置信息
	configStarSpyDiscover            = "primarydiscover"       ///初级三次固定抽取
	configTaskDayTaskRefresh         = "daytaskrefresh"        ///p1~p4每日任务刷新时间(小时),
	configTaskDayTaskInit            = "initdaytask"           ///日常任务中初始的5个任务类型id,对应tasktype表
)

type DrawGroupStaticData struct { ///抽奖分组表静态数据表
	ID          int ///记录id
	Group       int ///抽选分组
	AwardType   int ///奖励物类型
	AwardCount  int ///奖励数量
	ReqLevel    int ///需求等级
	ReqMoney    int ///人民币需求
	ReqVipLevel int ///vip等级需求
	ReqPass     int ///通关需求
	TakeWeight  int ///抽取权重
	ShowWeight  int ///显示权重
}

type ConfigStaticData struct { ///服务器配置静态数据表
	ID         int    ///记录id
	MasterType string ///主配置类型
	SubType    string ///次配置类型
	Param1     string ///参数1
	Param2     string ///参数2
	Param3     string ///参数3
	Param4     string ///参数1
	Param5     string ///参数2
	Param6     string ///参数3
	Desc       string ///描述 不读内存
}

//type IStaticDataMgrUnsafe interface { ///不安全接口
//	GetLevelType(levelType int) *LevelTypeStaticData       ///得到关卡类型信息指针
//	GetTaskType(taskType int) *TaskTypeStaticData          ///得到任务类型信息指针
//	GetNpcTeamType(npcTeamType int) *NpcTeamTypeStaticData ///得到npc球队类型信息
//	GetTaskTypeList() TaskTypePtrList                      ///得到所有任务类型对象指针列表
//	//GetStarTypeInfo(starType int) *StarTypeStaticData                                                  ///得到球员类型信息
//	GetSeatType(seatType int) *SeatTypeStaticData                                                      ///得到位置类型信息
//	GetFormationType(formationType int) *FormationTypeStaticData                                       ///得到阵形类型信息
//	GetFieldValueList(value interface{}, fieldNamePrefix string, beginIndex int, endIndex int) IntList ///得到字段值列表
//	GetStarType(starType int) *StarTypeStaticData                                                      ///得到球员类型信息
//	GetItemType(itemType int) *ItemTypeStaticData                                                      ///得到道具类型信息
//	GetTacticType(tacticType int) *TacticTypeStaticData                                                ///得到阵形战术类型信息
//	GetSkillType(skillType int) *SkillTypeStaticData                                                   ///得到技能类型信息
//	GetLevelExpType(expType int, expLevel int) *LevelExpStaticData                                     ///得到升级经验信息
//	GetStarFateType(starFateType int) *StarFateTypeStaticData                                          ///得到球员缘系统类型信息
//	GetStarLobbyType(starLobbyType int) *StarLobbyTypeStaticData                                       ///得到游说球员系统类型信息
//	GetVipShopItemInfo(commodityID int) *VipShopStaticData                                             ///得到商城物品类型信息
//	GetVipInfo(vipLevel int) *VipPrivilegeStaticData                                                   ///得到VIP类型信息
//}

//type IStaticDataMgr interface {
//	//getConfigStaticData(masterType string, subType string) *ConfigStaticData
//	GetConfigStaticData(masterType string, subType string) *ConfigStaticData
//	//GetConfigStaticDataReflect(masterType string, subType string) *reflect.Value
//	getConfigStaticDataParamIntList(masterType string, subType string) IntList     ///返回配置表的param整型数值列表
//	GetConfigStaticDataInt(masterType string, subType string, paramIndex int) int  ///得到int配置数据
//	getConfigStaticDataParamIntMap(masterType string, subType string) *map[int]int ///返回配置表的param整型数值列表
//	GetDrawGroupIndexList(drawGroup int) IntList                                   ///得到DrawGroupList
//	GetDrawGroupStaticData(drawGroupID int) *DrawGroupStaticData                   ///得到DrawGroupList数据
//	GetStaticData(tableName string, dataID int) interface{}                        ///得到静态数据
//	//GetStarTypeInfoCopy(starType int) *StarTypeStaticData                          ///得到球员类型信息复本
//	GetStarTypeBasePrice(starType int) int                    ///得到球员基础身价
//	GetItemTypeMerge(itemType int) int                        ///得到道具基础融合价值
//	GetLevelExpNeedExp(expType int, expLevel int) int         ///得到经验配置对象所需经验值
//	GetLevelExpNeedEvolveCount(expType int, expLevel int) int ///得到经验配置对象所需经验值
//	GetTacticTypeList(formType int, formLevel int) IntList    ///得到指定阵形和指定等级的战术列表
//	GetLevelExpCount(expType int) int                         ///得到经验配置对象总项目数,用作等级上限
//	GetTeamMaxStarCount() int                                 ///判断球队中的雇佣球员数上限数
//	Unsafe() IStaticDataMgrUnsafe                             ///得到静态管理器不安全接口
//	GetStaticDataList(tableName string) StaticData            ///得到静态数据表
//	GetStarTypeClass(starType int) int                        ///得到球员评价
//	GetLevelExpNeedCoin(expType int, expLevel int) int        ///得到经验配置对象所需金币
//	GetVipShopStaticDataMap() VipShopMap                      ///得到商城物品静态信息
//}

type StaticData map[int]interface{}        ///一张表,key为id,value为具体结构类型
type StaticDataList map[string]StaticData  ///静态数据表列表,key为表名,value为静态数据表
type DrawGroupIndexList map[int][]int      ///抽奖库索引列表
type IntMap map[int]int                    ///整数map列表
type LevelExpIndexList map[int]IntMap      ///升级经验索引列表,key 经验类型 v{key是等级 v是记录id}
type ScoreCalParamDic [][]int              ///球员评分参数字典
type VipShopMap map[int]*VipShopStaticData ///商城商品信息字典
type StaticDataMgr struct {                ///静态数据管理器
	staticDB           *DBServer          ///静态数据库组件
	staticDataList     StaticDataList     ///静态数据缓存组件
	drawGroupIndexList DrawGroupIndexList ///抽奖库索引列表
	levelExpIndexList  LevelExpIndexList  ///升级表索引表
}

func (self *StaticDataMgr) GetDrawGroupStaticData(drawGroupID int) *DrawGroupStaticData { ///得到DrawGroupList数据
	value, ok := self.staticDataList[tableDrawGroup][drawGroupID]
	if false == ok {
		return nil
	}
	drawGroupStaticData := value.(*DrawGroupStaticData)
	return drawGroupStaticData
}

func (self *StaticDataMgr) GetDrawGroupIndexList(drawGroup int) IntList { ///得到DrawGroupList
	if drawGroup <= 0 {
		return nil
	}
	drawGroupIndexList := self.drawGroupIndexList[drawGroup]
	if nil == drawGroupIndexList {
		return nil
	}
	drawGroupIndexListCopy := make(IntList, len(drawGroupIndexList))
	copy(drawGroupIndexListCopy, drawGroupIndexList)
	return drawGroupIndexListCopy
}

func (self *StaticDataMgr) GetTaskTypeList() TaskTypePtrList { ///得到所有任务类型对象指针列表
	if nil == self.staticDataList[tableTaskType] {
		return TaskTypePtrList{}
	}
	taskTypePtrList := TaskTypePtrList{}
	for _, v := range self.staticDataList[tableTaskType] {
		taskTypePtrList = append(taskTypePtrList, v.(*TaskTypeStaticData))
	}
	return taskTypePtrList
}

func (self *StaticDataMgr) GetTacticTypeList(formType int, formLevel int) IntList { ///得到指定阵形和指定等级的战术列表
	resultList := IntList{}
	if nil == self.staticDataList[tableTacticType] {
		return IntList{}
	}
	for _, v := range self.staticDataList[tableTacticType] {
		tacticType := v.(*TacticTypeStaticData)
		if tacticType.OpenFormType != formType {
			continue ///阵形不匹配
		}
		if tacticType.OpenFormLevel > formLevel {
			continue ///等级不匹配
		}
		resultList = append(resultList, tacticType.ID) ///放入结果
	}
	return resultList
}

///得到结构体内连续的整数数值列表
func (self *StaticDataMgr) GetFieldValueList(value interface{}, fieldNamePrefix string, beginIndex int, endIndex int) IntList {
	resultList := IntList{}
	element := reflect.ValueOf(value).Elem()
	for i := beginIndex; i <= endIndex; i++ {
		fieldName := fmt.Sprintf("%s%d", fieldNamePrefix, i)
		fieldValue := element.FieldByName(fieldName)
		if fieldValue.IsValid() == false || fieldValue.Kind() != reflect.Int {
			continue
		}
		resultList = append(resultList, int(fieldValue.Int()))
	}
	return resultList
}

func (self *StaticDataMgr) GetLevelExpCount(expType int) int { ///得到经验配置对象总项目数,用作等级上限
	if nil == self.levelExpIndexList[expType] {
		return 0
	}
	LevelExpCount := len(self.levelExpIndexList[expType])
	return LevelExpCount
}

func (self *StaticDataMgr) GetLevelExpType(expType int, expLevel int) *LevelExpStaticData { ///得到升级经验信息
	if nil == self.levelExpIndexList[expType] {
		return nil
	}
	if self.levelExpIndexList[expType][expLevel] <= 0 {
		return nil
	}
	levelExpStaticDataID := self.levelExpIndexList[expType][expLevel]
	levelExpStaticData := self.GetStaticData(tableLevelExp, levelExpStaticDataID)
	if nil == levelExpStaticData {
		return nil
	}
	return levelExpStaticData.(*LevelExpStaticData)
}

func (self *StaticDataMgr) GetLevelExpNeedExp(expType int, expLevel int) int { ///得到经验配置对象所需经验值
	if nil == self.levelExpIndexList[expType] {
		return 0
	}
	if self.levelExpIndexList[expType][expLevel] <= 0 {
		return 0
	}
	levelExpStaticDataID := self.levelExpIndexList[expType][expLevel]
	levelExpStaticData := self.GetStaticData(tableLevelExp, levelExpStaticDataID)
	if nil == levelExpStaticData {
		return 0
	}
	return levelExpStaticData.(*LevelExpStaticData).NeedExp
}

func (self *StaticDataMgr) GetLevelExpNeedEvolveCount(expType int, expLevel int) int { ///得到经验配置对象所需星级
	if nil == self.levelExpIndexList[expType] {
		return 0
	}
	if self.levelExpIndexList[expType][expLevel] <= 0 {
		return 0
	}
	levelExpStaticDataID := self.levelExpIndexList[expType][expLevel]
	levelExpStaticData := self.GetStaticData(tableLevelExp, levelExpStaticDataID)
	if nil == levelExpStaticData {
		return 0
	}
	return levelExpStaticData.(*LevelExpStaticData).NeedEvolveCount
}

func (self *StaticDataMgr) GetLevelExpNeedCoin(expType int, expLevel int) int { ///得到经验配置对象所需金币
	if nil == self.levelExpIndexList[expType] {
		return 0
	}
	if self.levelExpIndexList[expType][expLevel] <= 0 {
		return 0
	}
	levelExpStaticDataID := self.levelExpIndexList[expType][expLevel]
	levelExpStaticData := self.GetStaticData(tableLevelExp, levelExpStaticDataID)
	if nil == levelExpStaticData {
		return 0
	}
	return levelExpStaticData.(*LevelExpStaticData).PayCoin
}

func (self *StaticDataMgr) buildLevelExpIndexList() { ///建立经验升级索引表
	self.levelExpIndexList = make(LevelExpIndexList)
	for _, v := range self.staticDataList[tableLevelExp] {
		levelExpStaticData := v.(*LevelExpStaticData)
		if nil == self.levelExpIndexList[levelExpStaticData.Type] {
			self.levelExpIndexList[levelExpStaticData.Type] = make(IntMap)
		}
		self.levelExpIndexList[levelExpStaticData.Type][levelExpStaticData.Level] = levelExpStaticData.ID
	}
}

func (self *StaticDataMgr) buildDrawGroupIndexList() { ///建立抽奖索引表
	self.drawGroupIndexList = make(DrawGroupIndexList)
	for _, v := range self.staticDataList[tableDrawGroup] {
		drawGroupStaticData := v.(*DrawGroupStaticData)
		self.drawGroupIndexList[drawGroupStaticData.Group] = append(self.drawGroupIndexList[drawGroupStaticData.Group], drawGroupStaticData.ID)
	}

}

func (self *StaticDataMgr) Init() { ///初始化静态数据管理器
	self.staticDB = new(DBServer)
	self.staticDB.Init(&GetServer().GetConfig().StaticDNS) ///初始化静态数据组件

	///初始化静态数据容器组件
	self.staticDataList = make(StaticDataList)
	self.LoadStaticData(tableConfig, new(ConfigStaticData))
	self.LoadStaticData(tableStarType, new(StarTypeStaticData))
	self.LoadStaticData(tableDrawGroup, new(DrawGroupStaticData))
	self.LoadStaticData(tableItemType, new(ItemTypeStaticData))             ///加载道具类型表
	self.LoadStaticData(tableLevelExp, new(LevelExpStaticData))             ///加载经验升级类型表
	self.LoadStaticData(tableFormationType, new(FormationTypeStaticData))   ///加载经验升级类型表
	self.LoadStaticData(tableTacticType, new(TacticTypeStaticData))         ///加载经验升级类型表
	self.LoadStaticData(tableLevelType, new(LevelTypeStaticData))           ///加载关卡类型表
	self.LoadStaticData(tableNpcTeamType, new(NpcTeamTypeStaticData))       ///加载npc球队类型表
	self.LoadStaticData(tableSeatType, new(SeatTypeStaticData))             ///加载位置类型表
	self.LoadStaticData(tableTaskType, new(TaskTypeStaticData))             ///加载任务类型表
	self.LoadStaticData(tableSkillType, new(SkillTypeStaticData))           ///加载技能类型表
	self.LoadStaticData(tableStarFateType, new(StarFateTypeStaticData))     ///加载球员缘类型表
	self.LoadStaticData(tableStarLobbyType, new(StarLobbyTypeStaticData))   ///加载球员游说类型表
	self.LoadStaticData(tableTrainAward, new(TrainAwardStaticData))         ///加载训练赛奖励表
	self.LoadStaticData(tableArenaType, new(ArenaType))                     ///加载联赛奖励表
	self.LoadStaticData(tableVipShopType, new(VipShopStaticData))           ///加载商城类型表
	self.LoadStaticData(tableChallangeMatch, new(ChallangeMatchStaticData)) ///加载挑战赛表

	self.LoadStaticData(tableActivityType, new(ActivityType))       ///加载活动类型表
	self.LoadStaticData(tableActivityAward, new(ActivityAwardType)) ///加载活动奖励类型表
	self.LoadStaticData(tableActionType, new(ActionType))           ///加载活动奖励类型表

	self.LoadStaticData(tableVipPrivilege, new(VipPrivilegeStaticData))       ///加载VIP特权类型表
	self.LoadStaticData(tableLeagueAwardType, new(LeagueAwardTypeStaticData)) ///加载推图奖励表
	self.LoadStaticData(tableAwardType, new(ActivitCodeAward))                ///加载激活码奖励表
	self.LoadStaticData(tableMoney, new(MoneyType))                           ///加载充值类型表

	self.LoadStaticData(tableNpcStarType, new(NpcStarTypeStaticData)) ///加载npc球星表

	//	self.LoadStaticData(tableSDKType, new(SDKType)) ///加载SDK平台类型表
	///建立必要的索引
	self.buildDrawGroupIndexList() ///建立抽奖索引表
	self.buildLevelExpIndexList()  ///建立经验升级索引表
	//self.PrintStaticData(tableDrawGroup)
	//result := self.GetDrawGroupIndexList(1)
	//result := self.GetDrawGroupStaticData(1)
	//fmt.Println(*result)
}

func (self *StaticDataMgr) PrintStaticData(tableName string) { ///打印一张静态数据列表,用于调试
	fmt.Println("Static Table:", tableName)
	for _, value := range self.staticDataList[tableName] {
		switch value.(type) {
		case *ConfigStaticData:
			fmt.Println(*value.(*ConfigStaticData))
		case *StarTypeStaticData:
			fmt.Println(*value.(*StarTypeStaticData))
		case *DrawGroupStaticData:
			fmt.Println(*value.(*DrawGroupStaticData))
		}
	}
}

func (self *StaticDataMgr) Unsafe() *StaticDataMgr { ///得到静态管理器不安全接口
	return self
}

//func (self *StaticDataMgr) GetStaticData(tableName string) StaticData { ///得到球员类型信息
//	self.GetStaticData()
//}

func (self *StaticDataMgr) GetStarType(starType int) *StarTypeStaticData { ///得到球员类型信息
	element := self.GetStaticData(tableStarType, starType)
	if nil == element {
		return nil
	}
	return element.(*StarTypeStaticData)
}

func (self *StaticDataMgr) GetSkillType(skillType int) *SkillTypeStaticData { ///得到技能类型信息
	element := self.GetStaticData(tableSkillType, skillType)
	if nil == element {
		return nil
	}
	return element.(*SkillTypeStaticData)
}

func (self *StaticDataMgr) GetLevelType(levelType int) *LevelTypeStaticData { ///得到关卡类型信息
	element := self.GetStaticData(tableLevelType, levelType)
	if nil == element {
		return nil
	}
	return element.(*LevelTypeStaticData)
}

func (self *StaticDataMgr) GetLeagueAwardType(leagueType int) *LeagueAwardTypeStaticData { ///得到推图奖励类型信息
	element := self.GetStaticData(tableLeagueAwardType, leagueType)
	if nil == element {
		return nil
	}
	return element.(*LeagueAwardTypeStaticData)
}

func (self *StaticDataMgr) GetAwardType(award int) *ActivitCodeAward { ///得到激活码礼包奖励类型信息
	element := self.GetStaticData(tableAwardType, award)
	if nil == element {
		return nil
	}
	return element.(*ActivitCodeAward)
}

func (self *StaticDataMgr) GetFormationType(formationType int) *FormationTypeStaticData { ///得到阵形类型信息
	staticDataMgr := GetServer().GetStaticDataMgr()
	element := staticDataMgr.GetStaticData(tableFormationType, formationType)
	if nil == element {
		return nil
	}
	return element.(*FormationTypeStaticData)
}

func (self *StaticDataMgr) GetSeatType(seatType int) *SeatTypeStaticData { ///得到位置类型信息
	staticDataMgr := GetServer().GetStaticDataMgr()
	element := staticDataMgr.GetStaticData(tableSeatType, seatType)
	if nil == element {
		return nil
	}
	return element.(*SeatTypeStaticData)
}

func GetTaskType(taskType int) *TaskTypeStaticData { ///得到npc球队类型信息
	staticDataMgr := GetServer().GetStaticDataMgr()
	element := staticDataMgr.GetStaticData(tableTaskType, taskType)
	if nil == element {
		return nil
	}
	return element.(*TaskTypeStaticData)
}

func (self *StaticDataMgr) GetNpcTeamType(npcTeamType int) *NpcTeamTypeStaticData { ///得到npc球队类型信息
	element := self.GetStaticData(tableNpcTeamType, npcTeamType)
	if nil == element {
		return nil
	}
	return element.(*NpcTeamTypeStaticData)
}

func (self *StaticDataMgr) GetNpcTeamStarType(npcStarType int) *NpcStarTypeStaticData { ///得到npc球星类型信息
	element := self.GetStaticData(tableNpcStarType, npcStarType)
	if nil == element {
		return nil
	}
	return element.(*NpcStarTypeStaticData)
}

func (self *StaticDataMgr) GetStarFateType(starFateType int) *StarFateTypeStaticData { ///得到球员缘系统类型信息
	element := self.GetStaticData(tableStarFateType, starFateType)
	if nil == element {
		return nil
	}
	return element.(*StarFateTypeStaticData)
}

func (self *StaticDataMgr) GetStarLobbyType(starLobbyType int) *StarLobbyTypeStaticData {
	element := self.GetStaticData(tableStarLobbyType, starLobbyType)
	if nil == element {
		return nil
	}
	return element.(*StarLobbyTypeStaticData)
}

func (self *StaticDataMgr) GetVipShopItemInfo(commodityID int) *VipShopStaticData {
	element := self.GetStaticData(tableVipShopType, commodityID)
	if nil == element {
		return nil
	}
	return element.(*VipShopStaticData)
}

func (self *StaticDataMgr) GetTaskType(taskType int) *TaskTypeStaticData { ///得到任务类型信息
	element := self.GetStaticData(tableTaskType, taskType)
	if nil == element {
		return nil
	}
	return element.(*TaskTypeStaticData)
}

func (self *StaticDataMgr) GetTacticType(tacticType int) *TacticTypeStaticData { ///得到阵形战术类型信息
	element := self.GetStaticData(tableTacticType, tacticType)
	if nil == element {
		return nil
	}
	return element.(*TacticTypeStaticData)
}

func (self *StaticDataMgr) GetChallangeMatchType(challangeMatchType int) *ChallangeMatchStaticData { ///得到挑战赛类型信息
	element := self.GetStaticData(tableChallangeMatch, challangeMatchType)
	if nil == element {
		return nil
	}
	return element.(*ChallangeMatchStaticData)
}

func (self *StaticDataMgr) GetItemType(itemType int) *ItemTypeStaticData { ///得到道具类型信息
	element := self.GetStaticData(tableItemType, itemType)
	if nil == element {
		return nil
	}
	return element.(*ItemTypeStaticData)
}

func (self *StaticDataMgr) GetItemTypeMerge(itemType int) int { ///得到道具基础融合价值
	if self.staticDataList[tableItemType] == nil {
		return 0
	}
	staticData := self.staticDataList[tableItemType][itemType]
	if nil == staticData {
		return 0
	}
	return staticData.(*ItemTypeStaticData).Merge
}

func (self *StaticDataMgr) GetStarTypeBasePrice(starType int) int { ///得到球员基础身价
	if self.staticDataList[tableStarType] == nil {
		return 0
	}
	staticData := self.staticDataList[tableStarType][starType]
	if nil == staticData {
		return 0
	}
	return staticData.(*StarTypeStaticData).BasePrice
}

func (self *StaticDataMgr) GetStarTypeClass(starType int) int { ///得到球员评价
	if self.staticDataList[tableStarType] == nil {
		return 0
	}
	staticData := self.staticDataList[tableStarType][starType]
	if nil == staticData {
		return 0
	}
	return staticData.(*StarTypeStaticData).Class
}

func (self *StaticDataMgr) GetTeamMaxStarCount() int { ///判断球队中的雇佣球员数上限数
	const maxTeamStarCountIndex = 3
	maxStarCount := self.GetConfigStaticDataInt(configTeam, configItemDefaultTeamParam, maxTeamStarCountIndex)
	return maxStarCount
}

func (self *StaticDataMgr) GetStaticDataList(tableName string) StaticData { ///得到静态数据表
	staticData := self.staticDataList[tableName]
	return staticData
}

func (self *StaticDataMgr) GetAllTableIDList(tableName string) IntList { ///得到数据表所有记录id数组
	if self.staticDataList[tableName] == nil {
		return nil
	}
	idList := IntList{}
	for k, _ := range self.staticDataList[tableName] {
		idList = append(idList, k)
	}
	return idList
}

func (self *StaticDataMgr) GetStaticData(tableName string, dataID int) interface{} { ///得到静态数据复制指针
	if self.staticDataList[tableName] == nil {
		return nil
	}
	staticData := self.staticDataList[tableName][dataID]
	if nil == staticData {
		return nil
	}
	copyObj := reflect.New(reflect.TypeOf(staticData).Elem()).Elem()
	copyObj.Set(reflect.ValueOf(staticData).Elem())
	return copyObj.Addr().Interface()
}

func (self *StaticDataMgr) LoadStaticData(tableName string, structValue interface{}) { ///加载一张静态数据列表
	query := fmt.Sprintf("select * from %s", tableName)
	elementList := self.staticDB.fetchAllRows(query, structValue)
	if nil == elementList || len(elementList) <= 0 {
		GetServer().GetLoger().Fatal("StaticDataMgr LoadStaticData fetchAllRows fail! table:%s", tableName)
	}
	self.staticDataList[tableName] = make(StaticData)
	for i := range elementList {
		value := reflect.ValueOf(elementList[i]).Elem()
		id := int(value.FieldByName("ID").Int())
		descField := value.FieldByName("Desc")
		if descField.IsValid() == true {
			descField.SetString("N/A") ///如果有desc字段,将内容清掉
		}
		//descField = value.FieldByName("Name")
		//if descField.IsValid() == true {
		//	descField.SetString("N/A") ///如果有名字字段,将内容清掉
		//}
		self.staticDataList[tableName][id] = elementList[i]
	}
}

func (self *StaticDataMgr) GetConfigStaticDataInt(masterType string, subType string, paramIndex int) int {
	configStaticData := self.GetConfigStaticData(masterType, subType)
	if nil == configStaticData {
		return 0
	}
	configStaticDataReflect := reflect.ValueOf(configStaticData).Elem()
	fieldName := fmt.Sprintf("Param%d", paramIndex)
	fieldValue := configStaticDataReflect.FieldByName(fieldName)
	if fieldValue.IsValid() == false {
		return 0
	}
	result, _ := strconv.Atoi(fieldValue.String())
	return result
}

///得到服务器数据库静态配置表
func (self *StaticDataMgr) GetConfigStaticData(masterType string, subType string) *ConfigStaticData {
	var result *ConfigStaticData = nil
	for _, v := range self.staticDataList[tableConfig] {
		configData := v.(*ConfigStaticData)
		if configData.MasterType == masterType && configData.SubType == subType {
			result = configData
			break
		}
	}
	return result
}

func (self *StaticDataMgr) getConfigStaticDataParamIntList(masterType string, subType string) IntList { ///返回配置表的param整型数值列表
	result := []int{}
	configStaticData := self.GetConfigStaticData(masterType, subType)
	if nil == configStaticData {
		return result
	}
	configStaticDataReflect := reflect.ValueOf(configStaticData).Elem()
	for i := 1; i < 100; i++ {
		fieldName := fmt.Sprintf("Param%d", i)
		fieldValue := configStaticDataReflect.FieldByName(fieldName)
		if fieldValue.IsValid() == false {
			break
		}
		rateIntValue, _ := strconv.Atoi(fieldValue.String())
		result = append(result, rateIntValue)
	}
	return result
}

func (self *StaticDataMgr) getConfigStaticDataParamIntMap(masterType string, subType string) *map[int]int { ///返回配置表的param整型数值列表
	result := make(map[int]int)
	paramIntList := self.getConfigStaticDataParamIntList(masterType, subType)
	for i := range paramIntList {
		result[i+1] = paramIntList[i]
	}
	return &result
}

func (self *StaticDataMgr) GetVipShopStaticDataMap() VipShopMap {
	result := make(VipShopMap)
	for i, v := range self.staticDataList[tableVipShopType] {
		vipShopType := v.(*VipShopStaticData)
		result[i] = vipShopType
	}

	return result
}

func (self *StaticDataMgr) GetVipShopStaticDataList(type1 int) IntList {
	vipShopStaticDataList := IntList{}
	if self.staticDataList[tableVipShopType] == nil {
		return vipShopStaticDataList
	}
	for k, v := range self.staticDataList[tableVipShopType] {
		vipShopType := v.(*VipShopStaticData)
		if vipShopType.Type1 == type1 {
			vipShopStaticDataList = append(vipShopStaticDataList, k)
		}
	}
	return vipShopStaticDataList
}

func (self *StaticDataMgr) GetVipInfo(vipLevel int) *VipPrivilegeStaticData {
	if self.staticDataList[tableVipPrivilege] == nil {
		return nil
	}
	if self.staticDataList[tableVipPrivilege][vipLevel] != nil {
		return self.staticDataList[tableVipPrivilege][vipLevel].(*VipPrivilegeStaticData)
	}
	return nil
}

func (self *StaticDataMgr) GetVipNeedExp(vipLevel int) int {
	if self.staticDataList[tableVipPrivilege] == nil {
		return 0
	}

	for _, v := range self.staticDataList[tableVipPrivilege] {
		staticData := v.(*VipPrivilegeStaticData)
		if staticData.ID == vipLevel {
			return v.(*VipPrivilegeStaticData).Recharge
		}
	}

	return 0
}

func (self *StaticDataMgr) GetLevelFromNpcteamID(npcTeamID int) (int, int) { ///从NPC队伍id得到关卡id
	levelID, index := 0, 0
	for _, v := range self.staticDataList[tableLevelType] {
		levelType := v.(*LevelTypeStaticData)
		if levelType.Sid1 == npcTeamID {
			levelID = levelType.ID
			index = 1
		} else if levelType.Sid2 == npcTeamID {
			levelID = levelType.ID
			index = 2
		} else if levelType.Sid3 == npcTeamID {
			levelID = levelType.ID
			index = 3
		}
	}

	return levelID, index
}

//func (self *StaticDataMgr) FindSDKType(sdkName string) *SDKType {
//	if self.staticDataList[tableSDKType] == nil {
//		return nil
//	}
//	for _, v := range self.staticDataList[tableSDKType] {
//		staticData := v.(*SDKType)
//		if staticData.Name == sdkName {
//			return v.(*SDKType)
//		}
//	}
//	return nil
//}

func (self *StaticDataMgr) GetMoneyType(moneyID int) *MoneyType { ///得到套餐类型
	//staticDataMgr := GetServer().GetStaticDataMgr()
	element := self.GetStaticData(tableMoney, moneyID)
	if nil == element {
		return nil
	}
	return element.(*MoneyType)
}
