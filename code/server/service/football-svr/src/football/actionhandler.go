package football

import (
//"fmt"
//"time"
)

const ( ///ErrorType
	failLogin            = "loginfail"            ///登录失败
	failCreateTeam       = "createteamfail"       ///创建队伍失败
	failStarSpyDiscover  = "starspydiscoverfail"  ///球探发掘球员失败
	failFormationUplevel = "formationuplevelfail" ///球队阵形升级失败
	failChatWhispe       = "chatwhispefail"       ///私聊失败
)

const ( ///ErrorDesc
	failAccountNotExsit           = "AccountNotExsit"           ///帐户不存在
	failPasswordWrong             = "PasswordWrong"             ///密码错误
	failSameName                  = "SameName"                  ///同名冲突
	failInvalidName               = "InvalidName"               ///名字非法
	failInvalidAccountID          = "InvalidAccountID"          ///名字非法
	failInvalidMsg                = "InvalidMsg"                ///消息非法,非法的发送时机
	failInvalidParam              = "InvalidParam"              ///消息非法,非法的发送时机
	failReachLimit                = "ReachLimit"                ///超过限制
	failInsufficientTicket        = "InsufficientTicket"        ///球票不足
	failInsufficientDiscoverCount = "InsufficientDiscoverCount" ///发掘次数不足
	failInreachmaxlevel           = "ReachMaxLevel"             ///超过等级上限
	failNotFound                  = "NotFound"                  ///找不到对象
	failActivitCode               = "ActivitCodeError"          ///激活码错误或已被使用过
	failErrorVersion              = "ErrorVersion"              ///版本号错误
)

type ActionErrorMsg struct { ///错误提示消息
	MsgHead   `json:"head"`
	ErrorType string `json:"errortype"` ///错误类型
	ErrorDesc string `json:"errordesc"` ///错误描述
}

func (self *ActionErrorMsg) GetTypeAndAction() (string, string) {
	return "action", "error"
}

func NewActionErrorMsg(errorType *string, errorDesc *string) *ActionErrorMsg {
	msg := new(ActionErrorMsg)
	//msg.MsgType = "action"
	//msg.Action = "error"
	//msg.CreateTime = int(time.Now().Unix())
	msg.ErrorType = *errorType
	msg.ErrorDesc = *errorDesc
	return msg
}

type ActionHandler struct {
	MsgHandler
}

func (self *ActionErrorMsg) New() interface{} {
	return new(ActionErrorMsg)
}

func (self *ActionErrorMsg) getName() string {
	return "error"
}

func (self *ActionErrorMsg) processAction(client IClient) bool {
	//actionErrorMsg := msg.(*ActionErrorMsg)
	GetServer().GetLoger().Info("ActionMsg processAction msg:%v", self)
	return false
}

func (self *ActionHandler) getName() string { ///返回可处理的消息类型
	return "action"
}

func (self *ActionHandler) initHandler() { ///初始化消息处理器
	//self.addActionToList(new(ActionErrorMsg))
}
