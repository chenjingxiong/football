package football

import (
	"encoding/json"
	//	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime/debug"
	//	"sort"
	"strconv"
	"strings"
)

const ( ///生成公式：sig=MD5(token|app_key)  （中间有“|”）
	JiFeng_Name          = "JiFeng"
	JiFeng_AppID         = "25809660"
	JiFeng_AppKey        = "rsFnqrUD"
	JiFeng_PaymentKey    = "Md0G1R6M7Xde"
	JiFeng_GetUsrInfoURL = "http://connect.d.cn/open/member/info?app_id=%s&mid=%s&token=%s&sig=%s"
	JiFeng_SigFormat     = "order=%s&money=%s&mid=%s&time=%s&result=%s&ext=%s&key=%s"

	JiFeng_VerifyTokenURL = "http://api.gfan.com/uc1/common/verify_token?token=%s"

//	XiaoMi_CheckSignPostURL = "http://27.17.3.254:6060/IappDecryptDemo.php"
)

///处理机锋登录消息
func JiFengLogin(client IClient, userName string, sdkName string, accountName *string, userID string) bool {
	loger := GetServer().GetLoger()
	clientObj := client.GetElement()
	now := Now()
	verifyTokenURL := fmt.Sprintf(JiFeng_VerifyTokenURL, userName)
	fmt.Println(verifyTokenURL)
	response, err := http.Get(verifyTokenURL)
	if err != nil {
		loger.Error("JiFengLogin verifyTokenURL fail! %v", err)
		return false
	}
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	var f interface{} = nil
	err = json.Unmarshal(body, &f)
	result := f.(map[string]interface{})
	response.Body.Close()
	resultCode := result["resultCode"].(float64)
	uid := result["uid"]
	sdkUserID := fmt.Sprintf("%v", uid)
	if loger.CheckFail("resultCode==1", resultCode == 1, resultCode, 1) {
		return false ///会话检查失败
	}
	if loger.CheckFail("userID==sdkUserID", userID == sdkUserID, userID, sdkUserID) {
		return false ///userid检查失败
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

func JiFengPay(w http.ResponseWriter, req *http.Request) { ///处理机锋支付
	defer func() {
		x := recover()
		if x != nil {
			GetServer().GetLoger().Error("%v\r\n%s", x, debug.Stack())
		}
	}()
	responeOK := `<response>
						< ErrorCode>1</ErrorCode>
						< ErrorDesc></ErrorDesc>
				  </response>`
	loger := GetServer().GetLoger()
	fmt.Println(req.URL.String())
	payCallBackURL := req.FormValue("Order_id")      ///得到支付回调地址
	valueArray := strings.Split(payCallBackURL, "#") ///拆出真实支付回调地址和玩家球队id
	teamID := 0
	productID := 0
	if len(valueArray) >= 4 {
		payCallBackURL = valueArray[0]                                            ///真实回调地址
		payCallBackURL = strings.Replace(payCallBackURL, "/"+JiFeng_Name, "", -1) ///去掉多余的sdkname
		teamID, _ = strconv.Atoi(valueArray[1])                                   ///得到球队id
		productID, _ = strconv.Atoi(valueArray[2])                                ///套餐
		//		timeNow, _ = strconv.Atoi(valueArray[3])                                  ///支付时间
	}
	isEnd := strings.Contains(payCallBackURL, req.Host) ///判断此通知是否到达终点处理
	if false == isEnd {                                 ///需要分发给真正的处理地址
		realPayCallBackURL := payCallBackURL + req.URL.String()
		realPayCallBackURI, _ := url.ParseRequestURI(realPayCallBackURL)
		fmt.Println(realPayCallBackURI.String())
		responseRelay, errRelay := http.Get(realPayCallBackURI.String())
		if errRelay != nil {
			loger.Warn("JiFengPay Relay payCallBackURL fail! %v", errRelay)
			return
		}
		bodyRelay, _ := ioutil.ReadAll(responseRelay.Body)
		w.Write(bodyRelay)
		return
	}
	///正式开始处理
	w.Write([]byte(responeOK)) ///先让sdk滚
	///取表单数据
	uid := req.FormValue("mid")        ///本次支付用户的乐号，既登录后返回的 mid 参数。
	orderId := req.FormValue("order")  ///游戏平台订单ID
	appkey := req.FormValue("appkey")  ///支付结果，固定值。“1”代表成功，“0”代表失败
	payFee := req.FormValue("cost")    ///支付金额，单位：元。
	time := req.FormValue("time")      ///时间戳，格式：yyyymmddHH24mmss 月日小时分秒小于 10 前面补充 0
	signature := req.FormValue("sign") ///MD5 验证串，用于与接口生成的验证串做比较，保证计费通知的合法性。
	//if loger.CheckFail("appId==XiaoMi_AppID", appId == XiaoMi_AppID, appId, XiaoMi_AppID) {
	//	return ///appId不存在
	//}
	if loger.CheckFail("orderId!={nil}", orderId != "", orderId, "{nil}") {
		return ///orderId非法
	}
	if loger.CheckFail("appkey==JiFeng_AppKey", appkey == JiFeng_AppKey, appkey, JiFeng_AppKey) {
		return ///支付状态
	}
	payFeeMoney, _ := strconv.ParseFloat(payFee, 32)
	if loger.CheckFail("payFeeMoney>0", payFeeMoney > 0, payFeeMoney, 0) {
		return ///payFeeMoney非法
	}
	///验证签名
	signMine := CalcMD5(JiFeng_AppID + time)
	fmt.Println(signMine, signature)
	if loger.CheckFail("signMine==signature", signMine == signature, signMine, signature) {
		return ///签名无效l
	}
	money := int(payFeeMoney)
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
