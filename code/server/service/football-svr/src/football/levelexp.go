package football

const (
	levelExpTypeStarLevel          = 1 ///球员等级
	levelExpTypeStarEducationLevel = 2 ///球员培养加点等级
	levelExpTypeEquipMerge         = 3 ///装备融合等级
	levelExpTypeFormation          = 6 ///球队阵形等级
	levelExpTypeStarEvolve         = 7 ///球员突破等级
	levelExpTypeTeamLevel          = 8 ///球队等级

)

type LevelExpStaticData struct { ///服务器升级数据静态配置表
	ID              int ///记录id
	Type            int ///配置数据类型
	Level           int ///等级
	NeedExp         int ///所需经验上限,球员突破时重用为培养上限
	PayCoin         int ///升级所需游戏币
	NeedLevel       int ///所需等级
	NeedItemType1   int ///所需道具类型1,0表示不需要
	NeedItemCount1  int ///所需道具数量1,0表示不需要
	NeedItemType2   int ///所需道具类型2,0表示不需要
	NeedItemCount2  int ///所需道具数量2,0表示不需要
	NeedTeamLevel   int ///所需球队等级,0表示不需要
	NeedEvolveCount int ///所需星级
}
