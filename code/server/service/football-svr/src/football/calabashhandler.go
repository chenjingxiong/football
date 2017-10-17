package football

import ()

type QueryCalabashInfoMsg struct { /// 请求葫芦娃的领取状态
	MsgHead `json:"head"` ///"calabash", "queryinfo"
}

func (self *QueryCalabashInfoMsg) GetTypeAndAction() (string, string) { ///实现消息处理接口的返回类型方法
	return "calabash", "queryinfo"
}

func (self *QueryCalabashInfoMsg) processAction(client IClient) bool { ///实现消息处理接口的处理消息方法
	if self.checkAction(client) == false { ///检测
		return false
	}
	if self.payAction(client) == false { ///支付
		return false
	}
	if self.doAction(client) == false { ///发货
		return false
	}
	return true
}

///检测
func (self *QueryCalabashInfoMsg) checkAction(client IClient) bool {
	return true
}

///支付
func (self *QueryCalabashInfoMsg) payAction(client IClient) bool {
	return true
}

///发货
func (self *QueryCalabashInfoMsg) doAction(client IClient) bool {
	timeUNIX := Now()
	team := client.GetTeam()
	team.CheckLoginAward(timeUNIX)

	//QueryCalabashInfoResultMsg := new(QueryCalabashInfoResultMsg)
	//QueryCalabashInfoResultMsg.CthState, QueryCalabashInfoResultMsg.CthNeed = team.GetLoginAndPayAwardState(CthBit, timeUNIX)                ///建号3小时奖励状态
	//QueryCalabashInfoResultMsg.LTwoState, QueryCalabashInfoResultMsg.LTwoNeed = team.GetLoginAndPayAwardState(LtwoBit, timeUNIX)             ///登录2天奖励状态
	//QueryCalabashInfoResultMsg.LFiveState, QueryCalabashInfoResultMsg.LFiveNeed = team.GetLoginAndPayAwardState(LfiveBit, timeUNIX)          ///登录5天奖励状态
	//QueryCalabashInfoResultMsg.LTenState, QueryCalabashInfoResultMsg.LTenNeed = team.GetLoginAndPayAwardState(LtenBit, timeUNIX)             ///登录10天奖励状态
	//QueryCalabashInfoResultMsg.LFifteenState, QueryCalabashInfoResultMsg.LFifteenNeed = team.GetLoginAndPayAwardState(LfifteenBit, timeUNIX) ///登录15天奖励状态
	//QueryCalabashInfoResultMsg.LThirtyState, QueryCalabashInfoResultMsg.LThirtyNeed = team.GetLoginAndPayAwardState(LthirtyBit, timeUNIX)    ///登录30天奖励状态
	//QueryCalabashInfoResultMsg.InitPayState, _ = team.GetLoginAndPayAwardState(InitpayBit, timeUNIX)                                         ///首充奖励状态
	//QueryCalabashInfoResultMsg.Buy1980State, _ = team.GetLoginAndPayAwardState(Buy1980Bit, timeUNIX)                                         ///购买1980钻石奖励状态

	calabashstatelist := CalabashStateList{}
	tcalabashctate := CalabashState{}
	for i := 0; i < LfifteenBit; i++ {
		tcalabashctate.AwardState, tcalabashctate.AwardNeed = team.GetLoginAndPayAwardState(i+1, timeUNIX)
		calabashstatelist = append(calabashstatelist, tcalabashctate)
	}

	queryCalabashInfoResultMsg := new(QueryCalabashInfoResultMsg)
	queryCalabashInfoResultMsg.CalabashStates = calabashstatelist
	client.SendMsg(queryCalabashInfoResultMsg) ///发送给客户端返回信息
	return true
}

type CalabashState struct {
	AwardState int ///葫芦娃奖励状态，0不可领取1可领取2已领取
	AwardNeed  int ///还差的秒数 已达到为0
}

type CalabashStateList []CalabashState ///声明状态结构数组

type QueryCalabashInfoResultMsg struct { /// 返回葫芦娃的领取状态
	MsgHead `json:"head"` ///"calabash", "queryinforesult "

	CalabashStates CalabashStateList `json:"calabashstates"` ///葫芦娃状态列表

	//CthState int `json:"cthstate"` ///建号3小时奖励状态， 0不可领取1可领取2已领取
	//CthNeed  int `json:"cthneed"`  ///建号3小时仍需条件，还差的秒数，已达到为0

	//LTwoState int `json:"ltowstate"` ///登录2天奖励状态， 0不可领取1可领取2已领取
	//LTwoNeed  int `json:" ltowneed"` ///登录2天仍需条件，还差的天数，已达到为0

	//LFiveState int `json:"lfivestate"` ///登录5天奖励状态， 0不可领取1可领取2已领取
	//LFiveNeed  int `json:" lfiveneed"` ///登录5天仍需条件，还差的天数，已达到为0

	//LTenState int `json:"ltenstate"` ///登录10天奖励状态， 0不可领取1可领取2已领取
	//LTenNeed  int `json:" ltenneed"` ///登录10天仍需条件，还差的天数，已达到为0

	//LFifteenState int `json:"lfifteenstate"` ///登录15天奖励状态， 0不可领取1可领取2已领取
	//LFifteenNeed  int `json:" lfifteenneed"` ///登录15天仍需条件，还差的天数，已达到为0

	//LThirtyState int `json:"lthirtystate"` ///登录30天奖励状态， 0不可领取1可领取2已领取
	//LThirtyNeed  int `json:" lthirtyneed"` ///登录30天仍需条件，还差的天数，已达到为0

	//InitPayState int `json:"initpaystate"` ///首充奖励状态， 0不可领取1可领取2已领取

	//Buy1980State int `json:"buy1980state"` ///购买1980钻石奖励状态， 0不可领取1可领取2已领取

}

func (self *QueryCalabashInfoResultMsg) GetTypeAndAction() (string, string) { ///实现消息处理接口的返回类型方法
	return "calabash", "queryinforesult"
}

type AccpetCalabashMsg struct { /// 请求领取葫芦娃消息
	MsgHead      `json:"head"` ///" calabash ", "accpetaward"
	CalabashSite int           `json:"calabashsite"` ///葫芦娃的位置号
}

func (self *AccpetCalabashMsg) GetTypeAndAction() (string, string) { ///实现消息处理接口的返回类型方法
	return "calabash", "accpetaward"
}

func (self *AccpetCalabashMsg) processAction(client IClient) bool { ///实现消息处理接口的处理消息方法
	if self.checkAction(client) == false { ///检测
		return false
	}
	if self.payAction(client) == false { ///支付
		return false
	}
	if self.doAction(client) == false { ///发货
		return false
	}
	return true
}

///检测
func (self *AccpetCalabashMsg) checkAction(client IClient) bool {
	if self.CalabashSite <= 0 || self.CalabashSite > LfifteenBit {
		GetServer().GetLoger().Warn("AccpetCalabashMsg.CalabashSite error, CalabashSite = %d  teamid = %d", self.CalabashSite, client.GetTeam().ID)
		return false
	} else {
		return true
	}
}

///支付
func (self *AccpetCalabashMsg) payAction(client IClient) bool {
	return true
}

///发货
func (self *AccpetCalabashMsg) doAction(client IClient) bool {

	team := client.GetTeam()

	accpetCalabashResultMsg := new(AccpetCalabashResultMsg)
	accpetCalabashResultMsg.ResultCode = team.GetLoginAndPayAward(self.CalabashSite)
	accpetCalabashResultMsg.CalabashIndex = self.CalabashSite
	client.SendMsg(accpetCalabashResultMsg) ///发送给客户端返回信息
	return true
}

type AccpetCalabashResultMsg struct { /// 返回领取葫芦娃结果
	MsgHead       `json:"head"` ///" calabash ", "accpetresult"
	ResultCode    int           `json:"result"`        ///返回结果编号1成功，0不能领取，-1已领取，-2球队满
	CalabashIndex int           `json:"calabashindex"` ///返回之前请求领取的葫芦娃位置
}

func (self *AccpetCalabashResultMsg) GetTypeAndAction() (string, string) { ///实现消息处理接口的返回类型方法
	return "calabash", "accpetresult"
}
