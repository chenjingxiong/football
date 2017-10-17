package football

import (
	"encoding/json"
	"runtime/debug"
	"time"
)

const (
	msgResultOK   = "ok"
	msgResultFail = "fail"
)

type MsgHead struct {
	SeqID      int    `json:"seq"`    ///顺序号
	MsgType    string `json:"type"`   ///消息类型
	Action     string `json:"action"` ///操作类型
	CreateTime int    `json:"time"`   ///消息创建时间
}

func (self *MsgHead) FillMsgHead(seqID int, msgType *string, msgAction *string) { ///填写消息头
	if self.SeqID > 0 { ///忽略掉重复的调用
		return
	}
	self.SeqID = seqID
	self.CreateTime = int(time.Now().Unix())
	self.MsgType = *msgType
	self.Action = *msgAction
}

func (self *MsgHead) processAction(client IClient) bool { ///需要实现处理消息
	return true
}

///判断是否需要玩家登录完成并创建了球队后才处理此消息
func (self *MsgHead) IsNeedTeamHandle() bool {
	return true
}

func (self *MsgHead) setSeq(seqID int) { ///设置消息序列号
	self.SeqID = seqID
	self.CreateTime = int(time.Now().Unix())
}

func (self *MsgHead) GetItemKey() string { ///得到注册键值
	//msgType,msgAction:=self.
	itemKey := self.MsgType + "-" + self.Action
	return itemKey
}

func (self *MsgHead) broacastMsg(userMgr *UserMgr) bool {
	return true
}

type IMsgHead interface {
	FillMsgHead(seqID int, msgType *string, msgAction *string) ///填写消息头
	GetTypeAndAction() (string, string)                        ///得到消息大分类和小分类
	processAction(client IClient) bool                         ///需要实现处理消息
	GetItemKey() string                                        ///得到注册键值
	broacastMsg(userMgr *UserMgr) bool                         ///广播消息
	IsNeedTeamHandle() bool                                    ///判断此消息是否需要创建球队才能处理
}

type HeadMsg struct {
	MsgHead `json:"head"`
}

type IMsgHandler interface {
	getName() string                          ///能处理的消息类型
	initHandler()                             ///初始化消息处理器函数
	processMsg(IClient, string, *string) bool ///处理消息函数,以后需要优化
	Init()                                    ///初始化内部状态
}

type IActionHandler interface {
	GetTypeAndAction() (string, string)
	processAction(client IClient) bool
}

type ActionHandlerList map[string]IActionHandler

type MsgHandler struct {
	actionHandlerList ActionHandlerList
}

func (self *MsgHandler) Init() { ///初始化内部状态
	self.actionHandlerList = ActionHandlerList{}
}

func (self *MsgHandler) addActionToList(actionHandler IActionHandler) { ///初始化消息子类型处理器
	_, msgAction := actionHandler.GetTypeAndAction()
	self.actionHandlerList[msgAction] = actionHandler
}

func (self *MsgHandler) processMsg(client IClient, msgAction string, rawMsg *string) bool {
	defer func() {
		x := recover()
		if x != nil {
			GetServer().GetLoger().Error("%v\r\n%s", x, debug.Stack())
		}
	}()
	actionHandler, ok := self.actionHandlerList[msgAction]
	if ok != true {
		GetServer().GetLoger().Warn("MsgHandler processMsg get a unknow action!%v", *rawMsg)
		return false
	}
	newObj := CloneType(actionHandler)                       ///从注册信息中生成新的对象
	newActionHandler := newObj.(IActionHandler)              ///对新对象进行转型操作
	err := json.Unmarshal([]byte(*rawMsg), newActionHandler) ///解出消息
	if err != nil {
		GetServer().GetLoger().Warn("ChatHandler processMsg get a error!%v", err.Error())
		return false
	}
	result := newActionHandler.processAction(client) ///处理消息
	return result
}

type MsgRegistry map[string]IMsgHead ///消息注册表
type MsgHandlerList map[string]IMsgHandler
type MsgDispatch struct { ///消息分拣器
	msgHandlerList MsgHandlerList ///消息处理器列表
	msgRegistry    MsgRegistry    ///消息注册表
}

type IMsgDispatch interface { ///消息分拣器
	DispatchMsg(client IClient, msg *string) bool
}

func (self *MsgDispatch) addMsgRegistry(msg IMsgHead) { ///初始化新的消息子类型处理器
	if nil == self.msgRegistry {
		self.msgRegistry = make(MsgRegistry)
	}
	msgType, msgAction := msg.GetTypeAndAction()
	itemKey := msgType + "-" + msgAction
	_, ok := self.msgRegistry[itemKey]
	if true == ok {
		GetServer().GetLoger().Fatal("MsgDispatch addMsgRegistry duplicate msg! %s", itemKey)
	}
	self.msgRegistry[itemKey] = msg
}

//func (self *MsgDispatch) addMsgHandleToList(msgHandler IMsgHandler) { ///初始化消息子类型处理器
//	_, ok := self.msgHandlerList[msgHandler.getName()]
//	if true == ok {
//		GetServer().GetLoger().Fatal("MsgDispatch addMsgHandleToList duplicate handle! %s", msgHandler.getName())
//	}
//	msgHandler.Init()
//	self.msgHandlerList[msgHandler.getName()] = msgHandler
//}

func (self *MsgDispatch) initMsgDispatch() {
	self.msgHandlerList = make(MsgHandlerList)
	self.initMsgHandlerList() ///初始化消息处理器列表
	for name, handler := range self.msgHandlerList {
		if handler.getName() == name {
			handler.initHandler()
		}
	}
}

///将客户端传来的串消息转换成JSON消息
func (self *MsgDispatch) GetMsgHead(rawMsg *string) *HeadMsg {
	headMsg := new(HeadMsg)
	err := json.Unmarshal([]byte(*rawMsg), headMsg)
	if err != nil {
		GetServer().GetLoger().Warn("GetMsgHead get a unknow msghead!%v", err.Error())
		return nil
	}
	return headMsg
}

func (self *MsgDispatch) BroadcastMsg(userMgr *UserMgr, msg interface{}) bool {
	//defer func() {
	//	x := recover()
	//	if x != nil {
	//		GetServer().GetLoger().Error("%v\r\n%s", x, debug.Stack())
	//	}
	//}()
	needBroadcastMsg := msg.(IMsgHead)
	result := needBroadcastMsg.broacastMsg(userMgr) ///广播此消息
	return result
}

func (self *MsgDispatch) DispatchMsg(client IClient, rawMsg *string) bool {
	defer func() {
		x := recover()
		if x != nil {
			rawClient := client.GetElement()
			GetServer().GetLoger().Error("{%s:%s}%v\r\n%s", rawClient.accountName, rawClient.teamName, x, myStack(6))
		}
	}()
	headMsg := self.GetMsgHead(rawMsg) ///得到消息头
	itemKey := headMsg.GetItemKey()
	msgRegistry := self.msgRegistry[itemKey] ///取得消息实列
	if nil == msgRegistry {
		GetServer().GetLoger().Warn("dispatchMsg get a unknow msghead!%v", headMsg)
		return false
	}
	if msgRegistry.IsNeedTeamHandle() == true && client.HasInitTeam() == false {
		GetServer().GetLoger().Warn("dispatchMsg get a illegal msg!need team.%v", headMsg)
		return false
	}
	newObj := CloneType(msgRegistry)               ///从注册表中生成新的对象
	newMsg := newObj.(IMsgHead)                    ///对新对象进行转型操作
	err := json.Unmarshal([]byte(*rawMsg), newMsg) ///解出消息
	if err != nil {
		GetServer().GetLoger().Warn("MsgDispatch DispatchMsg Unmarshal get a error!%v", err.Error())
		return false
	}
	result := newMsg.processAction(client) ///处理此消息
	return result
}

//func (self *MsgDispatch) DispatchMsg(client IClient, msg *string) bool {
//	//GetServer().GetLoger().Debug("%v", *msg)
//	//var headMsg *HeadMsg = self.GetMsgHead(msg) ///得到消息头
//	headMsg := self.GetMsgHead(msg) ///得到消息头
//	if nil == headMsg {
//		GetServer().GetLoger().Warn("dispatchMsg get a unknow msghead!%v", msg)
//		return false
//	}
//	if false == client.checkMsgSeq(headMsg.SeqID) {
//		GetServer().GetLoger().Warn("dispatchMsg client.checkMsgSeq fail! %v", *msg)
//		return false ///非法的包序号,丢弃此包
//	}
//	result, isSkipMsg := self.DispatchMsg2(client, headMsg, msg) ///新的消息处理流程
//	if false == isSkipMsg {                                      ///此消息新流程已处理需要跳过的
//		return result
//	}
//	msgHandler, ok := self.msgHandlerList[headMsg.MsgType] ///取得消息处理器
//	if ok != true {
//		GetServer().GetLoger().Warn("dispatchMsg get a unknow msg!%v", *msg)
//		return false
//	}
//	result = msgHandler.processMsg(client, headMsg.Action, msg)
//	return result
//}
