package football

import (
	"fmt"
)

///道具管理器
//type IGameItemMgr interface {
//	IGameMgr
//	Init(teamID int) bool
//	GetStarItemSlice(starID int) ItemSlice                                    ///得到球员道具对象列表
//	GetItemInfoList() ItemInfoList                                            ///得到所有道具信息列表
//	GetItem(itemID int) IGameItem                                             ///得到道具对象
//	RemoveItem(removeItemIDList IntList) bool                                 ///销毁一个道具
//	HasEnoughItem(itemType int, itemCount int) bool                           ///判断是否拥有足够数量的道具
//	AwardItem(itemType int, itemCount int) int                                ///添加道具信息,成功返回道具对象id
//	GetItemSlice() ItemSlice                                                  ///得到所有道具的Slice
//	GetStarItemInfoList(starIDList IntList) ItemInfoList                      ///得到球员道具信息列表
//	PayItem(itemID int, itemCount int) bool                                   ///消费指定道具指定数量
//	PayItemType(itemType int, itemCount int) (IntList, bool)                  ///消费指定道具指定数量
//	GetStarItemTypeList(starID int) IntList                                   ///得到球员道具类型列表
//	GetItemCountByPos(itemPos int) int                                        ///计算指定位置的道具个数
//	TryComboItem(itemType int, itemCount int) (int, int)                      ///尝试此种数量的道具进行叠加,并返回剩于数量
//	GetEquipAttributeAddition(starIDList IntList) (float32, float32, float32) ///得到装备属性加成
//}

type ItemSlice []*Item
type ItemList map[int]*Item
type ItemMgr struct {
	GameMgr
	itemList ItemList ///道具列表
}

func (self *ItemMgr) GetType() int { ///得到管理器类型
	return mgrTypeItemMgr ///关卡管理器
}

func (self *ItemMgr) SaveInfo() { ///保存数据
	for _, v := range self.itemList {
		v.Save()
	}
}

func (self *ItemMgr) GetItemSlice() ItemSlice { ///得到所有道具的Slice
	itemSlice := ItemSlice{}
	for _, v := range self.itemList {
		itemSlice = append(itemSlice, v)
	}
	return itemSlice
}

func (self *ItemMgr) GetItemInfoList() ItemInfoList { ///外界提供道具信息列表由此函数填充
	itemInfoList := ItemInfoList{}
	for _, v := range self.itemList {
		itemInfoList = append(itemInfoList, *v.GetInfo())
	}
	return itemInfoList
}

func (self *ItemMgr) PayItemType(itemType int, itemCount int) (IntList, int, bool) { ///消费指定道具指定数量
	removeItemIDList := IntList{}
	//	changeItemIDList := IntList{}
	influenceItemID := 0
	//if self.HasEnoughItem(itemType, itemCount) == false {///这个检测是不排除已装备道具的
	if self.HasEnoughUnuseItem(itemType, itemCount) == false { ///改成了判断已装备道具的，因为下面也是不考虑已装备道具
		return removeItemIDList, 0, false ///不够扣
	}
	for _, v := range self.itemList {
		itemInfo := v.GetInfo()
		if itemInfo.StarID > 0 {
			continue ///已装备道具不得消费
		}
		if itemInfo.Type != itemType {
			continue ///非目标道具类型,不处理
		}
		if itemInfo.Count > itemCount {
			itemInfo.Count -= itemCount ///够扣了直接扣数量

			if itemCount != 0 {
				influenceItemID = itemInfo.ID ///该物品数量有变动
			}
			//changeItemIDList = append(changeItemIDList, itemInfo.ID)
			break
		}
		itemCount -= itemInfo.Count ///部分扣除
		removeItemIDList = append(removeItemIDList, itemInfo.ID)
	}

	///整理背包,尝试叠加
	/*
		if remainItemID != 0 {
			removeIDList, changeIDList := self.FinishingBag(remainItemID)
			removeItemIDList = append(removeItemIDList, removeIDList...)
			changeItemIDList = append(changeItemIDList, changeIDList...)
		}*/

	if removeItemIDList.Len() > 0 {
		self.RemoveItem(removeItemIDList)
	}
	return removeItemIDList, influenceItemID, true
}

// func (self *ItemMgr) FinishingBag(itemID int) (IntList, IntList) { //排序背包 暂时屏蔽
// 	item := self.GetItem(itemID)
// 	itemType := item.GetTypeInfo()
// 	removeItemIDList := IntList{}
// 	changeItemIDList := IntList{}
// 	if nil == item {
// 		return nil, nil //道具必须存在
// 	}

// 	for _, v := range self.itemList {
// 		itemInfo := v.GetInfo()
// 		if itemInfo.StarID > 0 {
// 			continue ///已装备道具不得消费
// 		}
// 		if itemInfo.Type != item.Type {
// 			continue ///非目标道具类型不处理
// 		}

// 		if itemInfo.Count == itemType.Overlay {
// 			continue ///为整组的不处理
// 		}

// 		if itemInfo.Count+item.Count > itemType.Overlay {
// 			remainItemCount := itemType.Overlay - itemInfo.Count
// 			item.Count += remainItemCount //补全该物品
// 			itemInfo.Count -= remainItemCount
// 			changeItemIDList = append(changeItemIDList, itemInfo.ID)

// 			break
// 		}

// 		item.Count += itemInfo.Count
// 		removeItemIDList = append(removeItemIDList, itemInfo.ID)
// 	}

// 	return removeItemIDList, changeItemIDList
// }

func (self *ItemMgr) PayItem(itemID int, itemCount int) bool { ///消费指定道具指定数量
	result := true
	item := self.itemList[itemID]
	if nil == item {
		return false ///找不到此道具
	}
	itemInfo := item.GetInfo()
	if itemInfo.Count < itemCount {
		return false ///数量不够扣
	}
	itemInfo.Count -= itemCount
	if itemInfo.Count <= 0 {
		result = self.RemoveItem(IntList{itemID})
	}
	return result
}

///尝试此种数量的道具进行叠加,并返回剩于数量
func (self *ItemMgr) TryComboItem(itemType int, itemCount int) (int, int) {
	remainItemCount := itemCount
	comboItemType := GetServer().GetStaticDataMgr().Unsafe().GetItemType(itemType)
	if comboItemType.IsNumberType() == true {
		return itemCount, 0 ///数值类道具直接返回全部数量,并返回占0格.
	}
	for _, v := range self.itemList {
		itemInfo := v.GetInfo()
		if itemInfo.Type != itemType {
			continue ///非目标道具类型,不处理
		}
		itemType := v.GetTypeInfo()
		if itemInfo.Count >= itemType.Overlay {
			continue ///已叠加满
		}
		needCount := itemType.Overlay - itemInfo.Count ///可继续叠加数
		if needCount >= remainItemCount {              ///可叠完
			remainItemCount = 0
			break
		}
		///不可叠完
		remainItemCount -= needCount
	}
	return remainItemCount, remainItemCount / comboItemType.Overlay
}

func (self *ItemMgr) ComboItem(itemTypeID int, itemCount int) int { ///叠加道具,返回叠加后道具剩余个数
	remainItemCount := itemCount
	syncItemList := SyncObjectList{} ///同步道具id列表
	for _, v := range self.itemList {
		itemInfo := v.GetInfo()
		if itemInfo.Type != itemTypeID {
			continue ///非目标道具类型,不处理
		}
		itemType := v.GetTypeInfo()
		if itemInfo.Count >= itemType.Overlay {
			continue ///已叠加满
		}
		syncItemList = append(syncItemList, v)
		needCount := itemType.Overlay - itemInfo.Count ///可继续叠加数
		if needCount >= remainItemCount {              ///可叠完
			itemInfo.Count += remainItemCount
			remainItemCount = 0
			break
		}
		///不可叠完
		itemInfo.Count += needCount
		remainItemCount -= needCount
	}
	if syncItemList.Len() > 0 {
		self.syncMgr.SyncObjectArray("ComboItem", syncItemList)
	}
	return remainItemCount
}

func (self *ItemMgr) AwardItem(itemType int, itemCount int) IntList { ///添加道具信息,成功返回道具对象id
	if 0 == itemType {
		return nil
	}
	if itemCount == 0 { ///如果要获得的道具数量为0，那就是什么也没加所以也返回nil
		return nil
	}
	awardItemType := GetServer().GetStaticDataMgr().Unsafe().GetItemType(itemType)
	if nil == awardItemType {
		return nil
	}
	if ItemTypeDummy == awardItemType.Type { ///判断新奖励的物品是否是数值类
		GetServer().GetLoger().Warn("ItemTypeDummy == awardItemType.Type  itemtype:%d", itemType)
		return nil ///为数值类则无法生成道具
	}
	if ItemTypeEquip != awardItemType.Type { ///判断新奖励的物品是否是非装备
		isStoreFull := self.team.IsStoreFull(itemType, itemCount) ///非装备物品需要判断仓库是否已满
		if true == isStoreFull {
			return nil
		}
	}
	itemCount = self.ComboItem(itemType, itemCount)
	if itemCount <= 0 {
		return nil ///已叠加完了没有剩余用于生成新的道具
	}
	itemPosType := itemPosStore ///默认奖励道具进球队仓库
	if ItemTypeEquip == awardItemType.Type {
		itemPosType = itemPosTeamEquip ///如果奖励道具是装备则进球队装备栏中
	}

	insertItemList := IntList{}
	needBoxCount := 0
	remainItemCount := itemCount
	if itemCount > awardItemType.Overlay {
		if awardItemType.Overlay <= 0 {
			return nil
		}

		needBoxCount = itemCount / awardItemType.Overlay
		remainItemCount = itemCount - (needBoxCount * awardItemType.Overlay)
		for i := 0; i < needBoxCount; i++ {
			awardItemQuery := fmt.Sprintf("insert %s set teamid=%d,type=%d,count=%d,color=%d,position=%d",
				tableItem, self.team.GetID(), itemType, awardItemType.Overlay, awardItemType.Color, itemPosType) ///组插入记录SQL
			lastInsertItemID, _ := GetServer().GetDynamicDB().Exec(awardItemQuery)
			if lastInsertItemID <= 0 {
				GetServer().GetLoger().Warn("ItemMgr AwardItem fail! itemType:%d itemCount:%d", itemType, awardItemType.Overlay)
				return nil
			}
			///创建item对象
			loadItemQuery := fmt.Sprintf("select * from %s where id=%d limit 1", tableItem, lastInsertItemID)
			itemInfo := new(ItemInfo)
			GetServer().GetDynamicDB().fetchOneRow(loadItemQuery, itemInfo)
			self.itemList[itemInfo.ID] = NewItem(itemInfo)

			insertItemList = append(insertItemList, itemInfo.ID)
		}
	}

	if remainItemCount == 0 {
		return insertItemList ///无残留物品
	}
	awardItemQuery := fmt.Sprintf("insert %s set teamid=%d,type=%d,count=%d,color=%d,position=%d",
		tableItem, self.team.GetID(), itemType, remainItemCount, awardItemType.Color, itemPosType) ///组插入记录SQL
	lastInsertItemID, _ := GetServer().GetDynamicDB().Exec(awardItemQuery)
	if lastInsertItemID <= 0 {
		GetServer().GetLoger().Warn("ItemMgr AwardItem fail! itemType:%d itemCount:%d", itemType, remainItemCount)
		return nil
	}
	///创建item对象
	loadItemQuery := fmt.Sprintf("select * from %s where id=%d limit 1", tableItem, lastInsertItemID)
	itemInfo := new(ItemInfo)
	GetServer().GetDynamicDB().fetchOneRow(loadItemQuery, itemInfo)
	self.itemList[itemInfo.ID] = NewItem(itemInfo)
	insertItemList = append(insertItemList, itemInfo.ID)
	//itemInfo := new(ItemInfo)
	//itemInfo.ID = lastInsertItemID
	//itemInfo.TeamID = self.team.GetID()
	//itemInfo.Type = itemType
	//itemInfo.Count = itemCount
	//itemInfo.Color = awardItemType.Color
	//item := NewItem(itemInfo) ///生成道具对象
	//self.itemList[lastInsertItemID] = item
	return insertItemList
}

func (self *ItemMgr) GetItem(itemID int) *Item { ///得到道具对象
	return self.itemList[itemID]
}

///判断是否拥有足够数量的道具
func (self *ItemMgr) HasEnoughItem(itemType int, itemCount int) bool {
	storeItemCount := 0
	for _, v := range self.itemList {
		itemInfo := v.GetInfo()
		if itemInfo.Type == itemType {
			storeItemCount += itemInfo.Count
		}
	}
	return storeItemCount >= itemCount
}

func (self *ItemMgr) HasEnoughUnuseItem(itemType int, itemCount int) bool { ///判断是否拥有足够数量的未使用道具
	storeItemCount := 0
	for _, v := range self.itemList {
		itemInfo := v.GetInfo()
		if itemInfo.StarID > 0 { ///排除已装备的道具
			continue
		}
		if itemInfo.Type == itemType {
			storeItemCount += itemInfo.Count
		}
	}
	return storeItemCount >= itemCount
}

func (self *ItemMgr) GetItemCountByPos(itemPos int) int { ///计算指定位置的道具个数
	itemCount := 0
	for _, item := range self.itemList {
		itemInfo := item.GetInfo()
		if itemInfo.Position == itemPos {
			itemCount++
		}
	}
	return itemCount
}

func (self *ItemMgr) GetItemCountByType(itemType int) int { ///得到背包中指定类型的道具的数量，不包括已装备道具
	storeItemCount := 0
	for _, v := range self.itemList {
		itemInfo := v.GetInfo()
		if itemInfo.StarID > 0 {
			continue
		}
		if itemInfo.Type == itemType {
			storeItemCount += itemInfo.Count
		}
	}
	return storeItemCount
}

func (self *ItemMgr) GetStarItemTypeList(starID int) IntList { ///得到球员道具类型列表
	itemTypeList := IntList{}
	for _, v := range self.itemList {
		itemInfo := v.GetInfo()
		if itemInfo.StarID == starID {
			itemTypeList = append(itemTypeList, itemInfo.Type)
		}
	}
	return itemTypeList
}

func (self *ItemMgr) GetStarItemSlice(starID int) ItemSlice { ///得到球员道具对象列表
	itemSlice := ItemSlice{}
	for _, v := range self.itemList {
		itemInfo := v.GetInfo()
		if itemInfo.StarID == starID {
			itemSlice = append(itemSlice, v)
		}
	}
	return itemSlice
}

//func (self *ItemMgr) GetStarItemList(starIDList IntList) ItemList { ///得到球员道具信息列表
//	ItemList := ItemList{}
//	for i := range starIDList {
//		starID := starIDList[i]
//		itemList := self.GetStarItemSlice(starID)
//		for j := range itemList {
//			item := itemList[j]
//			itemInfoList = append(itemInfoList, *item.GetInfo())
//		}
//	}
//	return itemInfoList
//}

func (self *ItemMgr) GetStarItemInfoList(starIDList IntList) ItemInfoList { ///得到球员道具信息列表
	itemInfoList := ItemInfoList{}
	for i := range starIDList {
		starID := starIDList[i]
		itemList := self.GetStarItemSlice(starID)
		for j := range itemList {
			item := itemList[j]
			itemInfoList = append(itemInfoList, *item.GetInfo())
		}
	}
	return itemInfoList
}

func (self *ItemMgr) Init(teamID int) bool {
	self.itemList = make(ItemList)
	//self.teamID = teamID ///存放自己的球队id
	itemListQuery := fmt.Sprintf("select * from %s where teamid=%d limit 3000", tableItem, teamID)
	itemInfo := new(ItemInfo)
	itemInfoList := GetServer().GetDynamicDB().fetchAllRows(itemListQuery, itemInfo)
	for i := range itemInfoList {
		itemInfo = itemInfoList[i].(*ItemInfo)
		self.itemList[itemInfo.ID] = NewItem(itemInfo)
	}
	return true
}

func (self *ItemMgr) RemoveItem(removeItemIDList IntList) bool { ///销毁道具列表
	removeItemIDListLen := len(removeItemIDList)
	if removeItemIDListLen <= 0 {
		return false ///空数列不处理
	}
	removeItemQuery := fmt.Sprintf("delete from %s where id in (", tableItem)
	for i := range removeItemIDList {
		removeItemID := removeItemIDList[i]
		if nil == self.itemList[removeItemID] {
			continue
		}
		delete(self.itemList, removeItemID)
		removeItemQuery += fmt.Sprintf("%d", removeItemID)
		if i < removeItemIDListLen-1 {
			removeItemQuery += ","
		}
	}
	removeItemQuery += fmt.Sprintf(") limit %d", removeItemIDListLen) ///限制只能删除指定数量记录
	///从数据库中删除对象
	_, rowsItemAffected := GetServer().GetDynamicDB().Exec(removeItemQuery)
	if rowsItemAffected != removeItemIDListLen {
		GetServer().GetLoger().Warn("ItemMgr RemoveItem removeStarQuery fail! itemID:%v", removeItemIDList)
		return false
	}
	return true
}

func (self *ItemMgr) GetEquipAttributeAddition(starIDList IntList) (float32, float32, float32) {
	attackscore := 0
	defensescore := 0
	organizescore := 0

	for i := range starIDList {
		starID := starIDList[i]
		itemList := self.GetStarItemSlice(starID)
		for j := range itemList {
			statictype := itemList[j].GetTypeInfo()
			attackscore += statictype.Attackscore
			defensescore += statictype.Defensescore
			organizescore += statictype.Organizescore
		}
	}
	return float32(attackscore), float32(defensescore), float32(organizescore)
}

func NewItemMgr(teamID int) IGameMgr {
	itemMgr := new(ItemMgr)
	if itemMgr.Init(teamID) == false {
		return nil
	}
	return itemMgr
}

//从Type得到item信息
func (self *ItemMgr) GetItemFromType(itemType int) *Item {
	for _, v := range self.itemList {
		itemInfo := v.GetInfo()
		if itemInfo.Type == itemType {
			return v
		}
	}
	return nil
}

//卸下该球员所有装备
func (self *ItemMgr) RemoveEquipment(client IClient, starID int) {
	itemSlice := self.GetStarItemSlice(starID)
	syncMgr := client.GetSyncMgr()
	for _, v := range itemSlice {
		v.Position = itemPosTeamEquip
		v.StarID = 0
		v.Cell = 0
		syncMgr.SyncObject("RemoveEquipment", v)
	}
}
