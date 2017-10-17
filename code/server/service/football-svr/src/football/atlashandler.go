package football

import ()

type QueryAtlasInfoMsg struct { //! 查询图鉴信息
	MsgHead `json:"head"` //! "atlas", "queryatlasinfo"
	Type    int           `json:"type"` //! 获取类型: 0为全部, 1为所有A级, 2为所有S级, 3为所有SS级
}

func (self *QueryAtlasInfoMsg) GetTypeAndAction() (string, string) {
	return "atlas", "queryatlasinfo"
}

func (self *QueryAtlasInfoMsg) processAction(client IClient) bool {
	team := client.GetTeam()
	atlasMgr := team.GetAtlasMgr()

	msg := new(QueryAtlasInfoMsgResult)

	if self.Type == 0 {
		msg.Atlas = atlasMgr.GetAllAtlas()
	} else {
		msg.Atlas = atlasMgr.GetTypeAtlas(self.Type)
	}

	client.SendMsg(msg)
	return true
}

type QueryAtlasInfoMsgResult struct { //! 查询图鉴消息返回结果
	MsgHead `json:"head"` //! "atlas", "queryatlasinforesult"
	Atlas   []AtlasInfo   `json:"atlas"`
}

func (self *QueryAtlasInfoMsgResult) GetTypeAndAction() (string, string) {
	return "atlas", "queryatlasinforesult"
}

type ReceiveAtlasAwardMsg struct { //! 领取图鉴奖励消息
	MsgHead  `json:"head"` //! "atlas", "receiveaward"
	StarType int           `json:"startype"` //! 领取奖励球星
}

func (self *ReceiveAtlasAwardMsg) GetTypeAndAction() (string, string) {
	return "atlas", "receiveaward"
}

func (self *ReceiveAtlasAwardMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	atlasMgr := team.GetAtlasMgr()
	loger := GetServer().GetLoger()

	hasStar := atlasMgr.HasStar(self.StarType)
	if loger.CheckFail("hasStar == true", hasStar == true, hasStar, true) {
		return false //! 必须曾经得到过该球星
	}

	atlasInfo := atlasMgr.GetAtlas(self.StarType)
	isReceived := atlasInfo.Received
	if loger.CheckFail("isReceived == 0", isReceived == 0, isReceived, 0) {
		return false //! 必须未领过该球星奖励
	}

	return true
}

func (self *ReceiveAtlasAwardMsg) payAction(client IClient) bool {
	team := client.GetTeam()
	atlasMgr := team.GetAtlasMgr()
	atlasInfo := atlasMgr.GetAtlas(self.StarType)
	atlasInfo.Received = 1
	return true
}

func (self *ReceiveAtlasAwardMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	atlasMgr := team.GetAtlasMgr()
	staticDataMgr := GetServer().GetStaticDataMgr()

	starType := staticDataMgr.GetStarType(self.StarType)

	if starType.Class < 450 { //! 暂定奖励
		team.AwardObject(400006, 50, 0, 0)
	} else if starType.Class < 500 {
		team.AwardObject(400006, 100, 0, 0)
	} else if starType.Class >= 500 {
		team.AwardObject(400006, 200, 0, 0)
	}

	msg := new(QueryAtlasInfoMsgResult)

	msg.Atlas = atlasMgr.GetAllAtlas()
	client.SendMsg(msg)
	return true
}

func (self *ReceiveAtlasAwardMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false {
		return false
	}
	if self.payAction(client) == false {
		return false
	}
	if self.doAction(client) == false {
		return false
	}

	return true
}
