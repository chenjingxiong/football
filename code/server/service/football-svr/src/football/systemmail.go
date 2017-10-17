package football

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
	"time"
)

const (
	SystemMailNoDeal   = 1 //没处理
	SystemMailDealDone = 2 //处理完
)

const (
	SysMailTargetPersonal = 0 //私人
	SysMailTargetOnline   = 1 //在线者
	SysMailTargetAllUser  = 2 //所有
)

//const (
//	NomalSend = 1 //使用接口发送
//	SpecSend  = 2 //使用存储数据库并广播的方式发送
//)

type OnlineTeamNum struct { ///在线人数
	ID     int
	Teamid int //teamid
}

func (self *OnlineTeamNum) GetReflectValue() reflect.Value { ///得到球队反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

type SystemMailType struct {
	ID       int //记录ID
	Pid      int
	Itemid   int
	Itemnum  int    //1 为未处理 2 处理完毕
	Theme    string // 0 个人 1 在线用户 2 全服用户
	Contents string //按照"sort:XXXX id1:XXXXXX num1:XXXXXX grade1:XXXXXX"格式解读
	Code     int
	ServerID int
}

type SystemDetailMailType struct {
	ID     int //记录ID
	State  int //状态: 0 未处理 1 处理中 2 处理完毕
	Target int //目标:0 个人 1 在线用户 2 全服用户
	Sort   int //邮件模板类型 1系统发物品奖励 2系统战报 3.补偿邮件 4.奖励邮件

	Sendtime int //定时发放

	Teamid1  int //用户1
	Teamid2  int //用户2
	Teamid3  int //用户3
	Teamid4  int //用户4
	Teamid5  int //用户5
	Teamid6  int //用户6
	Teamid7  int //用户7
	Teamid8  int //用户8
	Teamid9  int //用户9
	Teamid10 int //用户10

	AwardItem1  int ///道具类型
	AwardGrade1 int ///道具品质
	AwardCount1 int ///道具数量
	AwardItem2  int ///道具类型
	AwardGrade2 int ///道具品质
	AwardCount2 int ///道具数量
	AwardItem3  int ///道具类型
	AwardGrade3 int ///道具品质
	AwardCount3 int ///道具数量
	AwardItem4  int ///道具类型
	AwardGrade4 int ///道具品质
	AwardCount4 int ///道具数量
	AwardItem5  int ///道具类型
	AwardGrade5 int ///道具品质
	AwardCount5 int ///道具数量

	MakeTime string ///生成时间

}

type SystemMailInfo struct {
	SystemMailType
	DataUpdater
	bChg bool //修改标记
}

func (self *SystemMailInfo) GetReflectValue() reflect.Value { ///得到球队反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

type SystemMailInfoList map[int]*SystemMailInfo
type DispatchMailList map[int]*SystemMailInfo
type SystemMailMgr struct {
	systemMailList SystemMailInfoList
	userMgr        *UserMgr
	dealList       DispatchMailList
}

func (self *SystemMailMgr) Run() {
	serverUpdateTimer := time.NewTicker(time.Second * 1) ///邮件系统逻辑每1秒执行一次
	for {
		select {
		case now := <-serverUpdateTimer.C:
			self.OnTimer(now)

		}
	}
}

func (self *SystemMailMgr) OnTimer(now time.Time) {
	defer func() {
		x := recover()
		if x != nil {
			GetServer().GetLoger().Error("%v\r\n%s", x, debug.Stack())
		}
	}()
	if now.Unix()%60 == 0 {

		self.GetMailListFromDB()
		self.DealSystemMail()
	}
	//self.DealMailSend() //邮件发放
}

func (self *SystemMailMgr) Init(userMgr *UserMgr) { // 初始化参数
	self.systemMailList = make(SystemMailInfoList)
	self.userMgr = userMgr
	self.dealList = make(DispatchMailList)
}

//读取方式: 十分钟一次
func (self *SystemMailMgr) GetMailListFromDB() {
	if self.systemMailList == nil {
		self.systemMailList = make(SystemMailInfoList)
	}

	for i, _ := range self.systemMailList {
		delete(self.systemMailList, i) ///清空Map
	}
	//tableSystemMail
	mailListQuery := fmt.Sprintf("select * from %s where itemnum = %d and serverid = %d limit 10", "mmo2d_userljzm1.ht_mail", SystemMailNoDeal, GetServer().config.ServerID)
	mailType := new(SystemMailType)
	mailInfoList := GetServer().GetDynamicDB().fetchAllRows(mailListQuery, mailType)

	for _, v := range mailInfoList {
		mailType = v.(*SystemMailType)
		mailInfo := new(SystemMailInfo)
		mailInfo.SystemMailType = *mailType
		mailInfo.bChg = false
		self.systemMailList[mailInfo.ID] = mailInfo
		//fmt.Println("Get %v", mailInfo)
	}
}

//邮件处理消息派发队列 速度:500封/s
func (self *SystemMailMgr) DealMailSend(teamid int, mailInfo *SystemMailInfo) {
	// if len(self.dealList) <= 0 {
	// 	return //无邮件处理
	// }

	// mailNum := 0
	// //dealList := make(DispatchMailList)
	// for teamid, mailInfo := range self.dealList {
	// 	// dealList[teamid] = mailInfo
	// 	mailNum++
	// 	if mailNum == 500 {
	// 		break
	// 	}
	mailType := new(SystemDetailMailType)

	fmt.Sscanf(mailInfo.Contents,
		"sort:%d id1:%d num1:%d grade1:%d id2:%d num2:%d grade2:%d id3:%d num3:%d grade3:%d id4:%d num4:%d grade4:%d id5:%d num5:%d grade5:%d",
		&mailType.Sort,
		&mailType.AwardItem1, &mailType.AwardCount1, &mailType.AwardGrade1,
		&mailType.AwardItem2, &mailType.AwardCount2, &mailType.AwardGrade2,
		&mailType.AwardItem3, &mailType.AwardCount3, &mailType.AwardGrade3,
		&mailType.AwardItem4, &mailType.AwardCount4, &mailType.AwardGrade4,
		&mailType.AwardItem5, &mailType.AwardCount5, &mailType.AwardGrade5)

	awardList := IntList{mailType.AwardItem1, mailType.AwardItem2, mailType.AwardItem3, mailType.AwardItem4, mailType.AwardItem5}
	countList := IntList{mailType.AwardCount1, mailType.AwardCount2, mailType.AwardCount3, mailType.AwardCount4, mailType.AwardCount5}
	gradeList := IntList{mailType.AwardGrade1, mailType.AwardGrade2, mailType.AwardGrade3, mailType.AwardGrade4, mailType.AwardGrade5}

	client := self.userMgr.GetClientByTeamID(teamid)
	if client == nil {
		self.SendSysAwardMail(teamid, mailType.Sort, SystemMail, awardList, gradeList, countList) // 玩家不在线,直写数据库
		return
	}

	team := client.GetTeam()
	mailMgr := team.GetMailMgr()
	mailMgr.SendSysAwardMail(mailType.Sort, SystemMail, awardList, gradeList, countList, "", "")

	// 	delete(self.dealList, teamid)
	// }

	// self.MailSend(dealList)
}

//处理邮件发送
// func (self *SystemMailMgr) MailSend(dealList DispatchMailList) {

// 	// 发送
// 	for teamid, mailInfo := range dealList {
// 		awardList := IntList{mailInfo.AwardItem1, mailInfo.AwardItem2, mailInfo.AwardItem3, mailInfo.AwardItem4, mailInfo.AwardItem5}
// 		countList := IntList{mailInfo.AwardCount1, mailInfo.AwardCount2, mailInfo.AwardCount3, mailInfo.AwardCount4, mailInfo.AwardCount5}
// 		gradeList := IntList{mailInfo.AwardGrade1, mailInfo.AwardGrade2, mailInfo.AwardGrade3, mailInfo.AwardGrade4, mailInfo.AwardGrade5}
// 		client := self.userMgr.GetClientByTeamID(teamid)
// 		if client == nil {
// 			self.SendSysAwardMail(teamid, mailInfo.Sort, SystemMail, awardList, countList, gradeList)
// 			continue // 玩家不在线,直写数据库
// 		}

// 		team := client.GetTeam()
// 		mailMgr := team.GetMailMgr()
// 		mailMgr.SendSysAwardMail(mailInfo.Sort, SystemMail, awardList, countList, gradeList)

// 		//fmt.Println("邮件发送 teamID:%d", team.ID)
// 	}

// }

//存储方式: 即时存储
func (self *SystemMailMgr) SaveMailListToDB() {
	if self.systemMailList == nil {
		return //没有系统邮件
	}

	for i, v := range self.systemMailList {
		systemMailInfo := v
		if v.bChg == false {
			continue //无改变则跳过
		}

		//	mailListUpdate := fmt.Sprintf("update %s set state = %d where id = %d", tableSystemMail, systemMailInfo.State, systemMailInfo.ID)
		mailListUpdate := fmt.Sprintf("update %s set itemnum = %d where id = %d", "mmo2d_userljzm1.ht_mail", systemMailInfo.Itemnum, systemMailInfo.ID)
		GetServer().GetDynamicDB().Exec(mailListUpdate)

		// 如果处理完毕,则从内存删除该信息
		if systemMailInfo.Itemnum == SystemMailDealDone {
			delete(self.systemMailList, i)
		}
	}

}

func (self *SystemMailMgr) DealSystemMail() {

	if self.systemMailList == nil {
		return
	}
	checkLifeTimer := time.NewTicker(time.Second * 1)
	for _, v := range self.systemMailList {
		mailInfo := v
		if mailInfo.Itemnum == SystemMailDealDone {
			continue
		}

		// if mailInfo.Sendtime > Now() {
		// 	continue
		// }

		target, _ := strconv.Atoi(mailInfo.Theme)
		switch target { //目标
		case SysMailTargetPersonal: //个人
			// teamList := IntList{mailInfo.Teamid1, mailInfo.Teamid2, mailInfo.Teamid3, mailInfo.Teamid4, mailInfo.Teamid5,
			// 	mailInfo.Teamid6, mailInfo.Teamid7, mailInfo.Teamid8, mailInfo.Teamid9, mailInfo.Teamid10}

			// for i := range teamList {

			// 	if teamList[i] == 0 {
			// 		continue
			// 	}
			self.DealMailSend(mailInfo.Itemid, mailInfo)
			// }

		case SysMailTargetOnline: //在线
			teamidList := self.GetOnlineTeamID()
			nIndex := 0
			for i := range teamidList {
				nIndex += 1
				self.DealMailSend(teamidList[i], mailInfo)
				if nIndex >= 500 {
					nIndex = 0
					<-checkLifeTimer.C
				}
			}

		case SysMailTargetAllUser: //所有
			teamidList := self.GetAllTeamID()
			nIndex := 0
			for i := range teamidList {
				nIndex += 1
				self.DealMailSend(teamidList[i], mailInfo)
				if nIndex >= 500 {
					nIndex = 0
					<-checkLifeTimer.C
				}
			}
		}

		mailInfo.Itemnum = SystemMailDealDone
		mailInfo.bChg = true
		self.SaveMailListToDB()
	}

}

func (self *SystemMailMgr) GetAllTeamID() IntList { //得到所有在线或不在线的玩家队伍id
	staticDataMgr := GetServer().GetStaticDataMgr()
	serverLimit := staticDataMgr.GetConfigStaticDataInt("server", "commonconfig", 3)
	userQuery := fmt.Sprintf("select * from dy_team limit %d", serverLimit)
	teamInfo := new(TeamInfo)
	teamInfoList := GetServer().GetDynamicDB().fetchAllRows(userQuery, teamInfo)

	teamIDList := IntList{}
	for _, v := range teamInfoList {
		teamInfo = v.(*TeamInfo)
		teamIDList = append(teamIDList, teamInfo.ID)
	}

	return teamIDList
}

func (self *SystemMailMgr) GetOnlineTeamID() IntList { //得到所有在线玩家队伍id
	staticDataMgr := GetServer().GetStaticDataMgr()
	serverLimit := staticDataMgr.GetConfigStaticDataInt("server", "commonconfig", 3)
	userQuery := fmt.Sprintf("select * from %s limit %d", tableRecordOnline, serverLimit)
	onlineTeamNum := new(OnlineTeamNum)

	teamidList := GetServer().GetRecordDB().fetchAllRows(userQuery, onlineTeamNum)

	teamIDList := IntList{}
	for _, v := range teamidList {
		onlineTeamNum = v.(*OnlineTeamNum)
		teamIDList = append(teamIDList, onlineTeamNum.Teamid)
	}

	return teamIDList
}

///发送系统奖励邮件,共指定三个属性列表,类型,品质(可以为0),数量
func (self *SystemMailMgr) SendSysAwardMail(teamid int, mailSort int, senderType int, awardItemList IntList, awardGradeList IntList,
	awardCountList IntList) {
	mail := new(Mail)
	mail.TeamID = teamid
	mail.Type = SystemMail
	mail.Sort = mailSort
	mail.State = StateNewMail
	mail.MakeTime = Now()
	mail.SenderType = senderType

	for i := range awardItemList {
		awardItemType := awardItemList[i]
		awardItemGrade := awardGradeList[i]
		awardItemCount := awardCountList[i]
		mail.SetAwardInfo(i+1, awardItemType, awardItemGrade, awardItemCount)
	}

	mail.InitDataUpdater(tableMail, &mail.MailInfo)
	insertSql := mail.InsertSql()
	lastInsertItemID, _ := GetServer().GetDynamicDB().Exec(insertSql)
	if lastInsertItemID <= 0 {
		GetServer().GetLoger().Warn("MailMgr createNewMail is fail! teamid: %d", teamid)
	}
}
