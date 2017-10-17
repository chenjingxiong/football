package football

///道具消息处理器
func (self *ItemHandler) getName() string { ///返回可处理的消息类型
	return "item"
}

func (self *ItemHandler) initHandler() { ///初始化消息处理器
	self.addActionToList(new(QueryItemListMsg))
	self.addActionToList(new(EquipItemMsg))
	self.addActionToList(new(MergeItemMsg))
	self.addActionToList(new(ItemEvolveMsg))
}

type ItemHandler struct {
	MsgHandler
}

type MergeItemMsg struct { ///道具融合消息
	MsgHead         `json:"head"` ///"item", "queryitemlist"
	MasterItemID    int           `json:"equipitemid"`     ///受益装备的道具id
	MergeItemIDList IntList       `json:"mergeitemidlist"` ///被合并道具id列表
}

func (self *MergeItemMsg) GetTypeAndAction() (string, string) {
	return "item", "mergeitem"
}

func (self *MergeItemMsg) getNeedCoid(masterItem *Item) int {
	///得到当前等级升级所需金币
	staticDataMgr := GetServer().GetStaticDataMgr()
	itemInfo := masterItem.GetInfo()
	itemType := masterItem.GetTypeInfo()
	itemExpType := levelExpTypeEquipMerge + itemType.Sort - 1
	needCoin := staticDataMgr.GetLevelExpNeedCoin(itemExpType, itemInfo.MergeLevel+1)
	needExp := staticDataMgr.GetLevelExpNeedExp(itemExpType, itemInfo.MergeLevel+1)
	needCoin *= needExp ///需求金币 = 系数 * 融合点

	return needCoin
}

func (self *MergeItemMsg) checkAction(client IClient) bool { ///检测装备条件
	team := client.GetTeam()
	loger := GetServer().GetLoger()
	itemMgr := team.GetItemMgr()
	masterItem := itemMgr.GetItem(self.MasterItemID)
	if loger.CheckFail("masterItem!=nil", masterItem != nil, masterItem, nil) {
		return false ///不存在的融合主道具
	}
	masterItemInfo := masterItem.GetInfo()
	if loger.CheckFail("masterItemInfo.MergeLevel<ItemEvolveLevel", masterItemInfo.MergeLevel < ItemEvolveLevel,
		masterItemInfo.MergeLevel, ItemEvolveLevel) {
		return false ///请求融合主道具已达到当前品质的融合合等级上限
	}
	masterItemType := masterItem.GetTypeInfo()
	for i := range self.MergeItemIDList {
		mergeItemID := self.MergeItemIDList[i]
		mergeItem := itemMgr.GetItem(mergeItemID) ///得到被合并item
		if loger.CheckFail("mergeItem!=nil", mergeItem != nil, mergeItem, nil) {
			return false ///不存在的融合材料道具
		}
		mergeItemType := mergeItem.GetTypeInfo()
		if loger.CheckFail("mergeItem!=nil", masterItemType.Sort == mergeItemType.Sort,
			masterItemType.Sort, mergeItemType.Sort) {
			return false ///融合材料道具子类型与主道具不相符
		}
	}

	teamCoin := team.GetInfo().Coin
	needCoin := self.getNeedCoid(masterItem)
	if loger.CheckFail("teamCoin > needCoid", teamCoin > needCoin, teamCoin, needCoin) {
		return false ///金钱不足够支付融合费用
	}
	return true
}

func (self *MergeItemMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false {
		return false ///检测失败
	}

	team := client.GetTeam()
	syncMgr := client.GetSyncMgr() ///得到网络同步管理器
	ItemMgr := client.GetTeam().GetItemMgr()
	masterItem := ItemMgr.GetItem(self.MasterItemID) ///得到受益item
	needCoin := self.getNeedCoid(masterItem)
	team.PayCoin(needCoin)
	totalMergeExp := 0 ///总融合经验
	for i := range self.MergeItemIDList {
		mergeItemID := self.MergeItemIDList[i]
		mergeItem := ItemMgr.GetItem(mergeItemID)     ///得到被合并item
		totalMergeExp += mergeItem.GetTotalMergeExp() ///累加收益
	}
	ItemMgr.RemoveItem(self.MergeItemIDList)            ///销毁被合并的item
	syncMgr.SyncRemoveItem(self.MergeItemIDList)        ///同步客户端被销毁被合并的itemlist
	masterItem.AwardMergeExp(totalMergeExp)             ///得到受益融合值
	syncMgr.SyncObject(systemTypeItemMerge, masterItem) ///同步客户端此道具最终的属性
	syncMgr.SyncObject("MergeItemMsg", team)            ///同步客户端球队信息
	return true
}

const ( ///装备的操作类型
	equipItemOPTypeWield   = 1 ///装备
	equipItemOPTypeReplace = 2 ///替换
	equipItemOPTypeUnwield = 3 ///卸下
)

type EquipItemMsg struct { ///装备道具消息
	MsgHead       `json:"head"` ///"item", "queryitemlist"
	StarID        int           `json:"starid"`        ///需要装备的球员id
	EquipItemID   int           `json:"equipitemid"`   ///需要装备的道具id
	UnequipItemID int           `json:"unequipitemid"` ///需要卸下的道具id
	Type          int           `json:"type"`          ///1 装备 2 更换 3 卸下
}

func (self *EquipItemMsg) GetTypeAndAction() (string, string) {
	return "item", "equipitem"
}

///		球员品质	球衣1	球鞋3	个性物品2
///		绿			无		无		无
///		蓝			开启	无		无
///		紫			开启	开启	无
///		橙			开启	开启	开启
func (self *EquipItemMsg) checkItemOperable(client IClient) bool { ///检测球员装备操作条件
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	star := team.GetStar(self.StarID)
	starInfo := star.GetInfo()
	operItemID := self.EquipItemID
	if operItemID <= 0 {
		operItemID = self.UnequipItemID ///获取操作道具
	}
	operItem := itemMgr.GetItem(operItemID) ///尝试取得请求操作道具对象
	if loger.CheckFail("operItem!=nil", operItem != nil, operItem, nil) {
		return false ///不存在的操作道具
	}

	itemInfo := operItem.GetInfo()
	itemType := operItem.GetTypeInfo()
	result := true
	switch itemType.Sort { ///根据不同的操作道具来检测装备球员星级的合法性
	case ItemSortCloth:
		switch itemInfo.Color {
		case itemColorGreen:
			result = starInfo.EvolveCount >= 2
		case itemColorBlue:
			result = starInfo.EvolveCount >= 3
		case itemColorPurple:
			result = starInfo.EvolveCount >= 5
		case itemColorOrange:
			result = starInfo.EvolveCount >= 7
		}
	case ItemSortShoe:
		switch itemInfo.Color {
		case itemColorBlue:
			result = starInfo.EvolveCount >= 4
		case itemColorPurple:
			result = starInfo.EvolveCount >= 5
		case itemColorOrange:
			result = starInfo.EvolveCount >= 7
		}
	case ItemSortJewel:
		switch itemInfo.Color {
		case itemColorPurple:
			result = starInfo.EvolveCount >= 6
		case itemColorOrange:
			result = starInfo.EvolveCount >= 7
		}
	}

	return result
	//result := false
	//loger := GetServer().GetLoger()
	//team := client.GetTeam()
	//itemMgr := team.GetItemMgr()
	//star := team.GetStar(self.StarID)
	//starType := star.GetTypeInfo()
	//operItemID := self.EquipItemID
	//if operItemID <= 0 {
	//	operItemID = self.UnequipItemID
	//}
	//operItem := itemMgr.GetItem(operItemID) ///尝试取得请求操作道具对象
	//if loger.CheckFail("operItem!=nil", operItem != nil, operItem, nil) {
	//	return false ///不存在的操作道具
	//}
	//itemType := operItem.GetTypeInfo()
	//switch itemType.Sort { ///根据不同的操作道具来检测装备球员品质合法性
	//case ItemSortCloth: ///球衣
	//	result = loger.CheckFail("starType.Grade >= starGradeBlue",
	//		starType.Grade >= starGradeBlue, starType.Grade, starGradeBlue)
	//case ItemSortJewel: ///饰品
	//	result = loger.CheckFail("starType.Grade >= starGradeOrange",
	//		starType.Grade >= starGradeOrange, starType.Grade, starGradeOrange)
	//case ItemSortShoe: ///球鞋
	//	result = loger.CheckFail("starType.Grade >= starGradePurple",
	//		starType.Grade >= starGradePurple, starType.Grade, starGradePurple)
	//}
	//return false == result
}

func (self *EquipItemMsg) checkAction(client IClient) bool { ///检测装备条件
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	star := team.GetStar(self.StarID)
	if loger.CheckFail("star!=nil", star != nil, star, nil) {
		return false ///不存在的球员
	}
	if loger.CheckFail("self.checkItemOperable()==true", self.checkItemOperable(client) == true, true, false) {
		return false ///操作道具有效性检测失败
	}
	equipItem := itemMgr.GetItem(self.EquipItemID)     ///尝试取得请求装备道具对象
	unequipItem := itemMgr.GetItem(self.UnequipItemID) ///尝试取得请求装备道具对象
	switch self.Type {
	case equipItemOPTypeWield:
		if loger.CheckFail("equipItem!=nil", equipItem != nil, equipItem, nil) {
			return false ///不存在的装备道具
		}
		equipTypeInfo := equipItem.GetTypeInfo()
		equipItemInfo := equipItem.GetInfo()
		if loger.CheckFail("equipItemInfo.StarID==0", equipItemInfo.StarID == 0, equipItemInfo.StarID, 0) {
			return false ///装备拥有者必须为空
		}
		itemStored := star.GetItemFromSort(equipTypeInfo.Sort) ///查询是否已装备指定位置道具
		if loger.CheckFail("itemStored==nil", itemStored == nil, itemStored, nil) {
			return false ///装备空位必须为空,不得装备到已有装备的位置
		}
	case equipItemOPTypeUnwield:
		if loger.CheckFail("unequipItem!=nil", unequipItem != nil, unequipItem, nil) {
			return false ///不存在的拿下装备道具
		}
		unequipItemInfo := unequipItem.GetInfo()
		if loger.CheckFail("unequipItemInfo.StarID>0", unequipItemInfo.StarID > 0, unequipItemInfo.StarID, 0) {
			return false ///拿下的装备道具必须有拥有者
		}
		unequipTypeInfo := unequipItem.GetTypeInfo()
		unequipItemStored := star.GetItemFromSort(unequipTypeInfo.Sort) ///查询准备脱的装备是否已装备
		if loger.CheckFail("unequipItem==unequipItemStored", unequipItem == unequipItemStored,
			unequipItem, unequipItemStored) {
			return false ///拿下的装备道具并非是已装备的道具
		}
	case equipItemOPTypeReplace:
		if loger.CheckFail("equipItem!=nil", equipItem != nil, equipItem, nil) {
			return false ///不存在的装备道具
		}
		if loger.CheckFail("unequipItem!=nil", unequipItem != nil, unequipItem, nil) {
			return false ///不存在的拿下装备道具
		}
		equipTypeInfo := equipItem.GetTypeInfo()
		unequipTypeInfo := unequipItem.GetTypeInfo()
		if loger.CheckFail("equipTypeInfo.Sort== unequipTypeInfo.Sort",
			equipTypeInfo.Sort == unequipTypeInfo.Sort, equipTypeInfo.Sort, unequipTypeInfo.Sort) {
			return false ///替换的两个装备道具是不同类型的,禁止替换
		}
		equipItemInfo := equipItem.GetInfo()
		if loger.CheckFail("equipItemInfo.StarID<=0", equipItemInfo.StarID <= 0, equipItemInfo.StarID, 0) {
			return false ///替换上的道具必须没被其它球员装备
		}
	}
	return true
}

func (self *EquipItemMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false {
		return false ///检测失败
	}
	team := client.GetTeam()
	itemMgr := client.GetTeam().GetItemMgr()
	syncMgr := client.GetSyncMgr()                     ///得到网络同步管理器
	equipItem := itemMgr.GetItem(self.EquipItemID)     ///得到需要装备道具对象
	unequipItem := itemMgr.GetItem(self.UnequipItemID) ///得到需要卸下道具对象
	switch self.Type {
	case equipItemOPTypeWield:
		equipItem.SetStarID(self.StarID) ///设置装备拥有人
		equipTypeInfo := equipItem.GetTypeInfo()
		equipItem.SetCell(equipTypeInfo.Sort) ///设置装备单元格索引号
	case equipItemOPTypeUnwield:
		unequipItem.SetStarID(0) ///清零装备拥有人
		unequipItem.SetCell(0)   ///设置装备单元格索引号
	case equipItemOPTypeReplace:
		equipItem.SetStarID(self.StarID) ///设置装备拥有人
		equipTypeInfo := equipItem.GetTypeInfo()
		equipItem.SetCell(equipTypeInfo.Sort) ///设置装备单元格索引号
		unequipItem.SetStarID(0)              ///清零卸下装备拥有人
		unequipItem.SetCell(0)                ///设置装备单元格索引号
	}
	star := team.GetStar(self.StarID)
	star.CalcScore() ///装备信息变更后需要更新角色卡评分信息
	syncMgr.syncStarCalcInfo(star)
	syncObjectList := SyncObjectList{}
	if equipItem != nil {
		syncObjectList = append(syncObjectList, equipItem)
	}
	if unequipItem != nil {
		syncObjectList = append(syncObjectList, unequipItem)
	}
	syncMgr.SyncObjectArray(systemTypeItemEquip, syncObjectList) ///将道具属性变更消息发给客户端
	team.CalcScore()
	syncMgr.SyncObject("EquipItemMsg", team) ///同步客户端最新的战力评分
	syncMgr.SyncObject("EquipItemMsg", star) ///同步客户端最新的战力评分
	return true
}

type QueryItemListMsg struct { ///查询球队所拥有的所有道具列表消息
	MsgHead `json:"head"` ///"item", "queryitemlist"
	TeamID  int           `json:"teamid"` ///查询哪家球队的道具信息,查询自己也需要提供球队id
}

func (self *QueryItemListMsg) GetTypeAndAction() (string, string) {
	return "item", "queryitemlist"
}

func (self *QueryItemListMsg) processAction(client IClient) bool {
	itemMgr := client.GetTeam().GetItemMgr()
	msgQueryResult := new(QueryItemListResultMsg)
	msgQueryResult.TeamID = client.GetTeam().GetID()
	msgQueryResult.ItemList = itemMgr.GetItemInfoList()
	client.SendMsg(msgQueryResult)
	return true
}

//team := client.GetTeam()
//loger := GetServer().GetLoger()
//itemMgr := team.GetItemMgr()
//evolveItem := itemMgr.GetItem(self.ItemID)
//if loger.CheckFail("evolveItem!=nil", evolveItem != nil, evolveItem, nil) {
//	return false ///不存在的升阶道具
//}
//evolveItemInfo := evolveItem.GetInfo()
//if loger.CheckFail("(evolveItemInfo.MergeLevel%ItemEvolveLevel)==0", (evolveItemInfo.MergeLevel%ItemEvolveLevel) == 0,
//	evolveItemInfo.MergeLevel%ItemEvolveLevel, 0) {
//	return false ///未达到升阶条件
//}
//if loger.CheckFail("evolveItemInfo.Color<itemColorEnd", evolveItemInfo.Color < itemColorEnd,
//	evolveItemInfo.Color, itemColorEnd) {
//	return false ///已到达道具品质上限
//}
//mergeItem := itemMgr.GetItem(self.mergeItemID)
//if loger.CheckFail("mergeItem!=nil", mergeItem != nil, mergeItem, nil) {
//	return false ///不存在的升阶所需材料道具
//}
//evolveItemType := evolveItem.GetTypeInfo()
//if loger.CheckFail("evolveItemType!=nil", evolveItemType != nil, evolveItemType, nil) {
//	return false ///不存在的升阶道具类型
//}
//mergeItemType := mergeItem.GetTypeInfo()
//if loger.CheckFail("mergeItemType!=nil", mergeItemType != nil, mergeItemType, nil) {
//	return false ///不存在的升阶材料道具类型
//}
//if loger.CheckFail("evolveItemType.Type == mergeItemType.Type", evolveItemType.Type == mergeItemType.Type,
//	evolveItemType.Type, mergeItemType.Type) {
//	return false ///升阶道具与材料类型不相同
//}
//return true
type QueryItemListResultMsg struct { ///查询球队所拥有道具列表结果消息
	MsgHead  `json:"head"`
	TeamID   int          `json:"teamid"`   ///道具列表所属的球队id
	ItemList ItemInfoList `json:"itemlist"` ///道具信息列表
}

func (self *QueryItemListResultMsg) GetTypeAndAction() (string, string) {
	return "item", "itemlistresult"
}

func NewQueryItemListResultMsg() *QueryItemListResultMsg { ///生成查询消息
	msg := new(QueryItemListResultMsg)
	msg.ItemList = ItemInfoList{}
	return msg
}

type ItemEvolveMsg struct { ///请求道具升阶消息
	MsgHead     `json:"head"` ///"item", "itemevolve"
	ItemID      int           `json:"itemid"`    ///请求升阶道具id
	MergeItemID int           `json:"mergeItem"` ///被扣除道具id
}

func (self *ItemEvolveMsg) GetTypeAndAction() (string, string) {
	return "item", "itemevolve"
}

func (self *ItemEvolveMsg) checkAction(client IClient) bool { ///检测条件
	team := client.GetTeam()
	loger := GetServer().GetLoger()
	itemMgr := team.GetItemMgr()
	evolveItem := itemMgr.GetItem(self.ItemID)
	if loger.CheckFail("evolveItem!=nil", evolveItem != nil, evolveItem, nil) {
		return false ///不存在的升阶道具
	}
	evolveItemInfo := evolveItem.GetInfo()
	//if loger.CheckFail("evolveItemInfo.MergeLevel>=ItemEvolveLevel", evolveItemInfo.MergeLevel >= ItemEvolveLevel,
	//	evolveItemInfo.MergeLevel, ItemEvolveLevel) {
	//	return false ///装备融合等级必须大于等于ItemEvolveLevel
	//}
	//if loger.CheckFail("evolveItemInfo.MergeLevel%ItemEvolveLevel==0", (evolveItemInfo.MergeLevel%ItemEvolveLevel) == 0,
	//	evolveItemInfo.MergeLevel%ItemEvolveLevel, 0) {
	//	return false ///未达到升阶条件
	//}
	if loger.CheckFail("evolveItemInfo.Color<itemColorEnd", evolveItemInfo.Color < itemColorEnd,
		evolveItemInfo.Color, itemColorEnd) {
		return false ///已到达道具品质上限
	}
	mergeItem := itemMgr.GetItem(self.MergeItemID)
	if loger.CheckFail("mergeItem!=nil", mergeItem != nil, mergeItem, nil) {
		return false ///不存在的升阶所需材料道具
	}
	mergeItemInfo := mergeItem.GetInfo()
	if loger.CheckFail("mergeItemInfo.StarID == 0", mergeItemInfo.StarID == 0, mergeItemInfo.StarID, 0) {
		return false ///升阶材料道具不得是已装备的道具
	}
	evolveItemType := evolveItem.GetTypeInfo()
	if loger.CheckFail("evolveItemType!=nil", evolveItemType != nil, evolveItemType, nil) {
		return false ///不存在的升阶道具类型
	}
	mergeItemType := mergeItem.GetTypeInfo()
	if loger.CheckFail("mergeItemType!=nil", mergeItemType != nil, mergeItemType, nil) {
		return false ///不存在的升阶材料道具类型
	}
	if loger.CheckFail("evolveItemType.Type == mergeItemType.Type", evolveItemType.Type == mergeItemType.Type,
		evolveItemType.Type, mergeItemType.Type) {
		return false ///升阶道具与材料类型不相同
	}
	if loger.CheckFail("evolveItemInfo.Color == mergeItemInfo.Color", evolveItemInfo.Color == mergeItemInfo.Color,
		evolveItemInfo.Color, mergeItemInfo.Color) {
		return false ///升阶道具必须与材料颜色相同
	}

	if evolveItemInfo.StarID != 0 { ///判断该道具是否被穿戴
		star := team.GetStar(evolveItemInfo.StarID)
		starInfo := star.GetInfo()
		nextLevelColor := evolveItemInfo.Color + 1 ///下一阶品质颜色

		///对比星级与颜色,判断是否达到升阶要求
		switch nextLevelColor {
		case itemColorGreen:
			if loger.CheckFail("StarInfo.EvolveCount >= 1", starInfo.EvolveCount >= 1,
				starInfo.EvolveCount, 1) {
				return false
			}
		case itemColorBlue:
			if loger.CheckFail("StarInfo.EvolveCount >= 3", starInfo.EvolveCount >= 3,
				starInfo.EvolveCount, 3) {
				return false
			}

		case itemColorPurple:
			if loger.CheckFail("StarInfo.EvolveCount >= 5", starInfo.EvolveCount >= 5,
				starInfo.EvolveCount, 5) {
				return false
			}

		case itemColorOrange:
			if loger.CheckFail("StarInfo.EvolveCount >= 7", starInfo.EvolveCount >= 7,
				starInfo.EvolveCount, 7) {
				return false
			}
		}
	}
	return true
}

func (self *ItemEvolveMsg) payAction(client IClient) bool { ///支付代价
	team := client.GetTeam()
	syncMgr := client.GetSyncMgr()
	itemMgr := team.GetItemMgr()
	removeItemIDList := IntList{self.MergeItemID}
	itemMgr.RemoveItem(removeItemIDList)     ///扣除升阶材料装备
	syncMgr.SyncRemoveItem(removeItemIDList) ///同步客户端删除指定装备
	return true
}

func (self *ItemEvolveMsg) doAction(client IClient) bool { ///发货
	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	syncMgr := client.GetSyncMgr()
	evolveItem := itemMgr.GetItem(self.ItemID)
	itemInfo := evolveItem.GetInfo()
	itemInfo.Color++        ///品质提升一级
	itemInfo.MergeLevel = 0 ///融合等级提升一级
	syncMgr.SyncObject("ItemEvolveMsg", evolveItem)
	return true
}

func (self *ItemEvolveMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false { ///检测
		return false
	}
	if self.payAction(client) == false { ///支付
		return false
	}
	if self.doAction(client) == false { ///发货
		return false
	}
	return true
}

type QueryStoreResultMsg struct { ///查询球队所拥有道具列表结果消息
	MsgHead  `json:"head"` ///"item", "querystoreresult"
	ItemList ItemInfoList  `json:"storeitemlist"` ///仓库道具信息列表
}

func (self *QueryStoreResultMsg) GetTypeAndAction() (string, string) {
	return "item", "querystoreresult"
}

type QueryStoreMsg struct { ///查询仓库道具列表
	MsgHead `json:"head"` ///"item", "querystore"
}

func (self *QueryStoreMsg) GetTypeAndAction() (string, string) {
	return "item", "querystore"
}

func (self *QueryStoreMsg) processAction(client IClient) bool {
	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	itemSlice := itemMgr.GetItemSlice() ///除得球队所有道具列表
	queryStoreResultMsg := new(QueryStoreResultMsg)
	for i := range itemSlice {
		item := itemSlice[i]
		itemInfo := item.GetInfo()
		if itemInfo.Position != itemPosStore {
			continue ///忽略掉所有不在仓库中的道具
		}
		queryStoreResultMsg.ItemList = append(queryStoreResultMsg.ItemList, *itemInfo)
	}
	client.SendMsg(queryStoreResultMsg)
	return true
}

type ItemSellMsg struct { ///出售道具
	MsgHead `json:"head"` ///"item", "sell"
	ItemID  int           `json:"itemid"` ///售卖道具id
	Count   int           `json:"count"`  ///售卖道具数量
}

func (self *ItemSellMsg) GetTypeAndAction() (string, string) {
	return "item", "sell"
}

func (self *ItemSellMsg) checkAction(client IClient) bool {
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	item := itemMgr.GetItem(self.ItemID)
	if loger.CheckFail(" item!=nil", item != nil, item, nil) {
		return false ///升阶材料道具不得是已装备的道具
	}

	itemInfo := item.GetInfo()
	if loger.CheckFail("itemCount >= self.Count", itemInfo.Count >= self.Count, itemInfo.Count, self.Count) {
		return false //!物品数量必须足够
	}

	return true
}

func (self *ItemSellMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	syncMgr := client.GetSyncMgr()
	item := itemMgr.GetItem(self.ItemID)
	itemType := item.GetTypeInfo()
	awardCoin := itemType.SellCoin * self.Count ///计算总售价
	itemMgr.PayItem(self.ItemID, self.Count)    ///扣道具
	item = itemMgr.GetItem(self.ItemID)         ///再取一次道具
	if item != nil {
		syncMgr.SyncObject("ItemSellMsg", item) ///同步最新的道具信息
	} else {
		syncMgr.SyncRemoveItem(IntList{self.ItemID}) ///同步道具删除
	}
	team.AwardCoin(awardCoin)               ///加赚的球币
	syncMgr.SyncObject("ItemSellMsg", team) ///同步最新的球队信息
	return true
}

func (self *ItemSellMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false { ///检测
		return false
	}
	if self.doAction(client) == false { ///发货
		return false
	}
	return true
}

type ItemUseMsg struct { ///道具使用消息
	MsgHead   `json:"head"` ///"item", "use"
	ItemID    int           `json:"itemid"`    ///使用道具id
	UseNumber int           `json:"usenumber"` ///使用多少个
}

func (self *ItemUseMsg) GetTypeAndAction() (string, string) {
	return "item", "use"
}

func (self *ItemUseMsg) checkAction(client IClient) bool { ///检测
	loger := loger()
	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	item := itemMgr.GetItem(self.ItemID)
	if loger.CheckFail("item!=nil", item != nil, item, nil) {
		return false ///道具必须存在
	}
	itemType := item.GetTypeInfo()
	if loger.CheckFail("itemType!=nil", itemType != nil, itemType, nil) {
		return false ///道具类型必须存在
	}
	useActionType := item.GetUseAction()
	if loger.CheckFail("useActionType!=nil", useActionType != nil, useActionType, nil) {
		return false ///道具必须可使用
	}

	isEnough := itemMgr.HasEnoughItem(itemType.ID, self.UseNumber)
	if loger.CheckFail("isEnough == true", isEnough == true, isEnough, true) {
		return false ///道具数量不足
	}

	//特殊道具处理
	if item.Type >= 300050 && item.Type <= 300055 {
		self.UseNumber = 1
	}

	awardTypeList, _, _, _, _ := item.GetActionAwardList()
	needFreeCellCount := 0
	for i := range awardTypeList {
		awardType := awardTypeList[i]
		if awardType <= 0 {
			continue ///只关注道具类奖励
		}
		itemType := GetItemType(awardType)
		if loger.CheckFail("itemType!=nil", itemType != nil, itemType, nil) {
			return false ///道具类型必须存
		}
		if itemType.IsNumberType() == true {
			continue ///忽略掉数值类道具
		}

		itemCount := self.UseNumber
		_, needCellCount := itemMgr.TryComboItem(itemType.ID, itemCount) ///先判断是否能完全叠加
		if needCellCount <= 0 {
			continue ///如果可以完全叠加则不占用格子
		}

		needFreeCellCount++
	}
	if loger.CheckFail("team.StoreCapacity >=needFreeCellCount", team.StoreCapacity >= needFreeCellCount,
		team.StoreCapacity, needFreeCellCount) {
		return false ///仓库空格数必须以足够
	}
	return true
}

func (self *ItemUseMsg) payAction(client IClient) bool { ///支付
	return true
}

func (self *ItemUseMsg) doAction(client IClient) bool { ///发货
	team := client.GetTeam()
	syncMgr := client.GetSyncMgr()
	itemMgr := team.GetItemMgr()
	item := itemMgr.GetItem(self.ItemID)
	awardTypeList, awardCountList, awardGradeList, awardStarList, awardProbabilityList := item.GetActionAwardList()

	itemMgr.PayItem(self.ItemID, self.UseNumber) ///扣除道具
	item = itemMgr.GetItem(self.ItemID)
	if item != nil {
		syncMgr.SyncObject("ItemUseMsg", item)
	} else {
		syncMgr.SyncRemoveItem(IntList{self.ItemID})
	}

	if self.UseNumber > 1 {
		for i := 0; i < self.UseNumber; i++ {
			for j := range awardTypeList {
				awardType := awardTypeList[j]
				awardCount := awardCountList[j]
				team.awardNumberType(awardType, awardCount)
			}
		}

		syncMgr.SyncObject("ItemUseMsg", team)
		return true /// 处理完毕,返回
	}

	randNum := Random(0, 10000)
	nValue := 0
	for i := range awardTypeList {
		awardType := awardTypeList[i]
		awardCount := awardCountList[i]
		awardGrade := awardGradeList[i]
		awardStar := awardStarList[i]

		if awardProbabilityList[i] != 0 {
			if randNum >= nValue && randNum < nValue+awardProbabilityList[i] {
				team.AwardObject(awardType, awardCount, awardGrade, awardStar) ///给球队发送奖励
				break
			}
			nValue += awardProbabilityList[i]

		} else {
			team.AwardObject(awardType, awardCount, awardGrade, awardStar) ///给球队发送奖励
		}
		if awardType == awardTypeTicket && awardCount > 0 {
			client.RechargeRecord(Get_ItemUse, awardCount)
		}
	}
	return true
}

func (self *ItemUseMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false { ///检测
		return false
	}
	if self.payAction(client) == false { ///支付
		return false
	}
	if self.doAction(client) == false { ///发货
		return false
	}
	return true
}

type ItemCombineMsg struct { ///客户端请求道具合成消息
	MsgHead `json:"head"` ///"item", "combine"
	ItemID  int           `json:"itemid"`  ///使用道具id
	MakeNum int           `json:"makenum"` ///合成多少个
}

func (self *ItemCombineMsg) GetTypeAndAction() (string, string) {
	return "item", "combine"
}

func (self *ItemCombineMsg) checkAction(client IClient) bool { ///检测
	loger := loger()                                                 ///记录对象
	team := client.GetTeam()                                         ///球队对象
	itemMgr := team.GetItemMgr()                                     ///球队的道具管理器
	item := itemMgr.GetItem(self.ItemID)                             ///取得对应的道具对象
	if loger.CheckFail("MakeNum > 0", self.MakeNum > 0, self, nil) { ///要合成的数量必须大于0
		return false
	}
	if loger.CheckFail("item!=nil", item != nil, item, nil) { ///道具必须存在
		return false
	}
	itemType := item.GetTypeInfo()                                        ///取得道具类型对象
	if loger.CheckFail("itemType!=nil", itemType != nil, itemType, nil) { ///道具类型必须存在
		return false
	}

	if loger.CheckFail("itemType.BuyMask != 0", itemType.BuyMask != 0, itemType, nil) { ///判断有无合成物
		return false
	}

	maskItemType := GetItemType(itemType.BuyMask)
	if loger.CheckFail("itemType!=nil", maskItemType != nil, maskItemType, nil) { ///生成的道具类型必须存在
		return false
	}

	///检查数量
	needNumber := 0         ///合成需要的材料数量
	productNumber := 0      ///合成个数
	if itemType.Type != 5 { ///不是5类型的碎片就允许多个一起合成
		productNumber = self.MakeNum
	} else { ///5类型的碎片最多只能合成一个
		productNumber = 1
	}
	needNumber = productNumber * itemType.BuyTicket

	isEnough := itemMgr.HasEnoughUnuseItem(itemType.ID, needNumber) ///这是取得背包内全部的此类道具的数量再跟需求比

	if loger.CheckFail("isEnough == true", isEnough == true, isEnough, true) { ///判断所需数量
		return false ///道具数量不足
	}
	self.MakeNum = productNumber ///改写生成数量

	///类型为5的球员碎片要判断背包内是否有对应合成物的球员卡以及队伍内是否已有该球员
	if itemType.Type == 5 {
		maskitem := itemMgr.GetItemFromType(itemType.BuyMask)
		if loger.CheckFail("maskitem is't exist", maskitem == nil, maskitem, true) { ///检查背包内是否包含该球员卡
			return false
		}
		startypeid := itemType.ID % 10000            ///去除id的前面的50得到合成球员卡能产生的球员的id
		maskstar := team.GetStarFromType(startypeid) ///取得球队中此类型的球员

		if loger.CheckFail("maskstar is't exist", maskstar == nil, maskstar, true) { ///判断球队中有无此类型的球员
			return false
		}
	}

	_, needCellCount := itemMgr.TryComboItem(itemType.BuyMask, productNumber) ///先判断是否能完全叠加

	if loger.CheckFail("team.StoreCapacity >=needCellCount", team.StoreCapacity >= needCellCount,
		team.StoreCapacity, needCellCount) {
		return false ///仓库空格数必须以足够
	}

	return true
}

func (self *ItemCombineMsg) payAction(client IClient) bool { ///支付
	return true
}

func (self *ItemCombineMsg) doAction(client IClient) bool { ///发货
	team := client.GetTeam()
	syncMgr := client.GetSyncMgr()
	itemMgr := team.GetItemMgr()
	item := itemMgr.GetItem(self.ItemID)
	itemType := item.GetTypeInfo()                  ///没删之前先取得道具类型
	maskItemType := GetItemType(itemType.BuyMask)   ///取得合成物的道具类型
	needNumber := self.MakeNum * itemType.BuyTicket ///合成需要的材料数量
	productNumber := self.MakeNum                   ///计算要合成的数量

	removeList, updatItemID, _ := itemMgr.PayItemType(itemType.ID, needNumber) ///扣除道具
	updatItem := itemMgr.GetItem(updatItemID)

	team.AwardObject(maskItemType.ID, productNumber, 0, 0) ///给球队发送奖励

	if updatItem != nil { ///同步属性更改
		syncMgr.SyncObject("ItemCombineMsg", updatItem)
	}
	if removeList.Len() > 0 { ///有道具删除消息需要同步客户端
		syncMgr.SyncRemoveItem(removeList)
	}

	return true
}

func (self *ItemCombineMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false { ///检测
		return false
	}
	if self.payAction(client) == false { ///支付
		return false
	}
	if self.doAction(client) == false { ///发货
		return false
	}
	return true
}

//! 融合球员碎片消息
type MergeCardMsg struct {
	MsgHead   `json:"head"` //!"merge", "card"
	CardIDLst IntList       `json:"cardlst"` //! 请求融合的球员碎片
	groupNum  int
}

func (self *MergeCardMsg) GetTypeAndAction() (string, string) {
	return "merge", "card"
}

func (self *MergeCardMsg) checkAction(client IClient) bool {
	loger := GetServer().GetLoger()
	cardLength := len(self.CardIDLst)

	if loger.CheckFail("cardNum == 3", cardLength == 3, cardLength, 3) {
		return false //! 碎片数量不对
	}

	var cardColorLst IntList
	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	for i := 0; i < cardLength; i++ {
		if loger.CheckFail("self.CardIDLst[i] != 0", self.CardIDLst[i] != 0, self.CardIDLst[i], 0) {
			return false //! 碎片ID不允许为零
		}

		itemInfo := itemMgr.GetItem(self.CardIDLst[i])
		if loger.CheckFail("itemInfo != nil", itemInfo != nil, itemInfo, nil) {
			return false
		}

		staticDataMgr := GetServer().GetStaticDataMgr()
		cardType := staticDataMgr.GetItemType(itemInfo.Type)
		if loger.CheckFail("cardType != nil", cardType != nil, cardType, nil) {
			return false //! 静态表不存在球员碎片ID信息
		}

		if cardType.Color > 4 || cardType.Color < 3 {
			loger.Warn("Only S and A class star can merge")
			return false
		}

		if itemMgr.HasEnoughItem(cardType.ID, 1) == false {
			loger.Warn("Item num is not enough")
			return false
		}

		cardColorLst = append(cardColorLst, cardType.Color)
	}

	if (cardColorLst[0] == cardColorLst[1] && cardColorLst[0] == cardColorLst[2]) == false {
		loger.Warn("Merge need same color") //! 融合需要碎片品质相同
		return false
	}

	if self.CardIDLst[0] == self.CardIDLst[1] && self.CardIDLst[1] == self.CardIDLst[2] {

		itemInfo := itemMgr.GetItem(self.CardIDLst[0])

		if itemMgr.HasEnoughItem(itemInfo.Type, 3) == false {
			loger.Warn("Item num is not enough")
			return false
		}
	}
	return true
}

func (self *MergeCardMsg) payAction(client IClient) bool {

	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	syncMgr := client.GetSyncMgr()
	staticDataMgr := GetServer().GetStaticDataMgr()

	itemInfo := itemMgr.GetItem(self.CardIDLst[0])
	itemType := staticDataMgr.GetItemType(itemInfo.Type)

	if itemType.Color == 3 {
		self.groupNum = drawGroupTypeMergeStar1
	} else if itemType.Color == 4 {
		self.groupNum = drawGroupTypeMergeStar2
	}

	//! 扣除道具
	removeItemLst := IntList{}
	syncItemList := IntList{}
	for i := 0; i < len(self.CardIDLst); i++ {
		itemInfo = itemMgr.GetItem(self.CardIDLst[i])

		removeList, influenceItemID, _ := itemMgr.PayItemType(itemInfo.Type, 1)

		for i := 0; i < removeList.Len(); i++ {
			removeItemLst = append(removeItemLst, removeList[i])
		}

		syncItemList = append(syncItemList, influenceItemID)
		syncItemList = syncItemList.Unique()
	}

	if removeItemLst.Len() > 0 {
		syncMgr.SyncRemoveItem(removeItemLst) ///同步道具删除
	}

	if syncItemList.Len() > 0 {
		for i := 0; i < syncItemList.Len(); i++ {
			item := itemMgr.GetItem(syncItemList[i])
			if item != nil {
				syncMgr.SyncObject("MergeCardMsg", item)
			}
		}
	}

	return true
}

func (self *MergeCardMsg) doAction(client IClient) bool {

	team := client.GetTeam()
	staticDataMgr := GetServer().GetStaticDataMgr()
	drawGroupListItem := staticDataMgr.GetDrawGroupIndexList(self.groupNum)
	totalTakeWeightItem, totalShowWeightItem := discoverGetDrawWeightTotal(drawGroupListItem) ///得到权重总和
	drawShowOne, _ := discoverDrawOne(&drawGroupListItem, &totalTakeWeightItem, &totalShowWeightItem, 0, true)
	if drawShowOne != 0 {
		team.AwardObject(drawShowOne, 1, 0, 0)
	}

	return true
}

func (self *MergeCardMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false {
		return false
	}
	if self.payAction(client) == false {
		return false
	}
	if self.doAction(client) == false {
		return false
	}
	return true
}
