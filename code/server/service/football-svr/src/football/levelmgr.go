package football

import (
	"fmt"
)

type LeagueInfo struct { ///联赛信息
	Type            int  `json:"type"`            ///联赛类型
	LevelCount      int  `json:"levelcount"`      ///此联赛已交互过的关卡数
	StarCount       int  `json:"starcount"`       ///此联赛已获得星数
	IsPass          bool `json:"ispass"`          ///此联赛是否已通过,可以进下个联赛地图了
	ProcessAwardNum int  `json:"processawardnum"` ///此联赛的进度已领奖次数
	StarsAwardNum   int  `json:"starsawardnum"`   ///此联赛的星级已领奖次数
	LocksAwardNum   int  `json:"locksawardnum"`   ///此联赛的关卡锁已领奖次数
}

type LeagueInfoList []LeagueInfo       ///联赛信息列表
type LeagueInfoMap map[int]*LeagueInfo ///联赛信息列表

///关卡管理器
//type ILevelMgr interface {
//	IGameMgr
//	GetLeagueInfoList() LeagueInfoList                          ///得到联赛信息列表
//	GetLevelInfoList(leagueType int) LevelInfoList              ///得到关卡信息列表
//	GetLevel(levelID int) ILevel                                ///得到关卡对象
//	AddLevel(levelType int, starCount int) int                  ///添加关卡信息,成功返回关卡信息
//	FindLevel(levelType int) ILevel                             ///查找关卡对象
//	GetLevelStarCount(leagueType int) int                       //得到当前地图玩家拥有星级数量
//	GetSchedule(leagueType int) int                             //得到当前地图玩家总进度
//	GetLockNum(leagueType int) int                              //得到当前地图关卡锁数目
//	FindLeagueAward(leagueType int) *LeagueAwardInfo            ///得到当前地图奖励信息
//	AwardRecv(leagueInfo LeagueAwardInfo, awardType int)        ///设置奖励已领取
//	AwardIsRecv(leagueInfo LeagueAwardInfo, awardType int) bool ///判断奖励是否已领取
//	AddLeagueAwardInfo(infoID int, leagueAward *LeagueAward)    ///增加新奖励信息
//}

type LevelInfoList []LevelInfo
type LevelList map[int]*Level
type LevelAwardList map[int]*LeagueAward

type LevelMgr struct { ///关卡管理器
	GameMgr
	levelList      LevelList
	levelAwardList LevelAwardList
}

//func (self *LevelMgr) Save() { ///保存数据
//	for _, v := range self.levelList {
//		v.Save()
//	}
//}

func (self *LevelMgr) GetType() int { ///得到管理器类型
	return mgrTypeLevelMgr ///关卡管理器
}

func (self *LevelMgr) SaveInfo() { ///保存数据
	for _, v := range self.levelList {
		v.Save()
	}
	for _, v := range self.levelAwardList {
		v.Save()
	}
}

func (self *LevelMgr) FindLevel(levelType int) *Level { ///查找关卡对象
	for _, v := range self.levelList {
		levelInfoPtr := v.GetInfoPtr()
		if levelInfoPtr.Type == levelType {
			return v
		}
	}
	return nil
}

func (self *LevelMgr) GetLevel(levelID int) *Level { ///得到关卡对象
	return self.levelList[levelID]
}

func (self *LevelMgr) GetLevelInfoList(leagueType int) LevelInfoList { ///得到关卡信息列表
	levelInfoList := LevelInfoList{}
	staticDataMgrUnsafe := GetServer().GetStaticDataMgr().Unsafe()
	for _, v := range self.levelList {
		levelInfoPtr := v.GetInfoPtr()
		levelType := staticDataMgrUnsafe.GetLevelType(levelInfoPtr.Type)
		if levelType.League == leagueType {
			levelInfoList = append(levelInfoList, *levelInfoPtr)
		}
	}
	return levelInfoList
}

func (self *LevelMgr) GetLeagueInfoList() LeagueInfoList { ///得到联赛信息列表
	leagueInfoMap := make(LeagueInfoMap)
	staticDataMgrUnsafe := GetServer().GetStaticDataMgr().Unsafe()

	leagueTypeIDList := self.GetLeagueTypeIDList()
	for index := range leagueTypeIDList {
		leagueType := leagueTypeIDList[index]
		leagueInfoPtr := new(LeagueInfo)
		leagueInfoPtr.Type = leagueType
		leagueInfoMap[leagueType] = leagueInfoPtr
	}

	for _, v := range self.levelList {
		levelInfoPtr := v.GetInfoPtr()
		levelTypeStaticData := staticDataMgrUnsafe.GetLevelType(levelInfoPtr.Type)
		if nil == levelTypeStaticData {
			continue
		}
		leagueInfoPtr := leagueInfoMap[levelTypeStaticData.League]
		if nil == leagueInfoPtr {
			leagueInfoPtr = new(LeagueInfo)
			leagueInfoPtr.Type = levelTypeStaticData.League ///相同的type
			leagueInfoMap[levelTypeStaticData.League] = leagueInfoPtr
		}
		if levelTypeStaticData.Type == levelTypeFinish {
			leagueInfoPtr.IsPass = true ///已开终点关卡,整个联赛为通关状态
		}
		if levelTypeStaticData.Type != levelTypeMatch && levelTypeStaticData.Type != levelTypeLocker {
			continue
		}

		if levelTypeStaticData.Type == levelTypeLocker { //! 关卡锁必然 +1
			leagueInfoPtr.LevelCount += 1
		} else { //! 非关卡锁加入星
			leagueInfoPtr.LevelCount += (levelInfoPtr.StarCount + 2) / 3
		}
		leagueInfoPtr.StarCount += levelInfoPtr.StarCount
		leagueAward := self.FindLeagueAward(levelTypeStaticData.League)
		if leagueAward != nil {
			leagueInfoPtr.ProcessAwardNum = leagueAward.ProcessAwardNum
			leagueInfoPtr.StarsAwardNum = leagueAward.StarsAwardNum
			leagueInfoPtr.LocksAwardNum = leagueAward.LocksAwardNum
		}
	}
	///将map转成slice
	leagueInfoList := LeagueInfoList{}
	for _, leagueInfo := range leagueInfoMap {
		leagueInfoList = append(leagueInfoList, *leagueInfo)
	}
	return leagueInfoList
}

func NewLevelMgr(teamID int) IGameMgr {
	levelMgr := new(LevelMgr)
	levelMgr.levelList = make(LevelList)
	//	levelMgr.teamID = teamID ///存放自己的球队id
	levelListQuery := fmt.Sprintf("select * from %s where teamid=%d limit 3000", tableLevel, teamID)
	levelInfo := new(LevelInfo)
	levelInfoList := GetServer().GetDynamicDB().fetchAllRows(levelListQuery, levelInfo)
	for i := range levelInfoList {
		levelInfo = levelInfoList[i].(*LevelInfo)
		levelMgr.levelList[levelInfo.ID] = NewLevel(levelInfo)
	}

	levelMgr.levelAwardList = make(LevelAwardList)
	leagueAwardListQuery := fmt.Sprintf("select * from %s where teamid = %d limit 99", tableLeagueAward, teamID)
	leagueAwardInfo := new(LeagueAwardInfo)
	leagueAwardList := GetServer().GetDynamicDB().fetchAllRows(leagueAwardListQuery, leagueAwardInfo)
	for i := range leagueAwardList {
		leagueAwardInfo = leagueAwardList[i].(*LeagueAwardInfo)
		levelMgr.levelAwardList[leagueAwardInfo.ID] = NewLeagueAward(leagueAwardInfo)
	}

	return levelMgr
}

func (self *LevelMgr) AddLevel(levelType int, starCount int) int { ///添加关卡信息,成功返回关卡对象id
	addLevelQuery := fmt.Sprintf("insert %s set teamid=%d,type=%d,starcount=%d",
		tableLevel, self.team.GetID(), levelType, starCount) ///组插入记录SQL
	lastInsertLevelID, _ := GetServer().GetDynamicDB().Exec(addLevelQuery)
	if lastInsertLevelID <= 0 {
		GetServer().GetLoger().Warn("LevelMgr AddLevel fail! levelType:%d starCount:%d", levelType, starCount)
		return 0
	}
	levelInfo := new(LevelInfo)
	levelInfo.ID = lastInsertLevelID
	levelInfo.TeamID = self.team.GetID()
	levelInfo.Type = levelType
	levelInfo.StarCount = starCount
	level := NewLevel(levelInfo) ///生成关卡对象
	self.levelList[lastInsertLevelID] = level
	return lastInsertLevelID
}

func (self *LevelMgr) GetLevelStarCount(leagueType int) int { //得到当前地图玩家拥有星级数量
	levelInfoList := self.GetLevelInfoList(leagueType)
	starCount := 0
	for i := range levelInfoList {
		levelInfo := levelInfoList[i]
		starCount += levelInfo.StarCount
	}

	return starCount
}

func (self *LevelMgr) GetSchedule(leagueType int) int { //得到当前地图玩家总进度
	staticDataMgr := GetServer().GetStaticDataMgr()
	levelInfoList := self.GetLevelInfoList(leagueType)
	processCount := 0
	for i := range levelInfoList {
		levelInfo := levelInfoList[i]
		levelType := staticDataMgr.GetLevelType(levelInfo.Type)
		if levelType.Type == levelTypeLocker { //! 关卡锁+1
			processCount += 1
		} else if levelType.Type == levelTypeMatch {
			processCount += (levelInfo.StarCount + 2) / 3
		}
	}
	return processCount
}

func (self *LevelMgr) GetLockNum(leagueType int) int { //得到当前地图关卡锁数目
	staticDataMgr := GetServer().GetStaticDataMgr()
	levelInfoList := self.GetLevelInfoList(leagueType)
	lockCount := 0
	for i := range levelInfoList {
		levelInfo := levelInfoList[i]
		levelType := staticDataMgr.GetLevelType(levelInfo.Type)
		if 2 == levelType.Type {
			lockCount += 1
		}

	}

	return lockCount
}

func (self *LevelMgr) FindLeagueAward(leagueType int) *LeagueAward {
	//	if self.levelAwardList[leagueType] == nil {
	//		return nil
	//	}
	for _, v := range self.levelAwardList {
		if v.Leaguetype == leagueType {
			return v
		}
	}
	return nil
}

//func (self *LevelMgr) AwardRecv(leagueInfo LeagueAwardInfo, awardType int) { ///设置奖励已领取
//	leagueInfo.RecvAward |= 1 << uint(awardType)
//}

//func (self *LevelMgr) AwardIsRecv(leagueInfo LeagueAwardInfo, awardType int) bool { ///判断奖励是否已领取
//	return (leagueInfo.RecvAward & 1 << uint(awardType)) == 1
//}

func (self *LevelMgr) AddLeagueAwardInfo(leagueType int) *LeagueAward {
	if nil == self.levelAwardList {
		self.levelAwardList = make(LevelAwardList)
	}
	leagueAwardInfo := new(LeagueAwardInfo)
	leagueAwardInfo.Leaguetype = leagueType
	leagueAwardInfo.ProcessAwardNum = 0
	leagueAwardInfo.StarsAwardNum = 0
	leagueAwardInfo.LocksAwardNum = 0
	leagueAwardInfo.TeamID = self.team.GetID()
	leagueAward := NewLeagueAward(leagueAwardInfo)
	insertAwardInfoSQL := leagueAward.InsertSql()
	lastInsertID, _ := GetServer().GetDynamicDB().Exec(insertAwardInfoSQL)
	leagueAward.ID = lastInsertID
	self.levelAwardList[lastInsertID] = leagueAward
	return leagueAward
}

func (self *LevelMgr) GetLeagueTypeIDList() IntList { ///得到推图奖励类型编号列表
	leagueTypeIDList := IntList{}
	staticDataList := GetServer().GetStaticDataMgr().GetStaticDataList(tableLeagueAwardType)
	for k, _ := range staticDataList {
		leagueTypeIDList = append(leagueTypeIDList, k)
	}
	return leagueTypeIDList
}

func (self *LevelMgr) GetHasPassNpcTeamList() IntList { ///得到所有已经打过的npcteam球队类型列表
	hasPassNpcTeamList := IntList{}
	for _, v := range self.levelList {
		if v.IsMatchLevel() == false {
			continue ///排除非比赛关卡
		}
		levelPassNpcTeamList := v.GetPassNpcTeamList()
		hasPassNpcTeamList = append(hasPassNpcTeamList, levelPassNpcTeamList...)
	}
	return hasPassNpcTeamList
}
