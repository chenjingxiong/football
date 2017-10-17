package football

import (
	"fmt"
	"reflect"
)

type SyncObjectList []ISyncObject

type ISyncObject interface {
	GetReflectValue() reflect.Value ///得到反射对象
	Sync() DataValueChangList       ///得到属性变更列表
	GetID() int                     ///得到ID
}

const ( ///SystemType
	systemTypeStarSpyDiscover    = "StarSpyDiscover"
	systemTypeStarCenterTransfer = "StarCenterTransfer"
	systemTypeStarTrain          = "StarTrain"
	systemTypeStarVolunteer      = "StarVolunteer"
	systemTypeItemEquip          = "ItemEquip"
	systemTypeSkillEquip         = "SkillEquip"
	systemTypeItemMerge          = "ItemMerge"
	systemTypeCommon             = "Common"
)

const ( //attribtype
	primerStarSpyLuckAttribType                = "PrimerStarSpyLuck"                ///初级球探幸运
	middleStarSpyLuckAttribType                = "MiddleStarSpyLuck"                ///中级球探幸运
	expertStarSpyLuckAttribType                = "ExpertStarSpyLuck"                ///高级球探幸运
	primerStarSpyCDAttribType                  = "PrimerStarSpyCD"                  ///初级球探cd
	middleStarSpyCDAttribType                  = "MiddleStarSpyCD"                  ///中级球探cd
	expertStarSpyCDAttribType                  = "ExpertStarSpyCD"                  ///高级球探cd
	primerStarSpyDiscoverRemainCountAttribType = "PrimerStarSpyDiscoverRemainCount" ///初级球探发掘剩余次数
	middleStarSpyDiscoverRemainCountAttribType = "MiddleStarSpyDiscoverRemainCount" ///中级球探发掘剩余次数
	expertStarSpyDiscoverRemainCountAttribType = "ExpertStarSpyDiscoverRemainCount" ///高级球探发掘剩余次数
	teamTicketAttribType                       = "TeamTicket"                       ///球队球票
	teamCoinAttribType                         = "TeamCoin"                         ///球队金币
	teamTrainCellAttribType                    = "TeamTrainCell"                    ///球队训练位
	starExpAttribType                          = "StarExp"                          ///球员经验
	itemStarIDOwnAttribType                    = "ItemStarIDOwn"                    ///道具球星主人id
	itemMergeExpAttribType                     = "ItemMergeExp"                     ///道具当前融合经验值
	itemCellAttribType                         = "ItemCell"                         ///道具所在单元格索引号
	skillStarIDOwnAttribType                   = "SkillStarIDOwn"                   ///技能球星主人id
	skillPositionAttribType                    = "SkillPosition"                    ///技能所在装备位置
)

///同步管理器,专门用于将服务器最新的数据同步给客户端

type AttribChangeItem struct { ///属性变更项
	AttribType  string `json:"attribtype"`  ///属性类型
	AttribID    int    `json:"attribid"`    ///属性ID
	AttribValue string `json:"attribvalue"` ///属性值
}

type AttribChangeList []AttribChangeItem
type ActionAttribChangeMsg struct { ///属性变更消息
	MsgHead        `json:"head"`
	SystemType     string           `json:"systemtype"`
	AttribItemList AttribChangeList `json:"attribitemlist"`
}

func (self *ActionAttribChangeMsg) GetAttribChangeListSize() int {
	return len(self.AttribItemList) ///返回消息中已放入的属性列表长度
}

func (self *ActionAttribChangeMsg) AddAttribChange(attribType string, attribID int, attribValue int) {
	attribValueStr := fmt.Sprintf("%d", attribValue)
	self.AddAttribChangeStr(attribType, attribID, attribValueStr)
}

func (self *ActionAttribChangeMsg) AddAttribChangeStr(attribType string, attribID int, attribValue string) {
	attribChangeItem := new(AttribChangeItem)
	attribChangeItem.AttribType = attribType
	attribChangeItem.AttribID = attribID
	attribChangeItem.AttribValue = attribValue
	self.AttribItemList = append(self.AttribItemList, *attribChangeItem)
}

func (self *ActionAttribChangeMsg) AddAttribChangeInt(attribType string, attribValue int) {
	self.AddAttribChange(attribType, 0, attribValue)
}

func (self *ActionAttribChangeMsg) GetTypeAndAction() (string, string) {
	return "action", "attribchange"
}

func (self *SyncObjectList) Len() int {
	return len(*self)
}

type StarCalcInfoMsg struct { ///球员的二级属性
	MsgHead      `json:"head"` ///"star", "starcalcinfo"
	StarID       int           `json:"starid"`
	StarCalcInfo StarInfoCalc  `json:"starcalcinfo"`
}

func (self *StarCalcInfoMsg) GetTypeAndAction() (string, string) {
	return "star", "starcalcinfo"
}

func NewActionAttribChangeMsg(systemType string, attribType string, attribID int, attribValue int) *ActionAttribChangeMsg {
	msg := new(ActionAttribChangeMsg)
	msg.SystemType = systemType
	if attribValue > 0 {
		attribChangeItem := new(AttribChangeItem)
		attribChangeItem.AttribType = attribType
		attribChangeItem.AttribID = attribID
		attribChangeItem.AttribValue = fmt.Sprintf("%d", attribValue)
		msg.AttribItemList = append(msg.AttribItemList, *attribChangeItem)
	}
	return msg
}

type ItemRemoveMsg struct { ///道具删除消息
	MsgHead    `json:"head"`
	ItemIDList IntList `json:"itemidlist"` ///删除道具id列表
}

func (self *ItemRemoveMsg) GetTypeAndAction() (string, string) {
	return "item", "itemremove"
}

type TaskAddMsg struct { ///任务添加消息
	MsgHead  `json:"head"` ///"task", "taskadd"
	TaskList TaskInfoList  `json:"tasklist"` ///任务信息列表
}

func (self *TaskAddMsg) GetTypeAndAction() (string, string) {
	return "task", "taskadd"
}

type LevelAddMsg struct { ///关卡添加消息
	MsgHead   `json:"head"`
	LevelList LevelInfoList `json:"levellist"` ///关卡信息列表
}

func (self *LevelAddMsg) GetTypeAndAction() (string, string) {
	return "level", "leveladd"
}

type FormationAddMsg struct { ///道具添加消息
	MsgHead       `json:"head"`     ///"formation", "formationadd"
	FormationList FormationInfoList `json:"skilllist"` ///阵形信息列表
}

func (self *FormationAddMsg) GetTypeAndAction() (string, string) {
	return "formation", "formationadd"
}

type SkillAddMsg struct { ///道具添加消息
	MsgHead   `json:"head"` ///"skill", "skilladd"
	SkillList SkillInfoList `json:"skilllist"` ///道具列表
}

func (self *SkillAddMsg) GetTypeAndAction() (string, string) {
	return "skill", "skilladd"
}

type ItemAddMsg struct { ///道具添加消息
	MsgHead  `json:"head"` ///"item", "itemadd"
	ItemList ItemInfoList  `json:"itemlist"` ///道具列表
}

func (self *ItemAddMsg) GetTypeAndAction() (string, string) {
	return "item", "itemadd"
}

type StarAddMsg struct { ///球员添加消息
	MsgHead  `json:"head"` ///"star", "staradd"
	StarList []StarInfo    `json:"starlist"` ///球员列表
}

func (self *StarAddMsg) GetTypeAndAction() (string, string) {
	return "star", "staradd"
}

func (self *StarAddMsg) HasMember() bool {
	hasMember := len(self.StarList) > 0
	return hasMember
}

func NewStarAddMsg(starInfo *StarInfo) *StarAddMsg {
	msg := new(StarAddMsg)
	msg.AddMember(starInfo)
	return msg
}

func (self *StarAddMsg) AddMember(starInfo *StarInfo) {
	self.StarList = append(self.StarList, *starInfo)
}

//<<<<<<< .mine
////type ISyncMgr interface {
////	syncAddStar(starIDList IntList)                                                           ///同步客户端一组球员添加信息
////	syncAttribChangeList(systemType string, attribChangeList AttribChangeList)                ///同步客户端一组属性变更消息
////	syncAttribChangeItem(systemType string, attribType string, attribID int, AttribValue int) ///同步客户端一个属性变更消息
////	//SyncItem(systemType string, itemList ...IGameItem)                                        ///同步客户端一个道具属性变更消息
////	SyncSkill(systemType string, itemList ...ISkill)                     ///同步客户端一个技能属性变更消息
////	SyncRemoveItem(removeItemIDList IntList)                             ///同步客户端道具删除消息
////	SyncRemoveStarCenterMember(starCenterType int, memberIDList IntList) ///同步客户端球员中心删除成员消息
////	SyncObject(systemType string, syncObject ISyncObject)                ///同步属性变更消息
////	syncAddLevel(levelIDList IntList)                                    ///同步客户端一组关卡添加信息
////	syncAddItem(itemIDList IntList)                                      ///同步客户端一组道具添加信息
////	syncAddSkill(skillDList IntList)                                     ///同步客户端一组技能添加信息
////	syncAddFormation(formationIDList IntList)                            ///同步客户端一组阵形添加信息
////	SyncAddTask(taskIDList IntList)                                      ///同步客户端一组关卡添加信息
////	SyncObjectArray(systemType string, syncObjectList SyncObjectList)    ///同步对象列表属性变更消息
////	syncRemoveStar(starIDList IntList)                                   ///同步客户端一组球员删除信息
////	syncStarCalcInfo(star IStar)                                         ///同步球员的二级属性给客户端用于调试
////	Init(client IClient) bool                                            ///同步客户端一组球员添加信息
////}
//=======
//type ISyncMgr interface {
//	syncAddStar(starIDList IntList)                                                           ///同步客户端一组球员添加信息
//	syncAttribChangeList(systemType string, attribChangeList AttribChangeList)                ///同步客户端一组属性变更消息
//	syncAttribChangeItem(systemType string, attribType string, attribID int, AttribValue int) ///同步客户端一个属性变更消息
//	//SyncItem(systemType string, itemList ...IGameItem)                                        ///同步客户端一个道具属性变更消息
//	SyncSkill(systemType string, itemList ...ISkill)                     ///同步客户端一个技能属性变更消息
//	SyncRemoveItem(removeItemIDList IntList)                             ///同步客户端道具删除消息
//	SyncRemoveStarCenterMember(starCenterType int, memberIDList IntList) ///同步客户端球员中心删除成员消息
//	SyncObject(systemType string, syncObject ISyncObject)                ///同步属性变更消息
//	syncAddLevel(levelIDList IntList)                                    ///同步客户端一组关卡添加信息
//	syncAddItem(itemIDList IntList)                                      ///同步客户端一组道具添加信息
//	syncAddSkill(skillDList IntList)                                     ///同步客户端一组技能添加信息
//	syncAddFormation(formationIDList IntList)                            ///同步客户端一组阵形添加信息
//	SyncAddTask(taskIDList IntList)                                      ///同步客户端一组关卡添加信息
//	SyncObjectArray(systemType string, syncObjectList SyncObjectList)    ///同步对象列表属性变更消息
//	syncRemoveStar(starIDList IntList)                                   ///同步客户端一组球员删除信息
//	syncStarCalcInfo(star IStar)                                         ///同步球员的二级属性给客户端用于调试
//	SyncRemoveVipShopCommodity(commodityIDList IntList, changeType bool) ///同步客户端商城下架商品信息
//	Init(client IClient) bool                                            ///同步客户端一组球员添加信息
//}
//>>>>>>> .r1883

type SyncMgr struct {
	client IClient
	team   *Team
}

func (self *SyncMgr) SyncRemoveStarCenterMember(starCenterType int, memberIDList IntList) { ///同步客户端球员中心删除成员消息
	starCenterMemberRemoveMsg := new(StarCenterMemberRemoveMsg)
	starCenterMemberRemoveMsg.StarCenterType = starCenterType
	starCenterMemberRemoveMsg.MemberIDList = memberIDList
	self.client.SendMsg(starCenterMemberRemoveMsg)
}

func (self *SyncMgr) SyncRemoveItem(removeItemIDList IntList) { ///同步客户端道具删除消息
	itemRemoveMsg := new(ItemRemoveMsg)
	itemRemoveMsg.ItemIDList = removeItemIDList
	self.client.SendMsg(itemRemoveMsg)
}

func (self *SyncMgr) SyncNewMail(newMail MailInfo) {
	newMailGetMsg := new(NewMailGetMsg)
	newMailGetMsg.Mail = newMail
	self.client.SendMsg(newMailGetMsg)
}

func (self *SyncMgr) SyncRemoveMail(mailIDList IntList) {
	removeMailMsg := new(RemoveMailMsg)
	removeMailMsg.MailIDInfo = mailIDList
	self.client.SendMsg(removeMailMsg)
}

func (self *SyncMgr) SyncObjectArray(systemType string, syncObjectList SyncObjectList) { ///同步属性变更消息
	for i := range syncObjectList {
		syncObject := syncObjectList[i]
		if nil == syncObject {
			continue
		}
		self.SyncObject(systemType, syncObject)
	}
}

func (self *SyncMgr) SyncObject(systemType string, syncObject ISyncObject) { ///同步属性变更消息
	if nil == syncObject {
		return
	}
	objectAttribChangeMsg := NewActionAttribChangeMsg(systemType, "", 0, 0)
	dataValueChangList := syncObject.Sync()
	reflectValue := syncObject.GetReflectValue()
	for i := range dataValueChangList {
		fieldName := dataValueChangList[i]
		changFieldValue := reflectValue.FieldByName(fieldName)
		attribName := fmt.Sprintf("%s%sAttribType", reflectValue.Type().Name(), fieldName)
		switch changFieldValue.Kind() {
		case reflect.Int:
			changFieldValueInt := int(changFieldValue.Int())
			objectAttribChangeMsg.AddAttribChange(attribName, syncObject.GetID(), changFieldValueInt)
		case reflect.String:
			changFieldValueStr := changFieldValue.String()
			objectAttribChangeMsg.AddAttribChangeStr(attribName, syncObject.GetID(), changFieldValueStr)
		}
	}
	if objectAttribChangeMsg.GetAttribChangeListSize() > 0 {
		self.client.SendMsg(objectAttribChangeMsg)
	}
}

// func (self *SyncMgr) SyncSkill(systemType string, skillList ...ISkill) { ///同步客户端一个技能属性变更消息
// 	skillAttribChangeMsg := NewActionAttribChangeMsg(systemType, "", 0, 0)
// 	for k := range skillList {
// 		if nil == skillList[k] {
// 			continue
// 		}
// 		skill := skillList[k]
// 		dataValueChangList := skill.Sync()
// 		skillInfo := skill.GetInfo()
// 		for i := range dataValueChangList {
// 			fieldName := dataValueChangList[i]
// 			switch fieldName {
// 			case "StarID":
// 				skillAttribChangeMsg.AddAttribChange(skillStarIDOwnAttribType, skillInfo.ID, skillInfo.StarID)
// 			case "Position":
// 				skillAttribChangeMsg.AddAttribChange(skillPositionAttribType, skillInfo.ID, skillInfo.Position)
// 			}
// 		}
// 	}
// 	if skillAttribChangeMsg.GetAttribChangeListSize() > 0 {
// 		self.client.SendMsg(skillAttribChangeMsg)
// 	}
// }

//func (self *SyncMgr) SyncItem(systemType string, itemList ...IGameItem) { ///同步客户端一个道具属性变更消息
//	itemAttribChangeMsg := NewActionAttribChangeMsg(systemType, "", 0, 0)
//	for k := range itemList {
//		if nil == itemList[k] {
//			continue
//		}
//		item := itemList[k]
//		dataValueChangList := item.Sync()
//		itemInfo := item.GetInfo()
//		for i := range dataValueChangList {
//			fieldName := dataValueChangList[i]
//			switch fieldName {
//			case "StarID":
//				itemAttribChangeMsg.AddAttribChange(itemStarIDOwnAttribType, itemInfo.ID, itemInfo.StarID)
//			case "MergeExp":
//				itemAttribChangeMsg.AddAttribChange(itemMergeExpAttribType, itemInfo.ID, itemInfo.MergeExp)
//			case "Cell":
//				itemAttribChangeMsg.AddAttribChange(itemMergeExpAttribType, itemInfo.ID, itemInfo.MergeExp)
//			}
//		}
//	}
//	if itemAttribChangeMsg.GetAttribChangeListSize() > 0 {
//		self.client.SendMsg(itemAttribChangeMsg)
//	}
//}

func NewSyncMgr(client IClient) *SyncMgr {
	syncMgr := new(SyncMgr)
	syncMgr.Init(client)
	return syncMgr
}

func (self *SyncMgr) Init(client IClient) bool { ///同步客户端一组球员添加信息
	self.client = client
	self.team = client.GetTeam()
	return true
}

func (self *SyncMgr) syncStarCalcInfo(star *Star) { ///同步球员的二级属性给客户端用于调试
	starCalcInfoMsg := new(StarCalcInfoMsg)
	starCalcInfoMsg.StarID = star.GetID()
	starCalcInfoMsg.StarCalcInfo = *star.GetCalcInfo()
	self.client.SendMsg(starCalcInfoMsg)
}

func (self *SyncMgr) syncAttribChangeList(systemType string, attribChangeList AttribChangeList) { ///同步客户端一组属性变更消息
	actionAttribChangeMsg := new(ActionAttribChangeMsg)
	actionAttribChangeMsg.SystemType = systemType
	actionAttribChangeMsg.AttribItemList = attribChangeList
	self.client.SendMsg(actionAttribChangeMsg)
}

func (self *SyncMgr) syncAttribChangeItem(systemType string, attribType string, attribID int, attribValue int) { ///同步客户端一个属性变更消息
	actionAttribChangeMsg := new(ActionAttribChangeMsg)
	actionAttribChangeMsg.SystemType = systemType
	attribValueString := fmt.Sprintf("%d", attribValue)
	actionAttribChangeMsg.AttribItemList = AttribChangeList{{attribType, attribID, attribValueString}}
	self.client.SendMsg(actionAttribChangeMsg)
}

func (self *SyncMgr) syncAddFormation(formationIDList IntList) { ///同步客户端一组阵形添加信息
	syncAddFormationMsg := new(FormationAddMsg)
	for i := range formationIDList {
		formationID := formationIDList[i]
		formation := self.team.GetFormationMgr().GetFormation(formationID)
		if nil == formation {
			continue
		}
		formationInfoPtr := formation.GetInfo()
		syncAddFormationMsg.FormationList = append(syncAddFormationMsg.FormationList, *formationInfoPtr)
	}
	if len(syncAddFormationMsg.FormationList) > 0 {
		self.client.SendMsg(syncAddFormationMsg) ///通知客户端新添的关卡信息列表
	}
}

func (self *SyncMgr) syncAddItem(itemIDList IntList) { ///同步客户端一组道具添加信息
	syncAddItemMsg := new(ItemAddMsg)
	for i := range itemIDList {
		itemID := itemIDList[i]
		item := self.team.GetItemMgr().GetItem(itemID)
		if nil == item {
			continue
		}
		itemInfoPtr := item.GetInfo()
		syncAddItemMsg.ItemList = append(syncAddItemMsg.ItemList, *itemInfoPtr)
	}
	if len(syncAddItemMsg.ItemList) > 0 {
		self.client.SendMsg(syncAddItemMsg) ///通知客户端新添的关卡信息列表
	}
}

func (self *SyncMgr) syncAddSkill(skillDList IntList) { ///同步客户端一组技能添加信息
	syncAddSkillMsg := new(SkillAddMsg)
	for i := range skillDList {
		skillID := skillDList[i]
		skill := self.team.GetSkillMgr().GetSkill(skillID)
		if nil == skill {
			continue
		}
		skillInfoPtr := skill.GetInfo()
		syncAddSkillMsg.SkillList = append(syncAddSkillMsg.SkillList, *skillInfoPtr)
	}
	if len(syncAddSkillMsg.SkillList) > 0 {
		self.client.SendMsg(syncAddSkillMsg) ///通知客户端新添的技能信息列表
	}
}

func (self *SyncMgr) SyncAddTask(taskIDList IntList) { ///同步客户端一组关卡添加信息
	syncAddTaskMsg := new(TaskAddMsg)
	for i := range taskIDList {
		taskID := taskIDList[i]
		task := self.team.GetTaskMgr().GetTask(taskID)
		if nil == task {
			continue
		}
		taskInfoPtr := task.GetInfoPtr()
		syncAddTaskMsg.TaskList = append(syncAddTaskMsg.TaskList, *taskInfoPtr)
	}
	if len(syncAddTaskMsg.TaskList) > 0 {
		self.client.SendMsg(syncAddTaskMsg) ///通知客户端新添的关卡信息列表
	}
}

func (self *SyncMgr) syncAddLevel(levelIDList IntList) { ///同步客户端一组关卡添加信息
	syncAddLevelMsg := new(LevelAddMsg)
	for i := range levelIDList {
		levelID := levelIDList[i]
		level := self.team.GetLevelMgr().GetLevel(levelID)
		if nil == level {
			continue
		}
		levelInfoPtr := level.GetInfoPtr()
		syncAddLevelMsg.LevelList = append(syncAddLevelMsg.LevelList, *levelInfoPtr)
	}
	if len(syncAddLevelMsg.LevelList) > 0 {
		self.client.SendMsg(syncAddLevelMsg) ///通知客户端新添的关卡信息列表
	}
}

func (self *SyncMgr) syncAddStar(starIDList IntList) { ///同步客户端一组球员添加信息
	syncAddStarMsg := new(StarAddMsg)
	for i := range starIDList {
		starID := starIDList[i]
		star := self.team.GetStar(starID)
		if nil == star {
			continue
		}
		starInfo := star.GetInfo()
		syncAddStarMsg.StarList = append(syncAddStarMsg.StarList, *starInfo)
	}
	if len(syncAddStarMsg.StarList) > 0 {
		self.client.SendMsg(syncAddStarMsg) ///通知客户端新添了一个球员
		self.team.CalcScore()               ///添加一个球员成功后重新计算球队评分
		self.SyncObject("syncAddStar", self.team)
	}
}

type StarRemoveMsg struct { ///球员删除消息
	MsgHead    `json:"head"`
	StarIDList []int `json:"starlist"` ///球员id列表
}

func (self *StarRemoveMsg) GetTypeAndAction() (string, string) {
	return "star", "starremove"
}

func NewStarRemoveMsg(starID int) *StarRemoveMsg {
	msg := new(StarRemoveMsg)
	if starID > 0 {
		msg.AddMember(starID)
	}
	return msg
}

func (self *StarRemoveMsg) AddMember(starID int) {
	self.StarIDList = append(self.StarIDList, starID)
}

func (self *StarRemoveMsg) GetMemberSize() int {
	memberSize := len(self.StarIDList)
	return memberSize
}

func (self *SyncMgr) syncRemoveStar(starIDList IntList) { ///同步客户端一组球员删除信息
	syncRemoveStarMsg := new(StarRemoveMsg)
	for i := range starIDList {
		starID := starIDList[i]
		syncRemoveStarMsg.StarIDList = append(syncRemoveStarMsg.StarIDList, starID)
	}
	if len(syncRemoveStarMsg.StarIDList) > 0 {
		self.client.SendMsg(syncRemoveStarMsg) ///通知客户端删除一组球员
	}
}

func (self *SyncMgr) syncAllStars() { ///同步玩家所有球员属性信息
	allStarList := self.team.GetAllStarList()
	for index := range allStarList {
		starID := allStarList[index]
		star := self.team.GetStar(starID)
		star.CalcScore()
		self.SyncObject("syncAllStars", star)
	}
}

func (self *SyncMgr) syncMainStars() { ///同步客户端主力上场球员属性信息
	mainStarsList := self.team.GetMainStarsList()
	for index := range mainStarsList {
		star := mainStarsList[index]
		self.SyncObject("syncMainStars", star)
	}
}
