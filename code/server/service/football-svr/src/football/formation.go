package football

import (
	//	"encoding/json"
	"reflect"
)

const (
	formationMinStarCount = 11
)

type FormationInfo struct { ///阵型信息与数据库dy_formation表一一对应
	ID            int `json:"id"`            ///阵形id
	TeamID        int `json:"teamid"`        ///所属球队id
	Type          int `json:"type"`          ///阵型类型
	Pos1          int `json:"pos1"`          ///位置1
	Pos2          int `json:"pos2"`          ///位置2
	Pos3          int `json:"pos3"`          ///位置3
	Pos4          int `json:"pos4"`          ///位置4
	Pos5          int `json:"pos5"`          ///位置5
	Pos6          int `json:"pos6"`          ///位置6
	Pos7          int `json:"pos7"`          ///位置7
	Pos8          int `json:"pos8"`          ///位置8
	Pos9          int `json:"pos9"`          ///位置9
	Pos10         int `json:"pos10"`         ///位置10
	Pos11         int `json:"pos11"`         ///位置11
	CurrentTactic int `json:"currenttactic"` ///当前战术类型
	Tactic1       int `json:"tactic1"`       ///已开启战术1
	Tactic2       int `json:"tactic2"`       ///已开启战术2
	Tactic3       int `json:"tactic3"`       ///已开启战术3
}

type FormationTypeStaticData struct { ///球队阵型类型表
	ID           int    ///记录id
	Name         string ///阵型类型表
	Seat1        int    ///球员位1
	Seat2        int    ///球员位2
	Seat3        int    ///球员位3
	Seat4        int    ///球员位4
	Seat5        int    ///球员位4
	Seat6        int    ///球员位6
	Seat7        int    ///球员位7
	Seat8        int    ///球员位8
	Seat9        int    ///球员位9
	Seat10       int    ///球员位10
	Seat11       int    ///球员位11
	Pos1         int    ///位置1
	Pos2         int    ///位置2
	Pos3         int    ///位置3
	Pos4         int    ///位置4
	Pos5         int    ///位置5
	Pos6         int    ///位置6
	Pos7         int    ///位置7
	Pos8         int    ///位置8
	Pos9         int    ///位置9
	Pos10        int    ///位置10
	Pos11        int    ///位置11
	OvercomeForm int    ///剋制阵型类型
}

//type IFormation interface {
//	ISyncObject
//	Save()                               ///马上保存数据
//	GetInfo() *FormationInfo             ///得到球队反射对象
//	HasTactic(tacticType int) bool       ///判断此阵形是否已开启指定战术
//	IsOvercome(formationType int) bool   ///是否克制此阵形类型
//	IsBeOvercome(formationType int) bool ///是否被此阵形类型克制
//	GetStarIDList() IntList              ///得到阵形球员id列表
//	GetStarSeatType(starID int) int      ///得到球员在此阵形中的踢球位置
//	GetStarPos(starID int) int           ///返回球员在此阵形中的位置索引值,从1开始
//}

type FormationPtr *Formation
type Formation struct {
	FormationInfo
	DataUpdater ///信息保存组件
}

func (self *Formation) GetID() int {
	return self.ID
}

func (self *Formation) HasTactic(tacticType int) bool { ///判断此阵形是否已开启指定战术
	if tacticType == self.Tactic1 || tacticType == self.Tactic2 || tacticType == self.Tactic3 {
		return true
	}
	return false
}

func (self *Formation) GetInfo() *FormationInfo { ///得到球队反射对象
	return &self.FormationInfo
}

func (self *Formation) GetReflectValue() reflect.Value { ///得到球队反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

func (self *Formation) IsBeOverTactic(tacticType int) bool { ///是否克制此战术类型
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	currentTactic := self.CurrentTactic % 100
	tacticTypeStaticData := staticDataMgr.GetTacticType(currentTactic)
	if nil == tacticTypeStaticData || 0 == tacticTypeStaticData.Over || 0 == tacticType {
		return false ///任意一方战术为0则不产生克生关系
	}
	return tacticTypeStaticData.Over == tacticType
}

func (self *Formation) IsOverTactic(tacticType int) bool { ///是否克制此战术类型
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	tacticTypeStaticData := staticDataMgr.GetTacticType(tacticType)
	currentTactic := self.CurrentTactic % 100
	if nil == tacticTypeStaticData || 0 == tacticTypeStaticData.Over || 0 == currentTactic {
		return false ///任意一方战术为0则不产生克生关系
	}
	return tacticTypeStaticData.Over == currentTactic
}

func (self *Formation) IsOvercome(formationType int) bool { ///是否克制此阵形类型
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	formationTypeStaticData := staticDataMgr.GetFormationType(formationType)
	return formationTypeStaticData.OvercomeForm == self.Type
}

func (self *Formation) IsBeOvercome(formationType int) bool { ///是否被此阵形类型克制
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	formationTypeStaticData := staticDataMgr.GetFormationType(self.Type)
	return formationTypeStaticData.OvercomeForm == formationType
}

func NewFormation(formationInfo *FormationInfo) *Formation { ///新建阵型对象
	formation := new(Formation)
	formation.FormationInfo = *formationInfo ///放入球星属性信息
	formation.InitDataUpdater(tableFormation, &formation.FormationInfo)
	return formation
}

///得到球员在此阵型踢球位置,不在此阵形返回0
func (self *Formation) GetStarSeatType(starID int) int {
	seatType := 0
	starIDList := self.GetStarIDList()
	seatTypeList := self.GetSeatTypeList() ///得到阵形位置列表
	for i := range starIDList {
		storeStarID := starIDList[i]
		if storeStarID == starID {
			seatType = seatTypeList[i]
			break
		}
	}
	return seatType
}

func (self *Formation) GetStarPos(starID int) int { ///返回球员在此阵形中的位置索引值,从1开始
	starPos := 0
	starIDList := self.GetStarIDList()
	for i := range starIDList {
		storeStarID := starIDList[i]
		if storeStarID == starID {
			starPos = i + 1
			break
		}
	}
	return starPos
}

func (self *Formation) GetSeatTypeList() IntList { ///得到阵形位置列表
	formType := GetServer().GetStaticDataMgr().Unsafe().GetFormationType(self.Type)
	seatTypeList := IntList{formType.Pos1, formType.Pos2, formType.Pos3, formType.Pos4, formType.Pos5,
		formType.Pos6, formType.Pos7, formType.Pos8, formType.Pos9, formType.Pos10, formType.Pos11}
	return seatTypeList
}

func (self *Formation) SetStarPos(indexPos int, starID int) { ///设置球员在阵形中的位置
	starPosPtrList := IntPtrList{&self.Pos1, &self.Pos2, &self.Pos3, &self.Pos4, &self.Pos5, &self.Pos6,
		&self.Pos7, &self.Pos8, &self.Pos9, &self.Pos10, &self.Pos11}
	listLen := starPosPtrList.Len()
	if indexPos < 0 || indexPos >= listLen {
		return ///out of range
	}
	starPosPtr := starPosPtrList[indexPos]
	*starPosPtr = starID ///在指定位置放入新球员
}

func (self *Formation) GetStarIDList() IntList { ///得到阵形球员id列表
	starIDList := IntList{self.Pos1, self.Pos2, self.Pos3, self.Pos4, self.Pos5, self.Pos6,
		self.Pos7, self.Pos8, self.Pos9, self.Pos10, self.Pos11}
	return starIDList
}

func (self *Formation) IsTacticFull() bool { ///判断阵形是否已开启所有阵形
	isTacticFull := false
	tacticStoreList := IntList{self.Tactic1, self.Tactic2, self.Tactic3}
	tacticStoreListLen := tacticStoreList.Len()
	tacticCount := 0
	for i := range tacticStoreList {
		tacticType := tacticStoreList[i]
		if tacticType > 0 {
			tacticCount++
		}
	}
	isTacticFull = tacticCount >= tacticStoreListLen
	return isTacticFull
}

func (self *Formation) UpdateFormTactic(formLevel int) { ///根据当前阵形等级,升级阵形战术
	if self.IsTacticFull() == true {
		return ///所有战术已开启不用再升级了
	}
	staticDataMgr := GetServer().GetStaticDataMgr()
	tacticTypeList := staticDataMgr.GetTacticTypeList(self.Type, formLevel)
	tacticStoreList := IntList{self.Tactic1, self.Tactic2, self.Tactic3}
	for i := range tacticStoreList {
		tacticType := tacticStoreList[i]
		if tacticType > 0 {
			tacticTypeList = append(tacticTypeList, tacticType)
		}
	}
	tacticType := 0
	tacticTypeList = tacticTypeList.Unique()
	tacticTypeListLen := tacticTypeList.Len()
	tacticTypeStorePtrList := IntPtrList{&self.Tactic1, &self.Tactic2, &self.Tactic3}
	for i := range tacticTypeStorePtrList {
		tacticTypeStorePtr := tacticTypeStorePtrList[i]
		tacticType = 0
		if i < tacticTypeListLen {
			tacticType = tacticTypeList[i]
		}
		*tacticTypeStorePtr = tacticType
	}
}
