package football

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
)

const (
	awardTypeCoin        = 400001 ///球币
	awardTypeTeamExp     = 400002 ///球队经验池加经验
	awardTypeTactic      = 400003 ///球队战术点
	awardTypeTalent      = 400004 ///球队培养点数
	awardTypeActionPoint = 400005 ///球队行动点数
	awardTypeTicket      = 400006 ///球队钻石
	awardTypeVipExp      = 400007 ///vip经验
	awardTypeManagerExp  = 400008 ///球员经理经验
)

const (
	functionMaskStarEducation    = 1  ///开启球员培养
	functionMaskItemEquip        = 2  ///开启道具装备
	functionMaskStarTrain        = 3  ///开启训练功能
	functionMaskStarSpy          = 4  ///开启球员挖掘
	functionMaskSkill            = 5  ///开启技能系统
	functionMaskVolunteer        = 6  ///开启球星来投
	functionMaskStarConvince     = 7  ///开启球员游说
	functionMaskStarEvolve       = 8  ///开启球员突破
	functionMaskChangeStar       = 9  ///更换球员
	functionMaskFindStar1        = 10 ///初级寻找球员
	functionMaskFindStar2        = 11 ///中级寻找球员
	functionMaskFindStar3        = 12 ///高级寻找球员
	functionMaskTrainMatch       = 13 ///训练赛
	functionMaskFormationUplevel = 14 ///阵型升级
	functionMaskArenaMatch       = 15 ///天天联赛
	functionMaskItemMerge        = 16 ///装备融合
	functionMaskTrainStar        = 17 ///球员训练
	functionMaskDayTask          = 18 ///日常任务功能开启
	functionMaskChallangeMatch   = 19 ///挑战赛
	functionMaskStarSkill        = 20 //!球星技能
	functionMaskMannaStar        = 21 //!天赐球员
)

const (
	playerMaxActionPoint = 150 ///玩家行动点上限(非RMB玩家)
)

//<<<<<<< .mine
//<<<<<<< .mine
//type ITeam interface { ///球队接口
//	GetStarCenter() IStarCenter
//	GetStarSpy() IStarSpy                         ///得到球探接口
//	GetInfo() *TeamInfo                           ///得到球队属性信息
//	SpendTicket(payTicket int, reson string) bool ///扣球队球票,reson为原因
//	AwardTicket(addTicket int) bool               ///奖励球队球票
//	AwardStar(addStarType int) int                ///奖励球队球员
//	GetStarInfoCopy(starID int) *StarInfo         ///得到球员信息复本
//	PayCoin(payCoin int) (bool, int)              ///消费球队金币
//	GetStarType(starID int) int                   ///得到球员类型
//	RemoveStar(starIDList IntList) bool           ///删除球队球员
//	GetCoin() int                                 ///得到球队金币
//	HasStar(starType int) bool                    ///是否已拥有指定类型的球员
//	IsStarInCurrentFormation(starID int) bool     ///判断球员是否在首发阵形中
//	PayTicket(payTicket int) int                  ///消费球队球票
//	AwardTrainCell(addTrainCell int) int          ///奖励球队训练位
//	GetTrainCell() int                            ///得到球队当前拥有训练位总数
//	GetTicket() int                               ///得到球队球票数
//	GetProcessMgr() IProcessMgr                   ///得到球队处理中心管理器组件
//	GetID() int                                   ///得到球队ID
//	//AwardStarExp(starID int, addExp int) int                    ///奖励球队球员经验
//	GetLevel() int                                                          ///得到球队等级
//	GetDrawGroupFilterIndexList(drawGroupIndexList []int) []int             ///得到已排除已经在球队中的球员类型索引列表
//	GetResetAttribMgr() IResetAttribMgr                                     ///得到球队可重置属性管理器组件
//	GetStar(starID int) IStar                                               ///得到球队球员对象
//	IsStarFull() bool                                                       ///判断球队中的已雇佣球员数是否已达上限
//	GetStarFromType(starType int) IStar                                     ///根据球员类型得到球队中球员对象
//	GetItemMgr() IGameItemMgr                                               ///得到球队道具管理器组件
//	GetSkillMgr() ISkillMgr                                                 ///得到球队技能管理器组件
//	PayTalentPoint(payTalentPoint int) bool                                 ///消费球队潜力点
//	Sync() DataValueChangList                                               ///得到属性变更列表
//	GetReflectValue() reflect.Value                                         ///得到球队反射对象
//	PayActionPoint(payActionPoint int) bool                                 ///消费球队活力值
//	GetCurrentFormation() int                                               ///得到球队当前阵形
//	SetCurrentFormation(formationID int)                                    ///设置球队当前阵形
//	GetFormationMgr() IFormationMgr                                         ///得到阵形管理器
//	GetLevelMgr() ILevelMgr                                                 ///得到球队关卡管理器组件
//	PayTacticalPoint(payTacticalPoint int) bool                             ///消费球队战术点
//	GetTaskMgr() ITaskMgr                                                   ///得到球队任务管理器组件
//	TestFunctionMask(maskFunction int) bool                                 ///测试指定功能掩码是否已开启
//	SetFunctionMask(maskFunction int, maskValue int)                        ///设置指定功能掩码开启或关闭
//	GetScore() int                                                          ///得到评分
//	CalcMatchScore() int                                                    ///计算评分
//	AwardExp(addExp int) bool                                               ///奖励球队经验
//	AwardTalentPoint(addTalentPoint int) bool                               ///奖励球队潜力点
//	AwardCoin(addCoin int) bool                                             ///奖励球队球币
//	CalcScore() (float32, float32, float32)                                 ///计算球队评分,根据所有球员数据生成攻防组织力
//	GetCurrentFormObject() IFormation                                       ///得到当前阵形对象
//	AddFormationLevel(addLevel int) bool                                    ///设置球队当前阵形
//	FindStarCount(seatPos int, color int, starCount int, starScore int) int ///球队得到指定位置指定品质颜色球员个数
//	GetTotalStarCount() int                                                 ///得到当前球员所有球员总数
//	SpendStarsContractPoint(payContractPoint int)                           ///将球队所有上阵球员契约值减指定值
//	IsStarsContractPointEnough() bool                                       ///判断球队所有上阵球员契约值是否均大于0
//	GetStarFateMgr() IStarFateMgr                                           ///得到球员缘管理器
//	IsStoreFull(itemType int, itemCount int) bool                           ///判断仓库是否已满
//	SetLobbyMask(maskValue int)                                             ///设置球员游说功能掩码开启或关闭
//	AwardAttrib(attribType int, addAtrib int) bool                          ///根据不同属性类型奖励属性值
//}
//=======
//type ITeam interface { ///球队接口
//	GetStarCenter() IStarCenter
//	GetStarSpy() IStarSpy                         ///得到球探接口
//	GetInfo() *TeamInfo                           ///得到球队属性信息
//	SpendTicket(payTicket int, reson string) bool ///扣球队球票,reson为原因
//	AwardTicket(addTicket int) bool               ///奖励球队球票
//	AwardStar(addStarType int) int                ///奖励球队球员
//	GetStarInfoCopy(starID int) *StarInfo         ///得到球员信息复本
//	PayCoin(payCoin int) (bool, int)              ///消费球队金币
//	GetStarType(starID int) int                   ///得到球员类型
//	RemoveStar(starIDList IntList) bool           ///删除球队球员
//	GetCoin() int                                 ///得到球队金币
//	HasStar(starType int) bool                    ///是否已拥有指定类型的球员
//	IsStarInCurrentFormation(starID int) bool     ///判断球员是否在首发阵形中
//	PayTicket(payTicket int) int                  ///消费球队球票
//	AwardTrainCell(addTrainCell int) int          ///奖励球队训练位
//	GetTrainCell() int                            ///得到球队当前拥有训练位总数
//	GetTicket() int                               ///得到球队球票数
//	GetProcessMgr() IProcessMgr                   ///得到球队处理中心管理器组件
//	GetID() int                                   ///得到球队ID
//	//AwardStarExp(starID int, addExp int) int                    ///奖励球队球员经验
//	GetLevel() int                                                          ///得到球队等级
//	GetDrawGroupFilterIndexList(drawGroupIndexList []int) []int             ///得到已排除已经在球队中的球员类型索引列表
//	GetDrawGroupFullIndexList(drawGroupIndexList []int) []int               ///得到所有球员类型索引列表
//	GetResetAttribMgr() IResetAttribMgr                                     ///得到球队可重置属性管理器组件
//	GetStar(starID int) IStar                                               ///得到球队球员对象
//	IsStarFull() bool                                                       ///判断球队中的已雇佣球员数是否已达上限
//	GetStarFromType(starType int) IStar                                     ///根据球员类型得到球队中球员对象
//	GetItemMgr() IGameItemMgr                                               ///得到球队道具管理器组件
//	GetSkillMgr() ISkillMgr                                                 ///得到球队技能管理器组件
//	PayTalentPoint(payTalentPoint int) bool                                 ///消费球队潜力点
//	Sync() DataValueChangList                                               ///得到属性变更列表
//	GetReflectValue() reflect.Value                                         ///得到球队反射对象
//	PayActionPoint(payActionPoint int) bool                                 ///消费球队活力值
//	GetCurrentFormation() int                                               ///得到球队当前阵形
//	SetCurrentFormation(formationID int)                                    ///设置球队当前阵形
//	GetFormationMgr() IFormationMgr                                         ///得到阵形管理器
//	GetLevelMgr() ILevelMgr                                                 ///得到球队关卡管理器组件
//	PayTacticalPoint(payTacticalPoint int) bool                             ///消费球队战术点
//	GetTaskMgr() ITaskMgr                                                   ///得到球队任务管理器组件
//	TestFunctionMask(maskFunction int) bool                                 ///测试指定功能掩码是否已开启
//	SetFunctionMask(maskFunction int, maskValue int)                        ///设置指定功能掩码开启或关闭
//	GetScore() int                                                          ///得到评分
//	CalcMatchScore() int                                                    ///计算评分
//	AwardExp(addExp int) bool                                               ///奖励球队经验
//	AwardTalentPoint(addTalentPoint int) bool                               ///奖励球队潜力点
//	AwardCoin(addCoin int) bool                                             ///奖励球队球币
//	CalcScore() (float32, float32, float32)                                 ///计算球队评分,根据所有球员数据生成攻防组织力
//	GetCurrentFormObject() IFormation                                       ///得到当前阵形对象
//	AddFormationLevel(addLevel int) bool                                    ///设置球队当前阵形
//	FindStarCount(seatPos int, color int, starCount int, starScore int) int ///球队得到指定位置指定品质颜色球员个数
//	GetTotalStarCount() int                                                 ///得到当前球员所有球员总数
//	SpendStarsContractPoint(payContractPoint int)                           ///将球队所有上阵球员契约值减指定值
//	IsStarsContractPointEnough() bool                                       ///判断球队所有上阵球员契约值是否均大于0
//	GetStarFateMgr() IStarFateMgr                                           ///得到球员缘管理器
//	IsStoreFull(itemType int, itemCount int) bool                           ///判断仓库是否已满
//	SetLobbyMask(maskValue int)                                             ///设置球员游说功能掩码开启或关闭
//	AwardAttrib(attribType int, addAtrib int) bool                          ///根据不同属性类型奖励属性值
//}
//>>>>>>> .r808
//=======
//type ITeam interface { ///球队接口
//	GetStarCenter() IStarCenter
//	GetStarSpy() IStarSpy                         ///得到球探接口
//	GetInfo() *TeamInfo                           ///得到球队属性信息
//	SpendTicket(payTicket int, reson string) bool ///扣球队球票,reson为原因
//	AwardTicket(addTicket int) bool               ///奖励球队球票
//	AwardStar(addStarType int) int                ///奖励球队球员
//	GetStarInfoCopy(starID int) *StarInfo         ///得到球员信息复本
//	PayCoin(payCoin int) (bool, int)              ///消费球队金币
//	GetStarType(starID int) int                   ///得到球员类型
//	RemoveStar(starIDList IntList) bool           ///删除球队球员
//	GetCoin() int                                 ///得到球队金币
//	HasStar(starType int) bool                    ///是否已拥有指定类型的球员
//	IsStarInCurrentFormation(starID int) bool     ///判断球员是否在首发阵形中
//	PayTicket(payTicket int) int                  ///消费球队球票
//	AwardTrainCell(addTrainCell int) int          ///奖励球队训练位
//	GetTrainCell() int                            ///得到球队当前拥有训练位总数
//	GetTicket() int                               ///得到球队球票数
//	GetProcessMgr() IProcessMgr                   ///得到球队处理中心管理器组件
//	GetID() int                                   ///得到球队ID
//	//AwardStarExp(starID int, addExp int) int                    ///奖励球队球员经验
//	GetLevel() int                                                          ///得到球队等级
//	GetDrawGroupFilterIndexList(drawGroupIndexList []int) []int             ///得到已排除已经在球队中的球员类型索引列表
//	GetDrawGroupFullIndexList(drawGroupIndexList []int) []int               ///得到所有球员类型索引列表
//	GetResetAttribMgr() IResetAttribMgr                                     ///得到球队可重置属性管理器组件
//	GetStar(starID int) IStar                                               ///得到球队球员对象
//	IsStarFull() bool                                                       ///判断球队中的已雇佣球员数是否已达上限
//	GetStarFromType(starType int) IStar                                     ///根据球员类型得到球队中球员对象
//	GetItemMgr() IGameItemMgr                                               ///得到球队道具管理器组件
//	GetSkillMgr() ISkillMgr                                                 ///得到球队技能管理器组件
//	PayTalentPoint(payTalentPoint int) bool                                 ///消费球队潜力点
//	Sync() DataValueChangList                                               ///得到属性变更列表
//	GetReflectValue() reflect.Value                                         ///得到球队反射对象
//	PayActionPoint(payActionPoint int) bool                                 ///消费球队活力值
//	GetCurrentFormation() int                                               ///得到球队当前阵形
//	SetCurrentFormation(formationID int)                                    ///设置球队当前阵形
//	GetFormationMgr() IFormationMgr                                         ///得到阵形管理器
//	GetLevelMgr() ILevelMgr                                                 ///得到球队关卡管理器组件
//	PayTacticalPoint(payTacticalPoint int) bool                             ///消费球队战术点
//	GetTaskMgr() ITaskMgr                                                   ///得到球队任务管理器组件
//	TestFunctionMask(maskFunction int) bool                                 ///测试指定功能掩码是否已开启
//	SetFunctionMask(maskFunction int, maskValue int)                        ///设置指定功能掩码开启或关闭
//	GetScore() int                                                          ///得到评分
//	CalcMatchScore() int                                                    ///计算评分
//	AwardExp(addExp int) bool                                               ///奖励球队经验
//	AwardTalentPoint(addTalentPoint int) bool                               ///奖励球队潜力点
//	AwardCoin(addCoin int) bool                                             ///奖励球队球币
//	CalcScore() (float32, float32, float32)                                 ///计算球队评分,根据所有球员数据生成攻防组织力
//	GetCurrentFormObject() IFormation                                       ///得到当前阵形对象
//	AddFormationLevel(addLevel int) bool                                    ///设置球队当前阵形
//	FindStarCount(seatPos int, color int, starCount int, starScore int) int ///球队得到指定位置指定品质颜色球员个数
//	GetTotalStarCount() int                                                 ///得到当前球员所有球员总数
//	SpendStarsContractPoint(payContractPoint int)                           ///将球队所有上阵球员契约值减指定值
//	IsStarsContractPointEnough() bool                                       ///判断球队所有上阵球员契约值是否均大于0
//	GetStarFateMgr() IStarFateMgr                                           ///得到球员缘管理器
//	IsStoreFull(itemType int, itemCount int) bool                           ///判断仓库是否已满
//	SetLobbyMask(maskValue int)                                             ///设置球员游说功能掩码开启或关闭
//	AwardAttrib(attribType int, addAtrib int) bool                          ///根据不同属性类型奖励属性值
//	GetAverageLevel() int                                                   ///得到球队平均等级
//	GetMinimumLevelStar(minimumLevel int) IStar                             ///得到球队最低等级
//}
//>>>>>>> .r976

type TeamInfo struct { ///球队信息,和数据库表dy_team一一对应
	ID               int    `json:"id"`               ///球队id
	Name             string `json:"name"`             ///名字
	AccountID        int    `json:"accountid"`        ///账号名
	ClubID           int    `json:"clubid"`           ///球会id,未加入为0
	Level            int    `json:"level"`            ///等级
	FormationID      int    `json:"formationid"`      ///当前阵型id
	FormationLevel   int    `json:"formationlevel"`   ///阵型等级,所有阵型共享一个等级',
	Icon             int    `json:"icon"`             ///队徽
	Coin             int    `json:"coin"`             ///游戏币
	TeamShirts       int    `json:"teamshirts"`       ///球队球衣
	Exp              int    `json:"exp"`              ///经验值
	ActionPoint      int    `json:"actionpoint"`      ///行动点
	TacticalPoint    int    `json:"tacticalpoint"`    ///战术点
	Ticket           int    `json:"ticket"`           ///球票,人民币
	TrainCount       int    `json:"traincount"`       ///训练位
	TalentPoint      int    `json:"talentpoint"`      ///潜能点数,用于球员属性加点
	FunctionMask     int    `json:"functionmask"`     ///功能掩码,指定bit位是1开启0关闭
	Score            int    `json:"score"`            ///战力评分
	VipLevel         int    `json:"viplevel"`         ///VIP等级
	StoreCapacity    int    `json:"storecapacity"`    ///球队仓库容量,默认20格,最大256格
	LobbyMask        int64  `json:"lobbymask"`        ///球星游说掩码,指定bit位是1已领取0未领取
	StarExpPool      int    `json:"starexppool"`      ///球队经验池
	VipExp           int    `json:"vipexp"`           ///VIP经验
	MakeTime         int    `json:"maketime"`         ///球队创建时间 年月日时分 1408010123
	Restoretime      int    `json:"restoretime"`      ///下次回复行动点时间
	Operation        string `json:"operation"`        ///客户端操作串
	ActivitCodeAward int64  `json:"activitcodeaward"` ///激活码奖励
	LoginAndPayAward int64  `json:"loginandpayaward"` ///登陆和充值送球员奖励掩码
	AddStarPos       int    `json:"addstarpos"`       //! 增加的球星位
}

type StarList map[int]*Star

type VipPrivilegeStaticData struct { ///VIP特权信息静态数据
	ID       int ///VIP等级
	Recharge int ///充值金额
	Param1   int ///免费获得球币次数
	Param2   int ///购买球币次数
	Param3   int ///购买体力次数
	Param4   int ///购买杯赛次数
	Param5   int ///购买传奇赛次数
	Param6   int ///短程机票累计数量
	Param7   int ///中途机票累计数量
	Param8   int ///长途机票累计数量
	Param9   int ///球队仓库格
	Param10  int ///每天获得星卡福利
	Param11  int ///球星来投免费刷新次数
	Param12  int ///解雇返回训练点百分比
	Param13  int ///解雇返还经验百分比
	Param14  int ///球队整顿时间减少
	Param15  int ///杯赛复活次数
	Param16  int ///传奇赛额外增加
}

///type StarList map[int]IStar

type GameMgrList map[int]IGameMgr ///游戏逻辑管理器

type Team struct {
	TeamInfo
	DataUpdater              ///球队信息保存组件
	TeamInfoCalc             ///计算后二级属性
	gameMgrList  GameMgrList ///游戏逻辑管理器
	starList     StarList    ///球星列表
	starSpy      *StarSpy    ///球探对象
	starCenter   *StarCenter ///球员中心
	client       *Client     ///客户端	对象
	//vipShop      VipShopMgr  ///商城系统
	//processMgr   IProcessMgr ///处理中心管理器组件
	//resetAttribMgr IResetAttribMgr ///可重置属性数据管理器组件
	//itemMgr IGameItemMgr ///道具管理器

	//! cy
	PowerValue      *GetPower   //! 领取体力
	OSActivityValue *OSActivity //! 开服活动
}

//type ITeam *Team

type TeamInfoCalc struct { ///球员计算后属性
	AttackScoreCalc   float32 ///计算后攻击力评分
	DefenseScoreCalc  float32 ///计算后防御力评分
	OrganizeScoreCalc float32 ///计算后组织力评分
}

func (self *Team) GetFormation() int { //得到阵型
	return self.FormationID
}

func (self *Team) IsStoreFull(itemType int, itemCount int) bool { ///判断仓库是否已满
	isStoreFull := false
	itemMgr := self.GetItemMgr()
	_, needCellCount := itemMgr.TryComboItem(itemType, itemCount) ///先判断是否能完全叠加
	if needCellCount <= 0 {
		return false ///如果可以完全叠加则仓库没满
	}
	currentStoreItemCount := itemMgr.GetItemCountByPos(itemPosStore) ///得到当前仓库道具数
	///判断此时仓库道具总数加新道具是否到达上限
	isStoreFull = currentStoreItemCount+needCellCount >= self.StoreCapacity
	return isStoreFull
}

func (self *Team) CalcItemScore() { ///计算装备对球队的攻守加成
	///得到首发名单
	itemMgr := self.GetItemMgr()
	currentFormation := self.GetFormationMgr().GetFormation(self.FormationID)
	starIDList := currentFormation.GetStarIDList()
	attackscore, defensescore, organizescore := itemMgr.GetEquipAttributeAddition(starIDList) ///获取装备属性加成
	self.AttackScoreCalc += float32(attackscore)
	self.DefenseScoreCalc += float32(defensescore)
	self.OrganizeScoreCalc += float32(organizescore)
}

func (self *Team) CalcScore() (float32, float32, float32) { ///计算球队评分,根据所有球员数据生成攻防组织力
	self.Score = 0             ///球队评分清空
	self.AttackScoreCalc = 0   ///球队攻击力评分清空
	self.DefenseScoreCalc = 0  ///球队防御力评分清空
	self.OrganizeScoreCalc = 0 ///球队组织力评分清空
	staticDataMgr := GetServer().GetStaticDataMgr().Unsafe()
	currentFormation := self.GetCurrentFormObject()
	formationInfo := currentFormation.GetInfo()
	seatPosList := staticDataMgr.GetFieldValueList(formationInfo, "Pos", 1, 11)            ///得到每个位置的球员id列表
	formationTypeStaticData := staticDataMgr.GetFormationType(formationInfo.Type)          ///得到阵形类型对象
	seatTypeList := staticDataMgr.GetFieldValueList(formationTypeStaticData, "Pos", 1, 11) ///得到此阵形的位置类型列表
	for i := range seatPosList {
		starID := seatPosList[i]      ///得到球员id
		seatType := seatTypeList[i]   ///得到此球员当前踢的位置
		star := self.starList[starID] ///得到球员对象
		starScore := star.CalcScore() ///计算球员评分
		self.Score += int(starScore)  ///累加球员评分
		seatTypeStaticData := staticDataMgr.GetSeatType(seatType)
		self.AttackScoreCalc += (starScore * float32(seatTypeStaticData.AttackRate))     ///球员贡献的攻击力累加
		self.DefenseScoreCalc += (starScore * float32(seatTypeStaticData.DefenseRate))   ///球员贡献的防御力累加
		self.OrganizeScoreCalc += (starScore * float32(seatTypeStaticData.OrganizeRate)) ///球员贡献的组织力累加
	}
	self.CalcItemScore()          ///计算装备对攻击,防御,组织的加成
	self.AttackScoreCalc /= 100   ///缩放攻击力
	self.DefenseScoreCalc /= 100  ///缩放防御力
	self.OrganizeScoreCalc /= 100 ///缩放组织力
	self.Score = int(math.Pow(float64(self.AttackScoreCalc)*float64(self.DefenseScoreCalc)*math.Pow(float64(self.OrganizeScoreCalc), 1.2), 0.33333333))
	return self.AttackScoreCalc, self.DefenseScoreCalc, self.OrganizeScoreCalc
}

func (self *Team) CalcTeamPower() int {
	self.CalcScore()
	return self.Score
}

func (self *Team) GetCurrentFormObject() *Formation {
	return self.GetFormationMgr().GetFormation(self.FormationID)
}

//func (self *Team) IsGameMaster() bool {
//	isGM := strings.Index(self.Name, "_GM")
//	return isGM != -1
//}

func (self *Team) TestFunctionMask(maskFunction int) bool {
	result := self.FunctionMask & maskFunction
	return result > 0
}

func (self *Team) SetFunctionMask(maskFunction int, maskValue int) {
	self.FunctionMask |= maskValue << uint(maskFunction-1)
}

func (self *Team) SetLobbyMask(maskValue int) { ///设置球星游说掩码
	self.LobbyMask |= 1 << uint(maskValue)
}

func (self *Team) cloneCurrentFrom(formationID int) { ///设置球队当前阵形
	srcForm := self.GetCurrentFormObject()
	dstForm := self.GetFormationMgr().GetFormation(formationID)
	dstForm.Pos1 = srcForm.Pos1
	dstForm.Pos2 = srcForm.Pos2
	dstForm.Pos3 = srcForm.Pos3
	dstForm.Pos4 = srcForm.Pos4
	dstForm.Pos5 = srcForm.Pos5
	dstForm.Pos6 = srcForm.Pos6
	dstForm.Pos7 = srcForm.Pos7
	dstForm.Pos8 = srcForm.Pos8
	dstForm.Pos9 = srcForm.Pos9
	dstForm.Pos10 = srcForm.Pos10
	dstForm.Pos11 = srcForm.Pos11
}

func (self *Team) SetCurrentFormation(formationID int) { ///设置球队当前阵形
	if self.FormationID == formationID {
		return ///避免重复设置
	}
	formationMgr := self.GetFormationMgr()
	formationMgr.CorrectFormError(formationID)
	self.FormationID = formationID
}

func (self *Team) GetCurrentFormation() int { ///得到球队当前阵形
	return self.FormationID
}

func (self *Team) GetReflectValue() reflect.Value { ///得到球队反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

func (self *Team) GetTaskMgr() *TaskMgr { ///得到球队任务管理器组件
	return self.gameMgrList[mgrTypeTaskMgr].(*TaskMgr)
}

func (self *Team) AddFormationLevel(addLevel int) bool { ///设置球队当前阵形
	if addLevel <= 0 {
		return false
	}
	self.FormationLevel += addLevel
	return true
}

func (self *Team) GetStarFateMgr() IStarFateMgr { ///得到球员缘管理器
	return self.gameMgrList[mgrTypeStarFateMgr].(IStarFateMgr)
}

func (self *Team) GetFormationMgr() *FormationMgr { ///得到阵形管理器
	return self.gameMgrList[mgrTypeFormationMgr].(*FormationMgr)
}

func (self *Team) GetLevelMgr() *LevelMgr { ///得到球队关卡管理器组件
	return self.gameMgrList[mgrTypeLevelMgr].(*LevelMgr)
}

func (self *Team) GetSkillMgr() ISkillMgr { ///得到球队技能管理器组件
	return self.gameMgrList[mgrTypeSkillMgr].(ISkillMgr)
}

func (self *Team) GetResetAttribMgr() *ResetAttribMgr { ///得到球队可重置属性管理器组件
	return self.gameMgrList[mgrTypeResetAttribMgr].(*ResetAttribMgr)
}

func (self *Team) GetItemMgr() *ItemMgr { ///得到球队道具管理器组件
	return self.gameMgrList[mgrTypeItemMgr].(*ItemMgr)
}

func (self *Team) GetProcessMgr() IProcessMgr { ///得到球队处理中心管理器组件
	return self.gameMgrList[mgrTypeProcessMgr].(IProcessMgr)
}

func (self *Team) GetVipShopMgr() *VipShopMgr {
	return self.gameMgrList[mgrTypeVipShopMgr].(*VipShopMgr)
}

func (self *Team) GetArenaMgr() *ArenaMgr { ///得到球队竞技场管理器组件
	return self.gameMgrList[mgrTypeArenaMgr].(*ArenaMgr)
}

func (self *Team) GetActivityMgr() *ActivityMgr { ///得到球队活动管理器组件
	return self.gameMgrList[mgrTypeActivityMgr].(*ActivityMgr)
}

func (self *Team) GetGameMgr(gameMgrType int) IGameMgr { ///得到球队游戏逻辑管理器组件
	return self.gameMgrList[gameMgrType]
}

func (self *Team) GetMailMgr() *MailMgr {
	return self.gameMgrList[mgrTypeMailMgr].(*MailMgr)
}

func (self *Team) AddGameMgrList(gameMgr IGameMgr, syncMgr *SyncMgr) bool { ///注册游戏逻辑系统对象
	if nil == gameMgr {
		return false
	}
	gameMgrType := gameMgr.GetType()
	gameMgr.SetSyncMgr(syncMgr)
	gameMgr.SetTeam(self) ///设置球队对象
	self.gameMgrList[gameMgrType] = gameMgr
	gameMgr.onInit()
	return true
}

func (self *Team) InitGameMgrList(syncMgr *SyncMgr) bool { ///初始化球队游戏逻辑管理器组件
	self.gameMgrList = make(GameMgrList)
	if self.AddGameMgrList(NewSkillMgr(self.ID), syncMgr) == false { ///创建技能管理器
		return false
	}
	if self.AddGameMgrList(NewLevelMgr(self.ID), syncMgr) == false { ///创建关卡管理器
		return false
	}
	if self.AddGameMgrList(NewFormationMgr(self.ID), syncMgr) == false { ///创建阵型管理器
		return false
	}
	if self.AddGameMgrList(NewTaskMgr(self.ID), syncMgr) == false { ///创建任务管理器
		return false
	}
	if self.AddGameMgrList(NewResetAttribMgr(self.ID), syncMgr) == false { ///创建可重置属性管理器
		return false
	}
	if self.AddGameMgrList(NewItemMgr(self.ID), syncMgr) == false { ///创建道具管理器
		return false
	}
	if self.AddGameMgrList(NewProcessMgr(self.ID), syncMgr) == false { ///创建进程处理管理器
		return false
	}
	if self.AddGameMgrList(NewStarFateMgr(self.ID), syncMgr) == false { ///创建球员缘系统管理器
		return false
	}
	if self.AddGameMgrList(NewArenaMgr(self.ID), syncMgr) == false { ///创建竞技场管理器
		return false
	}
	if self.AddGameMgrList(NewVipShopMgr(self.ID), syncMgr) == false { ///创建商城管理器
		return false
	}
	if self.AddGameMgrList(NewMailMgr(self.ID), syncMgr) == false { ///创建邮件管理器
		return false
	}
	if self.AddGameMgrList(NewMannaStarMgr(self.ID), syncMgr) == false { //!创建天赐球员管理器
		return false
	}
	if self.AddGameMgrList(NewAtlasMgr(self.ID), syncMgr) == false { //! 创建图鉴管理器失败
		return false
	}
	self.GetFormationMgr().CorrectAllFormError()
	//	self.GetArenaMgr().RefreshArenaDate()
	return true
}

func (self *Team) OnLogout() { ///球队登出时逻辑处理
	self.Save() ///球队退出时需要保存数据
}

func (self *Team) Save() { ///保存球队信息
	self.DataUpdater.Save()
	self.starSpy.Save()                  ///球探保存数据
	self.starCenter.Save()               ///球员中心保存数据
	for _, star := range self.starList { ///保存所有球员信息
		star.Save()
	}
	//	self.resetAttribMgr.Save() ///保存了可重置属性数据
	//self.itemMgr.Save() ///道具管理器保存信息

	for _, v := range self.gameMgrList {
		v.SaveInfo()
	}

	//! cy
	self.PowerValue.Save()
	self.OSActivityValue.Save()
}

func (self *Team) IsStarInCurrentFormation(starID int) bool { ///判断球员是否在首发阵形中
	return self.GetFormationMgr().IsStarInFormation(self.FormationID, starID)
}

func (self *Team) GetStarType(starID int) int { ///得到球员类型
	star := self.starList[starID]
	if nil == star {
		return 0
	}
	starInfo := star.GetInfo()
	return starInfo.Type
}

func (self *Team) GetStarInfoCopy(starID int) *StarInfo { ///得到球员信息复本
	star := self.starList[starID]
	if nil == star {
		return nil
	}
	starInfo := new(StarInfo)
	*starInfo = *star.GetInfo()
	return starInfo
}

func (self *Team) RemoveStar(starIDList IntList) bool { ///删除球队球员列表
	starIDListLen := len(starIDList)
	if starIDListLen <= 0 {
		return false
	}
	removeStarQuery := fmt.Sprintf("delete from %s where id in(", tableStar)
	for i := range starIDList {
		starID := starIDList[i]
		if nil == self.starList[starID] {
			continue
		}
		delete(self.starList, starID) ///从内存中删除对象
		removeStarQuery += fmt.Sprintf("%d", starID)
		if i < starIDListLen-1 {
			removeStarQuery += ","
		}
	}
	removeStarQuery += fmt.Sprintf(") limit %d", starIDListLen) ///限制只能删除指定数量记录
	///从数据库中删除对象
	_, rowsStarAffected := GetServer().GetDynamicDB().Exec(removeStarQuery)
	if rowsStarAffected != starIDListLen {
		GetServer().GetLoger().Warn("Team RemoveStar removeStarQuery fail! starID:%v", starIDList)
		return false
	}
	return rowsStarAffected > 0
}

//func (self *Team) RemoveStar(starID int) bool { ///删除球队球员
//	///先删除内存防止玩家刷
//	if nil == self.starList[starID] {
//		return false
//	}
//	delete(self.starList, starID) ///从内存中删除对象
//	///从数据库中删除对象
//	removeStarQuery := fmt.Sprintf("delete from %s where id=%d limit 1", tableStar, starID)
//	_, rowsStarAffected := GetServer().GetDynamicDB().Exec(removeStarQuery)
//	if rowsStarAffected <= 0 {
//		GetServer().GetLoger().Warn("Team RemoveStar removeStarQuery fail! starID:%d", starID)
//		return false
//	}
//	return rowsStarAffected > 0
//}

func (self *Team) HasStar(starType int) bool { ///是否已拥有指定类型的球员
	for _, v := range self.starList {
		if v.GetInfo().Type == starType {
			return true
		}
	}
	return false
}

func (self *Team) AwardTactic(addTactic int) bool { ///奖励球队战术点
	if addTactic <= 0 {
		return false
	}
	self.TacticalPoint += addTactic
	return true
}

func (self *Team) AwardCoin(addCoin int) bool { ///奖励球队球币
	if addCoin <= 0 {
		return false
	}
	self.Coin += addCoin
	return true
}

func (self *Team) AwardExp(client IClient, addExp int) bool { ///奖励球队经验
	if addExp <= 0 {
		return false
	}
	self.Exp += addExp
	if self.Uplevel(client) {
		self.AwardObject(awardTypeActionPoint, 60, 0, 0)
	}
	return true
}

func (self *Team) AwardExpPool(addExp int) bool {
	if addExp <= 0 {
		return false
	}
	self.StarExpPool += addExp
	return true
}

//! 扩充球员
func (self *Team) AwardStarPos() int {
	maxpos := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configStarSpy, configItemDiscoverConfig, 5)
	if self.AddStarPos >= maxpos {
		return 1
	}

	needdiamond := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configStarSpy, configItemDiscoverConfig, 6)
	if self.Ticket < needdiamond {
		return 2
	}

	self.AddStarPos++
	self.Ticket -= needdiamond
	self.Save()

	self.GetItemMgr().syncMgr.SyncObject("AwardObject", self)

	return 0
}

func (self *Team) Uplevel(client IClient) bool { ///球队升级
	isLevel := false
	///得到当前等级升级所需经验
	staticDataMgr := GetServer().GetStaticDataMgr()
	levelExpCount := staticDataMgr.GetLevelExpCount(levelExpTypeTeamLevel)
	//	oldLevel := self.Level
	for i := 1; i <= levelExpCount; i++ {
		needExp := staticDataMgr.GetLevelExpNeedExp(levelExpTypeTeamLevel, self.Level) ///得到当前等级经验

		if self.Exp < needExp {
			break ///经验不足升级
		}

		if self.Level >= 100 {
			self.Exp = needExp ///满级后经验到达上限
			break              ///已满级
		}
		self.Level++
		isLevel = true
		self.Exp -= needExp
		if self.Exp < 0 {
			self.Exp = 0
		}
	}

	return isLevel

	//	if oldLevel != self.Level {
	//		///记录玩家升级
	//		client.LevelUpRecord(self.Level)
	//	}
}

//func (self *Team) AwardStarExp(starID int, addExp int) int { ///奖励球队球员经验
//	if addExp <= 0 {
//		return 0
//	}
//	star := self.starList[starID]
//	if star == nil {
//		return 0
//	}
//	starInfo := star.GetInfo()
//	starInfo.Exp += addExp
//	return starInfo.Exp
//}

//func (self *Team) AddStar(addStar *Star) bool { ///奖励球队球员
//	if nil == addStar {
//		return false
//	}
//	self.starList[addStar.ID] = addStar
//	return true
//}

func (self *Team) GetStar(starID int) *Star { ///得到球队球员对象
	return self.starList[starID]
}

func (self *Team) AwardStar(addStarType int) int { ///奖励球队球员
	if 0 == addStarType {
		return 0
	}
	if self.IsStarFull() == true {
		return 0
	}
	starType := GetServer().GetStaticDataMgr().Unsafe().GetStarType(addStarType)
	if nil == starType {
		return 0
	}
	star := self.GetStarFromType(addStarType)
	if nil != star {
		return 0
	}
	insertNewStarQuery := fmt.Sprintf("insert %s set teamid=%d,type=%d,grade=%d", tableStar, self.ID, addStarType, starType.Grade)
	lastInsertStarID, _ := GetServer().GetDynamicDB().Exec(insertNewStarQuery) ///创建新球员
	if lastInsertStarID <= 0 {
		GetServer().GetLoger().Warn("Team AwardStar Insert New Star fail! StarType:%d", addStarType)
		return 0
	}
	///创建star对象
	loadStarQuery := fmt.Sprintf("select * from %s where id=%d limit 1", tableStar, lastInsertStarID)
	starInfo := new(StarInfo)
	GetServer().GetDynamicDB().fetchOneRow(loadStarQuery, starInfo)
	starNew := NewStar(starInfo, self)
	starNew.CalcScore()
	self.starList[starInfo.ID] = starNew

	if starType.Class > 400 {
		//! A级以上加入图鉴
		atlasMgr := self.GetAtlasMgr()
		atlasMgr.AddAtlas(self.ID, starType.ID, 0)
	}
	//if false == ok {
	//	return 0
	//}
	//star.Init(nil) ///初始化Star对象
	//if self.AddStar(star) == false {
	//	return 0
	//}
	return lastInsertStarID
}

func (self *Team) AwardMannaStar(addStarType int) int { ///奖励球队球员
	if 0 == addStarType {
		return 0
	}
	if self.IsStarFull() == true {
		return 0
	}

	starType := self.GetMannaStarMgr().GetMannaStar(addStarType)
	if nil == starType {
		return 0
	}

	insertNewStarQuery := fmt.Sprintf("insert %s set teamid=%d,type=%d,grade=%d,ismannastar=1", tableStar, self.ID, addStarType, starType.Grade)
	lastInsertStarID, _ := GetServer().GetDynamicDB().Exec(insertNewStarQuery) ///创建新球员

	if lastInsertStarID <= 0 {
		GetServer().GetLoger().Warn("Team AwardStar Insert New Star fail! StarType:%d", addStarType)
		return 0
	}
	///创建star对象
	loadStarQuery := fmt.Sprintf("select * from %s where id=%d limit 1", tableStar, lastInsertStarID)
	starInfo := new(StarInfo)
	GetServer().GetDynamicDB().fetchOneRow(loadStarQuery, starInfo)
	starNew := NewStar(starInfo, self)
	starNew.CalcScore()
	self.starList[starInfo.ID] = starNew

	return lastInsertStarID
}

func (self *Team) GetLevel() int { ///得到球队等级
	return self.Level
}

func (self *Team) GetID() int { ///得到球队ID
	return self.ID
}

func (self *Team) GetTicket() int { ///得到球队球票数
	return self.Ticket
}

func (self *Team) GetCoin() int { ///得到球队金币
	return self.Coin
}

func (self *Team) PayTacticalPoint(payTacticalPoint int) bool { ///消费球队战术点
	if payTacticalPoint <= 0 {
		return false
	}
	if self.TacticalPoint < payTacticalPoint {
		return false
	}
	self.TacticalPoint -= payTacticalPoint
	return true
}

func (self *Team) AwardTalentPoint(addTalentPoint int) bool { ///奖励球队潜力点
	if addTalentPoint <= 0 {
		return false
	}
	//	self.TalentPoint += addTalentPoint
	self.Coin += addTalentPoint
	return true
}

func (self *Team) PayTalentPoint(payTalentPoint int) bool { ///消费球队潜力点
	if payTalentPoint <= 0 {
		return false
	}
	// if self.TalentPoint < payTalentPoint {
	// 	return false
	// }
	// self.TalentPoint -= payTalentPoint

	if self.Coin < payTalentPoint {
		return false
	}
	self.Coin -= payTalentPoint
	return true
}

func (self *Team) PayCoin(payCoin int) (bool, int) { ///消费球队金币
	if payCoin <= 0 {
		return false, 0
	}
	if self.Coin < payCoin {
		return false, 0
	}
	self.Coin -= payCoin
	return true, self.Coin
}

func (self *Team) PayActionPoint(payActionPoint int) bool { ///消费球队活力值
	if payActionPoint <= 0 {
		return false
	}
	if self.ActionPoint < payActionPoint {
		return false
	}
	self.ActionPoint -= payActionPoint
	//	self.CheckActionPoint()
	return true
}

func (self *Team) CreateShopTimesDefaultRefresh() *ResetAttrib {
	refreshShopTimesHours := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTeam,
		configTeamRestore, 2) ///取得配置表中每日刷新小时数
	resettime := GetHourUTC(refreshShopTimesHours)
	resetAttribMgr := self.GetResetAttribMgr()
	resetAttribMgr.AddResetAttrib(ResetAttribTypeTeamShopTimes, resettime, IntList{0, 0, 0})
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamShopTimes)
	return resetAttrib
}

//func (self *Team) CheckActionPoint() { ///检查行动点是否低于上限
//	if self.ActionPoint >= playerMaxActionPoint || self.Restoretime != 0 { ///若之前有恢复CD,则不设置恢复时间
//		return
//	}

//	///小于上限,则开始恢复  恢复速度: 12分钟1体力
//	now := Now()
//	restoreTime := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTeam, configTeamRestore, 1)
//	if restoreTime == 0 {
//		GetServer().GetLoger().Warn("Action point restore time is 0 teamid:%d  config subType:%s", self.ID, configTeamRestore)
//		self.Restoretime = 0
//		return
//	}
//	self.Restoretime = now + restoreTime
//}

func (self *Team) RestoreActionPoint(now int) { ///恢复行动点
	if now < self.Restoretime {
		return ///未到恢复时间
	}
	if self.ActionPoint >= playerMaxActionPoint {
		self.Restoretime = 0 ///不断重置时钟
		return               ///超过上限了
	}
	const addActionPoint = 1 ///每次增加1点行动力
	restoreTime := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTeam, configTeamRestore, 1)
	if self.Restoretime <= 0 {
		self.Restoretime = now + restoreTime ///第一次触发时间处理
		return
	}
	diffTimeSecs := now - self.Restoretime + restoreTime          ///得到相差秒数
	restoreCount := diffTimeSecs / restoreTime                    ///得到恢复次数
	remainSecs := diffTimeSecs - (restoreTime * restoreCount)     ///得到触发后剩余的秒数
	awardActionPoint := restoreCount * addActionPoint             ///得到总奖励行动点数
	addFullActionPoint := playerMaxActionPoint - self.ActionPoint ///得到补满所需行动点数
	awardActionPoint = Min(awardActionPoint, addFullActionPoint)  ///修正奖励行动点数以免超出上限
	self.Restoretime = now + restoreTime + remainSecs             ///生成新的加行动点时间
	self.AwardObject(awardTypeActionPoint, awardActionPoint, 0, 0)

	//if now < self.Restoretime || 0 == self.Restoretime {
	//	return ///不到恢复时间/不需要恢复
	//}

	//if now == self.Restoretime {
	//	self.AwardObject(awardTypeActionPoint, 1, 0, 0)
	//	// self.ActionPoint += 1 ///行动点恢复

	//	if self.ActionPoint >= playerMaxActionPoint {
	//		self.Restoretime = 0 ///达到上限
	//	} else {
	//		restoreTime := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTeam, configTeamRestore, 1)
	//		if restoreTime == 0 {
	//			GetServer().GetLoger().Warn("Action point restore time is 0 teamid:%d  config subType:%s", self.ID, configTeamRestore)
	//			self.Restoretime = 0
	//			return
	//		}
	//		self.Restoretime = now + restoreTime
	//	}

	//	client.GetSyncMgr().SyncObject("restoreActionPoint", self) ///同步客户端
	//	return
	//}

	//if now > self.Restoretime {
	//	// 隔夜情况
	//	restoreTime := GetServer().GetStaticDataMgr().GetConfigStaticDataInt(configTeam, configTeamRestore, 1)
	//	if restoreTime == 0 {
	//		GetServer().GetLoger().Warn("Action point restore time is 0 teamid:%d  config subType:%s", self.ID, configTeamRestore)
	//		self.Restoretime = 0
	//		return
	//	}

	//	restoreCount := (now - self.Restoretime + restoreTime) / restoreTime //有多少次恢复
	//	curActionPoint := Min(self.ActionPoint+restoreCount, 150)
	//	awardActionPoint := curActionPoint - self.ActionPoint
	//	self.AwardObject(awardTypeActionPoint, awardActionPoint, 0, 0)

	//	curTime := (now - self.Restoretime + restoreTime) % restoreTime //当前循环剩余秒数
	//	self.Restoretime = (restoreTime - curTime) + now
	//}
}

func (self *Team) GetDrawGroupFullIndexList(drawGroupIndexList []int) []int { ///得到已排除已经在球队中的球员类型索引列表
	staticDataMgr := GetServer().GetStaticDataMgr()
	///建立查询索引
	storeStarTypeList := make(map[int]bool)
	for _, v := range self.starList {
		starInfo := v.GetInfo()
		storeStarTypeList[starInfo.Type] = true
	}
	drawGroupFilterIndexList := []int{}
	for i := range drawGroupIndexList {
		drawGroupIndex := drawGroupIndexList[i]
		drawGroupStaticData := staticDataMgr.GetDrawGroupStaticData(drawGroupIndex)
		if nil == drawGroupStaticData {
			break
		}

		if drawGroupStaticData.ReqLevel > self.Level {
			continue
		}

		if drawGroupStaticData.ReqPass != 0 {

			//得到前置球队对应关卡
			levelID, index := staticDataMgr.GetLevelFromNpcteamID(drawGroupStaticData.ReqPass)
			if levelID == 0 {
				GetServer().GetLoger().Warn("noone in level reward  RollVolunteerStarTypeList() levelType = %d", drawGroupStaticData.ReqPass)
				continue
			}

			levelMgr := self.GetLevelMgr()
			level := levelMgr.FindLevel(levelID)
			if level == nil {
				GetServer().GetLoger().Warn("level is not exist func:RollVolunteerStarTypeList levelType = %d level", drawGroupStaticData.ReqPass)
				continue
			}
			if level.IsPass() == false { //未通过关卡,则淘汰
				continue
			}

			levelInfo := level.GetInfoPtr()
			starCount := index * 3
			if levelInfo.StarCount != starCount {
				continue
			}
		}

		if drawGroupStaticData.ReqMoney > self.Coin {
			continue
		}

		drawGroupFilterIndexList = append(drawGroupFilterIndexList, drawGroupIndex)
	}
	return drawGroupFilterIndexList
}

func (self *Team) GetDrawGroupFilterIndexList(drawGroupIndexList []int) []int { ///得到已排除已经在球队中的球员类型索引列表
	staticDataMgr := GetServer().GetStaticDataMgr()
	///建立查询索引
	storeStarTypeList := make(map[int]bool)
	for _, v := range self.starList {
		storeStarTypeList[v.GetInfo().Type] = true
	}
	drawGroupFilterIndexList := []int{}
	for i := range drawGroupIndexList {
		drawGroupIndex := drawGroupIndexList[i]
		drawGroupStaticData := staticDataMgr.GetDrawGroupStaticData(drawGroupIndex)
		if nil == drawGroupStaticData {
			GetServer().GetLoger().Warn("Team GetDrawGroupFilterIndexList fail! drawGroupIndex %d is Invalid!", drawGroupIndex)
			break
		}
		_, ok := storeStarTypeList[drawGroupStaticData.AwardType]
		if false == ok {
			drawGroupFilterIndexList = append(drawGroupFilterIndexList, drawGroupIndex)
		}
	}
	return drawGroupFilterIndexList
}

func (self *Team) IsGM() bool { ///判断自己是不是GM
	isGM := strings.Contains(self.Name, "[GM]")
	return isGM
}

func (self *Team) PayTicket(payTicket int) int { ///消费球队球票
	if payTicket <= 0 {
		return 0
	}
	if payTicket > self.Ticket {
		return 0 ///余额不足
	}
	self.Ticket -= payTicket
	self.Save()
	return self.Ticket
}

func (self *Team) GetTrainCell() int { ///得到球队当前拥有训练位总数
	return self.TrainCount
}

func (self *Team) GetScore() int {
	return self.Score
}

func (self *Team) CalcMatchScore() int {
	return self.GetScore()
}

func (self *Team) AwardTrainCell(addTrainCell int) int { ///奖励球队训练位并返回最新训练位数量
	if addTrainCell <= 0 {
		return 0
	}
	self.TrainCount += addTrainCell
	return self.TrainCount
}

func (self *Team) AwardTicket(addTicket int) bool { ///奖励球队球票
	if addTicket <= 0 {
		return false
	}
	self.Ticket += addTicket
	self.Save()
	return true
}

func (self *Team) Create(accountID int, teamID int, syncMgr *SyncMgr) bool { ///加载球队信息
	if self.ID != 0 {
		return false ///避免重复加载GetProcessMgr()
	}
	query := fmt.Sprintf("select * from %s where accountid=%d limit 1", tableTeam, accountID)
	if teamID > 0 {
		query = fmt.Sprintf("select * from %s where id=%d limit 1", tableTeam, teamID)
	}
	ok := GetServer().GetDynamicDB().fetchOneRow(query, &self.TeamInfo)
	if false == ok {
		return false
	}

	if self.createStarSpy() == false { ///加载所属球探信息
		GetServer().GetLoger().Warn("Team Create createStarSpy fail! teamID:%d", self.ID)
		return false ///加载所属阵型失败//避免重复加载
	}
	if self.createStarCenter() == false { ///加载所属球员中心信息
		GetServer().GetLoger().Warn("Team Create createStarCenter fail! teamID:%d", self.ID)
		return false ///加载所属球员中心失败
	}
	if self.loadStarList() == false { ///加载所属球员信息
		GetServer().GetLoger().Warn("Team Create LoadStarList fail! accountID:%d", accountID)
		return false ///加载所属球员失败
	}

	//! cy
	if self.createGetPower() == false {
		GetServer().GetLoger().Warn("Team Create createGetPower fail! teamID:%d", self.ID)
		return false
	}

	if self.createOSActivity() == false {
		GetServer().GetLoger().Warn("Team Create createOSActivity fail! teamID:%d", self.ID)
		return false
	}

	// if self.createVipShop() == false { ///加载商城信息
	// 	GetServer().GetLoger().Warn("Team Create createvipshop fail! teamID:%d", self.ID) ///加载所属商城信息失败
	// }
	//if self.loadFormationList() == false { ///加载所属阵型信息
	//	GetServer().GetLoger().Warn("Team Create loadFormationList fail! accountID:%d", accountID)
	//	return false ///加载所属阵型失败//避免重复加载
	//}
	//self.processMgr = NewProcessMgr(self.ID)
	//if nil == self.processMgr {
	//	GetServer().GetLoger().Warn("Team Create NewProcessMgr fail! accountID:%d", accountID)
	//	return false ///加载所属处理中心失败
	//}
	//self.resetAttribMgr = NewResetAttribMgr(self.ID)
	//if nil == self.resetAttribMgr {
	//	GetServer().GetLoger().Warn("Team Create NewResetAttribMgr fail! accountID:%d", accountID)
	//	return false ///加载所属处理中心失败
	//}
	//self.itemMgr = NewItemMgr(self.ID)
	//if nil == self.itemMgr {
	//	GetServer().GetLoger().Warn("Team Create NewItemMgr fail! accountID:%d", accountID)
	//	return false ///加载所属道具管理器失败
	//}
	self.InitDataUpdater(tableTeam, &self.TeamInfo) ///创建球队信息更新器
	if self.InitGameMgrList(syncMgr) == false {     ///初始化球队游戏逻辑管理器组件列表
		GetServer().GetLoger().Warn("Team Create InitGameMgrList fail! accountID:%d", accountID)
		return false
	}

	self.client = syncMgr.client.GetElement()
	return true
}

func (self *Team) GetCurrentFormStarItemInfo() ItemInfoList { ///得到首发球员装备信息
	currentFormation := self.GetFormationMgr().GetFormation(self.FormationID)
	starIDList := currentFormation.GetStarIDList()
	itemInfoList := self.GetItemMgr().GetStarItemInfoList(starIDList)
	return itemInfoList
}

func (self *Team) GetTeamInfoMsg() *TeamInfoMsg { ///加载球队所属型阵列表
	msg := NewTeamInfoMsg()
	msg.TeamInfo = self.TeamInfo                  ///放入球队信息
	msg.TeamInfo.Operation = ""                   ///串置空
	msg.StarSpyInfo = *self.starSpy.GetInfoCopy() ///放入球探信息
	for _, star := range self.starList {          ///放入所有球星信息
		star.CalcScore() ///计算球员评分和二级属性
		// if star.IsMannaStar == 1 {
		// 	node := *star.GetInfo()
		// 	node.ID += 10000
		// 	msg.StarList = append(msg.StarList, node)
		// 	continue
		// }
		msg.StarList = append(msg.StarList, *star.GetInfo())
	}
	msg.FormationList = self.GetFormationMgr().GetFormationInfoList() ///得到阵形信息列表
	msg.EquipmentList = self.GetCurrentFormStarItemInfo()             ///得到首发球员装备信息列表
	self.CalcScore()
	msg.Score = self.Score
	return msg
}

func (self *Team) createStarCenter() bool { ///创建球队所属球员中心组件
	if self.starCenter != nil {
		return false
	}
	self.starCenter = new(StarCenter)
	ok := self.starCenter.Init(self.ID)
	return ok
}

func (self *Team) createStarSpy() bool { ///创建球队所属球探组件
	if self.starSpy != nil {
		return false
	}
	self.starSpy = new(StarSpy)
	ok := self.starSpy.Init(self.ID, self)
	return ok
}

// func (self *Team) createVipShop() bool {
// 	if self.vipshop != nil {
// 		return false
// 	}
// 	self.vipShop = new(VipShop)
// 	ok := self.vipShop.Init(self.ID)
// 	return ok
// }

//! cy
//! 创建得到体力对象
func (self *Team) createGetPower() bool {
	if self.PowerValue != nil {
		return false
	}

	var ok bool
	self.PowerValue, ok = NewGetPower(self.ID)

	return ok
}

//! 创建开服活动
func (self *Team) createOSActivity() bool {
	if self.OSActivityValue != nil {
		return false
	}

	var ok bool
	self.OSActivityValue, ok = NewOSActivity(self.ID)

	return ok
}

///奖励一个球星进球队,如果球队内已有此球员则尝试提升星级,如果星级低于现有球员则转化
///成经验进入经验池之中
func (self *Team) AwardStarEx(starType int, starCount int) bool {
	itemMgr := self.GetItemMgr()
	syncMgr := self.GetItemMgr().syncMgr
	star := self.GetStarFromType(starType)
	loger := GetServer().GetLoger()
	itemIDList := IntList{}
	if nil == star { ///如果球队中没有此球队直接给球员入球队
		if loger.CheckFail("self.IsStarFull() == false", self.IsStarFull() == false, self.IsStarFull(), false) {
			return false
		}
		starID := self.AwardStar(starType)
		//fmt.Println(starType)
		star = self.GetStar(starID)
		if starCount > 0 {
			star.SetStarCount(starCount)
		}
		syncMgr.syncAddStar(IntList{starID}) ///通知客户端得到一个新球员star
	} else if starCount > star.EvolveCount {
		starTypeInfo := star.GetTypeInfo()
		///给星卡
		awardStarCardCount := self.GetStarCardCount(starTypeInfo, star.EvolveCount)
		itemIDList = itemMgr.AwardItem(ItemStarCard, awardStarCardCount)
		///可提升星级
		star.SetStarCount(starCount)
		syncMgr.SyncObject("AwardStarEx", star)
	} else {
		starTypeInfo := star.GetTypeInfo()
		///给星卡
		awardStarCardCount := self.GetStarCardCount(starTypeInfo, star.EvolveCount)
		itemIDList = itemMgr.AwardItem(ItemStarCard, awardStarCardCount)

		//starType := star.GetTypeInfo()
		//firstValue, growValue := starType.GetFirstAndGrowValue()
		//exp := (firstValue + growValue*3) * (starCount * starCount)
		//newExp := self.StarExpPool + exp
		//if newExp >= 0 {
		//	self.StarExpPool = newExp
		//} else {
		//	self.StarExpPool = ExpPoolLimit
		//}
		//syncMgr.SyncObject("AwardStarEx", self) ///同步球队经验池属性值变更
	}
	if itemIDList.Len() > 0 {
		syncMgr.syncAddItem(itemIDList)
	}
	return true
}

///通用奖励流程
func (self *Team) AwardObject(awardType int, awardCount int, awardGrade int, awardStar int) bool {
	syncMgr := self.GetItemMgr().syncMgr
	itemMgr := self.GetItemMgr()
	///首先排除球员
	if awardStar != 0 {
		return self.AwardStarEx(awardStar, awardGrade)
	}
	switch awardType {
	case awardTypeCoin: ///球币
		self.AwardCoin(awardCount)
	case awardTypeTeamExp: ///球队经验池加经验
		self.StarExpPool += awardCount
	case awardTypeTactic: ///球队战术点
		self.TacticalPoint += awardCount
	case awardTypeTalent: ///球队培养点数
		//		self.TalentPoint += awardCount
		self.AwardCoin(awardCount)
	case awardTypeActionPoint: ///球队行动点数
		self.ActionPoint += awardCount
	case awardTypeTicket: ///球队钻石
		self.AwardTicket(awardCount)
	case awardTypeVipExp: //      = 400007 ///vip经验
		self.AwardVipExp(awardCount) ///充钻同时加vip经验
	case awardTypeManagerExp: //     = 400008 ///球员经理经验
		self.AwardExp(nil, awardCount)
	default:
		itemIDList := itemMgr.AwardItem(awardType, awardCount)
		if itemIDList == nil {
			return false
		}

		for i := range itemIDList {
			item := itemMgr.GetItem(itemIDList[i])
			if awardGrade > 0 {
				item.Color = awardGrade
			}
		}

		syncMgr.syncAddItem(itemIDList)
	}
	syncMgr.SyncObject("AwardObject", self) ///通知客户端球队信息属性变更
	return true
}

func (self *Team) awardNumberType(awardType int, awardCount int) {
	switch awardType {
	case awardTypeCoin: ///球币
		self.AwardCoin(awardCount)
	case awardTypeTeamExp: ///球队经验池加经验
		self.StarExpPool += awardCount
	case awardTypeTactic: ///球队战术点
		self.TacticalPoint += awardCount
	case awardTypeTalent: ///球队培养点数
		//		self.TalentPoint += awardCount
		self.AwardCoin(awardCount)
	case awardTypeActionPoint: ///球队行动点数
		self.ActionPoint += awardCount
	case awardTypeTicket: ///球队钻石
		self.AwardTicket(awardCount)
	case awardTypeManagerExp: ///球队经理经验
		self.AwardExp(nil, awardCount)
	case awardTypeVipExp: //!VIP经验
		self.AwardVipExp(awardCount)
	}

	return
}

func (self *Team) loadStarList() bool { ///加载球队所属球星列表
	if self.starList != nil {
		return false
	}
	self.starList = make(StarList)
	starInfo := new(StarInfo)
	starListQuery := fmt.Sprintf("select * from %s where teamid=%d", tableStar, self.ID)
	elmentList := GetServer().GetDynamicDB().fetchAllRows(starListQuery, starInfo)
	if nil == elmentList {
		return false
	}
	for i := range elmentList {
		starInfo = elmentList[i].(*StarInfo)
		self.starList[starInfo.ID] = NewStar(starInfo, self)
	}
	numStar := len(self.starList)
	return numStar > 0
}

func (self *Team) GetStarFromType(starType int) *Star { ///判断球队中的已雇佣球员数是否已达上限
	for _, v := range self.starList {
		if v.GetInfo().Type == starType {
			return v
		}
	}
	return nil
}

func (self *Team) GetTotalStarCount() int { ///得到当前球员所有球员总数
	currentStarCount := len(self.starList) ///得到当前球员数
	return currentStarCount
}

func (self *Team) IsStarFull() bool { ///判断球队中的已雇佣球员数是否已达上限
	const maxTeamStarCountIndex = 3
	staticDataMgr := GetServer().GetStaticDataMgr()
	currentStarCount := len(self.starList) ///得到当前球员数
	maxStarCount := staticDataMgr.GetConfigStaticDataInt(configTeam, configItemDefaultTeamParam, maxTeamStarCountIndex)
	maxStarCount += self.AddStarPos
	return currentStarCount >= maxStarCount
}

func (self *Team) GetStarCenter() *StarCenter {
	return self.starCenter
}

func (self *Team) GetStarSpy() *StarSpy { ///得到球探接口
	return self.starSpy
}

func (self *Team) GetMannaStarMgr() *MannaStarMgr { //!得到原创球员管理器
	return self.gameMgrList[mgrTypeMannaStarMgr].(*MannaStarMgr)
}

func (self *Team) GetAtlasMgr() *AtlasMgr { //! 得到图鉴管理器
	return self.gameMgrList[mgrTypeAtlasMgr].(*AtlasMgr)
}

// func (self *Team) GetVipShop() *VipShop { ///得到商城接口
// 	return self.vipShop
// }

func (self *Team) SpendTicket(payTicket int, reson string) bool { ///扣球队球票,reson为原因
	if self.Ticket < payTicket || payTicket <= 0 { ///余额不足或是扣负数均为失败
		return false
	}
	payTicket -= payTicket ///扣钱
	self.Save()
	return true
}

func (self *Team) AwardAttrib(attribType int, addAtrib int) bool { ///根据不同属性类型奖励属性值
	result := false
	switch attribType {
	case awardTypeCoin:
		result = self.AwardCoin(addAtrib)
	case awardTypeTicket:
		result = self.AwardTicket(addAtrib)
	}
	return result
}

func (self *Team) RedressPay() { ///玩家上线时补偿可能产生的充值时掉线情况
	redressPayQuery := fmt.Sprintf("select * from %s where teamid=%d and state=1 limit 1", tablePayOrder, self.ID)
	payOrderInfo := new(PayOrderInfo)
	GetServer().GetDynamicDB().fetchOneRow(redressPayQuery, payOrderInfo)
	if payOrderInfo.ID <= 0 {
		return ///没有需要补偿的充值记录
	}
	vipShopBuyDiamondMsg := new(VipShopBuyDiamondMsg)
	vipShopBuyDiamondMsg.PayOrderID = payOrderInfo.ID
	vipShopBuyDiamondMsg.MoneyID = payOrderInfo.ProductID
	vipShopBuyDiamondMsg.TeamID = self.ID
	vipShopBuyDiamondMsg.PayMoney = payOrderInfo.Money
	vipShopBuyDiamondMsg.processAction(self.client) ///转给服务器处理
}

func (self *Team) GetInfo() *TeamInfo { ///得到球队属性复本
	return &self.TeamInfo
}

func (self *Team) GetMainStarsList() StarSlice { ///得到球队首发球员列表
	starIDList := self.GetCurrentFormObject().GetStarIDList()
	starSlice := StarSlice{}
	for i := range starIDList {
		starID := starIDList[i]
		star := self.starList[starID]
		starSlice = append(starSlice, star)
	}
	return starSlice
}

func (self *Team) GetStartersList() StarInfoList { ///得到球队首发球员列表
	starIDList := self.GetCurrentFormObject().GetStarIDList()
	starInfoList := StarInfoList{}
	for i := range starIDList {
		starID := starIDList[i]
		star := self.starList[starID]
		starInfo := star.GetInfo()
		starInfoList = append(starInfoList, *starInfo)
	}
	return starInfoList
}

func (self *Team) Update(now int, client IClient) { ///球队自身更新状态
	//return ///暂时屏蔽
	self.starSpy.Update(now, client) ///球探更新自身状态
	//self.starCenter.Update(now, client)          ///球员中心更新自身状态
	self.GetProcessMgr().Update(now, client)     ///处理中心更新自身状态
	self.GetResetAttribMgr().Update(now, client) ///可重置值管理中心更新自身状态
	self.RestoreActionPoint(now)                 ///自动恢复行动点
	self.AutoSave(now)
	//self.GetVipShopMgr().Update(now, client)     ///商城更新自身状态
	//	GetServer().GetSDKMgr().ProcessQihoo360Pay()
}

func (self *Team) AutoSave(now int) { ///每隔五分钟自动存储
	const autoSaveTime = 300
	if now%autoSaveTime <= 0 {
		self.Save()
	}
}

func (self *Team) FindStarCount(fieldType int, color int, starCount int, starScore int) int { ///球队得到指定位置指定品质颜色球员个数
	result := 0
	for _, star := range self.starList {
		starType := star.GetTypeInfo()
		starInfo := star.GetInfo()
		if color > 0 && starType.Grade < color {
			continue
		}
		if starCount > 0 && starInfo.EvolveCount < starCount {
			continue
		}
		if starScore > 0 && starInfo.Score < starScore {
			continue
		}
		if star.CanKickFieldType(fieldType) == false {
			continue ///判断踢球场位是否满足要求
		}
		result++
	}
	return result
}

func (self *Team) IsStarsContractPointEnough() bool { ///判断球队所有上阵球员契约值是否均大于0
	result := true
	formation := self.GetCurrentFormObject()
	starIDList := formation.GetStarIDList()
	for starID := range starIDList {
		star := self.starList[starID]
		if nil == star {
			continue
		}
		starInfo := star.GetInfo()
		if starInfo.ContractPoint <= 0 {
			result = false
			break
		}
	}
	return result
}

func (self *Team) SpendStarsContractPoint(payContractPoint int) { ///将球队所有上阵球员契约值减指定值
	formation := self.GetCurrentFormObject()
	starIDList := formation.GetStarIDList()
	for i := range starIDList {
		starID := starIDList[i]
		star := self.starList[starID]
		if nil == star {
			continue
		}
		starInfo := star.GetInfo()
		if starInfo.ContractPoint >= 1000 {
			continue ///永久球员不扣契约值
		}
		if (starInfo.ContractPoint - payContractPoint) >= 0 {
			starInfo.ContractPoint -= payContractPoint
		}
	}
}

///计算比赛结果,返回x,y表示几比几x为玩家进球数,y为npc进球数
func (self *Team) CalcMatchResult(userTeam *Team) (int, int) {
	const overcomeAwardTurns = 2 ///克制奖励进攻次数
	///得到玩家球队比赛三围
	attackUserScoreCalc, defenseUserScoreCalc, organizeUserScoreCalc := self.CalcScore()
	///得到目标球队比赛三围
	attackTargetScoreCalc, defenseTargetScoreCalc, organizeTargetScoreCalc := userTeam.CalcScore()
	///计算玩家球队进攻次数
	attackUserTurns := int(organizeUserScoreCalc * 10 / (organizeUserScoreCalc + organizeTargetScoreCalc))
	///计算npc球队进攻次数
	attackTargetTurns := int(organizeTargetScoreCalc * 10 / (organizeUserScoreCalc + organizeTargetScoreCalc))
	///计算双方阵形相克对进攻次数的加成
	formOverNum, tacticOverNum := float32(0), float32(0) ///攻方阵型克制系数与战术克制系数
	formationUser := self.GetCurrentFormObject()
	targetFormation := userTeam.GetCurrentFormObject().GetInfo().Type
	if formationUser.IsOvercome(targetFormation) {
		attackUserTurns += overcomeAwardTurns
		formOverNum = -1 * (0.1 + float32(self.FormationLevel)*0.005)
	} else if formationUser.IsBeOvercome(targetFormation) {
		attackTargetTurns += overcomeAwardTurns
		formOverNum = 0.1 + float32(userTeam.FormationLevel)*0.005
	}
	tacticalTarget := userTeam.GetCurrentFormObject().CurrentTactic
	///战术克制
	if formationUser.IsOverTactic(tacticalTarget) {
		tacticOverNum = 0.1 + float32(self.FormationLevel)*0.005
	} else if formationUser.IsBeOverTactic(tacticalTarget) {
		tacticOverNum = -1 * (0.1 + float32(userTeam.FormationLevel)*0.005)
	}
	userGoalRate, targetGoalRate := 0, 0
	userGoalCount, targetGoalCount := 0, 0
	///计算玩家进球数
	//userGoalCount, userGoalRate = CalcGoalCount(attackUserTurns, attackUserScoreCalc, defenseTargetScoreCalc, tacticOverNum)
	///计算目标球队进球数
	//targetGoalCount, targetGoalRate = CalcGoalCount(attackTargetTurns, attackTargetScoreCalc, defenseUserScoreCalc, formOverNum)

	userGoalCount, userGoalRate, targetGoalCount, targetGoalRate, matchList := CalcGoalCountTurns(attackUserTurns, attackUserScoreCalc, defenseTargetScoreCalc, tacticOverNum,
		attackTargetTurns, attackTargetScoreCalc, defenseUserScoreCalc, formOverNum, self, targetFormation, 0, organizeUserScoreCalc+organizeTargetScoreCalc)

	//! 天天联赛添加潜规则
	isLost := true

	if self.VipLevel >= 1 && self.VipLevel <= 3 {
		if userTeam.Score < self.Score*75/100 {
			isLost = false
		}
	} else if self.VipLevel >= 4 && self.VipLevel <= 5 {
		if userTeam.Score < self.Score*85/100 {
			isLost = false
		}
	} else {
		if userTeam.Score < self.Score*95/100 {
			isLost = false
		}
	}
	//! 不会输
	if !isLost {
		if userGoalCount == targetGoalCount {
			userGoalCount += 1
			matchFlowInfo := NewAttackTurns(OurTeam, 1, 0)
			matchList = append(matchList, matchFlowInfo)
		} else if userGoalCount < targetGoalCount {
			tmp := userGoalCount
			userGoalCount = targetGoalCount
			targetGoalCount = tmp

			for i := 0; i < len(matchList); i++ {
				if matchList[i].Offensive == OurTeam {
					matchList[i].Offensive = EnemyTeam
				} else {
					matchList[i].Offensive = OurTeam
				}
			}

		}
	}

	client := GetServer().userMgr.GetClientByTeamID(self.ID)
	SendMatchFlowMsg(client, matchList)

	GetServer().GetLoger().Print("%s{攻%f 防%f 组%f 攻次%d 球率%d}---%s{攻%f 防%f 组%f 攻次%d 球率%d}\n",
		userTeam.GetInfo().Name, attackUserScoreCalc, defenseUserScoreCalc, organizeUserScoreCalc, attackUserTurns, userGoalRate,
		self.Name, attackTargetScoreCalc, defenseTargetScoreCalc, organizeTargetScoreCalc, attackTargetTurns, targetGoalRate)
	return userGoalCount, targetGoalCount
}

func (self *Team) GetAverageLevel() int {
	nSumLevel := 0
	nSumNum := 0
	for _, v := range self.starList {
		starInfo := v.GetInfo()
		nSumLevel += starInfo.Level
		nSumNum++
	}

	nAverageLevel := 0
	if 0 != nSumNum {
		nAverageLevel = nSumLevel / (nSumNum)
	}
	return nAverageLevel
}

func (self *Team) GetMinimumLevelStar(minimumLevel int) *Star { ///返回第一个小于平均等级的球员
	if nil == self.starList {
		return nil
	}

	///返回第一个小于平均等级的球星
	var randomStar *Star = nil
	randomStar = nil
	for _, v := range self.starList {
		starInfo := v.GetInfo()
		if starInfo.Level < minimumLevel {
			return v
		}

		if randomStar == nil {
			randomStar = v
		}
	}

	///若都不小于平均等级,则返回随机球星
	return randomStar
}

//VIP功能涉及
func (self *Team) AwardVipExp(vipExp int) bool {
	if vipExp <= 0 {
		return false
	}

	//获取当前VIP经验上限
	//	staticDataMgr := GetServer().GetStaticDataMgr()
	//	expUpperLimit := staticDataMgr.GetConfigStaticDataInt(configVip, configVipData, 2)
	//	if self.VipExp >= expUpperLimit {
	//		return false
	//	}
	staticDataMgr := GetServer().GetStaticDataMgr()
	vipNeedExp := staticDataMgr.GetVipNeedExp(self.VipLevel + 1)
	if vipNeedExp <= 0 {
		self.VipExp = 0 ///满级手经验置0
		return false    ///满级
	}

	self.VipExp += vipExp
	//	if self.VipExp >= expUpperLimit {
	//		self.VipExp = expUpperLimit
	//	}

	self.UpVipLevel()
	self.Save()
	return true
}

func (self *Team) UpVipLevel() {
	//获取VIP等级上限
	staticDataMgr := GetServer().GetStaticDataMgr()
	//	levelUpperLimit := staticDataMgr.GetConfigStaticDataInt(configVip, configVipData, 1)
	//	if self.VipLevel >= levelUpperLimit {
	//		self.VipLevel = levelUpperLimit
	//	}
	vipNeedExp := staticDataMgr.GetVipNeedExp(self.VipLevel + 1)
	if vipNeedExp <= 0 {
		self.VipExp = 0 ///满级手经验置0
		return          ///满级
	}

	vipOldLevel := self.VipLevel
	for i := 0; i < 100; i++ {
		vipNeedExp = staticDataMgr.GetVipNeedExp(self.VipLevel + 1)
		if 0 == vipNeedExp {
			//GetServer().GetLoger().Warn("VIPTypeInfo is nil teamid:%d", self.GetID())
			break
		}
		if self.VipExp < vipNeedExp {
			break ///经验不足升级
		}

		//		if self.VipLevel > levelUpperLimit {
		//			self.VipLevel = levelUpperLimit
		//			break ///VIP等级最大
		//		}

		self.VipLevel++
		//self.VipExp -= vipNeedExp ///扣除已升级消耗经验值
	}

	if vipOldLevel != self.VipLevel {
		syncMgr := self.GetMailMgr().syncMgr
		//		vipInfo := GetServer().GetStaticDataMgr().GetVipInfo(self.VipLevel)
		//		if vipInfo == nil {
		//			GetServer().GetLoger().Warn("vipInfo == nil teamid: %d", self.GetID())
		//		}

		//self.StoreCapacity = vipInfo.Param9
		syncMgr.SyncObject("viplevelup", self)
	}
}

func (self *Team) GetVipLevel() int {
	return self.VipLevel
}

func (self *Team) GetAllStarList() IntList {
	allStarList := IntList{}
	for k, _ := range self.starList {
		allStarList = append(allStarList, k)
	}
	return allStarList
}

func (self *Team) GetStarCardCount(starType *StarTypeStaticData, starCount int) int {
	// 		公式：读取startype表ticket字段的值为基础星卡数
	// 			= Int（（1+(球员星级-1)*(球员星级-2)/4）*ticket/100）
	/// int awardStarCardCount=(1+(pStarInfo->evolvecount-1)*(pStarInfo->evolvecount-2)/4.0f)*pStarInfo->pCsvData->ticket/100;
	//awardStarCardCount := (1 + (float32(starInfo.EvolveCount)-1)*(float32(starInfo.EvolveCount)-2)/float32(4)) * float32(starType.Ticket) / 100
	//starCardCount := int(awardStarCardCount)

	awardStarCardCount := (1 + (float32(starCount)-1)*(float32(starCount)-2)/float32(4)) * float32(starType.Ticket) / 100
	starCardCount := int(awardStarCardCount)
	return starCardCount
}

func (self *Team) GetMannaStarCardCount(starType *MannaStar, starCount int) int {
	// 		公式：读取startype表ticket字段的值为基础星卡数
	// 			= Int（（1+(球员星级-1)*(球员星级-2)/4）*ticket/100）
	/// int awardStarCardCount=(1+(pStarInfo->evolvecount-1)*(pStarInfo->evolvecount-2)/4.0f)*pStarInfo->pCsvData->ticket/100;
	//awardStarCardCount := (1 + (float32(starInfo.EvolveCount)-1)*(float32(starInfo.EvolveCount)-2)/float32(4)) * float32(starType.Ticket) / 100
	//starCardCount := int(awardStarCardCount)

	awardStarCardCount := (1 + (float32(starCount)-1)*(float32(starCount)-2)/float32(4)) * float32(starType.Ticket) / 100
	starCardCount := int(awardStarCardCount)
	return starCardCount
}

func (self *Team) GetVipStarSackRepayRate() (int, int) { ///得到球队vip解雇球员训练点返还率和经验值返还率
	staticDataMgr := GetServer().GetStaticDataMgr()
	sackRepayTalentRate := staticDataMgr.GetConfigStaticDataInt(configStar,
		configItemStarCommonConfig, 5) ///vip0普通解雇球员返还率
	sackRepayExpRate := sackRepayTalentRate
	vipTypeInfo := staticDataMgr.GetVipInfo(self.VipLevel)
	if vipTypeInfo != nil {
		sackRepayTalentRate = vipTypeInfo.Param12
		sackRepayExpRate = vipTypeInfo.Param13
	}
	return sackRepayTalentRate, sackRepayExpRate
}

func (self *Team) GetVipFreeVolunteerUpdateCount() int { ///得到球队vip解雇球员训练点返还率和经验值返还率
	staticDataMgr := GetServer().GetStaticDataMgr()
	///先取vip0的球星来投免费刷新次数
	maxUpdateCount := staticDataMgr.GetConfigStaticDataInt(configStarCenter, configItemStarVolunteer, 3)
	vipTypeInfo := staticDataMgr.GetVipInfo(self.VipLevel)
	if vipTypeInfo != nil {
		maxUpdateCount = vipTypeInfo.Param11
	}
	return maxUpdateCount
}

func (self *Team) OnEnterMap() { ///进大地图事件

	mannaStarMgr := self.GetMannaStarMgr()
	starList := MannaStarSlice{}
	for i := MannaStarSeatOne; i <= MannaStarSeatThree; i++ {
		star := mannaStarMgr.GetMannaStarFromSeat(i)
		if star == nil {
			continue
		}

		node := star.MannaStarType

		//! 获取自创球员信息
		starList = append(starList, &node)
	}

	mannaStarMsg := new(QueryMannaStarResultMsg)
	mannaStarMsg.StarList = starList
	self.client.SendMsg(mannaStarMsg)
	self.client.UpdateSeqID()

	msg := self.GetTeamInfoMsg()
	msg.AccessToken = self.client.accessToken
	msg.SDKUserID = self.client.userID
	self.client.SendMsg(msg)

	// starList := MannaStarSlice{}
	// for i := MannaStarSeatOne; i <= MannaStarSeatThree; i++ {
	// 	star := mannaStarMgr.GetMannaStarFromSeat(i)
	// 	if star == nil {
	// 		continue
	// 	}

	// 	//! 获取自创球员信息
	// 	starList = append(starList, &star.MannaStarType)
	// }

	// if len(starList) == 0 {
	// 	//! 创建默认自创球员

	// 	starList = append(starList, mannaStarMgr.GetTemplate(self.GetID(), 1))
	// 	starList = append(starList, mannaStarMgr.GetTemplate(self.GetID(), 2))
	// 	starList = append(starList, mannaStarMgr.GetTemplate(self.GetID(), 3))
	// }

	// mannaStarMsg := new(QueryMannaStarResultMsg)
	// mannaStarMsg.StarList = starList
	// self.client.SendMsg(mannaStarMsg)
}

const (
	createthreehour = 3 * int(time.Hour/time.Second)  //建号3小时的时间间隔秒数
	oneday          = 24 * int(time.Hour/time.Second) //一天的时间间隔秒数
	logintwoday     = 1 * oneday                      //登录2天
	loginfiveday    = 4 * oneday                      //登录5天
	logintenday     = 9 * oneday                      //登录10天
	loginfifteenday = 14 * oneday                     //登录15天
	loginthirtyday  = 29 * oneday                     //登录30天

	InitpayBit     = 1  //首次充值掩码位
	Buy1980Bit     = 2  //购买1980钻石掩码位
	NewBit         = 3  //! 新人福利
	CthBit         = 4  //建号3小时掩码位
	LtwoBit        = 5  //登录2天掩码位
	LfiveBit       = 6  //登录5天掩码位
	LtenBit        = 7  //登录10天掩码位
	LfifteenBit    = 8  //登录15天掩码位
	InitpayGetBit  = 9  //首次充值领取标识掩码位
	Buy1980GetBit  = 10 //购买1980钻石领取标识掩码位
	NewGetBit      = 11 //! 新人福利掩码位
	CthGetBit      = 12 //建号3小时领取标识掩码位
	LtwoGetBit     = 13 //登录2天领取标识掩码位
	LfiveGetBit    = 14 //登录5天领取标识掩码位
	LtenGetBit     = 15 //登录10天领取标识掩码位
	LfifteenGetBit = 16 //登录15天领取标识掩码位

	awardmoneyid = 4 //1980钻石商品id

	InitpayAward  = 5074 //首次充值励球员id
	Buy1980Award  = 3043 //购买1980钻石励球员id
	NewAward      = 4165
	CthAward      = 3084 //建号3小时奖励球员id
	LtwoAward     = 3126 //登录2天励球员id
	LfiveAward    = 3071 //登录5天励球员id
	LtenAward     = 5078 //登录10天励球员id
	LfifteenAward = 4191 //登录15天励球员id

	AwardGetSuccess  = 1  //领取成功
	CannotGetAward   = 0  //不能领取
	AllreadyGetAward = -1 //已领取
	CannotGetFull    = -2 //球队满

)

//登陆活动检测方法
func (self *Team) CheckLoginAward(logintimeUTC int) {
	//取得与当前的时间间隔秒数，都是UTC时间
	durationsec := logintimeUTC - self.MakeTime                                                                   ///这是与当前utc时间真实的差值
	makeTime := time.Unix(int64(self.MakeTime), 0)                                                                ///用创建的utc时间妙取得创建时间对象
	trimMakeTime := time.Date(makeTime.Year(), makeTime.Month(), makeTime.Day(), 0, 0, 0, 0, makeTime.Location()) ///把天内的时间去掉
	//durationWholeDaysec := durationsec + self.MakeTime%86400 ///加上24小时内的模值就是当前时间与注册天0时的差值,这里少加了8小时
	durationWholeDaysec := logintimeUTC - int(trimMakeTime.Unix()) ///当前时间与注册天0时的差值
	//掩码位
	var bit int = 0
	//判断奖励
	if durationWholeDaysec >= loginfifteenday {
		bit = LfifteenBit
	} else if durationWholeDaysec >= logintenday {
		bit = LtenBit
	} else if durationWholeDaysec >= loginfiveday {
		bit = LfiveBit
	} else if durationWholeDaysec >= logintwoday {
		bit = LtwoBit
	} else if durationsec >= createthreehour {
		bit = CthBit
	} else {
		bit = NewBit
	}
	//判断获得奖励
	//needrecord := false //是否需要写数据库
	for i := 3; i <= bit; i++ {
		if !TestMask64(self.LoginAndPayAward, i) { //先测试有没有置位，主要是为了避免无谓的写库操作
			self.LoginAndPayAward = SetMask64(self.LoginAndPayAward, i, 1) //置位
			//needrecord = true
		}
	}
	//if needrecord { //写入数据库
	//	updateawardsql := fmt.Sprintf("update %s set loginandpayaward=%d where `id`=%d", tableTeam, self.LoginAndPayAward, self.ID)
	//	_, rowsStarAffected := GetServer().GetDynamicDB().Exec(updateawardsql)
	//	if !(rowsStarAffected > 0) {
	//		GetServer().GetLoger().Warn("update %s loginandpayaward error, teamid = %d", tableTeam, self.ID)
	//	}
	//}
}

////检测3小时创建奖励，单独拿出来是为了在请求消息到达时做一次判断，其它奖励就只登陆判断
////返回结果为达到领取奖励还差的秒数
//func (self *Team) CheckCreatThreeHour(querytimeUTC int) int {
//	if !TestMask64(self.LoginAndPayAward, CthBit) { //先测试3小时奖励有没有置位，没有才进行后续比较等操作
//		//取得与当前的时间间隔秒数，都是UTC时间
//		durationsec := querytimeUTC - self.MakeTime
//		//转换为纳秒
//		durationnsec := int64(durationsec) * int64(time.Second)
//		if durationnsec >= createthreehour {
//			self.LoginAndPayAward = SetMask64(self.LoginAndPayAward, CthBit, 1) //置位
//			updateawardsql := fmt.Sprintf("update %s set loginandpayaward=%d where `id`=%d", tableTeam, self.LoginAndPayAward, self.ID)
//			_, rowsStarAffected := GetServer().GetDynamicDB().Exec(updateawardsql)
//			if !(rowsStarAffected > 0) {
//				GetServer().GetLoger().Warn("update %s loginandpayaward error, teamid = %d", tableTeam, self.ID)
//			}
//			return 0
//		} else {
//			return int((createthreehour - durationnsec) / int64(time.Second))
//		}
//	} else {
//		return 0
//	}
//}

//充值活动检测方法
func (self *Team) CheckPayAward(moneyid int) {
	//先测试是否充过值
	//needrecord := false //是否需要写数据库
	if !TestMask64(self.LoginAndPayAward, InitpayBit) {
		self.LoginAndPayAward = SetMask64(self.LoginAndPayAward, InitpayBit, 1) //充值了但没置位过就要置位
		//needrecord = true
	}
	//测试是否充过1980钻石
	if !TestMask64(self.LoginAndPayAward, Buy1980Bit) {
		if moneyid == awardmoneyid { //判断是不是充的指定卡
			self.LoginAndPayAward = SetMask64(self.LoginAndPayAward, Buy1980Bit, 1) //充值了但没置位过就要置位
			//needrecord = true
		}
	}
	//if needrecord { //写入数据库
	//	updateawardsql := fmt.Sprintf("update %s set loginandpayaward=%d where `id`=%d", tableTeam, self.LoginAndPayAward, self.ID)
	//	_, rowsStarAffected := GetServer().GetDynamicDB().Exec(updateawardsql)
	//	if !(rowsStarAffected > 0) {
	//		GetServer().GetLoger().Warn("update %s loginandpayaward error, teamid = %d", tableTeam, self.ID)
	//	}
	//}
}

//取得奖励状态
func (self *Team) GetLoginAndPayAwardState(awardbit int, querytimeUTC int) (awardstate int, stillneed int) {
	//检测参数
	if awardbit <= 0 || awardbit > LfifteenBit {
		awardstate = 0
		stillneed = 0
		return
	}
	var GetBit int = 0   //领取标识掩码位
	var needterm int = 0 //领取条件
	switch awardbit {
	case CthBit:
		GetBit = CthGetBit
		needterm = createthreehour
	case LtwoBit:
		GetBit = LtwoGetBit
		needterm = logintwoday
	case LfiveBit:
		GetBit = LfiveGetBit
		needterm = loginfiveday
	case LtenBit:
		GetBit = LtenGetBit
		needterm = logintenday
	case LfifteenBit:
		GetBit = LfifteenGetBit
		needterm = loginfifteenday
	case NewBit:
		GetBit = NewGetBit
		needterm = 0
	case InitpayBit:
		GetBit = InitpayGetBit
	case Buy1980Bit:
		GetBit = Buy1980GetBit
	}
	achievestate := TestMask64(self.LoginAndPayAward, awardbit) //判断能否领取
	if achievestate {                                           //判断能否领取
		fetchstate := TestMask64(self.LoginAndPayAward, GetBit) //判断是否已领
		if fetchstate {
			awardstate = 2
		} else {
			awardstate = 1
		}
		stillneed = 0
	} else {
		awardstate = 0
		//检查差距
		//取得与当前的时间间隔秒数，都是UTC时间
		durationsec := querytimeUTC - self.MakeTime
		if awardbit != CthBit {
			makeTime := time.Unix(int64(self.MakeTime), 0)                                                                ///用创建的utc时间妙取得创建时间对象
			trimMakeTime := time.Date(makeTime.Year(), makeTime.Month(), makeTime.Day(), 0, 0, 0, 0, makeTime.Location()) ///把天内的时间去掉
			durationsec = querytimeUTC - int(trimMakeTime.Unix())                                                         ///当前时间与注册天0时的差值
		}
		stillneed = Max(0, needterm-durationsec) //这是秒数
	}
	return
}

//领取奖励，领取成功返回AwardGetSuccess，领取失败返回CannotGetAward，AllreadyGetAward，CannotGetFull
func (self *Team) GetLoginAndPayAward(awardbit int) int {
	//检测参数
	if awardbit <= 0 || awardbit > LfifteenBit {
		return CannotGetAward
	}
	var awardid int = 0 //奖励编号
	var GetBit int = 0  //领取标识掩码位
	switch awardbit {
	case CthBit:
		awardid = CthAward
		GetBit = CthGetBit
	case LtwoBit:
		awardid = LtwoAward
		GetBit = LtwoGetBit
	case LfiveBit:
		awardid = LfiveAward
		GetBit = LfiveGetBit
	case LtenBit:
		awardid = LtenAward
		GetBit = LtenGetBit
	case LfifteenBit:
		awardid = LfifteenAward
		GetBit = LfifteenGetBit
	case NewBit:
		awardid = NewAward
		GetBit = NewGetBit
	case InitpayBit:
		awardid = InitpayAward
		GetBit = InitpayGetBit
	case Buy1980Bit:
		awardid = Buy1980Award
		GetBit = Buy1980GetBit
	}

	if TestMask64(self.LoginAndPayAward, awardbit) { //判断能否领取
		if awardbit == NewBit {
			levelMgr := self.GetLevelMgr()
			if nil == levelMgr.FindLevel(1102) {
				GetServer().GetLoger().Warn("Level not pass ------  GetLoginAndPayAward")
				return CannotGetFull
			}
		}
		if !TestMask64(self.LoginAndPayAward, GetBit) { //判断是否已领取
			awardresult := self.AwardObject(0, 0, 1, awardid) //获得奖励
			if awardresult {                                  //如果获得成功就的修改对应数据并存库
				self.LoginAndPayAward = SetMask64(self.LoginAndPayAward, GetBit, 1)
				//updateawardsql := fmt.Sprintf("update %s set loginandpayaward=%d where `id`=%d", tableTeam, self.LoginAndPayAward, self.ID)
				//_, rowsStarAffected := GetServer().GetDynamicDB().Exec(updateawardsql)
				//if !(rowsStarAffected > 0) {
				//	GetServer().GetLoger().Warn("update %s loginandpayaward error, teamid = %d", tableTeam, self.ID)
				//}
				return AwardGetSuccess
			} else { //球队满
				return CannotGetFull
			}

		} else { //已领取
			return AllreadyGetAward
		}

	} else { //不能领取
		return CannotGetAward
	}

}

//登陆时修复球队创建时间
