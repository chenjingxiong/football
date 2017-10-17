package football

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"time"
)

const (
	ArenaSettleAccounts = 3 ///三天一结算
)

const (
	MaxPlayerInGroup         = 18 ///每个分组最多人数
	refreshArenaDataInterval = 60 ///每位玩家刷新竞技场数据时间间隔
)

const (
	AwardPromoteD = 1
	AwardPromoteC = 2
	AwardPromoteB = 3 ///联赛晋级奖励掩码位
	AwardPromoteA = 4
	AwardPromoteS = 5
)

type ArenaType struct {
	ID              int    ///联赛类型id
	Name            string ///联赛名字
	GroupNum        int    ///拥有最大分组数
	NpcTeamType     int    ///联赛机器人球队类型
	UpAwardCoin     int    ///晋级奖励球币
	Upawardtactic   int    ///晋级奖励战术点数
	KeepAwardCoin   int    ///保级奖励球币
	KeepAwardTactic int    ///保级奖励战术点数
	DownAwardCoin   int    ///降级奖励球币',
	DownAwardTactic int    ///降级奖励战术点数',
	AwardItem       int    ///奖励随机宝箱
}

type ArenaInfo struct {
	ID               int    `json:"id"`               ///id编号
	ArenaType        int    `json:"arenatype"`        ///当前所在联赛类型
	GroupNum         int    `json:"group"`            ///所在分组号
	TeamID           int    `json:"teamid"`           ///球队编号
	TeamName         string `json:"teamname"`         ///球队名字
	TeamIcon         int    `json:"teamicon"`         ///队徽
	RemainMatchCount int    `json:"remainmatchcount"` ///剩余比赛次数
	PlayMask         int    `json:"playmask"`         ///已攻击球队掩码,1表示指定位已比赛,16个球队16个bit
	Score            int    `json:"score"`            ///联赛成绩积分
	WinCount         int    `json:"wincount"`         ///获胜次数
	DrawCount        int    `json:"drawcount"`        ///平局次数
	LostCount        int    `json:"lostcount"`        ///负场次数
	AwardTicket      int    `json:"awardticket"`      ///奖卷标识 0无领奖 1已领奖 2晋级奖 3保级奖 4降级奖
	LastArenaType    int    `json:"lastarenatype"`    ///上次联赛类型
	LastRank         int    `json:"lastrank"`         ///上次联赛排名
	UpdateUTC        int    `json:"updateutc"`        ///比分更新utc时间
	AwardMask        int    `json:"awardmask"`        ///晋级领取钻石奖励掩码
}

func (self *ArenaInfo) GetReflectValue() reflect.Value { ///得到反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

type ArenaInfoSlice []ArenaInfo       ///联赛信息数组
type ArenaInfoList map[int]*ArenaInfo ///竞技场信息列表

type ArenaMgr struct {
	GameMgr                       ///逻辑系统管理器
	DataUpdater                   ///数据更新组件
	arenaInfoList   ArenaInfoList ///与自己相关的竞技场信息列表
	selfArenaInfo   *ArenaInfo    ///本队联赛信息
	nextRefreshTime time.Time     ///下次更新时间
	arenaIDList     IntList       ////竞技场信息id列表,用于比赛掩码用
}

func (self *ArenaMgr) GetArenaInfoSlice() ArenaInfoSlice {
	arenaInfoSlice := ArenaInfoSlice{}
	for i := range self.arenaIDList {
		arenaID := self.arenaIDList[i]
		arenaInfo := self.arenaInfoList[arenaID]
		arenaInfoSlice = append(arenaInfoSlice, *arenaInfo)
	}
	return arenaInfoSlice
	//arenaInfoSlice := ArenaInfoSlice{}
	//for _, v := range self.arenaInfoList {
	//	arenaInfoSlice = append(arenaInfoSlice, *v)
	//}
	//return arenaInfoSlice
}

func (self *ArenaMgr) createDefault() { ///初次查询时创建默认记录
	const defaultArenaType = 7 ///默认竞技场类型
	createArenaInfoQuery := fmt.Sprintf("insert %s (teamid,teamname,teamicon,arenatype,groupnum) select %d,'%s',%d,%d,ceil((count(*)+1)/%d) from %s where arenatype=%d",
		tableArena, self.team.ID, self.team.Name, self.team.Icon, defaultArenaType, MaxPlayerInGroup, tableArena, defaultArenaType) ///组插入记录SQL
	lastInsertItemID, _ := GetServer().GetDynamicDB().Exec(createArenaInfoQuery)
	if lastInsertItemID <= 0 {
		GetServer().GetLoger().Warn("ArenaMgr createDefault fail! teamid:%d", self.team.ID)
	}
}

func (self *ArenaMgr) GetType() int { ///得到管理器类型
	return mgrTypeArenaMgr ///竞技场管理器
}

func (self *ArenaMgr) SaveInfo() { ///保存数据
	if nil == self.selfArenaInfo {
		return
	}
	self.Save()
}

func (self *ArenaMgr) GetInfo() *ArenaInfo {
	return self.selfArenaInfo
}

//func (self *ArenaMgr) RefreshTeamOrderSlice() { ///刷新球队顺序列表
//	for i := range arenaInfoList {
//		arenaInfoStore := arenaInfoList[i].(*ArenaInfo)
//		self.arenaInfoList[arenaInfoStore.ID] = arenaInfoStore
//	}
//}

func (self *ArenaMgr) RefreshArenaDate() bool {
	now := time.Now()
	isExpire := now.After(self.nextRefreshTime)
	if false == isExpire {
		return false ///未到刷新时间忽略返回
	}
	self.nextRefreshTime = GetExpireTime(refreshArenaDataInterval) ///一分钟更新一次
	if self.selfArenaInfo != nil {
		self.Save() ///保存上次数据
	}
	self.selfArenaInfo = nil
	self.arenaIDList = IntList{}
	self.arenaInfoList = ArenaInfoList{}
	arenaInfo := new(ArenaInfo)
	arenaInfoQuery := fmt.Sprintf("select * from %s where teamid=%d limit 1", tableArena, self.team.ID)
	result := GetServer().GetDynamicDB().fetchOneRow(arenaInfoQuery, arenaInfo)
	if false == result {
		///初始使用竞技场系统
		self.createDefault() ///初次查询时创建默认记录
		///再次查旬记录
		result = GetServer().GetDynamicDB().fetchOneRow(arenaInfoQuery, arenaInfo)
	}
	arenaInfoQueryList := fmt.Sprintf("select * from %s where arenatype=%d and groupnum=%d  limit %d", tableArena,
		arenaInfo.ArenaType, arenaInfo.GroupNum, MaxPlayerInGroup)
	arenaSelfID := arenaInfo.ID
	arenaInfoList := GetServer().GetDynamicDB().fetchAllRows(arenaInfoQueryList, arenaInfo)
	if nil == arenaInfoList {
		return false
	}
	for i := range arenaInfoList {
		arenaInfoStore := arenaInfoList[i].(*ArenaInfo)
		self.arenaInfoList[arenaInfoStore.ID] = arenaInfoStore
		self.arenaIDList = append(self.arenaIDList, arenaInfoStore.ID) ///将记录号放到索引数组中
	}
	//fmt.Println(self.arenaInfoList)
	sort.Ints(self.arenaIDList) ///对索引值进行升序排列
	//fmt.Println(self.arenaIDList)
	self.selfArenaInfo = self.arenaInfoList[arenaSelfID] ///保存指定自己的指针
	self.InitDataUpdater(tableArena, self.selfArenaInfo)
	return true
}

func (self *ArenaMgr) Init(teamID int) bool {
	return true
	//self.t
	//result := self.RefreshArenaDate()
	//arenaInfo := new(ArenaInfo)
	//arenaInfoQuery := fmt.Sprintf("select * from %s where teamid=%d limit 1", tableArena, teamID)
	//result := GetServer().GetDynamicDB().fetchOneRow(arenaInfoQuery, arenaInfo)
	//if false == result {
	//	///初始使用竞技场系统
	//	self.createDefault() ///初次查询时创建默认记录
	//	///再次查旬记录
	//	result = GetServer().GetDynamicDB().fetchOneRow(arenaInfoQuery, arenaInfo)
	//}
	//arenaInfoQueryList := fmt.Sprintf("select * from %s where arenatype=%d and group=%d limit %d", tableArena,
	//	arenaInfo.ArenaType, arenaInfo.GroupNum, MaxPlayerInGroup)
	//arenaInfoList := GetServer().GetDynamicDB().fetchAllRows(arenaInfoQueryList, arenaInfo)
	//if nil == arenaInfoList {
	//	return false
	//}
	//for i := range arenaInfoList {
	//	arenaInfoStore := arenaInfoList[i].(*ArenaInfo)
	//	self.arenaInfoList[arenaInfoStore.ID] = arenaInfoStore
	//}
	//self.selfArenaInfo = self.arenaInfoList[arenaInfo.ID] ///保存指定自己的指针
	//self.InitDataUpdater(tableArena, &self.selfArenaInfo)
	//return result
}

func NewArenaMgr(teamID int) IGameMgr {
	arenaMgr := new(ArenaMgr)
	if arenaMgr.Init(teamID) == false {
		return nil
	}
	return arenaMgr
}

func GetArenaType(arenaType int) *ArenaType {
	staticDataMgr := GetServer().GetStaticDataMgr()
	element := staticDataMgr.GetStaticData(tableArenaType, arenaType)
	if nil == element {
		return nil
	}
	return element.(*ArenaType)
}

func GetLevelupNumber(levelType int) int { //得到该晋级等级还能接纳人数
	server := GetServer()
	dyDBMgr := server.GetDynamicDB()
	arenaInfo := new(ArenaInfo)
	queryArenaSql := fmt.Sprintf("select * from %s where arenatype = %d", tableArena, levelType)
	arenaInfoLst := dyDBMgr.fetchAllRows(queryArenaSql, arenaInfo)
	return len(arenaInfoLst)
}

func CalcArenaResult() { ///结算竞争场成绩,纯数据库操作
	///根据每个玩家的成绩确定领奖字段数据
	//remainMatchCount := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configArenaMatch,
	//	configArenaMatchCommon, 1) ///确定每日重置比赛次数
	//arenaInfoCalcQuery := fmt.Sprintf(`set @rownum = 1;update %s as a,
	//(select id,@rownum:=CASE when @rownum=NULL then @rownum:=1 when @prev!=groupnum  then @rownum:=1 else @rownum+1 end as rank,
	//@prev:=groupnum as prev,
	//@newaward:=CASE when @rownum=1 then 2 when @rownum between 2 and 15 then 3 else 4 end as newaward,
	//CASE when @newaward=2 and arenatype>1 then arenatype-1 when @newaward=4 and arenatype<6 then arenatype+1 else arenatype end as newtype,
	//CASE when @newaward=2 and arenatype>1 then CEIL(groupnum/3) when @newaward=4 and arenatype<6 then groupnum+18-@rownum else groupnum end as newgroup,
	//arenatype as newlastrenatype
	//from dy_arena  order by arenatype desc,groupnum desc,score desc,updateutc asc) as t
	//set a.lastarenatype=t.newlastrenatype,a.arenatype=t.newtype,a.groupnum=t.newgroup,a.remainmatchcount=%d,
	//a.playmask=0,a.score=0,a.wincount=0,drawcount=0,lostcount=0,updateutc=0,a.awardticket=t.newaward,a.lastrank=t.rank
	//where a.id=t.id`, tableArena, remainMatchCount)
	//GetServer().GetDynamicDB().Exec(arenaInfoCalcQuery)

	// 修改: 联赛积分超过30分都可晋级。但晋级后联赛人数不可超过联赛上限人数。
	levelUpScore := 30
	server := GetServer()
	dyDBMgr := server.GetDynamicDB()

	// 普通晋级且积分暂不清零
	dyDBMgr.Exec("call CalcArenaResult()")

	// 查询没有晋级/降级的活跃用户
	arenaInfo := new(ArenaInfo)
	queryArenaSql := fmt.Sprintf("select * from %s where score >= %d and awardticket = 3", tableArena, levelUpScore)
	arenaInfoLst := dyDBMgr.fetchAllRows(queryArenaSql, arenaInfo)
	levelUpTeamList := []*ArenaInfo{}
	for i := 2; i <= 7; i++ { //S级联赛无晋级 取值为A-F联赛
		number := GetLevelupNumber(i - 1) //得到阶级现有人数
		//fmt.Printf("number : %d \r\n", number)
		number = int(math.Pow(3, float64(i-1)))*18 - number //得到可接纳人数
		//fmt.Printf("number : %d \r\n", number)
		for _, v := range arenaInfoLst {
			info := v.(*ArenaInfo)
			if info.ArenaType == i && number > 0 {
				levelUpTeamList = append(levelUpTeamList, info)

				number-- //可接纳人数-1
			}
		}
	}

	for _, v := range levelUpTeamList {
		groupNum := int(math.Ceil(float64(v.GroupNum) / 3)) //得到新组号
		updateArenaSql := fmt.Sprintf("update %s set lastarenatype = %d, arenatype = %d, groupnum = %d, awardticket = 2 where id = %d", tableArena, v.ArenaType, v.ArenaType-1, groupNum, v.ID)
		dyDBMgr.Exec(updateArenaSql)
	}

	//整理竞技场分组积分清零
	dyDBMgr.Exec("call ArrangeArena()")
}
