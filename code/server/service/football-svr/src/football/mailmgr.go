package football

import (
	"fmt"
	"reflect"
)

///单个球队所有拥有的邮件最大数为200,超过时自动删除过期邮件
///删除顺序为：已读已收附件——》未读已收附件——》已读未收附件——》未读未收附件
const (
	MaxMailCount = 200
)

//邮件状态
const (
	StateNewMail     = 1 ///未读
	StateAlreadyRead = 2 ///已读
	StateReadAndGet  = 3 ///已读且领取
)

//邮件类型
const (
	SystemMail       = 1 //系统邮件
	ReportMail       = 2 //战报邮件
	RemedyMail       = 3 ///补偿邮件
	ArenaFirstUpMail = 4 ///联赛首次晋级奖励
)

//发起方类型
const (
	ArenaSend        = 1 //竞技场系统发送
	CompensatingMail = 2 //补偿邮件
	LuckAwardMail    = 3 //奖励邮件
	ArenaRankUp      = 4 //竞技场晋级发送
)

type MailInfo struct {
	ID              int    ///对象编号
	TeamID          int    ///球队id
	Type            int    ///邮件类型 1系统消息 2战报
	Sort            int    ///邮件模板类型 1系统发物品奖励 2系统战报
	SenderType      int    ///邮件发起方类型,客户端对应名字
	State           int    ///邮件状态 1新邮 2已读 3已读并领取所有道具
	MakeTime        int    ///邮件生成时间
	TargetName      string ///对手名字
	TargetID        int    ///对方队徽
	TargetNpcID     int    ///对阵npc球队id
	TargetFormation int    ///对手阵型类型
	TargetTactic    int    ///对手战术类型
	HomeGoal        int    ///已方进球数
	TargetGoal      int    ///对方进球数
	HomeScore       int    ///已方战力
	TargetScore     int    ///对方战力

	AwardItem1  int ///道具类型
	AwardGrade1 int ///道具品质
	AwardCount1 int ///道具数量
	AwardItem2  int ///道具类型
	AwardGrade2 int ///道具品质
	AwardCount2 int ///道具数量
	AwardItem3  int ///道具类型
	AwardGrade3 int ///道具品质
	AwardCount3 int ///道具数量
	AwardItem4  int ///道具类型
	AwardGrade4 int ///道具品质
	AwardCount4 int ///道具数量
	AwardItem5  int ///道具类型
	AwardGrade5 int ///道具品质
	AwardCount5 int ///道具数量

	MailTitle string ///邮件标题
	MailText  string ///邮件内容

}

type Mail struct {
	MailInfo    ///邮件信息对象
	DataUpdater ///信息更新组件
}

type MailInfoList []MailInfo ///邮件列表
type MailList map[int]*Mail  ///邮件列表

func (self *Mail) GetReflectValue() reflect.Value { ///得到球队反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

func (self *Mail) GetID() int {
	return self.ID
}

///设置邮件中的奖品信息
func (self *Mail) SetAwardInfo(itemIndex int, awardItemType int, awardItemGrade int, awardItemCount int) {
	reflectValue := reflect.ValueOf(self).Elem()
	awardItemFieldName := fmt.Sprintf("AwardItem%d", itemIndex)   ///奖品类型
	awardGradeFieldName := fmt.Sprintf("AwardGrade%d", itemIndex) ///奖品品质
	awardCountFieldName := fmt.Sprintf("AwardCount%d", itemIndex) ///奖品数量
	awardItemField := reflectValue.FieldByName(awardItemFieldName)
	awardGradeField := reflectValue.FieldByName(awardGradeFieldName)
	awardCountField := reflectValue.FieldByName(awardCountFieldName)
	awardItemField.SetInt(int64(awardItemType))
	awardGradeField.SetInt(int64(awardItemGrade))
	awardCountField.SetInt(int64(awardItemCount))
}

type MailMgr struct {
	GameMgr  ///游戏管理器基类
	mailList MailList
}

func (self *MailMgr) GetType() int { ///得到管理器类型
	return mgrTypeMailMgr ///关卡管理器
}

func (self *MailMgr) SaveInfo() { ///保存数据
	for _, v := range self.mailList {
		v.Save()
	}
}

func NewMailMgr(teamID int) IGameMgr { ///邮件管理器
	mailMgr := new(MailMgr)
	mailMgr.mailList = make(MailList) ///创建邮件列表
	mailListQuery := fmt.Sprintf("select * from %s where teamid=%d limit %d", tableMail, teamID, MaxMailCount)
	mailInfo := new(MailInfo)
	mailInfoList := GetServer().GetDynamicDB().fetchAllRows(mailListQuery, mailInfo)
	for i := range mailInfoList {
		mailInfo = mailInfoList[i].(*MailInfo)
		mail := new(Mail)
		mail.MailInfo = *mailInfo
		mail.InitDataUpdater(tableMail, &mail.MailInfo)
		mailMgr.mailList[mail.ID] = mail
	}
	return mailMgr
}

func (self *MailMgr) createNewMail(newMail *Mail) { ///向数据库中插入一件新邮件记录
	//createMailQuery := fmt.Sprintf(`insert %s set teamid=%d,type=%d,sort=%d,sendertype=%d,state=%d,maketime=%d,
	//targetname='%s',targetid=%d,targetnpcid=%d,targetformation=%d,targettactic=%d,
	//`
	//tableItem, self.team.GetID(), itemType, itemCount, awardItemType.Color, itemPosType) ///组插入记录SQL
	//lastInsertItemID, _ := GetServer().GetDynamicDB().Exec(awardItemQuery)
	//if lastInsertItemID <= 0 {
	//GetServer().GetLoger().Warn("ItemMgr AwardItem fail! itemType:%d itemCount:%d", itemType, itemCount)
	//return 0
	//}
	newMail.InitDataUpdater(tableMail, &newMail.MailInfo)
	insertSql := newMail.InsertSql()
	lastInsertItemID, _ := GetServer().GetDynamicDB().Exec(insertSql)
	if lastInsertItemID <= 0 {
		GetServer().GetLoger().Warn("MailMgr createNewMail is fail! teamid: %d", self.team.GetID())
		return
	}
	newMail.ID = lastInsertItemID ///更新mail id
	newMail.InitDataUpdater(tableMail, &newMail.MailInfo)
	//加入内存
	if self.mailList == nil {
		self.mailList = make(MailList)
	}
	self.mailList[newMail.ID] = newMail

	//通知客户端新邮件
	self.syncMgr.SyncNewMail(newMail.MailInfo)
}

///发送系统奖励邮件,共指定三个属性列表,类型,品质(可以为0),数量
func (self *MailMgr) SendSysAwardMail(mailSort int, senderType int, awardItemList IntList, awardGradeList IntList,
	awardCountList IntList, mailtitle string, mailtext string) {
	mail := new(Mail)
	mail.TeamID = self.team.GetID()
	mail.Type = SystemMail
	mail.Sort = mailSort
	mail.State = StateNewMail
	mail.MakeTime = Now()
	mail.SenderType = senderType
	mail.MailTitle = mailtitle
	mail.MailText = mailtext
	//awardItemListLen := awardItemList.Len()
	//awardGradeListLen := awardGradeList.Len()
	//awardCountListLen := awardCountList.Len()
	//isVaildListLen := true
	//if awardItemListLen != awardGradeListLen || awardGradeListLen != awardCountListLen || awardCountListLen != awardItemListLen {
	//	return ///无效的列表数量,它们中有不匹配的
	//}
	for i := range awardItemList {
		awardItemType := awardItemList[i]
		awardItemGrade := awardGradeList[i]
		awardItemCount := awardCountList[i]
		mail.SetAwardInfo(i+1, awardItemType, awardItemGrade, awardItemCount)
	}
	self.createNewMail(mail)
}

///发送系统战报
//func (self *MailMgr) SendMatchReport(mailSort int, senderType int, targetName string, targetID int, targetNpcID int,
//	targetFormation int, targetTactic int, homeGoal int, targetGoal int, homeScore int, targetScore int) {
//	mail := new(Mail)
//	mail.Sort = mailSort
//	mail.Type = ReportMail
//	mail.TeamID = self.team.GetID()
//	mail.SenderType = senderType
//	mail.MakeTime = Now()
//	mail.State = StateNewMail
//	mail.TargetName = targetName
//	mail.TargetID = targetID
//	mail.TargetNpcID = targetNpcID
//	mail.TargetFormation = targetFormation
//	mail.TargetTactic = targetTactic
//	mail.HomeGoal = homeGoal
//	mail.TargetGoal = targetGoal
//	mail.HomeScore = homeScore
//	mail.TargetScore = targetScore
//	self.createNewMail(mail)
//}

///得到玩家所有邮件
func (self *MailMgr) GetAllMail() MailInfoList {
	var mailInfoList = MailInfoList{}
	for _, v := range self.mailList {
		mailInfoList = append(mailInfoList, v.MailInfo)
	}
	return mailInfoList
}

///得到玩家一封邮件
func (self *MailMgr) GetMail(mailID int) *Mail {
	mail := self.mailList[mailID]
	return mail
}

func (self *MailMgr) deleteMail(mailID int) bool {
	deleteSql := fmt.Sprintf("delete from %s where id = %d", tableMail, mailID)
	_, rowsAffected := GetServer().GetDynamicDB().Exec(deleteSql)
	if rowsAffected <= 0 {
		GetServer().GetLoger().Warn("MailMgr deleteMail is fail! teamid: %d mailID: %d", self.team.GetID(), mailID)
		return false
	}
	delete(self.mailList, mailID) //删除内存
	return true
}

///删除玩家一封邮件
func (self *MailMgr) RemoveMail(mailID int) {
	mail := self.mailList[mailID]
	if nil == mail {
		return
	}

	self.deleteMail(mailID)
	self.syncMgr.SyncRemoveMail(IntList{mailID})
	return
}

///删除玩家所有邮件
func (self *MailMgr) RemoveAllMail() {
	removeList := IntList{}
	for i, v := range self.mailList {
		mail := v
		if mail.State == StateNewMail {
			continue //不删除新邮件
		}

		if mail.AwardItem1 != 0 && mail.State == StateAlreadyRead {
			continue //不删除尚未领取奖励的邮件
		}
		self.deleteMail(mail.GetID())
		removeList = append(removeList, i)
	}

	if removeList.Len() > 0 {
		self.syncMgr.SyncRemoveMail(removeList)
	}
}

//func NewStar(starInfo *StarInfo, team *Team) *Star {
//	star := new(Star)
//	star.StarInfo = *starInfo
//	star.InitDataUpdater(tableStar, &star.StarInfo)
//	star.team = team
//	return star
//}

//const ( ///ErrorType
//	failLogin            = "loginfail"            ///登录失败
//	failCreateTeam       = "createteamfail"       ///创建队伍失败
//	failStarSpyDiscover  = "starspydiscoverfail"  ///球探发掘球员失败
//	failFormationUplevel = "formationuplevelfail" ///球队阵形升级失败
//	failChatWhispe       = "chatwhispefail"       ///私聊失败
//)

//const ( ///ErrorDesc
//	failAccountNotExsit           = "AccountNotExsit"           ///帐户不存在
//	failPasswordWrong             = "PasswordWrong"             ///密码错误
//	failSameName                  = "SameName"                  ///同名冲突
//	failInvalidName               = "InvalidName"               ///名字非法
//	failInvalidAccountID          = "InvalidAccountID"          ///名字非法
//	failInvalidMsg                = "InvalidMsg"                ///消息非法,非法的发送时机
//	failInvalidParam              = "InvalidParam"              ///消息非法,非法的发送时机
//	failReachLimit                = "ReachLimit"                ///超过限制
//	failInsufficientTicket        = "InsufficientTicket"        ///球票不足
//	failInsufficientDiscoverCount = "InsufficientDiscoverCount" ///发掘次数不足
//	failInreachmaxlevel           = "ReachMaxLevel"             ///超过等级上限
//	failNotFound                  = "NotFound"                  ///找不到对象
//)

//type ActionErrorMsg struct { ///错误提示消息
//	MsgHead   `json:"head"`
//	ErrorType string `json:"errortype"` ///错误类型
//	ErrorDesc string `json:"errordesc"` ///错误描述
//}

//func (self *ActionErrorMsg) GetTypeAndAction() (string, string) {
//	return "action", "error"
//}

//func NewActionErrorMsg(errorType *string, errorDesc *string) *ActionErrorMsg {
//	msg := new(ActionErrorMsg)
//	//msg.MsgType = "action"
//	//msg.Action = "error"
//	//msg.CreateTime = int(time.Now().Unix())
//	msg.ErrorType = *errorType
//	msg.ErrorDesc = *errorDesc
//	return msg
//}

//type ActionHandler struct {
//	MsgHandler
//}

//func (self *ActionErrorMsg) New() interface{} {
//	return new(ActionErrorMsg)
//}

//func (self *ActionErrorMsg) getName() string {
//	return "error"
//}

//func (self *ActionErrorMsg) processAction(client IClient) bool {
//	//actionErrorMsg := msg.(*ActionErrorMsg)
//	GetServer().GetLoger().Info("ActionMsg processAction msg:%v", self)
//	return false
//}

//func (self *ActionHandler) getName() string { ///返回可处理的消息类型
//	return "action"
//}

//func (self *ActionHandler) initHandler() { ///初始化消息处理器
//	//self.addActionToList(new(ActionErrorMsg))
//}
