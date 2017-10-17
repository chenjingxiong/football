package football

const (
	mgrTypeSkillMgr       = 1  ///技能管理器
	mgrTypeLevelMgr       = 2  ///关卡管理器
	mgrTypeTaskMgr        = 3  ///任务管理器
	mgrTypeFormationMgr   = 4  ///阵形管理器
	mgrTypeResetAttribMgr = 5  ///可重置属性管理器
	mgrTypeItemMgr        = 6  ///游戏道具管理器
	mgrTypeProcessMgr     = 7  ///任务处理管理器
	mgrTypeStarFateMgr    = 8  ///球员缘管理器
	mgrTypeArenaMgr       = 9  ///竞技场管理器
	mgrTypeVipShopMgr     = 10 ///商城管理器
	mgrTypeMailMgr        = 11 ///邮件管理器
	mgrTypeActivityMgr    = 12 ///活动管理器
	mgrTypeMannaStarMgr   = 13 //!天赐球员管理器
	mgrTypeAtlasMgr       = 14 //!球员图鉴管理器
)

///游戏逻辑管理器
type IGameMgr interface {
	SaveInfo()                   ///保存数据到数据库
	GetType() int                ///得到管理器类型
	SetSyncMgr(syncMgr *SyncMgr) ///设置同步管理器
	SetTeam(team *Team)          ///设置球队对象
	onInit()                     ///初始化函数
}

type GameMgr struct {
	//teamID  int
	team    *Team
	syncMgr *SyncMgr
}

func (self *GameMgr) SetSyncMgr(syncMgr *SyncMgr) {
	self.syncMgr = syncMgr
}

func (self *GameMgr) SetTeam(team *Team) {
	self.team = team
}

func (self *GameMgr) onInit() {
	k := 1
	k++
}
