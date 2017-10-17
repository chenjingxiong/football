package football

///球员位置类型数据
type SeatTypeStaticData struct { ///位置类型,对应动态表中的dy_seattype
	ID                   int    ///任务类型id
	Name                 string ///位置名
	ScorePassRate        int    ///评分传球加成
	ScoreStealsRate      int    ///评分抢断加成
	ScoreDribblingRate   int    ///评分盘带加成
	ScoreSlidingRate     int    ///评分铲球加成
	ScoreShootingRate    int    ///评分射门加成
	ScoreGoalKeepingRate int    ///评分守门加成
	ScoreBodyRate        int    ///评分身体加成
	ScoreSpeedRate       int    ///评分速度加成
	AttackRate           int    ///攻击力加成
	DefenseRate          int    ///防御力加成
	OrganizeRate         int    ///组织力加成
}
