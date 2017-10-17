package football

import (
	"code.google.com/p/go.net/websocket"
	//"log"
	"runtime"
)

//const addClientChannelSize = 0        ///添加客户端频道缓冲区长度
//const removeClientChannelSize = 0     ///删除客户端频道缓冲区长度
const msgBroadcastChannelSize = 1000  ///服务器广播消息缓冲区长度
const clientMsgSendChannelSize = 1000 ///服务器广播消息缓冲区长度

type ClientList map[int]*Client ///定义客户列表类型

type UserMgr struct { ///用户管理器
	clientList ClientList ///用户列表
	//ownServer  *Server    ///服务器组件
}

func NewUserMgr() *UserMgr {
	userMgr := new(UserMgr)
	userMgr.clientList = ClientList{}
	//userMgr.ownServer = server
	return userMgr
}

func (self *UserMgr) getNewClient(clientID int, ownServer *Server, ws *websocket.Conn) *Client { ///创建一个客户端实例
	newClient := new(Client)
	newClient.id = clientID ///使用当前纳秒时间做客户端唯一标识号
	newClient.ws = ws
	newClient.sendMsgChannel = make(chan interface{}, clientMsgSendChannelSize) ///用于发送协程发送消息
	//newClient.ownServer = ownServer                                             ///保存上级指针
	newClient.seqID = 0 ///消息序号从0开始
	return newClient
}

func (self *UserMgr) addClient(newClient *Client) {
	self.clientList[newClient.id] = newClient
	GetServer().GetLoger().Debug("addClient client id %d done,server remain %d client NumGoroutine:%d", newClient.id, len(self.clientList), runtime.NumGoroutine())
}

//func (self *Server) notifyLeave(goneClientID int) { ///处理协议客户端协程
//	notifyText := fmt.Sprintf("client %d has left the chat room.", goneClientID)
//	msg := ChatMsg{"system", notifyText}
//	self.broadcastMsg(&msg)
//}

func (self *UserMgr) removeClient(client *Client) {
	_, ok := self.clientList[client.id]
	if false == ok {
		return ///避免重复删除
	}
	client.ws.Close()                  ///关闭客户端套接字,可能有错误
	close(client.sendMsgChannel)       ///删除客户间关闭它的发送频道
	delete(self.clientList, client.id) ///从队列中删除此客户端
	GetServer().GetLoger().Debug("removeClient client id %d done,server remain %d client,NumGoroutine:%d", client.id, len(self.clientList), runtime.NumGoroutine())
}

func (self *UserMgr) GetClientByTeamID(teamID int) IClient {
	for _, client := range self.clientList {
		if client.teamID == teamID {
			return client
		}
	}
	return nil
}

func (self *UserMgr) GetClientByTeamName(teamName string) IClient {
	for _, client := range self.clientList {
		if client.teamName == teamName {
			return client
		}
	}
	return nil
}

func (self *UserMgr) broadcastMsg(msg interface{}) {
	for _, client := range self.clientList {
		select {
		case client.sendMsgChannel <- msg:
		default:
			GetServer().GetLoger().Debug("get a disconnect when broadcastMsg to client id %d msg:%v", client.id, msg)
			client.ws.Close()
			self.removeClient(client)
		}
	}
}
