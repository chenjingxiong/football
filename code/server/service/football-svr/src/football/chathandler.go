package football

//import (
////	"encoding/json"
////"log"
//)

import (
	"os/exec"
	"strconv"
	"strings"
)

const (
	chatTypeSystem = 1 ///系统聊天
	chatTypeWhispe = 2 ///私聊
	chatTypeWorld  = 3 ///世界聊天
	chatGMKick     = 4 ///GM踢人
)

//type ChatHandler struct {
//	MsgHandler
//}

//func (self *ChatHandler) initHandler() { ///初始化消息处理器
//	self.addActionToList(new(ChatMsg))
//}

type ChatMsg struct { ///私聊消息
	MsgHead  `json:"head"`
	Sender   string `json:"sender"`   ///发送者名字
	Receiver string `json:"receiver"` ///接收者名字
	Type     int    `json:"type"`     ///频道类型
	Text     string `json:"text"`     ///正文
}

func (self *ChatMsg) GetTypeAndAction() (string, string) {
	return "chat", "chat"
}

func (self *ChatMsg) checkAction(client IClient) bool { ///通用检测
	return true
}

func (self *ChatMsg) processWhispe(client IClient) { ///处理私聊逻辑
	client.BroadcastMsg(self) ///广播自己
}

func (self *ChatMsg) processWorld(client IClient) { ///处理世界聊天逻辑
	client.BroadcastMsg(self) ///广播自己
}

func (self *ChatMsg) GMCommand(client IClient) bool { ///检测是否为GM命令
	team := client.GetTeam()
	itemMgr := team.GetItemMgr()
	loger := GetServer().GetLoger()
	//	gmMgr := GetServer().GetGMMgr()
	sync := client.GetSyncMgr()

	isOpenGMTalk := GetServer().IsOpenGMTalk()
	if isOpenGMTalk == false {
		return false ///如果没有开启接受gm指令的功能则忽略gm指令解析过程
	}
	///实现GM命令
	///1.格式:greedisgood_10000  加钱
	// isGMCommand := gmMgr.IsGMCommand(self.Text)
	// if isGMCommand == false {
	// 	return false
	// }
	isGM := team.IsGM()
	if loger.CheckFail("isGM==true", isGM == true, isGM, true) {
		return false //Gm指令只限GM号使用
	}

	switch {
	case strings.Index(self.Text, "@Money") != -1:
		value := strings.Replace(self.Text, "@Money", "", 1)
		count, _ := strconv.Atoi(value)
		team.AwardObject(awardTypeCoin, count, 0, 0)
		team.AwardObject(awardTypeTicket, count, 0, 0)
		return true
	case strings.Index(self.Text, "@Item") != -1:
		value := strings.Replace(self.Text, "@Item", "", 1)
		itemType, _ := strconv.Atoi(value)
		itemMgr.AwardItem(itemType, 1)
		team.AwardObject(itemType, 99, 0, 0)
		return true
	case strings.Index(self.Text, "@Exppool") != -1: ///经验池
		value := strings.Replace(self.Text, "@Exppool", "", 1)
		count, _ := strconv.Atoi(value)
		team.AwardExpPool(count)
		sync.SyncObject("GMCommand", team)
		return true
	case strings.Index(self.Text, "@clsfunmask") != -1: ///关闭所有功能
		team.FunctionMask = 0
		sync.SyncObject("GMCommand", team)
	case strings.Index(self.Text, "@mail") != -1:
		mailMgr := client.GetTeam().GetMailMgr()
		//		targetName := "什么叫nb"
		//		targetID := 1070106
		for i := 0; i < 5; i++ {
			//			mailMgr.SendMatchReport(2, 1, targetName, targetID, 0, 1, 101, 3, 3, 360, 359)
			mailMgr.SendSysAwardMail(1, 1, IntList{200001}, IntList{0}, IntList{5}, "", "")
		}
	case strings.Index(self.Text, "@star") != -1:
		value := strings.Replace(self.Text, "@star", "", 1)
		starType, _ := strconv.Atoi(value)
		team.AwardObject(0, 0, 1, starType)
	case strings.Index(self.Text, "@time") != -1:
		value := strings.Replace(self.Text, "@time", "", -1)
		cmd := exec.Command("date", "-s", value)
		cmd.Run()
	case strings.Index(self.Text, "@form") != -1:
		value := strings.Replace(self.Text, "@form", "", -1)
		formType, _ := strconv.Atoi(value)
		formationID := team.GetFormationMgr().AwardFormation(formType)
		sync.syncAddFormation(IntList{formationID}) ///同步新加的阵型到客户端
	case strings.Index(self.Text, "@monthcard") != -1:
		vipBuyMonthCardMsg := new(VipBuyMonthCardMsg)
		vipBuyMonthCardMsg.processAction(client)
	}
	//gmMgr.ExecCommand(client, self.Text)
	return false

}

func (self *ChatMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false { ///通用检测
		return false
	}

	if self.GMCommand(client) == true {
		return true
	}

	switch self.Type {
	case chatTypeWhispe: ///私聊
		self.processWhispe(client)
	case chatTypeWorld: ///世界聊天
		self.processWorld(client)
	}
	//	client.SendMsg(self) ///所有聊天消息则服务器应答回显示,用来测试网络连通性
	//client.SendSysNotice("hello world!")
	return true
}

func (self *ChatMsg) sendWhispeMsg(userMgr *UserMgr) bool {
	receiverClient := userMgr.GetClientByTeamName(self.Receiver) ///得到接收者
	senderClient := userMgr.GetClientByTeamName(self.Sender)     ///得到发送者
	if receiverClient != nil {
		receiverClient.SendMsg(self) ///可以找到接收者直接发送
	} else if senderClient != nil {
		senderClient.SendErrorMsg(failChatWhispe, failNotFound) ///找不到接收者
	}
	return true
}

func (self *ChatMsg) broacastWorldMsg(userMgr *UserMgr) bool {
	userMgr.broadcastMsg(self)
	return true
}

func (self *ChatMsg) KickSelf(UserMgr *UserMgr) bool {
	selfClient := UserMgr.GetClientByTeamName(self.Receiver)
	GetServer().Kickout("Repeat Login", selfClient.GetElement())
	return true
}

func (self *ChatMsg) broacastMsg(userMgr *UserMgr) bool {
	switch self.Type {
	case chatTypeSystem: ///系统公告
		self.broacastWorldMsg(userMgr)
	case chatTypeWhispe: ///私聊
		self.sendWhispeMsg(userMgr)
	case chatTypeWorld: ///世界聊天
		self.broacastWorldMsg(userMgr)
	case chatGMKick:
		self.KickSelf(userMgr)
	}
	return true
}
