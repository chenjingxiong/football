package football

import (
	"reflect"
)

/////道具对象
//const (
//	itemTypeNormal = 1 ///1为衣服 2为个性物品 3为战靴 4
//)

/////1为衣服 2为个性物品 3为战靴
//const (
//	ItemSortCloth = 1 ///球衣
//	ItemSortJewel = 2 ///饰品
//	ItemSortShoe  = 3 ///球鞋
//)

// Type
//   技能的触发类型
//   1=根据比分状况触发，次类型下
//     Sort=1  为比分领先下触发
//     Sort=2  为比分持平下触发
//     Sort=3  为比分落后下触发
//   2=根据我方阵型触发，此类型下sort字段定义如下，同formationtype表阵型定义
//      1	442平行
// 2	433
// 3	4321
// 4	352
// 5	442菱形
// 6	343菱形
// 7	4231
// 8	343平行
//  3=根据敌方阵型触发，sort字段定义同上。
//  4=被动永久生效类型，此类型下sort字段无效。

// Sort
//   触发效果的分类，详细见type字段说明

// Tartype
//   触发技能后生效的对象类型
//   1=本方球员
//   2=敌方球员
// Tarsort
//   触发技能后生效的具体对象
//     编号	球员位置	备注
// 0	技能释放者	tartype为1时才有效
// 1	CF	全部中锋
// 2	LWF	左边锋
// 3	RWF	右边锋
// 4	SS	全部影锋
// 5	LMF	左前卫
// 6	RMF	右前卫
// 7	CMF	全部中前卫
// 8	AMF	全部前腰
// 9	DMF	全部后腰
// 10	CB	全部中后卫
// 11	LB	左后卫
// 12	RB	右后卫
// 13	GK	门将
// 14	1.2.3.4	全部中锋、边锋、影锋
// 15	5.6.7.8.9.	全部边前卫、前腰、中前卫、后腰
// 16	10.11.12.13	全部边后卫、中后卫、门将
// 17	全部球员
// Func
//    技能产生的功能效果
//    1=增加属性
//    2=降低属性
//    3=减少克制效果
// 4以后为特殊类效果

// Att
//    属性类型

// func字段	att字段	说明
// 1	1	增加传球属性
// 	2	增加抢断属性
// 	3	增加盘带属性
// 	4	增加铲球属性
// 	5	增加射门属性
// 	6	增加守门属性
// 2	1	降低传球属性
// 	2	降低抢断属性
// 	3	降低盘带属性
// 	4	降低铲球属性
// 	5	降低射门属性
// 	6	降低守门属性
// 3	1	降低阵型克制效果
// 	2	降低战术克制效果

type SkillInfoList []SkillInfo

type SkillTypeStaticData struct {
	ID      int    //!技能类型id
	Name    string //!技能名
	Icon    int    //!技能图标
	Open    int    //!开启技能所需星级
	Type    int    //!技能触发类型 1->比分 2->我方阵型 3->敌方阵型 4->被动永久
	Sort    int    //!Type触发分类 详细见上方注释
	Tartype int    //!技能生效对象 1->我方  2->敌方
	Tarsort int    //!技能生效分类 详细见上方注释
	Func    int    //!技能产生效果 1->增加属性 2->降低属性 3->减少克制效果 4->特殊处理
	Attr    int    //!技能效果 详细见上方注释
	Power   int    //!技能效果百分比
	Skill   int    //!技能携带特效 详细参见skillnew表
	Time    int    //!技能学习所需花费时间 (单位: 秒)
	Desc1   string //!技能界面描述
	Desc2   string //!比赛界面描述
}

type SkillInfo struct { ///游戏中的技能信息
	ID     int `json:"id"`     ///技能id
	TeamID int `json:"teamid"` ///拥有球队id
	StarID int `json:"starid"` ///拥有者球星id
	Type   int `json:"type"`   ///技能类型
}

// type ISkill interface {
// 	GetInfo() *SkillInfo      ///得到信息对象指针
// 	SetStarID(starID int)     ///设置技能球星主人id
// 	Sync() DataValueChangList ///得到属性变更列表
// 	//GetTotalMergeExp() int           ///得到道具总融合经验值
// 	//AwardMergeExp(mergeExp int) bool ///加道具融合经验值,需要处理升级
// 	Save()                             ///马上保存数据
// 	SetPosition(position int)          ///设置技能装备位置
// 	GetID() int                        ///得到技能编号
// 	GetReflectValue() reflect.Value    ///得到反射对象
// 	GetTypeInfo() *SkillTypeStaticData ///取得技能类型静态数据信息
// }

type Skill struct { ///游戏中ItemInfo的道具
	SkillInfo
	DataUpdater
}

func (self *Skill) GetInfo() *SkillInfo {
	return &self.SkillInfo
}

func (self *Skill) GetTypeInfo() *SkillTypeStaticData { ///取得技能类型静态数据信息
	skillTypeStaticData := GetServer().GetStaticDataMgr().Unsafe().GetSkillType(self.Type)
	return skillTypeStaticData
}

func (self *Skill) GetReflectValue() reflect.Value { ///得到反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

func NewSkill(info *SkillInfo) *Skill {
	skill := new(Skill)
	skill.SkillInfo = *info
	skill.InitDataUpdater(tableSkill, &skill.SkillInfo)
	return skill
}
