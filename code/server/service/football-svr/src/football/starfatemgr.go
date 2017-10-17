package football

const (
	fateTypeBothPlay         = 1 ///同时指定球员同时上场生效,多个球员类型与关系
	fateTypeFormationSeatPos = 2 ///指定阵型与位置生效,可以指定多个位置,多个位置或关系
	fateTypeWearEquip        = 3 ///球员穿戴指定装备生效,可以指定多个装备,多个道具与关系
	fateTypeTactic           = 4 ///在指定战术中生效,多个战术是或关系
)

type StarFateTypeStaticData struct { ///球员缘静态配置数据
	ID          int    ///缘类型id
	Name        string ///缘名称
	Type        int    ///缘类型
	Param       int    ///类型相关参数
	Pos1        int    ///对象1
	Pos2        int    ///对象2
	Pos3        int    ///对象3
	Pos4        int    ///对象4
	Pos5        int    ///对象5
	AddType     string ///属性加成类型 number 绝对值 percent百分比
	Pass        int    ///传球加成
	Steals      int    ///抢夺加成
	Dribbling   int    ///盘带加成
	Sliding     int    ///铲球加成
	Shooting    int    ///射门加成
	GoalKeeping int    ///守门加成
	Body        int    ///身体值加成
	Speed       int    ///速度加成
	Desc        string ///道具描述
}

///球员缘系统管理器
type IStarFateMgr interface {
	//CalcStarFateAddition(star IStar, starFateType int) ///计算球员缘系统效果
	IsStarFateMeetCondition(starFateType int, starID int) bool ///判断此缘条件是否满足
}

type StarFateMgr struct {
	GameMgr
}

func (self *StarFateMgr) GetType() int { ///得到管理器类型
	return mgrTypeStarFateMgr
}

func NewStarFateMgr(teamID int) IGameMgr {
	starFateMgr := new(StarFateMgr)
	return starFateMgr
}

///判断同时上场是否满足
func (self *StarFateMgr) IsBothPlayMeetCondition(starFatePosList IntList) bool {
	for i := range starFatePosList {
		starFatePos := starFatePosList[i]
		if 0 == starFatePos {
			continue
		}

		if starFatePos > 10000 { //! 天赐球员缘分特殊处理
			starSeat := starFatePos % 10
			starMannaSeat := (starFatePos - 10000 - starFatePos%10) / 100
			mannaStarMgr := self.team.GetMannaStarMgr()
			mannaStar := mannaStarMgr.GetMannaStarFromSeat(starMannaSeat)

			if mannaStar == nil { //! 判断是否拥有该球员
				return false
			}

			if mannaStar.Seat1 != starSeat { //! 判断球员是否踢指定位置
				return false
			}

			continue
		}

		if self.team.HasStar(starFatePos) == false {
			return false
		}
	}
	return true
}

///判断同时阵型位置是否满足
func (self *StarFateMgr) IsFormationSeatPosMeetCondition(starID int, starFatePosList IntList) bool {
	currentForm := self.team.GetCurrentFormObject()
	starSeatType := currentForm.GetStarSeatType(starID)
	for i := range starFatePosList {
		starFatePos := starFatePosList[i]
		if 0 == starFatePos {
			continue
		}
		if starFatePos == starSeatType {
			return true
		}
	}
	return false
}

///判断穿戴装备是否满足
func (self *StarFateMgr) IsWearEquipMeetCondition(starID int, starFatePosList IntList) bool {
	starItemTypeList := self.team.GetItemMgr().GetStarItemTypeList(starID)
	for i := range starFatePosList {
		starFatePos := starFatePosList[i]
		if 0 == starFatePos {
			continue
		}
		if starItemTypeList.Search(starFatePos) <= 0 {
			return false
		}
	}
	return true
}

///判断当前战术是否满足
func (self *StarFateMgr) IsTacticMeetCondition(starFatePosList IntList) bool {
	currentFormation := self.team.GetCurrentFormObject()
	currentFormInfo := currentFormation.GetInfo()
	if currentFormInfo.CurrentTactic <= 0 {
		return false
	}
	for i := range starFatePosList {
		starFatePos := starFatePosList[i]
		if 0 == starFatePos {
			continue
		}
		if currentFormInfo.CurrentTactic == starFatePos {
			return true
		}
	}
	return false
}

func (self *StarFateMgr) IsStarFateMeetCondition(starFateType int, starID int) bool { ///判断此缘条件是否满足
	isStarFateMeetCondition := false
	fateType := GetServer().GetStaticDataMgr().Unsafe().GetStarFateType(starFateType)
	starFatePosList := IntList{fateType.Pos1, fateType.Pos2, fateType.Pos3, fateType.Pos4, fateType.Pos5}
	currentFormInfo := self.team.GetCurrentFormObject().GetInfo()
	switch fateType.Type {
	case fateTypeBothPlay: ///同时指定球员同时上场生效,多个球员类型与关系
		isStarFateMeetCondition = self.IsBothPlayMeetCondition(starFatePosList)
	case fateTypeFormationSeatPos: ///指定阵型与位置生效,可以指定多个位置,多个位置或关系、
		if currentFormInfo.Type != fateType.Param {
			return false ///当前阵形不匹配
		}
		isStarFateMeetCondition = self.IsFormationSeatPosMeetCondition(starID, starFatePosList)
	case fateTypeWearEquip: ///球员穿戴指定装备生效,可以指定多个装备,多个道具与关系
		isStarFateMeetCondition = self.IsWearEquipMeetCondition(starID, starFatePosList)
	case fateTypeTactic: ///在指定战术中生效,多个战术是或关系
		isStarFateMeetCondition = self.IsTacticMeetCondition(starFatePosList)
	}
	return isStarFateMeetCondition
}

func (self *StarFateMgr) SaveInfo() { ///保存数据

}

//func (self *StarFateMgr) CalcStarFateAddition(star IStar, starFateType int) { ///计算球员缘系统效果
//	passFateAddition := float32(0)
//	stealsFateAddition := float32(0)
//	dribblingFateAddition := float32(0)
//	slidingFateAddition := float32(0)
//	shootingFateAddition := float32(0)
//	goalKeepingFateAddition := float32(0)
//	bodyFateAddition := float32(0)
//	speedFateAddition := float32(0)
//	team := self.team
//	fateType := GetServer().GetStaticDataMgr().Unsafe().GetStarFateType(starFateType)
//	starFatePosList := IntList{fateType.Pos1, fateType.Pos2, fateType.Pos3, fateType.Pos4, fateType.Pos5}
//	for i:=range starFatePosList{
//		starFatePos:=starFatePosList[i]
//		if 0==starFatePos{
//			continue
//		}
//		switch starFateType {
//		case fateTypeBothPlay: ///同时指定球员同时上场生效,多个球员类型与关系
//			if false==team.IsStarInCurrentFormation(starFatePos){
//				fateType=nil///未满足条件
//				break
//			}
//		case fateTypeFormationSeatPos: ///指定阵型与位置生效,可以指定多个位置,多个位置或关系
//		case fateTypeWearEquip: ///球员穿戴指定装备生效,可以指定多个装备,多个道具与关系
//		case fateTypeTactic: ///在指定战术中生效,多个战术是或关系
//		}
//	}
//	if fateType!=nil{

//	}
//}
