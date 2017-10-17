package football

import (
	"fmt"
	"reflect"
)

type MannaStarType struct { ///自创球员属性信息
	ID              int    `json:"id"`              ///球员类型id
	Teamid          int    `json:"teamid"`          ///所属球队id
	MannaSeat       int    `json:"mannaseat"`       //!天赐球星所在位置
	Name            string `json:"name"`            ///球员名
	Grade           int    `json:"grade"`           ///品质
	Card            int    `json:"card"`            ///卡类
	Icon            int    `json:"icon"`            ///头像
	Face            int    `json:"face"`            ///外观
	Seat1           int    `json:"seat1"`           ///踢球位置1
	Seat2           int    `json:"seat2"`           ///踢球位置2
	Seat3           int    `json:"seat3"`           ///踢球位置3
	Nationality     string `json:"nationality"`     ///国藉
	Pass            int    `json:"pass"`            ///传球
	Steals          int    `json:"steals"`          ///抢断
	Dribbling       int    `json:"dribbling"`       ///盘带
	Sliding         int    `json:"sliding"`         ///铲球
	Shooting        int    `json:"shooting"`        ///射门
	GoalKeeping     int    `json:"goalkeeping"`     ///守门
	Body            int    `json:"body"`            ///身体值
	Speed           int    `json:"speed"`           ///速度
	PassGrow        int    `json:"passgrow"`        ///传球成长
	StealsGrow      int    `json:"stealsgrow"`      ///抢断成长
	DribblingGrow   int    `json:"dribblinggrow"`   ///盘带成长
	SlidingGrow     int    `json:"slidinggrow"`     ///铲球成长
	ShootingGrow    int    `json:"shootinggrow"`    ///射门成长
	GoalKeepingGrow int    `json:"goalkeepinggrow"` ///守门成长
	Skill1          int    `json:"skill1"`          ///初始技能1
	Skill2          int    `json:"skill2"`          ///初始技能2
	Skill3          int    `json:"skill3"`          ///初始技能3
	Skill4          int    `json:"skill4"`          ///初始技能4
	BasePrice       int    `json:"baseprice"`       ///基础身价
	BaseScore       int    `json:"basescore"`       ///基础评分
	Ticket          int    `json:"ticket"`          ///球票价格
	Fate1           int    `json:"fate1"`           ///球员缘类型1
	Fate2           int    `json:"fate2"`           ///球员缘类型2
	Fate3           int    `json:"fate3"`           ///球员缘类型3
	Fate4           int    `json:"fate4"`           ///球员缘类型4
	Fate5           int    `json:"fate5"`           ///球员缘类型5
	Fate6           int    `json:"fate6"`           ///球员缘类型6
	Item            int    `json:"item"`            ///升星所需道具
	Team            int    `json:"team"`            ///升星所需道具获取途径
	Hair            int    `json:"hair"`            ///头发
	Eyebrow         int    `json:"eyebrow"`         ///眉毛 (字段无用)
	Mouth           int    `json:"mouth"`           ///嘴巴
	Eye             int    `json:"eye"`             ///眼睛
	Skin            int    `json:"skin"`            ///皮肤
	Clothes         int    `json:"clothes"`         ///衣服
	Desc            string `json:"desc"`            ///球员描述
}

type MannaStarLst map[int]*MannaStar
type MannaStarSlice []*MannaStarType

type MannaStar struct {
	MannaStarType
	DataUpdater ///信息更新组件
}

func (self *MannaStar) GetReflectValue() reflect.Value { ///得到球队反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

type MannaStarMgr struct {
	GameMgr
	starList MannaStarLst
}

func (self *MannaStarMgr) SaveInfo() { ///保存数据
	if self.starList != nil {
		for _, v := range self.starList {
			v.Save()
		}
	}
}

//!根据自创球员类型返回自创球员信息
func (self *MannaStarMgr) GetMannaStar(starID int) *MannaStar {

	for _, v := range self.starList {
		if v.ID == starID {
			return v
		}
	}

	for _, v := range self.starList {
		fmt.Println(v.MannaStarType)
	}

	return nil
}

//!根据自创球员的位置得到球员信息
func (self *MannaStarMgr) GetMannaStarFromSeat(seat int) *MannaStar {
	return self.starList[seat]
}

//!得到管理器类型
func (self *MannaStarMgr) GetType() int {
	return mgrTypeMannaStarMgr
}

//!插入新自创球星
func (self *MannaStarMgr) AddMannaStar(starInfo *MannaStarType) int {
	insertSql := fmt.Sprintf(`insert into %s (teamid,mannaseat,name,grade,class,icon,face,seat1,seat2,seat3,
	nationality,pass,steals,dribbling,sliding,shooting,goalkeeping,body,speed,passgrow,
	stealsgrow,dribblinggrow,slidinggrow,shootinggrow,goalkeepinggrow,skill1,skill2,skill3,
	skill4,baseprice,basescore,ticket,fate1,fate2,fate3,fate4,fate5,fate6,item,team, hair, eyebrow, mouth, eye, skin, clothes)
	values(%d, %d, '%s', %d, %d, %d, %d, %d, %d, %d, '%s', %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d,
	%d, %d, %d, %d, %d, %d, %d, %d,%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d)`,
		tableMannaStar, starInfo.Teamid, starInfo.MannaSeat, starInfo.Name,
		starInfo.Grade, starInfo.Card, starInfo.Icon, starInfo.Face, starInfo.Seat1, starInfo.Seat2, starInfo.Seat3,
		starInfo.Nationality, starInfo.Pass, starInfo.Steals, starInfo.Dribbling, starInfo.Sliding, starInfo.Shooting,
		starInfo.GoalKeeping, starInfo.Body, starInfo.Speed, starInfo.PassGrow, starInfo.StealsGrow, starInfo.DribblingGrow,
		starInfo.SlidingGrow, starInfo.ShootingGrow, starInfo.GoalKeepingGrow, starInfo.Skill1, starInfo.Skill2, starInfo.Skill3,
		starInfo.Skill4, starInfo.BasePrice, starInfo.BaseScore, starInfo.Ticket, starInfo.Fate1, starInfo.Fate2, starInfo.Fate3,
		starInfo.Fate4, starInfo.Fate5, starInfo.Fate6, starInfo.Item, starInfo.Team, starInfo.Hair, starInfo.Eyebrow,
		starInfo.Mouth, starInfo.Eye, starInfo.Skin, starInfo.Clothes)

	dynamicDBMgr := GetServer().GetDynamicDB()
	insertID, _ := dynamicDBMgr.Exec(insertSql)
	return insertID
}

func (self *MannaStarMgr) GetTemplate(teamid int, seat int) *MannaStarType {
	staticDBMgr := GetServer().GetStaticDataMgr()

	seatType := seat*100 + 10001

	starType := staticDBMgr.GetStarType(seatType)
	star := new(MannaStarType)

	star.Teamid = teamid
	star.Name = fmt.Sprintf("%s", starType.Name)
	star.MannaSeat = seat
	star.Grade = 1
	star.Card = starType.Class
	star.Icon = starType.Icon
	star.Face = 1
	star.Seat1 = starType.Seat1
	star.Seat2 = 0
	star.Seat3 = 0
	star.Nationality = starType.Nationality
	star.Pass = starType.Pass
	star.Steals = starType.Steals
	star.Dribbling = starType.Dribbling
	star.Sliding = starType.Sliding
	star.Shooting = starType.Shooting
	star.GoalKeeping = starType.GoalKeeping
	star.Body = starType.Body
	star.Speed = starType.Speed
	star.PassGrow = starType.PassGrow
	star.StealsGrow = starType.StealsGrow
	star.DribblingGrow = starType.DribblingGrow
	star.SlidingGrow = starType.SlidingGrow
	star.ShootingGrow = starType.ShootingGrow
	star.GoalKeepingGrow = starType.GoalKeepingGrow
	star.Skill1 = starType.Skill1
	star.Skill2 = starType.Skill2
	star.Skill3 = starType.Skill3
	star.Skill4 = starType.Skill4
	star.BasePrice = starType.BasePrice
	star.BaseScore = starType.BaseScore
	star.Ticket = starType.Ticket
	star.Fate1 = starType.Fate1
	star.Fate2 = starType.Fate2
	star.Fate3 = starType.Fate3
	star.Fate4 = starType.Fate4
	star.Fate5 = starType.Fate5
	star.Fate6 = starType.Fate6
	star.Item = starType.Item
	star.Team = starType.Team

	//!默认模型数据
	star.Hair = 1
	star.Eyebrow = 1
	star.Mouth = 1
	star.Eye = 1
	star.Skin = 1
	star.Clothes = 1

	star.Desc = starType.Desc

	star.ID = self.AddMannaStar(star)

	starInfo := new(MannaStar)
	starInfo.MannaStarType = *star
	starInfo.InitDataUpdater(tableMannaStar, &starInfo.MannaStarType)

	self.starList[seat] = starInfo

	return star
}

func NewMannaStar(starInfo *MannaStarType) *MannaStar {
	star := new(MannaStar)
	star.MannaStarType = *starInfo
	star.InitDataUpdater(tableMannaStar, &star.MannaStarType)
	return star
}

//!初始化自创球员管理器
func NewMannaStarMgr(teamID int) IGameMgr {
	mannaStarMgr := new(MannaStarMgr)
	mannaStarMgr.starList = make(MannaStarLst)
	mannaStarQuery := fmt.Sprintf("select * from %s where teamid = %d limit 3", tableMannaStar, teamID)
	mannaStarType := new(MannaStarType)
	mannaStarList := GetServer().GetDynamicDB().fetchAllRows(mannaStarQuery, mannaStarType)
	for v := range mannaStarList {
		mannaStarType = mannaStarList[v].(*MannaStarType)
		mannaStarMgr.starList[mannaStarType.MannaSeat] = NewMannaStar(mannaStarType)
	}

	if len(mannaStarMgr.starList) == 0 {
		//! 若该玩家不存在数据,则创建默认模板
		mannaStarMgr.GetTemplate(teamID, 1)
		mannaStarMgr.GetTemplate(teamID, 2)
		mannaStarMgr.GetTemplate(teamID, 3)
	}

	return mannaStarMgr
}

func (self *MannaStarMgr) UpdateMannaStarSeat() { //! 更新自创球员可踢位置
	staticDataMgr := GetServer().GetStaticDataMgr()
	starList := self.team.GetAllStarList()
	for i := 0; i < starList.Len(); i++ {
		starInfo := self.team.GetStar(starList[i])
		if starInfo.IsMannaStar != 1 {
			continue //!过滤非自创球员
		}

		//! 四星开启场上位置2
		if starInfo.EvolveCount >= 4 {
			mannaStar := self.GetMannaStar(starInfo.Type)
			starTypeID := 10000 + mannaStar.MannaSeat*100 + mannaStar.Seat1
			starTemplate := staticDataMgr.GetStarType(starTypeID)

			if mannaStar.Seat2 == 0 { //!未设置的情况下改为默认位置
				mannaStar.Seat2 = starTemplate.Seat2
			}

			//! 六星开启场上位置3
			if starInfo.EvolveCount >= 6 && mannaStar.Seat3 == 0 {
				mannaStar.Seat3 = starTemplate.Seat3
			}
		}

	}
}
