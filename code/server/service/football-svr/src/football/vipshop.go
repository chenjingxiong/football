package football

import (
	"fmt"
	"reflect"
	"time"
)

const (
	VipShopType     = 1 //商城
	DayTaskBuyGoods = 2 ///日常任务用到的购买物品条件.
)

const (
	TimeLimit    = 1 //有时限类型
	TimeInfinite = 2 //无时限类型
)

const (
	FirstVer = 1 //初始版本号
)

const (
	MonthCardMoneyID = 7 //月卡在money表中的id号
)

type MoneyType struct { // st_money
	ID       int `json:"id"`       ///记录id
	Type     int `json:"type"`     ///1为标准充值   2为月卡
	Icon     int `json:"icon"`     ///图标id
	Diamonds int `json:"diamonds"` ///type为1时，为充值给与的标准钻石数。type为2时，为月卡每日领取的钻石数
	Money    int `json:"money"`    ///充值所需人民币金额，以元为单位
	Num      int `json:"num"`      ///超额赠送的限定次数
	NumGive  int `json:"numgive"`  ///超额赠送的钻石金额
	Give     int `json:"give"`     ///标准赠送的钻石金额
}

func (self *MoneyType) IsMonthCard() bool { ///判断此money是否是月卡
	isMonthCard := 2 == self.Type
	return isMonthCard
}

type VipShopInfo struct { // dy_vipshop
	ID          int `json:"id"`          //记录id
	Teamid      int `json:"teamid"`      //球队id
	Commodityid int `json:"commodityid"` //物品记录id
	Limittimes  int `json:"limittimes"`  //当前剩下的限购次数
	Cycle       int `json:"cycle"`       //记录当前商品循环次数
}

type VipShop struct { // dy_vipshop
	goodsinfo VipShopInfo
	// goodsList  LimitList
	DataUpdater
}

func (self *VipShop) GetReflectValue() reflect.Value { ///得到反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

func (self *VipShop) GetInfo() *VipShopInfo {
	if self == nil {
		return nil
	}
	return &self.goodsinfo
}

type VipShopStaticData struct {
	ID        int //记录id
	ItemType  int //商品类型
	Type1     int //商店类型: 1.商城 2.(暂无)
	Type2     int //商店分类: 1.限时促销 2.日常商品
	Class     int //道具分类:1.(待定)
	Num       int //商品包含道具数量
	Discount  int //折扣
	MoneyID   int //货币id
	MoneyNum  int //支付货币数量
	StartTime int //上架时间 (例:12251200 12月25日12点0分)
	Duration  int //持续时间 0为永久 单位: 秒
	CD        int //下架冷却时间 0为永久 单位: 秒
	Limit     int //限购次数
}

type LimitList map[int]*VipShop

// type IVipShop interface {
// 	GetStampFromStaticData(startTime int) int64    //转换静态数据到时间戳格式
// 	CheckShoppingTime(commodityID int) bool        //检测是否为购买期内
// 	PutCommodityUp(commodityList IntList)          //上架商品
// 	PullCommodityDown(commodityList IntList)       //下架商品
// 	GetCommodityInfo(commodityID int) *VipShopInfo //获取商品信息
// 	GetVersion() int                               //获取版本号
// 	SetVersion(currentVer int)                     //设置版本号
// 	GetGoogsList() *LimitList                      //得到商品信息链表
// }

//type IVipShopMgr interface {
//	IGameMgr
//	Update(now int, client IClient) ///处理中心更新自己状态
//	GetVersion() int
//	SetVersion(version int)
//	GetCommodityInfo(commodityID int) *VipShopInfo
//	CheckShoppingTime(commodityID int) (bool, int)
//	GetStampFromStaticData(startTime int) int64
//}

type VipShopMgr struct {
	GameMgr
	goodsList  LimitList //商品信息
	currentVer int       //版本号

}

func (self *VipShopMgr) SaveInfo() { ///保存数据
	for _, v := range self.goodsList {
		v.Save()
	}
}

func (self *VipShopMgr) GetType() int { ///得到管理器类型
	return mgrTypeVipShopMgr ///任务管理器
}

func NewCommodity(vipShopInfo *VipShopInfo) *VipShop {
	vipShop := new(VipShop)
	vipShop.goodsinfo = *vipShopInfo
	vipShop.InitDataUpdater(tableVipShop, &vipShop.goodsinfo)
	return vipShop
}

func NewVipShopMgr(teamID int) IGameMgr {
	vipShopMgr := new(VipShopMgr)
	if vipShopMgr.Init(teamID) == false {
		return nil
	}
	return vipShopMgr
}

// func (self *VipShopMgr) HasCommodity(commodityID int) bool {
// 	for _, v := range self.goodsList {
// 		commodityInfo := v.GetInfo()
// 		if commodityInfo.Commodityid == commodityID {
// 			return true
// 		}
// 	}
// 	return false
// }

func (self *VipShopMgr) GetCommodityInfo(commodityID int) *VipShopInfo {
	if nil != self.goodsList[commodityID] {
		return self.goodsList[commodityID].GetInfo()
	}

	// 若动态库不存在该商品信息(可能动态增加了新的商品) 查询静态库
	vipShopType := GetServer().GetStaticDataMgr().Unsafe().GetVipShopItemInfo(commodityID)
	if nil == vipShopType {
		return nil
	}

	if self.goodsList == nil {
		self.goodsList = make(LimitList)
	}

	vipshopInfo := new(VipShopInfo)
	vipshopInfo.Commodityid = vipShopType.ID
	vipshopInfo.Limittimes = vipShopType.Limit
	vipshopInfo.Teamid = self.team.GetID()
	vipshopInfo.Cycle = 0
	self.goodsList[commodityID] = NewCommodity(vipshopInfo)
	self.InsertDefault(vipshopInfo)
	return vipshopInfo
}

func (self *VipShopMgr) Init(teamID int) bool { //初始化商城信息
	vipShopQuery := fmt.Sprintf("select * from dy_vipshop where teamid = %d limit 1000", teamID)
	vipShopInfo := new(VipShopInfo)
	vipShopList := GetServer().GetDynamicDB().fetchAllRows(vipShopQuery, vipShopInfo)
	if nil == vipShopList {
		return false
	}

	//取得必要信息,存入内存
	for i := range vipShopList {
		vipShopInfo = vipShopList[i].(*VipShopInfo)
		recordID := vipShopInfo.Commodityid
		if nil == self.goodsList {
			self.goodsList = make(LimitList)
		}
		self.goodsList[recordID] = NewCommodity(vipShopInfo)
	}

	//设置初始版本号
	self.currentVer = FirstVer

	if len(self.goodsList) > 0 {
		return true
	}

	//若数据库无信息,则初始化玩家动态数据
	// staticDataMgr := GetServer().GetStaticDataMgr()
	// commodityList := staticDataMgr.GetVipShopStaticDataMap()
	// for _, v := range commodityList {
	// 	commodityType := v
	// 	if commodityType.Type2 == TimeLimit {
	// 		vipShopInfo.Commodityid = commodityType.ID
	// 		vipShopInfo.Limittimes = commodityType.Limit
	// 		vipShopInfo.Teamid = teamID
	// 		if nil == self.goodsList {
	// 			self.goodsList = make(LimitList)
	// 		}
	// 		self.goodsList[commodityType.ID] = NewCommodity(vipShopInfo)
	// 		self.InsertDefault(vipShopInfo)
	// 	}
	// }

	return true
}

func (self *VipShopMgr) InsertDefault(vipShopInfo *VipShopInfo) { //插入数据

	insertVipShopInfoSql := fmt.Sprintf("insert into %s (teamid,commodityid,limittimes,cycle) value (%d,%d,%d,%d)",
		tableVipShop, vipShopInfo.Teamid, vipShopInfo.Commodityid, vipShopInfo.Limittimes, vipShopInfo.Cycle) ///组插入记录SQL

	lastInsertItemID, _ := GetServer().GetDynamicDB().Exec(insertVipShopInfoSql)
	if lastInsertItemID <= 0 {
		GetServer().GetLoger().Warn("VipShop InsertDefault fail! teamid:%d", vipShopInfo.Teamid)
	}
}

// func (self *VipShopMgr) updateLimitTimes(now int, client IClient) { //更新商品信息
// 	const checkSec = 60
// 	if now%checkSec != 0 {
// 		return //每60秒检查一次
// 	}

// 	//检查当前时间是否为商品动态下架时间
// 	removeCommodityIDList := IntList{}
// 	for i := range self.goodsList {
// 		vipShopInfo := self.goodsList[i].GetInfo()
// 		//fmt.Printf("pulltime = %d and now = %d \r\n", vipShopInfo.PullTime, now)
// 		if now >= vipShopInfo.PullTime && 0 != vipShopInfo.PullTime &&
// 			vipShopInfo.PullTime > vipShopInfo.PutTime {
// 			removeCommodityIDList = append(removeCommodityIDList, vipShopInfo.Commodityid)
// 		}
// 	}

// 	if removeCommodityIDList.Len() > 0 {
// 		//同步客户端
// 		self.PullCommodityDown(removeCommodityIDList)
// 		vipShopChangeMsg := new(VipShopCommodityChangeMsg)
// 		vipShopChangeMsg.CommodityIDList = removeCommodityIDList
// 		vipShopChangeMsg.ChangeType = false
// 		client.SendMsg(vipShopChangeMsg)
// 		self.currentVer += 1
// 	}

// 	//检查当前时间是否为商品动态上架时间
// 	addCommodityIDList := IntList{}
// 	for i := range self.goodsList {
// 		vipShopInfo := self.goodsList[i].GetInfo()
// 		//fmt.Printf("puttime = %d and now = %d \r\n", vipShopInfo.PutTime, now)
// 		if now == vipShopInfo.PutTime {
// 			addCommodityIDList = append(addCommodityIDList, vipShopInfo.Commodityid)
// 		}
// 	}

// 	if addCommodityIDList.Len() > 0 {
// 		//同步客户端
// 		self.PutCommodityUp(addCommodityIDList)
// 		vipShopChangeMsg := new(VipShopCommodityChangeMsg)
// 		vipShopChangeMsg.CommodityIDList = addCommodityIDList
// 		vipShopChangeMsg.ChangeType = true
// 		client.SendMsg(vipShopChangeMsg)
// 		self.currentVer += 1
// 	}

// }

func (self *VipShopMgr) Update(now int, client IClient) { //商城更新信息
	//	self.updateLimitTimes(now, client)
}

func (self *VipShopMgr) GetStampFromStaticData(startTime int) int64 {
	startYear := time.Now().Year() //默认年份为当前年份

	startTimeList := IntList{}
	var tempValue = 0

	//计算具体年月日时分
	for i := 0; i < 8; i++ {
		if (i+1)%2 == 0 {
			tempValue = startTime%10*10 + tempValue
			startTimeList = append(startTimeList, tempValue)
			tempValue = 0
		} else {
			tempValue += startTime % 10
		}
		startTime /= 10
	}

	if startTimeList.Len() <= 0 {
		return 0
	}

	startMin := startTimeList[0]
	startHour := startTimeList[1]
	startDay := startTimeList[2]
	startMon := startTimeList[3]

	startTimeStamp := time.Date(startYear, time.Month(startMon), startDay, startHour, startMin, 0, 0, time.Local).Unix() //整分刷新
	return startTimeStamp
}

func (self *VipShopMgr) CheckShoppingTime(commodityID int) (bool, int) {
	//取得该商品静态信息
	vipShopType := GetServer().GetStaticDataMgr().GetVipShopItemInfo(commodityID)
	if TimeInfinite == vipShopType.Type2 {
		return true, 0 //若该商品为永久在架,则直接返回可以购买
	}

	if 0 == vipShopType.Duration {
		return true, 0 //物品不存在下架可能
	}

	// goodsInfo := self.goodsList[commodityID]
	// if nil == goodsInfo {
	// 	return true
	// }

	nowTime := Now()                                                     //当前时间
	startTime := int(self.GetStampFromStaticData(vipShopType.StartTime)) //开始时间
	endTime := startTime + vipShopType.Duration

	if nowTime < startTime {
		return false, 0 /// 商品未上架
	}

	if 0 == vipShopType.CD && nowTime >= endTime {
		return false, 0 ///商品已下架
	}

	curTime := (nowTime - startTime) % (vipShopType.CD + vipShopType.Duration) //当前轮回中时间
	curturn := (nowTime - startTime) / (vipShopType.CD + vipShopType.Duration) //当前轮回次数
	if curTime < vipShopType.Duration {
		return true, curturn
	}

	// startTime := goodsInfo.GetInfo().PutTime
	// endTime := goodsInfo.GetInfo().PullTime
	// currentTime := Now()

	// if currentTime >= startTime && currentTime <= endTime {
	// 	return true //在贩卖时间内则返回可以购买
	// }

	return false, 0
}

// func (self *VipShopMgr) PutCommodityUp(commodityList IntList) { //上架商品
// 	for _, i := range commodityList {
// 		vipShopInfo := self.goodsList[i].GetInfo()
// 		if nil == vipShopInfo {
// 			return
// 		}

// 		vipShopType := GetServer().GetStaticDataMgr().Unsafe().GetVipShopItemInfo(vipShopInfo.Commodityid)
// 		vipShopInfo.PutTime = Now()
// 		if 0 == vipShopType.Duration {
// 			vipShopInfo.PullTime = 0
// 		} else {
// 			vipShopInfo.PullTime = vipShopInfo.PutTime + vipShopType.Duration //设置下架时间
// 		}

// 	}
// }

// func (self *VipShopMgr) PullCommodityDown(commodityList IntList) { //下架商品
// 	for _, i := range commodityList {
// 		vipShopInfo := self.goodsList[i].GetInfo()
// 		if nil == vipShopInfo {
// 			return
// 		}

// 		vipShopType := GetServer().GetStaticDataMgr().Unsafe().GetVipShopItemInfo(vipShopInfo.Commodityid)
// 		vipShopInfo.PullTime = Now()
// 		if 0 == vipShopType.CD {
// 			vipShopInfo.PutTime = 0
// 		} else {
// 			vipShopInfo.PutTime = vipShopInfo.PullTime + vipShopType.CD
// 		}

// 	}
// }

func (self *VipShopMgr) GetVersion() int { //获取版本号
	return self.currentVer
}

func (self *VipShopMgr) SetVersion(currentVer int) { //设置版本号
	self.currentVer = currentVer
}
