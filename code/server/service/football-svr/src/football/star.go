package football

import (
	"fmt"
	"math"
	"reflect"
)

const LifeTimeContractPoint = 1000 ///永久球员契约值
const ContractPointMax = 100       ///普通契约值上限

const (
	starAttribPass        = 1 ///传球
	starAttribSteals      = 2 ///抢断
	starAttribDribbling   = 3 ///盘带
	starAttribSliding     = 4 ///铲球
	starAttribShooting    = 5 ///射门
	starAttribGoalKeeping = 6 ///守门
	starAttribBody        = 7 ///身体值
	starAttribSpeed       = 8 ///速度
)

const (
	fieldTypeNone   = 0 ///开始符
	fieldTypeFront  = 1 ///前场
	fieldTypeMiddle = 2 ///中场
	fieldTypeBack   = 3 ///后场
	fieldTypeAll    = 4 ///全场
)

const (
	seatTypeCF  = 1  ///中锋
	seatTypeLWF = 2  ///左边锋
	seatTypeRWF = 3  ///右边锋
	seatTypeSS  = 4  ///影锋
	seatTypeLMF = 5  ///左边卫
	seatTypeRMF = 6  ///右边卫
	seatTypeCMF = 7  ///中前卫
	seatTypeAMF = 8  ///前腰
	seatTypeDMF = 9  ///后腰
	seatTypeCB  = 10 ///中后卫
	seatTypeLB  = 11 ///左边后卫
	seatTypeRB  = 12 ///右边后卫
	seatTypeGK  = 13 ///门将
)

type StarInfoList []StarInfo

type StarTypeStaticData struct { ///服务器配置静态数据表
	ID              int    ///球员类型id
	Name            string ///球员名
	Grade           int    ///品质
	Class           int    ///卡类
	Icon            int    ///头像
	Face            int    ///外观
	Seat1           int    ///踢球位置1
	Seat2           int    ///踢球位置2
	Seat3           int    ///踢球位置3
	Nationality     string ///国藉
	Pass            int    ///传球 IsSkillFull() int { ///判断装备技能栏已满
	Steals          int    ///抢断
	Dribbling       int    ///盘带
	Sliding         int    ///铲球
	Shooting        int    ///射门
	GoalKeeping     int    ///守门
	Body            int    ///身体值
	Speed           int    ///速度
	PassGrow        int    ///传球成长
	StealsGrow      int    ///抢断成长
	DribblingGrow   int    ///盘带成长
	SlidingGrow     int    ///铲球成长
	ShootingGrow    int    ///射门成长
	GoalKeepingGrow int    ///守门成长
	Skill1          int    ///初始技能1
	Skill2          int    ///初始技能2
	Skill3          int    ///初始技能3
	Skill4          int    ///初始技能4
	BasePrice       int    ///基础身价
	BaseScore       int    ///基础评分
	Ticket          int    ///球票价格
	Fate1           int    ///球员缘类型1
	Fate2           int    ///球员缘类型2
	Fate3           int    ///球员缘类型3
	Fate4           int    ///球员缘类型4
	Fate5           int    ///球员缘类型5
	Fate6           int    ///球员缘类型6
	Item            int    ///升星所需道具
	Team            int    ///升星所需道具获取途径
	Desc            string ///球员描述
}

func (self *StarTypeStaticData) GetFirstAndGrowValue() (int, int) { //得到初始值p与成长值
	firstValue := self.Pass + self.Steals +
		self.Dribbling + self.Sliding +
		self.Shooting + self.GoalKeeping +
		self.Body + self.Speed
	growValue := self.PassGrow + self.StealsGrow +
		self.DribblingGrow + self.SlidingGrow +
		self.ShootingGrow + self.GoalKeepingGrow
	return firstValue, growValue
}

const (
	starGradeBegin  = 1 ///开始符
	starGradeGreen  = 1 ///绿色品质球员
	starGradeBlue   = 2 ///蓝色品质球员
	starGradePurple = 3 ///紫色品质球员
	starGradeOrange = 4 ///橙色品质球员
	starGradeRed    = 5 ///红色品质球员
	starGradeEnd        ///结束符
)

const (
	permanentContract = 1000 ///永久契约值表示球员每场比赛均不扣契约值
)

type StarInfo struct { ///球员信息,和dy_star一一对应
	ID                   int `json:"id"`                   ///球员id
	TeamID               int `json:"teamid"`               ///所属球队id
	Type                 int `json:"type"`                 ///球员类型
	Grade                int `json:"grade"`                ///球员品质 1绿2蓝3紫4橙5红
	Level                int `json:"level"`                ///球员等级
	Exp                  int `json:"exp"`                  ///球员经验值
	ContractPoint        int `json:"contractpoint"`        ///合约点数
	EvolveCount          int `json:"evolvecount"`          ///突破次数
	PassTalentAdd        int `json:"passtalentadd"`        ///传球潜力加点
	StealsTalentAdd      int `json:"stealstalentadd"`      ///抢断潜力加点
	DribblingTalentAdd   int `json:"dribblingtalentadd"`   ///盘带潜力加点
	SlidingTalentAdd     int `json:"slidingtalentadd"`     ///铲球潜力加点
	ShootingTalentAdd    int `json:"shootingtalentadd"`    ///射门潜力加点
	GoalKeepingTalentAdd int `json:"goalkeepingtalentadd"` ///守门潜力加点
	TotalPayTalentPoint  int `json:"totalpaytalentpoint"`  ///总消耗培养点,用于洗点返还
	Score                int `json:"score"`                ///战力评分
	IsMannaStar          int `json:"ismannastar"`          //!是否为天赐球员
}

//type IStar interface {
//	//	IObject
//	ISyncObject
//	SetContractPoint(contractPoint int)                                                                                    ///设置球员合约点数
//	SetPermanentContract()                                                                                                 ///设置球员为永久合约
//	GetInfo() *StarInfo                                                                                                    ///得到球员信息对象指针
//	Save()                                                                                                                 ///马上保存数据
//	AddTalentPoint(passAdd int, stealsAdd int, dribblingAdd int, slidingAdd int, shootingAdd int, goalKeepingAdd int) bool ///增加球员培养加成点数
//	AddTotalPayTalentPoint(totalPayTalentPoint int) bool                                                                   ///加球员已消费培养点数
//	CalcScore() float32                                                                                                    ///计算球员评分,根据所踢的位置
//	GetTotalPrice() int                                                                                                    ///得到球员身价
//	GetTotalPrice() int                                                                                                    ///得到球员身价
//	AwardExp(addExp int) int                                                                                               ///奖励球员经验
//	isMaxLevel() bool                  TypeNpcTeamFormWin                                                                                    ///判断球员是否已经到达最高等级上限
//	GetItemFromSort(sortType int) *Item                                                                                    ///通过道具小分类查找球员拥有的装备道具对象
//	GetCalcInfo() *StarInfoCalc
//	IsSkillFull() bool                   ///判断装备技能栏已满                                                                                ///得到二级属性
//	GetTypeInfo() *StarTypeStaticData    ///取得类型静态数据信息
//	IsReachEvolveLimit() bool            ///判断球员是否已到达突破次数限制,受品质颜色影响
//	CanKickFieldType(fieldType int) bool ///判断是否能踢指定的场上位置,前场,中场,后场
//}

type StarInfoCalc struct { ///球员计算后属性
	///基础值由静态数+升级加成+培养加点组成
	PassBase        float32 `json:"passbase"`        ///基础值传球
	StealsBase      float32 `json:"stealsbase"`      ///基础值抢断
	DribblingBase   float32 `json:"dribblingbase"`   ///基础值盘带
	SlidingBase     float32 `json:"slidingbase"`     ///基础值铲球
	ShootingBase    float32 `json:"shootingbase"`    ///基础值射门
	GoalKeepingBase float32 `json:"goalkeepingbase"` ///基础值守门
	BodyBase        float32 `json:"bodybase"`        ///基础值身体值
	SpeedBase       float32 `json:"speedbase"`       ///基础值速度
	///计算后的二级属性
	PassCalc        float32 `json:"passcalc"`        ///计算后传球
	StealsCalc      float32 `json:"stealscalc"`      ///计算后抢断
	DribblingCalc   float32 `json:"dribblingcalc"`   ///计算后盘带
	SlidingCalc     float32 `json:"slidingcalc"`     ///计算后铲球
	ShootingCalc    float32 `json:"shootingcalc"`    ///计算 IsSkillFull() int { ///判断装备技能栏已满后射门
	GoalKeepingCalc float32 `json:"goalkeepingcalc"` ///计算后守门
	BodyCalc        float32 `json:"bodycalc"`        ///计算后身体值
	SpeedCalc       float32 `json:"speedcalc"`       ///计算后速度
	ScoreCalc       float32 `json:"scorecalc"`       ///计算后高精度评分
}

type StarSlice []*Star

type Star struct {
	StarInfo
	StarInfoCalc       ///计算后属性
	DataUpdater        ///信息更新组件
	team         *Team ///所属球队对象
}

func (self *Star) SetStarCount(starCount int) { ///设置球员星级
	self.EvolveCount = starCount
	self.Grade = starCount/2 + (starCount % 2)
}

func (self *Star) GetFieldType(seatType int) int { ///判断是否能踢指定的场上位置,前场,中场,后场
	starFieldType := fieldTypeNone
	if seatType >= seatTypeCF && seatType <= seatTypeSS {
		starFieldType = fieldTypeFront
	}
	if seatType >= seatTypeLMF && seatType <= seatTypeDMF {
		starFieldType = fieldTypeMiddle
	}
	if seatType >= seatTypeCB && seatType <= seatTypeGK {
		starFieldType = fieldTypeBack
	}
	return starFieldType
}

func (self *Star) CanKickFieldType(fieldType int) bool { ///判断是否能踢指定的场上位置,前场,中场,后场
	if self.IsMannaStar == 1 {
		starTypeInfo := self.GetMannaTypeInfo()
		if fieldType == fieldTypeNone {
			return true
		}
		if self.GetFieldType(starTypeInfo.Seat1) == fieldType {
			return true
		}
		if self.GetFieldType(starTypeInfo.Seat2) == fieldType {
			return true
		}
		if self.GetFieldType(starTypeInfo.Seat3) == fieldType {
			return true
		}
		return false
	}

	starTypeInfo := self.GetTypeInfo()
	if fieldType == fieldTypeNone {
		return true
	}
	if self.GetFieldType(starTypeInfo.Seat1) == fieldType {
		return true
	}
	if self.GetFieldType(starTypeInfo.Seat2) == fieldType {
		return true
	}
	if self.GetFieldType(starTypeInfo.Seat3) == fieldType {
		return true
	}
	return false
}

func (self *Star) IsNeedSeatPunish() bool { ///是否需要进行踢球位置惩罚

	if self.IsMannaStar == 1 {
		starTypeInfo := self.team.GetMannaStarMgr().GetMannaStar(self.Type) ///得到球员类型信息
		formation := self.team.GetCurrentFormObject()
		seatType := formation.GetStarSeatType(self.ID) ///尝试在首发阵形中寻找位置
		if seatType <= 0 {
			return false ///未在首发阵形中不进行惩罚
		}
		if seatType != starTypeInfo.Seat1 && seatType != starTypeInfo.Seat2 &&
			seatType != starTypeInfo.Seat3 {
			return true ///未踢习惯的位置需要进行惩罚 IsSkillFull() int { ///判断装备技能栏已满
		}
		return false ///不处罚
	}

	starTypeInfo := GetServer().GetStaticDataMgr().Unsafe().GetStarType(self.Type) ///得到球员类型信息
	formation := self.team.GetCurrentFormObject()
	seatType := formation.GetStarSeatType(self.ID) ///尝试在首发阵形中寻找位置
	if seatType <= 0 {
		return false ///未在首发阵形中不进行惩罚
	}
	if seatType != starTypeInfo.Seat1 && seatType != starTypeInfo.Seat2 &&
		seatType != starTypeInfo.Seat3 {
		return true ///未踢习惯的位置需要进行惩罚 IsSkillFull() int { ///判断装备技能栏已满
	}
	return false ///不处罚
}

func (self *Star) GetSeatType() int { ///计算球员评分,根据所踢的位置
	if self.IsMannaStar == 1 {

		starTypeInfo := self.team.GetMannaStarMgr().GetMannaStar(self.Type) ///得到球员类型信息
		formation := self.team.GetCurrentFormObject()
		seatType := formation.GetStarSeatType(self.ID) ///尝试在首发阵形中寻找位置
		if seatType <= 0 {
			seatType = starTypeInfo.Seat1 ///如果首发阵形中
		}
		return seatType

	}

	starTypeInfo := GetServer().GetStaticDataMgr().Unsafe().GetStarType(self.Type) ///得到球员类型信息
	formation := self.team.GetCurrentFormObject()
	seatType := formation.GetStarSeatType(self.ID) ///尝试在首发阵形中寻找位置
	if seatType <= 0 {
		seatType = starTypeInfo.Seat1 ///如果首发阵形中
	}
	return seatType
}

func (self *Star) GetTotalSkillScore() float32 { ///得到球员总的技能评分 IsSkillFull() int { ///判断装备技能栏已满
	//staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	totalSkillScore := 0
	//skillMgr := self.team.GetSkillMgr()
	//starSkillInfoList := skillMgr.GetStarSkillInfoList(self.ID)
	//for i := range starSkillInfoList {
	//	starSkillInfo := starSkillInfoList[i]
	//skillTypeStaticData := staticDataMgr.GetSkillType(starSkillInfo.Type)
	//totalSkillScore += 1
	//}
	return float32(totalSkillScore)
}

func (self *Star) CalcScore() float32 { ///计算球员评分,根据所踢的位置
	//fmt.Println(self.GetTypeInfo().Name, self.ID, "Star::CalcInfo1", self.StarInfoCalc)
	self.CalcInfo() ///计算球员的二级属性
	seatType := self.GetSeatType()
	calcScore := float32(0.0)
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	seatTypeStaticData := staticDataMgr.GetSeatType(seatType)
	calcScore += self.PassCalc * float32(seatTypeStaticData.ScorePassRate)               ///传球
	calcScore += self.StealsCalc * float32(seatTypeStaticData.ScoreStealsRate)           ///抢断
	calcScore += self.DribblingCalc * float32(seatTypeStaticData.ScoreDribblingRate)     ///盘带
	calcScore += self.SlidingCalc * float32(seatTypeStaticData.ScoreSlidingRate)         ///铲球
	calcScore += self.ShootingCalc * float32(seatTypeStaticData.ScoreShootingRate)       ///射门
	calcScore += self.GoalKeepingCalc * float32(seatTypeStaticData.ScoreGoalKeepingRate) ///守门
	calcScore += self.BodyCalc * float32(seatTypeStaticData.ScoreBodyRate)               ///身体
	calcScore += self.SpeedCalc * float32(seatTypeStaticData.ScoreSpeedRate)             ///速度

	///四舍五入
	// nTemp := int(calcScore) / 10
	// nTemp %= 10
	// if nTemp >= 5 {
	// 	calcScore += 100
	// }

	self.ScoreCalc = calcScore / 100                           ///存入高精度球员评分
	self.ScoreCalc += self.GetTotalSkillScore()                ///加上所有技能评分
	self.Score = int(TakeAccuracy(float32(self.ScoreCalc), 0)) ///存入球员评分

	//fmt.Println(self.GetTypeInfo().Name, self.ID, "Star::CalcInfo2", self.StarInfoCalc)
	return self.ScoreCalc ///返回外界高精度评分
}

func (self *Star) CalcSeatPunish() { ///计算球员踢球位置惩罚
	if self.IsNeedSeatPunish() == false {
		return
	}
	staticDataMgr := GetServer().GetStaticDataMgr()
	///从配置文件中读出球员不在擅长位置踢球所受属性惩罚百分数
	seatPunishValue := staticDataMgr.GetConfigStaticDataInt(configStar, configItemStarCommonConfig, 1)
	seatPunishRate := float32(seatPunishValue) / 100 ///得到惩罚百分小数
	self.PassCalc *= seatPunishRate
	self.StealsCalc *= seatPunishRate
	self.DribblingCalc *= seatPunishRate
	self.SlidingCalc *= seatPunishRate
	self.ShootingCalc *= seatPunishRate
	self.GoalKeepingCalc *= seatPunishRate
	self.BodyCalc *= seatPunishRate
	self.SpeedCalc *= seatPunishRate
}

///球员的缘数量上限由球员品质来决定，如下：绿2、蓝3、紫4、橙5、红6。
func (self *Star) CalcStarFateAddition() { ///计算球员缘系统属性点加成
	//	if self.ID == 11810 {
	//		self.Exp++
	//	}
	PassStarFateAddition := float32(0)
	StealsStarFateAddition := float32(0)
	DribblingStarFateAddition := float32(0)
	SlidingStarFateAddition := float32(0)
	ShootingStarFateAddition := float32(0)
	GoalKeepingStarFateAddition := float32(0)
	BodyStarFateAddition := float32(0)
	SpeedStarFateAddition := float32(0)

	starFateList := IntList{}
	starFateMgr := self.team.GetStarFateMgr()
	if self.IsMannaStar == 1 {
		starTypeInfo := self.GetMannaTypeInfo()
		if starTypeInfo == nil {
			return
		}
		starFateList = IntList{starTypeInfo.Fate1, starTypeInfo.Fate2,
			starTypeInfo.Fate3, starTypeInfo.Fate4, starTypeInfo.Fate5, starTypeInfo.Fate6}
	} else {
		starTypeInfo := self.GetTypeInfo()

		starFateList = IntList{starTypeInfo.Fate1, starTypeInfo.Fate2,
			starTypeInfo.Fate3, starTypeInfo.Fate4, starTypeInfo.Fate5, starTypeInfo.Fate6}
	}

	for i := range starFateList {
		starFate := starFateList[i]
		if 0 == starFate {
			continue ///忽略掉空缘
		}
		//if i > self.Grade+1 {
		//	continue ///忽略掉品质控制之外的缘
		//}
		isStarFateMeetCondition := starFateMgr.IsStarFateMeetCondition(starFate, self.ID)
		if false == isStarFateMeetCondition {
			continue
		}
		starFateType := GetServer().GetStaticDataMgr().Unsafe().GetStarFateType(starFate)
		if equipAddValueNumber == starFateType.AddType { ///数值加成
			PassStarFateAddition += float32(starFateType.Pass)
			StealsStarFateAddition += float32(starFateType.Steals)
			DribblingStarFateAddition += float32(starFateType.Dribbling)
			SlidingStarFateAddition += float32(starFateType.Sliding)
			ShootingStarFateAddition += float32(starFateType.Shooting)
			GoalKeepingStarFateAddition += float32(starFateType.GoalKeeping)
			BodyStarFateAddition += float32(starFateType.Body)
			SpeedStarFateAddition += float32(starFateType.Speed)
		} else if equipAddValueRate == starFateType.AddType { ///百分比加成
			PassStarFateAddition += float32(starFateType.Pass) * self.PassBase / 100
			StealsStarFateAddition += float32(starFateType.Steals) * self.StealsBase / 100
			DribblingStarFateAddition += float32(starFateType.Dribbling) * self.DribblingBase / 100
			SlidingStarFateAddition += float32(starFateType.Sliding) * self.SlidingBase / 100
			ShootingStarFateAddition += float32(starFateType.Shooting) * self.ShootingBase / 100
			GoalKeepingStarFateAddition += float32(starFateType.GoalKeeping) * self.GoalKeepingBase / 100
			BodyStarFateAddition += float32(starFateType.Body) * self.BodyBase / 100
			SpeedStarFateAddition += float32(starFateType.Speed) * self.SpeedBase / 100
		}
	}
	//fmt.Println(self.GetTypeInfo().Name, self.ID, "Star::CalcInfo1", self.StarInfoCalc)
	///将装备加成到二级属性上
	self.PassCalc += PassStarFateAddition
	self.StealsCalc += StealsStarFateAddition
	self.DribblingCalc += DribblingStarFateAddition
	self.SlidingCalc += SlidingStarFateAddition
	self.ShootingCalc += ShootingStarFateAddition
	self.GoalKeepingCalc += GoalKeepingStarFateAddition
	self.BodyCalc += BodyStarFateAddition
	self.SpeedCalc += SpeedStarFateAddition
	//fmt.Println(self.GetTypeInfo().Name, self.ID, "Star::CalcInfo2", self.StarInfoCalc)
}

func (self *Star) CalcEquipAddition() { ///计算球员的装备属性加成
	///开始计算装备加成数
	PassEquipAddition := float32(0)
	StealsEquipAddition := float32(0)
	DribblingEquipAddition := float32(0)
	SlidingEquipAddition := float32(0)
	ShootingEquipAddition := float32(0)
	GoalKeepingEquipAddition := float32(0)
	BodyEquipAddition := float32(0)
	SpeedEquipAddition := float32(0)
	itemMgr := self.team.GetItemMgr()
	equipList := itemMgr.GetStarItemSlice(self.ID)
	for i := range equipList {
		equipItem := equipList[i]
		equipTypeStatic := equipItem.GetTypeInfo()
		equipAddtionFactor := float32(1)                    //equipItem.GetEquipAddtionFactor() ///得到装备加成系数
		if equipAddValueNumber == equipTypeStatic.AddType { ///数值加成
			PassEquipAddition += float32(equipTypeStatic.Pass) * equipAddtionFactor
			StealsEquipAddition += float32(equipTypeStatic.Steals) * equipAddtionFactor
			DribblingEquipAddition += float32(equipTypeStatic.Dribbling) * equipAddtionFactor
			SlidingEquipAddition += float32(equipTypeStatic.Sliding) * equipAddtionFactor
			ShootingEquipAddition += float32(equipTypeStatic.Shooting) * equipAddtionFactor
			GoalKeepingEquipAddition += float32(equipTypeStatic.GoalKeeping) * equipAddtionFactor
			BodyEquipAddition += float32(equipTypeStatic.Body) * equipAddtionFactor
			SpeedEquipAddition += float32(equipTypeStatic.Speed) * equipAddtionFactor
		} else if equipAddValueRate == equipTypeStatic.AddType { ///百分比加成
			PassEquipAddition += float32(equipTypeStatic.Pass) * self.PassBase * equipAddtionFactor / 100
			StealsEquipAddition += float32(equipTypeStatic.Steals) * self.StealsBase * equipAddtionFactor / 100
			DribblingEquipAddition += float32(equipTypeStatic.Dribbling) * self.DribblingBase * equipAddtionFactor / 100
			SlidingEquipAddition += float32(equipTypeStatic.Sliding) * self.SlidingBase * equipAddtionFactor / 100
			ShootingEquipAddition += float32(equipTypeStatic.Shooting) * self.ShootingBase * equipAddtionFactor / 100
			GoalKeepingEquipAddition += float32(equipTypeStatic.GoalKeeping) * self.GoalKeepingBase * equipAddtionFactor / 100
			BodyEquipAddition += float32(equipTypeStatic.Body) * self.BodyBase * equipAddtionFactor / 100
			SpeedEquipAddition += float32(equipTypeStatic.Speed) * self.SpeedBase * equipAddtionFactor / 100
		}
	}
	///将装备加成到二级属性上
	self.PassCalc += PassEquipAddition
	self.StealsCalc += StealsEquipAddition
	self.DribblingCalc += DribblingEquipAddition
	self.SlidingCalc += SlidingEquipAddition
	self.ShootingCalc += ShootingEquipAddition
	self.GoalKeepingCalc += GoalKeepingEquipAddition
	self.BodyCalc += BodyEquipAddition
	self.SpeedCalc += SpeedEquipAddition
}

func (self *Star) CalcTacticEffect(tacticEffectType int) { ///计算阵形战术对球员属性加成
	if tacticEffectType <= 0 {
		return ///不需要加成
	}
	teamInfo := self.team.GetInfo()
	staticDataMgr := GetServer().GetStaticDataMgr()
	///得到阵形战术每等级给玩家球员属性的加成百分比
	tacticAdditionNumer := staticDataMgr.GetConfigStaticDataInt(configStar, configItemStarCommonConfig, 2)
	if tacticAdditionNumer <= 0 {
		return ///战术加成百分比未配置,不加成
	}
	starTypeInfo := staticDataMgr.Unsafe().GetStarType(self.Type)
	if nil == starTypeInfo && self.IsMannaStar == 0 {
		return ///球员类型静态信息不存在,不处理
	}
	//adddTacticEffectValue := float32(tacticAdditionNumer * teamInfo.FormationLevel)

	// 修改:阵型对战术数值加成计算公式需要修改
	//以前是阵型每升一级所有战术加的数值提高5
	//现在是阵型每升一级所有战术加的数值提高值为当前阵型级别的值

	adddTacticEffectValue := float32((1+teamInfo.FormationLevel)*teamInfo.FormationLevel/2 + 5)

	switch tacticEffectType {
	case starAttribPass: ///传球
		self.PassCalc += adddTacticEffectValue //float32(starTypeInfo.Pass * (tacticAdditionRate + 100) * teamInfo.FormationLevel / 100)
	case starAttribSteals: ///抢断
		self.StealsCalc += adddTacticEffectValue // float32(starTypeInfo.Steals * (tacticAdditionRate + 100) * teamInfo.FormationLevel / 100)
	case starAttribDribbling: ///盘带
		self.DribblingCalc += adddTacticEffectValue //float32(starTypeInfo.Dribbling * (tacticAdditionRate + 100) * teamInfo.FormationLevel / 100)
	case starAttribSliding: ///铲球
		self.SlidingCalc += adddTacticEffectValue //float32(starTypeInfo.Sliding * (tacticAdditionRate + 100) * teamInfo.FormationLevel / 100)
	case starAttribShooting: ///射门
		self.ShootingCalc += adddTacticEffectValue //float32(starTypeInfo.Shooting * (tacticAdditionRate + 100) * teamInfo.FormationLevel / 100)
	}
}

func (self *Star) CalcTacticAddition() { ///计算阵形战术对球员属性加成
	staticDataMgr := GetServer().GetStaticDataMgr()
	formation := self.team.GetCurrentFormObject() ///得到球队当前阵形对象
	starPos := formation.GetStarPos(self.ID)
	if starPos <= 0 {
		return ///非首发球员不计算阵形战术加成
	}
	formationInfo := formation.GetInfo()
	if formationInfo.CurrentTactic <= 0 {
		return ///当前阵形未设置战术时不计算阵形战术加成
	}
	tacticTypStaicData := staticDataMgr.Unsafe().GetTacticType(formationInfo.CurrentTactic) ///得到阵型类型信息
	if nil == tacticTypStaicData {
		return ///查到此战术类型信息
	}
	effectPosList := IntList{tacticTypStaicData.Pos1, tacticTypStaicData.Pos2, tacticTypStaicData.Pos3,
		tacticTypStaicData.Pos4, tacticTypStaicData.Pos5, tacticTypStaicData.Pos6, tacticTypStaicData.Pos7,
		tacticTypStaicData.Pos8, tacticTypStaicData.Pos9, tacticTypStaicData.Pos10, tacticTypStaicData.Pos11}
	for i := range effectPosList {
		effectPos := effectPosList[i]
		if effectPos <= 0 {
			continue
		}
		if i == starPos-1 {
			self.CalcTacticEffect(tacticTypStaicData.Attrib1) ///进行属性1加成
			self.CalcTacticEffect(tacticTypStaicData.Attrib2) ///进行属性2加成
			break
		}
		//if i ==  && effectPos > 0 {
		//	self.CalcTacticEffect(tacticTypStaicData.Attrib1) ///进行属性1加成
		//	self.CalcTacticEffect(tacticTypStaicData.Attrib2) ///进行属性2加成
		//	break
		//}
	}
}

func (self *Star) ResetInfoToStatic() { ///将球员的属性重置为静态基础属性
	loger := GetServer().GetLoger()
	if self.IsMannaStar == 1 {
		mannaStarMgr := self.team.GetMannaStarMgr()
		starTypeInfo := mannaStarMgr.GetMannaStar(self.Type)
		if nil == starTypeInfo {
			loger.Warn("starTypeInfo == nil starType = %d starID = %d", self.Type, self.ID)
			return
		}

		///将基础属性恢复成静态基础属性
		self.PassBase = float32(starTypeInfo.Pass)
		self.StealsBase = float32(starTypeInfo.Steals)
		self.DribblingBase = float32(starTypeInfo.Dribbling)
		self.SlidingBase = float32(starTypeInfo.Sliding)
		self.ShootingBase = float32(starTypeInfo.Shooting)
		self.GoalKeepingBase = float32(starTypeInfo.GoalKeeping)
		self.BodyBase = float32(starTypeInfo.Body)
		self.SpeedBase = float32(starTypeInfo.Speed)
		///将二级属性重置为基础信息
		self.PassCalc = self.PassBase
		self.StealsCalc = self.StealsBase
		self.DribblingCalc = self.DribblingBase
		self.SlidingCalc = self.SlidingBase
		self.ShootingCalc = self.ShootingBase
		self.GoalKeepingCalc = self.GoalKeepingBase
		self.BodyCalc = self.BodyBase
		self.SpeedCalc = self.SpeedBase

		return
	}

	staticDataMgr := GetServer().GetStaticDataMgr()
	starTypeInfo := staticDataMgr.Unsafe().GetStarType(self.Type)
	if nil == starTypeInfo {
		loger.Warn("starTypeInfo == nil starType = %d starID = %d", self.Type, self.ID)
		return
	}

	///将基础属性恢复成静态基础属性
	self.PassBase = float32(starTypeInfo.Pass)
	self.StealsBase = float32(starTypeInfo.Steals)
	self.DribblingBase = float32(starTypeInfo.Dribbling)
	self.SlidingBase = float32(starTypeInfo.Sliding)
	self.ShootingBase = float32(starTypeInfo.Shooting)
	self.GoalKeepingBase = float32(starTypeInfo.GoalKeeping)
	self.BodyBase = float32(starTypeInfo.Body)
	self.SpeedBase = float32(starTypeInfo.Speed)
	///将二级属性重置为基础信息
	self.PassCalc = self.PassBase
	self.StealsCalc = self.StealsBase
	self.DribblingCalc = self.DribblingBase
	self.SlidingCalc = self.SlidingBase
	self.ShootingCalc = self.ShootingBase
	self.GoalKeepingCalc = self.GoalKeepingBase
	self.BodyCalc = self.BodyBase
	self.SpeedCalc = self.SpeedBase
}

func (self *Star) CalcLevelAddition() { ///计算球员的属性等级加成
	if self.Level <= 1 {
		return ///等级1球员不计算等级加成 IsSkillFull() int { ///判断装备技能栏已满
	}
	///开始计算等级加成
	//staticDataMgr := GetServer().GetStaticDataMgr()
	//starTypeInfo := staticDataMgr.Unsafe().GetStarType(self.Type)
	//starGradeEffectSlice := staticDataMgr.getConfigStaticDataParamIntList(configStar, configStarGradeEffect)
	//starGradeEffect := float32(0)
	//if starTypeInfo.Grade > 0 && starTypeInfo.Grade <= len(starGradeEffectSlice) {
	//	starGradeEffect = float32(starGradeEffectSlice[starTypeInfo.Grade-1]) ///得到品质加成系数
	//}

	//等级成长加成公式修改:等级成长=初始值*（1+（等级-1）（0.03+星级*0.02））
	//starLevelGradeEffect := 1.0 + (float32(self.Level)-1.0)*(0.03+float32(self.EvolveCount)*0.02)

	//等级成长加成公式修改:等级成长=初始值*（1+（等级-1）（0.04*星级））
	starLevelGradeEffect := 1.0 + (float32(self.Level)-1.0)*(float32(self.EvolveCount)*0.04)
	self.PassCalc = self.PassBase * starLevelGradeEffect
	self.StealsCalc = self.StealsBase * starLevelGradeEffect
	self.DribblingCalc = self.DribblingBase * starLevelGradeEffect
	self.SlidingCalc = self.SlidingBase * starLevelGradeEffect
	self.ShootingCalc = self.ShootingBase * starLevelGradeEffect
	self.GoalKeepingCalc = self.GoalKeepingBase * starLevelGradeEffect
	self.BodyCalc = self.BodyBase * starLevelGradeEffect
	self.SpeedCalc = self.SpeedBase * starLevelGradeEffect
}

func (self *Star) CalcEducationAddition() { ///计算球员培养属性点加成
	///开始计算等级加成
	//self.PassBase += float32(self.PassTalentAdd)
	//self.StealsBase += float32(self.StealsTalentAdd)
	//self.DribblingBase += float32(self.DribblingTalentAdd)
	//self.SlidingBase += float32(self.SlidingTalentAdd)
	//self.ShootingBase += float32(self.ShootingTalentAdd)
	//self.GoalKeepingBase += float32(self.GoalKeepingTalentAdd)
	///将二级属性重置为基础信息
	//self.PassCalc = self.PassBase
	//self.StealsCalc = self.StealsBase
	//self.DribblingCalc = self.DribblingBase
	//self.SlidingCalc = self.SlidingBase
	//self.ShootingCalc = self.ShootingBase
	//self.GoalKeepingCalc = self.GoalKeepingBase
	//self.BodyCalc = self.BodyBase
	//self.SpeedCalc = self.SpeedBase

	if self.IsMannaStar == 1 {
		//培养点数不会变
		starType := self.team.GetMannaStarMgr().GetMannaStar(self.Type)

		//!获取培养上限
		passTalentAdd := self.GetAddPointMax(float32(self.EvolveCount), float32(starType.PassGrow))
		self.PassTalentAdd = int(Min_float(float32(self.PassTalentAdd), passTalentAdd))

		stealsTalentAdd := self.GetAddPointMax(float32(self.EvolveCount), float32(starType.StealsGrow))
		self.StealsTalentAdd = int(Min_float(float32(self.StealsTalentAdd), stealsTalentAdd))

		dribblingTalentAdd := self.GetAddPointMax(float32(self.EvolveCount), float32(starType.DribblingGrow))
		self.DribblingTalentAdd = int(Min_float(float32(self.DribblingTalentAdd), dribblingTalentAdd))

		slidingTalentAdd := self.GetAddPointMax(float32(self.EvolveCount), float32(starType.SlidingGrow))
		self.SlidingTalentAdd = int(Min_float(float32(self.SlidingTalentAdd), slidingTalentAdd))

		shootingTalentAdd := self.GetAddPointMax(float32(self.EvolveCount), float32(starType.ShootingGrow))
		self.ShootingTalentAdd = int(Min_float(float32(self.ShootingTalentAdd), shootingTalentAdd))

		goalKeepingTalentAdd := self.GetAddPointMax(float32(self.EvolveCount), float32(starType.GoalKeepingGrow))
		self.GoalKeepingTalentAdd = int(Min_float(float32(self.GoalKeepingTalentAdd), goalKeepingTalentAdd))

	}

	self.PassCalc += float32(self.PassTalentAdd)
	self.StealsCalc += float32(self.StealsTalentAdd)
	self.DribblingCalc += float32(self.DribblingTalentAdd)
	self.SlidingCalc += float32(self.SlidingTalentAdd)
	self.ShootingCalc += float32(self.ShootingTalentAdd)
	self.GoalKeepingCalc += float32(self.GoalKeepingTalentAdd)
}

func (self *Star) CalcBaseInfo() { ///计算球员的基础属性
	self.ResetInfoToStatic()     ///将球员的属性重置为静态基础属性
	self.CalcLevelAddition()     ///计算球员等级属性加成
	self.CalcEducationAddition() ///计算球员培养属性点加成

	self.PassBase = self.PassCalc
	self.StealsBase = self.StealsCalc
	self.DribblingBase = self.DribblingCalc
	self.SlidingBase = self.SlidingCalc
	self.ShootingBase = self.ShootingCalc
	self.GoalKeepingBase = self.GoalKeepingCalc
	self.BodyBase = self.BodyCalc
	self.SpeedBase = self.SpeedCalc
}

func (self *Star) CalcInfo() { ///计算球员的属性
	//if self.ID == 2537 {
	//	self.ResetInfoToStatic() ///将球员的属性重置为静态基础属性
	//}
	self.CalcBaseInfo()         ///计算球员的基础属性
	self.CalcStarFateAddition() ///计算球员缘属性加成
	self.CalcEquipAddition()    ///计算球员的装备属性加成
	self.CalcTacticAddition()   ///计算阵形战术对球员属性加成
	self.CalcSeatPunish()       ///最后计算球员踢球位置惩罚

	if self.IsMannaStar == 1 {
		fmt.Println(self.GetMannaTypeInfo().Name, self.ID, "Star::CalcInfo", self.StarInfoCalc)
		return
	}
	fmt.Println(self.GetTypeInfo().Name, self.ID, "Star::CalcInfo", self.StarInfoCalc)
}

func (self *Star) GetReflectValue() reflect.Value { ///得到球队反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

func NewStar(starInfo *StarInfo, team *Team) *Star {
	star := new(Star)
	star.StarInfo = *starInfo
	star.InitDataUpdater(tableStar, &star.StarInfo)
	star.team = team
	return star
}

func (self *Star) GetCalcInfo() *StarInfoCalc { ///得到二级属性
	return &self.StarInfoCalc
}

func (self *Star) GetTotalPrice() int { ///得到球员身价
	//	starTypeStaticData := GetServer().GetStaticDataMgr().Unsafe().GetStarType(self.Type) ////得到球员类型数据
	//totalPrice := self.Score * starTypeStaticData.Grade
	//公式改变:球员身价=球员评分*星级/100
	starLevel := self.EvolveCount
	totalPrice := self.Score * starLevel / 100
	return totalPrice
}

func (self *Star) IsSkillFull() bool { ///判断装备技能栏已满
	const MAXSKILLCOUNT = 3 ///一个球员最多可装备3道具
	isSkillFull := self.GetSkillCount() >= MAXSKILLCOUNT
	return isSkillFull
}

func (self *Star) GetSkillCount() int { ///得到球员已装备技能数
	skillMgr := self.team.GetSkillMgr()
	starSkillSlice := skillMgr.GetStarSkillSlice(self.ID)
	skillCount := len(starSkillSlice)
	return skillCount
}

func (self *Star) GetID() int { ///得到ID
	return self.ID
}

func (self *Star) isMaxLevel() bool { ///判断球员是否已经到达最高等级上限
	staticDataMgr := GetServer().GetStaticDataMgr()
	levelExpCount := staticDataMgr.GetLevelExpCount(levelExpTypeStarLevel)
	return self.Level >= levelExpCount
}

func (self *Star) Uplevel() { ///球员升级
	///得到当前等级升级所需经验
	staticDataMgr := GetServer().GetStaticDataMgr()
	levelExpCount := staticDataMgr.GetLevelExpCount(levelExpTypeStarLevel)
	oldLevel := self.Level ///保存旧等级
	for i := 1; i <= levelExpCount; i++ {
		needExp := staticDataMgr.GetLevelExpNeedExp(levelExpTypeStarLevel, self.Level) ///得到当前等级经验

		nextLevel := self.Level + 1
		if nextLevel > 100 {
			nextLevel = 100
		}

		needEvolveCount := staticDataMgr.GetLevelExpNeedEvolveCount(levelExpTypeStarLevel, nextLevel) ///得到下一等级需求星级
		if self.Exp < needExp {
			break ///经验不足升级
		}

		if self.EvolveCount < needEvolveCount {
			break ///没有达到星级需求
		}
		if self.isMaxLevel() == true {
			self.Exp = needExp ///满级后经验到达上限
			break              ///已满级
		}
		self.Level++
	}
	//	teamLevel := self.team.GetLevel()
	//	self.Level = Min(teamLevel, self.Level+1) ///限制球员最大等级不得超过球队等级
	if self.Level > oldLevel {
		self.CalcScore() ///升级后需要重新计算评分
	}
}

///通过道具小分类查找球员拥有的装备道具对象
func (self *Star) GetItemFromSort(sortType int) *Item {
	itemMgr := self.team.GetItemMgr()
	starItemSlice := itemMgr.GetStarItemSlice(self.ID)
	for i := range starItemSlice {
		starItem := starItemSlice[i]
		itemTypeInfo := starItem.GetTypeInfo()
		if itemTypeInfo.Sort == sortType {
			return starItem
		}
	}
	return nil
}

func (self *Star) AwardExp(addExp int) int { ///奖励球员经验
	if addExp <= 0 {
		return 0
	}
	if self.isMaxLevel() == true { ///判断球员是否已经到达最高等级上限
		return self.Exp ///满级后不在获得经验
	}
	self.Exp += addExp
	self.Uplevel() ///球员升级
	return self.Exp
}

func (self *Star) GetAddPointMax(evolveCount float32, growValue float32) float32 { ///根据公式得到当前培养上限
	param1 := (evolveCount * (evolveCount + 1.0) / 2.0)
	param2 := (growValue*growValue/90.0 + 0.5*growValue - 30.0)

	//	param2 = TakeAccuracy(param2, 3) // 取三位小数精度

	addPointMax := float64(param1 * param2)
	addPointMax = math.Ceil(addPointMax)
	if addPointMax < 0 {
		addPointMax = 0
	}
	return float32(addPointMax)
}

func (self *Star) ChangeStarAttribute(starCount int) { ///设置球员属性

	if self.IsMannaStar == 1 {
		starType := self.GetMannaTypeInfo()
		self.PassTalentAdd = int(self.GetAddPointMax(float32(starCount), float32(starType.PassGrow)))
		self.StealsTalentAdd = int(self.GetAddPointMax(float32(starCount), float32(starType.StealsGrow)))
		self.DribblingTalentAdd = int(self.GetAddPointMax(float32(starCount), float32(starType.DribblingGrow)))
		self.SlidingTalentAdd = int(self.GetAddPointMax(float32(starCount), float32(starType.SlidingGrow)))
		self.ShootingTalentAdd = int(self.GetAddPointMax(float32(starCount), float32(starType.ShootingGrow)))
		self.GoalKeepingTalentAdd = int(self.GetAddPointMax(float32(starCount), float32(starType.GoalKeepingGrow)))

		fmt.Println("mannastar.PassTalentAdd", self.PassTalentAdd)

		self.CalcScore()
		return
	}

	starType := self.GetTypeInfo()

	self.PassTalentAdd = int(self.GetAddPointMax(float32(starCount), float32(starType.PassGrow)))
	self.StealsTalentAdd = int(self.GetAddPointMax(float32(starCount), float32(starType.StealsGrow)))
	self.DribblingTalentAdd = int(self.GetAddPointMax(float32(starCount), float32(starType.DribblingGrow)))
	self.SlidingTalentAdd = int(self.GetAddPointMax(float32(starCount), float32(starType.SlidingGrow)))
	self.ShootingTalentAdd = int(self.GetAddPointMax(float32(starCount), float32(starType.ShootingGrow)))
	self.GoalKeepingTalentAdd = int(self.GetAddPointMax(float32(starCount), float32(starType.GoalKeepingGrow)))
	self.CalcScore()
}

func (self *Star) SetPermanentContract() { ///设置球员为永久合约
	self.ContractPoint = permanentContract
}

func (self *Star) GetInfo() *StarInfo { ///得到球员信息对象指针
	return &self.StarInfo
}

func (self *Star) IsReachEvolveLimit() bool { ///判断球员是否已到达突破次数限制,受品质颜色影响
	result := false
	staticDataMgr := GetServer().GetStaticDataMgr()
	starEvolveLevelLimitList := staticDataMgr.getConfigStaticDataParamIntList(configStar, configStarEvolveLevelLimit)

	if self.IsMannaStar == 1 {
		starType := self.GetMannaTypeInfo()
		starEvolveLevelLimitCount := starEvolveLevelLimitList[starType.Grade-1] ///索引到品质对应突破限制次数
		result = self.EvolveCount >= starEvolveLevelLimitCount
		return result
	}
	starType := self.GetTypeInfo()
	starEvolveLevelLimitCount := starEvolveLevelLimitList[starType.Grade-1] ///索引到品质对应突破限制次数
	result = self.EvolveCount >= starEvolveLevelLimitCount
	return result
}

func (self *Star) SetContractPoint(contractPoint int) { ///设置球员合约点数
	if contractPoint < 0 {
		return ///要求为正数或零
	}
	self.ContractPoint = contractPoint
}

func (self *Star) AddTotalPayTalentPoint(totalPayTalentPoint int) bool { ///加球员已消费培养点数
	if totalPayTalentPoint <= 0 {
		return false ///要求为正数或零
	}
	self.TotalPayTalentPoint += totalPayTalentPoint
	return true
}

func (self *Star) AddTalentPoint(passAdd int, stealsAdd int, dribblingAdd int,
	slidingAdd int, shootingAdd int, goalKeepingAdd int) bool { ///增加球员培养加成点数
	if passAdd < 0 || stealsAdd < 0 || dribblingAdd < 0 || slidingAdd < 0 || shootingAdd < 0 || goalKeepingAdd < 0 {
		return false ///不允许为负数
	}
	self.PassTalentAdd += passAdd
	self.StealsTalentAdd += stealsAdd
	self.DribblingTalentAdd += dribblingAdd
	self.SlidingTalentAdd += slidingAdd
	self.ShootingTalentAdd += shootingAdd
	self.GoalKeepingTalentAdd += goalKeepingAdd
	return true
}

func (self *Star) GetTypeInfo() *StarTypeStaticData { ///取得类型静态数据信息
	StarTypeStaticData := GetServer().GetStaticDataMgr().Unsafe().GetStarType(self.Type)
	return StarTypeStaticData
}

func (self *Star) GetMannaTypeInfo() *MannaStarType {
	mannaStar := self.team.GetMannaStarMgr()
	if mannaStar.GetMannaStar(self.Type) != nil {
		return &mannaStar.GetMannaStar(self.Type).MannaStarType
	}

	return nil
}

const (
	ssPotentialCode = 1 ///ss级潜力代号
	sPotentialode   = 2 ///s级潜力代号
	aPotentialCode  = 3 ///a级潜力代号
	bPotentialCode  = 4 ///b级潜力代号
	cPotentialCode  = 5 ///c级潜力代号
	dPotentialCode  = 6 ///d级潜力代号

	ssEvolveCountMax = 9 ///ss级球员升星上限
	sEvolveCountMax  = 8 ///s级球员升星上限
	aEvolveCountMax  = 7 ///a级球员升星上限
	bEvolveCountMax  = 7 ///b级球员升星上限
	cEvolveCountMax  = 7 ///c级球员升星上限
	dEvolveCountMax  = 7 ///d级球员升星上限

	ssEvolveCommonItem = 510001 ///ss级万能碎片ID
	sEvolveCommonItem  = 510002 ///s级万能碎片ID
	aEvolveCommonItem  = 510003 ///a级万能碎片ID
	bEvolveCommonItem  = 510004 ///b级万能碎片ID
	cEvolveCommonItem  = 510005 ///c级万能碎片ID
	dEvolveCommonItem  = 510006 ///d级万能碎片ID
)

func GetPotentialCode(grow int) int { ///通过静态球员模板的class球员评分数据取得潜力评级
	if grow >= 550 {
		return ssPotentialCode
	} else if grow >= 500 {
		return ssPotentialCode
	} else if grow >= 450 {
		return sPotentialode
	} else if grow >= 400 {
		return aPotentialCode
	} else if grow >= 350 {
		return bPotentialCode
	} else if grow >= 300 {
		return cPotentialCode
	} else if grow >= 250 {
		return dPotentialCode
	} else if grow >= 200 {
		return dPotentialCode
	} else {
		return dPotentialCode
	}
}

func getEvolveCommonItem(potentialCode int) int { ///取得潜力评级对应的万能碎片id
	res := 0
	switch potentialCode {
	case ssPotentialCode:
		res = ssEvolveCommonItem
	case sPotentialode:
		res = sEvolveCommonItem
	case aPotentialCode:
		res = aEvolveCommonItem
	case bPotentialCode:
		res = bEvolveCommonItem
	case cPotentialCode:
		res = cEvolveCommonItem
	case dPotentialCode:
		res = dEvolveCommonItem
	}
	return res
}

func (self *Star) GetEvolveCountMax() int { ///取得球星的升星上限
	res := 0
	if self.IsMannaStar == 1 {
		starTypeStaticData := self.GetMannaTypeInfo()
		switch GetPotentialCode(starTypeStaticData.Card) {
		case ssPotentialCode:
			res = ssEvolveCountMax
		case sPotentialode:
			res = sEvolveCountMax
		case aPotentialCode:
			res = aEvolveCountMax
		case bPotentialCode:
			res = bEvolveCountMax
		case cPotentialCode:
			res = cEvolveCountMax
		case dPotentialCode:
			res = dEvolveCountMax
		}
		return res
	}
	starTypeStaticData := self.GetTypeInfo()
	switch GetPotentialCode(starTypeStaticData.Class) {
	case ssPotentialCode:
		res = ssEvolveCountMax
	case sPotentialode:
		res = sEvolveCountMax
	case aPotentialCode:
		res = aEvolveCountMax
	case bPotentialCode:
		res = bEvolveCountMax
	case cPotentialCode:
		res = cEvolveCountMax
	case dPotentialCode:
		res = dEvolveCountMax
	}
	return res
}

func (self *Star) GetEvolveNeedItem() (needItem int, insteadItem int) { ///取得球星升星需求道具和替代道具

	if self.IsMannaStar == 1 {

		starTypeStaticData := self.GetMannaTypeInfo()              ///取得对应的球星模板数据
		potentialCode := GetPotentialCode(starTypeStaticData.Card) ///取得球星潜力代号 天赐球员只能用固定碎片升星
		if starTypeStaticData.Item != 0 {
			needItem = starTypeStaticData.Item
			//		insteadItem = getEvolveCommonItem(potentialCode)
		} else {
			needItem = 0
			insteadItem = 0
		}
		if potentialCode > aPotentialCode { ///把b级以后的替代道具设置为0
			insteadItem = 0
		}
		return
	}

	starTypeStaticData := self.GetTypeInfo()                    ///取得对应的球星模板数据
	potentialCode := GetPotentialCode(starTypeStaticData.Class) ///取得球星潜力代号
	if starTypeStaticData.Item != 0 {
		needItem = starTypeStaticData.Item
		insteadItem = getEvolveCommonItem(potentialCode)
	} else {
		needItem = getEvolveCommonItem(potentialCode)
		insteadItem = 0
	}
	if potentialCode > aPotentialCode { ///把b级以后的替代道具设置为0
		insteadItem = 0
	}
	return
}
