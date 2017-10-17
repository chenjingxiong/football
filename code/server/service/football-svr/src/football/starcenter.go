package football

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	starCenterTypeBegin       = 1 ///开始符
	starCenterTypeDiscover    = 1 ///球探发掘出的球员所在的球员中心类型
	starCenterTypeLobby       = 1 ///球员游说中心类型
	starCenterTypeStarTenDraw = 2 ///球星十连抽
	starCenterTypeEnd             ///结束符

	starLobbyLimitMax = 64 ///球员游说上限
)

const (
	ItemCanvassCard = 200004 ///球星刷新卡
	ItemStarCard    = 200005 ///球星签约卡
)

const (
	classSS = 500 //评分决定等级
	classS  = 450
	classA  = 400
	classB  = 350
	classC  = 300
	classD  = 250
	classE  = 200
)

const StarTenDrawMaxCount = 3 ///十连抽最大十个球星 (修改: 最大抽3个球星)

type StarLobbyTypeStaticData struct { ///球员游说静态配置数据
	ID            int    ///记录id
	Name          string ///游说名称
	Awardstartype int    ///奖励球星ID
	Needstartype1 int    ///条件球员1
	Needstartype2 int    ///条件球员2
	Needstartype3 int    ///条件球员3
	Needstartype4 int    ///条件球员4
	Needstartype5 int    ///条件球员5
	Desc          string ///描述
}

//type IStarCenter interface { ///球队接口
//	GetStarCenterCount() int                                                                  ///得到球员中心个数
//	GetStarCenterMemberCount(starCenterType int) int                                          ///得到指定球员中心中会员的个数
//	AwardMember(centerType int, starType int) int                                             ///奖励一个球员到球员中心
//	AwardSpecialMember(centerType int, starType int, evolvelimit int, spyType int) int        ///奖励一个特殊球员(含有星级)到球员中心
//	GetStarCenterMember(centerType int, memberID int) *StarCenterMember                       ///得到球员中心中一位成员信息
//	GetDrawGroupFilterIndexList(centerType int, drawGroupIndexList []int) []int               ///得到已排除已经在球员中心中的球员类型索引列表
//	GetStarCenterMemberList(centerType int) []int                                             ///得到指定球员中心类型的成员列表
//	RemoveMember(centerType int, memberIDList IntList) bool                                   ///从球员中心删除一个成员
//	RollVolunteerStarTypeList(client IClient, needStarCount int, excludeList IntList) IntList ///从球星来投集合中抽选指定个数的球员
//	IsTypeExistStarCenter(centerType int, starType int) bool                                  ///判断球星是否存在于转会中心
//	GetStarEvolveCount(centerType int, starType int) int                                      ///得到球员默认星级
//	GetProbabilityEvolveCount(spyType int) int                                                ///根据概率得到球员随机星级
//}

type StarCenterMember struct { ///球员中心球员信息
	ID          int `json:"starcentermemberid"` ///记录id
	TeamID      int `json:"teamid"`             ///球队id
	Type        int `json:"type"`               ///球员中心类型 1表示球探发掘球员中心 2表示球星来投所在不球员中心
	StarType    int `json:"startype"`           ///球员类型
	ExpireTime  int `json:"expiretime"`         ///过期删除球员期限时间utc秒
	EvolveCount int `json:"evolvecount"`        ///转会中心球员星级
}

func (self *StarCenterMember) Save() {
	sqlUpdate := fmt.Sprintf("update %s set StarType=%d,ExpireTime=%d,EvolveCount=%d where id=%d",
		tableStarCenter, self.StarType, self.ExpireTime, self.EvolveCount, self.ID)
	GetServer().GetDynamicDB().Exec(sqlUpdate)
}

func (self *StarCenter) pureExpireMember(now int, client IClient) { ///球员中心删除已过期的球员
	const pureExpireMemberInterval = 60 ///每60秒清理一次过期球员
	nowTime := now % pureExpireMemberInterval
	if nowTime > 0 {
		return ///每60秒判断一次
	}
	removeMemberIDList := IntList{}
	for k, v := range self.starCenterList[starCenterTypeDiscover] {
		if now >= v.ExpireTime {
			removeMemberIDList = append(removeMemberIDList, k)
		}
	}
	if len(removeMemberIDList) > 0 {
		self.RemoveMember(starCenterTypeDiscover, removeMemberIDList)
		client.GetSyncMgr().SyncRemoveStarCenterMember(starCenterTypeDiscover, removeMemberIDList)
	}
}

func (self *StarCenter) Save() { ///保存数据
	///目前只保存了球员十连抽信息
	if self.starCenterList[starCenterTypeStarTenDraw] != nil {
		for _, v := range self.starCenterList[starCenterTypeStarTenDraw] {
			v.Save()
		}
	}
}

func (self *StarCenter) Update(now int, client IClient) { ///球队自身更新状态
	self.pureExpireMember(now, client)
}

type StarCenterMemberInfoList []StarCenterMember ///球员中心会员信息列表

type StarCenterMemberList map[int]*StarCenterMember ///球员中心会员列表
type StarCenterList map[int]StarCenterMemberList    ///球员中心列表
type StarCenter struct {                            ///球员中心,球队在此获得球员
	starCenterList StarCenterList ///球员中心列表
	teamID         int            ///本球员中心所隶属的球队id
}

func (self *StarCenter) GetStarCenterCount() int { ///得到球员中心个数
	starCenterCount := len(self.starCenterList)
	return starCenterCount
}

func (self *StarCenter) GetStarCenterMemberCount(starCenterType int) int { ///得到指定球员中心中会员的个数
	starCenterMemberCount := 0
	if self.starCenterList[starCenterType] != nil {
		starCenterMemberCount = len(self.starCenterList[starCenterType])
	}
	return starCenterMemberCount
}

func (self *StarCenter) GetDrawGroupFilterIndexList(centerType int, drawGroupIndexList []int) []int { ///得到已排除已经在球员中心中的球员类型索引列表
	if nil == self.starCenterList[centerType] {
		return drawGroupIndexList
	}
	staticDataMgr := GetServer().GetStaticDataMgr()
	///建立查询索引
	storeStarTypeList := make(map[int]bool)
	for _, v := range self.starCenterList[centerType] {
		storeStarTypeList[v.StarType] = true
	}
	drawGroupFilterIndexList := []int{}
	for i := range drawGroupIndexList {
		drawGroupIndex := drawGroupIndexList[i]
		drawGroupStaticData := staticDataMgr.GetDrawGroupStaticData(drawGroupIndex)
		if nil == drawGroupStaticData {
			GetServer().GetLoger().Warn("StarCenter GetDrawGroupFilterIndexList fail! drawGroupIndex %d is Invalid!", drawGroupIndex)
			break
		}
		_, ok := storeStarTypeList[drawGroupStaticData.AwardType]
		if false == ok {
			drawGroupFilterIndexList = append(drawGroupFilterIndexList, drawGroupIndex)
		}
	}
	return drawGroupFilterIndexList
}

///得到指定球员中心类型的成员列表
func (self *StarCenter) GetStarCenterMemberList(centerType int) []int {
	memberList := []int{}
	if nil == self.starCenterList[centerType] {
		return []int{}
	}
	for k, _ := range self.starCenterList[centerType] {
		memberList = append(memberList, k)
	}
	return memberList
}

func (self *StarCenter) GetStarCenterMember(centerType int, memberID int) *StarCenterMember { ///得到球员中心中一位成员信息
	if nil == self.starCenterList[centerType] {
		return nil
	}
	return self.starCenterList[centerType][memberID]
}

func (self *StarCenter) RemoveMember(centerType int, memberIDList IntList) bool { ///从球员中心删除一个成员
	///先删除内存防止玩家刷
	if nil == self.starCenterList[centerType] {
		return false
	}
	memberIDListLen := len(memberIDList)
	if memberIDListLen <= 0 {
		return false
	}
	///从数据库中删除对象
	removeMemberQuery := fmt.Sprintf("delete from %s where id in (", tableStarCenter)
	for i := range memberIDList {
		memberID := memberIDList[i]
		_, ok := self.starCenterList[centerType][memberID]
		if false == ok {
			return false
		}
		delete(self.starCenterList[centerType], memberID) ///从内存中删除对象
		removeMemberQuery += fmt.Sprintf("%d", memberID)
		if i < memberIDListLen-1 {
			removeMemberQuery += ","
		}
	}
	removeMemberQuery += fmt.Sprintf(") limit %d", memberIDListLen)
	_, rowsStarAffected := GetServer().GetDynamicDB().Exec(removeMemberQuery)
	if rowsStarAffected <= 0 {
		GetServer().GetLoger().Warn("StarCenter RemoveMember removeMemberQuery fail! centerType:%d memberIDList:%v",
			centerType, memberIDList)
		return false
	}
	return rowsStarAffected == memberIDListLen
}

func (self *StarCenter) GetStarEvolveCount(centerType int, starType int) int {
	if nil == self.starCenterList[centerType] {
		return 0
	}
	for _, v := range self.starCenterList[centerType] {
		if v.StarType == starType {
			return v.EvolveCount
		}
	}

	return 0
}

func (self *StarCenter) AwardMember(centerType int, starType int) int { ///奖励一个球员到球员中心
	defaultExpireTime := 259200 ///默认球员中心成员失效时间///72小时,259200秒
	configExpireTime := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configStarCenter,
		configItemStarCenterCommonConfig, 1)
	if configExpireTime > 0 {
		defaultExpireTime = configExpireTime
	}
	memberExpireTime := int(time.Now().Unix()) + defaultExpireTime ///生成过期时间
	awardMemberQuery := fmt.Sprintf("insert %s set teamid=%d,type=%d,startype=%d,expiretime=%d",
		tableStarCenter, self.teamID, centerType, starType, memberExpireTime)
	lastInsertAwardMemberID, _ := GetServer().GetDynamicDB().Exec(awardMemberQuery)
	if lastInsertAwardMemberID <= 0 {
		GetServer().GetLoger().Warn("StarCenter AwardMember fail! centerType:%d starType:%d", centerType, starType)
		return 0
	}
	starCenterMember := new(StarCenterMember)
	starCenterMember.ID = lastInsertAwardMemberID
	starCenterMember.TeamID = self.teamID
	starCenterMember.Type = centerType
	starCenterMember.StarType = starType
	starCenterMember.ExpireTime = memberExpireTime
	if nil == self.starCenterList[centerType] {
		self.starCenterList[centerType] = make(StarCenterMemberList)
	}
	self.starCenterList[centerType][lastInsertAwardMemberID] = starCenterMember
	return lastInsertAwardMemberID
}

func (self *StarCenter) GetProbabilityEvolveCount(spyType int) int {
	starEvolveCount := 0
	evolvelimit := 0
	evolveCountList := []int{3, 5, 7} ///球探决定最高星级
	if spyType >= 1 && spyType <= 4 {
		evolvelimit = evolveCountList[spyType-1]
	} else {
		return 0
	}

	/// 得到一级星级概率
	fOneCountProbability := math.Pow(3.0, float64(evolvelimit-1)) / ((math.Pow(3.0, float64(evolvelimit)) - 1.0) / 2.0)
	/// 计算球员抽取星级概率
	fProbabilitylist := make([]float64, 7)
	iProbabilitylist := make([]int, 7)
	fProbabilitylist[0] = fOneCountProbability
	iCurProbability := 0
	///计算所有概率并转换为万分比
	for i := 1; i < evolvelimit; i++ {
		fProbabilitylist[i] = fProbabilitylist[i-1] / 3.0
		iProbabilitylist[i-1] = int(fProbabilitylist[i-1] * 10000)
		iCurProbability += iProbabilitylist[i-1]
		if i == evolvelimit-1 {
			iProbabilitylist[i] = 10000 - iCurProbability
		}
	}

	///根据概率随机出当前球星星级
	randNumber := Random(0, 10000)
	//randNumber := 9875
	currentProbability := 0
	nTemp := 0
	for i := 0; i < evolvelimit; i++ {
		nTemp = iProbabilitylist[i] + currentProbability
		if randNumber < nTemp && randNumber >= currentProbability {
			starEvolveCount = i + 1
			break
		}
		currentProbability += iProbabilitylist[i]
	}

	return starEvolveCount
}

func (self *StarCenter) AwardSpecialMember(centerType int, starType int, evolvelimit int, spyType int) int { ///奖励一个特殊球员到球员中心
	defaultExpireTime := 259200 ///默认球员中心成员失效时间///72小时,259200秒
	configExpireTime := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configStarCenter,
		configItemStarCenterCommonConfig, 1)
	if configExpireTime > 0 {
		defaultExpireTime = configExpireTime
	}

	starEvolveCount := 0
	if spyType == 0 {
		starEvolveCount = evolvelimit ///当spyType==0时 evolvelimit 为必出星级
	} else {
		if evolvelimit == 0 {
			return 0
		}

		starEvolveCount = self.GetProbabilityEvolveCount(spyType)
	}

	memberExpireTime := int(time.Now().Unix()) + defaultExpireTime ///生成过期时间
	awardMemberQuery := fmt.Sprintf("insert %s set teamid=%d,type=%d,startype=%d,expiretime=%d, evolvecount=%d",
		tableStarCenter, self.teamID, centerType, starType, memberExpireTime, evolvelimit)
	lastInsertAwardMemberID, _ := GetServer().GetDynamicDB().Exec(awardMemberQuery)
	if lastInsertAwardMemberID <= 0 {
		GetServer().GetLoger().Warn("StarCenter AwardMember fail! centerType:%d starType:%d", centerType, starType)
		return 0
	}

	starMemberQuery := fmt.Sprintf("select * from %s where id=%d limit 1", tableStarCenter, lastInsertAwardMemberID)
	starCenterMember := new(StarCenterMember)
	GetServer().GetDynamicDB().fetchOneRow(starMemberQuery, starCenterMember)
	starCenterMember.EvolveCount = starEvolveCount
	if self.starCenterList[centerType] == nil {
		self.starCenterList[centerType] = make(StarCenterMemberList)
	}
	self.starCenterList[centerType][lastInsertAwardMemberID] = starCenterMember

	return lastInsertAwardMemberID
}

func (self *StarCenter) Init(teamID int) bool { ///初始化球员中心
	self.starCenterList = make(StarCenterList)
	self.teamID = teamID ///存放自己的球队id
	starCenterQuery := fmt.Sprintf("select * from %s where teamid=%d limit 3000", tableStarCenter, teamID)
	starCenterMember := new(StarCenterMember)
	starCenterMemberList := GetServer().GetDynamicDB().fetchAllRows(starCenterQuery, starCenterMember)
	if nil == starCenterMemberList {
		return true
	}
	for i := range starCenterMemberList {
		starCenterMember = starCenterMemberList[i].(*StarCenterMember)
		centerType := starCenterMember.Type
		memberID := starCenterMember.ID
		if nil == self.starCenterList[centerType] {
			self.starCenterList[centerType] = make(StarCenterMemberList)
		}
		self.starCenterList[centerType][memberID] = starCenterMember
	}
	return true
}

func (self *StarCenter) FilterVolunteerStarTypeList(drawGroupIndexList IntList, excludeList IntList) IntList {
	if nil == excludeList {
		return drawGroupIndexList
	}
	staticDataMgr := GetServer().GetStaticDataMgr()
	resultList := IntList{}
	starTypeList := make(map[int]bool)
	for i := range excludeList {
		starType := excludeList[i]
		starTypeList[starType] = true
	}
	for k := range drawGroupIndexList {
		drawGroupIndex := drawGroupIndexList[k]
		drawGroupStaticData := staticDataMgr.GetDrawGroupStaticData(drawGroupIndex)
		if nil == drawGroupStaticData {
			break
		}
		_, found := starTypeList[drawGroupStaticData.AwardType]
		if true == found {
			continue
		}
		resultList = append(resultList, drawGroupIndex)
	}
	return resultList
}

///随机生成一组球星来投球星类型列表
func (self *StarCenter) RollVolunteerStarTypeList(client IClient, needStarCount int, excludeList IntList) IntList {
	team := client.GetTeam()
	levelMgr := team.GetLevelMgr()
	staticDataMgr := GetServer().GetStaticDataMgr()
	drawGroupIndexList := staticDataMgr.GetDrawGroupIndexList(drawGroupTypeVolunteer)      ///得到索引列表
	drawGroupIndexList = self.FilterVolunteerStarTypeList(drawGroupIndexList, excludeList) ///处理排除列表
	//drawGroupIndexList = team.GetDrawGroupFilterIndexList(drawGroupIndexList)              ///过滤掉已在球队中的球员列表
	drawGroupIndexListLen := len(drawGroupIndexList) ///得到索引列表长度
	if drawGroupIndexListLen < needStarCount {
		needStarCount = drawGroupIndexListLen ///修正需求个数为样本数
	}
	starCount := 0
	resultList := IntList{}
	for i := 0; i < drawGroupIndexListLen; i++ {
		diceDraw := rand.Intn(drawGroupIndexListLen) ///生成抽取随机值
		drawGroupIndex := drawGroupIndexList[diceDraw]
		drawGroupStaticData := staticDataMgr.GetDrawGroupStaticData(drawGroupIndex)
		if nil == drawGroupStaticData {
			GetServer().GetLoger().Warn("QueryVolunteerInfoMsg rollVolunteerStarTypeList fail! drawGroupIndex %d is Invalid!", drawGroupIndex)
			break
		}

		if drawGroupStaticData.ReqPass != 0 {

			//得到前置球队对应关卡
			levelID, index := staticDataMgr.GetLevelFromNpcteamID(drawGroupStaticData.ReqPass)
			if levelID == 0 {
				GetServer().GetLoger().Warn("noone in level reward  RollVolunteerStarTypeList() levelType = %d", drawGroupStaticData.ReqPass)
				continue
			}

			level := levelMgr.FindLevel(levelID)
			if level == nil {
				//GetServer().GetLoger().Warn("level is not exist func:RollVolunteerStarTypeList levelType = %d level", drawGroupStaticData.ReqPass)
				continue
			}
			if level.IsPass() == false { //未通过关卡,则淘汰
				continue
			}

			levelInfo := level.GetInfoPtr()
			starCount := index * 3
			if levelInfo.StarCount < starCount {
				continue
			}
		}
		starCount++
		resultList = append(resultList, drawGroupStaticData.AwardType)                                 ///放入抽中的球星类型
		drawGroupIndexListLen--                                                                        ///减去已抽取的
		drawGroupIndexList = append(drawGroupIndexList[:diceDraw], drawGroupIndexList[diceDraw+1:]...) ///去掉被抽中的项
		if starCount >= needStarCount {
			break
		}
	}
	resultListLen := len(resultList)
	if resultListLen != needStarCount {
		return nil
	}
	return resultList
}

func (self *StarCenter) IsTypeExistStarCenter(centerType int, starType int) bool {
	if nil == self.starCenterList[centerType] {
		return false
	}
	for _, v := range self.starCenterList[centerType] {
		if v.StarType == starType {
			return true
		}
	}
	return false
}

///得到指定球员中心类型的成员列表
func (self *StarCenter) GetStarCenterMemberInfoList(centerType int) StarCenterMemberInfoList {
	starCenterMemberInfoList := StarCenterMemberInfoList{}
	if nil == self.starCenterList[centerType] {
		return starCenterMemberInfoList
	}
	for _, v := range self.starCenterList[centerType] {
		starCenterMemberInfoList = append(starCenterMemberInfoList, *v)
	}
	return starCenterMemberInfoList
}

func (self *StarCenter) createStarTenDraw() {
	///创建初始的球星十连抽信息
	for i := 1; i <= 10; i++ {
		///放十个不存在的球员进去,通过后面刷新逻辑纠正数据
		self.AwardMember(starCenterTypeStarTenDraw, i)
	}
}

func (self *StarCenter) CalcStarTenDrawStarCount() int { ///计算球星十连抽星级
	//evolveCount := Random(3, 7)
	//return evolveCount

	///简化公式3^(7-m)/40  m = 星级 星级上限为7 m >= 4
	///概率 0.675 0.225 0.075 0.025
	iProbabilityList := []int{669, 225, 74, 24, 8}
	nTemp := 0
	randNum := Random(0, 1000)
	evolveCount := 3
	for i := 3; i <= 7; i++ {
		if randNum >= nTemp && randNum < nTemp+iProbabilityList[i-3] {
			evolveCount = i
			break
		}
		nTemp += iProbabilityList[i-3]
	}
	return evolveCount
}

func (self *StarCenter) GetRemainTenDrawStarCount() int { ///计算十连抽剩余未取得球星数
	remainTenDrawStarCount := 0
	for _, v := range self.starCenterList[starCenterTypeStarTenDraw] {
		if 0 == v.ExpireTime {
			remainTenDrawStarCount++
		}
	}
	return remainTenDrawStarCount
}

///更新一批球星十连抽
func (self *StarCenter) UpdateStarTenDraw() {
	drawResultList := IntList{}
	staticDataMgr := GetServer().GetStaticDataMgr()
	///得到十连抽已有球员数
	starCenterMemberCount := self.GetStarCenterMemberCount(starCenterTypeStarTenDraw)
	if starCenterMemberCount <= 0 {
		self.createStarTenDraw() ///创建初始的球星十连抽信息
	}
	starCenterMemberCount = self.GetStarCenterMemberCount(starCenterTypeStarTenDraw)
	///从A组抽一个S级球员
	drawGroupListRare := staticDataMgr.GetDrawGroupIndexList(drawGroupTypeStarTenDrawSuper)                     ///取得抽卡索引列表
	totalTakeWeightRare, totalShowWeightRare := discoverGetDrawWeightTotal(drawGroupListRare)                   ///得到权重总和
	drawTakeRare, _ := discoverDrawOne(&drawGroupListRare, &totalTakeWeightRare, &totalShowWeightRare, 0, true) ///抽取稀有卡
	///从B组抽9个杂鱼(修改: 抽2个球员)
	drawResultList = append(drawResultList, drawTakeRare) ///放入抽取普通球员type
	drawGroupListNormal := staticDataMgr.GetDrawGroupIndexList(drawGroupTypeStarTenDrawNormal)
	totalTakeWeightNormal, totalShowWeightNormal := discoverGetDrawWeightTotal(drawGroupListNormal) ///得到权重总和
	for i := 2; i <= StarTenDrawMaxCount; i++ {                                                     ///一共抽剩下9张卡
		drawShowOne, _ := discoverDrawOne(&drawGroupListNormal, &totalTakeWeightNormal, &totalShowWeightNormal, 0, true) ///后面四张为展示
		if 0 == drawShowOne {
			break
		}
		drawResultList = append(drawResultList, drawShowOne) ///取得结果
	}

	///从C组抽7个道具
	drawGroupListItem := staticDataMgr.GetDrawGroupIndexList(drawGroupTypeTenDrawItem)
	totalTakeWeightItem, totalShowWeightItem := discoverGetDrawWeightTotal(drawGroupListItem) ///得到权重总和
	itemNumLst := IntList{}
	for i := 3; i < 9; i++ {
		drawShowOne, drawShowNum := discoverDrawOne(&drawGroupListItem, &totalTakeWeightItem, &totalShowWeightItem, 0, true)
		if 0 == drawShowOne {
			break
		}

		drawResultList = append(drawResultList, drawShowOne) ///抽取七个道具并记录对应数量
		itemNumLst = append(itemNumLst, drawShowNum)
	}

	drawGroupListSP := staticDataMgr.GetDrawGroupIndexList(drawGroupTypeMannaStar)
	totalTakeWeightSP, totalShowWeightSP := discoverGetDrawWeightTotal(drawGroupListSP) ///得到权重总和
	drawShowOne, drawShowNum := discoverDrawOne(&drawGroupListSP, &totalTakeWeightSP, &totalShowWeightSP, 0, true)
	//	fmt.Println("drawGroupListSP:", drawGroupListSP, " totalTakeWeightSP:", totalTakeWeightSP, " totalShowWeightSP:", totalShowWeightSP)
	if 0 != drawShowOne {
		drawResultList = append(drawResultList, drawShowOne) ///抽取七个道具并记录对应数量
		itemNumLst = append(itemNumLst, drawShowNum)

		//		fmt.Println("特殊道具:", drawShowOne, " 数量:", drawShowNum)
	}

	//fmt.Println(drawResultList)

	drawResultListLen := drawResultList.Len()

	//fmt.Println(drawResultListLen)
	if drawResultListLen != 10 {
		return ///长度不符不能更新
	}
	newStarTypeIndex := 0
	newItemNumIndex := 0
	for _, v := range self.starCenterList[starCenterTypeStarTenDraw] {
		v.StarType = drawResultList[newStarTypeIndex]

		if v.StarType < 110000 {
			v.EvolveCount = self.CalcStarTenDrawStarCount()
		} else {
			//!若抽到的是物品,则星级 = 物品的数量
			v.EvolveCount = itemNumLst[newItemNumIndex]
			newItemNumIndex++
		}
		v.ExpireTime = 0
		newStarTypeIndex++ ///取下一个球星
	}

	//fmt.Println(self.starCenterList)
}
