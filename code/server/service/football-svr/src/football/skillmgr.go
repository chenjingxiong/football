package football

import (
	"fmt"
)

///技能管理器
type ISkillMgr interface {
	//Init(teamID int) bool
	GetSkillList() *SkillList                                                                                             ///得到技能列表
	GetSkillInfoList() SkillInfoList                                                                                      ///得到所有技能信息列表
	GetSkill(skillID int) *Skill                                                                                          ///得到技能对象
	GetStarSkillSlice(starID int) SkillSlice                                                                              ///得到指定球员所装备的技能对象列表
	AddSkill(starID int, skillType int)                                                                                   ///为球员增加技能
	SkillTime(starID int, teamGoal int, tarTeamGoal int, tarTeamFormation int) IntList                                    //! 使用技能
	SkillEffectNPC(npcTeamID int, skillType int, useSkillStar int, targetTeam *Team) (float32, float32, float32, float32) //! PVE中使用技能
	RemoveStarSkill(skillID int)
}

type SkillSlice []*Skill
type SkillList map[int]*Skill
type SkillMgr struct {
	GameMgr
	skillList SkillList ///技能列表
}

func (self *SkillMgr) SaveInfo() { ///保存数据
	for _, v := range self.skillList {
		v.Save()
	}
}

func (self *SkillMgr) GetType() int { ///得到管理器类型
	return mgrTypeSkillMgr
}

func (self *SkillMgr) GetStarSkillSlice(starID int) SkillSlice { ///得到指定球员所装备的技能对象列表
	skillSlice := SkillSlice{}
	for _, v := range self.skillList {
		skillInfo := v.GetInfo()
		if skillInfo.StarID == starID {
			skillSlice = append(skillSlice, v)
		}
	}
	return skillSlice
}

func (self *SkillMgr) GetSkillInfoList() SkillInfoList { ///外界提供技能信息列表由此函数填充
	skillInfoList := SkillInfoList{}
	for _, v := range self.skillList {
		skillInfoList = append(skillInfoList, *v.GetInfo())
	}
	return skillInfoList
}

func (self *SkillMgr) GetSkill(skillID int) *Skill { ///得到技能对象
	return self.skillList[skillID]
}

func (self *SkillMgr) GetSkillList() *SkillList {
	return &self.skillList
}

//!增加技能
func (self *SkillMgr) AddSkill(starID int, skillType int) {
	skillInsertSql := fmt.Sprintf("insert into %s (teamid, starid, type) value(%d, %d, %d)", tableSkill, self.team.GetID(), starID, skillType)
	lastID, _ := GetServer().GetDynamicDB().Exec(skillInsertSql)
	if lastID == 0 {
		GetServer().GetLoger().Warn("%s is fail", skillInsertSql)
		return
	}

	//!若插入成功则加入内存
	skillInfo := new(SkillInfo)
	skillInfo.ID = lastID
	skillInfo.StarID = starID
	skillInfo.TeamID = self.team.GetID()
	skillInfo.Type = skillType
	self.skillList[lastID] = NewSkill(skillInfo)

}

//! 初始化技能管理器
func NewSkillMgr(teamID int) IGameMgr {
	skillMgr := new(SkillMgr)
	skillMgr.skillList = make(SkillList)
	//skillMgr.teamID = teamID ///存放自己的球队id	levelID := levelMgr.AddLevel(self.LevelType, 0)
	skillListQuery := fmt.Sprintf("select * from %s where teamid=%d limit 3000", tableSkill, teamID)
	skillInfo := new(SkillInfo)
	skillInfoList := GetServer().GetDynamicDB().fetchAllRows(skillListQuery, skillInfo)
	for i := range skillInfoList {
		skillInfo = skillInfoList[i].(*SkillInfo)
		skillMgr.skillList[skillInfo.ID] = NewSkill(skillInfo)
	}
	return skillMgr
}

//! 判断比分触发
func (self *SkillMgr) scoreOpenSkill(sort int, teamGoal int, tarTeamGoal int) bool {
	skillIsOpen := false
	switch sort {
	case 1:
		skillIsOpen = teamGoal > tarTeamGoal //!比分领先则触发
	case 2:
		skillIsOpen = teamGoal == tarTeamGoal //!比分持平则触发
	case 3:
		skillIsOpen = teamGoal < tarTeamGoal //!比分落后则触发
	}

	//fmt.Printf("teamGoal : %d   tarTeamGoal : %d  skillIsOpen = %v", teamGoal, tarTeamGoal, skillIsOpen)
	return skillIsOpen
}

//! 判断当前球员能使用的技能
func (self *SkillMgr) SkillTime(starID int, teamGoal int, tarTeamGoal int, tarTeamFormation int) IntList {

	//! 得到该球星已学技能
	skillLst := self.GetStarSkillSlice(starID)
	if len(skillLst) == 0 {
		return nil
	}

	//fmt.Println(skillLst)
	starUseSkill := IntList{}

	//! 根据技能类型得到技能信息
	for i := 0; i < len(skillLst); i++ {
		skillType := GetServer().GetStaticDataMgr().GetSkillType(skillLst[i].Type)
		if skillType == nil {
			break
		}

		switch skillType.Type {
		case 1: //!根据比分触发
			if self.scoreOpenSkill(skillType.Sort, teamGoal, tarTeamGoal) == true {
				starUseSkill = append(starUseSkill, skillType.ID)
			}
		case 2: //!根据我方阵型触发
			formationID := self.team.GetCurrentFormation()
			currentFormation := self.team.GetFormationMgr().GetFormation(formationID)
			if currentFormation.Type == skillType.Sort {
				starUseSkill = append(starUseSkill, skillType.ID)
			}

		case 3: //!根据敌方阵型触发
			if tarTeamFormation == skillType.Sort {
				starUseSkill = append(starUseSkill, skillType.ID)
			}

		case 4: //!被动永久触发
			starUseSkill = append(starUseSkill, skillType.ID)
		}

	}

	//fmt.Println(starUseSkill)
	return starUseSkill
}

//! 得到NPC球星信息
func (self *SkillMgr) getNpcStarList(npcTeamID int) NpcStarList {
	if npcTeamID == 0 {
		return nil
	}
	staticDataMgr := GetServer().GetStaticDataMgr()
	npcStarList := NpcStarList{}
	for i := 1; i <= 11; i++ {
		starType := staticDataMgr.GetNpcTeamStarType(npcTeamID*100 + i)
		if starType != nil {
			npcStarList = append(npcStarList, starType)
		}
	}

	return npcStarList
}

//! 得到技能生效对象 _ NPC
func (self *SkillMgr) skillEffectTarget_NPC(tarSort int, npcStarList NpcStarList) NpcStarList {
	targetList := NpcStarList{}
	if tarSort <= 13 {
		for i := 0; i < len(npcStarList); i++ {
			if tarSort == npcStarList[i].Seat {
				targetList = append(targetList, npcStarList[i])
			}

		}
	} else {
		switch tarSort {
		case 14: //! 全部中锋,边锋,影锋
			for i := 0; i < len(npcStarList); i++ {
				if npcStarList[i].Seat <= 4 {
					targetList = append(targetList, npcStarList[i])
				}

			}
		case 15: //! 全部边前卫、前腰、中前卫、后腰
			for i := 0; i < len(npcStarList); i++ {
				if npcStarList[i].Seat <= 9 && npcStarList[i].Seat > 4 {
					targetList = append(targetList, npcStarList[i])
				}

			}
		case 16: //! 全部边后卫、中后卫、门将
			for i := 0; i < len(npcStarList); i++ {
				if npcStarList[i].Seat <= 13 && npcStarList[i].Seat > 10 {
					targetList = append(targetList, npcStarList[i])
				}

			}
		case 17: //! 全部球员
			for i := 0; i < len(npcStarList); i++ {
				targetList = append(targetList, npcStarList[i])
			}
		}
	}

	return targetList
}

//! 得到技能生效对象  team
func (self *SkillMgr) skillEffectTarget_Team(tarSort int, targetTeam *Team) (IntList, IntList) {
	formationMgr := targetTeam.GetFormationMgr()
	formationInfo := formationMgr.GetFormation(targetTeam.FormationID)

	starIDList := IntList{}
	starSeatList := IntList{}
	starList := IntList{formationInfo.Pos1, formationInfo.Pos2, formationInfo.Pos3, formationInfo.Pos4, formationInfo.Pos5,
		formationInfo.Pos6, formationInfo.Pos7, formationInfo.Pos8, formationInfo.Pos9, formationInfo.Pos10, formationInfo.Pos11}
	seatTypeList := formationInfo.GetSeatTypeList()

	if tarSort < 14 {
		for i := 0; i < seatTypeList.Len(); i++ {
			if seatTypeList[i] == tarSort {
				starIDList = append(starIDList, starList[i])
				starSeatList = append(starSeatList, seatTypeList[i])
			}
		}
	} else {
		switch tarSort {
		case 14: //! 全部中锋,边锋,影锋
			for i := 0; i < seatTypeList.Len(); i++ {
				if seatTypeList[i] <= 4 {
					starIDList = append(starIDList, starList[i])
					starSeatList = append(starSeatList, seatTypeList[i])
				}

			}
		case 15: //! 全部边前卫、前腰、中前卫、后腰
			for i := 0; i < seatTypeList.Len(); i++ {
				if seatTypeList[i] <= 9 && seatTypeList[i] > 4 {
					starIDList = append(starIDList, starList[i])
					starSeatList = append(starSeatList, seatTypeList[i])
				}

			}
		case 16: //! 全部边后卫、中后卫、门将
			for i := 0; i < seatTypeList.Len(); i++ {
				if seatTypeList[i] <= 13 && seatTypeList[i] > 10 {
					starIDList = append(starIDList, starList[i])
					starSeatList = append(starSeatList, seatTypeList[i])
				}

			}
		case 17: //! 全部球员
			return starList, seatTypeList
		}
	}

	return starIDList, starSeatList
}

//! 技能---改变球员属性  返回对球队三围的影响绝对值
func (self *SkillMgr) SkillAttr(attrType int, power int, pass int, steals int,
	dribbling int, sliding int, shooting int, goalkeeping int, seatType int) (float32, float32, float32) {
	staticDataMgr := GetServer().GetStaticDataMgr()
	seatTypeStaticData := staticDataMgr.GetSeatType(seatType)

	//attackScore := 0
	var attackScore float32
	var defenseScore float32
	var organizeScore float32
	var starScore float32

	switch attrType {
	case 1: //!传球
		passCalc := float32(pass) * float32(power) / 100.0
		passCalc = passCalc * float32(seatTypeStaticData.ScorePassRate)
		passCalc /= 100
		starScore = float32(TakeAccuracy(float32(passCalc), 0))
	case 2: //!抢断
		stealsCalc := float32(pass) * float32(power) / 100.0
		stealsCalc = stealsCalc * float32(seatTypeStaticData.ScorePassRate)
		stealsCalc /= 100
		starScore = TakeAccuracy(float32(stealsCalc), 0)

	case 3: //!盘带
		dribblingCalc := float32(dribbling) * float32(power) / 100.0
		dribblingCalc = dribblingCalc * float32(seatTypeStaticData.ScorePassRate)
		dribblingCalc /= 100
		starScore = TakeAccuracy(float32(dribblingCalc), 0)

	case 4: //!铲球
		slidingCalc := float32(sliding) * float32(power) / 100.0
		slidingCalc = slidingCalc * float32(seatTypeStaticData.ScorePassRate)
		slidingCalc /= 100
		starScore = TakeAccuracy(float32(slidingCalc), 0)

	case 5: //!射门
		shootCalc := float32(shooting) * float32(power) / 100.0
		shootCalc = shootCalc * float32(seatTypeStaticData.ScorePassRate)
		shootCalc /= 100
		starScore = TakeAccuracy(float32(shootCalc), 0)

	case 6: //!守门
		goalkeepingCalc := float32(goalkeeping) * float32(power) / 100.0
		goalkeepingCalc = goalkeepingCalc * float32(seatTypeStaticData.ScorePassRate)
		goalkeepingCalc /= 100
		starScore = TakeAccuracy(float32(goalkeepingCalc), 0)

	}

	attackScore = starScore * float32(seatTypeStaticData.AttackRate)
	defenseScore = starScore * float32(seatTypeStaticData.DefenseRate)
	organizeScore = starScore * float32(seatTypeStaticData.OrganizeRate)
	return attackScore, defenseScore, organizeScore
}

//! 技能对攻防数值影响
func (self *SkillMgr) SkillEffectNPC(npcTeamID int, skillType int, useSkillStar int, targetTeam *Team) (float32, float32, float32, float32) {

	npcStarLst := self.getNpcStarList(npcTeamID)
	if len(npcStarLst) != 11 && npcTeamID != 0 {
		return 0, 0, 0, 0
	}

	staticDataMgr := GetServer().GetStaticDataMgr()
	skillTypeInfo := staticDataMgr.GetSkillType(skillType)

	var attackscore float32
	var defensescore float32
	var organizescore float32
	var restraintFormation float32

	if skillTypeInfo.Tartype == 1 {
		//!技能触发后作用于本方球员
		effectStarList := IntList{}
		if skillTypeInfo.Tarsort != 0 {
			effectStarList, _ = self.skillEffectTarget_Team(skillTypeInfo.Tarsort, self.team)
		} else {
			effectStarList = append(effectStarList, useSkillStar)
		}

		formationID := self.team.GetCurrentFormation()
		formationInfo := self.team.GetFormationMgr().GetFormation(formationID)

		//!得到技能效果
		switch skillTypeInfo.Func {
		case 1: //增益Buff
			for i := 0; i < effectStarList.Len(); i++ {
				star := self.team.GetStar(effectStarList[i])

				if star.IsMannaStar == 1 {
					mannaStarMgr := self.team.GetMannaStarMgr()
					starType := mannaStarMgr.GetMannaStar(star.Type)
					attackScoreCalc, defenseScoreCalc, organizeScoreCalc := self.SkillAttr(skillTypeInfo.Attr, skillTypeInfo.Power,
						starType.Pass+star.PassTalentAdd,
						starType.Steals+star.StealsTalentAdd,
						starType.Dribbling+star.DribblingTalentAdd,
						starType.Sliding+star.SlidingTalentAdd,
						starType.Shooting+star.ShootingTalentAdd,
						starType.GoalKeeping+star.GoalKeepingTalentAdd,
						formationInfo.Type)

					attackscore += attackScoreCalc
					defensescore += defenseScoreCalc
					organizescore += organizeScoreCalc

				} else {
					starType := staticDataMgr.GetStarType(star.Type)

					attackScoreCalc, defenseScoreCalc, organizeScoreCalc := self.SkillAttr(skillTypeInfo.Attr, skillTypeInfo.Power,
						starType.Pass+star.PassTalentAdd,
						starType.Steals+star.StealsTalentAdd,
						starType.Dribbling+star.DribblingTalentAdd,
						starType.Sliding+star.SlidingTalentAdd,
						starType.Shooting+star.ShootingTalentAdd,
						starType.GoalKeeping+star.GoalKeepingTalentAdd,
						formationInfo.Type)

					attackscore += attackScoreCalc
					defensescore += defenseScoreCalc
					organizescore += organizeScoreCalc
				}

			}

		case 2: //增益Buff
			for i := 0; i < effectStarList.Len(); i++ {
				star := self.team.GetStar(effectStarList[i])
				if star.IsMannaStar == 1 {
					starType := self.team.GetMannaStarMgr().GetMannaStar(star.Type)
					attackScoreCalc, defenseScoreCalc, organizeScoreCalc := self.SkillAttr(skillTypeInfo.Attr, skillTypeInfo.Power,
						starType.Pass+star.PassTalentAdd,
						starType.Steals+star.StealsTalentAdd,
						starType.Dribbling+star.DribblingTalentAdd,
						starType.Sliding+star.SlidingTalentAdd,
						starType.Shooting+star.ShootingTalentAdd,
						starType.GoalKeeping+star.GoalKeepingTalentAdd,
						formationInfo.Type)

					attackscore += attackScoreCalc
					defensescore += defenseScoreCalc
					organizescore += organizeScoreCalc
				} else {
					starType := staticDataMgr.GetStarType(star.Type)
					attackScoreCalc, defenseScoreCalc, organizeScoreCalc := self.SkillAttr(skillTypeInfo.Attr, skillTypeInfo.Power,
						starType.Pass+star.PassTalentAdd,
						starType.Steals+star.StealsTalentAdd,
						starType.Dribbling+star.DribblingTalentAdd,
						starType.Sliding+star.SlidingTalentAdd,
						starType.Shooting+star.ShootingTalentAdd,
						starType.GoalKeeping+star.GoalKeepingTalentAdd,
						formationInfo.Type)

					attackscore += attackScoreCalc
					defensescore += defenseScoreCalc
					organizescore += organizeScoreCalc
				}

			}

			attackscore *= -1
			defensescore *= -1
			organizescore = -1
		case 3: //战术阵型克制效果削弱
			restraintFormation = float32(skillTypeInfo.Power) / 100.0
		}
	} else {
		//!技能触发后作用于对方球员
		if npcTeamID != 0 {
			//!NPC
			effectStarList := self.skillEffectTarget_NPC(skillTypeInfo.Tarsort, npcStarLst)
			switch skillTypeInfo.Func {
			case 1:
				for i := 0; i < len(effectStarList); i++ {
					curStar := effectStarList[i]
					attackScoreCalc, defenseScoreCalc, organizeScoreCalc := self.SkillAttr(skillTypeInfo.Attr, skillTypeInfo.Power,
						curStar.Pass, curStar.Steals, curStar.Dribblinig, curStar.Sliding, curStar.Shooting, curStar.GoalKeeping, curStar.Seat)
					attackscore += attackScoreCalc
					defensescore += defenseScoreCalc
					organizescore += organizeScoreCalc
				}

			case 2:
				for i := 0; i < len(effectStarList); i++ {
					curStar := effectStarList[i]
					attackScoreCalc, defenseScoreCalc, organizeScoreCalc := self.SkillAttr(skillTypeInfo.Attr, skillTypeInfo.Power,
						curStar.Pass, curStar.Steals, curStar.Dribblinig, curStar.Sliding, curStar.Shooting, curStar.GoalKeeping, curStar.Seat)
					attackscore += attackScoreCalc
					defensescore += defenseScoreCalc
					organizescore += organizeScoreCalc
				}

				attackscore *= -1
				defensescore *= -1
				organizescore = -1

			}
		} else {
			//!玩家
			effectStarList, starSeatList := self.skillEffectTarget_Team(skillTypeInfo.Tarsort, targetTeam)
			switch skillTypeInfo.Func {
			case 1:
				for i := 0; i < len(effectStarList); i++ {
					curStar := targetTeam.GetStar(effectStarList[i])
					if curStar.IsMannaStar == 1 {
						starType := targetTeam.GetMannaStarMgr().GetMannaStar(curStar.Type)
						attackScoreCalc, defenseScoreCalc, organizeScoreCalc := self.SkillAttr(skillTypeInfo.Attr, skillTypeInfo.Power,
							curStar.PassTalentAdd+starType.Pass,
							curStar.StealsTalentAdd+starType.Steals,
							curStar.DribblingTalentAdd+starType.Dribbling,
							curStar.SlidingTalentAdd+starType.Sliding,
							curStar.ShootingTalentAdd+starType.Shooting,
							curStar.GoalKeepingTalentAdd+starType.GoalKeeping,
							starSeatList[i])
						attackscore += attackScoreCalc
						defensescore += defenseScoreCalc
						organizescore += organizeScoreCalc

					} else {
						starType := staticDataMgr.GetStarType(curStar.Type)
						attackScoreCalc, defenseScoreCalc, organizeScoreCalc := self.SkillAttr(skillTypeInfo.Attr, skillTypeInfo.Power,
							curStar.PassTalentAdd+starType.Pass,
							curStar.StealsTalentAdd+starType.Steals,
							curStar.DribblingTalentAdd+starType.Dribbling,
							curStar.SlidingTalentAdd+starType.Sliding,
							curStar.ShootingTalentAdd+starType.Shooting,
							curStar.GoalKeepingTalentAdd+starType.GoalKeeping,
							starSeatList[i])
						attackscore += attackScoreCalc
						defensescore += defenseScoreCalc
						organizescore += organizeScoreCalc
					}

				}

			case 2:
				for i := 0; i < len(effectStarList); i++ {
					curStar := targetTeam.GetStar(effectStarList[i])

					if curStar.IsMannaStar == 1 {
						starType := targetTeam.GetMannaStarMgr().GetMannaStar(curStar.Type)
						attackScoreCalc, defenseScoreCalc, organizeScoreCalc := self.SkillAttr(skillTypeInfo.Attr, skillTypeInfo.Power,
							curStar.PassTalentAdd+starType.Pass,
							curStar.StealsTalentAdd+starType.Steals,
							curStar.DribblingTalentAdd+starType.Dribbling,
							curStar.SlidingTalentAdd+starType.Sliding,
							curStar.ShootingTalentAdd+starType.Shooting,
							curStar.GoalKeepingTalentAdd+starType.GoalKeeping,
							starSeatList[i])
						attackscore += attackScoreCalc
						defensescore += defenseScoreCalc
						organizescore += organizeScoreCalc

					} else {
						starType := staticDataMgr.GetStarType(curStar.Type)
						attackScoreCalc, defenseScoreCalc, organizeScoreCalc := self.SkillAttr(skillTypeInfo.Attr, skillTypeInfo.Power,
							curStar.PassTalentAdd+starType.Pass,
							curStar.StealsTalentAdd+starType.Steals,
							curStar.DribblingTalentAdd+starType.Dribbling,
							curStar.SlidingTalentAdd+starType.Sliding,
							curStar.ShootingTalentAdd+starType.Shooting,
							curStar.GoalKeepingTalentAdd+starType.GoalKeeping,
							starSeatList[i])
						attackscore += attackScoreCalc
						defensescore += defenseScoreCalc
						organizescore += organizeScoreCalc
					}
				}

				attackscore *= -1
				defensescore *= -1
				organizescore = -1

			}
		}

	}

	return restraintFormation, attackscore, defensescore, organizescore
}

func (self *SkillMgr) RemoveStarSkill(skillID int) {
	sql := fmt.Sprintf("delete from %s where id = %d", tableSkill, skillID)
	dynamicDBMgr := GetServer().GetDynamicDB()
	dynamicDBMgr.Exec(sql)

	delete(self.skillList, skillID)
}
