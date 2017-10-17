package football

import (
	// "time"
	"fmt"
	"math/rand"
)

///npc球队
type NpcTeamTypeStaticData struct { ///npc球队类型表
	ID                   int    ///记录id
	Name                 string ///npc球队名字
	Icon                 int    ///球队标志
	TeamShirts           int    ///球队球衣
	Formation            int    ///阵型
	Tactical             int    ///战术
	AddBaseExp           int    ///奖励经验基础值
	AddCoin              int    ///奖励球币
	AddTeamExp           int    ///奖励经验池经验
	AddItemType1         int    ///奖励道具类型
	AddItemRate1         int    ///奖励道具1机率
	AddItemMinCount1     int    ///最小奖励数,可以为0
	AddItemMaxCount1     int    ///最大奖励道具数
	AddItemType2         int    ///奖励道具类型2
	AddItemRate2         int    ///奖励道具2机率
	AddItemMinCount2     int    ///最小奖励数2,可以为0
	AddItemMaxCount2     int    ///最大奖励道具数2
	AddStarType          int    ///奖励球星进球队
	AddVolunteerStarType int    ///奖励球星来投球星类型
	AddTicket            int    ///奖励钻石
	AttackScore          int    ///攻击评分
	DefenseScore         int    ///防御评分
	OrganizeScore        int    ///组织评分
	Desc                 string ///描述
	Score                int    ///战斗力评分
}

//! npc球星
type NpcStarTypeStaticData struct { //!npc球星类型表
	ID          int    //!记录ID
	TeamID      int    //!所属球队ID
	Name        string //!npc球星名字
	Icon        int    //!球队标识
	Face        int    //!脸
	Pass        int    //!传球
	Steals      int    //!抢断
	Dribblinig  int    //!盘带
	Sliding     int    //!铲球
	Shooting    int    //!射门
	GoalKeeping int    //!守门
	Body        int    //!身体
	Speed       int    //!速度
	Skill       int    //!技能
	Seat        int    //!位置
	Score       int    //!得分
	Level       int    //!等级
	Grade       int    //!品质颜色
}

const (
	OurTeam   = 1
	EnemyTeam = 2
)

type MatchFlowInfo struct {
	Offensive int `json:"offensive"` //!进攻方  1为我方  2为对方
	Goals     int `json:"goals"`     //!进球数量 为零则无进球
	ScorerID  int `json:"scorerid"`  //!进球球星ID
	//SkillInfo StarSkillList `json:"skillinfo"` //! 本回合技能使用信息
}

type StarUseSkillInfo struct {
	StarID    int     `json:"starid"`
	SkillList IntList `json:"skilllist"`
}

type MatchFlowList []*MatchFlowInfo
type StarSkillList []*StarUseSkillInfo
type NpcStarList []*NpcStarTypeStaticData

func NewAttackTurns(offensive int, goals int, scorerid int /*, skillList StarSkillList*/) *MatchFlowInfo {
	matchFlowInfo := new(MatchFlowInfo)
	matchFlowInfo.Offensive = offensive
	matchFlowInfo.Goals = goals
	matchFlowInfo.ScorerID = scorerid
	//	matchFlowInfo.SkillInfo = skillList
	//fmt.Println(skillList)
	// for i := 0; i < len(skillList); i++ {
	// 	for j := 0; j < len(skillList[i].SkillList); j++ {
	// 		fmt.Printf("starID: %d  skillList: %d  \r\n", skillList[i].StarID, skillList[i].SkillList[j])
	// 	}
	// }

	//fmt.Println("Look there~~~~~")
	return matchFlowInfo
}

func (self *NpcTeamTypeStaticData) CalcLoseResult() (int, int) { ///随机生成一个输了的比分
	userGoalCount := Random(0, 7)
	npcGoalCount := Random(0, 7)
	if userGoalCount == npcGoalCount {
		npcGoalCount += 1
	}
	if userGoalCount > npcGoalCount {
		SwapInt(&userGoalCount, &npcGoalCount) ///如果比分大则交换位置
	}
	return userGoalCount, npcGoalCount
}

///模拟比赛过程
func CalcGoalCountTurns(attackTurns int, attackScore float32, defenseScore float32, overNum float32,
	attackTurns2 int, attackScore2 float32, defenseScore2 float32, overNum2 float32, team *Team,
	tarTeamFormation int, npcTeamID int, organizeSum float32) (int, int, int, int, MatchFlowList) { ///计算进球数

	attackgoalCount := 0
	attackgoalRate := int(attackScore*100/(attackScore+defenseScore) + (overNum * 100)) ///由数值策划林之冠要求更改

	denfensegoalCount := 0
	denfensegoalRate := int(attackScore2*100/(attackScore2+defenseScore2) + (overNum2 * 100)) ///由数值策划林之冠要求更改

	matchFlowLst := MatchFlowList{}

	attackgoalRate = Min(attackgoalRate, 100) ///不得超出100

	//fmt.Printf("team : %v", team)

	currentFormation := team.GetFormationMgr().GetFormation(team.FormationID)
	starIDList := currentFormation.GetStarIDList()

	//fmt.Printf("starIDList : %v", starIDList)

	staticDataMgr := GetServer().GetStaticDataMgr()
	skillMgr := team.GetSkillMgr()

	//!技能属性修正值
	var attackTemp float32
	var defenseTemp float32
	var overNumTemp float32
	var attackTemp2 float32
	var defenseTemp2 float32
	var overNumTemp2 float32
	attackTurnTemp := 0
	attackTurnTemp2 := 0
	turnsTemp := 0

	//fmt.Printf("attackgoalRate:%d  \r\n  denfensegoalRate:%d  \r\n", attackgoalRate, denfensegoalRate)

	turns := Max(attackTurns, attackTurns2)

	attackTurnTemp = attackTurns
	attackTurnTemp2 = attackTurns2 //! 保存无技能标准回合数
	turnsTemp = turns

	for i := 1; i <= turnsTemp; i++ { ///计算玩家进球数
		fmt.Println("turnsTemp: ", turnsTemp)
		//!偏差值清零
		attackTemp = attackScore
		defenseTemp = defenseScore
		overNumTemp = overNum
		attackTemp2 = attackScore2
		defenseTemp2 = defenseScore2
		overNumTemp2 = overNum2
		attackTurnTemp = attackTurns
		attackTurnTemp2 = attackTurns2
		turnsTemp = turns

		//!本回合使用技能表
		skillUseList := StarSkillList{}

		for j := 0; j < starIDList.Len(); j++ {
			userSkillList := skillMgr.SkillTime(starIDList[j], attackgoalCount, denfensegoalCount, tarTeamFormation)
			if userSkillList.Len() <= 0 {
				continue
			}

			fmt.Printf("Skill Open!! Goal:  %d  :  %d \r\n", attackgoalCount, denfensegoalCount)

			skilldata := new(StarUseSkillInfo)
			skilldata.StarID = starIDList[j]
			skilldata.SkillList = userSkillList
			skillUseList = append(skillUseList, skilldata)
			//fmt.Printf("skillList: %v          ", skillUseList)
			//fmt.Printf("skilldata.skillList: %v \r\n", skilldata.SkillList)

			for n := 0; n < len(userSkillList); n++ {
				fmt.Printf("starID: %d Skill: %d ", starIDList[j], userSkillList[n])

			}

			fmt.Println("Turn:", i)

		}

		//! 使用技能
		for b := 0; b < len(skillUseList); b++ {

			for n := 0; n < len(skillUseList[b].SkillList); n++ {
				//	fmt.Printf("b: %d   n: %d  len(skillUseList): %d ", b, n, len(skillUseList))
				restraintFormation, attackCalc, defenceCalc, organizeCalc := skillMgr.SkillEffectNPC(npcTeamID, skillUseList[b].SkillList[n], skillUseList[b].StarID, team)
				//	fmt.Printf("skillUseList[b].skillList[n]: %d   ", skillUseList[b].SkillList[n])
				skillType := staticDataMgr.GetSkillType(skillUseList[b].SkillList[n])
				//	fmt.Println(attackCalc)
				//	fmt.Println(defenceCalc)
				if skillType.Tartype == 1 {
					attackTemp += attackCalc
					defenseTemp += defenceCalc

					overNumTemp = overNum * (1.0 - restraintFormation)

					attackTurnTemp = attackTurns + int(organizeCalc*10.0/float32(organizeSum))

					//! 修正攻击次数不得超过9
					attackTurnTemp = Min(attackTurnTemp, 9)
				} else if skillType.Tartype == 2 {
					attackTemp2 += attackCalc
					defenseTemp2 += defenceCalc

					overNumTemp2 = overNum2 * (1.0 - restraintFormation)

					attackTurnTemp2 = attackTurns2 + int(organizeCalc*10.0/float32(organizeSum))

					//! 额外攻击次数
					attackTurnTemp2 = Min(attackTurnTemp2, 9)
				}

				//! 额外攻击回合
				turnsTemp = Max(attackTurnTemp, attackTurnTemp2)
				//fmt.Printf("restraintFormation:%v  attackCalc:%v  defenceCalc:%v", restraintFormation, attackCalc, defenceCalc)
			}

			//! 改变属性后重新计算胜率
			attackgoalRate = int(attackTemp*100/(attackTemp+defenseTemp) + (overNumTemp * 100))
			denfensegoalRate = int(attackTemp2*100/(attackTemp2+defenseTemp2) + (overNumTemp2 * 100))

			//! 进球率不得为负
			attackgoalRate = Max(attackgoalRate, 0)
			denfensegoalRate = Max(denfensegoalRate, 0)

			fmt.Printf("attackgoalRate:%d  \r\n  denfensegoalRate:%d  \r\n", attackgoalRate, denfensegoalRate)
		}

		if i <= attackTurnTemp {
			//!主动方回合
			randRate := rand.Intn(100)
			if randRate <= attackgoalRate {
				attackgoalCount++ ///进一球

				matchFlowInfo := NewAttackTurns(OurTeam, 1, 0)
				matchFlowLst = append(matchFlowLst, matchFlowInfo)

			} else {
				//!无进球
				matchFlowInfo := NewAttackTurns(OurTeam, 0, 0)
				matchFlowLst = append(matchFlowLst, matchFlowInfo)
			}
		}

		if i <= attackTurnTemp2 {
			randRate := rand.Intn(100)
			if randRate <= denfensegoalRate {
				denfensegoalCount++ ///进一球
				matchFlowInfo := NewAttackTurns(EnemyTeam, 1, 0)
				matchFlowLst = append(matchFlowLst, matchFlowInfo)
			} else {
				matchFlowInfo := NewAttackTurns(EnemyTeam, 0, 0)
				matchFlowLst = append(matchFlowLst, matchFlowInfo)
			}

		}

	}

	return attackgoalCount, attackgoalRate, denfensegoalCount, denfensegoalRate, matchFlowLst
}

//根据结果模拟过程
func SimulationMatch(ourTeamCount int, enemyTeamCount int) MatchFlowList {
	attackgoalCount := 0
	attackgoalRate := 50

	denfensegoalCount := 0
	denfensegoalRate := 50

	matchFlowLst := MatchFlowList{}

	const loopTimes = 2000

	//	turns := ourTeamCount + enemyTeamCount

	for i := 1; i <= loopTimes; i++ { ///计算玩家进球数

		if attackgoalCount < ourTeamCount {
			//!主动方回合
			randRate := rand.Intn(100)
			if randRate <= attackgoalRate {
				attackgoalCount++ ///进一球

				matchFlowInfo := NewAttackTurns(OurTeam, 1, 0)
				matchFlowLst = append(matchFlowLst, matchFlowInfo)

			} else {
				//!无进球
				matchFlowInfo := NewAttackTurns(OurTeam, 0, 0)
				matchFlowLst = append(matchFlowLst, matchFlowInfo)
			}
		}

		if denfensegoalCount < enemyTeamCount {
			randRate := rand.Intn(100)
			if randRate <= denfensegoalRate {
				denfensegoalCount++ ///进一球
				matchFlowInfo := NewAttackTurns(EnemyTeam, 1, 0)
				matchFlowLst = append(matchFlowLst, matchFlowInfo)
			} else {
				matchFlowInfo := NewAttackTurns(EnemyTeam, 0, 0)
				matchFlowLst = append(matchFlowLst, matchFlowInfo)
			}

		}

		if ourTeamCount == attackgoalCount && enemyTeamCount == denfensegoalCount {
			break
		}

	}

	return matchFlowLst
}

///计算比赛结果,返回x,y表示几比几x为玩家进球数,y为npc进球数
func (self *NpcTeamTypeStaticData) CalcMatchResult(userTeam *Team) (int, int) {
	const overcomeAwardTurns = 2 ///克制奖励进攻次数
	///得到玩家球队比赛三围
	attackUserScoreCalc, defenseUserScoreCalc, organizeUserScoreCalc := userTeam.CalcScore()
	///得到npc球队比赛三围
	attackNpcScoreCalc := float32(self.AttackScore)
	defenseNpcScoreCalc := float32(self.DefenseScore)
	organizeNpcScoreCalc := float32(self.OrganizeScore)
	///计算玩家球队进攻次数
	attackUserTurns := int(organizeUserScoreCalc * 10 / (organizeUserScoreCalc + organizeNpcScoreCalc))
	///计算npc球队进攻次数
	attackNpcTurns := int(organizeNpcScoreCalc * 10 / (organizeUserScoreCalc + organizeNpcScoreCalc))
	formOverNum, tacticOverNum := float32(0), float32(0) ///攻方阵型克制系数与战术克制系数
	formationLevel := float32(userTeam.FormationLevel)
	///计算双方阵形相克对进攻次数的加成
	formationUser := userTeam.GetCurrentFormObject()
	if formationUser.IsOvercome(self.Formation) {
		attackUserTurns += overcomeAwardTurns
		formOverNum = -1 * (0.1 + formationLevel*0.005)
	} else if formationUser.IsBeOvercome(self.Formation) {
		attackNpcTurns += overcomeAwardTurns
		formOverNum = 0.1 + formationLevel*0.005 ///被克
	}
	///战术克制
	if formationUser.IsOverTactic(self.Tactical) {
		tacticOverNum = 0.1 + formationLevel*0.005
	} else if formationUser.IsBeOverTactic(self.Tactical) {
		tacticOverNum = -1 * (0.1 + formationLevel*0.005) //被克
	}
	userGoalRate, npcGoalRate := 0, 0
	userGoalCount, npcGoalCount := 0, 0
	userScore := float32(userTeam.Score) //float32(math.Pow(float64(attackUserScoreCalc)*float64(defenseUserScoreCalc)*math.Pow(float64(organizeUserScoreCalc), 1.2), 0.33333333)) ///玩家战斗力
	npcScore := float32(self.Score)      //float32(math.Pow(float64(attackNpcScoreCalc)*float64(defenseNpcScoreCalc)*math.Pow(float64(organizeNpcScoreCalc), 1.2), 0.33333333))     ///npc战斗力
	///改动 3：PVE战斗计算前加入另一个判定条件“NPC战斗力低于1000，且低于我方战斗力直接判负。
	//	if userScore > (npcScore * 0.9) {
	//		npcGoalCount, userGoalCount = self.CalcLoseResult() ///直接判胜
	//	} else

	//	if userScore <= (npcScore * 0.9) {
	//userGoalCount, npcGoalCount = self.CalcLoseResult() ///直接判负,在01,12,02之间随机
	//		userGoalCount, npcGoalCount = RandMatchResult([][2]int{{0, 1}, {1, 2}, {0, 2}})
	//	} else { ///随机生成比较结果
	///计算玩家进球数
	//userGoalCount, userGoalRate = CalcGoalCount(attackUserTurns, attackUserScoreCalc, defenseNpcScoreCalc, tacticOverNum)
	///计算npc进球数
	//npcGoalCount, npcGoalRate = CalcGoalCount(attackNpcTurns, attackNpcScoreCalc, defenseUserScoreCalc, formOverNum)
	//	}

	//修改: 新回合制模拟比赛
	organizeSum := organizeUserScoreCalc + organizeNpcScoreCalc
	userGoalCount, userGoalRate, npcGoalCount, npcGoalRate, matchList := CalcGoalCountTurns(attackUserTurns, attackUserScoreCalc, defenseNpcScoreCalc, tacticOverNum,
		attackNpcTurns, attackNpcScoreCalc, defenseUserScoreCalc, formOverNum, userTeam, self.Formation, self.ID, organizeSum)

	//	GetServer().GetLoger().Debug("userGoalCount = %v, userGoalRate = %v, npcGoalCount = %v, npcGoalRate = %v, List: %v", test1, test2, test3, test4, test5)
	// fmt.Printf("userGoalCount = %v, userGoalRate = %v, npcGoalCount = %v, npcGoalRate = %v", test1, test2, test3, test4)

	//改动: 当NPC战力小于1050,且正常流程我方输了,则判我方胜利
	//再次改动: 战力高于NPC  必胜
	//npcScore < 1050 &&
	if userGoalCount <= npcGoalCount && userScore >= npcScore {
		if userGoalCount == npcGoalCount {
			userGoalCount += 1 //平局 我方进球加一
			matchFlowInfo := NewAttackTurns(OurTeam, 1, 0)
			matchList = append(matchList, matchFlowInfo)
		} else {
			//SwapInt(&userGoalCount, &npcGoalCount) //否则交换比分,在21,31,32随机
			userGoalCount, npcGoalCount = RandMatchResult([][2]int{{2, 1}, {3, 1}, {3, 2}})
			matchList = SimulationMatch(userGoalCount, npcGoalCount)
		}
	} else if userGoalCount > npcGoalCount && userScore <= npcScore*0.9 {
		userGoalCount, npcGoalCount = RandMatchResult([][2]int{{0, 1}, {1, 2}, {0, 2}})
		matchList = SimulationMatch(userGoalCount, npcGoalCount)
	}

	client := GetServer().userMgr.GetClientByTeamID(userTeam.ID)
	SendMatchFlowMsg(client, matchList)

	GetServer().GetLoger().Print("%s{分%d 攻%f 防%f 组%f 攻次%d 球率%d}---%s{分%d 攻%f 防%f 组%f 攻次%d 球率%d}\n",
		userTeam.GetInfo().Name, userTeam.Score, attackUserScoreCalc, defenseUserScoreCalc, organizeUserScoreCalc, attackUserTurns, userGoalRate,
		self.Name, self.AttackScore+self.DefenseScore+self.OrganizeScore, attackNpcScoreCalc, defenseNpcScoreCalc, organizeNpcScoreCalc, attackNpcTurns, npcGoalRate)
	return userGoalCount, npcGoalCount
}
