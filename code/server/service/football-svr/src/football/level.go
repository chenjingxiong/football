package football

import (
	"reflect"
)

///关卡信息

const (
	levelTypeMatch        = 1 ///比赛关卡
	levelTypeLocker       = 2 ///关卡锁
	levelTypeStart        = 3 ///起点
	levelTypeFinish       = 4 ///终点
	levelTypeStarCard     = 5 ///球员卡
	levelTypeFunctionCard = 6 ///功能卡
	levelTypeItemCard     = 7 ///道具卡
	levelTypeSkillCard    = 8 ///技能卡
)

const (
	leagueAwardStar1 = 1 //星级奖励
	leagueAwardStar2 = 2
	leagueAwardStar3 = 3
	leagueAwardLock1 = 4 //锁奖励
	leagueAwardLock2 = 5
	leagueAwardLock3 = 6
	leagueAwardChe1  = 7 //总进度奖励
	leagueAwardChe2  = 8
	leagueAwardChe3  = 9
)

const (
	leagueStarAward     = 1 //星级奖励
	leagueLockAward     = 2 //锁奖励
	leagueScheduleAward = 3 //总进度奖励
)

//type ILevel interface {
//	IDataUpdater
//	GetInfoPtr() *LevelInfo         ///得到信息接口
//	GetID() int                     ///得到编号
//	GetReflectValue() reflect.Value ///得到反射对象
//	IsPass() bool                   ///判断此关卡是否已通关
//}

type LevelInfo struct { ///关卡信息,对应动态表中的dy_level
	ID        int `json:"id"`        ///阵形id
	TeamID    int `json:"teamid"`    ///拥有关卡信息的球队id
	Type      int `json:"type"`      ///关卡类型',
	StarCount int `json:"starcount"` ///当前关卡星级进度
}

type LeagueAwardInfo struct {
	ID              int `json:"id"`              //ID
	TeamID          int `json:"teamid"`          //队伍ID
	Leaguetype      int `json:"leaguetype"`      //地图种类
	ProcessAwardNum int `json:"processawardnum"` //进度奖励已领次数
	StarsAwardNum   int `json:"starsawardnum"`   //星级奖励已领次数
	LocksAwardNum   int `json:"locksawardnum"`   //关卡锁奖励已领次数
}

type LeagueAward struct {
	LeagueAwardInfo
	DataUpdater
}

func (self *LeagueAward) GetID() int { ///得到反射对象
	return self.ID
}

func (self *LeagueAward) GetReflectValue() reflect.Value { ///得到反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

type LeagueAwardTypeStaticData struct {
	ID     int    //记录ID
	Name   string //名字
	Logo   int    //标志
	Stars1 int    //星级奖励达标要求
	Item1  int    //星级奖励
	Num1   int    //星级奖励数量
	Stars2 int
	Item2  int
	Num2   int
	Stars3 int
	Item3  int
	Num3   int
	Lock1  int //锁奖励达标要求
	Item4  int //锁奖励
	Num4   int //锁奖励数量
	Lock2  int
	Item5  int
	Num5   int
	Lock3  int
	Item6  int
	Num6   int
	Che1   int //总进度奖励达标要求
	Item7  int //总进度奖励
	Num7   int //总进度奖励数量
	Che2   int
	Item8  int
	Num8   int
	Che3   int
	Item9  int
	Num9   int
	Score  int    //推荐战力
	Desc   string //描述
}

type LevelTypeStaticData struct { ///关卡类型,对应动态表中的dy_leveltype
	ID         int    ///关卡类型id ID规则：Type*1000+顺序位
	Name       string ///关卡名字
	Type       int    ///关卡类型(1关卡 2关卡锁 3起点 4终点 5球员牌 6战术牌 7功能牌 8道具牌 9技能牌
	League     int    ///联赛
	Untie      int    ///前置关卡条件
	Logo       int    ///图标
	Sid1       int    ///参数1
	Sid2       int    ///参数2
	Sid3       int    ///参数3
	Sid4       int    ///参数4
	Sid5       int    ///参数5
	BlockColor int    ///块颜色
	Desc       string ///描述
}

type Level struct { ///关卡对象
	LevelInfo ///关卡信息
	DataUpdater
}

func (self *Level) GetReflectValue() reflect.Value { ///得到反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

func (self *Level) GetID() int {
	return self.ID
}

func (self *Level) IsPass() bool { ///判断此关卡是否已通关
	levelType := GetServer().GetStaticDataMgr().Unsafe().GetLevelType(self.Type)
	if nil == levelType {
		return false ///关卡类型无效
	}
	if levelTypeMatch != levelType.Type {
		return true ///非战斗关卡,有对象就算通关了
	}
	///得到前置关卡类型
	if self.StarCount < 1 {
		return false ///战斗关卡只有赢得1星以上才算通关
	}
	return true
}

func (self *Level) GetInfoPtr() *LevelInfo {
	return &self.LevelInfo
}

func (self *Level) GetType() *LevelTypeStaticData {
	levelTypeStaticData := GetServer().GetStaticDataMgr().GetLevelType(self.Type)
	return levelTypeStaticData
}

func (self *Level) IsMatchLevel() bool {
	levelTypeStaticData := self.GetType()
	isMatchLevel := levelTypeStaticData.Type == levelTypeMatch
	return isMatchLevel
}

func (self *Level) GetPassNpcTeamList() IntList { ///得到已经打开的npcteam类型
	levelTypeStaticData := self.GetType()
	passNpcTeamList := IntList{}
	if self.StarCount >= 1 && levelTypeStaticData.Sid1 > 0 {
		passNpcTeamList = append(passNpcTeamList, levelTypeStaticData.Sid1)
	}
	if self.StarCount >= 4 && levelTypeStaticData.Sid2 > 0 {
		passNpcTeamList = append(passNpcTeamList, levelTypeStaticData.Sid2)
	}
	if self.StarCount >= 7 && levelTypeStaticData.Sid3 > 0 {
		passNpcTeamList = append(passNpcTeamList, levelTypeStaticData.Sid3)
	}
	return passNpcTeamList
}

func NewLevel(levelInfo *LevelInfo) *Level {
	level := new(Level)
	level.LevelInfo = *levelInfo
	level.InitDataUpdater(tableLevel, &level.LevelInfo)
	return level
}

func NewLeagueAward(leagueAwardInfo *LeagueAwardInfo) *LeagueAward {
	leagueaward := new(LeagueAward)
	leagueaward.LeagueAwardInfo = *leagueAwardInfo
	leagueaward.InitDataUpdater(tableLeagueAward, &leagueaward.LeagueAwardInfo)
	return leagueaward
}

func (self *LevelTypeStaticData) GetNpcIndexInLevel(npcID int) int { ///取得npcteam在leveltype中的顺序位
	teamOrder := 0 ///球队在关卡中的位置
	if npcID == self.Sid1 {
		teamOrder = 1
	} else if npcID == self.Sid2 {
		teamOrder = 2
	} else if npcID == self.Sid3 {
		teamOrder = 3
	}
	return teamOrder
}
