package football

import (
	"fmt"
	"reflect"
)

//type IFormationMgr interface { ///阵型管理器
//	GetFormationList() FormationArray                         ///设置球队当前阵形
//	GetFormation(formationID int) IFormation                  ///得到球队阵形对象
//	IsStarInFormation(formationID int, starID int) bool       ///判断球员是否在阵形中
//	GetFormationInfoList() FormationInfoList                  ///得到阵形信息列表
//	AwardFormation(formationType int) int                     ///添加阵形信息,成功返回阵形对象id
//	IsStarTypeInFormation(formationID int, starType int) bool ///判断球员类型是否在阵形中
//}

type FormationList map[int]*Formation
type FormationArray []*Formation       ///阵形对象数组
type FormationInfoList []FormationInfo ///阵形信息数组

type FormationMgr struct { ///阵型管理器
	GameMgr
	formationList FormationList ///阵型列表
}

func (self *FormationMgr) GetFormation(formationID int) *Formation { ///得到球队阵形对象
	return self.formationList[formationID]
}

func (self *FormationMgr) GetType() int { ///得到管理器类型
	return mgrTypeFormationMgr
}

func (self *FormationMgr) IsStarTypeInFormation(formationID int, starType int) bool { ///判断球员类型是否在阵形中
	currentFormation := self.formationList[formationID]
	if nil == currentFormation {
		return false
	}
	formationValue := reflect.ValueOf(currentFormation).Elem()
	for i := 1; i <= 11; i++ {
		fieldName := fmt.Sprintf("Pos%d", i)
		value := formationValue.FieldByName(fieldName)
		if value.IsValid() == false {
			return false
		}
		starID := int(value.Int())
		star := self.team.GetStar(starID)
		if nil == star {
			continue
		}
		starInfo := star.GetInfo()
		if starInfo.Type == starType {
			return true
		}
	}
	return false
}

func (self *FormationMgr) IsStarInFormation(formationID int, starID int) bool { ///判断球员是否在阵形中
	currentFormation := self.formationList[formationID]
	if nil == currentFormation {
		return false
	}
	formationValue := reflect.ValueOf(currentFormation).Elem()
	for i := 1; i <= 11; i++ {
		fieldName := fmt.Sprintf("Pos%d", i)
		value := formationValue.FieldByName(fieldName)
		if value.IsValid() == false {
			return false
		}
		if value.Int() == int64(starID) {
			return true
		}
	}
	return false
}

func (self *FormationMgr) GetFormationInfoList() FormationInfoList { ///得到阵形信息列表
	formationInfoList := FormationInfoList{}
	for _, v := range self.formationList {
		formationInfo := v.GetInfo()
		formationInfoList = append(formationInfoList, *formationInfo)
	}
	return formationInfoList
}

func (self *FormationMgr) GetFormationList() FormationArray { ///设置球队当前阵形
	formationArray := FormationArray{}
	for _, v := range self.formationList {
		formationArray = append(formationArray, v)
	}
	return formationArray
}

func (self *FormationMgr) HasFormation(formationID int) bool { ///判断球队是否拥有指定阵形
	_, ok := self.formationList[formationID]
	return ok
}

func (self *FormationMgr) loadFormationList() bool { ///加载球队所属型阵列表
	if self.formationList != nil {
		return false
	}
	self.formationList = make(FormationList)
	formationInfo := new(FormationInfo)
	formationListQuery := fmt.Sprintf("select * from %s where teamid=%d limit 8", tableFormation, self.team.GetID())
	elmentList := GetServer().GetDynamicDB().fetchAllRows(formationListQuery, formationInfo)
	if nil == elmentList {
		return false
	}
	for i := range elmentList {
		formationInfo = elmentList[i].(*FormationInfo)
		//newFormation := new(Formation)
		//newFormation.Create(formationInfo)
		self.formationList[formationInfo.ID] = NewFormation(formationInfo)
	}
	numFormation := len(self.formationList)
	return numFormation > 0
}

func NewFormationMgr(teamID int) IGameMgr {
	formationMgr := new(FormationMgr)
	formationMgr.formationList = make(FormationList)
	//formationMgr.teamID = teamID ///存放自己的球队id
	formationListQuery := fmt.Sprintf("select * from %s where teamid=%d limit 3000", tableFormation, teamID)
	formationInfo := new(FormationInfo)
	formationInfoList := GetServer().GetDynamicDB().fetchAllRows(formationListQuery, formationInfo)
	for i := range formationInfoList {
		formationInfo = formationInfoList[i].(*FormationInfo)
		formationMgr.formationList[formationInfo.ID] = NewFormation(formationInfo)
	}
	return formationMgr
}

func (self *FormationMgr) CorrectAllFormError() { ///纠正所有阵形中潜在的错误
	for k, _ := range self.formationList {
		self.CorrectFormError(k)
	}
}

func (self *FormationMgr) CorrectFormError(formID int) { ///纠正阵形中潜在的错误
	dstForm := self.GetFormation(formID)
	curForm := self.team.GetCurrentFormObject()  ///得到首发阵形对象
	curFormStarIDList := curForm.GetStarIDList() ///得到首发阵形球员列表
	dstFormStarIDList := dstForm.GetStarIDList()
	for i := range dstFormStarIDList {
		starID := dstFormStarIDList[i]
		star := self.team.GetStar(starID)
		if star != nil {
			continue ///有效的球员直接忽略
		}
		///此球员已经被解雇了star对象为nil
		starIDInCurForm := curFormStarIDList[i]       ///取得当前阵形中指定索引中的球员	id
		starNew := self.team.GetStar(starIDInCurForm) ///再次验证
		if nil == starNew {
			continue ///无效的球员直接忽略
		}
		dstForm.SetStarPos(i, starIDInCurForm) ///更新一个正确的球员到阵形中
	}
}

func (self *FormationMgr) SaveInfo() { ///保存数据
	for _, v := range self.formationList {
		v.Save()
	}
}

func (self *FormationMgr) AwardFormation(formationType int) int { ///添加阵形信息,成功返回阵形对象id
	currentFormationID := self.team.GetCurrentFormation()
	currentFormation := self.GetFormation(currentFormationID)
	starIDList := currentFormation.GetStarIDList() ///得到首发球员id列表
	awardFormationQuery := fmt.Sprintf("insert %s set teamid=%d,type=%d", tableFormation, self.team.GetID(), formationType)
	for i := range starIDList {
		awardFormationQuery += fmt.Sprintf(",pos%d=%d", i+1, starIDList[i])
	}
	lastInsertFormationID, _ := GetServer().GetDynamicDB().Exec(awardFormationQuery)
	if lastInsertFormationID <= 0 {
		GetServer().GetLoger().Warn("FormationMgr AwardFormation fail! formationType:%d teamID:%d",
			formationType, self.team.GetID())
		return 0
	}
	formationInfo := new(FormationInfo)
	formationInfo.ID = lastInsertFormationID
	formationInfo.TeamID = self.team.GetID()
	formationInfo.Type = formationType
	formationInfo.Pos1 = starIDList[0]
	formationInfo.Pos2 = starIDList[1]
	formationInfo.Pos3 = starIDList[2]
	formationInfo.Pos4 = starIDList[3]
	formationInfo.Pos5 = starIDList[4]
	formationInfo.Pos6 = starIDList[5]
	formationInfo.Pos7 = starIDList[6]
	formationInfo.Pos8 = starIDList[7]
	formationInfo.Pos9 = starIDList[8]
	formationInfo.Pos10 = starIDList[9]
	formationInfo.Pos11 = starIDList[10]
	formation := NewFormation(formationInfo)
	self.formationList[lastInsertFormationID] = formation ///生成阵型对象
	formation.UpdateFormTactic(self.team.FormationLevel)  ///获得阵型后马上更新战术
	return lastInsertFormationID
}
