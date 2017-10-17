package football

import ()

const (
	NormalChallangeNumber    = 3 ///普通队伍可挑战的次数
	Vip6ChallangeNumber      = 5 ///vip6队伍可挑战的次数
	ChallangeMatchResetClock = 5 ///重置钟点暂，时硬编码
)

const (
	ChallangeMatchTypePerfect = 1 ///完美艺术
	ChallangeMatchTypeCrazy   = 2 ///疯狂轰炸
	ChallangeMatchTypeDefend  = 3 ///无懈可击
)

///挑战赛
type ChallangeMatchStaticData struct { ///挑战赛类型表
	ID      int    ///挑战赛id
	Type    int    ///挑战赛类型1完美艺术，2疯狂轰炸，3无懈可击
	TeamID  int    ///npc球队ID
	Name    string ///npc球队名
	LV      int    ///开放等级
	Item    int    ///固定掉落类型
	Num     int    ///固定掉落数量
	Item1   int    ///可能奖励道具1
	Item2   int    ///可能奖励道具2
	Item3   int    ///可能奖励道具3
	Item4   int    ///可能奖励道具4
	Item5   int    ///可能奖励道具5
	Num2    int    ///可能奖励道具数量
	Item6   int    ///额外奖励道具
	Num3    int    ///可能奖励道具数量
	Type1   int    ///条件类型
	TypeNum int    ///条件数量
	TypeRes int    ///条件结果
	IGong   string ///我方进攻
	IZU     string ///我方组织
	IFang   string ///我方防守
	EnGong  string ///敌方进攻
	EnZU    string ///敌方组织
	EnFang  string ///敌方防守
	Tips    string
}

///取得队伍的可挑战次数
func GetChallangeNumber(client IClient) int {
	team := client.GetTeam()  ///取得玩家队伍
	vipLevel := team.VipLevel ///取得玩家的VIP等级
	if vipLevel >= 6 {
		return Vip6ChallangeNumber
	} else {
		return NormalChallangeNumber
	}
}

func GetChallangeMatchAttribType(ChallangeMatchType int) int { ///根据挑战赛类型取得对应的可重置数据类型
	res := 0
	switch ChallangeMatchType {
	case ChallangeMatchTypePerfect:
		res = ResetAttribTypeChallangeMatchPerfect
	case ChallangeMatchTypeCrazy:
		res = ResetAttribTypeChallangeMatchCrazy
	case ChallangeMatchTypeDefend:
		res = ResetAttribTypeChallangeMatchDefend
	}
	return res
}

///计算奖励

func CountChallangeMatchAward(userGoalCount int, npcGoalCount int, challangeData *ChallangeMatchStaticData) (awardIDList IntList, awardNumList IntList) {
	awardIDList = IntList{0, 0, 0}
	awardNumList = IntList{0, 0, 0}
	if userGoalCount > npcGoalCount { ///只有胜利了才能得到奖励
		///计算固定奖励
		awardIDList[0] = challangeData.Item
		awardNumList[0] = challangeData.Num
		///计算随机奖励
		randomAward := IntList{}
		if challangeData.Item1 > 0 {
			randomAward = append(randomAward, challangeData.Item1)
		}
		if challangeData.Item2 > 0 {
			randomAward = append(randomAward, challangeData.Item2)
		}
		if challangeData.Item3 > 0 {
			randomAward = append(randomAward, challangeData.Item3)
		}
		if challangeData.Item4 > 0 {
			randomAward = append(randomAward, challangeData.Item4)
		}
		if challangeData.Item5 > 0 {
			randomAward = append(randomAward, challangeData.Item5)
		}
		if randomAward.Len() > 0 { ///有的完全没有额外奖励
			randomTop := randomAward.Len() - 1
			randomAwardNumber := Random(0, randomTop) ///随机取得道具数量
			awardIDList[1] = randomAward[randomAwardNumber]
			awardNumList[1] = challangeData.Num2
		}
		///计算额外奖励
		if challangeData.TypeNum != 0 { ///后面要作为除数，所以检测下避免服务器崩溃
			countNum := 0 ///额外奖励计数
			switch challangeData.Type1 {
			case 1: ///静胜球数
				cGDNum := userGoalCount - npcGoalCount ///净胜球数
				countNum = cGDNum / challangeData.TypeNum
				if countNum > 0 {
					awardIDList[2] = challangeData.Item6
					awardNumList[2] = challangeData.TypeRes * countNum
				}
			case 2: ///进球数
				countNum = userGoalCount / challangeData.TypeNum
				if countNum > 0 {
					awardIDList[2] = challangeData.Item6
					awardNumList[2] = challangeData.TypeRes * countNum
				}
			case 3: ///失球数
				countNum = npcGoalCount / challangeData.TypeNum
				if countNum > 0 {
					awardNumList[0] = Max(0, awardNumList[0]-challangeData.TypeRes*countNum)
				}
			}
		}
	}
	return
}
