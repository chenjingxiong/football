package football

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
	//	"time"
)

const (
	///动态表
	tableStar           = "dy_star"           ///球员表
	tableFormation      = "dy_formation"      ///阵型表
	tableStarSpy        = "dy_starspy"        ///球探表
	tableStarCenter     = "dy_starcenter"     ///球员中心表
	tableTeam           = "dy_team"           ///球队表
	tableProcessCenter  = "dy_processcenter"  ///球队处理中心表
	tableResetAttrib    = "dy_resetattrib"    ///可重置属性表
	tableItem           = "dy_item"           ///道具对象表
	tableSkill          = "dy_skill"          ///对象表
	tableLevel          = "dy_level"          ///关卡对象表
	tableTask           = "dy_task"           ///任务对象表
	tableVipShop        = "dy_vipshop"        //商城限购记录表
	tableArena          = "dy_arena"          ///竞技场信息表
	tableActivity       = "dy_activity"       ///球队活动信息表
	tableMail           = "dy_mail"           ///球队邮件列表
	tableSystemMail     = "dy_systemmail"     ///系统邮件表
	tableActivationCode = "dy_activationcode" ///激活码表
	tableLeagueAward    = "dy_leagueaward"    ///推图奖励表
	tableAccount        = "dy_account"        ///帐号表
	tablePayOrder       = "dy_payorder"       ///支付日志表
	tableMannaStar      = "dy_mannastar"      //!天赐球员信息表
	tableAtlas          = "dy_atlas"          //!球星图鉴

	///静态表
	tableConfig          = "st_config"           ///服务器配置表
	tableStarType        = "st_startype"         ///球员类型表
	tableDrawGroup       = "st_drawgroup"        ///球员类型表
	tableItemType        = "st_itemtype"         ///道具类型表
	tableLevelExp        = "st_levelexp"         ///道具类型表
	tableFormationType   = "st_formationtype"    ///球队阵形类型表
	tableTacticType      = "st_tactictype"       ///阵形战术类型表
	tableLevelType       = "st_leveltype"        ///关卡类型表
	tableTaskType        = "st_tasktype"         ///任务类型表
	tableNpcTeamType     = "st_npcteam"          ///npc球队类型表
	tableSeatType        = "st_seattype"         ///球员位置类型表
	tableSkillType       = "st_skilltype"        ///技能类型表
	tableStarFateType    = "st_fatetype"         ///球员缘类型表
	tableStarLobbyType   = "st_starlobby"        ///球员游说类型表
	tableTrainAward      = "st_trainaward"       ///训练赛奖励表
	tableVipShopType     = "st_vipshop"          ///商城信息表
	tableArenaType       = "st_arenatype"        ///联赛类型表
	tableVipPrivilege    = "st_vipprivilege"     ///VIP特权类型表
	tableActivityType    = "st_activitytype"     ///活动类型表
	tableActivityAward   = "st_activityaward"    ///活动类型表
	tableActionType      = "st_action"           ///动作类型表
	tableLeagueAwardType = "st_leaguematch"      ///推图奖励表
	tableAwardType       = "st_activitcodeaward" ///激活码奖励表
	tableMoney           = "st_money"            ///充值配置表
	tableSDKList         = "st_sdklist"          ///sdk平台列表
	tableServerList      = "st_serverlist"       ///服务器列表配置
	tableChallangeMatch  = "st_challangematch"   ///挑战赛表
	tableNpcStarType     = "st_npcstar"          ///npc球星类型表
	//	tableUpdateVersion   = "st_updateversion"    ///客户端版本更新表

	///记录表
	tableRecordCreate = "record_create" ///玩家创建表
	tableRecordAction = "record_action" ///玩家行为表
	tableRecordPay    = "record_pay"    ///玩家消费表
	tableRecordOnline = "record_online" ///玩家在线表
)

const MaxSQLCostMS = 100 ///SQL最慢容忍毫秒数,默认为100毫秒

type DBServer struct { ///数据库组件
	db     *sql.DB ///直接从mysqldb派生
	dbName string
}

func (self *DBServer) deamonWriteGo() { ///后台写入协程

}

func (self *DBServer) GetDBName() string {
	return self.dbName
}

func (self *DBServer) parseDBName(dbDNS *string) {
	dbName := *dbDNS
	beginIndex := strings.Index(dbName, "/")
	endIndex := strings.Index(dbName, "?")
	self.dbName = dbName[beginIndex+1 : endIndex]
}

func (self *DBServer) Init(dbDNS *string) {
	GetServer().GetLoger().Info("connecting db please wait... %s", *dbDNS)
	db, err := sql.Open("mysql", *dbDNS)

	if err != nil {
		db.Close()
		GetServer().GetLoger().Fatal("newDBServer Open fail! err:%v dns:%s", err, *dbDNS)
	}
	err = db.Ping() // Open doesn't open a connection. Validate DSN data:
	if err != nil {
		db.Close()
		GetServer().GetLoger().Fatal("newDBServer Open Ping  fail!  err:%v dns:%s", err, *dbDNS)
	}
	db.SetMaxIdleConns(GetServer().GetConfig().MaxIdleDBConns) ///设置连接池库存连接数
	self.db = db
	self.parseDBName(dbDNS)
	GetServer().GetLoger().Info("connect db success!%s", *dbDNS)
}

func (self *DBServer) Exec(query string, args ...interface{}) (int, int) {
	sql := fmt.Sprintf(query, args...)
	beginMS := NowMS()
	result, err := self.db.Exec(sql)
	endMS := NowMS()
	spendMS := endMS - beginMS ///计算sql执行毫秒成本
	if spendMS >= MaxSQLCostMS {
		GetServer().GetLoger().Warn("DBServer Exec slow! cost:%d>%d ms sql:%s", spendMS, MaxSQLCostMS, sql)
	}
	//GetServer().GetLoger().Debug("DBServer Exec SQL:%s", query)
	if err != nil {
		GetServer().GetLoger().Warn("DBServer Exec fail! err:%v sql:%s", err, sql)
		return 0, 0
	}
	lastInsertID := int64(0)
	lastInsertID, err = result.LastInsertId()
	if err != nil {
		GetServer().GetLoger().Warn("DBServer LastInsertId fail! err:%v sql:%s", err, sql)
	}
	rowsAffected := int64(0)
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		GetServer().GetLoger().Warn("DBServer RowsAffected fail! err:%v sql:%s", err, sql)
	}
	return int(lastInsertID), int(rowsAffected)
}

func (self *DBServer) fetchOneRow(query string, struc interface{}) bool {
	ok := false
	beginMS := NowMS()
	rows, err := self.db.Query(query)
	endMS := NowMS()
	spendMS := endMS - beginMS ///计算sql执行毫秒成本
	if spendMS >= MaxSQLCostMS {
		GetServer().GetLoger().Warn("DBServer fetchOneRow slow! cost:%d>%d ms sql:%s", spendMS, MaxSQLCostMS, query)
	}
	//rows, err := self.db.Query(query)
	if nil == rows || err != nil {
		GetServer().GetLoger().Warn("DBServer fetchOneRow fail! err:%v sql:%s", err, query)
		return false
	}
	s := reflect.ValueOf(struc).Elem()            ///得到结构体的反射类型
	numFieldStruct := s.NumField()                ///得到结构体的字段数
	onerow := make([]interface{}, numFieldStruct) ///生成一行记录接口
	for i := 0; i < numFieldStruct; i++ {
		onerow[i] = s.Field(i).Addr().Interface() ///将结构体与接口绑定
	}
	for rows.Next() {
		err = rows.Scan(onerow...) ///生成一条记录///得到结构体的反射类
		if err != nil {
			GetServer().GetLoger().Warn("DBServer fetchAllRowsPtr rows.Scan fail! err:%v sql:%s", err, query)
		}
		ok = true
		break ///只取一行记录
	}
	rows.Close()
	return ok
}

func (self *DBServer) fetchAllRows(query string, struc interface{}) []interface{} {
	beginMS := NowMS()
	rows, err := self.db.Query(query)
	endMS := NowMS()
	spendMS := endMS - beginMS ///计算sql执行毫秒成本
	if spendMS >= MaxSQLCostMS {
		GetServer().GetLoger().Warn("DBServer fetchAllRows slow! cost:%d>%d ms sql:%s", spendMS, MaxSQLCostMS, query)
	}
	//rows, err := self.db.Query(query)
	if nil == rows || err != nil {
		GetServer().GetLoger().Warn("DBServer fetchAllRowsPtr fail! err:%v sql:%s", err, query)
		return nil
	}
	s := reflect.ValueOf(struc).Elem()            ///得到结构体的反射类型
	numFieldStruct := s.NumField()                ///得到结构体的字段数
	result := make([]interface{}, 0)              ///生成记录集列表
	onerow := make([]interface{}, numFieldStruct) ///生成一行记录接口
	for i := 0; i < numFieldStruct; i++ {
		onerow[i] = s.Field(i).Addr().Interface() ///将结构体与接口绑定
	}
	for rows.Next() {
		err = rows.Scan(onerow...) ///生成一条记录
		if err != nil {
			GetServer().GetLoger().Warn("DBServer fetchAllRowsPtr rows.Scan fail! err:%v sql:%s", err, query)
			break
		}
		newObj := reflect.New(reflect.TypeOf(struc).Elem()).Elem()
		newObj.Set(s)
		result = append(result, newObj.Addr().Interface()) ///将记录放入记录集中
	}
	rows.Close()
	return result
}

func (self *DBServer) Query(query string, args ...interface{}) *sql.Rows {
	sql := fmt.Sprintf(query, args...)
	beginMS := NowMS()
	rows, err := self.db.Query(sql)
	endMS := NowMS()
	spendMS := endMS - beginMS ///计算sql执行毫秒成本
	if spendMS >= MaxSQLCostMS {
		GetServer().GetLoger().Warn("DBServer Query slow! cost:%d>%d ms sql:%s", spendMS, MaxSQLCostMS, sql)
	}
	if err != nil {
		GetServer().GetLoger().Warn("DBServer Query fail! err:%v sql:%s", err, sql)
		return nil
	}
	return rows
}
