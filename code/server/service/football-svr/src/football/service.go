package football

import (
	"code.google.com/p/go.net/websocket"
	//	"fmt"
	//	"io"
	"log"
	//	"net"
	"net/http"
	"runtime" ///服务端版本号
	"strconv"

	"os"
	"os/signal"
	"time"

	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net"
	"runtime/debug"
	"sync"
	"syscall"
)

type ServerConfig struct { ///服务器配置文件,用于记录WS侦听地址,数据库权限等
	WSHost            string `json:"wshost"`            ///WS侦听地址与端口
	DynamicDNS        string `json:"dynamicdns"`        ///动态库DNS
	StaticDNS         string `json:"staticdns"`         ///静态库DNS
	LoginDNS          string `json:"logindns"`          ///登录库DNS
	RecordDNS         string `json:"recorddns"`         ///记录库DNS
	MinLogLevel       int    `json:"minloglevel"`       ///最小日志记录等级 1 调试 2 信息 3 警告  4 错误
	MaxIdleDBConns    int    `json:"dbmaxidleconns"`    ///数据库连接池最大空闲连接数
	RunKey            string `json:"runkey"`            ///运行键值
	OpenServerTime    string `json:"openservertime"`    ///开服时间
	ServerID          int    `json:"serverid"`          ///服务器编号
	IsNeedActivitCode int    `json:"isneedactivitcode"` ///当前服务器是否需要激活码
	OpenGMTalk        int    `json:"opengmtalk"`        ///是否允许接受gm指令
}

func (self *ServerConfig) defaultConfig() {
	self.WSHost = "0.0.0.0:8080"
	self.DynamicDNS = "root:1234@tcp(192.168.20.229:3306)/football_dynamic?charset=utf8"
	self.StaticDNS = "root:1234@tcp(192.168.20.229:3306)/football_static?charset=utf8"
	self.LoginDNS = "root:1234@tcp(192.168.20.229:3306)/football_login?charset=utf8"
	self.RecordDNS = "root:1234@tcp(192.168.20.229:3306)/football_record?charset=utf8"
	self.MinLogLevel = logInfo ///默认日志信息等级
	self.MaxIdleDBConns = 500  ///默认500,根据单服最大人数考虑
	self.RunKey = "demo"       ///演示模式
	self.ServerID = 10001      ///默认服务器编号
	self.IsNeedActivitCode = 0 ///默认不需要激活码
	self.OpenGMTalk = 0        ///默认不接受gm指令
}

type Server struct {
	//	Object
	addClientChannel    chan *Client     ///添加用户频道
	removeClientChannel chan *Client     ///删除用户频道
	msgBroadcastChannel chan interface{} ///广播聊天频道
	kickoutALLChannel   chan bool        ///主动踢掉所有客户端
	userMgr             *UserMgr         ///用户管理器self.config.defaultConfig() ///生成默认配置
	msgDispatch         *MsgDispatch     ///消息分发器
	loger               *Loger           ///日志记录器
	//	checker             *Checker         ///代码逻辑检查器
	config             *ServerConfig   ///服务器配置
	dynamicDB          *DBServer       ///动态数据库组件
	staticDataMgr      *StaticDataMgr  ///静态数据组件
	loginDB            *DBServer       ///数据库组件
	recordDB           *DBServer       ///行为记录数据库组件
	nextUpdateGameTime time.Time       ///服务器下次更新全局游戏逻辑时间
	waitGroup          *sync.WaitGroup ///定义关服时同步对象
	sysMailMgr         *SystemMailMgr  ///系统邮件
	activitCodeMgr     *ActivitCodeMgr ///激活码组件
	//gmMgr              *GameMasterMgr  ///GM命令管理器
	sdkMgr          *SDKMgr          ///sdk管理器
	dataCalcMgr     *DataCalcMgr     ///数据统计管理器
	serverIDChecker *net.TCPListener ///检测serverid有效性
	ipIntelnetWS    string           ///服务器外网地址与端口
}

//type IServer interface {
//	//	IObject
//	GetConfig() *ServerConfig
//	GetLoger() ILoger
//	//	GetChecker() IChecker
//	GetDynamicDB() IDBServer
//	GetStaticDataMgr() IStaticDataMgr
//	GetLoginDB() IDBServer
//	GetVersion() string
//	GetMsgDispatch() IMsgDispatch
//	Kickout(reason string, client IClient) ///踢掉一个客户端并和它断开连接
//	BroadcastMsg(msg IMsgHead)             ///广播消息
//	SendSysNotice(textSysNotice string)    ///发送系统公告
//}

var serverSingleton *Server

//func (self *Server) IsNil() bool {
//	return nil == self
//}

func (self *Server) GetConfig() *ServerConfig {
	return self.config
}

func (self *Server) GetMsgDispatch() IMsgDispatch {
	return self.msgDispatch
}

func (self *Server) GetDynamicDB() *DBServer {
	return self.dynamicDB
}

func (self *Server) GetRecordDB() *DBServer {
	return self.recordDB
}

func (self *Server) GetStaticDataMgr() *StaticDataMgr {
	return self.staticDataMgr
}

func (self *Server) GetLoginDB() *DBServer {
	return self.loginDB
}

func (self *Server) GetLoger() *Loger {
	return self.loger
}

func (self *Server) GetSystemMail() *SystemMailMgr {
	return self.sysMailMgr
}

// func (self *Server) GetGMMgr() *GameMasterMgr {
// 	return self.gmMgr
// }

func (self *Server) GetDataCalcMgr() *DataCalcMgr {
	return self.dataCalcMgr
}

func (self *Server) GetActivitCode() *ActivitCodeMgr {
	return self.activitCodeMgr
}

func (self *Server) GetSDKMgr() *SDKMgr {
	return self.sdkMgr
}

func GetServer() *Server {
	return serverSingleton
}

func (self *Server) IsOpenGMTalk() bool { ///判断此服务器是否接受gm指令
	isOpenGMTalk := self.config.OpenGMTalk > 0
	return isOpenGMTalk
}

func (self *Server) CheckMachineCode() {
	return
	machineCodeContext := ""
	macAddress := GetMacAddress()
	sqlQueryMachineCode := fmt.Sprintf("select password(UCASE('%s'))", macAddress)
	rows := self.dynamicDB.Query(sqlQueryMachineCode)
	if rows != nil {
		for rows.Next() {
			rows.Scan(&machineCodeContext)
		}
		rows.Close()
	}
	isCheckMachineCodeOK := (self.config.RunKey == machineCodeContext)
	if self.loger.CheckFail("isCheckMachineCodeOK==true", isCheckMachineCodeOK == true,
		isCheckMachineCodeOK, true) { ///检测机器码有效性
		self.loger.Fatal("isCheckMachineCodeOK==false")
	}
}

func (self *Server) cleanOnlineTable() {
	recordDB := self.GetRecordDB()
	cleanOnlineTableSQL := fmt.Sprintf("truncate %s", tableRecordOnline)
	recordDB.Exec(cleanOnlineTableSQL)
}

func (self *Server) CheckServerID() { ///检测配置文件中的sid是否填的正确
	service := fmt.Sprintf("0.0.0.0:%d", self.config.ServerID)
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		log.Fatalf("Server CheckServerID ResolveTCPAddr fail.serverid:%d,%v", self.config.ServerID, err)
	}
	self.serverIDChecker, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("Server CheckServerID ListenTCP fail.serverid:%d,%v", self.config.ServerID, err)
	}
}

func (self *Server) GetIntenetIP() string { ///启动时从数据库中取得自己的外网ip
	return self.ipIntelnetWS
}

func (self *Server) ReadIntenetIP() { ///启动时从数据库中取得自己的外网ip
	sqlQuery := fmt.Sprintf("select wsurl from %s where serverid=%d limit 1", tableServerList, self.config.ServerID)
	rows := self.GetLoginDB().Query(sqlQuery)
	if nil == rows {
		return
	}
	WSUrl := ""
	for rows.Next() {
		rows.Scan(&WSUrl)
	}
	rows.Close()
	self.ipIntelnetWS = WSUrl
}

func (self *Server) InitServer(serverIndex int) {
	self.userMgr = NewUserMgr() ///用户管理器初始化
	//self.msgDispatch = NewMsgDispatch()            ///消息分发器初始化
	logDirName := fmt.Sprintf("./log%d", serverIndex)
	self.loger = NewLoger(logDirName, logDebug, true) ///日志组件初始化
	//	self.checker = new(Checker)                    ///创建代码逻辑检查器
	self.NewServerConfig(serverIndex)    ///服务器配置信息初始化
	self.DisplayStartupInfo(serverIndex) ///显示启动信息
	//	self.InitDBServer()                                      ///数据库初始化
	self.staticDataMgr.Init()                                ///初始化静态数据管理器
	self.dynamicDB.Init(&GetServer().GetConfig().DynamicDNS) ///初始化动态数据库组件
	self.loginDB.Init(&GetServer().GetConfig().LoginDNS)     ///初始化登录数据库组件
	self.recordDB.Init(&GetServer().GetConfig().RecordDNS)   ///初始化记录数据库组件
	self.msgDispatch.initMsgDispatch()
	self.sysMailMgr.Init(self.userMgr)
	self.cleanOnlineTable() ///清理即时在线表
	//	self.gmMgr.Init()
	self.CheckMachineCode()
	self.CheckServerID()
	self.ReadIntenetIP() ///启动时从数据库中取得自己的外网ip
}

func NewServer() *Server {
	newServer := new(Server)
	newServer.addClientChannel = make(chan *Client)                                 ///添加用户频道
	newServer.removeClientChannel = make(chan *Client)                              ///删除用户频道
	newServer.msgBroadcastChannel = make(chan interface{}, msgBroadcastChannelSize) ///创建广播聊天消息频道
	newServer.kickoutALLChannel = make(chan bool)                                   ///生成关闭服务器通道
	newServer.staticDataMgr = new(StaticDataMgr)
	newServer.dynamicDB = new(DBServer)
	newServer.loginDB = new(DBServer)
	newServer.recordDB = new(DBServer)
	newServer.msgDispatch = new(MsgDispatch)
	newServer.waitGroup = new(sync.WaitGroup)
	newServer.sysMailMgr = new(SystemMailMgr)
	//newServer.gmMgr = new(GameMasterMgr) ///生成GM命令管理器
	newServer.dataCalcMgr = new(DataCalcMgr)       ///生成数据统计管理器
	newServer.activitCodeMgr = new(ActivitCodeMgr) ///生成激活码管理器
	newServer.sdkMgr = new(SDKMgr)                 ///生成sdk管理器
	return newServer
}

func (self *Server) NewServerConfig(serverIndex int) {
	configFileName := fmt.Sprintf("./serverconfig%d.json", serverIndex) ///服务器配置文件路径及文件名
	self.config = new(ServerConfig)
	configFile, err := ioutil.ReadFile(configFileName) ///尝试打开配置文件
	if err != nil {
		self.loger.Fatal("Server NewServerConfig serverIndex:%d fail!%v", serverIndex, err)

		//		self.config.defaultConfig() ///文件不存在,服务器重新建新一份生成默认配置
		//		configString, err := json.Marshal(self.config)
		//		if err != nil {
		//			self.loger.Fatal("Server Marshal fail!%v", err)
		//		}
		//		err = ioutil.WriteFile(configFileName, configString, os.ModePerm)
		//		if err != nil {
		//			self.loger.Fatal("Server WriteFile fail!%v", err)
		//		}
		//		return
	}
	err = json.Unmarshal(configFile, self.config)
	if err != nil {
		self.loger.Fatal("Server loadConfig fail!%v", err)
	}
	self.loger.logMinLevel = self.config.MinLogLevel ///配置文件日志等级生效
}

func (self *Server) Kickout(reason string, client *Client) { ///踢掉一个客户端并和它断开连接
	client.ws.Close()
}

func (self *Server) KickoutALLPlayer() { ///正常关闭服务器
	self.kickoutALLChannel <- true
}

func (self *Server) ShutDown() { ///正常关闭服务器
	self.kickoutALLChannel <- true
	self.waitGroup.Wait()
	self.loger.Fatal("Server Normal ShutDown!Cause handleINT.")
}

func (self *Server) GetVersion() string {
	return "1.2.3.4"
}

func (self *Server) GetServerInfo() *ServerInfo {
	serverInfo := new(ServerInfo)
	sqlQueryServer := fmt.Sprintf("select * from %s where serverid=%d limit 1", tableServerList, self.config.ServerID)
	self.GetLoginDB().fetchOneRow(sqlQueryServer, serverInfo)
	return serverInfo
}

func (self *Server) DisplayStartupInfo(serverIndex int) { ///显示启动信息
	self.loger.Info("********************************************************************")
	self.loger.Info("FootBall Server(SVN:%s,ServerIndex:%d,ServerID:%d,PID:%d) Is Starting...", ServerVersion, serverIndex, self.config.ServerID, os.Getpid())
}

func (self *Server) SendSysNotice(textSysNotice string) { ///发送系统公告
	msgSysNotice := new(ChatMsg)
	msgSysNotice.Type = chatTypeSystem
	msgSysNotice.Text = textSysNotice
	msgType, msgAction := msgSysNotice.GetTypeAndAction()
	msgSysNotice.FillMsgHead(0, &msgType, &msgAction)
	self.BroadcastMsg(msgSysNotice)
}

func (self *Server) BroadcastMsg(msg IMsgHead) { ///广播消息
	self.msgBroadcastChannel <- msg
}

func (self *Server) updateNextLogicUpdateTime() { ///更新下次处理服务器更新逻辑时间
	serverUpdateConfig := GetServer().GetStaticDataMgr().GetConfigStaticData(configServer, configServerCommon)
	updateLogicHours, _ := strconv.Atoi(serverUpdateConfig.Param1)   ///更新小时数
	updateLogicMinutes, _ := strconv.Atoi(serverUpdateConfig.Param2) ///更新分钟数

	currentDay := GetOpenServerSamsaraDay(ArenaSettleAccounts)

	now := time.Now()
	nowHour := now.Hour()
	nowMinute := now.Minute()
	if nowHour > updateLogicHours {
		now = now.AddDate(0, 0, ArenaSettleAccounts-currentDay) ///下一天
	} else if nowHour == updateLogicHours && nowMinute >= updateLogicMinutes {
		now = now.AddDate(0, 0, ArenaSettleAccounts-currentDay) ///下一天
	}
	self.nextUpdateGameTime = time.Date(now.Year(), now.Month(), now.Day(),
		updateLogicHours, updateLogicMinutes, 0, 0, now.Location())
	self.loger.Info("updateNextLogicUpdateTime at %v.", self.nextUpdateGameTime)
}

func (self *Server) Update(now time.Time) {
	if self.nextUpdateGameTime.IsZero() == true { ///初次启动时设置更新时间
		self.updateNextLogicUpdateTime() ///更新下次处理服务器更新逻辑时间
	}
	if now.Before(self.nextUpdateGameTime) {
		return ///更新时间未到直接返回
	}
	self.loger.Info("Start ServerUpdate at %v...", now)
	CalcArenaResult()                ///结算竞争场成绩,纯数据库操作
	self.updateNextLogicUpdateTime() ///更新下次处理服务器更新逻辑时间
	self.loger.Info("Finish ServerUpdate at %v.", time.Now())
}

//func (self *Server) IsTeamOnline(teamID int) bool { ///判断此球队是否在线
//	queryTeamOnlineSQL := fmt.Sprintf("select teamid from %s where teamid=%d limit 1", tableRecordOnline, teamID)
//	rows := GetServer().GetRecordDB().Query(queryTeamOnlineSQL)
//	if nil == rows { ///这里不可能为nil记录
//		return false
//	}
//	teamOnlineID := 0
//	for rows.Next() {
//		rows.Scan(&teamOnlineID)
//	}
//	return teamOnlineID > 0
//}

//team.AwardObject(awardTypeTicket, totalAwardNum, 0, 0)      ///给钻石
//team.AwardObject(awardTypeVipExp, moneyType.Diamonds, 0, 0) ///给vip经验
//clientObj := client.GetElement()
//clientObj.RechargeRecord(Get_VIPShopMoney, totalAwardNum)

//func (self *Server) OffLineRecharge(teamID int, actionType int, rechargeMoney int) {///离线充值

//}

//func (self *Server) RechargeRecord(accountID int, teamID int, actionType int, rechargeMoney int, totalMoney int) {
//	//money := self.team.Ticket
//	playerRechargeRecord := fmt.Sprintf("insert %s set accountid = %d ,teamid = %d ,type = %d, param = %d",
//		tableRecordAction, accountID, teamID, PlayerRecharge, rechargeMoney)

//	lastInsertRecordid, _ := GetServer().GetRecordDB().Exec(playerRechargeRecord)
//	if lastInsertRecordid <= 0 {
//		GetServer().GetLoger().Warn("insert %s recharge money error, accountid = %d", tableRecordPay, accountID)
//	}

//	//充值记录后,记录玩家获取钻石数
//	self.SetMoneyRecord(accountID, teamID, PlayerGetMoney, actionType, rechargeMoney, totalMoney)
//}

//func (self *Server) SetMoneyRecord(accountID int, teamID int, state int, action int, value int, money int) {
//	///记录客户端获取/花费钻石信息
//	clientPayRecord := fmt.Sprintf("insert %s set accountid = %d ,teamid = %d ,type = %d,action = %d, param = %d, money = %d",
//		tableRecordPay, accountID, teamID, state, action, value, money)

//	lastInsertRecordid, _ := GetServer().GetRecordDB().Exec(clientPayRecord)
//	if lastInsertRecordid <= 0 {
//		GetServer().GetLoger().Warn("insert %s error, accountid = %d", tableRecordPay, accountID)
//	}
//}

func (self *Server) kickoutRepeatLogin(teamName string) {
	msgKick := new(ChatMsg)
	msgKick.Type = chatGMKick
	msgKick.Receiver = teamName
	msgType, msgAction := msgKick.GetTypeAndAction()
	msgKick.FillMsgHead(0, &msgType, &msgAction)
	self.BroadcastMsg(msgKick)
}

func (self *Server) kickoutAll() { ///其它线程不要调用
	for _, v := range self.userMgr.clientList {
		v.ws.Close() ///关闭所有客户端的套接字
	}
}

func (self *Server) OnTime(serverUpdateTimer *time.Ticker) {
	defer func() {
		x := recover()
		if x != nil {
			GetServer().GetLoger().Error("%v\r\n%s", x, debug.Stack())
		}
	}()
	select {
	case client := <-self.addClientChannel: ///从频道中读出新连进来的客户端
		self.userMgr.addClient(client)
	case client := <-self.removeClientChannel: ///删除用户实例
		self.userMgr.removeClient(client) ///处理底层反馈上来的一个断线的客户端
	case msg := <-self.msgBroadcastChannel: ///对收到消息进行广播
		self.msgDispatch.BroadcastMsg(self.userMgr, msg)
	case now := <-serverUpdateTimer.C:
		self.Update(now) ///运行服务器更新逻辑,全局更新逻辑
	case <-self.kickoutALLChannel: ///收到主动关闭服务器的信号
		self.kickoutAll() ///主动关闭服务器
	}
}

func (self *Server) Listen() {
	self.loger.Info("start game service at %s init NumGoroutine:%d", self.config.WSHost, runtime.NumGoroutine())
	serverUpdateTimer := time.NewTicker(time.Second * 1) ///服务器更新逻辑每秒更新一次
	for {
		//select {
		//case client := <-self.addClientChannel: ///从频道中读出新连进来的客户端
		//	self.userMgr.addClient(client)
		//case client := <-self.removeClientChannel: ///删除用户实例
		//	self.userMgr.removeClient(client) ///处理底层反馈上来的一个断线的客户端
		//case msg := <-self.msgBroadcastChannel: ///对收到消息进行广播
		//	self.msgDispatch.BroadcastMsg(self.userMgr, msg)
		//case now := <-serverUpdateTimer.C:
		//	self.Update(now) ///运行服务器更新逻辑,全局更新逻辑
		//case <-self.kickoutALLChannel: ///收到主动关闭服务器的信号
		//	self.kickoutAll() ///主动关闭服务器
		//}
		self.OnTime(serverUpdateTimer)
	}
	self.loger.Info("stop game service...")
}

func getConnectHandler(s *Server) websocket.Handler {
	clientInitID := 0
	getNewClientID := func() int {
		clientInitID++
		return clientInitID
	}
	connectHandler := func(ws *websocket.Conn) {
		newClientID := getNewClientID()
		newClient := s.userMgr.getNewClient(newClientID, s, ws)
		//loger.Check1("newClient==nil", newClient == nil, newClient, nil)
		s.addClientChannel <- newClient ///通知器了服务器一个客户端连上服务
		newClient.Run()
		s.removeClientChannel <- newClient ///断线时通知服务器协程删除此客户端
	}
	return websocket.Handler(connectHandler)
}

type SignalHandler func(*chan os.Signal)

func handlePIPE(ch *chan os.Signal) {
	for {
		<-*ch
		log.Println("handlePIPE get a SIGPIPE!")
	}
}

func handleTRAP(ch *chan os.Signal) {
	for {
		<-*ch
		log.Println("handleTRAP get a SIGTRAP!")
	}
}

func handleTSTP(ch *chan os.Signal) {
	for {
		<-*ch
		log.Println("handleTSTP get a SIGTSTP!")
		GetServer().KickoutALLPlayer() ///踢掉所有玩家
	}
}

func handleINT(ch *chan os.Signal) {
	for {
		<-*ch
		log.Println("handleINT get a SIGINT!")
		GetServer().ShutDown()
	}
}

func makeHandleSignal(signalType os.Signal, handleFun SignalHandler) {
	ch := make(chan os.Signal)
	signal.Notify(ch, signalType)
	go handleFun(&ch)
}

//var handlerPtr interface{} = nil

func RunGame(serverIndex int) {
	//runtime.GOMAXPROCS(-1)                        ///通知go使用所有cpu调度协程
	makeHandleSignal(syscall.SIGPIPE, handlePIPE) ///忽略SIGPIPE信号
	makeHandleSignal(syscall.SIGTSTP, handleTSTP) ///建立即时断点机制
	makeHandleSignal(syscall.SIGTRAP, handleTRAP) ///建立即时断点机制
	makeHandleSignal(syscall.SIGINT, handleINT)   ///建立关闭服务器机制
	serverSingleton = NewServer()
	serverSingleton.InitServer(serverIndex)
	onConnectHandler := getConnectHandler(serverSingleton)
	http.Handle("/", onConnectHandler)
	http.HandleFunc("/"+Qihoo360_Name, Qihoo360Pay) ///处理360回应的支付成功消息
	http.HandleFunc("/"+Lenovo_Name, LenovoPay)     ///处理360回应的支付成功消息
	http.HandleFunc("/"+XiaoMi_Name, XiaoMiPay)     ///处理360回应的支付成功消息
	http.HandleFunc("/"+DangLe_Name, DangLePay)     ///处理360回应的支付成功消息
	//http.HandleFunc("/"+JiFeng_Name, JiFengPay)                  ///处理360回应的支付成功消息
	http.HandleFunc("/"+Tencent_Name, TencentPay)                ///处理360回应的支付成功消息
	http.HandleFunc("/"+Tencent_GetBuyToken, TencentGetBuyToken) ///处理腾迅获得交易token请求
	http.HandleFunc("/QueryServerList", QueryServerList)         ///查询服务器列表
	//	http.HandleFunc("/QueryUpdateVersion", QueryUpdateVersion) ///查询版本更新
	go serverSingleton.Listen()
	go serverSingleton.GetSystemMail().Run()
	go serverSingleton.GetDataCalcMgr().Run()
	log.Fatal(http.ListenAndServe(serverSingleton.config.WSHost, nil))
}

///天天联赛比赛计算sp
//CREATE DEFINER=`root`@`%` PROCEDURE `CalcArenaResult`()
//BEGIN
//SET @prevarena = 0;
//SET @prevgroup = 0;
//update dy_arena as a,
//	(select id,@rownum:=CASE when arenatype!=@prevarena or groupnum!=@prevgroup then @rownum:=1 else @rownum+1 end as rank,
//	@prevarena:=arenatype as prevarena,
//  @prevgroup:=groupnum as prevgroup,
//	@newaward:=CASE when @rownum=1 then 2 when @rownum between 2 and 15 then 3 else 4 end as newaward,
//	CASE when @newaward=2 and arenatype>1 then arenatype-1 when @newaward=4 and arenatype<6 then arenatype+1 else arenatype end as newtype,
//	CASE when @rownum=1 and arenatype>1 then CEIL(groupnum/3)
//  when @rownum=16 then groupnum*3-2 when @rownum=17  then groupnum*3-1 when @rownum=18  then groupnum*3 else groupnum end as newgroup,
//	arenatype as newlastrenatype
//	from dy_arena  order by arenatype asc,groupnum asc,score desc,updateutc asc) as t
//	set a.lastarenatype=t.newlastrenatype,a.arenatype=t.newtype,a.groupnum=t.newgroup,a.remainmatchcount=17,
//	a.playmask=0,a.score=0,a.wincount=0,drawcount=0,lostcount=0,updateutc=0,a.awardticket=t.newaward,a.lastrank=t.rank
//	where a.id=t.id;

//SET @prevarena = 0;
//SET @gruoupnum = 1;
//set @rownum=1;
//set @lastrownum=0;
//update dy_arena as ar,(select id,arenatype,@rownum:=CASE when arenatype!=@prevarena or @rownum>=18 then @rownum:=1 else @rownum+1 end as rank,
//@prevarena:=arenatype as prevarena,
//@gruoupnum :=CASE when @lastrownum>=18 then @gruoupnum+1 else @gruoupnum end as newgroupnum,
//@lastrownum:=@rownum as lastrownum
// from dy_arena order by arenatype) as t set ar.groupnum=t.newgroupnum where ar.id=t.id;

//END;

//CREATE DEFINER=`root`@`%` FUNCTION `func_range_string_mod`(
//    f_num INT UNSIGNED -- Total strings.
//    ) RETURNS varchar(200) CHARSET latin1
//BEGIN

//      DECLARE i INT UNSIGNED DEFAULT 0;
//      DECLARE v_result VARCHAR(200) DEFAULT '';
//      DECLARE v_dict VARCHAR(200) DEFAULT '';
//      SET v_dict = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
//      SET v_dict = LPAD(v_dict,200,v_dict);
//      WHILE i < f_num
//      DO
//SET v_result = CONCAT(v_result,SUBSTR(v_dict,CEIL(RAND()*200),1));
//        SET i = i + 1;
//      END WHILE;
//      RETURN v_result;
//END;

//CREATE DEFINER=`root`@`%` PROCEDURE `InsertActivationCode`(`codetype` int(4),`number` int(4))
//BEGIN
//    /* type = 1: 激活码   type = 2: 奖励兑换码 */
//     declare prefix varchar(64);
//     declare codeID int;
//     declare newCode varchar(255);
//     declare i int;
//     declare seed int;
//     set i = 0;
//    select IFNULL(max(id),0) into seed  from dy_activationcode;
//      loop_label:LOOP
//           /*set newCode = func_range_string_mod(6);*/
//           set newCode="";
//           select CONV(crc32(seed),10,16) into newCode;
//           set i = i + 1;
//           INSERT dy_activationcode VALUES(0, lower(newCode),codetype, 0, 0,0);
//           select last_insert_id() into seed;
//           IF i>=number THEN
//              LEAVE loop_label;
//           END IF;
//     END LOOP;
//END;
