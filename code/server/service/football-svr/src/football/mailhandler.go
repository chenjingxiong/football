package football

import (
//"fmt"
//"time"
)

type NewMailGetMsg struct { ///新邮件消息
	MsgHead `json:"head"` ///"mail", "newmail"
	Mail    MailInfo      `json:"mailinfo"`
}

func (self *NewMailGetMsg) GetTypeAndAction() (string, string) {
	return "mail", "newmail"
}

type RemoveMailMsg struct { //通知客户端删除一封邮件
	MsgHead    `json:"head"` //"mail", "removemail"
	MailIDInfo IntList       `json:"mailidinfo"` // 删除邮件id int[]
}

func (self *RemoveMailMsg) GetTypeAndAction() (string, string) {
	return "mail", "removemail"
}

type MailQueryMsg struct { //球队邮件查询
	MsgHead `json:"head"` // "mail", "query"
}

func (self *MailQueryMsg) GetTypeAndAction() (string, string) {
	return "mail", "query"
}

func (self *MailQueryMsg) processAction(client IClient) bool {

	self.sendMailQueryMsg(client)
	return true
}

func (self *MailQueryMsg) sendMailQueryMsg(client IClient) {
	team := client.GetTeam()
	mailMgr := team.GetMailMgr()
	mailList := mailMgr.GetAllMail()
	msg := NewMailQueryResultMsg(mailList)
	client.SendMsg(msg)
}

func NewMailQueryResultMsg(mailList MailInfoList) *MailQueryResultMsg {
	msg := new(MailQueryResultMsg)
	msg.MailList = mailList
	return msg
}

type MailQueryResultMsg struct { //球队邮件查询结果
	MsgHead  `json:"head"` // "mail", "queryresult"
	MailList MailInfoList  `json:"maillist"`
}

func (self *MailQueryResultMsg) GetTypeAndAction() (string, string) {
	return "mail", "queryresult"
}

type MailAwardReceiveMsg struct {
	MsgHead `json:"head"` //"mail", "receiveaward"
	MailID  int           `json:"mailid"`
}

func (self *MailAwardReceiveMsg) GetTypeAndAction() (string, string) {
	return "mail", "receiveaward"
}

func (self *MailAwardReceiveMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	mailMgr := team.GetMailMgr()
	loger := GetServer().GetLoger()

	mail := mailMgr.GetMail(self.MailID)
	if loger.CheckFail("mail != nil", mail != nil, mail, nil) {
		return false //邮件不存在
	}

	awardTypeList := IntList{mail.AwardItem1, mail.AwardItem2, mail.AwardItem3, mail.AwardItem4, mail.AwardItem5}       //奖励物品类型
	awardCountList := IntList{mail.AwardCount1, mail.AwardCount2, mail.AwardCount3, mail.AwardCount4, mail.AwardCount5} //奖励物品数量

	if loger.CheckFail("awardTypeList[0] != 0", awardTypeList[0] != 0, awardTypeList[0], 0) {
		return false //邮件中不存在奖励
	}

	if loger.CheckFail("mail.State != StateReadAndGet", mail.State != StateReadAndGet, mail.State, StateReadAndGet) {
		return false //邮件已经读取并领取了所有道具
	}

	for i, v := range awardTypeList {
		if 0 == v {
			continue
		}

		//检查背包是否足够存放物品
		isStoreFull := team.IsStoreFull(v, awardCountList[i])
		if loger.CheckFail("isStoreFull == false", isStoreFull == false, isStoreFull, false) {
			return false //仓库已经满载
		}
	}

	return true
}

func (self *MailAwardReceiveMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	mailMgr := team.GetMailMgr()
	//领取单封邮件
	mail := mailMgr.GetMail(self.MailID)
	mail.State = StateReadAndGet                                                                                        //设置状态为已读且领取
	awardTypeList := IntList{mail.AwardItem1, mail.AwardItem2, mail.AwardItem3, mail.AwardItem4, mail.AwardItem5}       //奖励物品类型
	awardCountList := IntList{mail.AwardCount1, mail.AwardCount2, mail.AwardCount3, mail.AwardCount4, mail.AwardCount5} //奖励物品数量
	awardGradeList := IntList{mail.AwardGrade1, mail.AwardGrade2, mail.AwardGrade3, mail.AwardGrade4, mail.AwardGrade5} //奖励物品品质
	for i := range awardTypeList {
		if 0 == awardTypeList[i] {
			continue
		}
		team.AwardObject(awardTypeList[i], awardCountList[i], awardGradeList[i], 0)
		if awardTypeList[i] == awardTypeTicket && awardCountList[i] > 0 {
			client.RechargeRecord(Get_MailItem, awardCountList[i])
		}
	}
	return true

	//领取全部
	//for _, v := range mailMgr.mailList {
	//	mail := v.MailInfo
	//	if mail.State == StateReadAndGet {
	//		continue // 跳过已领取
	//	}

	//	if mail.AwardItem1 == 0 {
	//		continue // 跳过无奖励的邮件
	//	}
	//	mail.State = StateReadAndGet                                                                                        //设置邮件状态
	//	awardTypeList := IntList{mail.AwardItem1, mail.AwardItem2, mail.AwardItem3, mail.AwardItem4, mail.AwardItem5}       //奖励物品类型
	//	awardCountList := IntList{mail.AwardCount1, mail.AwardCount2, mail.AwardCount3, mail.AwardCount4, mail.AwardCount5} //奖励物品数量
	//	awardGradeList := IntList{mail.AwardGrade1, mail.AwardGrade2, mail.AwardGrade3, mail.AwardGrade4, mail.AwardGrade5} //奖励物品品质
	//	for i, j := range awardTypeList {
	//		if 0 == i {
	//			continue
	//		}
	//		team.AwardObject(j, awardCountList[i], awardGradeList[i], 0)
	//	}
	//}

	return true
}

func (self *MailAwardReceiveMsg) processAction(client IClient) bool {
	msgResult := msgResultFail
	if false == self.checkAction(client) {
		self.NewMailAwardReceiveResultMsg(client, &msgResult)
		return false
	}

	if false == self.doAction(client) {
		self.NewMailAwardReceiveResultMsg(client, &msgResult)
		return false
	}

	msgResult = msgResultOK
	self.NewMailAwardReceiveResultMsg(client, &msgResult)
	return true
}

type MailAwardReceiveResultMsg struct {
	MsgHead `json:"head"` //"mail", "receiveresult"
	Result  string
}

func (self *MailAwardReceiveResultMsg) GetTypeAndAction() (string, string) {
	return "mail", "receiveresult"
}

func (self *MailAwardReceiveMsg) NewMailAwardReceiveResultMsg(client IClient, result *string) {
	msg := new(MailAwardReceiveResultMsg)
	msg.Result = *result
	client.SendMsg(msg)
}

type MailReadMsg struct {
	MsgHead `json:"head"` //"mail", "read"
	MailID  int           `json:"id"` //查看邮件ID
}

func (self *MailReadMsg) GetTypeAndAction() (string, string) {
	return "mail", "read"
}

func (self *MailReadMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	mailMgr := team.GetMailMgr()
	loger := GetServer().GetLoger()

	mail := mailMgr.GetMail(self.MailID)
	if loger.CheckFail("mail != nil", mail != nil, mail, nil) {
		return false //邮件不存在
	}
	mailState := mail.State
	if loger.CheckFail("StateNewMail==mailState", StateNewMail == mailState,
		StateNewMail, mailState) {
		return false //邮件当前不是未读状态
	}
	return true
}

func (self *MailReadMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	mailMgr := team.GetMailMgr()

	mail := mailMgr.GetMail(self.MailID)
	mail.State = StateAlreadyRead
	return true
}

func (self *MailReadMsg) processAction(client IClient) bool {
	if false == self.checkAction(client) {
		return false
	}

	if false == self.doAction(client) {
		return false
	}

	return true
}

type MailDeleteMsg struct {
	MsgHead     `json:"head"`
	MailID      int  `json:"mailid"`      //删除邮件ID
	IsDeleteAll bool `json:"isdeleteall"` //是否删除所有   是/否(true/false)
}

func (self *MailDeleteMsg) GetTypeAndAction() (string, string) {
	return "mail", "delete"
}

func (self *MailDeleteMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	mailMgr := team.GetMailMgr()
	loger := GetServer().GetLoger()
	if true == self.IsDeleteAll {
		return true
	}
	if loger.CheckFail("self.MailID>0", self.MailID > 0, self.MailID, 0) {
		return false ///MailID必须大于0
	}
	mail := mailMgr.GetMail(self.MailID)
	if loger.CheckFail("mail != nil", mail != nil, mail, nil) {
		return false
	}
	return true
}

func (self *MailDeleteMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	mailMgr := team.GetMailMgr()
	if true == self.IsDeleteAll {
		mailMgr.RemoveAllMail()
		return true
	}

	mailMgr.RemoveMail(self.MailID)
	return true
}

func (self *MailDeleteMsg) processAction(client IClient) bool {
	if false == self.checkAction(client) {
		return false
	}

	if false == self.doAction(client) {
		return false
	}

	return true
}
