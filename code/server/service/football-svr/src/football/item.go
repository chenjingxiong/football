package football

import (
	"fmt"
	"reflect"
)

const (
	ItemActionAwardMaxCount = 10 ///活动奖励道具最大种类数
)

type ActionType struct {
	ID           int ///行为id
	NeedType1    int ///条件类型1
	NeedCount1   int ///条件数量
	NeedType2    int ///条件类型1
	NeedCount2   int ///条件类型1
	NeedType3    int ///条件类型1
	NeedCount3   int ///条件类型1
	NeedType4    int ///条件类型1
	NeedCount4   int ///条件类型1
	AwardType1   int ///奖励道具类型 正数为奖 负数为扣
	AwardCount1  int ///奖励数量
	AwardGrade1  int ///奖励品质
	AwardStar1   int ///奖励球星类型
	Probability1 int ///概率(万分比)
	AwardType2   int ///条件类型1
	AwardCount2  int ///条件数量
	AwardGrade2  int ///条件类型1
	AwardStar2   int ///条件类型1
	Probability2 int ///概率(万分比)
	AwardType3   int ///条件类型1
	AwardCount3  int ///条件数量
	AwardGrade3  int ///条件类型1
	AwardStar3   int ///条件类型1
	Probability3 int ///概率(万分比)
	AwardType4   int ///条件类型1
	AwardCount4  int ///条件数量
	AwardGrade4  int ///条件类型1
	AwardStar4   int ///条件类型1
	Probability4 int ///概率(万分比)
	AwardType5   int ///条件类型1
	AwardCount5  int ///条件数量
	AwardGrade5  int ///条件类型1
	AwardStar5   int ///条件类型1
	Probability5 int ///概率(万分比)
	AwardType6   int ///条件类型1
	AwardCount6  int ///条件数量
	AwardGrade6  int ///条件类型1
	AwardStar6   int ///条件类型1
	Probability6 int ///概率(万分比)
	AwardType7   int ///条件类型1
	AwardCount7  int ///条件数量
	AwardGrade7  int ///条件类型1
	AwardStar7   int ///条件类型1
	Probability7 int ///概率(万分比)
	AwardType8   int ///条件类型1
	AwardCount8  int ///条件数量
	AwardGrade8  int ///条件类型1
	AwardStar8   int ///条件类型1
	Probability8 int ///概率(万分比)
	AwardType9   int ///条件类型1
	AwardCount9  int ///条件数量
	AwardGrade9  int ///条件类型1
	AwardStar9   int ///条件类型1
	Probability9 int ///概率(万分比)
}

/////道具对象
//const (
//	itemTypeNormal = 1 ///1为衣服 2为个性物品 3为战靴 4
//)

const (
	equipAddValueNumber = "number"
	equipAddValueRate   = "rate"
)

const ( ///道具所在位置
	itemPosTeamEquip = 1 ///球队经理装备栏
	itemPosClothes   = 2 ///球员球衣装备栏
	itemPosShoe      = 3 ///球员球鞋装备栏
	itemPosAccessory = 4 ///球员附件装备栏
	itemPosStore     = 5 ///球队经理仓库
)

const ItemEvolveLevel = 5 ///每达到一个5级可进行升阶

const (
	itemColorBegin  = 1 ///开始符
	itemColorWhite  = 1 ///白色品质道具
	itemColorGreen  = 2 ///绿色品质道具
	itemColorBlue   = 3 ///蓝色品质道具
	itemColorPurple = 4 ///紫色品质道具
	itemColorOrange = 5 ///橙色品质道具
	//	itemColorRed    = 6 ///红色品质道具
	itemColorEnd ///结束符
)

///1为衣服 2为个性物品 3为战靴
const (
	ItemSortCloth = 1 ///球衣
	ItemSortJewel = 2 ///饰品
	ItemSortShoe  = 3 ///球鞋
)

const (
	ItemTypeEquip = 1 ///球员装备
	ItemTypeStuff = 2 ///材料
	ItemTypeProp  = 3 ///可使用道具
	ItemTypeDummy = 4 ///数值类道具,并不会产生对象
)

type ItemInfoList []ItemInfo

type ItemTypeStaticData struct {
	ID            int    ///道具类型id
	Name          string ///道具名
	Icon          int    ///道具图标id
	Type          int    ///道具类型 1 装备 2 材料 3 道具 4数值类
	Sort          int    ///小分类 装备(1为衣服 2为个性物品 3为战靴)
	Color         int    ///道具品质 1=绿色 2=蓝色 3=紫色 4=橙色
	Overlay       int    ///道具叠加数量上限
	Level         int    ///道具使用等级
	BuyTicket     int    ///具买入球票价格
	SellCoin      int    ///道具出售球票价格
	BuyMask       int    ///购买限制
	UseAction     int    ///使用行为触发,对应action表中一条记录
	AddType       string ///属性加成类型 number 绝对值 percent百分比
	Pass          int    ///传球加成
	Steals        int    ///抢夺加成
	Dribbling     int    ///盘带加成
	Sliding       int    ///铲球加成MergeLevel
	Shooting      int    ///射门加成
	GoalKeeping   int    ///守门加成
	Body          int    ///身体值加成
	Speed         int    ///速度加成
	Merge         int    ///融合值,用于道具合成
	Mod           int    ///装备样式
	Skill         int    ///绑定技能类型,0表示无技能
	Attackscore   int    ///攻击加成
	Defensescore  int    ///防守加成const (
	Organizescore int    ///组织加成
	Desc          string ///道具描述
}

func (self *ItemTypeStaticData) IsNumberType() bool { ///判断是否是数值类道具
	result := (self.Type == ItemTypeDummy)
	return result
}

type ItemInfo struct { ///游戏中的道具信息
	ID         int `json:"id"`         ///道具id
	TeamID     int `json:"teamid"`     ///拥有球队id
	StarID     int `json:"starid"`     ///拥有者球星id
	Type       int `json:"type"`       ///道具类型
	Color      int `json:"color"`      ///道具品质 1=白色 2=绿色 3=蓝色 4=紫色 5=橙色
	Count      int `json:"count"`      ///道具叠加数
	Position   int `json:"postion"`    ///道具存放的位置 1球队背包 2球员装备栏
	Cell       int `json:"cell"`       ///所在位置中的单元格子索引 1球员装备球衣 2球员装备不过鞋 3球员装备饰品
	MergeExp   int `json:"mergeexp"`   ///当前累计融合值
	MergeLevel int `json:"mergelevel"` ///融合等级
}

//type IGameItem interface {
//	GetID() int                       ///得到道具ID
//	GetInfo() *ItemInfo               ///得到信息对象指针
//	SetStarID(starID int)             ///设置道具球星主人id
//	Sync() DataValueChangList         ///得到属性变更列表
//	GetTotalMergeExp() int            ///得到道具总融合经验值
//	AwardMergeExp(mergeExp int) bool  ///加道具融合经验值,需要处理升级
//	Save()                            ///马上保存数据
//	SetCell(cell int)                 ///设置道具所在单元格索引号
//	GetTypeInfo() *ItemTypeStaticData ///取得道具类型静态数据信息
//	GetReflectValue() reflect.Value   ///得到反射对象
//	GetEquipAddtionFactor() float32   ///得到装备加成系数
//}

type Item struct { ///游戏中的道具
	//	Object
	ItemInfo
	DataUpdater
}

func (self *Item) GetReflectValue() reflect.Value { ///得到反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

func GetActionType(actionType int) *ActionType {
	staticDataMgr := GetServer().GetStaticDataMgr()
	element := staticDataMgr.GetStaticData(tableActionType, actionType)
	if nil == element {
		return nil
	}
	return element.(*ActionType)
}

func (self *Item) GetUseAction() *ActionType {
	typeInfo := self.GetTypeInfo()
	actionType := GetActionType(typeInfo.UseAction)
	return actionType
}

///得到action奖励列表,奖励道具类型,奖励数量,奖励品质,奖励球星
func (self *Item) GetActionAwardList() (IntList, IntList, IntList, IntList, IntList) {
	awardTypeList, awardCountList, awardGradeList, awardStarList, awardProbabilityList := IntList{}, IntList{}, IntList{}, IntList{}, IntList{}
	actionType := self.GetUseAction()
	value := reflect.ValueOf(actionType).Elem()
	awardTypeFieldName, awardCountFieldName, awardGradeFieldName, awardStarFieldName, awardProbabilityName := "", "", "", "", ""
	for i := 1; i < ItemActionAwardMaxCount; i++ {
		awardTypeFieldName = fmt.Sprintf("AwardType%d", i)
		awardTypeField := value.FieldByName(awardTypeFieldName)
		if awardTypeField.IsValid() == false {
			break ///没有此字段直接跳出
		}
		awardCountFieldName = fmt.Sprintf("AwardCount%d", i)
		awardCountField := value.FieldByName(awardCountFieldName)
		if awardCountField.IsValid() == false {
			break ///没有此字段直接跳出
		}
		awardGradeFieldName = fmt.Sprintf("AwardGrade%d", i)
		awardGradeField := value.FieldByName(awardGradeFieldName)
		if awardGradeField.IsValid() == false {
			break ///没有此字段直接跳出
		}
		awardStarFieldName = fmt.Sprintf("AwardStar%d", i)
		awardStarField := value.FieldByName(awardStarFieldName)
		if awardStarField.IsValid() == false {
			break ///没有此字段直接跳出
		}
		awardProbabilityName = fmt.Sprintf("Probability%d", i)
		awardProbabilityField := value.FieldByName(awardProbabilityName)
		if awardStarField.IsValid() == false {
			break ///没有此字段直接跳出
		}

		awardType := int(awardTypeField.Int())
		awardCount := int(awardCountField.Int())
		awardGrade := int(awardGradeField.Int())
		awardStar := int(awardStarField.Int())
		awardProbability := int(awardProbabilityField.Int())
		if awardType <= 0 && awardStar <= 0 {
			continue
		}
		awardTypeList = append(awardTypeList, awardType)
		awardCountList = append(awardCountList, awardCount)
		awardGradeList = append(awardGradeList, awardGrade)
		awardStarList = append(awardStarList, awardStar)
		awardProbabilityList = append(awardProbabilityList, awardProbability)
	}
	return awardTypeList, awardCountList, awardGradeList, awardStarList, awardProbabilityList
}

func (self *Item) GetInfo() *ItemInfo {
	return &self.ItemInfo
}

func (self *Item) GetTotalMergeExp() int {
	totalergeExp := GetServer().GetStaticDataMgr().GetItemTypeMerge(self.Type) ///得到基础融合经验
	totalergeExp += self.MergeExp                                              ////累加当前融合经验
	return totalergeExp
}

func (self *Item) GetEquipAddtionFactor() float32 { ///得到装备加成系数
	result := self.Color
	return float32(result)

	//	paramColorDic := []float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6} ///颜色品质对应的加成系数
	//	paramTotalDic := []float32{0, 0.5, 1.5, 3, 5, 7.5}       ///累计加成系
	staticDataMgr := GetServer().GetStaticDataMgr()
	paramColorParamList := staticDataMgr.getConfigStaticDataParamIntList(configItem, configItemColorAdditionParam)
	paramSumColorParamList := staticDataMgr.getConfigStaticDataParamIntList(configItem, configItemSumColorAdditionParam)
	currentColorLevel := float32(self.MergeLevel) ///通过融合等级计算当前颜色等级
	colorIndex := self.Color - 1                  ///颜色品质索引
	paramColorParam := float32(paramColorParamList[colorIndex]) / 100
	paramSumColorParam := float32(paramSumColorParamList[colorIndex]) / 100
	equipAddtionFactor := 1 + (currentColorLevel * paramColorParam) + paramSumColorParam
	return equipAddtionFactor
}

func (self *Item) Uplevel() { ///道具升级
	///得到当前等级升级所需经验
	staticDataMgr := GetServer().GetStaticDataMgr()
	itemType := self.GetTypeInfo()
	itemExpType := levelExpTypeEquipMerge + itemType.Sort - 1
	levelExpCount := staticDataMgr.GetLevelExpCount(itemExpType)
	colorMaxLevel := self.Color * ItemEvolveLevel ///当前品质颜色升融合等级上限
	for i := 1; i < levelExpCount; i++ {
		needExp := staticDataMgr.GetLevelExpNeedExp(itemExpType, self.MergeLevel+1)
		if self.MergeExp < needExp {
			break ///经验不足升级
		}
		self.MergeLevel++
		if self.MergeLevel >= levelExpCount || self.MergeLevel >= colorMaxLevel {
			self.MergeExp = 0 ///满级后经验到达上限
			break             ///已满级
		}
	}
}

func (self *Item) AwardMergeExp(mergeExp int) bool { ///加道具融合经验值,需要处理升级
	if mergeExp <= 0 {
		return false
	}
	self.MergeExp += mergeExp
	self.Uplevel() ///道具升级
	return true
}

func (self *Item) SetMergeExp(mergeExp int) { ///设置道具融合经验值
	self.MergeExp = mergeExp
}

func (self *Item) SetCell(cell int) { ///设置道具所在单元格索引号
	self.Cell = cell
}

func (self *Item) GetTypeInfo() *ItemTypeStaticData { ///取得道具类型静态数据信息
	itemTypeStaticData := GetServer().GetStaticDataMgr().Unsafe().GetItemType(self.Type)
	return itemTypeStaticData
}

func (self *Item) GetID() int { ///得到道具ID
	return self.ID
}

func (self *Item) SetStarID(starID int) { ///设置道具球星主人id
	self.StarID = starID
}

func NewItem(itemInfo *ItemInfo) *Item {
	item := new(Item)
	item.ItemInfo = *itemInfo
	item.InitDataUpdater(tableItem, &item.ItemInfo)
	return item
}

func GetItemType(itemType int) *ItemTypeStaticData {
	staticDataMgr := GetServer().GetStaticDataMgr()
	element := staticDataMgr.GetStaticData(tableItemType, itemType)
	if nil == element {
		return nil
	}
	return element.(*ItemTypeStaticData)
}
