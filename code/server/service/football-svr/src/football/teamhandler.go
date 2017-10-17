package football

import (
	"fmt"
	//	"math"
	"strconv"
	"strings"
	"time"
)

const (
	GetOperationString    = 1
	SetOperationString    = 2
	OperationStringMaxLen = 255
)

type TeamInfoMsg struct {
	MsgHead       `json:"head"`
	TeamInfo      `json:"teaminfo"`    ///球队信息
	StarSpyInfo   `json:"starspyinfo"` ///球探信息
	StarList      []StarInfo           `json:"starlist"`      ///球员列表
	FormationList []FormationInfo      `json:"formationlist"` ///阵型信息
	EquipmentList ItemInfoList         `json:"equipmentlist"` ///首发球员装备道具信息
	AccessToken   string               `json:"accesstoken"`   ///sessionkey
	SDKUserID     string               `json:"sdkuserid"`     ///sdkuserid
}

func (self *TeamInfoMsg) GetTypeAndAction() (string, string) {
	return "team", "teaminfo"
}

func NewTeamInfoMsg() *TeamInfoMsg {
	msg := new(TeamInfoMsg)
	msg.MsgType = "team"
	msg.Action = "teaminfo"
	msg.CreateTime = int(time.Now().Unix())
	return msg
}

type TeamHandler struct {
	MsgHandler
}

func (self *TeamHandler) getName() string { ///返回可处理的消息类型
	return "team"
}

func (self *TeamHandler) initHandler() { ///初始化消息处理器
	self.addActionToList(new(TeamCreateMsg))
	self.addActionToList(new(GetStarExpPoolMsg))
	//self.addActionToList(new(TeamAddTrainCellMsg))
	//self.addActionToList(new(TeamQueryStarTrainListMsg))
	//self.addActionToList(new(TeamTrainStarMsg))
	//self.addActionToList(new(TeamAbortTrainStarMsg))
}

//type TeamAddTrainCellMsg struct { ///请求增加球队训练位
//	MsgHead `json:"head"`
//}

//func (self *TeamAddTrainCellMsg) GetTypeAndAction() (string, string) {
//	return "team", "addtraincell"
//}

//type TeamQueryStarTrainListMsg struct { ///查询球队球员训练位列表请求
//	MsgHead `json:"head"`
//}

//type TeamQueryStarTrainListResultMsg struct { ///球队球员训练位列表请求结果
//	MsgHead       `json:"head"`
//	StarTrainList ProcessInfoList `json:"startrainlist"`
//	//GainExpList   IntList         `json:"gainexplist"` ///已获得经验列表,和StarTrainList顺序对应
//}

//func (self *TeamQueryStarTrainListResultMsg) GetTypeAndAction() (string, string) {
//	return "team", "teamquerystartrainlistresult"
//}

//func NewTeamQueryStarTrainListResultMsg(starTrainList ProcessInfoList) *TeamQueryStarTrainListResultMsg {
//	msg := new(TeamQueryStarTrainListResultMsg)
//	msg.StarTrainList = starTrainList
//	return msg
//}

//func (self *TeamQueryStarTrainListMsg) GetTypeAndAction() (string, string) {
//	return "team", "teamquerystartrainlist"
//}

//func (self *TeamQueryStarTrainListMsg) processAction(client IClient) bool {
//	if client.GetTeam() == nil {
//		return false
//	}
//	processMgr := client.GetTeam().GetProcessMgr()
//	porcessList := processMgr.GetProcessInfoList(ProcessTypeStarTrain, getProcessListAll)
//	msg := NewTeamQueryStarTrainListResultMsg(porcessList)
//	client.SendMsg(msg)
//	return true
//}

//type TeamAbortTrainStarResultMsg struct { ///中止球队训练一位球员消息,由客户端发起
//	MsgHead   `json:"head"`
//	Result    string `json:"result"`    ///ok or fail
//	ProcessID int    `json:"processid"` ///处理对象id
//}

//func (self *TeamAbortTrainStarResultMsg) GetTypeAndAction() (string, string) {
//	return "team", "aborttrainstarresult"
//}

//type TeamAbortTrainStarMsg struct { ///中止球队训练一位球员消息,由客户端发起
//	MsgHead   `json:"head"`
//	ProcessID int `json:"processid"` ///处理对象id
//}

//func (self *TeamAbortTrainStarMsg) SendResultMsg(client IClient, result string, processID int) {
//	teamAbortTrainStarResultMsg := new(TeamAbortTrainStarResultMsg)
//	teamAbortTrainStarResultMsg.Result = result
//	teamAbortTrainStarResultMsg.ProcessID = processID
//	client.SendMsg(teamAbortTrainStarResultMsg)
//}

//func (self *TeamAbortTrainStarMsg) processAction(client IClient) bool {
//	processMgr := client.GetTeam().GetProcessMgr()
//	process := processMgr.GetProcess(ProcessTypeStarTrain, self.ProcessID)
//	process.Reset() ///重置状态,清空训练信息
//	self.SendResultMsg(client, msgResultOK, self.ProcessID)
//	return true
//}

//func (self *TeamAbortTrainStarMsg) GetTypeAndAction() (string, string) {
//	return "team", "aborttrainstar"
//}

//type TeamTrainStarMsg struct { ///球队请求训练一位球员消息
//	MsgHead   `json:"head"`
//	StarID    int  `json:"starid"`    ///需要训练的球员id
//	ProcessID int  `json:"processid"` ///处理对象id
//	Immediate bool `json:"immediate"` ///是否立即完成,默认为false
//}

//func (self *TeamTrainStarMsg) GetTypeAndAction() (string, string) {
//	return "team", "teamtrainstar"
//}

//func (self *TeamTrainStarMsg) checkAction(client IClient) bool {
//	processMgr := client.GetTeam().GetProcessMgr()
//	team := client.GetTeam()
//	loger := GetServer().GetLoger()
//	staticDataMgr := GetServer().GetStaticDataMgr()
//	configStarTrain := staticDataMgr.GetConfigStaticData(configTeam, configItemStarTrain)
//	payTicket, _ := strconv.Atoi(configStarTrain.Param6) ///得到立即完成所需球票数
//	currentTicket := client.GetTeam().GetTicket()
//	if loger.CheckFail("configStarTrain!=nil", configStarTrain != nil, configStarTrain, nil) {
//		return false
//	}
//	star := team.GetStar(self.StarID) ///判断此球员有效性
//	if loger.CheckFail("star!=nil", star != nil, star, nil) {
//		return false ///球员不存在
//	}
//	if loger.CheckFail("star.isMaxLevel() == false", star.isMaxLevel() == false, star.isMaxLevel(), false) {
//		return false ///球员满级后不准训练
//	}
//	processID := processMgr.FindProcessByObjID(ProcessTypeStarTrain, self.StarID)
//	if processID > 0 {
//		if loger.CheckFail("processID==self.ProcessID", processID == self.ProcessID, processID, self.ProcessID) {
//			return false ///尝试将一个球员连续放到不同的训练位中
//		}
//	}
//	process := processMgr.GetProcess(ProcessTypeStarTrain, self.ProcessID)
//	if loger.CheckFail("process!=nil", process != nil, process, nil) {
//		return false ///训练位不存在
//	}
//	processInfo := process.GetInfo() ///得到训练位信息
//	if processInfo.ObjID > 0 {       ///判断训练位有人的情况,只能为相同的球员
//		if loger.CheckFail("processInfo.ObjID==self.StarID", processInfo.ObjID == self.StarID,
//			processInfo.ObjID, self.StarID) {
//			return false ///训练位被其它球员占用
//		}
//		if false == self.Immediate {
//			now := int(time.Now().Unix())
//			if loger.CheckFail("now>=processInfo.ExpireTime", now >= processInfo.ExpireTime,
//				now, processInfo.ExpireTime) {
//				return false ///重复训练时必须已完成上次训练
//			}
//		}
//	}
//	if true == self.Immediate { ///客户端希望立即完成
//		if loger.CheckFail("payTicket>0", payTicket > 0, payTicket, 0) {
//			return false ///配置数据必须存在
//		}
//		if loger.CheckFail("currentTicket>=payTicket", currentTicket >= payTicket,
//			currentTicket, payTicket) {
//			return false ///当要求立即完成训练时球队球票余额需要够扣
//		}
//	}
//	return true
//}

//func (self *TeamTrainStarMsg) payAction(client IClient) bool {
//	if false == self.Immediate {
//		return true
//	}
//	team := client.GetTeam()
//	staticDataMgr := GetServer().GetStaticDataMgr()
//	syncMgr := client.GetSyncMgr()
//	configStarTrain := staticDataMgr.GetConfigStaticData(configTeam, configItemStarTrain)
//	payTicket, _ := strconv.Atoi(configStarTrain.Param6) ///得到立即完成所需球票数
//	team.PayTicket(payTicket)
//	syncMgr.SyncObject("TeamTrainStarMsg", team)
//	return true
//}

//func (self *TeamTrainStarMsg) doAction(client IClient) bool {
//	team := client.GetTeam()
//	processMgr := client.GetTeam().GetProcessMgr()
//	syncMgr := client.GetSyncMgr()
//	staticDataMgr := GetServer().GetStaticDataMgr()
//	configStarTrain := staticDataMgr.GetConfigStaticData(configTeam, configItemStarTrain)
//	expireTime, _ := strconv.Atoi(configStarTrain.Param5) ///得到配置数据中每次训练周期时间
//	expireTime = expireTime + int(time.Now().Unix())      ///得到过期时间秒数
//	process := processMgr.GetProcess(ProcessTypeStarTrain, self.ProcessID)
//	processInfo := process.GetInfo()
//	if true == self.Immediate {
//		expireTime = int(time.Now().Unix()) ///得到过期时间秒数
//		awardStarExp := process.CalcTrainAwardExp(team.GetLevel(), true)
//		star := team.GetStar(self.StarID)
//		star.AwardExp(awardStarExp)
//		processInfo.Param1 += awardStarExp ///放入累计经验
//		syncMgr.SyncObject(systemTypeStarTrain, star)
//	} else {
//		process.Reset() ///非立即完成时重置状态
//	}
//	processInfo.ObjID = self.StarID
//	processInfo.ExpireTime = expireTime
//	processInfo.NextProcessTime = Now() + trainAwardExpInterval ///更新下次给经验的时间
//	processInfoList := processMgr.GetProcessInfoList(ProcessTypeStarTrain, self.ProcessID)
//	teamQueryStarTrainListResultMsg := NewTeamQueryStarTrainListResultMsg(processInfoList)
//	client.SendMsg(teamQueryStarTrainListResultMsg) ///发送最新的
//	return true
//}

//func (self *TeamTrainStarMsg) processAction(client IClient) (result bool) {
//	defer func() {
//		if false == result {
//			self.SendResultMsg(client, msgResultFail)
//		} else {
//			self.SendResultMsg(client, msgResultOK)
//		}
//	}()
//	if self.checkAction(client) == false { ///检测
//		return false
//	}
//	if self.payAction(client) == false { ///支付
//		return false
//	}
//	if self.doAction(client) == false { ///发货
//		return false
//	}
//	return true
//}

//func (self *TeamTrainStarMsg) SendResultMsg(client IClient, result string) {
//	msgResult := NewTeamTrainStarResultMsg(result)
//	client.SendMsg(msgResult)
//}

//type TeamTrainStarResultMsg struct { ///训练位结果
//	MsgHead `json:"head"`
//	Result  string `json:"result"` ///ok成功 fail失败
//}

//func (self *TeamTrainStarResultMsg) GetTypeAndAction() (string, string) {
//	return "team", "teamtrainstarresult"
//}

//func NewTeamTrainStarResultMsg(result string) *TeamTrainStarResultMsg {
//	msg := new(TeamTrainStarResultMsg)
//	msg.Result = result
//	return msg
//}

//type TeamAddTrainCellResultMsg struct { ///增加球队训练位结果
//	MsgHead `json:"head"`
//	Result  string `json:"result"` ///ok成功 fail失败
//}

//func (self *TeamAddTrainCellResultMsg) GetTypeAndAction() (string, string) {
//	return "team", "addtraincellresult"
//}

//func NewTeamAddTrainCellResultMsg(result string) *TeamAddTrainCellResultMsg {
//	msg := new(TeamAddTrainCellResultMsg)
//	msg.Result = result
//	return msg
//}

//func (self *TeamAddTrainCellMsg) SendResultMsg(client IClient, result string) {
//	msgResult := NewTeamAddTrainCellResultMsg(result)
//	client.SendMsg(msgResult)
//}

//func (self *TeamAddTrainCellMsg) checkAction(client IClient) bool {
//	processMgr := client.GetTeam().GetProcessMgr()
//	team := client.GetTeam()
//	staticDataMgr := GetServer().GetStaticDataMgr()
//	configStarTrain := staticDataMgr.GetConfigStaticData(configTeam, configItemStarTrain)
//	if nil == configStarTrain {
//		return false ///没找到配置信息
//	}
//	defaultTrainCellCount, _ := strconv.Atoi(configStarTrain.Param1)        ///得到默认训练位个数
//	currentTrainCell := processMgr.GetProcessMaxIndex(ProcessTypeStarTrain) ///得到已拥有训练位总数
//	currentTicket := team.GetTicket()                                       ///得到当前球票数
//	maxTrainCellCount, _ := strconv.Atoi(configStarTrain.Param2)            ///最大训练格数
//	perCellPrice, _ := strconv.Atoi(configStarTrain.Param3)
//	maxPerCellPrice, _ := strconv.Atoi(configStarTrain.Param4)
//	needPay := (currentTrainCell - defaultTrainCellCount + 1) * perCellPrice ///得到球票总花费
//	if needPay > maxPerCellPrice {
//		needPay = maxPerCellPrice ///花费如果超过上限则设置为上限
//	}
//	if needPay <= 0 || currentTicket < needPay {
//		return false ///余额不足
//	}
//	if currentTrainCell > maxTrainCellCount {
//		return false ///训练位已满无法继续开格
//	}
//	return true
//}

//func (self *TeamAddTrainCellMsg) payAction(client IClient) bool {
//	processMgr := client.GetTeam().GetProcessMgr()
//	team := client.GetTeam()
//	staticDataMgr := GetServer().GetStaticDataMgr()
//	currentTrainCell := processMgr.GetProcessMaxIndex(ProcessTypeStarTrain) ///得到已拥有训练位总数
//	configStarTrain := staticDataMgr.GetConfigStaticData(configTeam, configItemStarTrain)
//	defaultTrainCellCount, _ := strconv.Atoi(configStarTrain.Param1) ///得到默认训练位个数
//	perCellPrice, _ := strconv.Atoi(configStarTrain.Param3)
//	maxPerCellPrice, _ := strconv.Atoi(configStarTrain.Param4)
//	needPay := (currentTrainCell - defaultTrainCellCount + 1) * perCellPrice ///得到球票总花费
//	if needPay > maxPerCellPrice {
//		needPay = maxPerCellPrice ///花费如果超过上限则设置为上限
//	}
//	currentTicket := team.PayTicket(needPay) ///扣球票
//	actionAttribChangeMsg := NewActionAttribChangeMsg(systemTypeStarTrain, teamTicketAttribType, 0, currentTicket)
//	client.SendMsg(actionAttribChangeMsg)
//	return true
//}

//func (self *TeamAddTrainCellMsg) doAction(client IClient) bool {
//	processMgr := client.GetTeam().GetProcessMgr()
//	processID := processMgr.AddProcess(ProcessTypeStarTrain, 1) ///添加一个训练位
//	if processID <= 0 {
//		return false
//	}
//	processInfoList := processMgr.GetProcessInfoList(ProcessTypeStarTrain, processID)
//	msg := NewTeamQueryStarTrainListResultMsg(processInfoList)
//	client.SendMsg(msg)
//	self.SendResultMsg(client, msgResultOK)
//	return true
//}

//func (self *TeamAddTrainCellMsg) handleAction(client IClient) bool {
//	if self.checkAction(client) == false { ///检测
//		return false
//	}
//	if self.payAction(client) == false { ///支付
//		return false
//	}
//	if self.doAction(client) == false { ///发货
//		return false
//	}
//	return true
//}

//func (self *TeamAddTrainCellMsg) processAction(client IClient) bool {
//	ok := self.handleAction(client)
//	if ok == false {
//		self.SendResultMsg(client, msgResultFail)
//	}
//	return ok
//}

type TeamCreateMsg struct { ///请求创建球队消息
	MsgHead      `json:"head"`
	AccountID    int     `json:"accountid"`    ///玩家帐号id
	TeamName     string  `json:"teamname"`     ///球队名
	Icon         int     `json:"icon"`         ///球队队徽
	TeamShirts   int     `json:"teamshirts"`   ///球队球衣
	StarTypeList IntList `json:"startypelist"` ///默认初始十一人大名单
}

func (self *TeamCreateMsg) GetTypeAndAction() (string, string) {
	return "team", "createteam"
}

func (self *TeamCreateMsg) IsNeedTeamHandle() bool { ///此消息不需要team创建
	return false
}

func (self *TeamCreateMsg) New() interface{} {
	return new(TeamCreateMsg)
}

func (self *TeamCreateMsg) CheckClientAccount(client IClient) bool { ///检测客户端账号合法性
	//if self.AccountID != client.accountID {
	//	GetServer().GetLoger().Warn("CreateNewTeam fail! self.AccountID != client.accountID!")
	//	client.SendErrorMsg(failCreateTeam, failInvalidAccountID)
	//	return false
	//}
	return true
}

func (self *TeamCreateMsg) CheckStarTypeList() bool { ///检测初始十一人名单的合法性
	loger := GetServer().GetLoger()
	starTypeListLen := self.StarTypeList.Len()
	loger.Info("starTypeListLen:%d", starTypeListLen)
	if loger.CheckFail("starTypeListLen==formationMinStarCount", starTypeListLen == formationMinStarCount,
		starTypeListLen, formationMinStarCount) { ///初始名单次数不等于11人
		return false
	}
	groupDic := IntList{drawGroupTypeDefaultMasterStar, drawGroupTypeDefaultStarCF,
		drawGroupTypeDefaultStarRMF, drawGroupTypeDefaultStarCMF,
		drawGroupTypeDefaultStarCMF, drawGroupTypeDefaultStarLMF,
		drawGroupTypeDefaultStarRB, drawGroupTypeDefaultStarCB,
		drawGroupTypeDefaultStarCB, drawGroupTypeDefaultStarLB,
		drawGroupTypeDefaultStarGK}
	staticDataMgr := GetServer().GetStaticDataMgr()
	drawItemGroup := IntList{}
	foundStarType := false ///已找到对应球员类型
	groupType := 0
	for i := range self.StarTypeList {
		foundStarType = false ///重置寻找标识
		starType := self.StarTypeList[i]
		groupType = groupDic[i]
		drawItemGroup = staticDataMgr.GetDrawGroupIndexList(groupType)
		for k := range drawItemGroup {
			drawGroupItemID := drawItemGroup[k]
			drawGroupItem := staticDataMgr.GetDrawGroupStaticData(drawGroupItemID)
			drawStarType := drawGroupItem.AwardType
			if drawStarType == starType {
				foundStarType = true ///找到对应球员
				break
			}
		}
		///初始名单中的球员不在对应抽选表中
		if loger.CheckFail("create team starType not in draw group", foundStarType == true,
			starType, drawGroupTypeDefaultStarCF-i+1) {
			return false
		}
	}
	return true
}

func (self *TeamCreateMsg) CheckAction(client IClient) bool { ///检测消息合法性
	loger := GetServer().GetLoger()
	isGM := strings.Contains(self.TeamName, "[GM]")
	if loger.CheckFail("isGM==false", isGM == false, isGM, false) {
		return false ///不允许尝试创建gm权限的号
	}
	if self.CheckClientAccount(client) == false { ///检测客户端账号合法性
		client.SendErrorMsg(failCreateTeam, failInvalidAccountID)
		return false
	}
	if client.HasInitTeam() == true { ///检测客户端是是否已经创建过球队了
		GetServer().GetLoger().Warn("CreateNewTeam fail! client's team is not nil!")
		client.SendErrorMsg(failCreateTeam, failInvalidMsg)
		return false
	}
	checkStarTypeList := self.CheckStarTypeList()
	if loger.CheckFail("checkStarTypeList==true", checkStarTypeList == true, checkStarTypeList, true) {
		return false
	}
	return true
}

func (self *TeamCreateMsg) getCreateTeamDefaultFormationType() int { ///得到创建球队时默认的阵形type
	configStaticData := GetServer().GetStaticDataMgr().GetConfigStaticData(configTeam, configItemDefaultTeamParam)
	if nil == configStaticData {
		return 0 ///返回空列表
	}
	defaultFormationType, _ := strconv.Atoi(configStaticData.Param2) ///从配置数据中得到默认阵形type
	return defaultFormationType
}

func (self *TeamCreateMsg) getCreateTeamDefaultSkillList() []int { ///得到创建球队时默认的技能列表
	result := []int{}
	configStaticData := GetServer().GetStaticDataMgr().GetConfigStaticData(configTeam, configItemDefaultTeamParam)
	if nil == configStaticData {
		return []int{} ///返回空列表
	}
	skillStringList := strings.Split(configStaticData.Param5, ",")
	for i := range skillStringList {
		skillType, _ := strconv.Atoi(skillStringList[i])
		result = append(result, skillType)
	}
	return result
}

func (self *TeamCreateMsg) getCreateTeamDefaultItemList() []int { ///得到创建球队时默认的球员列表
	result := []int{}
	configStaticData := GetServer().GetStaticDataMgr().GetConfigStaticData(configTeam, configItemDefaultTeamParam)
	if nil == configStaticData {
		return []int{} ///返回空列表
	}
	itemStringList := strings.Split(configStaticData.Param4, ",")
	for i := range itemStringList {
		itemType, _ := strconv.Atoi(itemStringList[i])
		if itemType <= 0 {
			continue
		}
		result = append(result, itemType)
	}
	return result
}

func (self *TeamCreateMsg) getCreateTeamDefaultStarList() []int { ///得到创建球队时默认的球员列表
	result := []int{}
	configStaticData := GetServer().GetStaticDataMgr().GetConfigStaticData(configTeam, configItemDefaultTeamParam)
	if nil == configStaticData {
		return []int{} ///返回空列表
	}
	startStringList := strings.Split(configStaticData.Param1, ",")
	for i := range startStringList {
		starType, _ := strconv.Atoi(startStringList[i])
		if starType <= 0 {
			continue
		}
		result = append(result, starType)
	}
	return result
}

func (self *TeamCreateMsg) CreateDefaultStarTrainProcessCenter(client IClient, teamID int) bool { ///创建默认的球员训练处理中心
	staticDataMgr := GetServer().GetStaticDataMgr()
	configStarTrain := staticDataMgr.GetConfigStaticData(configTeam, configItemStarTrain)
	if nil == configStarTrain {
		return false ///没找到配置信息
	}
	insertDefaultStarTrainProcessQuery := fmt.Sprintf("Insert %s (teamid,type,pos) VALUES ", tableProcessCenter)
	///得到默认训练位数
	defaultTrainCellCount, _ := strconv.Atoi(configStarTrain.Param1)
	for i := 1; i <= defaultTrainCellCount; i++ {
		if defaultTrainCellCount == i {
			///最后一条记录特殊处理
			insertDefaultStarTrainProcessQuery = insertDefaultStarTrainProcessQuery + fmt.Sprintf("(%d,%d,%d)",
				teamID, ProcessTypeStarTrain, i)
		} else {
			insertDefaultStarTrainProcessQuery = insertDefaultStarTrainProcessQuery + fmt.Sprintf("(%d,%d,%d),",
				teamID, ProcessTypeStarTrain, i)
		}
	}
	///执行插入语句
	_, rowsProcessAffected := GetServer().GetDynamicDB().Exec(insertDefaultStarTrainProcessQuery)
	if rowsProcessAffected != defaultTrainCellCount {
		GetServer().GetLoger().Warn("TeamCreateMsg CreateDefaultStarTrainProcessCenter insertDefaultStarTrainProcessQuery fail! msg:%v", self)
		return false
	}
	return true
}

func (self *TeamCreateMsg) createNewTeam(client IClient) bool { ///创建新球员逻辑
	///测试代码: 不允许创建球队名时增加GM标志
	nameLegal := strings.Index(self.TeamName, "_GM")
	if nameLegal != -1 {
		GetServer().GetLoger().Warn("Team name has GM sign  msg:%v", self)
		client.SendErrorMsg(failCreateTeam, failInvalidName)
		return false
	}

	if len(self.TeamName) <= 0 {
		GetServer().GetLoger().Warn("Team name is nil  msg:%v", self)
		client.SendErrorMsg(failCreateTeam, failInvalidName)
		return false
	}

	///创建新球队
	createTeamQuery := fmt.Sprintf("insert %s set name='%s',accountid=%d,icon=%d,teamshirts=%d,maketime=%d",
		tableTeam, self.TeamName, self.AccountID, self.Icon, self.TeamShirts, Now())
	lastInsertTeamID, _ := GetServer().GetDynamicDB().Exec(createTeamQuery) ///创建新帐号
	if lastInsertTeamID <= 0 {
		GetServer().GetLoger().Warn("TeamCreateMsg processAction Insert New Team fail! msg:%v", self)
		client.SendErrorMsg(failCreateTeam, failSameName)
		return false
	}

	///创建新球队所附属球员,默认球员为
	defaultStarTypeList := self.StarTypeList //self.getCreateTeamDefaultStarList() ///从配置表中读取默认球员列表
	defaultStarTypeListLen := len(defaultStarTypeList)
	if defaultStarTypeListLen != 11 {
		GetServer().GetLoger().Warn("TeamCreateMsg processAction getCreateTeamDefaultStarList len !=11! msg:%v", self)
		return false
	}

	createDefaultStarQuery := fmt.Sprintf("insert %s (teamid,type,grade,evolvecount) VALUES ", tableStar)
	for i := range defaultStarTypeList {
		if 10 == i {
			createDefaultStarQuery = createDefaultStarQuery + fmt.Sprintf("(%d,%d,%d,%d)", lastInsertTeamID, defaultStarTypeList[i], 1, 1) ///最后一条记录特殊处理
			//else if 0 == i { 删除初始球星星级
			//}createDefaultStarQuery = createDefaultStarQuery + fmt.Sprintf("(%d,%d,%d,%d),", lastInsertTeamID, defaultStarTypeList[i], 2, 4) ///当家球星默认为蓝色(4星级)
		} else {
			createDefaultStarQuery = createDefaultStarQuery + fmt.Sprintf("(%d,%d,%d,%d),", lastInsertTeamID, defaultStarTypeList[i], 1, 1)
		}
	}
	_, rowsStarAffected := GetServer().GetDynamicDB().Exec(createDefaultStarQuery)
	if rowsStarAffected != 11 {
		GetServer().GetLoger().Warn("TeamCreateMsg processAction createDefaultStarQuery fail! msg:%v", self)
		return false
	}

	///赠送球员
	defaultAddStarTypeList := self.getCreateTeamDefaultStarList() ///从配置表中读取默认球员列表
	defaultAddStarTypeListLen := len(defaultAddStarTypeList)
	if defaultAddStarTypeListLen > 0 {
		createDefaultAddStarQuery := fmt.Sprintf("insert %s (teamid,type) VALUES ", tableStar)
		for i := range defaultAddStarTypeList {
			createDefaultAddStarQuery = createDefaultAddStarQuery + fmt.Sprintf("(%d,%d)",
				lastInsertTeamID, defaultAddStarTypeList[i])
			if i < defaultAddStarTypeListLen-1 {
				createDefaultAddStarQuery += ","
			}
		}
		_, rowsAddStarAffected := GetServer().GetDynamicDB().Exec(createDefaultAddStarQuery)
		if rowsAddStarAffected <= 0 {
			GetServer().GetLoger().Warn("TeamCreateMsg processAction createDefaultAddStarQuery fail! msg:%v", self)
			return false
		}
	}
	///创建新球队所附属的阵型
	defaultFormationType := self.getCreateTeamDefaultFormationType()
	if defaultFormationType <= 0 {
		GetServer().GetLoger().Warn("TeamCreateMsg processAction getCreateTeamDefaultFormationType is nil! msg:%v", self)
		return false
	}
	defaultStarListQuery := fmt.Sprintf("select id from %s where teamid=%d order by id asc limit 11", tableStar, lastInsertTeamID) ///取得默认球员列表
	rowsDefaultStar := GetServer().GetDynamicDB().Query(defaultStarListQuery)
	if nil == rowsDefaultStar { ///这里不可能为nil记录
		GetServer().GetLoger().Warn("TeamCreateMsg processAction defaultStarListQuery rows is nil! msg:%v", self)
		return false
	}
	createDefaultFormation := fmt.Sprintf("insert %s VALUES (0,%d,%d", tableFormation, lastInsertTeamID, defaultFormationType) ///默认第一个阵型为新建球队的初始阵型
	starID := 0
	for rowsDefaultStar.Next() {
		rowsDefaultStar.Scan(&starID)
		createDefaultFormation = createDefaultFormation + fmt.Sprintf(",%d", starID)
	}
	rowsDefaultStar.Close()
	createDefaultFormation = createDefaultFormation + ",101,101,0,0)" ///进行sql闭合
	lastInsertFormationID, _ := GetServer().GetDynamicDB().Exec(createDefaultFormation)
	if lastInsertFormationID <= 0 {
		GetServer().GetLoger().Warn("TeamCreateMsg processAction createDefaultFormation fail! msg:%v", self)
		return false
	}

	///更新球队当前阵型字段数据
	updateTeamCurrentFormationQuery := fmt.Sprintf("update dy_team set formationid=%d where id=%d", lastInsertFormationID, lastInsertTeamID)
	_, rowsTeamAffected := GetServer().GetDynamicDB().Exec(updateTeamCurrentFormationQuery) ///创建新帐号
	if rowsTeamAffected <= 0 {
		GetServer().GetLoger().Warn("TeamCreateMsg processAction updateTeamCurrentFormationQuery fail! msg:%v", self)
		return false
	}

	///创建新球队所属球探记录
	createDefaultStarSpyQuery := fmt.Sprintf("insert %s set teamid=%d, discoverluck1 = 95, discoverluck2 = 95, discoverluck3 = 95", tableStarSpy, lastInsertTeamID)
	_, rowsStarSpyAffected := GetServer().GetDynamicDB().Exec(createDefaultStarSpyQuery)
	if rowsStarSpyAffected <= 0 {
		GetServer().GetLoger().Warn("TeamCreateMsg processAction createDefaultStarSpyQuery fail! msg:%v", self)
		return false
	}

	///创建新球队所属球员训练中心
	if self.CreateDefaultStarTrainProcessCenter(client, lastInsertTeamID) == false {
		GetServer().GetLoger().Warn("TeamCreateMsg processAction CreateDefaultStarTrainProcessCenter fail! msg:%v", self)
		return false
	}

	///创建新球队时赠送道具列表,调试功能
	defaultItemTypeList := self.getCreateTeamDefaultItemList() ///从配置表中读取默认球员列表
	defaultItemTypeListLen := len(defaultItemTypeList)
	if defaultItemTypeListLen > 0 {
		createDefaultItemQuery := fmt.Sprintf("insert %s (teamid,type,color) VALUES ", tableItem)
		for i := range defaultItemTypeList {
			createDefaultItemQuery = createDefaultItemQuery + fmt.Sprintf("(%d,%d,%d)", lastInsertTeamID,
				defaultItemTypeList[i], Random(1, 4)) ///最后一条记录特殊处理
			if i < defaultItemTypeListLen-1 {
				createDefaultItemQuery += ","
			}
		}
		_, rowsItemAffected := GetServer().GetDynamicDB().Exec(createDefaultItemQuery)
		if rowsItemAffected <= 0 {
			GetServer().GetLoger().Warn("TeamCreateMsg processAction createDefaultItemQuery fail! msg:%v", self)
			return false
		}
	}

	///创建新球队时赠送技能列表,调试功能
	//defaultSkillTypeList := self.getCreateTeamDefaultSkillList() ///从配置表中读取默认球员列表
	//defaultSkillTypeListLen := len(defaultSkillTypeList)
	//createDefaultSkillQuery := fmt.Sprintf("insert %s (teamid,type) VALUES ", tableSkill)
	//for i := range defaultSkillTypeList {
	//	createDefaultSkillQuery += fmt.Sprintf("(%d,%d)", lastInsertTeamID, defaultSkillTypeList[i]) ///最后一条记录特殊处理
	//	if i < defaultSkillTypeListLen-1 {
	//		createDefaultSkillQuery += ","
	//	}
	//}
	//_, rowsSkillAffected := GetServer().GetDynamicDB().Exec(createDefaultSkillQuery)
	//if rowsSkillAffected <= 0 {
	//	GetServer().GetLoger().Warn("TeamCreateMsg processAction createDefaultSkillQuery fail! msg:%v", self)
	//	return false
	//}

	//! cy
	//! 领取记录
	strSql := fmt.Sprintf("insert into `dy_getpower`(`teamid`, `time1`, `time2`) values (%d, '1990-01-01 00:00:00', '1990-01-01 00:00:00')", lastInsertTeamID)
	_, rows := GetServer().GetDynamicDB().Exec(strSql)
	if rows <= 0 {
		GetServer().GetLoger().Warn("TeamCreateMsg processAction strSql fail! msg:%s", strSql)
		return false
	}

	return true
}

func (self *TeamCreateMsg) processAction(client IClient) bool {
	if self.CheckAction(client) == false { ///检测客户端请求合法性
		return false
	}
	if self.createNewTeam(client) == false {
		return false
	}
	if client.LoadTeam(self.AccountID) != true { ///加载球队信息,一般不会失败的
		//GetServer().GetLoger().Warn("TeamCreateMsg processAction client.CreateTeam fail! msg:%v", self)
		return false
	}
	client.CreateTeamRecord() ///添加创建球队记录

	team := client.GetTeam()
	starList := team.GetAllStarList()
	atlasMgr := team.GetAtlasMgr()
	staticDataMgr := GetServer().GetStaticDataMgr()
	for i := 0; i < len(starList); i++ {
		starInfo := team.GetStar(starList[i])

		starType := staticDataMgr.GetStarType(starInfo.Type)
		if starType.Class > 400 {
			//! A级以上的球员则加入图鉴
			atlasMgr.AddAtlas(team.GetID(), starInfo.Type, 0)
		}

	}

	client.SendTeam() ///创建组队成功,向客户端发送球队信息
	return true
}

type GetStarExpPoolMsg struct { ///使用经验池消息
	MsgHead   `json:"head"` /// "team", "getstarexppoolmsg"
	StarID    int           `json:"starid"`    ///训练球员ID
	GrowLevel int           `json:"growlevel"` ///成长到多少级
	IsOneKey  bool          `json:"isonekey"`  ///true 使用一键训练  false 不使用一键训练
}

func (self *GetStarExpPoolMsg) GetTypeAndAction() (string, string) {
	return "team", "getstarexppoolmsg"
}

func (self *GetStarExpPoolMsg) checkExpEnough(client IClient, expPool int, lastLevel int) (bool, int) {
	team := client.GetTeam()
	staticDataMgr := GetServer().GetStaticDataMgr()
	star := team.GetStar(self.StarID)
	starInfo := star.GetInfo()
	curExp := starInfo.Exp
	// curLevel := starInfo.Level

	needExp := staticDataMgr.GetLevelExpNeedExp(levelExpTypeStarLevel, lastLevel-1) ///需要多少经验
	//	redundancyExp := staticDataMgr.GetLevelExpNeedExp(levelExpTypeStarLevel, curLevel)
	//	redundancyExp = redundancyExp - curExp ///得到当前等级冗余经验
	needExp -= curExp

	if needExp < 0 {
		needExp = 0
	}

	return expPool >= needExp, needExp
}

func (self *GetStarExpPoolMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	loger := GetServer().GetLoger()

	///检查经验池数值是否合法
	expPool := team.GetInfo().StarExpPool
	if loger.CheckFail("expPool <= ExpPoolLimit", expPool <= ExpPoolLimit,
		expPool, ExpPoolLimit) {
		return false
	}

	if false == self.IsOneKey {
		///检查球员id是否合法
		star := team.GetStar(self.StarID)
		if loger.CheckFail("GetStar() != nil", star != nil, star, nil) {
			return false
		}

		///检查经验池是否足够成长指定级别
		isExpEnough, _ := self.checkExpEnough(client, expPool, self.GrowLevel)
		if loger.CheckFail("isExpEnough == true", isExpEnough == true,
			isExpEnough, true) {
			return false
		}
	}

	return true
}

func (self *GetStarExpPoolMsg) payAction(client IClient) bool {

	if false == self.IsOneKey {
		///扣除所需经验
		sync := client.GetSyncMgr()
		team := client.GetTeam()
		teamInfo := team.GetInfo()
		expPool := teamInfo.StarExpPool
		_, needExp := self.checkExpEnough(client, expPool, self.GrowLevel)
		teamInfo.StarExpPool -= needExp
		sync.SyncObject("GetStarExpPoolMsg", team)
	}

	return true
}

func (self *GetStarExpPoolMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	teamInfo := team.GetInfo()
	syncMgr := client.GetSyncMgr()
	expPool := teamInfo.StarExpPool
	loger := GetServer().GetLoger()
	syncObjectList := SyncObjectList{}
	if false == self.IsOneKey { ///不使用一键训练
		//给予玩家经验
		_, needExp := self.checkExpEnough(client, expPool, self.GrowLevel)
		star := team.GetStar(self.StarID)
		star.AwardExp(needExp)
		syncObjectList = append(syncObjectList, star)
		star.CalcScore()
	} else if true == self.IsOneKey { ///使用训练
		nTimes := 0
		for nTimes <= 100100 {
			//得到队伍球星平均等级
			nAverageLevel := team.GetAverageLevel()
			if loger.CheckFail("nAverageLevel < 100", nAverageLevel < 100,
				nAverageLevel, 100) {
				return false
			}

			minimumStar := team.GetMinimumLevelStar(nAverageLevel)
			if minimumStar == nil {
				break
			}
			starInfo := minimumStar.GetInfo()
			self.StarID = starInfo.ID
			upLevel := 0
			//判断经验池经验是否足够让该球员升级
			if nAverageLevel <= starInfo.Level {
				upLevel = starInfo.Level + 1
			} else {
				upLevel = (nAverageLevel - starInfo.Level) + starInfo.Level
			}
			isEnough, needExp := self.checkExpEnough(client, expPool, upLevel)
			if isEnough == false {
				minimumStar.AwardExp(expPool)
				teamInfo.StarExpPool = 0
				expPool = 0
				break

			} else {
				minimumStar.AwardExp(needExp)
				expPool -= needExp
				teamInfo.StarExpPool -= needExp
				syncObjectList = append(syncObjectList, minimumStar)
			}
			nTimes++
		}
	}
	team.CalcScore()                              ///计算新的球队评分
	syncMgr.SyncObject("GetStarExpPoolMsg", team) ///同步客户端

	if len(syncObjectList) > 0 { ///同步球员属性变更
		syncMgr.SyncObjectArray("GetStarExpPoolMsg", syncObjectList)
	}

	///更新球员训练中的日常任务
	team.GetTaskMgr().UpdateDayTaskFuntion(client.GetElement(), dayTaskFunctionStarTrain)

	return true
}

func (self *GetStarExpPoolMsg) processAction(client IClient) bool {
	if false == self.checkAction(client) {
		return false
	}

	if false == self.payAction(client) {
		return false
	}

	if false == self.doAction(client) {
		return false
	}

	return true
}

type BuyActionPointMsg struct {
	MsgHead `json:"head"` ///"team", "buyactionpoint"
}

func (self *BuyActionPointMsg) GetTypeAndAction() (string, string) {
	return "team", "buyactionpoint"
}

func (self *BuyActionPointMsg) checkAction(client IClient) bool {
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamShopTimes)
	if resetAttrib == nil {
		resetAttrib = team.CreateShopTimesDefaultRefresh() ///创建默认购买次数
	}

	shoppingTimes := resetAttrib.Value1 ///取得已购买次数
	if loger.CheckFail("shoppingTimes >= 0", shoppingTimes >= 0, shoppingTimes, 0) {
		return false
	}
	vipLevel := team.GetVipLevel()
	if loger.CheckFail("vipLevel >0", vipLevel > 0, vipLevel, 0) {
		return false
	}
	vipType := GetServer().GetStaticDataMgr().GetVipInfo(vipLevel)
	if loger.CheckFail("vipType !=nil", vipType != nil, vipType, nil) {
		return false ///非vip不能购买
	}
	canBuyTimes := vipType.Param3 ///可购买次数
	if loger.CheckFail("shoppingTimes<canBuyTimes", shoppingTimes < canBuyTimes, shoppingTimes, canBuyTimes) {
		return false ///vip购买次数不足
	}

	needMoney := 50 * (1 + shoppingTimes/2)
	needMoney = Min(300, needMoney) ///上限是300钻
	if loger.CheckFail("team.Ticket >= needMoney", team.Ticket >= needMoney, team.Ticket, needMoney) {
		return false ///球票不足以支付
	}
	return true
}

func (self *BuyActionPointMsg) payAction(client IClient) bool {
	team := client.GetTeam()
	syncMgr := client.GetSyncMgr()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamShopTimes)
	shoppingTimes := resetAttrib.Value1 ///取得购买次数
	needMoney := 50 * (1 + shoppingTimes/2)
	needMoney = Min(300, needMoney) ///上限是300钻
	team.PayTicket(needMoney)
	client.SetMoneyRecord(PlayerCostMoney, Pay_BuyActionPoint, needMoney, team.GetTicket()) ///记入后台记录
	curTimes := shoppingTimes + 1
	resetAttribMgr.UpdateResetAttrib(ResetAttribTypeTeamShopTimes, resetAttrib.ResetTime, IntList{curTimes}) ///累加购买次数
	syncMgr.SyncObject("BuyActionPointMsg", team)
	return true
}

func (self *BuyActionPointMsg) doAction(client IClient) bool {
	// /给予玩家行动点
	team := client.GetTeam()
	team.AwardObject(awardTypeActionPoint, 150, 0, 0)

	///返回客户端信息
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamShopTimes)

	shoppingTimes := resetAttrib.Value1
	restoreTime := team.Restoretime

	msg := NewQueryShoppingActionPointResultMsg(shoppingTimes, restoreTime)
	client.SendMsg(msg)

	///更新天天联赛中的日常任务
	team.GetTaskMgr().UpdateDayTaskFuntion(client.GetElement(), dayTaskFunctionBuyActionPoint)
	return true
}

func (self *BuyActionPointMsg) processAction(client IClient) bool {
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

type QueryShoppintActionPointInfoMsg struct { ///查询购买行动点信息
	MsgHead `json:"head"` //"team", "queryshoppingactionpointinfo"
}

func (self *QueryShoppintActionPointInfoMsg) GetTypeAndAction() (string, string) {
	return "team", "queryshoppingactionpointinfo"
}

func (self *QueryShoppintActionPointInfoMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	//	loger := GetServer().GetLoger()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamShopTimes)
	if resetAttrib == nil {
		resetAttrib = team.CreateShopTimesDefaultRefresh() ///创建默认刷新时间与购买次数
	}
	/*
		//检测购买次数是否足够
		currentBuyTimes := resetAttrib.Value1
		vipLevel := team.GetVipLevel()

		//限购次数
		vipPrivil := GetServer().GetStaticDataMgr().GetVipInfo(vipLevel)
		if loger.CheckFail("vipPrivil != nil", vipPrivil != nil, vipPrivil, nil) {
			return false
		}
		limitTimes := vipPrivil.Param3
		if loger.CheckFail("currentBuyTimes <= limitTimes", currentBuyTimes <= limitTimes, currentBuyTimes, limitTimes) {
			return false
		}
	*/
	return true
}

func (self *QueryShoppintActionPointInfoMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false {
		return false
	}

	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamShopTimes)

	shoppingTimes := resetAttrib.Value1
	restoreTime := team.Restoretime

	msg := NewQueryShoppingActionPointResultMsg(shoppingTimes, restoreTime)
	client.SendMsg(msg)
	return true
}

type QueryShoppingActionPointResultMsg struct { ///查询购买行动点界面信息返回结果
	MsgHead       `json:"head"` // "team", "queryshoppingactionpointresult"
	ShoppingTimes int           `json:"shoppingTimes"` ///当前购买次数
	RestoreTime   int           `json:"restoretime"`   ///下次恢复行动点时间  720秒/恢复1点  等于0时代表行动点恢复已至上限
}

func (self *QueryShoppingActionPointResultMsg) GetTypeAndAction() (string, string) {
	return "team", "queryshoppingactionpointresult"
}

func NewQueryShoppingActionPointResultMsg(shoppingTimes int, restoreTime int) *QueryShoppingActionPointResultMsg {
	msg := new(QueryShoppingActionPointResultMsg)
	msg.ShoppingTimes = shoppingTimes
	msg.RestoreTime = restoreTime
	return msg
}

type QueryGoldFingerInfoMsg struct { ///查询点金指
	MsgHead `json:"head"` // "team", "querygoldfingerinfo"
}

func (self *QueryGoldFingerInfoMsg) GetTypeAndAction() (string, string) {
	return "team", "querygoldfingerinfo"
}

type QueryGoldFingerInfoResultMsg struct { ///查询金手指回应消息
	MsgHead      `json:"head"` // "team", "querygoldfingerinforesult"
	NeedMoney    int           `json:"paymoney"`     // 需支付钻石数量
	GetCoin      int           `json:"getcoin"`      // 得到金币数量
	CurrentTimes int           `json:"currenttimes"` // 当前购买次数
	IsDouble     int           `json:"isdouble"`     //是否暴击
	Multiple     int           `json:"multiple"`     //暴击倍数
}

func (self *QueryGoldFingerInfoResultMsg) GetTypeAndAction() (string, string) {
	return "team", "querygoldfingerinforesult"
}

func SendQueryGoldFingerInfoResultMsg(client IClient, needMoney int, getCoin int, currentTimes int, isDouble int, multiple int) {
	msg := new(QueryGoldFingerInfoResultMsg)
	msg.NeedMoney = needMoney
	msg.GetCoin = getCoin
	msg.CurrentTimes = currentTimes
	msg.IsDouble = isDouble
	msg.Multiple = multiple
	client.SendMsg(msg)
}

func (self *QueryGoldFingerInfoMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamShopTimes)
	if resetAttrib == nil {
		resetAttrib = team.CreateShopTimesDefaultRefresh() ///创建默认刷新时间与购买次数
	}
	return true
}

func GetNeedAndPay(currentTimes int, freeTimes int, level int) (int, int) { //得到获取金币数量与支付钻石数量
	currentTimes += 1 ///下一点所需
	needMoney := 0    //支付钻石
	getMoney := 0     //获取金币
	if currentTimes >= freeTimes {
		//向上取整((次数-1)/3)*20+20
		//needMoney = int(math.Ceil(float64((currentTimes-1)/3)))*20 + 20
		needMoney = Min(5*(currentTimes-freeTimes), 30)
	} else {
		needMoney = 0
	}

	//10000*（1+0.1*玩家等级）*（1+0.2*次数）
	//getMoney = int(10000 * (1.0 + 0.1*float32(level)) * (1.0 + 0.2*float32(currentTimes)))
	getMoney = int(10000 * (1.0 + 0.1*float32(currentTimes)) * (1.0 + 0.05*float32(level)))
	return needMoney, getMoney

}
func (self *QueryGoldFingerInfoMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false {
		return false
	}
	team := client.GetTeam()
	staticDataMgr := GetServer().GetStaticDataMgr()
	//得到当前已购买次数
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamShopTimes)
	currentTimes := resetAttrib.Value2 ///得到下次所需金额

	//得到可免费使用金手指次数
	vipLevel := team.GetVipLevel()
	vipType := GetServer().GetStaticDataMgr().GetVipInfo(vipLevel)

	freeTimes := 1
	if vipType != nil {
		freeTimes = vipType.Param1
	}
	needMoney, getCoin := GetNeedAndPay(currentTimes, freeTimes, team.GetLevel())

	needMoney = Max(needMoney, 0)
	multiple := staticDataMgr.GetConfigStaticDataInt(configGoldFinger, configGoldFingerCommonConfig, 2) //倍数

	SendQueryGoldFingerInfoResultMsg(client, needMoney, getCoin, currentTimes, resetAttrib.Value3, multiple)
	return true
}

type UseGoldFingerMsg struct {
	MsgHead `json:"head"` //"team", "usegoldfinger"
}

func (self *UseGoldFingerMsg) GetTypeAndAction() (string, string) {
	return "team", "usegoldfinger"
}

func (self *UseGoldFingerMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	loger := GetServer().GetLoger()
	//得到当前已购买次数
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamShopTimes)
	currentBuyTimes := resetAttrib.Value2

	//得到可免费使用金手指次数
	vipLevel := team.GetVipLevel()
	vipType := GetServer().GetStaticDataMgr().GetVipInfo(vipLevel)

	//	if loger.CheckFail("vipType != nil", vipType != nil, vipType, nil) {
	//		return false //VIP
	//	}
	freeTimes := 1 ///vip0的免费次数
	buyTimes := 2  ///vip0的购买次数
	if vipType != nil {
		freeTimes = vipType.Param1 //免费次数
		buyTimes = vipType.Param2  //购买次数
	}

	needMoney, _ := GetNeedAndPay(currentBuyTimes, freeTimes, team.GetLevel())
	needMoney = Max(needMoney, 0)

	if loger.CheckFail("currentBuyTimes < freeTimes+buyTimes", currentBuyTimes < (freeTimes+buyTimes), currentBuyTimes, freeTimes+buyTimes) {
		return false //购买次数已到上限
	}

	if loger.CheckFail("team.GetTicket() >= needMoney", team.GetTicket() >= needMoney, team.GetTicket(), needMoney) {
		return false //钱不够
	}
	return true
}

func (self *UseGoldFingerMsg) payAction(client IClient) bool {
	//得到当前已购买次数
	team := client.GetTeam()
	syncMgr := client.GetSyncMgr()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamShopTimes)
	currentBuyTimes := resetAttrib.Value2

	//得到可免费使用金手指次数
	vipLevel := team.GetVipLevel()
	vipType := GetServer().GetStaticDataMgr().GetVipInfo(vipLevel)

	freeTimes := 1 ///vip0的免费次数
	if vipType != nil {
		freeTimes = vipType.Param1 //免费次数
	}
	needMoney, _ := GetNeedAndPay(currentBuyTimes, freeTimes, team.GetLevel())
	needMoney = Max(needMoney, 0)

	if needMoney != 0 {
		client.SetMoneyRecord(PlayerCostMoney, Pay_GodFinger, needMoney, team.GetTicket())
	}

	team.PayTicket(needMoney)
	syncMgr.SyncObject("UseGoldFingerMsg", team)
	return true
}

func (self *UseGoldFingerMsg) isDouble() bool { //是否暴击
	staticDataMgr := GetServer().GetStaticDataMgr()
	doubleRate := staticDataMgr.GetConfigStaticDataInt(configGoldFinger, configGoldFingerCommonConfig, 1) //爆率

	rand := Random(0, 100)
	return rand < doubleRate
}

func (self *UseGoldFingerMsg) doAction(client IClient) bool {
	//得到当前已购买次数
	team := client.GetTeam()
	staticDataMgr := GetServer().GetStaticDataMgr()
	resetAttribMgr := team.GetResetAttribMgr()
	resetAttrib := resetAttribMgr.GetResetAttrib(ResetAttribTypeTeamShopTimes)
	currentBuyTimes := resetAttrib.Value2

	//得到可免费使用金手指次数
	vipLevel := team.GetVipLevel()
	vipType := GetServer().GetStaticDataMgr().GetVipInfo(vipLevel)
	multiple := staticDataMgr.GetConfigStaticDataInt(configGoldFinger, configGoldFingerCommonConfig, 2) //倍数

	freeTimes := 1 ///vip0的免费次数
	if vipType != nil {
		freeTimes = vipType.Param1 //免费次数
	}
	needMoney, getCoin := GetNeedAndPay(currentBuyTimes, freeTimes, team.GetLevel())
	needMoney = Max(needMoney, 0)

	if resetAttrib.Value3 == 1 {
		getCoin *= multiple
	}
	team.AwardObject(awardTypeCoin, getCoin, 0, 0)

	resetAttrib.Value2 += 1

	// 重新计算下次点金手参数
	currentBuyTimes = resetAttrib.Value2
	needMoney, getCoin = GetNeedAndPay(currentBuyTimes, freeTimes, team.GetLevel())
	needMoney = Max(needMoney, 0)

	// 计算爆率
	isDouble := self.isDouble()
	value := 0
	if isDouble == true {
		value = 1 //下次暴击
	}
	resetAttribMgr.UpdateResetAttrib(ResetAttribTypeTeamShopTimes, resetAttrib.ResetTime, IntList{resetAttrib.Value1, resetAttrib.Value2, value})
	SendQueryGoldFingerInfoResultMsg(client, needMoney, getCoin, currentBuyTimes, value, multiple)

	///更新天天联赛中的日常任务
	team.GetTaskMgr().UpdateDayTaskFuntion(client.GetElement(), dayTaskFunctionBuyCoin)
	return true
}

func (self *UseGoldFingerMsg) processAction(client IClient) bool {
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

type ClientOperationMsg struct { //服务端与客户端都会发送,用于Set与Get存储串
	MsgHead         `json:"head"` //"team", "clientoperation"
	OperationType   int           //1为Get 2为Set
	OperationString string        //操作串
}

func (self *ClientOperationMsg) GetTypeAndAction() (string, string) {
	return "team", "clientoperation"
}

func (self *ClientOperationMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	loger := GetServer().GetLoger()
	OperationStringLen := len(team.Operation)
	if loger.CheckFail("OperationStringLen <= OperationStringMaxLen", OperationStringLen <= OperationStringMaxLen,
		OperationStringLen, OperationStringMaxLen) {
		return false //超过最大上限
	}
	return true
}

func (self *ClientOperationMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false {
		return false
	}

	team := client.GetTeam()
	switch self.OperationType {
	case GetOperationString: ///得到串
		self.sendClientOperationMsg(client, self.OperationType, team.Operation)
	case SetOperationString: ///设置串
		team.Operation = self.OperationString
	}
	return true
}

func (self *ClientOperationMsg) sendClientOperationMsg(client IClient, operationType int, operationString string) {
	msg := new(ClientOperationMsg)
	msg.OperationType = operationType
	msg.OperationString = operationString
	client.SendMsg(msg)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//! 消息
//! client2server
type AddStarPosMsg struct {
	MsgHead `json:"head"` //! "team", "addstarpos"
}

func (self *AddStarPosMsg) GetTypeAndAction() (string, string) {
	return "team", "addstarpos"
}

func (self *AddStarPosMsg) processAction(client IClient) bool {
	if self.doAction(client) == false {
		return false
	}
	return true
}

func (self *AddStarPosMsg) doAction(client IClient) bool {
	team := client.GetTeam()

	result := team.AwardStarPos()

	addStarPosResultMsg := new(AddStarPosResultMsg)
	addStarPosResultMsg.Result = result
	team.client.SendMsg(addStarPosResultMsg)

	return true
}

//! server2client
type AddStarPosResultMsg struct {
	MsgHead `json:"head"` //! "team", "addstarposresult"
	Result  int           `json:"result"` //! 0成功 1已到最大 2钻石不足
}

func (self *AddStarPosResultMsg) GetTypeAndAction() (string, string) {
	return "team", "addstarposresult"
}

type MatchFlowMsg struct {
	MsgHead  `json:"head"` ///"team", "matchflow"
	Flowlist MatchFlowList `json:"matchflowlist"`
}

func (self *MatchFlowMsg) GetTypeAndAction() (string, string) {
	return "team", "matchflow"
}

func SendMatchFlowMsg(client IClient, matchList MatchFlowList) {
	msg := new(MatchFlowMsg)
	msg.Flowlist = matchList
	client.SendMsg(msg)
}

type UpdateTeamIconAndShirtMsg struct { //! 更新队徽与队服
	MsgHead `json:"head"` //! "team", "updateteamiconandshirt"
	Icon    int           `json:"icon"`  //!队徽
	Shirt   int           `json:"shirt"` //!队服
}

func (self *UpdateTeamIconAndShirtMsg) GetTypeAndAction() (string, string) {
	return "team", "updateteamiconandshirt"
}

func (self *UpdateTeamIconAndShirtMsg) processAction(client IClient) bool {

	syncMgr := client.GetSyncMgr()
	team := client.GetTeam()
	team.Icon = self.Icon
	team.TeamShirts = self.Shirt

	syncMgr.SyncObject("UpdateTeamIconAndShirtMsg", team) //!同步改变的队服与队徽
	return true
}
