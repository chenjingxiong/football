package football

import (
	"fmt"
	"reflect"
)

type AtlasInfo struct {
	ID       int `json:"id"`
	TeamID   int `json:"teamid"`
	StarType int `json:"type"`
	Received int `json:"received"`
}

type Atlas struct {
	AtlasInfo
	DataUpdater
}

type AtlasLst map[int]*Atlas

type AtlasMgr struct {
	GameMgr
	atlasLst AtlasLst //! 球星图鉴
}

func (self *Atlas) GetReflectValue() reflect.Value { ///得到反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

func (self *AtlasMgr) GetType() int {
	return mgrTypeAtlasMgr
}

func (self *AtlasMgr) SaveInfo() {
	for _, v := range self.atlasLst {
		v.Save() //! 玩家下线保存
	}
}

func (self *AtlasMgr) AddAtlas(teamID int, starType int, received int) { //! 添加一条新的图鉴
	if self.HasStar(starType) == true { //! 曾经得到过就不再录入数据库
		return
	}

	info := new(AtlasInfo)
	info.TeamID = teamID
	info.StarType = starType
	info.Received = received

	atlasInser := fmt.Sprintf("insert into %s (teamid,startype,received) values(%d, %d, %d)", tableAtlas, teamID, starType, received)
	insertID, _ := GetServer().GetDynamicDB().Exec(atlasInser)

	if insertID > 0 {

		info.ID = insertID

		self.atlasLst[starType] = NewAtlas(info)
	}
}

func (self *AtlasMgr) HasStar(starType int) bool { //! 根据球星类型判断图鉴存在
	return self.atlasLst[starType] != nil
}

func (self *AtlasMgr) GetAtlas(starType int) *AtlasInfo { //! 根据球星类型得到图鉴信息
	return &self.atlasLst[starType].AtlasInfo
}

func (self *AtlasMgr) GetTypeAtlas(classType int) []AtlasInfo { //! 根据类型获取图鉴信息
	staticDataMgr := GetServer().GetStaticDataMgr()
	var atlasLst []AtlasInfo
	minClass := 0
	maxClass := 0
	switch classType {
	case 1: //! 所有A级
		maxClass = 450
		minClass = 400
	case 2: //! 所有S级
		maxClass = 500
		minClass = 450
	case 3: //! 所有SS级
		maxClass = 600
		minClass = 500
	}

	for _, v := range self.atlasLst {
		info := v.AtlasInfo
		starType := staticDataMgr.GetStarType(info.StarType)
		if starType.Class >= minClass && starType.Class < maxClass {
			atlasLst = append(atlasLst, info)
		}
	}

	return atlasLst
}

func (self *AtlasMgr) GetAllAtlas() []AtlasInfo {
	var atlasLst []AtlasInfo
	for _, v := range self.atlasLst {
		info := v.AtlasInfo
		atlasLst = append(atlasLst, info)
	}

	return atlasLst
}

//! 初始化球员图鉴管理器
func NewAtlasMgr(teamID int) IGameMgr {
	atlasMgr := new(AtlasMgr)
	atlasMgr.atlasLst = make(AtlasLst)
	atlasQuery := fmt.Sprintf("select * from %s where teamid = %d", tableAtlas, teamID)

	info := new(AtlasInfo)
	atlasLst := GetServer().GetDynamicDB().fetchAllRows(atlasQuery, info)
	for v := range atlasLst {
		info = atlasLst[v].(*AtlasInfo)
		atlasMgr.atlasLst[info.StarType] = NewAtlas(info)
	}

	return atlasMgr
}

func NewAtlas(info *AtlasInfo) *Atlas {
	atlas := new(Atlas)
	atlas.AtlasInfo = *info
	atlas.InitDataUpdater(tableAtlas, &atlas.AtlasInfo)
	return atlas
}
