package football

///这里将注册所有消息处理器
func (self *MsgDispatch) initMsgHandlerList() {
	//self.addMsgHandleToList(new(LoginHandler))      ///登录消息处理器
	//self.addMsgHandleToList(new(TeamHandler))       ///球队消息处理器
	//self.addMsgHandleToList(new(StarSpyHandler))    ///球探消息处理器
	//self.addMsgHandleToList(new(StarCenterHandler)) ///球员中心消息处理器
	//self.addMsgHandleToList(new(ItemHandler))      ///道具消息处理器
	//self.addMsgHandleToList(new(SkillHandler))     ///技能消息处理器 技能消息被废弃
	//self.addMsgHandleToList(new(FormationHandler)) ///球队阵形消息处理器
	//self.addMsgHandleToList(new(StarHandler))  ///球员消息处理器
	//self.addMsgHandleToList(new(LevelHandler)) ///关卡消息处理器

	///新的处理流程
	//	self.addMsgRegistry(new(ClientLoginMsg)) ///处理登录消息

	self.addMsgRegistry(new(RegisteringAccount)) //! 注册账号消息

	self.addMsgRegistry(new(TeamCreateMsg)) ///处理球队创建消息
	///self.addMsgRegistry(new(TeamAddTrainCellMsg))       ///处理增加训练位消息
	///self.addMsgRegistry(new(TeamQueryStarTrainListMsg)) ///查询球员训练列表消息
	///self.addMsgRegistry(new(TeamTrainStarMsg))          ///训练一个球员消息
	///self.addMsgRegistry(new(TeamAbortTrainStarMsg))     ///放弃训练一个球员消息

	self.addMsgRegistry(new(StarSpyDiscoverMsg)) ///转会中心发掘球员消息

	//self.addMsgRegistry(new(QueryStarCenterMemberListMsg)) ///查询转会中心球员列表消息 (删除)
	//self.addMsgRegistry(new(StarCenterTransferMsg))        ///签约转会中心球员消息
	self.addMsgRegistry(new(QueryVolunteerInfoMsg)) ///查询球星来投信息列表
	self.addMsgRegistry(new(VolunteerSignMsg))      ///签约球星来投系统中的球员
	self.addMsgRegistry(new(GetStarLobbyAwardMsg))  ///领取球员游说中的奖励

	self.addMsgRegistry(new(QueryItemListMsg)) ////查询球队所拥有的所有道具列表消息
	self.addMsgRegistry(new(EquipItemMsg))     ///装备球员一个道具消息
	self.addMsgRegistry(new(MergeItemMsg))     ///请求融合一个道具消息
	self.addMsgRegistry(new(ItemEvolveMsg))    ///请求对一个道具进行升星操作

	self.addMsgRegistry(new(FormationSetCurrentMsg)) ///设置球队阵形消息
	self.addMsgRegistry(new(FormationUplevelMsg))    ///球队阵形升级消息
	self.addMsgRegistry(new(FormationChangeStarMsg)) ///球队阵形换人消息

	self.addMsgRegistry(new(StarEducationMsg)) ///球员培养消息
	self.addMsgRegistry(new(StarEvolveMsg))    ///球员升星消息
	self.addMsgRegistry(new(StarSackMsg))      ///解雇球员

	self.addMsgRegistry(new(QueryLeagueInfoMsg)) ///查询联赛地图信息
	self.addMsgRegistry(new(QueryLevelListMsg))  ///查询关卡列表消息
	self.addMsgRegistry(new(PassLevelMsg))       ///请求通过关卡消息
	self.addMsgRegistry(new(SkipLevelMsg))       ///请求跳过比赛消息

	self.addMsgRegistry(new(ChatMsg))             ///处理聊天消息
	self.addMsgRegistry(new(QueryStoreMsg))       ///注册仓库查询消息
	self.addMsgRegistry(new(ItemSellMsg))         ///出售道具消息
	self.addMsgRegistry(new(QueryTrainMatchMsg))  ///查询训练赛当前信息
	self.addMsgRegistry(new(RefeshTrainMatchMsg)) ///客户端请求刷新训练赛训练项目列表
	self.addMsgRegistry(new(AwardTrainMatchMsg))  ///客户端请求领取训练赛积分奖励
	self.addMsgRegistry(new(PlayTrainMatchMsg))   ///客户端请求领取训练赛积分奖励
	self.addMsgRegistry(new(QueryArenaMatchMsg))  ///客户端请求查询竞技场信息
	self.addMsgRegistry(new(QueryTeamInfoMsg))    ///客户端请求查询球队详细信息
	self.addMsgRegistry(new(PlayArenaMatchMsg))   ///客户端请求打竞技场比赛信息
	self.addMsgRegistry(new(AcceptArenaAwardMsg)) ///客户端请求领取竞技场奖励

	self.addMsgRegistry(new(GetStarExpPoolMsg)) ///客户端请求经验池分配奖励
	self.addMsgRegistry(new(StarSpyOperateMsg)) ///客户端球员抽取处理操作

	self.addMsgRegistry(new(VipShopCommodityQueryMsg)) ///客户端查询商城商品消息
	self.addMsgRegistry(new(VipShopCommodityBuyMsg))   ///客户端购买商品消息
	self.addMsgRegistry(new(ItemUseMsg))               ///客户端使用物品消息

	self.addMsgRegistry(new(MailQueryMsg))        ///邮件查询
	self.addMsgRegistry(new(MailAwardReceiveMsg)) ///邮件领奖
	self.addMsgRegistry(new(MailReadMsg))         ///阅读邮件
	self.addMsgRegistry(new(MailDeleteMsg))       ///邮件删除

	self.addMsgRegistry(new(QueryShoppintActionPointInfoMsg)) ///购买行动点界面信息查询
	self.addMsgRegistry(new(BuyActionPointMsg))               ///购买行动点消息

	self.addMsgRegistry(new(QueryGoldFingerInfoMsg)) ///查询金手指信息
	self.addMsgRegistry(new(UseGoldFingerMsg))       ///使用金手指
	self.addMsgRegistry(new(ClientOperationMsg))     ///客户端操作串消息
	self.addMsgRegistry(new(ActivitCodeMsg))         ///激活码消息

	self.addMsgRegistry(new(GetPassLevelAwardMsg)) ///领取推图奖励消息
	//self.addMsgRegistry(new(VipShopBuyDiamondMsg))           ///领取推图奖励消息
	self.addMsgRegistry(new(VipShopQueryBuyDiamondCountMsg)) ///查询已购买套餐次数
	self.addMsgRegistry(new(VipAccpetDayAwardMsg))           ///领取vip每日礼包
	self.addMsgRegistry(new(VipQueryDayAwardMsg))            ///查询vip每日礼包领奖状态
	//	self.addMsgRegistry(new(VipBuyMonthCardMsg))             ///购买vip月卡
	self.addMsgRegistry(new(VipQueryMonthCardMsg))       ///查询vip月卡状态
	self.addMsgRegistry(new(VipAccpetMonthCardAwardMsg)) ///领取vip月卡每日礼包

	self.addMsgRegistry(new(TaskQueryDayTaskMsg))       ///查询日常任务信息
	self.addMsgRegistry(new(TaskAccpetDayTaskAwardMsg)) ///领取日常任务奖励

	self.addMsgRegistry(new(TaskQueryTenDrawStarsMsg)) ///请求查询十连抽信息
	self.addMsgRegistry(new(TaskTenDrawStarsMsg))      ///请求进行十连抽信息
	self.addMsgRegistry(new(TaskTakeTenDrawStars))     ///请求获得十连抽产生的球员
	self.addMsgRegistry(new(SDKLoginMsg))              ///客户端sdk请求登录消息

	self.addMsgRegistry(new(QueryCalabashInfoMsg)) ///请求葫芦娃的领取状态
	self.addMsgRegistry(new(AccpetCalabashMsg))    /// 请求领取葫芦娃消息

	self.addMsgRegistry(new(ItemCombineMsg)) /// 请求合成道具消息

	self.addMsgRegistry(new(QueryMaxStartTimesMsg)) ///请求冠军之路三星后剩余可挑战次数

	self.addMsgRegistry(new(QueryChallengeTimesMsg)) ///请求挑战赛日剩余挑战次数
	self.addMsgRegistry(new(QueryChallengeFightMsg)) ///请求挑战赛挑战

	self.addMsgRegistry(new(GetPowerMsg))   //! 请求领取体力
	self.addMsgRegistry(new(AddStarPosMsg)) //! 扩充球员位

	self.addMsgRegistry(new(StudySkillMsg))          //! 客户端请求球员学习技能消息
	self.addMsgRegistry(new(QuerySkillStudyInfoMsg)) //! 查询球员学习技能状态消息
	self.addMsgRegistry(new(QueryStarSkillInfoMsg))  //! 查询球员当前拥有技能

	self.addMsgRegistry(new(StarFateSignMsg)) //! 缘分系统签约球星消息

	self.addMsgRegistry(new(UpdateMannaStarMsg)) //!更新自创球员消息
	self.addMsgRegistry(new(QueryMannaStarMsg))  //!查询天赐球员信息

	self.addMsgRegistry(new(UpdateTeamIconAndShirtMsg)) //! 更新球队队徽与队服

	self.addMsgRegistry(new(MergeCardMsg)) //! 融合球员碎片消息

	self.addMsgRegistry(new(QueryAtlasInfoMsg))    //! 查询图鉴消息
	self.addMsgRegistry(new(ReceiveAtlasAwardMsg)) //! 领取图鉴奖励消息
}
