package football

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
)

//以下是超玩的
//	Tencent_AppID       = "1101766949"
//	Tencent_AppKey      = "EcFyc5HKsd5g4tD5"

const ( ///腾迅应用宝sdk &sig=%s
	Tencent_Name        = "Tencent"
	Tencent_AppID       = "1101766949"
	Tencent_AppKey      = "EcFyc5HKsd5g4tD5"
	Tencent_GetBuyToken = "TencentGetBuyToken"
	//	Tencent_GetUsrInfoURL = "http://openapi.tencentyun.com/v3/user/get_info?openid=%s&openkey=%s&pf=%s&appid=%s&format=json&userip=10.0.0.1"
	//	Tencent_IsLoginURL = "http://119.147.19.43/v3/user/is_login?openid=%s&openkey=%s&pf=%s&appid=%s&format=json&userip=10.0.0.1"
	Tencent_GetUserInfoURL = "https://graph.qq.com/user/get_user_info?oauth_consumer_key=%s&access_token=%s&openid=%s&format=json"
	//	Tencent_GetUserInfoQuery = ""
	//	Tencent_IsLoginOAUT = "/v3/user/is_login"

	Tencent_SignFormat = "appid=%s&format=json&openid=%s&openkey=%s&pf=%s&userip=10.0.0.1"
	Tencent_MoneyType  = "钻石*VIP货币"
	//Tencent_MoneyType          = "money*coin"
	//	Tencent_BuyGoodsURL        = "https://119.147.19.43/mpay/buy_goods_m?openid=%s&openkey=%s&pf=%s&pfkey=%s&pay_token=%s&ts=%s&payitem=%s&goodsmeta=%s&goodsurl=&zoneid=1000&app_metadata=%s&format=json&appid=%s&appmode=1"
	Tencent_BuyGoodsURL        = "https://openapi.tencentyun.com/mpay/buy_goods_m?openid=%s&openkey=%s&pf=%s&pfkey=%s&pay_token=%s&ts=%s&payitem=%s&goodsmeta=%s&goodsurl=&zoneid=1000&app_metadata=%s&format=json&appid=%s&appmode=1"
	Tencent_BuyGoodsOAUT       = "/mpay/buy_goods_m"
	Tencent_BuyGoodsSignFormat = "app_metadata=%s&appid=%s&appmode=1&format=json&goodsmeta=%s&goodsurl=&openid=%s&openkey=%s&pay_token=%s&payitem=%s&pf=%s&pfkey=%s&ts=%s&zoneid=1000"

//	JiFeng_PaymentKey    = "Md0G1R6M7Xde"

//	JiFeng_SigFormat     = "order=%s&money=%s&mid=%s&time=%s&result=%s&ext=%s&key=%s"

//	JiFeng_VerifyTokenURL = "http://api.gfan.com/uc1/common/verify_token?token=%s"

//	XiaoMi_CheckSignPostURL = "http://27.17.3.254:6060/IappDecryptDemo.php"
)

///处理腾迅应用宝登录消息
func TencentLogin(client IClient, userName string, sdkName string, accountName *string, userID string) bool {
	//	userName = "10736143F57AB12FE5648B986ACA0E38"
	//	sdkName = "Tencent"
	//	accountName = "test"
	//	userID = "2DB24AC1684846D15EC49184E39A24B6#qzone_m"
	loger := GetServer().GetLoger()
	clientObj := client.GetElement()
	now := Now()
	//	valueArray := strings.Split(userID, "#") ///拆出真实支付回调地址和玩家球队id
	//	userID = valueArray[0]
	//	pfText := valueArray[1]
	//isLoginURL := fmt.Sprintf(Tencent_IsLoginURL, userID, userName, pfText, Tencent_AppID)
	//fmt.Println(isLoginURL)
	//signFormat := fmt.Sprintf(Tencent_SignFormat, Tencent_AppID, userID, userName, pfText)
	//fmt.Println(signFormat)
	//signFormatURI := url.QueryEscape(signFormat)
	//urlEncode := url.QueryEscape(Tencent_IsLoginOAUT)
	////源串是由3部分内容用“&”拼接起来的： HTTP请求方式 & urlencode(uri) & urlencode(a=x&b=y&...)
	//signFormatURI = "GET&" + urlEncode + "&" + signFormatURI
	//fmt.Println(signFormatURI)
	//signKey := Tencent_AppKey + "&"
	//signText := CalcHmacRaw(signFormatURI, signKey)
	//signFinal := base64.StdEncoding.EncodeToString(signText)
	//signFinal = url.QueryEscape(signFinal)
	//isLoginURL += "&sig=" + signFinal
	//fmt.Println(isLoginURL)
	//	Tencent_GetUserInfoURL="https://graph.qq.com/user/get_user_info?oauth_consumer_key=%s&access_token=%s&openid=%s&format=json"

	getUserInfoURL := fmt.Sprintf(Tencent_GetUserInfoURL, Tencent_AppID, userName, userID)
	//	fmt.Println(getUserInfoURL)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	peer := &http.Client{Transport: tr}
	response, err := peer.Get(getUserInfoURL)
	if err != nil {
		loger.Error("TencentLogin getUsrInfoURL fail! %v", err)
		return false
	}
	body, _ := ioutil.ReadAll(response.Body)
	//	fmt.Println(string(body))
	var f interface{} = nil
	err = json.Unmarshal(body, &f)
	result := f.(map[string]interface{})
	response.Body.Close()
	ret := result["ret"].(float64)
	if loger.CheckFail("ret==0", ret == 0, ret, 0) {
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

///处理腾迅应用宝购买秘钥获取流程
func TencentGetBuyToken(w http.ResponseWriter, req *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			GetServer().GetLoger().Error("%v\r\n%s", x, debug.Stack())
		}
	}()
	loger := GetServer().GetLoger()
	fmt.Println(req.URL.String())
	///取表单数据
	openID := req.FormValue("openid")   ///与APP通信的用户key
	openKey := req.FormValue("openkey") ///session key
	pf := req.FormValue("pf")           ///应用的来源平台
	pfKey := req.FormValue("pfkey")     ///表示平台的信息加密串，根据openid，openkey，pf，appid等生成。
	//	ts := req.FormValue("ts")           ///linux时间戳，以秒为单位。
	payItem := req.FormValue("payitem")   ///请使用ID*price*num的格式，ID表示物品ID，price表示单价,num数量
	payToken := req.FormValue("paytoken") ///1表示用户不可以修改物品数量，2表示用户可以选择购买物品的数量。
	target := req.FormValue("target")     ///需要真实的回调地址
	//	productID:= req.FormValue("productid")     ///需要真实的回调地址
	//	teamID:= req.FormValue("userid")     ///需要真实的回调地址
	//	goodsMeta := req.FormValue("goodsmeta") ///物品信息，格式必须是“name*des”，给出该商品的名称和描述
	//	goodsurl := req.FormValue("goodsurl")   ///物品的图片url，用户购买物品的确认支付页面将显示该物品图片。长度<=512字符,注意图片规格最大为：116*116 px。
	//	zoneid := req.FormValue("zoneid")       ///如果应用不分区，请输入0。

	//	openID = "2DB24AC1684846D15EC49184E39A24B6" ///与APP通信的用户key
	//	openKey = "10736143F57AB12FE5648B986ACA0E38"
	//	pf = "desktop_m_qq-10000144-android-2002-" ///应用的来源平台
	//	pfKey = "f51a43d69405841a26c2fc470ac36c2b"
	//	payToken = "30066082E609490DEC72B1487E432052"
	ts := fmt.Sprint(Now())
	//	payItem = "1*60*1"
	//	appmode = "1"

	//	goodsurl = ""
	//	zoneid = "0"
	//Tencent_BuyGoodsURL = `http://openapi.tencentyun.com/v3/pay/buy_goods?openid=%s&openkey=%s&pf=%s
	//						   &appid=%s&format=json&pfkey=%s&ts=%s&payitem=%s
	//						  &appmode=1&goodsmeta=%s&goodsurl=&zoneid=0`
	//Tencent_BuyGoodsSignFormat = `appid=%s&appmode=1&format=json&goodsmeta=%s&goodsurl=&openid=%s
	//								&openkey=%s&payitem=%s&pf=%s&pfkey=%s&ts=%s&zoneid=0`
	//	Tencent_BuyGoodsURL = "https://119.147.19.43/mpay/buy_goods_m?openid=%s&openkey=%s&pf=%s&pfkey=%s&pay_token=%s&ts=%s&payitem=%s&goodsmeta=%s&goodsurl=&zoneid=0&app_metadata=customkey&format=json&appid=%s&appmode=1"
	//target = ""
	//app_metadata := url.QueryEscape(target)
	monetyType := url.QueryEscape(Tencent_MoneyType)
	app_metadata := url.QueryEscape(target)
	BuyGoodsURL := fmt.Sprintf(Tencent_BuyGoodsURL, openID, openKey, pf, pfKey, payToken, ts, payItem, monetyType, app_metadata, Tencent_AppID)
	fmt.Println(BuyGoodsURL)

	//	Tencent_BuyGoodsSignFormat = "app_metadata=customkey&appid=%s&appmode=1&format=json&goodsmeta=%s&goodsurl=&openid=%s&openkey=%s&pay_token=%s&payitem=%s&pf=%s&pfkey=%s&ts=%s&zoneid=0"
	//	goodsMeta := url.QueryEscape("钻石*VIP货币")
	signFormat := fmt.Sprintf(Tencent_BuyGoodsSignFormat, target, Tencent_AppID, Tencent_MoneyType, openID, openKey, payToken, payItem, pf, pfKey, ts)
	fmt.Println(signFormat)
	signFormatURI := url.QueryEscape(signFormat)
	urlEncode := url.QueryEscape(Tencent_BuyGoodsOAUT)
	//源串是由3部分内容用“&”拼接起来的： HTTP请求方式 & urlencode(uri) & urlencode(a=x&b=y&...)
	signFormatURI = "GET&" + urlEncode + "&" + signFormatURI
	//if err != nil {
	//	fmt.Println(err)
	//	return false
	//}
	fmt.Println(signFormatURI)
	signKey := Tencent_AppKey + "&"
	signText := CalcHmacRaw(signFormatURI, signKey)
	signFinal := base64.StdEncoding.EncodeToString(signText)
	signFinal = url.QueryEscape(signFinal)
	BuyGoodsURL += "&sig=" + signFinal
	//fmt.Println(getUsrInfoURL)
	//	getUsrInfoURI, _ := url.ParseRequestURI(getUsrInfoURL)
	//	getUsrInfoURIText := getUsrInfoURI.String()
	fmt.Println(BuyGoodsURL)
	//response, err := http.Get(BuyGoodsURL)

	//cookieJar.SetCookies()

	req, err := http.NewRequest("GET", BuyGoodsURL, nil)
	if err != nil {
		loger.Error("TencentGetBuyToken BuyGoodsURL fail! %v", err)
		return
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cookieJar, _ := cookiejar.New(nil)
	expire := time.Now().AddDate(0, 0, 1)
	cookieValue := url.QueryEscape("openid")
	cookie1 := http.Cookie{Name: "session_id", Value: cookieValue, Path: "/", Expires: expire, MaxAge: 86400}
	cookieValue = url.QueryEscape("kp_actoken")
	cookie2 := http.Cookie{Name: "session_type", Value: cookieValue, Path: "/", Expires: expire, MaxAge: 86400}
	cookieValue = url.QueryEscape("/mpay/buy_goods_m")
	cookie3 := http.Cookie{Name: "org_loc", Value: cookieValue, Path: "/", Expires: expire, MaxAge: 86400}
	cookies := []*http.Cookie{&cookie1, &cookie2, &cookie3}
	cookieJar.SetCookies(req.URL, cookies)
	client := &http.Client{
		Jar:       cookieJar,
		Transport: tr,
	}

	response, errDo := client.Do(req)
	if errDo != nil {
		loger.Error("TencentLogin getUsrInfoURL fail! %v", err)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	var f interface{} = nil
	err = json.Unmarshal(body, &f)
	result := f.(map[string]interface{})
	response.Body.Close()
	ret := result["ret"].(float64)
	if loger.CheckFail("ret==0", ret == 0, ret, 0) {
		w.Write([]byte("ok")) ///先让sdk滚
		return                ///会话检查失败
	}
	url_params := result["url_params"].(string)
	w.Write([]byte(url_params))
}

//http://27.17.3.254:8080/Tencent?amt=600&appid=1101766949&appmeta=http://27.17.3.254:8080/Tencent*qdqb*qq&billno=-APPDJSX21271-20140715-1933429831&clientver=android&openid=277C2A4815C4B244B0594C8E4190E7BB&payamt_coins=0&payitem=1*60*1&providetype=5&pubacct_payamt_coins=&token=54EAFAD33CE43C875C484544A35169CF13172&ts=1405424023&version=v3&zoneid=1000&sig=H3z4Lp4RsKEbBJ0HpFATw4vjf4I%3D
func TencentPay(w http.ResponseWriter, req *http.Request) { ///处理机锋支付
	defer func() {
		x := recover()
		if x != nil {
			GetServer().GetLoger().Error("%v\r\n%s", x, debug.Stack())
		}
	}()
	loger := GetServer().GetLoger()
	fmt.Println(req.URL.String())
	payCallBackURL := req.FormValue("appmeta")       ///得到支付回调地址
	valueArray := strings.Split(payCallBackURL, "*") ///拆出真实支付回调地址和玩家球队id
	teamID := 0
	productID := 0
	if len(valueArray) >= 5 {
		payCallBackURL = valueArray[0]                                             ///真实回调地址
		payCallBackURL = strings.Replace(payCallBackURL, "/"+Tencent_Name, "", -1) ///去掉多余的sdkname
		teamID, _ = strconv.Atoi(valueArray[1])                                    ///得到球队id
		productID, _ = strconv.Atoi(valueArray[2])                                 ///套餐
	}
	intenetIP := GetServer().GetIntenetIP()
	isEnd := strings.Contains(payCallBackURL, intenetIP) ///判断此通知是否到达终点处理
	if false == isEnd {                                  ///需要分发给真正的处理地址
		realPayCallBackURL := payCallBackURL + req.URL.String()
		realPayCallBackURI, _ := url.ParseRequestURI(realPayCallBackURL)
		fmt.Println("TencentPay Relay", realPayCallBackURI.String())
		responseRelay, errRelay := http.Get(realPayCallBackURI.String())
		if errRelay != nil {
			loger.Warn("TencentPay Relay payCallBackURL fail! %v", errRelay)
			return
		}
		bodyRelay, _ := ioutil.ReadAll(responseRelay.Body)
		w.Write(bodyRelay)
		return
	}
	/////正式开始处理
	w.Write([]byte(`{"ret":0,"msg":"OK"}`)) ///先让sdk滚

	///取表单数据
	openid := req.FormValue("openid")   ///本次支付用户的乐号，既登录后返回的 mid 参数。
	appid := req.FormValue("appid")     ///游戏平台订单ID
	ts := req.FormValue("ts")           ///支付结果，固定值。“1”代表成功，“0”代表失败
	payitem := req.FormValue("payitem") ///支付金额，单位：元。
	token := req.FormValue("token")     ///时间戳，格式：yyyymmddHH24mmss 月日小时分秒小于 10 前面补充 0
	billno := req.FormValue("billno")   ///MD5 验证串，用于与接口生成的验证串做比较，保证计费通知的合法性。
	version := req.FormValue("version") ///本次支付用户的乐号，既登录后返回的 mid 参数。
	zoneid := req.FormValue("zoneid")   ///游戏平台订单ID
	//	providetype := req.FormValue("providetype")   ///支付结果，固定值。“1”代表成功，“0”代表失败
	amt := req.FormValue("amt")                   ///支付金额，单位：元。
	payamt_coins := req.FormValue("payamt_coins") ///时间戳，格式：yyyymmddHH24mmss 月日小时分秒小于 10 前面补充 0
	//	pubacct_payamt_coins := req.FormValue("pubacct_payamt_coins") ///时间戳，格式：yyyymmddHH24mmss 月日小时分秒小于 10 前面补充 0
	sig := req.FormValue("sig") ///MD5 验证串，用于与接口生成的验证串做比较，保证计费通知的合法性。
	if loger.CheckFail("appId==Tencent_AppID", appid == Tencent_AppID, appid, Tencent_AppID) {
		return ///appId不存在
	}
	if loger.CheckFail("billno!={nil}", billno != "", billno, "{nil}") {
		return ///orderId非法
	}
	if loger.CheckFail("zoneid==1000", zoneid == "1000", zoneid, "1000") {
		return ///支付状态
	}
	if loger.CheckFail("version==v3", version == "v3", zoneid, "v3") {
		return ///支付状态
	}
	if loger.CheckFail("openid!={nil}", openid != "", openid, "nil") {
		return ///支付状态
	}
	if loger.CheckFail("teamID>0", teamID > 0, teamID, 0) {
		return ///支付状态
	}
	if loger.CheckFail("productID>0", productID > 0, productID, 0) {
		return ///支付状态
	}
	tsUnix, _ := strconv.Atoi(ts)
	if loger.CheckFail("tsUnix>0", tsUnix > 0, tsUnix, 0) {
		return ///支付状态
	}
	if loger.CheckFail("payitem!={nil}", payitem != "", payitem, "nil") {
		return ///支付状态
	}
	payArray := strings.Split(payitem, "*") ///拆出真实支付回调地址和玩家球队id
	payArrayLen := len(payArray)
	if loger.CheckFail("payArrayLen==3", payArrayLen == 3, payArrayLen, 3) {
		return ///支付状态
	}
	//if loger.CheckFail("payitem!={nil}", payitem != "", payitem, "nil") {
	//	return ///支付状态
	//}
	if loger.CheckFail("token!={nil}", token != "", token, "nil") {
		return ///支付状态
	}
	//if loger.CheckFail("providetype==0", providetype == "0", providetype, "0") {
	//	return ///支付状态
	//}
	if loger.CheckFail("amt!={nil}", amt != "", amt, "nil") {
		return ///支付状态
	}
	if loger.CheckFail("payamt_coins!={nil}", payamt_coins != "", payamt_coins, "nil") {
		return ///支付状态
	}
	//if loger.CheckFail("pubacct_payamt_coins!={nil}", pubacct_payamt_coins != "", pubacct_payamt_coins, "nil") {
	//	return ///支付状态
	//}
	if loger.CheckFail("sig!={nil}", sig != "", sig, "nil") {
		return ///支付状态
	}
	digitArrage := []string{}
	for k, v := range req.Form {
		if len(v) > 0 && k != "sig" {
			digitArrage = append(digitArrage, k)
		}
	}
	//	fmt.Println(digitArrage)
	sort.Strings(digitArrage)
	//	fmt.Println(digitArrage)
	digitText := ""
	for i := range digitArrage {
		if i > 0 {
			digitText += "&"
		}
		fieldName := digitArrage[i]
		fieldText := req.FormValue(fieldName)
		fieldText = EncodeValue(fieldText)
		digitText += fieldName + "=" + fieldText
	}
	//fmt.Println(digitText)
	sigA := url.QueryEscape("GET")
	sigB := url.QueryEscape("/" + Tencent_Name)
	sigC := url.QueryEscape(digitText)
	digitText = sigA + "&" + sigB + "&" + sigC
	//fmt.Println(digitText)
	//	digitText = url.QueryEscape(digitText)
	//	fmt.Println(digitText)
	signKey := Tencent_AppKey + "&"
	signText := CalcHmacRaw(digitText, signKey)
	signMine := base64.StdEncoding.EncodeToString(signText)
	//fmt.Println(signMine)
	//sig = signMine
	if loger.CheckFail("signMine==sign", signMine == sig, signMine, sig) {
		return ///签名无效
	}
	//	money, _ := strconv.Atoi(amt)
	money, _ := strconv.Atoi(payArray[1])
	money = money * 10 ///转换成0.1Q点单位
	user_id := openid
	orderId := billno
	payOrderID := GetServer().GetSDKMgr().CreatePayOrder(teamID, productID, money, user_id, Tencent_Name, orderId)
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
