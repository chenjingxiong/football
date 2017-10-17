package football

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

type AccountInfo struct {
	ID               int    `json:"id"`               ///帐号id
	Name             string `json:"name"`             ///登录名/平台用户编号
	Password         string `json:"password"`         ///登录密码
	Channel          int    `json:"channel"`          ///平台渠道编号,对应platform表
	IP               string `json:"ip"`               ///客户端ip地址
	LastServer       int    `json:"lastserver"`       ///最近登录游戏区编号
	MakeTime         int    `json:"maketime"`         ///帐号创建时间 年月日时分 1408010123
	SDKName          string `json:"sdkname"`          ///平台名
	SDKUserID        string `json:"sdkuserid"`        ///平台用户编号
	AuthCode         string `json:"authcode"`         ///认证码,用户名与密码映射结果
	AccessKey        string `json:"accesskey"`        ///访问密钥
	RefreshKey       string `json:"refreshkey"`       ///更新密钥
	AccessKeyExpires int    `json:"accesskeyexpires"` ///访问密钥期限,utc
}

type IClient interface { ///客户端受限接口
	GetTeam() *Team                                                   ///得到客户端拥有的team组件受限接口
	checkMsgSeq(msgSeq int) bool                                      ///验证包序列号
	checkClientVersion(versionClient *string) bool                    ///检测客户端版本正确性
	createClientAccount(loginName *string, loginPwd *string) int      ///创建客户端帐号,仅用于测试
	checkClientLogin(loginName *string, loginPwd *string) (int, bool) ///检测客户端登录信息正确性,1 帐号id 2 密码是否正确
	SendErrorMsg(errorType string, errorDesc string)                  ///发送错误信息到客户端
	LoadTeam(accountID int) bool                                      ///通过账号id查询所拥有的teamid
	SendTeam()                                                        ///发送球队信息和球员信息给客户端
	SendLoginResultMsg(loginResult string)                            ///发送登录结果到客户端
	GetElement() *Client                                              ///特权接口
	HasInitTeam() bool                                                ///判断客户端是否已创建了team组件
	SendMsg(msg IMsgHead)                                             ///发送消息
	GetSyncMgr() *SyncMgr                                             ///得到客户端拥有的同步数据组件
	BroadcastMsg(msg IMsgHead)                                        ///发送消息
	SendSysNotice(textSysNotice string)                               ///发送系统公告
	LoginRecord(action int)                                           ///登录/登出记录
	SetMoneyRecord(state int, action int, value int, money int)       ///金钱变动记录
	LevelUpRecord(level int)                                          ///玩家升级记录
	CreateTeamRecord()                                                ///创建球队
	RechargeRecord(actionType int, rechargeMoney int)                 ///添加获得钻石记录
}

type Client struct { ///客户端核心数据结构
	id             int              ///玩家编号,从1开始
	ws             *websocket.Conn  ///客户端网络底层套接字
	sendMsgChannel chan interface{} ///用于发送协程发送消息
	accountID      int              ///帐号id,用于创建新球队时使用
	seqID          int              ///消息顺序,防止客户端发送重复消息
	team           *Team            ///球队对象
	teamName       string           ///用于广播需求,此变理可能存在线程不安全可能
	teamID         int              ///用于广播需求,此变理可能存在线程不安全可能
	accountName    string           ///保存玩家帐号名用于排错
	syncMgr        *SyncMgr         ///同步客户端服务器最新状态与数据组件
	//gmMgr          *GameMaster      ///游戏GM功能
	accessToken  string
	refreshToken string
	expiresIn    int
	userID       string
	authCode     string ///认证码
	sdkName      string ///sdk平台名字
}

const (
	PlayerRegistration = 1 ///玩家注册
	PlayerCreateTeam   = 2 ///玩家创建
)

const (
	PlayerLogin    = 1 ///玩家登录/登出
	PlayerRecharge = 2 ///玩家充值
	PlayerLevelUp  = 3 ///玩家等级提升
)

const (
	PlayerLogout = 2 ///玩家登出
)

const (
	PlayerGetMoney  = 1 ///获取钻石
	PlayerCostMoney = 2 ///花费钻石
)

const (
	//获取钻石
	Get_Recharge       = 1 //1: 充值
	Get_GMCommand      = 2 //2: GM命令
	Get_VIPShopMoney   = 3 //3: vip充值商店
	Get_MailItem       = 4 //4: 邮件道具获得
	Get_PassLevelAward = 5 //5:冠军之路过关奖励
	Get_ItemUse        = 6 //6:使用道具获得钻石
	Get_MonthCard      = 7 //6:月卡福利获得钻石

	//失去钻石
	Pay_BuyTactic         = 200 // 200: 购买战术
	Pay_SpecialAward      = 201 // 201: 关卡特殊奖励
	Pay_StarSack          = 202 // 202: 解雇球员
	Pay_BuyActionPoint    = 203 // 203: 购买行动点
	Pay_VipShop           = 204 // 204: 商城
	Pay_RefreshTrainMatch = 205 // 205: 训练赛刷新
	Pay_GodFinger         = 206 // 206: 使用金手指
)

type Record_Create struct { ///记录_注册创建
	ID        int ///记录ID
	AccountID int ///帐号ID
	TeamID    int ///队伍ID
	Type      int ///类型: 1.注册 2.创建
}

type Record_Action struct { ///记录_玩家行为
	ID        int ///记录ID
	AccountID int ///帐号ID
	TeamID    int ///队伍ID
	Type      int ///类型: 1.登录 2.充值 3.升级
	Param     int ///p1上线/下线(1 2) p2充值钻石数 p3级别
}

type Record_Pay struct { ///记录_消耗
	ID        int ///记录ID
	AccountID int ///帐号ID
	TeamID    int ///队伍ID
	Type      int ///类型: 1.登录 2.充值 3.升级
	Action    int ///获取/消耗玩家所做行为
	Param     int ///p1上线/下线(1 2) p2充值钻石数 p3升级花费时间
}

func (self *Client) Run() { ///发送消息到频道中，等着发送协程写到网络层
	go self.SendMsgGo() ///开启读消息协程
	self.ReceiveMsgGo() ///运行读消息协程
}

func (self *Client) GetElement() *Client {
	return self
}

func (self *Client) HasInitTeam() bool { ///判断客户端是否已创建了team组件
	return self.team != nil
}

func (self *Client) GetSyncMgr() *SyncMgr { ///得到客户端拥有的同步数据组件
	return self.syncMgr
}
func (self *Client) GetTeam() *Team { ///得到客户端拥有的team组件受限接口
	return self.team
}

//func (self *Client) GetGMMgr() *GameMaster { ///得到游戏管理员功能
//	return self.gmMgr
//}
//
func (self *Client) LoginRecord(action int) { ///记录玩家登入/登出
	///记录客户端登出信息
	clientLogoutRecord := fmt.Sprintf("insert %s set accountid = %d ,teamid = %d ,type = %d, param = %d,sdkname='%s'", tableRecordAction,
		self.accountID, self.teamID, PlayerLogin, action, self.sdkName)

	lastInsertRecordid, _ := GetServer().GetRecordDB().Exec(clientLogoutRecord)
	if lastInsertRecordid <= 0 {
		GetServer().GetLoger().Warn("insert %s login/logout error, accountid = %d", tableRecordAction, self.accountID)
	}

	///记录玩家在线表
	if action == PlayerLogin {
		LoginRecord := fmt.Sprintf("insert %s set teamid = %d,sdkname='%s'", tableRecordOnline, self.teamID, self.sdkName)

		lastInsertRecordid, _ := GetServer().GetRecordDB().Exec(LoginRecord)
		if lastInsertRecordid <= 0 {
			GetServer().Kickout("Repeat Login", self)
			GetServer().kickoutRepeatLogin(self.teamName)
		}
	} else if action == PlayerLogout {
		LogoutRecord := fmt.Sprintf("delete from %s where teamid = %d", tableRecordOnline, self.teamID)

		_, rowsAffected := GetServer().GetRecordDB().Exec(LogoutRecord)
		if rowsAffected <= 0 {
			GetServer().GetLoger().Warn("insert %s login/logout error, teamid = %d", tableRecordOnline, self.teamID)
		}
	}
}

func (self *Client) CreateTeamRecord() { /// 记录创建球队信息
	createTeamRecord := fmt.Sprintf("insert %s set accountid = %d ,teamid = %d ,type = %d,sdkname='%s'",
		tableRecordCreate, self.accountID, self.teamID, PlayerCreateTeam, self.sdkName)

	///记录创建队伍信息
	lastInsertRecordid, _ := GetServer().GetRecordDB().Exec(createTeamRecord)
	if lastInsertRecordid <= 0 {
		GetServer().GetLoger().Warn("insert %s create team error, accountid = %d", tableRecordCreate, self.accountID)
	}

	///玩家登录表
	self.LoginRecord(PlayerLogin)
}

func (self *Client) RechargeRecord(actionType int, rechargeMoney int) {
	money := self.team.Ticket
	playerRechargeRecord := fmt.Sprintf("insert %s set accountid = %d ,teamid = %d ,type = %d, param = %d,sdkname='%s'",
		tableRecordPay, self.accountID, self.teamID, PlayerRecharge, rechargeMoney, self.sdkName)

	lastInsertRecordid, _ := GetServer().GetRecordDB().Exec(playerRechargeRecord)
	if lastInsertRecordid <= 0 {
		GetServer().GetLoger().Warn("insert %s recharge money error, accountid = %d", tableRecordPay, self.accountID)
	}

	//充值记录后,记录玩家获取钻石数
	self.SetMoneyRecord(PlayerGetMoney, actionType, rechargeMoney, money)
}

func (self *Client) SetMoneyRecord(state int, action int, value int, money int) {
	///记录客户端获取/花费钻石信息
	clientPayRecord := fmt.Sprintf("insert %s set accountid = %d ,teamid = %d ,type = %d,action = %d, param = %d, money = %d,sdkname='%s'",
		tableRecordPay, self.accountID, self.teamID, state, action, value, money, self.sdkName)

	lastInsertRecordid, _ := GetServer().GetRecordDB().Exec(clientPayRecord)
	if lastInsertRecordid <= 0 {
		GetServer().GetLoger().Warn("insert %s error, accountid = %d", tableRecordPay, self.accountID)
	}
}

func (self *Client) LevelUpRecord(level int) { ///玩家升级记录
	playerLevelUpRecord := fmt.Sprintf("insert %s set accountid = %d ,teamid = %d ,type = %d, param = %d,sdkname='%s'",
		tableRecordAction, self.accountID, self.teamID, PlayerLevelUp, level, self.sdkName)

	lastInsertRecordid, _ := GetServer().GetRecordDB().Exec(playerLevelUpRecord)
	if lastInsertRecordid <= 0 {
		GetServer().GetLoger().Warn("insert %s level up error, accountid = %d", tableRecordPay, self.accountID)
	}
}

func (self *Client) OnLogout() { ///客户端登出时逻辑处理
	if self.team != nil {
		self.team.GetStarSpy().OnLogout(self.team)
		self.team.OnLogout() ///通知球队调用退出逻辑
		self.LoginRecord(PlayerLogout)
	}
}

func (self *Client) ReceiveMsgGo() { ///读客户端协程
	GetServer().waitGroup.Add(1)     ///添加一个客户端
	rand.Seed(time.Now().UnixNano()) ///每个客户端初始化一次随机种子
	msg := ""
	for {
		self.ws.SetReadDeadline(time.Now().Add(1 * time.Second)) ///设置1秒超时,间隔1秒进行update
		err := websocket.Message.Receive(self.ws, &msg)
		if nil == err {
			fmt.Println("get a client msg:", msg)
			self.ProcessMsg(&msg) ///无错时处理消息并进行下一次侦听
			continue
		}
		neterr, ok := err.(net.Error)
		if true == ok && neterr.Timeout() == true {
			self.Update() ///如果是IO超时错误则进行一次update并进行下一次侦听
			continue
		}
		///其它错误出日志并断开连接
		GetServer().GetLoger().Warn("client %d get a error when Receive:%v", self.id, err)
		break
	}
	GetServer().GetLoger().Info("client %d exit ReceiveMsgGo", self.id)
	self.ws.Close()
	self.OnLogout()              ///客户端登出时逻辑处理
	GetServer().waitGroup.Done() ///删除一个客户端
}

func (self *Client) SendMsgGo() { ///读客户端协程
	for msg := range self.sendMsgChannel {
		err := websocket.JSON.Send(self.ws, msg)
		msgText, _ := json.Marshal(msg) ///解出消息
		//fmt.Println("msg size:", len(msgText))
		fmt.Println("send client msg:", "msg size:", len(msgText), string(msgText))
		if err != nil {
			GetServer().GetLoger().Warn("client %d get a error when Send:%v", self.id, err)
			break
		}
	}
	GetServer().GetLoger().Info("client %d exit SendMsgGo", self.id)
	self.ws.Close()
}

func (self *Client) checkClientLogin(loginName *string, loginPwd *string) (int, bool) { ///检测客户端登录信息正确性,1 帐号id 2 密码是否正确
	ok := false
	accountID := 0
	accountPassword := ""
	rows := GetServer().GetLoginDB().Query("select id,password from dy_account where name='%s' limit 1", *loginName)
	if nil == rows {
		return 0, false
	}
	for rows.Next() {
		rows.Scan(&accountID, &accountPassword)
		ok = (accountPassword == *loginPwd)
	}
	rows.Close()
	self.accountName = *loginName
	nameList := strings.Split(self.accountName, "_")
	if len(nameList) >= 1 {
		self.sdkName = nameList[0]
	}
	return accountID, ok
}

func (self *Client) SendLoginResultMsg(loginResult string) { ///发送登录结果到客户端
	msg := NewLoginResultMsg(&loginResult)
	msg.AccountID = self.accountID
	self.SendMsg(msg)
}

func (self *Client) SendSysNotice(textSysNotice string) { ///发送系统公告
	GetServer().SendSysNotice(textSysNotice)
}

func (self *Client) SendErrorMsg(errorType string, errorDesc string) { ///发送错误信息到客户端
	msg := NewActionErrorMsg(&errorType, &errorDesc)
	self.SendMsg(msg)
}

func (self *Client) SendCreatTeam() { ///发送客户端创建球队命令
	msg := NewTeamInfoMsg()
	msg.TeamInfo = self.team.TeamInfo
	self.SendMsg(msg)
}

func (self *Client) SendTeam() { ///发送球队信息和球员信息给客户端
	self.team.OnEnterMap()
}

func (self *Client) LoadTeam(accountID int) bool { ///通过账号id查询所拥有的teamid
	if self.team != nil {
		GetServer().GetLoger().Warn("Client LoadTeam duplicate error! accountID:%d", accountID)
		return false
	}
	syncMgr := NewSyncMgr(self) ///初始化信息同步组件
	///加载球队信息
	team := new(Team)
	if team.Create(accountID, 0, syncMgr) == false {
		self.accountID = accountID
		GetServer().GetLoger().Debug("Client LoadTeam not exsit.accountID:%d", accountID)
		return false
	}
	self.team = team
	self.accountID = accountID
	self.teamName = team.GetInfo().Name ///放入用于广播的名字
	self.teamID = team.GetInfo().ID     ///放入用于广播的名字
	self.syncMgr = syncMgr
	self.syncMgr.Init(self)
	return true
}

func (self *Client) HasTeam(accountID int) bool { ///判断此账号是否已经创建球队了
	teamInfo := new(TeamInfo)
	query := fmt.Sprintf("select * from dy_team where accountid=%d limit 1", accountID)
	GetServer().GetDynamicDB().fetchOneRow(query, teamInfo)
	return teamInfo.AccountID == accountID
}

func (self *Client) queryTeamID(accountID int) int { ///通过账号id查询所拥有的teamid
	teamID := 0
	query := fmt.Sprintf("select id from dy_team where accountid=%d", accountID)
	rows := GetServer().GetDynamicDB().Query(query)
	if nil == rows {
		return 0
	}
	for rows.Next() {
		rows.Scan(&teamID)
	}
	rows.Close()
	return teamID
}

func (self *Client) createClientAccount(loginName *string, loginPwd *string) int { ///创建客户端帐号,仅用于测试
	accountID, _ := self.checkClientLogin(loginName, loginPwd) ///首先检测此用户是否存在
	if accountID != 0 {
		return accountID ///有此用户
	}
	createAccountQuery := fmt.Sprintf("insert dy_account set name='%s',password='%s'", *loginName, *loginPwd)
	lastInsertID, _ := GetServer().GetLoginDB().Exec(createAccountQuery) ///创建新帐号

	///记录创建帐号信息
	createAccountRecord := fmt.Sprintf("insert %s set accountid = %d ,type = %d,sdkname='%s'", tableRecordCreate, lastInsertID, PlayerRegistration, self.sdkName)
	lastRecordID, _ := GetServer().GetRecordDB().Exec(createAccountRecord) ///创建记录: 玩家注册
	if lastRecordID <= 0 {
		GetServer().GetLoger().Warn("insert %s registration error, accountid = %d", tableRecordCreate, lastInsertID)
	}
	self.accountID = lastInsertID
	return int(lastInsertID)
}

func (self *Client) checkClientVersion(versionClient *string) bool { ///检测客户端版本正确性
	serverVersion := GetServer().GetVersion()
	result := *versionClient == serverVersion
	if false == result {
		GetServer().GetLoger().Warn("client version error!%v!=1.2.3.4", *versionClient)
	}
	return result
}

func (self *Client) checkMsgSeq(msgSeq int) bool {
	if msgSeq != self.seqID {
		GetServer().GetLoger().Warn("dispatchMsg get seqID error msgSeq %d != expect:%d", msgSeq, self.seqID)
		return false ///非法
	}
	return true ///合法
}

func (self *Client) BroadcastMsg(msg IMsgHead) { ///发送消息
	GetServer().BroadcastMsg(msg) ///转给服务器处理
}

func (self *Client) SendMsg(msg IMsgHead) { ///发送消息
	//msg.setSeq(self.seqID)
	msgType, msgAction := msg.GetTypeAndAction()
	msg.FillMsgHead(self.seqID, &msgType, &msgAction)
	self.sendMsgChannel <- msg
}

func (self *Client) Update() { ///客户端自身更新状态
	///更新球队状态
	if self.team != nil {
		now := int(time.Now().Unix())
		self.team.Update(now, self)
	}
}

func (self *Client) ProcessMsg(msg *string) { ///处理协议客户端协程
	self.seqID++ ///接受新的消息时自动更新序列号
	GetServer().GetMsgDispatch().DispatchMsg(self, msg)
	//self.Update() ///调用自更新状态
}

func (self *Client) UpdateSeqID() {
	self.seqID++
}
