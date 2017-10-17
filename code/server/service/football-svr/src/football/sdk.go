package football

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
)

type SDKType struct { ///sdk平台类型,对应动态表中的st_sdk
	ID             int    ///记录id
	Name           string ///平台名字
	AppID          string ///应用编号
	AppKey         string ///应用密钥
	AppSecret      string ///签名私钥
	AccessTokenURL string ///获取 access token–服务器端接口
	GetUsrInfoURL  string ///获取用户信息–服务器端接口
	UpdateTokenURL string ///刷新 access token–服务器端接口
}

type PayOrderInfo struct {
	ID        int    `json:"id"`        ///流水号
	Type      int    `json:"type"`      ///订单类型 1普通订单 2保留
	State     int    `json:"state"`     ///订单状态 1已支付 2已发货,未发货状态可能出现错误
	TeamID    int    `json:"teamid"`    ///拥有球队id
	TeamName  string `json:"teamname"`  ///球队名字
	SDKUserID string `json:"sdkyserid"` ///sdk用户id
	SDKName   string `json:"sdkname"`   ///sdk名字
	OrderID   string `json:"orderid"`   ///订单号
	ProductID int    `json:"productid"` ///商品id
	Money     int    `json:"money"`     ///支付货币数,不同平台单位可能不同,单位是分
	Balance   int    `json:"balance"`   ///完成交易后的余额
	PayTime   string `json:"paytime"`   ///支付时间
}

type SDKMgr struct { ///任务管理器
}

///艺龙世纪
//const (
//	Qihoo360_Name            = "Qihoo360"
//	Qihoo360_AppID           = "201634806"
//	Qihoo360_AppKey          = "a36ce0be05f58a5ea6c782c7872a7a73"
//	Qihoo360_AppSecret       = "f750197a7d163cada56a2a7ff3dbd44a"
//	Qihoo360_PrivateKey      = "1ad8e68a8130419961df43f3ba4a74b4"
//	Qihoo360_AccessTokenURL  = "https://openapi.360.cn/oauth2/access_token?grant_type=authorization_code&code=%s&client_id=%s&client_secret=%s&redirect_uri=oob"
//	Qihoo360_RefreshTokenURL = "https://openapi.360.cn/oauth2/refresh_token?refresh_token=%s&client_id=%s&client_secret=%s&scope=basic"
//	Qihoo360_GetUsrInfoURL   = "https://openapi.360.cn/user/me.json?access_token=%s&fields=id"
//)

//指尖sdk提供的介入参数
const (
	Qihoo360_Name            = "Qihoo360"
	Qihoo360_AppID           = "201603906"
	Qihoo360_AppKey          = "658d90fea256cbd1c901d58212ff61b7"
	Qihoo360_AppSecret       = "1d0ff35dc7ccc6dc04bd49f13a9d6c79"
	Qihoo360_PrivateKey      = "947c4afe24cd803c9c06357d3ea09377"
	Qihoo360_AccessTokenURL  = "https://openapi.360.cn/oauth2/access_token?grant_type=authorization_code&code=%s&client_id=%s&client_secret=%s&redirect_uri=oob"
	Qihoo360_RefreshTokenURL = "https://openapi.360.cn/oauth2/refresh_token?refresh_token=%s&client_id=%s&client_secret=%s&scope=basic"
	Qihoo360_GetUsrInfoURL   = "https://openapi.360.cn/user/me.json?access_token=%s&fields=id"
)

const (
	YiLongShiJi_Name = "ylsj"
)

//const (
//	Qihoo360_Name            = "Qihoo360"
//	Qihoo360_AppID           = "201603906"
//	Qihoo360_AppKey          = "658d90fea256cbd1c901d58212ff61b7"
//	Qihoo360_AppSecret       = "1d0ff35dc7ccc6dc04bd49f13a9d6c79"
//	Qihoo360_AccessTokenURL  = "https://openapi.360.cn/oauth2/access_token?grant_type=authorization_code&code=%s&client_id=%s&client_secret=%s&redirect_uri=oob"
//	Qihoo360_RefreshTokenURL = "https://openapi.360.cn/oauth2/refresh_token?refresh_token=%s&client_id=%s&client_secret=%s&scope=basic"
//	Qihoo360_GetUsrInfoURL   = "https://openapi.360.cn/user/me.json?access_token=%s&fields=id"
//)

///艺龙世纪
//const (
//	Lenovo_Name      = "Lenovo"
//	Lenovo_AppID     = "20043400000001200434"
//	Lenovo_AppKey    = "ZXHUFBWNPSSK"
//	Lenovo_AppSecret = "REU3QTBENDMwOTQ4MUQ2M0I3MTUyNUE0NEY4RDA4MkVGRTZDRTdFNU1UQTVNamMzTkRreU1qazNNVFEyTmpnek9ETXJNVFkxTnpRNU16RTJPREF5TXpneE16QTVOREkyTnpVd09UQXlNamsxTURVNU5UVTFNamc1"
//	//	Lenovo_AccessTokenURL  = "https://openapi.360.cn/oauth2/access_token?grant_type=authorization_code&code=%s&client_id=%s&client_secret=%s&redirect_uri=oob"
//	//	Lenovo_RefreshTokenURL = "https://openapi.360.cn/oauth2/refresh_token?refresh_token=%s&client_id=%s&client_secret=%s&scope=basic"
//	Lenovo_Realm            = "10000950.realm.lenovoidapps.com"
//	Lenovo_GetUsrInfoURL    = "http://passport.lenovo.com/interserver/authen/1.2/getaccountid?lpsust=%s&realm=%s"
//	Lenovo_CheckSignURL     = "http://27.17.3.254:6060/IappDecryptDemo.php?trans_data=%s&key=%s&sign=%s"
//	Lenovo_CheckSignPostURL = "http://27.17.3.254:6060/IappDecryptDemo.php"
//)

//指尖sdk
const (
	Lenovo_Name             = "Lenovo"
	Lenovo_AppID            = "20042200000002200422"
	Lenovo_AppKey           = "ZXHUFBWNPSSK"
	Lenovo_AppSecret        = "RTZGRjQyQTBDOTUzNjNBQjREQ0EyMDQ4MDYzMkNGQzg5MTM4NzhFMU1UWXpNVFV4T0RNeE5qTXdNVGswTlRnek5UY3JNVFl4T1RNNU1UVXlOekl3T1RFek1qTTJPRGd5TXpZNU16STBNakUxTWprNE56QTJNVGN4"
	Lenovo_Realm            = "10000950.realm.lenovoidapps.com"
	Lenovo_GetUsrInfoURL    = "http://passport.lenovo.com/interserver/authen/1.2/getaccountid?lpsust=%s&realm=%s"
	Lenovo_CheckSignURL     = "http://27.17.3.254:6060/IappDecryptDemo.php?trans_data=%s&key=%s&sign=%s"
	Lenovo_CheckSignPostURL = "http://27.17.3.254:6060/IappDecryptDemo.php"
)

const (
	XiaoMi_Name            = "XiaoMi"
	XiaoMi_AppID           = "26598"
	XiaoMi_AppKey          = "8d93eb1e-e3a7-1a60-0eb2-539576031a67"
	XiaoMi_CheckSessionURL = "http://mis.migc.xiaomi.com/api/biz/service/verifySession.do?%s" //appId=%s&session=%s&uid=%s&signature=%s"
	XiaoMi_CheckSessionStr = "appId=%s&session=%s&uid=%s"

//	XiaoMi_CheckSignURL     = "http://27.17.3.254:6060/IappDecryptDemo.php?trans_data=%s&key=%s&sign=%s"
//	XiaoMi_CheckSignPostURL = "http://27.17.3.254:6060/IappDecryptDemo.php"
)

const ( ///生成公式：sig=MD5(token|app_key)  （中间有“|”）
	DangLe_Name  = "DangLe"
	DangLe_AppID = "1790"
	//	DangLe_MerchantID    = "791"
	DangLe_AppKey        = "rsFnqrUD"
	DangLe_PaymentKey    = "Md0G1R6M7Xde"
	DangLe_GetUsrInfoURL = "http://connect.d.cn/open/member/info?app_id=%s&mid=%s&token=%s&sig=%s"
	DangLe_SigFormat     = "order=%s&money=%s&mid=%s&time=%s&result=%s&ext=%s&key=%s"

//	XiaoMi_CheckSignURL     = "http://27.17.3.254:6060/IappDecryptDemo.php?trans_data=%s&key=%s&sign=%s"
//	XiaoMi_CheckSignPostURL = "http://27.17.3.254:6060/IappDecryptDemo.php"
)

func (self *SDKMgr) ProcessLoginDangLe(client IClient, userName string, sdkName string, accountName *string, userID string) bool { ///处理登录消息
	loger := GetServer().GetLoger()
	clientObj := client.GetElement()
	now := Now()
	signature := CalcMD5(userName + "|" + DangLe_AppKey)
	fmt.Println(signature)
	getUsrInfoURLStr := fmt.Sprintf(DangLe_GetUsrInfoURL, DangLe_AppID, userID, userName, signature)
	fmt.Println(getUsrInfoURLStr)
	///需要urlencode
	//	getUsrInfoURLStr = url.QueryEscape(getUsrInfoURLStr)
	//	fmt.Println(getUsrInfoURLStr)
	response, err := http.Get(getUsrInfoURLStr)
	if err != nil {
		loger.Error("ProcessLoginDangLe getUsrInfoURL fail! %v", err)
		return false
	}
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	var f interface{} = nil
	err = json.Unmarshal(body, &f)
	result := f.(map[string]interface{})
	response.Body.Close()
	errCode := result["error_code"].(float64)
	error_msg := result["error_msg"]
	if loger.CheckFail("errCode==0", errCode == 0, errCode, 0) {
		loger.Warn("ProcessLoginDangLe getUsrInfoURL fail! errCode=%v error_msg=%v", errCode, error_msg)
		return false ///会话检查失败
	}
	clientObj.authCode = userName ///保存认证码
	clientObj.userID = userID     ///得到平台用户id
	*accountName = fmt.Sprintf("%s_%s", sdkName, clientObj.userID)
	accountID := client.createClientAccount(accountName, accountName) ///创建帐号
	sqlUpdateAccount := fmt.Sprintf(`update %s set sdkname='%s',sdkuserid='%s',
	authcode='%s',accesskey='%s',refreshkey='%s',accesskeyexpires=%d where id=%d`, tableAccount,
		sdkName, clientObj.userID, userName, clientObj.accessToken,
		clientObj.refreshToken, now+clientObj.expiresIn, accountID)
	GetServer().GetLoginDB().Exec(sqlUpdateAccount)
	return true
}

func (self *SDKMgr) ProcessLogin360Qihoo(client IClient, userName string, sdkName string, accountName *string) bool { ///处理登录消息
	loger := GetServer().GetLoger()
	clientObj := client.GetElement()
	accountInfo := new(AccountInfo)
	now := Now()
	accountSQL := fmt.Sprintf("select * from %s where sdkname='%s' and authcode='%s' limit 1",
		tableAccount, sdkName, userName)
	GetServer().GetLoginDB().fetchOneRow(accountSQL, accountInfo)
	if accountInfo.AuthCode == userName { ///认证码如果相同
		if now >= accountInfo.AccessKeyExpires { ///AccessKey已经过期了,使用RefreshKey更新AccessKey
			accessUrl := fmt.Sprintf(Qihoo360_RefreshTokenURL, accountInfo.RefreshKey,
				Qihoo360_AppKey, Qihoo360_AppSecret)
			response, err := http.Get(accessUrl)
			if err != nil {
				loger.Error("%v", err)
				return false
			}
			body, _ := ioutil.ReadAll(response.Body)
			var f interface{} = nil
			err = json.Unmarshal(body, &f)
			result := f.(map[string]interface{})
			response.Body.Close()
			if result["error_code"] != nil {
				loger.Error("%v", result)
				return false
			}
			accountInfo.AccessKey = result["access_token"].(string)
			accountInfo.RefreshKey = result["refresh_token"].(string)
			expiresIn, _ := strconv.Atoi(result["expires_in"].(string))
			accountInfo.AccessKeyExpires = now + expiresIn
			sqlUpdateAccount := fmt.Sprintf(`update %s set sdkname='%s',authcode='%s',
			accesskey='%s',refreshkey='%s',accesskeyexpires=%d where id=%d`, sdkName, userName,
				accountInfo.AccessKey, accountInfo.RefreshKey, now+expiresIn, accountInfo.ID)
			GetServer().GetLoginDB().Exec(sqlUpdateAccount)
		}
		*accountName = accountInfo.Name
		clientObj.userID = accountInfo.SDKUserID
		clientObj.accessToken = accountInfo.AccessKey
		clientObj.refreshToken = accountInfo.RefreshKey
		clientObj.expiresIn = accountInfo.AccessKeyExpires
		clientObj.authCode = userName ///保存认证码
		return true
	}
	///生成取得Token的URL
	accessUrl := fmt.Sprintf(Qihoo360_AccessTokenURL, userName,
		Qihoo360_AppKey, Qihoo360_AppSecret)
	fmt.Println(accessUrl) ///for debug
	response, err := http.Get(accessUrl)
	if err != nil {
		loger.Error("%v", err)
		return false
	}
	body, _ := ioutil.ReadAll(response.Body)
	var f interface{} = nil
	err = json.Unmarshal(body, &f)
	result := f.(map[string]interface{})
	response.Body.Close()
	if result["error_code"] != nil {
		loger.Error("%v", result)
		return false
	}
	///save sessionInfo
	clientObj.accessToken = result["access_token"].(string)
	clientObj.refreshToken = result["refresh_token"].(string)
	clientObj.expiresIn, _ = strconv.Atoi(result["expires_in"].(string))
	clientObj.authCode = userName ///保存认证码
	///取得用户信息
	getUserInfoUrl := fmt.Sprintf(Qihoo360_GetUsrInfoURL, clientObj.accessToken)
	response, err = http.Get(getUserInfoUrl)
	if err != nil {
		loger.Error("%v", err)
		return false
	}
	body, _ = ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &f)
	result = f.(map[string]interface{})
	response.Body.Close()
	if result["error_code"] != nil {
		loger.Error("%v", result)
		return false
	}
	clientObj.userID = result["id"].(string) ///得到平台用户id
	*accountName = fmt.Sprintf("%s_%s", sdkName, clientObj.userID)
	accountID := client.createClientAccount(accountName, accountName) ///创建帐号
	sqlUpdateAccount := fmt.Sprintf(`update %s set sdkname='%s',sdkuserid='%s',
	authcode='%s',accesskey='%s',refreshkey='%s',accesskeyexpires=%d where id=%d`, tableAccount,
		sdkName, clientObj.userID, userName, clientObj.accessToken,
		clientObj.refreshToken, now+clientObj.expiresIn, accountID)
	GetServer().GetLoginDB().Exec(sqlUpdateAccount)
	return true
}

func (self *SDKMgr) ProcessLoginXiaoMi(client IClient, userName string, sdkName string, accountName *string, userID string) bool { ///处理登录消息
	loger := GetServer().GetLoger()
	clientObj := client.GetElement()
	now := Now()
	checkSessionStr := fmt.Sprintf(XiaoMi_CheckSessionStr, XiaoMi_AppID, userName, userID)
	fmt.Println(checkSessionStr)
	signature := CalcHmac(checkSessionStr, XiaoMi_AppKey)
	fmt.Println(signature)
	checkSessionStr += "&signature=" + signature
	///可能需要urlencode
	checkSessionURL := fmt.Sprintf(XiaoMi_CheckSessionURL, checkSessionStr)
	//url, err := url.ParseQuery(checkSessionURL)
	checkSessionURI, _ := url.ParseRequestURI(checkSessionURL)
	fmt.Println(checkSessionURI.String())
	response, err := http.Get(checkSessionURI.String())
	if err != nil {
		loger.Error("ProcessLoginXiaoMi checkSession fail! %v", err)
		return false
	}
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	var f interface{} = nil
	err = json.Unmarshal(body, &f)
	result := f.(map[string]interface{})
	response.Body.Close()
	errCode := result["errcode"].(float64)
	errMsg := result["errMsg"]
	if loger.CheckFail("errCode==200", errCode == 200, errCode, 200) {
		loger.Warn("ProcessLoginXiaoMi checkSession fail! errCode=%v,errMsg=%s", errCode, errMsg)
		return false ///签名检查失败
	}
	clientObj.authCode = userName ///保存认证码
	clientObj.userID = userID     ///得到平台用户id
	*accountName = fmt.Sprintf("%s_%s", sdkName, clientObj.userID)
	accountID := client.createClientAccount(accountName, accountName) ///创建帐号
	sqlUpdateAccount := fmt.Sprintf(`update %s set sdkname='%s',sdkuserid='%s',
	authcode='%s',accesskey='%s',refreshkey='%s',accesskeyexpires=%d where id=%d`, tableAccount,
		sdkName, clientObj.userID, userName, clientObj.accessToken,
		clientObj.refreshToken, now+clientObj.expiresIn, accountID)
	GetServer().GetLoginDB().Exec(sqlUpdateAccount)
	return true
}

type Lenovo_IdentityInfo struct {
	AccountID string `xml:"AccountID"`
	Username  string `xml:"Username"`
	DeviceID  string `xml:"DeviceID"`
	Verified  string `xml:"verified"`
}

type Lenovo_Identity struct {
	lenovo_Identity Lenovo_IdentityInfo `xml:"IdentityInfo"`
}

func (self *SDKMgr) ProcessLoginLenovo(client IClient, userName string, sdkName string, accountName *string) bool { ///处理登录消息
	loger := GetServer().GetLoger()
	clientObj := client.GetElement()
	accountInfo := new(AccountInfo)
	now := Now()
	accountSQL := fmt.Sprintf("select * from %s where sdkname='%s' and authcode='%s' limit 1",
		tableAccount, sdkName, userName)
	GetServer().GetLoginDB().fetchOneRow(accountSQL, accountInfo)
	if accountInfo.AuthCode == userName { ///认证码如果相同l
		*accountName = accountInfo.Name
		clientObj.userID = accountInfo.SDKUserID
		clientObj.accessToken = accountInfo.AccessKey
		clientObj.refreshToken = accountInfo.RefreshKey
		clientObj.expiresIn = accountInfo.AccessKeyExpires
		clientObj.authCode = userName ///保存认证码
		return true
	}
	/////save sessionInfo
	clientObj.accessToken = userName
	//clientObj.refreshToken = result["refresh_token"].(string)
	//clientObj.expiresIn, _ = strconv.Atoi(result["expires_in"].(string))
	//clientObj.authCode = userName ///保存认证码
	///取得用户信息
	//	var f interface{} = nil
	getUserInfoUrl := fmt.Sprintf(Lenovo_GetUsrInfoURL, clientObj.accessToken, Lenovo_Realm)
	response, err := http.Get(getUserInfoUrl)
	fmt.Println(getUserInfoUrl)
	if err != nil {
		loger.Error("%v", err)
		return false
	}
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	identityInfo := new(Lenovo_IdentityInfo)
	err = xml.Unmarshal(body, identityInfo)
	response.Body.Close()
	if identityInfo.AccountID == "" {
		loger.Error("%v", body)
		return false
	}
	clientObj.userID = identityInfo.AccountID ///得到平台用户id
	*accountName = fmt.Sprintf("%s_%s", sdkName, clientObj.userID)
	accountID := client.createClientAccount(accountName, accountName) ///创建帐号
	sqlUpdateAccount := fmt.Sprintf(`update %s set sdkname='%s',sdkuserid='%s',
	authcode='%s',accesskey='%s',refreshkey='%s',accesskeyexpires=%d where id=%d`, tableAccount,
		sdkName, clientObj.userID, userName, clientObj.accessToken,
		clientObj.refreshToken, now+clientObj.expiresIn, accountID)
	GetServer().GetLoginDB().Exec(sqlUpdateAccount)
	return true
}

func (self *SDKMgr) ProcessLogin(client IClient, userName string, sdkName string, userID string) (string, bool) { ///处理登录消息
	accountName := ""
	loginResult := false
	// if sdkName == Qihoo360_Name {
	// 	loginResult = self.ProcessLogin360Qihoo(client, userName, sdkName, &accountName)
	// } else if sdkName == Lenovo_Name {
	// 	loginResult = self.ProcessLoginLenovo(client, userName, sdkName, &accountName)
	// } else if sdkName == XiaoMi_Name {
	// 	loginResult = self.ProcessLoginXiaoMi(client, userName, sdkName, &accountName, userID)
	// } else if sdkName == DangLe_Name {
	// 	loginResult = self.ProcessLoginDangLe(client, userName, sdkName, &accountName, userID)
	// } else if sdkName == Tencent_Name {
	// 	loginResult = TencentLogin(client, userName, sdkName, &accountName, userID)
	// } else if sdkName == YiLongShiJi_Name {
	accountName = userName
	loginResult = true
	//}
	//	} else if sdkName == JiFeng_Name {
	//		loginResult = JiFengLogin(client, userName, sdkName, &accountName, userID)
	return accountName, loginResult
}

///执行完订单后设置完成状态
func (self *SDKMgr) DonePayOrder(payOrderID int, teamName string, sdkUserID string, balance int) bool {
	if len(sdkUserID) <= 0 {
		sdkUserID = "0"
	}
	sqlUpdate := fmt.Sprintf("update %s set state=2,teamname='%s',sdkuserid='%s',balance=%d where id=%d",
		tablePayOrder, teamName, sdkUserID, balance, payOrderID)
	_, rowsAffected := GetServer().GetDynamicDB().Exec(sqlUpdate)
	return rowsAffected > 0
}

///生成订单
func (self *SDKMgr) CreatePayOrder(teamID int, product_id int, amount int, user_id string, sdkName string, orderid string) int {
	sqlInsert := fmt.Sprintf(`insert %s set teamid=%d,sdkname='%s',orderid='%s',
	productid=%d,money=%d,sdkuserid='%s'`, tablePayOrder, teamID, sdkName, orderid, product_id, amount, user_id)
	payOrderID, _ := GetServer().GetDynamicDB().Exec(sqlInsert)
	if payOrderID <= 0 {
		GetServer().GetLoger().Warn("SDKMgr CreatePayOrder fail! teamid:%d,sdkname:%s,orderid:%s ", teamID, sdkName, orderid)
	}
	return payOrderID
}

type ServerInfo struct {
	ID          int    `json:"id"`          ///记录id
	Title       string `json:"title"`       ///区服务器标题
	Name        string `json:"name"`        ///区服务器名字
	State       int    `json:"state"`       ///状态 1新服 2爆满
	WSURL       string `json:"wsurl"`       ///游戏区登录地址
	MaxPlayer   int    `json:"maxplayer"`   ///允许最大人数
	Enable      int    `json:"enable"`      ///有效开关
	ServerID    int    `json:"serverid"`    ///服务器id,例如10001,10002,与服务器配置文件中对应
	BearVersion string `json:"bearversion"` ///最低容忍版本号
	SDKLimit    string `json:"sdklimit"`    ///sdk登录限制,sdk的名字,N/A表示此服所有平台都可以登录,多sdk用分号连接
}

type ServerList []ServerInfo ///球员中心会员信息列表

func QueryServerList(w http.ResponseWriter, req *http.Request) { ///处理查询服务器列表
	defer func() {
		x := recover()
		if x != nil {
			GetServer().GetLoger().Error("%v\r\n%s", x, debug.Stack())
		}
	}()
	serverList := ServerList{}
	serverInfo := new(ServerInfo)
	fetchAllRowsSql := fmt.Sprintf("select * from %s where enable>0 limit 200", tableServerList)
	elmentList := GetServer().GetLoginDB().fetchAllRows(fetchAllRowsSql, serverInfo)
	for i := range elmentList {
		serverInfo = elmentList[i].(*ServerInfo)
		serverList = append(serverList, *serverInfo)
	}
	responseText, err := json.Marshal(serverList)
	if err != nil {
		GetServer().GetLoger().Error("QueryServerList error %v", err)
	}
	w.Write(responseText)
}

type UpdateVersionInfo struct {
	ID          int    `json:"id"`          ///记录id
	Version     string `json:"version"`     ///外部展示版本号
	FileUrl     string `json:"fileurl"`     ///更新文件url',
	ForceUpdate int    `json:"forceupdate"` ///强制更新 0非强制 1强制
	Enable      int    `json:"enable"`      ///有效开关 0表示忽略此版本
	Desc        string `json:"desc"`        ///版本更新描述
}

func LenovoPay(w http.ResponseWriter, req *http.Request) { ///联想充值回调
	defer func() {
		x := recover()
		if x != nil {
			GetServer().GetLoger().Error("%v\r\n%s", x, debug.Stack())
		}
	}()
	loger := GetServer().GetLoger()
	transData := req.FormValue("transdata")
	signData := req.FormValue("sign")
	//fmt.Println(transData)
	//transData = strings.Replace(transData, `"waresid":1`, `"waresid":2`, -1)
	//fmt.Println(transData)
	///检测签名
	data := make(url.Values)
	data["trans_data"] = []string{transData}
	data["key"] = []string{Lenovo_AppSecret}
	data["sign"] = []string{signData}
	fmt.Println(transData, Lenovo_AppSecret, signData)
	res, err := http.PostForm(Lenovo_CheckSignPostURL, data)
	if err != nil {
		w.Write([]byte("FAILURE")) ///通过平台失败
		fmt.Println(err.Error())
		return
	}
	body, _ := ioutil.ReadAll(res.Body)

	checkSignResult := string(body)
	fmt.Println(checkSignResult)
	if loger.CheckFail("checkSignResult==SUCCESS", checkSignResult == "SUCCESS", checkSignResult, "SUCCESS") {
		w.Write([]byte("FAILURE")) ///通过平台失败
		return                     ///签名检查失败
	}
	var f interface{} = nil
	err = json.Unmarshal([]byte(transData), &f)
	if err != nil {
		loger.Error("%v", err)
		w.Write([]byte("FAILURE")) ///通过平台失败
		return
	}
	result := f.(map[string]interface{})

	appID := result["appid"].(string)
	if loger.CheckFail("appID==Lenovo_AppID", appID == Lenovo_AppID, appID, Lenovo_AppID) {
		//w.Write([]byte("FAILURE")) ///通过平台失败
		return ///appkey不存在
	}

	order_id := result["transid"].(string)
	if loger.CheckFail("order_id!={empty}", order_id != "", order_id, "{empty}") {
		//w.Write([]byte("FAILURE")) ///通过平台失败
		return ///订单号不存在
	}

	///平台用户编号不提供
	user_id := "N/A" // req.PostFormValue("user_id")

	payResult := result["result"].(float64)
	if loger.CheckFail("payResult<=0", payResult <= 0, payResult, 0) {
		//w.Write([]byte("FAILURE")) ///通过平台失败
		return ///支付不是成功标识
	}

	transtype := result["transtype"].(float64)
	if loger.CheckFail("transtype<=0", transtype <= 0, transtype, 0) {
		//w.Write([]byte("FAILURE")) ///通过平台失败
		return ///只接受交易,拒绝冲正
	}

	product_id := result["cpprivate"].(string)
	productID, _ := strconv.Atoi(product_id)
	//if loger.CheckFail(" productID>0", productID > 0, productID, 0) {
	//	w.Write([]byte("FAILURE")) ///通过平台失败
	//	return                     ///商品编号
	//}

	amount := result["money"].(float64)
	money := int(amount) // strconv.Atoi(amount) ///平台货币
	//if loger.CheckFail(" money>0", money > 0, money, 0) {
	//	w.Write([]byte("FAILURE")) ///通过平台失败
	//	return                     ///商品编号money不存在或是非法
	//}

	app_uid := result["exorderno"].(string)
	//if loger.CheckFail("app_uid!={empty}", app_uid != "", app_uid, "{empty}") {
	//	w.Write([]byte("FAILURE")) ///通过平台失败
	//	return                     ///平台用户id不存在
	//}

	teamID, _ := strconv.Atoi(app_uid)
	//if loger.CheckFail(" teamID>0", teamID > 0, teamID, 0) {
	//	w.Write([]byte("FAILURE")) ///通过平台失败
	//	return                     ///球队id不存在
	//}

	sdkname := "Lenovo"
	payOrderID := GetServer().GetSDKMgr().CreatePayOrder(teamID, productID, money, user_id, sdkname, order_id)
	if loger.CheckFail("payOrderID>0", payOrderID > 0, payOrderID, 0) {
		//w.Write([]byte("FAILURE")) ///通过平台失败
		return ///重复订单
	}

	vipShopBuyDiamondMsg := new(VipShopBuyDiamondMsg)
	vipShopBuyDiamondMsg.MoneyID = productID
	vipShopBuyDiamondMsg.TeamID = teamID
	vipShopBuyDiamondMsg.PayOrderID = payOrderID
	vipShopBuyDiamondMsg.PayMoney = money
	GetServer().BroadcastMsg(vipShopBuyDiamondMsg) ///转给服务器处理
	w.Write([]byte("SUCCESS"))                     ///先让sdk滚
}

func Qihoo360Pay(w http.ResponseWriter, req *http.Request) { ///处理登录消息
	defer func() {
		x := recover()
		if x != nil {
			GetServer().GetLoger().Error("%v\r\n%s", x, debug.Stack())
		}
	}()
	w.Write([]byte("ok")) ///先让sdk滚
	loger := GetServer().GetLoger()
	fmt.Println(req.URL.String())
	app_key := req.FormValue("app_key")
	if loger.CheckFail("app_key==Qihoo360_AppKey", app_key == Qihoo360_AppKey, app_key, Qihoo360_AppKey) {
		return ///appkey不存在
	}

	product_id := req.FormValue("product_id")
	productID, _ := strconv.Atoi(product_id)
	//if loger.CheckFail(" productID>0", productID > 0, productID, 0) {
	//	return ///productID非法
	//}

	//staticDataMgr := GetServer().GetStaticDataMgr()
	//moneyType := staticDataMgr.GetMoneyType(productID)
	//if loger.CheckFail("moneyType!=nil", moneyType != nil, moneyType, nil) {
	//	return //商品类型不存在
	//}
	amount := req.FormValue("amount")
	money, _ := strconv.Atoi(amount) ///平台货币
	//if loger.CheckFail(" money>0", money > 0, money, 0) {
	//	return ///money不存在或是非法
	//}

	//moneyPrice := moneyType.Money * 100 ///转换将人民币"元"价格转换成"分"价格
	//if loger.CheckFail("money==moneyPrice", money == moneyPrice, money, moneyPrice) {
	//	return //商品价格与平台扣费不符,客户端被篡改?回调消息被破解?
	//}

	order_id := req.FormValue("order_id")
	//if loger.CheckFail("order_id!={empty}", order_id != "", order_id, "{empty}") {
	//	return ///订单号不存在
	//}

	app_uid := req.FormValue("app_uid")
	//if loger.CheckFail("app_uid!={empty}", app_uid != "", app_uid, "{empty}") {
	//	return ///平台用户id不存在
	//}
	user_id := req.FormValue("user_id")
	teamID, _ := strconv.Atoi(app_uid)
	//if loger.CheckFail(" teamID>0", teamID > 0, teamID, 0) {
	//	return ///球队id不存在
	//}

	gateway_flag := req.FormValue("gateway_flag")
	if loger.CheckFail(" gateway_flag==success", gateway_flag == "success", gateway_flag, "success") {
		return ///支付不是成功标识
	}

	sign := req.FormValue("sign")
	if loger.CheckFail("sign!={empty}", sign != "", sign, "{empty}") {
		return ///签名不能为空
	}

	digitArrage := []string{}
	for k, v := range req.Form {
		if len(v) > 0 && k != "sign" && k != "sign_return" {
			digitArrage = append(digitArrage, k)
		}
	}
	//	fmt.Println(digitArrage)
	sort.Strings(digitArrage)
	//	fmt.Println(digitArrage)
	digitText := ""
	for i := range digitArrage {
		fieldName := digitArrage[i]
		fieldText := req.FormValue(fieldName)
		digitText += fieldText + "#"
	}
	digitText += Qihoo360_AppSecret
	signMine := CalcMD5(digitText)
	if loger.CheckFail("signMine==sign", signMine == sign, signMine, sign) {
		return ///签名无效
	}
	//	app_ext1 := req.FormValue("app_ext1")

	//	sign_type := req.FormValue("sign_type")
	//	gateway_flag := req.FormValue("gateway_flag")
	//
	//	sign_return := req.FormValue("sign_return")
	sdkname := "Qihoo360"
	payOrderID := GetServer().GetSDKMgr().CreatePayOrder(teamID, productID, money, user_id, sdkname, order_id)
	if loger.CheckFail("payOrderID>0", payOrderID > 0, payOrderID, 0) {
		return ///重复订单
	}
	vipShopBuyDiamondMsg := new(VipShopBuyDiamondMsg)
	vipShopBuyDiamondMsg.MoneyID = productID
	vipShopBuyDiamondMsg.TeamID = teamID
	vipShopBuyDiamondMsg.PayOrderID = payOrderID
	vipShopBuyDiamondMsg.PayMoney = money
	GetServer().BroadcastMsg(vipShopBuyDiamondMsg) ///转给服务器处理
}

func XiaoMiPay(w http.ResponseWriter, req *http.Request) { ///处理小米支付
	defer func() {
		x := recover()
		if x != nil {
			GetServer().GetLoger().Error("%v\r\n%s", x, debug.Stack())
		}
	}()
	loger := GetServer().GetLoger()
	fmt.Println(req.URL.String())
	payCallBackURL := req.FormValue("cpUserInfo")    ///得到支付回调地址
	valueArray := strings.Split(payCallBackURL, "#") ///拆出真实支付回调地址和玩家球队id
	teamID := 0
	if len(valueArray) >= 2 {
		payCallBackURL = valueArray[0]                                            ///真实回调地址
		payCallBackURL = strings.Replace(payCallBackURL, "/"+XiaoMi_Name, "", -1) ///去掉多余的sdkname
		teamID, _ = strconv.Atoi(valueArray[1])                                   ///得到球队id
	}
	isEnd := strings.Contains(payCallBackURL, req.Host) ///判断此通知是否到达终点处理
	if false == isEnd {                                 ///需要分发给真正的处理地址
		realPayCallBackURL := payCallBackURL + req.URL.String()
		realPayCallBackURI, _ := url.ParseRequestURI(realPayCallBackURL)
		fmt.Println(realPayCallBackURI.String())
		responseRelay, errRelay := http.Get(realPayCallBackURI.String())
		if errRelay != nil {
			loger.Warn("XiaoMiPay Relay payCallBackURL fail! %v", errRelay)
		}
		bodyRelay, _ := ioutil.ReadAll(responseRelay.Body)
		w.Write(bodyRelay)
		return
	}
	///正式开始处理
	w.Write([]byte(`{"errcode":200}`)) ///先让sdk滚
	///取表单数据
	appId := req.FormValue("appId")         ///游戏ID
	cpOrderId := req.FormValue("cpOrderId") ///开发商订单ID
	//cpUserInfo := req.FormValue("cpUserInfo")///开发商透传信息
	uid := req.FormValue("uid")                 ///用户ID
	orderId := req.FormValue("orderId")         ///游戏平台订单ID
	orderStatus := req.FormValue("orderStatus") ///订单状态，TRADE_SUCCESS 代表成功
	payFee := req.FormValue("payFee")           ///支付金额,单位为分,即0.01 米币。
	//	productCode := req.FormValue("productCode")           ///商品代码
	//	productName := req.FormValue("productName")           ///商品名称
	//	productCount := req.FormValue("productCount")         ///商品数量
	//	payTime := req.FormValue("payTime")                   ///支付时间,格式 yyyy-MM-dd HH:mm:ss
	//	orderConsumeType := req.FormValue("orderConsumeType") ///订单类型：10：普通订单11：直充直消订单
	signature := req.FormValue("signature") ///签名,签名方法见后面说明
	if loger.CheckFail("appId==XiaoMi_AppID", appId == XiaoMi_AppID, appId, XiaoMi_AppID) {
		return ///appId不存在
	}
	if loger.CheckFail("orderId!={nil}", orderId != "", orderId, "{nil}") {
		return ///orderId非法
	}
	if loger.CheckFail("orderStatus==TRADE_SUCCESS", orderStatus == "TRADE_SUCCESS", orderStatus, "TRADE_SUCCESS") {
		return ///orderStatus非法
	}
	payFeeMoney, _ := strconv.ParseFloat(payFee, 32)
	if loger.CheckFail("payFeeMoney>0", payFeeMoney > 0, payFeeMoney, 0) {
		return ///payFeeMoney非法
	}
	///验证签名
	digitText, _ := url.QueryUnescape(req.URL.String())
	digitText = strings.Replace(digitText, "/"+XiaoMi_Name+"?", "", -1)
	digitText = strings.Replace(digitText, "&signature="+signature, "", -1)
	fmt.Println(digitText)
	signMine := CalcHmac(digitText, XiaoMi_AppKey)
	if loger.CheckFail("signMine==signature", signMine == signature, signMine, signature) {
		return ///签名无效
	}
	productID, _ := strconv.Atoi(cpOrderId)
	money := int(payFeeMoney)
	user_id := uid
	payOrderID := GetServer().GetSDKMgr().CreatePayOrder(teamID, productID, money, user_id, XiaoMi_Name, orderId)
	if loger.CheckFail("payOrderID>0", payOrderID > 0, payOrderID, 0) {
		return ///重复订单
	}
	vipShopBuyDiamondMsg := new(VipShopBuyDiamondMsg)
	vipShopBuyDiamondMsg.MoneyID = productID
	vipShopBuyDiamondMsg.TeamID = teamID
	vipShopBuyDiamondMsg.PayOrderID = payOrderID
	vipShopBuyDiamondMsg.PayMoney = money
	GetServer().BroadcastMsg(vipShopBuyDiamondMsg) ///转给服务器处理
}

func DangLePay(w http.ResponseWriter, req *http.Request) { ///处理当乐支付
	defer func() {
		x := recover()
		if x != nil {
			GetServer().GetLoger().Error("%v\r\n%s", x, debug.Stack())
		}
	}()
	loger := GetServer().GetLoger()
	fmt.Println(req.URL.String())
	payCallBackURL := req.FormValue("ext")           ///得到支付回调地址
	ext := req.FormValue("ext")                      ///得到支付回调地址
	valueArray := strings.Split(payCallBackURL, "#") ///拆出真实支付回调地址和玩家球队id
	teamID := 0
	productID := 0
	if len(valueArray) >= 3 {
		payCallBackURL = valueArray[0]                                            ///真实回调地址
		payCallBackURL = strings.Replace(payCallBackURL, "/"+DangLe_Name, "", -1) ///去掉多余的sdkname
		teamID, _ = strconv.Atoi(valueArray[1])                                   ///得到球队id
		productID, _ = strconv.Atoi(valueArray[2])                                ///套餐
	}
	isEnd := strings.Contains(payCallBackURL, req.Host) ///判断此通知是否到达终点处理
	if false == isEnd {                                 ///需要分发给真正的处理地址
		realPayCallBackURL := payCallBackURL + req.URL.String()
		realPayCallBackURI, _ := url.ParseRequestURI(realPayCallBackURL)
		fmt.Println(realPayCallBackURI.String())
		responseRelay, errRelay := http.Get(realPayCallBackURI.String())
		if errRelay != nil {
			loger.Warn("DangLePay Relay payCallBackURL fail! %v", errRelay)
		}
		bodyRelay, _ := ioutil.ReadAll(responseRelay.Body)
		w.Write(bodyRelay)
		return
	}
	///正式开始处理
	w.Write([]byte("success")) ///先让sdk滚
	///取表单数据
	//appId := req.FormValue("appId") ///游戏ID
	//cpOrderId := req.FormValue("cpOrderId") ///开发商订单ID
	//cpUserInfo := req.FormValue("cpUserInfo")///开发商透传信息
	uid := req.FormValue("mid")       ///本次支付用户的乐号，既登录后返回的 mid 参数。
	orderId := req.FormValue("order") ///游戏平台订单ID
	result := req.FormValue("result") ///支付结果，固定值。“1”代表成功，“0”代表失败
	payFee := req.FormValue("money")  ///支付金额，单位：元。
	time := req.FormValue("time")     ///时间戳，格式：yyyymmddHH24mmss 月日小时分秒小于 10 前面补充 0
	//	productCode := req.FormValue("productCode")           ///商品代码
	//	productName := req.FormValue("productName")           ///商品名称
	//	productCount := req.FormValue("productCount")         ///商品数量
	//	payTime := req.FormValue("payTime")                   ///支付时间,格式 yyyy-MM-dd HH:mm:ss
	//	orderConsumeType := req.FormValue("orderConsumeType") ///订单类型：10：普通订单11：直充直消订单
	signature := req.FormValue("signature") ///MD5 验证串，用于与接口生成的验证串做比较，保证计费通知的合法性。
	//if loger.CheckFail("appId==XiaoMi_AppID", appId == XiaoMi_AppID, appId, XiaoMi_AppID) {
	//	return ///appId不存在
	//}
	if loger.CheckFail("orderId!={nil}", orderId != "", orderId, "{nil}") {
		return ///orderId非法
	}
	if loger.CheckFail("result==1", result == "1", result, "1") {
		return ///支付状态
	}
	payFeeMoney, _ := strconv.ParseFloat(payFee, 32)
	if loger.CheckFail("payFeeMoney>0", payFeeMoney > 0, payFeeMoney, 0) {
		return ///payFeeMoney非法
	}
	///验证签名
	//DangLe_SigFormat="order=%s&money=%s&mid=%s&time=%s&result=%s&ext=%s&key=%s"
	digitText := fmt.Sprintf(DangLe_SigFormat, orderId, payFee, uid, time, result, ext, DangLe_PaymentKey)
	fmt.Println(digitText)
	signMine := CalcMD5(digitText)
	if loger.CheckFail("signMine==signature", signMine == signature, signMine, signature) {
		return ///签名无效
	}
	//productID, _ := strconv.Atoi(cpOrderId)
	money := int(payFeeMoney) * 100 ///将当乐的元单位转换成分单位
	user_id := uid
	payOrderID := GetServer().GetSDKMgr().CreatePayOrder(teamID, productID, money, user_id, DangLe_Name, orderId)
	if loger.CheckFail("payOrderID>0", payOrderID > 0, payOrderID, 0) {
		return ///重复订单
	}
	vipShopBuyDiamondMsg := new(VipShopBuyDiamondMsg)
	vipShopBuyDiamondMsg.MoneyID = productID
	vipShopBuyDiamondMsg.TeamID = teamID
	vipShopBuyDiamondMsg.PayOrderID = payOrderID
	vipShopBuyDiamondMsg.PayMoney = money
	GetServer().BroadcastMsg(vipShopBuyDiamondMsg) ///转给服务器处理
}
