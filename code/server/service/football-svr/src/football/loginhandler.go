package football

import (
//"fmt"
//	"io/ioutil"
//	"net/http"
)

type RegisteringAccount struct { //! 注册账号协议
	MsgHead  `json:"head"` //! "regeist", "registeringaccount"
	UserName string        `json:"username"` //! 用户名
	UserPwd  string        `json:"password"` //! 密码
}

func (self *RegisteringAccount) IsNeedTeamHandle() bool { ///此消息不需要team创建
	return false
}

func (self *RegisteringAccount) GetTypeAndAction() (string, string) {
	return "regeist", "registeringaccount"
}

func (self *RegisteringAccount) checkAction(client IClient) bool {

	//! 账号位数必须在6-16位以内
	//fmt.Println("Len(UserName)", len(self.UserName), "Len(UserPwd)", len(self.UserPwd))
	if len(self.UserName) < 6 {
		SendLoginResult(client, 8)
		return false
	}

	if len(self.UserName) > 16 {
		SendLoginResult(client, 1)
		return false
	}

	//! 密码位数必须在6-18位以内
	if len(self.UserPwd) < 6 {
		SendLoginResult(client, 9)
		return false
	}

	if len(self.UserPwd) > 16 {
		SendLoginResult(client, 2)
		return false
	}

	//! 账号密码只允许为英文与数字
	//! 账号密码不需要区分大小写
	for i := 0; i < len(self.UserName); i++ {
		if (self.UserName[i] >= 'a' && self.UserName[i] <= 'z') ||
			(self.UserName[i] >= 'A' && self.UserName[i] <= 'Z') ||
			(self.UserName[i] >= '0' && self.UserName[i] <= '9') {
			continue
		}

		SendLoginResult(client, 4) //! 账号非法
		return false
	}

	for i := 0; i < len(self.UserPwd); i++ {
		if (self.UserPwd[i] >= 'a' && self.UserPwd[i] <= 'z') ||
			(self.UserPwd[i] >= 'A' && self.UserPwd[i] <= 'Z') ||
			(self.UserPwd[i] >= '0' && self.UserPwd[i] <= '9') {
			continue
		}

		SendLoginResult(client, 5) //! 密码非法
		return false
	}

	accountID, _ := client.checkClientLogin(&self.UserName, &self.UserPwd) ///检测此用户是否存在
	if accountID != 0 {
		SendLoginResult(client, 3)
		return false //! 账号已存在
	}

	return true
}

func (self *RegisteringAccount) doAction(client IClient) bool {
	//! 创建账号
	client.createClientAccount(&self.UserName, &self.UserPwd)

	SendLoginResult(client, 0) //! 成功创建
	return true
}

func (self *RegisteringAccount) processAction(client IClient) bool {
	if self.checkAction(client) == false { ///检测
		return false
	}
	if self.doAction(client) == false { ///发货
		return false
	}
	return true
}

func SendLoginResult(client IClient, code int) {
	msg := new(RegisteringAccountResult)
	msg.Result = code
	client.SendMsg(msg)
}

type RegisteringAccountResult struct { //! 注册账号返回结果
	MsgHead `json:"head"` //! "regeist", "registeringaccountresult"
	Result  int           //! 0为注册成功 1账号超长 2密码超长 3账号已存在
	//! 4账号非法 5密码非法(只允许英语字母与数字) 8账号过短 9密码过短
	//! 登录时错误码: 6密码错误  7账号不存在
}

func (self *RegisteringAccountResult) GetTypeAndAction() (string, string) {
	return "regeist", "registeringaccountresult"
}

type SDKLoginMsg struct { ///专供平台sdk使用的登录消息
	MsgHead       `json:"head"`
	UserName      string `json:"username"`      ///用户名/平台登录上下文
	SDKName       string `json:"sdkname"`       ///渠道编号,ucweb,91,qq等等
	UserID        string `json:"userid"`        ///平台的userid,某些平台要求验证
	ClientVersion string `json:"clientversion"` ///客户端版本号,期望服务端能向下兼容协议
	ActivitCode   string `json:"activitcode"`   ///激活码
}

func (self *SDKLoginMsg) GetTypeAndAction() (string, string) {
	return "login", "sdklogin"
}

func (self *SDKLoginMsg) IsNeedTeamHandle() bool { ///此消息不需要team创建
	return false
}

///比较两个版本号,0表示完成相等,>0表示目标版本号大于源版本号,<0表示目标版本号小于源版本号
func (self *SDKLoginMsg) checkAction(client IClient) bool {
	loger := GetServer().GetLoger()
	serverInfo := GetServer().GetServerInfo()
	result := CompareVersion(serverInfo.BearVersion, self.ClientVersion)
	if loger.CheckFail("result>=0", result >= 0, result, 0) {
		client.SendErrorMsg(failLogin, GetModErrorText(failErrorVersion))
		return false ///客户端版本号大于服务器
	}
	////strings.self.ClientVersion
	//result := client.checkClientVersion(&self.ClientVersion) ///检测客户版本号
	//if false == result {
	//	GetServer().Kickout("error version!", client.GetElement()) ///客户端版本号错误,踢它下线
	//	return false                                               ///消息处理函数错误
	//}
	return true
}

func (self *SDKLoginMsg) doAction(client IClient) (result bool) {
	sdkMgr := GetServer().GetSDKMgr()
	AccountName, isOK := sdkMgr.ProcessLogin(client, self.UserName, self.SDKName, self.UserID)
	if false == isOK {
		return false
	}
	// PassWords := AccountName
	// if self.SDKName == YiLongShiJi_Name {
	PassWords := self.SDKName
	//}

	// accountID, isOK := client.checkClientLogin(&AccountName, &PassWords) ///检测此用户是否存在
	// if accountID == 0 {
	// 	SendLoginResult(client, 7)
	// 	return false //! 账号不存在
	// }

	// if isOK == false {
	// 	SendLoginResult(client, 6)
	// 	return false //! 密码错误
	// }

	accountID := client.createClientAccount(&AccountName, &PassWords) ///首先检测此用户是否存在
	ok := client.LoadTeam(accountID)                                  ///尝试创建球队
	if ok != false {
		//这里还需要修复球队创建时间的方法*****************
		client.SendTeam() ///有球队直接发送球队信息及球员信息
		client.LoginRecord(PlayerLogin)
		client.GetTeam().CheckLoginAward(Now())
	} else {
		client.SendLoginResultMsg(loginResultCreateTeam) ///没球队
	}
	return true
}

func (self *SDKLoginMsg) processAction(client IClient) (result bool) {
	if self.checkAction(client) == false { ///检测
		return false
	}
	if self.doAction(client) == false { ///发货
		return false
	}
	return true
}

type ClientLoginMsg struct {
	MsgHead       `json:"head"`
	UserName      string `json:"username"`      ///用户名
	PassWord      string `json:"password"`      ///密码
	channelID     int    `json:"channelid"`     ///渠道编号,ucweb,91,qq等等
	ClientVersion string `json:"clientversion"` ///客户端版本号,期望服务端能向下兼容协议
	ActivitCode   string `json:"activitcode"`   ///激活码
}

const (
	loginResultFail            = "loginfail"       ///登录失败
	loginResultCreateTeam      = "createteam"      ///要求创建队伍
	loginResultNeedActivitCode = "needactivitcode" ///需求激活码
)

type LoginResultMsg struct {
	MsgHead       `json:"head"`
	LoginResult   string `json:"loginresult"`   ///登录结果
	AccountID     int    `json:"accountid"`     ///帐号id
	ServerVersion string `json:"serverversion"` ///服务端版本号
}

func (self *LoginResultMsg) GetTypeAndAction() (string, string) {
	return "login", "loginresult"
}

func NewLoginResultMsg(loginResult *string) *LoginResultMsg {
	msg := new(LoginResultMsg)
	//msg.MsgType = "login"
	//msg.Action = "loginresult"
	//msg.CreateTime = int(time.Now().Unix())
	msg.ServerVersion = GetServer().GetVersion()
	msg.LoginResult = *loginResult
	return msg
}

type LoginHandler struct {
	MsgHandler
}

func (self *ClientLoginMsg) GetTypeAndAction() (string, string) {
	return "login", "clientlogin"
}

func (self *ClientLoginMsg) IsNeedTeamHandle() bool { ///此消息不需要team创建
	return false
}

func (self *ClientLoginMsg) CheckNameAndPassWordLegal() (bool, string) { ///检查名字与密码合法性

	if len(self.UserName) <= 0 || len(self.PassWord) <= 0 {
		return false, failInvalidAccountID
	}

	return true, ""
}

func GetModErrorText(strText string) string {
	strErrorText := "未知错误"
	if "loginfail" == strText {
		strErrorText = "登录失败"
	} else if "PasswordWrong" == strText {
		strErrorText = "登录密码错误!"
	} else if "ActivitCodeError" == strText {
		strErrorText = "激活码错误,激活码不存在或已被使用!"
	} else if failErrorVersion == strText {
		strErrorText = "版本号过低,请更新至最新版本."
	}
	return strErrorText
}

func (self *ClientLoginMsg) processAction(client IClient) bool {
	//loger := GetServer().GetLoger()
	//result := client.checkClientVersion(&self.ClientVersion) ///检测客户版本号
	//if false == result {
	//	GetServer().Kickout("error version!", client.GetElement()) ///客户端版本号错误,踢它下线
	//	return false                                               ///消息处理函数错误
	//}

	//isLegalNameAndPassWord, errorCode := self.CheckNameAndPassWordLegal()
	//if loger.CheckFail("isLegalNameAndPassWord == true", isLegalNameAndPassWord == true, isLegalNameAndPassWord, true) {
	//	client.SendErrorMsg(failLogin, errorCode)
	//	GetServer().Kickout("Invalid AccountID", client.GetElement()) ///客户端版本号错误,踢它下线
	//	return false                                                  ///名字/密码不合法
	//}

	//accountID, _ := client.checkClientLogin(&self.UserName, &self.PassWord)
	//if accountID == 0 {
	//	///判断当前服务器激活码
	//	//	staticDataMgr := GetServer().GetStaticDataMgr()
	//	//		isNeedActivatCode := staticDataMgr.GetConfigStaticDataInt(configServer, configServerCommon, 4)
	//	isNeedActivatCode := GetServer().config.IsNeedActivitCode
	//	if isNeedActivatCode != 0 {
	//		if len(self.ActivitCode) <= 0 {
	//			result := loginResultNeedActivitCode
	//			msg := NewLoginResultMsg(&result)
	//			client.SendMsg(msg) ///索要激活码
	//			return false
	//		}

	//		// loginDB := GetServer().GetLoginDB()
	//		// activitCode := new(ActivitCode)
	//		// activitCodeQuery := fmt.Sprintf("select * from dy_activationcode where code = '%s'", self.ActivitCode)
	//		// isExist := loginDB.fetchOneRow(activitCodeQuery, activitCode)

	//		// if loger.CheckFail("CodeisExist == true", isExist == true, isExist, true) {
	//		// 	return false //无此激活码
	//		// }
	//		// if activitCode.State != 0 {
	//		// 	loger.Warn("The activitCode is be using by accountID: %d", activitCode.UserAccountID)
	//		// 	return false //激活码已被使用
	//		// }

	//		//activitCodeMgr := GetServer().GetActivitCode()
	//		//			activitCodeCanUse := activitCodeMgr.CanUse(self.ActivitCode)
	//		//if loger.CheckFail("activitCodeCanUse == true", activitCodeCanUse == true, activitCodeCanUse, true) {
	//		//	client.SendErrorMsg(failLogin, GetModErrorText(failActivitCode))
	//		//	return false ///激活码不可用
	//		//}

	//		accoundID := client.createClientAccount(&self.UserName, &self.PassWord) ///创建帐号,并标记激活码已被使用

	//		//activitCodeMgr.Use(self.ActivitCode, accoundID)

	//		// activitCode.State = activitCodeBeUsed
	//		// activitCode.UserAccountID = accoundID
	//		// //激活码被使用,即时存储
	//		// activitCodeUpdate := fmt.Sprintf("update %s set state = %d, useraccountid = %d where id = %d", tableActivationCode, activitCode.State, activitCode.UserAccountID, activitCode.ID)
	//		// loginDB.Exec(activitCodeUpdate)
	//	} else {
	//		client.createClientAccount(&self.UserName, &self.PassWord) ///for debug,自动创建账号
	//	}
	//}

	//accountID, ok := client.checkClientLogin(&self.UserName, &self.PassWord)
	//if 0 == accountID { ///用户不存在错误
	//	client.SendErrorMsg(failLogin, failAccountNotExsit)
	//	return false
	//}
	//if false == ok { ///用户密码错误
	//	client.SendErrorMsg(failLogin, GetModErrorText(failPasswordWrong))
	//	return false
	//}
	//ok = client.LoadTeam(accountID) ///尝试创建球队
	//if ok != false {
	//	client.SendTeam() ///有球队直接发送球队信息及球员信息
	//	client.LoginRecord(PlayerLogin)
	//} else {
	//	client.SendLoginResultMsg(loginResultCreateTeam) ///没球队
	//}
	return true
}

func (self *LoginHandler) getName() string { ///返回可处理的消息类型
	return "login"
}

func (self *LoginHandler) initHandler() { ///初始化消息处理器
	self.addActionToList(new(ClientLoginMsg))
}
