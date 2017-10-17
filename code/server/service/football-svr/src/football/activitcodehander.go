package football

import (
	"strings"
)

type ActivitCodeMsg struct { ///激活码
	MsgHead         `json:"head"` //"activitcode", "award"
	ActivitCodeInfo string        `json:"activitcodeinfo"` //
}

func (self *ActivitCodeMsg) GetTypeAndAction() (string, string) {
	return "activitcode", "award"
}

type ActivitCodeResultMsg struct {
	MsgHead   `json:"head"` //"activitcode", "News"
	ReturnMsg string        `json:"ReturnMsg"` ///接受到激活码后向客户端返回的消息
}

func (self *ActivitCodeResultMsg) GetTypeAndAction() (string, string) {
	return "activitcode", "News"
}

func (self *ActivitCodeMsg) checkAction(client IClient) bool {
	loger := loger() ///定义日志对象
	activitCodeMgr := GetServer().GetActivitCode()
	activeCode := activitCodeMgr.GetActiveCode(self.ActivitCodeInfo)
	if loger.CheckFail("activeCode!=nil", activeCode != nil, activeCode, nil) {
		return false //激活码必须存在
	}
	if loger.CheckFail("activeCode.State==0", activeCode.State == 0, activeCode.State, 0) {
		return false //激活码必须为未用状态
	}
	activeCodeAwardType := activitCodeMgr.GetActiveCodeAwardType(activeCode.Type)
	if loger.CheckFail("activeCodeAwardType!=nil", activeCodeAwardType != nil, activeCodeAwardType, nil) {
		return false
	}
	if activeCodeAwardType.SDKName != "N/A" { ///判断是否需要指定sdk
		clientObj := client.GetElement()
		isMatchSDK := strings.Contains(clientObj.accountName, activeCodeAwardType.SDKName)
		if loger.CheckFail("isMatchSDK==true", isMatchSDK == true, clientObj.accountName, activeCodeAwardType.SDKName) {
			return false
		}
	}
	///判断玩家是否已领取此类奖励
	team := client.GetTeam()
	isAccpetedAward := TestMask64(team.ActivitCodeAward, activeCodeAwardType.ID)
	if loger.CheckFail("isAccpetedAward == false", isAccpetedAward == false, isAccpetedAward, false) {
		return false
	}
	return true
}

func (self *ActivitCodeMsg) payAction(client IClient) bool {
	activitCodeMgr := GetServer().GetActivitCode()
	team := client.GetTeam()
	activitCodeMgr.Use(self.ActivitCodeInfo, team.AccountID, team.ID)
	///打上已领取的掩码
	activeCode := activitCodeMgr.GetActiveCode(self.ActivitCodeInfo)
	activeCodeAwardType := activitCodeMgr.GetActiveCodeAwardType(activeCode.Type)
	team.ActivitCodeAward = SetMask64(team.ActivitCodeAward, activeCodeAwardType.ID, 1)
	return true
}

func (self *ActivitCodeMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	activitcodeMgr := GetServer().GetActivitCode()
	activeCode := activitcodeMgr.GetActiveCode(self.ActivitCodeInfo)
	activeCodeAwardType := activitcodeMgr.GetActiveCodeAwardType(activeCode.Type)
	if activeCodeAwardType.Awarditem1 > 0 {
		team.AwardObject(activeCodeAwardType.Awarditem1,
			activeCodeAwardType.Awardcount1, activeCodeAwardType.Awardgrade1, 0)
	}
	if activeCodeAwardType.Awarditem2 > 0 {
		team.AwardObject(activeCodeAwardType.Awarditem2,
			activeCodeAwardType.Awardcount2, activeCodeAwardType.Awardgrade2, 0)
	}
	if activeCodeAwardType.Awarditem3 > 0 {
		team.AwardObject(activeCodeAwardType.Awarditem3,
			activeCodeAwardType.Awardcount3, activeCodeAwardType.Awardgrade3, 0)
	}
	if activeCodeAwardType.Awarditem4 > 0 {
		team.AwardObject(activeCodeAwardType.Awarditem4,
			activeCodeAwardType.Awardcount4, activeCodeAwardType.Awardgrade4, 0)
	}
	if activeCodeAwardType.Awarditem5 > 0 {
		team.AwardObject(activeCodeAwardType.Awarditem5,
			activeCodeAwardType.Awardcount5, activeCodeAwardType.Awardgrade5, 0)
	}
	if activeCodeAwardType.Awarditem6 > 0 {
		team.AwardObject(activeCodeAwardType.Awarditem6,
			activeCodeAwardType.Awardcount6, activeCodeAwardType.Awardgrade6, 0)
	}
	return true
}

func (self *ActivitCodeMsg) processAction(client IClient) (result bool) {
	activitCodeResultMsg := new(ActivitCodeResultMsg)
	activitCodeResultMsg.ReturnMsg = "ok"
	defer func() {
		if false == result {
			activitCodeResultMsg.ReturnMsg = "fail"
		}
		client.SendMsg(activitCodeResultMsg)
	}()
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
