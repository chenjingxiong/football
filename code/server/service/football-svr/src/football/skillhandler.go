package football

import (
	"fmt"
)

//{"head":{"seq":2,"type":"skill","action":"studyskill","time":1409735838},"starid":9999,"skilltype":9999}
// {"head":{"action":"studyskill","seq":21,"skilltype":100005,"starid":17067,"time":1409735644,"type":"skill"}}
type StudySkillMsg struct { //!球星学习技能信息
	MsgHead   `json:"head"` //!"skill", "studyskill"
	StarID    int           `json:"starid"`    //!球星ID
	SkillType int           `json:"skilltype"` //!技能类型
}

type StudySkillResultMsg struct { //!球星学习技能信息
	MsgHead `json:"head"` //!"skill", "studyskillresult"
	Result  int           `json:"result"` //! 0为正常 1 有技能正在学习  2 已学习   3 星级不够
}

func (self *StudySkillMsg) GetTypeAndAction() (string, string) {
	return "skill", "studyskill"
}

func (self *StudySkillResultMsg) GetTypeAndAction() (string, string) {
	return "skill", "studyskillresult"
}

func (self *StudySkillMsg) checkAction(client IClient) bool {
	team := client.GetTeam()
	star := team.GetStar(self.StarID)
	loger := GetServer().GetLoger()
	staticDataMgr := GetServer().GetStaticDataMgr()

	fmt.Println("starid: ", self.StarID)
	fmt.Println("star: ", star)
	if loger.CheckFail("star != nil", star != nil, star, nil) { //! 球队必须有此球员
		return false
	}

	//!获取球员可学习技能
	isCanStudy := false
	starInfo := star.GetInfo()
	if star.IsMannaStar == 1 {
		mannaStarMgr := team.GetMannaStarMgr()
		starType := mannaStarMgr.GetMannaStar(starInfo.Type)

		if loger.CheckFail("starType != nil", starType != nil, starType, nil) { //! 静态表需有该球员类型
			return false
		}

		starSkillLst := IntList{starType.Skill1, starType.Skill2, starType.Skill3, starType.Skill4}

		for i := 0; i < starSkillLst.Len(); i++ {
			if starSkillLst[i] == self.SkillType {
				isCanStudy = true
				break
			}
		}
	} else {

		starType := staticDataMgr.GetStarType(starInfo.Type)

		if loger.CheckFail("starType != nil", starType != nil, starType, nil) { //! 静态表需有该球员类型
			return false
		}

		starSkillLst := IntList{starType.Skill1, starType.Skill2, starType.Skill3, starType.Skill4}

		for i := 0; i < starSkillLst.Len(); i++ {
			if starSkillLst[i] == self.SkillType {
				isCanStudy = true
				break
			}
		}

	}

	if loger.CheckFail("isCanStudy == true", isCanStudy == true, isCanStudy, true) { //! 该球星必须能够学习此技能
		return false
	}

	skillType := GetServer().GetStaticDataMgr().GetSkillType(self.SkillType)
	if loger.CheckFail("skillType != nil", skillType != nil, skillType, nil) { //! 技能ID必须存在于技能表
		return false
	}

	skillMgr := team.GetSkillMgr()
	skillLst := skillMgr.GetStarSkillSlice(self.StarID)
	isHasStudy := false
	for i := 0; i < len(skillLst); i++ {
		if skillLst[i].Type == self.SkillType && skillLst[i].TeamID == team.GetID() {
			fmt.Println("skillLst[i].Type", skillLst[i].Type)
			fmt.Println("self.SkillType", self.SkillType)
			isHasStudy = true
			break
		}
	}

	if loger.CheckFail("isHasStudy == false", isHasStudy == false, isHasStudy, false) { //! 已经学过该技能
		msg := new(StudySkillResultMsg)
		msg.Result = 2
		client.SendMsg(msg)

		return false
	}

	//!球星一次只能学习一种技能
	attrMgr := team.GetResetAttribMgr()
	resetAttrSkill := attrMgr.GetResetAttrib(ResetAttribTypeSkillStudyInfo)
	if resetAttrSkill == nil {
		//! 不存在则创建默认
		attrMgr.AddResetAttrib(ResetAttribTypeSkillStudyInfo, -1, IntList{0, 0, 0})
	} else {
		if loger.CheckFail("resetAttrSkill.Value3 == 0", resetAttrSkill.Value3 == 0, resetAttrSkill.Value3, 0) {
			msg := new(StudySkillResultMsg)
			msg.Result = 1
			client.SendMsg(msg)

			return false //! 有其他球员在训练
		}
	}

	//! 判断星级是否达到学习技能条件
	if loger.CheckFail("starInfo.EvolveCount >= skillType.Open", starInfo.EvolveCount >= skillType.Open, starInfo.EvolveCount, skillType.Open) {

		msg := new(StudySkillResultMsg)
		msg.Result = 3
		client.SendMsg(msg)
		return false
	}

	return true
}

//! 学习技能暂无代价
// func (self *StudySkillMsg) payAction(client IClient) bool {
// 	return true
// }

func (self *StudySkillMsg) doAction(client IClient) bool {

	team := client.GetTeam()
	attrMgr := team.GetResetAttribMgr()
	skillType := GetServer().GetStaticDataMgr().GetSkillType(self.SkillType)
	skillAttrInfo := attrMgr.GetResetAttrib(ResetAttribTypeSkillStudyInfo)
	skillAttrInfo.Value1 = self.StarID
	skillAttrInfo.Value2 = self.SkillType
	skillAttrInfo.Value3 = 1
	skillAttrInfo.ResetTime = Now() + skillType.Time //! 单位:秒

	msg := new(StudySkillResultMsg)
	msg.Result = 0
	client.SendMsg(msg)

	// msg := NewQuerySkillInfoMsg(self.StarID, self.SkillType, Now()+skillType.Time)
	// client.SendMsg(msg)

	return true
}

func (self *StudySkillMsg) processAction(client IClient) bool {

	if self.checkAction(client) == false {
		return false
	}

	// if self.payAction(client) == false {
	// 	return false
	// }

	if self.doAction(client) == false {
		return false
	}

	return true
}

type SkillState struct {
	StarID     int `json:"starid"`     //! 球员ID
	SkillType  int `json:"skilltype"`  //! 技能种类
	SkillState int `json:"skillstate"` //! 技能状态 0 = 未学习  1 = 学习中  2 = 学习完
	Time       int `json:"time"`       //! CD时间
}

type QuerySkillStudyInfoMsg struct {
	MsgHead `json:"head"` //! "skill" "querystudyinfo"
}

type QuerySkillStudyInfoResultMsg struct {
	MsgHead   `json:"head"` //! "skill" "querystudyinforesult"
	StarSkill []*SkillState `json:"starid"`
}

func (self *QuerySkillStudyInfoResultMsg) GetTypeAndAction() (string, string) {
	return "skill", "querystudyinforesult"
}

func NewQuerySkillInfoMsg(starSkill []*SkillState) *QuerySkillStudyInfoResultMsg {
	msg := new(QuerySkillStudyInfoResultMsg)
	msg.StarSkill = starSkill

	return msg
}

func (self *QuerySkillStudyInfoMsg) GetTypeAndAction() (string, string) {
	return "skill", "querystudyinfo"
}

func (self *QuerySkillStudyInfoMsg) skillIsStudy(team *Team, skillType int, starID int) bool {
	skillMgr := team.GetSkillMgr()
	isStudy := false
	starSkillLst := skillMgr.GetStarSkillSlice(starID)
	for i := 0; i < len(starSkillLst); i++ {
		v := starSkillLst[i]
		if v.Type == skillType {
			isStudy = true
			break
		}
	}
	return isStudy
}

func (self *QuerySkillStudyInfoMsg) checkSkillState(team *Team, skillType int, starID int) (int, int) {
	attrMgr := team.GetResetAttribMgr()
	skillAttrMgr := attrMgr.GetResetAttrib(ResetAttribTypeSkillStudyInfo)

	state := 0
	time := 0

	//!判断技能是否已学习
	isStudy := self.skillIsStudy(team, skillType, starID)
	if isStudy == true {
		state = 2
		time = 0

		return state, time
	}

	//!判断技能是否在学习中
	if skillAttrMgr.Value1 == starID && skillAttrMgr.Value2 == skillType && skillAttrMgr.Value3 != 0 {
		state = 1
		time = skillAttrMgr.ResetTime - Now()

		return state, time
	}

	//!否则为未学习技能
	return state, time
}

func (self *QuerySkillStudyInfoMsg) processAction(client IClient) bool {
	team := client.GetTeam()
	attrMgr := team.GetResetAttribMgr()
	skillAttrMgr := attrMgr.GetResetAttrib(ResetAttribTypeSkillStudyInfo)
	if skillAttrMgr == nil {
		//! 不存在则创建默认
		skillAttrMgr = attrMgr.AddResetAttrib(ResetAttribTypeSkillStudyInfo, -1, IntList{0, 0, 0})
	}

	starSkill := []*SkillState{}
	staticDataMgr := GetServer().GetStaticDataMgr()
	allStarList := team.GetAllStarList()
	for i := 0; i < allStarList.Len(); i++ {
		star := team.GetStar(allStarList[i])
		if star.IsMannaStar == 1 {
			starType := team.GetMannaStarMgr().GetMannaStar(star.Type)

			if starType.Skill1 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill1
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill1, star.ID)
				starSkill = append(starSkill, node)
			}

			if starType.Skill2 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill2
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill2, star.ID)
				starSkill = append(starSkill, node)
			}

			if starType.Skill3 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill3
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill3, star.ID)

				starSkill = append(starSkill, node)
			}

			if starType.Skill4 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill4
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill4, star.ID)

				starSkill = append(starSkill, node)
			}
		} else {
			starType := staticDataMgr.GetStarType(star.Type)

			if starType.Skill1 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill1
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill1, star.ID)
				starSkill = append(starSkill, node)
			}

			if starType.Skill2 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill2
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill2, star.ID)
				starSkill = append(starSkill, node)
			}

			if starType.Skill3 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill3
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill3, star.ID)

				starSkill = append(starSkill, node)
			}

			if starType.Skill4 != 0 {
				node := new(SkillState)
				node.StarID = star.ID
				node.SkillType = starType.Skill4
				node.SkillState, node.Time = self.checkSkillState(team, starType.Skill4, star.ID)

				starSkill = append(starSkill, node)
			}
		}

	}

	//!调试输出
	// for i := 0; i < len(starSkill); i++ {
	// 	v := starSkill[i]
	// 	fmt.Println("starID: ", v.StarID)
	// 	fmt.Println("skillType: ", v.SkillType)
	// 	fmt.Println("skillState: ", v.SkillState)
	// 	fmt.Println("time: ", v.Time)

	// }

	msg := NewQuerySkillInfoMsg(starSkill)
	client.SendMsg(msg)

	return true
}

type QueryStarSkillInfoMsg struct {
	MsgHead `json:"head"` //! "skill" "querystarskillinfo"
	StarID  int           `json:"starid"`
}

type QueryStarSkillInfoResultMsg struct {
	MsgHead  `json:"head"` //! "skill" "querystarskillinforesult"
	SkillLst IntList       `json:"skilllst"` //!技能列表
}

func (self *QueryStarSkillInfoMsg) GetTypeAndAction() (string, string) {
	return "skill", "querystarskillinfo"
}

func (self *QueryStarSkillInfoResultMsg) GetTypeAndAction() (string, string) {
	return "skill", "querystarskillinforesult"
}

func (self *QueryStarSkillInfoMsg) checkAction(client IClient) bool {
	loger := GetServer().GetLoger()
	team := client.GetTeam()
	star := team.GetStar(self.StarID)
	if loger.CheckFail("star != nil", star != nil, star, nil) {
		return false
	}
	return true
}

func (self *QueryStarSkillInfoMsg) doAction(client IClient) bool {
	team := client.GetTeam()
	skillMgr := team.GetSkillMgr()
	skill := skillMgr.GetStarSkillSlice(self.StarID)
	skillLst := IntList{}
	for i := 0; i < len(skill); i++ {
		skillLst = append(skillLst, skill[i].Type)
	}

	self.sendQuerySkillResult(client, skillLst)
	return true
}

func (self *QueryStarSkillInfoMsg) sendQuerySkillResult(client IClient, skillLst IntList) {
	msg := new(QueryStarSkillInfoResultMsg)
	msg.SkillLst = skillLst
	client.SendMsg(msg)
}

func (self *QueryStarSkillInfoMsg) processAction(client IClient) bool {
	if self.checkAction(client) == false {
		return false
	}

	if self.doAction(client) == false {
		return false
	}

	return true
}
