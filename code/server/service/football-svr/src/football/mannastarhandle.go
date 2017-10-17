package football

import (
// "fmt"
)

const (
	MannaStarSeatOne   = 1 //! 天赐球员一号位
	MannaStarSeatTwo   = 2 //! 天赐球员二号位
	MannaStarSeatThree = 3 //! 天赐球员三号位
)

type UpdateMannaStarMsg struct { //! 修改自创球员
	MsgHead `json:"head"` //! "mannastar", "updatemannastar"
	ID      int           `json:"id"`      //! 自创球员X号位
	Name    string        `json:"name"`    //! 球员名称
	Icon    int           `json:"icon"`    //! 头像
	Face    int           `json:"face"`    //! 脸
	Seat1   int           `json:"seat1"`   //! 位置
	Seat2   int           `json:"seat2"`   //! 位置
	Seat3   int           `json:"seat3"`   //! 位置
	Hair    int           `json:"hair"`    ///头发
	Eyebrow int           `json:"eyebrow"` ///眉毛 (字段无用)
	Mouth   int           `json:"mouth"`   ///嘴巴
	Eye     int           `json:"eye"`     ///眼睛
	Skin    int           `json:"skin"`    ///皮肤
	Clothes int           `json:"clothes"` ///衣服
}

type UpdateMannaStarResultMsg struct {
	MsgHead   `json:"head"` //! "mannastar", "updatemannastarresult"
	MannaStar Star          `json:"mannastar"` //!天赐球员
	Result    int           `json:"result"`    //! 1成功 2 星级不到 3 队伍已满 4 名字过长  5 未知
}

func (self *UpdateMannaStarResultMsg) GetTypeAndAction() (string, string) {
	return "mannastar", "updatemannastarresult"
}

func (self *UpdateMannaStarMsg) GetTypeAndAction() (string, string) {
	return "mannastar", "updatemannastar"
}

func (self *UpdateMannaStarMsg) checkOpenCondition(team *Team) bool { //!检查自创开启条件是否满足

	//! 得到球队拥有的自创球员的最大星级
	loger := GetServer().GetLoger()
	starIDList := team.GetAllStarList()
	maxStarCount := 0
	for i := 0; i < len(starIDList); i++ {
		starID := starIDList[i]
		star := team.GetStar(starID)
		if star.IsMannaStar == 1 { //! 若该球星为自创球星
			maxStarCount = Max(star.EvolveCount, maxStarCount)
		}
	}

	if self.ID == MannaStarSeatTwo {
		if loger.CheckFail("maxStarCount >= NeedStarCount", maxStarCount >= 5, maxStarCount, 5) {
			return false //! 开启二号位需要自创球星达到5星级
		}
	}

	if self.ID == MannaStarSeatThree {
		if loger.CheckFail("maxStarCount >= NeedStarCount", maxStarCount >= 7, maxStarCount, 7) {
			return false //! 开启三号位需要自创球星达到7星级
		}
	}

	return true
}

func (self *UpdateMannaStarMsg) checkAction(client IClient) bool {
	loger := GetServer().GetLoger()
	team := client.GetTeam()

	//! 检查自创球员位置合理性
	if self.ID > MannaStarSeatThree || self.ID < MannaStarSeatOne {
		loger.Warn("UpdateMannaStarMsg ID is not legal")
		return false
	}

	//! 名字不得超过六个汉字
	if loger.CheckFail("len(self.Name) < 32", len(self.Name) < 32, len(self.Name), 32) {
		msg := new(UpdateMannaStarResultMsg)
		msg.Result = 4
		client.SendMsg(msg)

		return false
	}

	//! 自己的自创球员不许重名
	// nameSame := false
	// mannaStarMgr := team.GetMannaStarMgr()
	// for i := MannaStarSeatOne; i <= MannaStarSeatThree; i++ {
	// 	star := mannaStarMgr.GetMannaStarFromSeat(i)
	// 	if star == nil {
	// 		continue
	// 	}

	// 	if star.Name == self.Name && self.ID != star.MannaSeat {
	// 		nameSame = true
	// 		break
	// 	}
	// }

	// if loger.CheckFail("nameSame == false", nameSame == false, nameSame, false) {
	// 	return false
	// }

	//! 计算球员取值ID
	starTypeID := 10000 + self.ID*100 + self.Seat1

	staticDataMgr := GetServer().GetStaticDataMgr()
	starType := staticDataMgr.GetStarType(starTypeID)

	//! 球员属性数据不得为nil
	if loger.CheckFail("starType != nil", starType != nil, starType, nil) {
		return false
	}

	//! 检测星级是否满足开启条件
	isCanOpen := self.checkOpenCondition(team)
	if isCanOpen == false {
		msg := new(UpdateMannaStarResultMsg)
		msg.Result = 2
		client.SendMsg(msg)

		return false
	}

	//! 检测队伍是否已经满了
	// isFull := team.IsStarFull()
	// if loger.CheckFail("isFull == false", isFull == false, isFull, false) {
	// 	msg := new(UpdateMannaStarResultMsg)
	// 	msg.Result = 3
	// 	client.SendMsg(msg)

	// 	return false
	// }

	return true
}

// func (self *UpdateMannaStarMsg) payAction(client IClient) bool {
// 	return true
// }

// func (self *UpdateMannaStarMsg) UpdateStar(star *MannaStarType, starAttr *StarTypeStaticData) *MannaStarType {
// 	star.Grade = starAttr.Grade
// 	star.Card = starAttr.Class
// 	star.Icon = self.Icon
// 	star.Face = self.Face
// 	star.Seat1 = self.Seat1
// 	star.Seat2 = self.Seat2 //starAttr.Seat2
// 	star.Seat3 = self.Seat3 //starAttr.Seat3
// 	star.Nationality = starAttr.Nationality
// 	star.Pass = starAttr.Pass
// 	star.Steals = starAttr.Steals
// 	star.Dribbling = starAttr.Dribbling
// 	star.Sliding = starAttr.Sliding
// 	star.Shooting = starAttr.Shooting
// 	star.GoalKeeping = starAttr.GoalKeeping
// 	star.Body = starAttr.Body
// 	star.Speed = starAttr.Speed
// 	star.PassGrow = starAttr.PassGrow
// 	star.StealsGrow = starAttr.StealsGrow
// 	star.DribblingGrow = starAttr.DribblingGrow
// 	star.SlidingGrow = starAttr.SlidingGrow
// 	star.ShootingGrow = starAttr.ShootingGrow
// 	star.GoalKeepingGrow = starAttr.GoalKeepingGrow
// 	star.Skill1 = starAttr.Skill1
// 	star.Skill2 = starAttr.Skill2
// 	star.Skill3 = starAttr.Skill3
// 	star.Skill4 = starAttr.Skill4
// 	star.BasePrice = starAttr.BasePrice
// 	star.BaseScore = starAttr.BaseScore
// 	star.Ticket = starAttr.Ticket
// 	star.Fate1 = starAttr.Fate1
// 	star.Fate2 = starAttr.Fate2
// 	star.Fate3 = starAttr.Fate3
// 	star.Fate4 = starAttr.Fate4
// 	star.Fate5 = starAttr.Fate5
// 	star.Fate6 = starAttr.Fate6
// 	star.Item = starAttr.Item
// 	star.Team = starAttr.Team

// 	star.Hair = self.Hair
// 	star.Eyebrow = self.Eyebrow //! 保存模型数据
// 	star.Mouth = self.Mouth
// 	star.Eye = self.Eye
// 	star.Skin = self.Skin
// 	star.Clothes = self.Clothes

// 	star.Desc = starAttr.Desc

// 	return star
// }

func (self *UpdateMannaStarMsg) doAction(client IClient) bool {
	//! 根据球员位置取得对应属性信息
	staticDataMgr := GetServer().GetStaticDataMgr()
	team := client.GetTeam()
	mannaStarMgr := team.GetMannaStarMgr()
	mannaStar := mannaStarMgr.GetMannaStarFromSeat(self.ID)
	// if mannaStar == nil { //! 原创球星不存在
	// 	star := new(MannaStarType)
	// 	starAttr := staticDataMgr.GetStarType(10000 + self.ID*100 + self.Seat1)
	// 	star = self.UpdateStar(star, starAttr)
	// 	star.Name = self.Name
	// 	star.Teamid = team.GetID()
	// 	star.MannaSeat = self.ID
	// 	star.ID = mannaStarMgr.AddMannaStar(star)
	// 	mannaStarMgr.starList[star.ID] = NewMannaStar(star)
	// 	newStarID := team.AwardMannaStar(star.ID)

	// 	//!创建时,根据开启位给予星级
	// 	starCount := 0
	// 	if self.ID == 2 {
	// 		starCount = 4
	// 	} else if self.ID == 3 {
	// 		starCount = 6
	// 	}

	// 	//!天赐球员第一位没有初始属性与星际
	// 	if starCount != 0 {
	// 		starInfo := team.GetStar(star.ID)
	// 		starInfo.EvolveCount = starCount
	// 		starInfo.ChangeStarAttribute(starCount)
	// 	}

	// 	msg := new(UpdateMannaStarResultMsg)
	// 	msg.Result = "ok"
	// 	client.SendMsg(msg)

	// 	syncMgr := client.GetSyncMgr()
	// 	syncMgr.syncAddStar(IntList{newStarID})

	// 	mannaStarMgr.UpdateMannaStarSeat()

	// 	return true
	// }

	//! 更新数据库
	syncMgr := client.GetSyncMgr()
	starAttr := staticDataMgr.GetStarType(10000 + self.ID*100 + self.Seat1)

	//	mannaStar.MannaStarType = *self.UpdateStar(&star, starAttr)
	mannaStar.Grade = 1
	mannaStar.Card = starAttr.Class
	mannaStar.Icon = self.Icon
	mannaStar.Face = self.Face

	if self.Seat1 == mannaStar.Seat1 {
		mannaStar.Seat1 = self.Seat1
		mannaStar.Seat2 = self.Seat2 //starAttr.Seat2
		mannaStar.Seat3 = self.Seat3 //starAttr.Seat3
	} else {
		mannaStar.Seat1 = self.Seat1
		mannaStar.Seat2 = starAttr.Seat2 //starAttr.Seat2
		mannaStar.Seat3 = starAttr.Seat3 //starAttr.Seat3

		//! 切换主位置时候 删除球员以前学习技能
		skillMgr := team.GetSkillMgr()

		star := team.GetStarFromType(mannaStar.ID)
		skillList := skillMgr.GetStarSkillSlice(star.ID)

		for i := 0; i < len(skillList); i++ {
			skillInfo := skillList[i]
			skillMgr.RemoveStarSkill(skillInfo.ID)
		}

		//! 清空正在学习的技能
		resetAttrMgr := team.GetResetAttribMgr()
		skillAttr := resetAttrMgr.GetResetAttrib(ResetAttribTypeSkillStudyInfo)
		if skillAttr != nil {
			skillAttr.ResetTime = -1
			skillAttr.Value1 = 0
			skillAttr.Value2 = 0
			skillAttr.Value3 = 0
			skillAttr.Save()
		}

	}

	mannaStar.Nationality = starAttr.Nationality
	mannaStar.Pass = starAttr.Pass
	mannaStar.Steals = starAttr.Steals
	mannaStar.Dribbling = starAttr.Dribbling
	mannaStar.Sliding = starAttr.Sliding
	mannaStar.Shooting = starAttr.Shooting
	mannaStar.GoalKeeping = starAttr.GoalKeeping
	mannaStar.Body = starAttr.Body
	mannaStar.Speed = starAttr.Speed
	mannaStar.PassGrow = starAttr.PassGrow
	mannaStar.StealsGrow = starAttr.StealsGrow
	mannaStar.DribblingGrow = starAttr.DribblingGrow
	mannaStar.SlidingGrow = starAttr.SlidingGrow
	mannaStar.ShootingGrow = starAttr.ShootingGrow
	mannaStar.GoalKeepingGrow = starAttr.GoalKeepingGrow
	mannaStar.Skill1 = starAttr.Skill1
	mannaStar.Skill2 = starAttr.Skill2
	mannaStar.Skill3 = starAttr.Skill3
	mannaStar.Skill4 = starAttr.Skill4
	mannaStar.BasePrice = starAttr.BasePrice
	mannaStar.BaseScore = starAttr.BaseScore
	mannaStar.Ticket = starAttr.Ticket
	mannaStar.Fate1 = starAttr.Fate1
	mannaStar.Fate2 = starAttr.Fate2
	mannaStar.Fate3 = starAttr.Fate3
	mannaStar.Fate4 = starAttr.Fate4
	mannaStar.Fate5 = starAttr.Fate5
	mannaStar.Fate6 = starAttr.Fate6
	mannaStar.Item = starAttr.Item
	mannaStar.Team = starAttr.Team

	mannaStar.Hair = self.Hair
	mannaStar.Eyebrow = self.Eyebrow //! 保存模型数据
	mannaStar.Mouth = self.Mouth
	mannaStar.Eye = self.Eye
	mannaStar.Skin = self.Skin
	mannaStar.Clothes = self.Clothes

	mannaStar.Desc = starAttr.Desc

	//	fmt.Println("mannaStar: ", mannaStar)

	mannaStar.Name = self.Name
	mannaStar.Teamid = team.GetID()
	mannaStar.MannaSeat = self.ID

	// fmt.Println("mannaStar.Name: ", mannaStar.Name)
	// fmt.Println("mannaStar: ", mannaStar)
	// fmt.Println("name: ", self.Name)

	//	testStar := mannaStarMgr.GetMannaStarFromSeat(self.ID)

	//fmt.Println("test: ", testStar)
	//mannaStar.Save()

	isHasStar := false
	starList := team.GetAllStarList()
	for i := 0; i < starList.Len(); i++ {
		starInfo := team.GetStar(starList[i])
		if starInfo.Type == mannaStar.ID {
			isHasStar = true
		}
	}

	if isHasStar == false {
		newStarID := team.AwardMannaStar(mannaStar.ID)
		if newStarID == 0 {
			GetServer().GetLoger().Warn("AwardMannaStar fail!")
			return true
		}

		//syncMgr.syncAddStar(IntList{newStarID})
		mannaStarMgr.UpdateMannaStarSeat() //!更新球员可踢位置
		newStar := team.GetStar(newStarID)

		starCount := 0
		if self.ID == 2 {
			starCount = 5
		} else if self.ID == 3 {
			starCount = 7
		}

		//!天赐球员第一位没有初始属性与星际
		if starCount != 0 {
			newStar.SetStarCount(starCount)
			newStar.ChangeStarAttribute(starCount - 1)
		}

		mannaStarMgr.UpdateMannaStarSeat() //!更新球员可踢位置
		mannaStar.Seat2 = 0
		mannaStar.Seat3 = 0

		msg := new(UpdateMannaStarResultMsg)
		msg.MannaStar = *newStar
		msg.Result = 1
		client.SendMsg(msg)

		return true
	}

	updateStar := team.GetStarFromType(mannaStar.ID)

	mannaStarMgr.UpdateMannaStarSeat() //!更新球员可踢位置
	if updateStar.EvolveCount < 4 {
		mannaStar.Seat2 = 0
	}

	if updateStar.EvolveCount < 6 {
		mannaStar.Seat3 = 0
	}

	msg := new(UpdateMannaStarResultMsg)
	msg.MannaStar = *updateStar
	msg.Result = 1
	client.SendMsg(msg)

	syncMgr.SyncObject("mannastar", updateStar)

	return true
}

func (self *UpdateMannaStarMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false {
		return false
	}

	// if self.payAction(client) == false {
	// 	return false
	// }

	if self.doAction(client) == false {
		msg := new(UpdateMannaStarResultMsg)
		msg.Result = 5
		client.SendMsg(msg)
		return false
	}

	return true
}

type QueryMannaStarMsg struct { //!  查询所有原创球员信息
	MsgHead `json:"head"` //! "mannastar", "querymannastar"
}

func (self *QueryMannaStarMsg) GetTypeAndAction() (string, string) {
	return "mannastar", "querymannastar"
}

func (QueryMannaStarMsg) processAction(client IClient) bool {
	team := client.GetTeam()
	starList := MannaStarSlice{}
	mannaStarMgr := team.GetMannaStarMgr()
	for i := MannaStarSeatOne; i <= MannaStarSeatThree; i++ {
		star := mannaStarMgr.GetMannaStarFromSeat(i)
		if star == nil {
			continue
		}

		//! 获取自创球员信息
		starList = append(starList, &star.MannaStarType)
	}

	msg := new(QueryMannaStarResultMsg)
	msg.StarList = starList
	client.SendMsg(msg)

	return true
}

type QueryMannaStarResultMsg struct { //! 查询所有原创球员结果反馈
	MsgHead  `json:"head"`  //! "mannastar", "querymannastarresult"
	StarList MannaStarSlice `json:"starlist"`
}

func (self *QueryMannaStarResultMsg) GetTypeAndAction() (string, string) {
	return "mannastar", "querymannastarresult"
}
