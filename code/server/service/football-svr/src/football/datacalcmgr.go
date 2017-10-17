package football

import (
	"fmt"
	"runtime/debug"
	"time"
)

type DataCalcMgr struct { ///数据统计管理器
	calcTimeDelaySecs int64 ///每十分钟统计一次数据
}

func (self *DataCalcMgr) calcServerOnline(startTime string, endTime string, serverID int) { ///统计服务器在线数据
	loger().Info("DataCalcMgr::calcServerOnline begin at %s", time.Now().Format("2006-01-02 15:04:05"))
	targetTableName := "mmo2d_recordljzm.onlinedata" ///后台对应的表
	todayTime := time.Now().Format("2006/01/02")     ///取得今天日期串
	calcDataSQL := fmt.Sprintf(`replace %s select 0 as id,%d as sid,'%s' as 'date',a.sdkname,a.loginnum,0 as maxonline,0 as aveonline,IFNULL(b.num,0) as  curonline from 
	(select sdkname,count(distinct teamid) as loginnum from %s  
	where type=1 and param=1  and maketime between '%s' and '%s' group by sdkname) as a left join
	(select sdkname,count(id) as num from %s group by sdkname) as b on a.sdkname=b.sdkname`,
		targetTableName, serverID, endTime, tableRecordAction, startTime, endTime, tableRecordOnline)
	_, rowsAffected := GetServer().GetRecordDB().Exec(calcDataSQL) ///写入服务器在线人数

	calcDataSQL = fmt.Sprintf(`replace %s select 0 as id,%d as sid,'%s' as 'date',a.sdkname,b.loginnum,0 as maxonline,0 as aveonline,IFNULL(a.num,0) as  curonline from 
	(select sdkname,count(id) as num from %s group by sdkname) as a left join
	(select sdkname,count(distinct teamid) as loginnum from %s  
	where type=1 and param=1  and maketime between '%s' and '%s' group by sdkname) as b	 on a.sdkname=b.sdkname`,
		targetTableName, serverID, endTime, tableRecordOnline, tableRecordAction, startTime, endTime)
	_, rowsAffected = GetServer().GetRecordDB().Exec(calcDataSQL) ///写入服务器在线人数

	//	if lastInsertItemID > 0 {
	calcDataSQL = fmt.Sprintf(`update %s as o,(select sdkname,max(curonline) as maxonline,avg(curonline) as aveonline from
		%s where date between '%s 00:00:00' and '%s 23:59:59' group by sdkname) as t set o.maxonline=t.maxonline,o.aveonline=t.aveonline
			where o.date='%s' and o.sdkname=t.sdkname`, targetTableName, targetTableName, todayTime, todayTime, endTime)
	GetServer().GetRecordDB().Exec(calcDataSQL) ///更新服务器最大在线人数和平均在线人数
	//	}
	loger().Info("DataCalcMgr::calcServerOnline end at %s insert rs count:%d",
		time.Now().Format("2006-01-02 15:04:05"), rowsAffected)
}

//update %s as us,(select count(distinct ra.teamid) remaincount,rc.sdkname from %s as ra,
//      (select distinct teamid,sdkname from record_create where type=2 and maketime between '%s 00:00:00' and '%s 23:59:59') as rc
//where ra.teamid=rc.teamid and ra.maketime between '%s 00:00:00' and '%s 23:59:59' group by rc.sdkname) as t
//set us.day%dnum=t.remaincount*10000/us.createnum where date='%s' and us.sdkname=t.sdkname;

//tableRecordCreate = "record_create" ///玩家创建表
//tableRecordAction = "record_action" ///玩家行为表
//tableRecordPay    = "record_pay"    ///玩家消费表
//tableRecordOnline = "record_online" ///玩家在线表
//留存
func (self *DataCalcMgr) calcRemainCountAndRate(pastTime string, pastTimeNew string, todayTime string, serverID int, pastDay int) { ///计算指定天留存数与留存率
	userstayTableName := "mmo2d_recordljzm.userstay" ///用户留存数据对应的表
	calcDataSQL := fmt.Sprintf(`
				update %s as us,(select c.sdkname,count(1) as num from 
					(select DISTINCT teamid,sdkname from %s where type=2 and maketime between '%s 00:00:00' and '%s 23:59:59') as c,
					(select DISTINCT teamid,sdkname from %s where type=1 and maketime between '%s 00:00:00' and '%s 23:59:59') as a
 					where a.teamid=c.teamid group by sdkname) as t
					set us.day%dnum=t.num where us.sid=%d and date='%s' and us.sdkname=t.sdkname;`,
		userstayTableName, tableRecordCreate, pastTimeNew, pastTimeNew, tableRecordAction, todayTime, todayTime, pastDay, serverID, pastTime)

	//	fmt.Println(calcDataSQL)
	GetServer().GetRecordDB().Exec(calcDataSQL)

}

//select sid,'%s',sdkname,sum(createnum) from %s where sid=%d and  date between '%s 00:00:00' and '%s 23:59:59' group by sdkname;`,

func (self *DataCalcMgr) calcRemainUserData(startTime string, endTime string, serverID int) { ///统计留存用户数据
	loger().Info("DataCalcMgr::calcRemainUserData begin at %s", time.Now().Format("2006-01-02 15:04:05"))
	todayTime := time.Now().Format("2006/01/02") ///取得今天日期串
	todayTimeNew := time.Now().Format("2006-01-02")
	loginDBName := GetServer().GetLoginDB().GetDBName()
	createTableName := "mmo2d_recordljzm.createdata" ///用户创建数据对应的表
	userstayTableName := "mmo2d_recordljzm.userstay" ///用户留存数据对应的表
	///生成新的记录存在今日总注册帐号数和总创建角色数
	calcDataSQL := fmt.Sprintf(`replace %s (sid,date,sdkname,createnum) 
		select %d,'%s',name,IFNULL(t.num,0) from %s.%s as s left join (select sdkname,sum(createnum) as num from %s where sid=%d and  date between '%s 00:00:00' and '%s 23:59:59' group by sdkname) as t on s.name=t.sdkname where s.enable=1;`,
		userstayTableName, serverID, todayTime, loginDBName, tableSDKList, createTableName, serverID, todayTime, todayTime)
	_, rowsAffected := GetServer().GetRecordDB().Exec(calcDataSQL)

	// 生成留存数据
	userStayList := IntList{1, 2, 3, 4, 5, 6, 7, 14, 30}
	for i := range userStayList {
		//self.calcRemainCountAndRate("2014/04/12", 2"2014/04/14", serverID)
		day := userStayList[i]
		queryTime := time.Now().AddDate(0, 0, day*-1)
		pastTime := queryTime.Format("2006/01/02")
		pastTimeNew := queryTime.Format("2006-01-02")
		self.calcRemainCountAndRate(pastTime, pastTimeNew, todayTimeNew, serverID, day)

	}

	loger().Info("DataCalcMgr::calcRemainUserData end at %s insert rs count:%d",
		time.Now().Format("2006-01-02 15:04:05"), rowsAffected)
}

func (self *DataCalcMgr) calcComplexData(startTime string, endTime string, serverID int) { ///统计注册登录数据
	loger().Info("DataCalcMgr::calcComplexData begin at %s", time.Now().Format("2006-01-02 15:04:05"))
	todayTime := time.Now().Format("2006/01/02")       ///取得今天日期串
	createTableName := "mmo2d_recordljzm.createdata"   ///用户创建数据对应的表
	onlineTableName := "mmo2d_recordljzm.onlinedata"   ///后台对应的表
	complexTableName := "mmo2d_recordljzm.complexdata" ///综合统计数据对应的表
	//	dynamicName := GetServer().GetDynamicDB().GetDBName()
	///生成新的记录存在今日总注册帐号数和总创建角色数
	calcDataSQL := fmt.Sprintf(`replace %s (sid,date,sdkname,registernum,createnum) 
		select sid,'%s',sdkname,sum(registernum)  as registernum,sum(createnum) as createnum from %s 
		where sid=%d and date between '%s 00:00:00' and '%s 23:59:59' group by sdkname;`,
		complexTableName, todayTime, createTableName, serverID, todayTime, todayTime)
	_, rowsAffected := GetServer().GetRecordDB().Exec(calcDataSQL)
	///更新记录中日登录总数,日最大在线人数,日平均在线人数
	calcDataSQL = fmt.Sprintf(`update %s as c,(select sid,'%s' as date,sum(loginnum) as loginnum,
	max(curonline) as maxonline,avg(curonline) as aveonline,sdkname  from %s where sid=%d and 
	date between '%s 00:00:00' and '%s 23:59:59' group by sdkname) as t 
	set c.loginnum=t.loginnum,c.maxonline=t.maxonline,c.aveonline=t.aveonline where c.sid=t.sid and c.date=t.date and c.sdkname=t.sdkname;`,
		complexTableName, todayTime, onlineTableName, serverID, todayTime, todayTime)
	_, rowsAffected = GetServer().GetRecordDB().Exec(calcDataSQL)
	///更新记录中的日消费钻石总费type=2为消费钻石
	calcDataSQL = fmt.Sprintf(`update %s as c,
	(select %d as sid,'%s' as date,sum(param) as consumption,sdkname from %s 
	where type=2 and maketime between  '%s 00:00:00' and '%s 23:59:59' group by sdkname) as t
	set c.consumption=t.consumption where c.sid=t.sid and c.date=t.date and c.sdkname=t.sdkname;`,
		complexTableName, serverID, todayTime, tableRecordPay, todayTime, todayTime)
	_, rowsAffected = GetServer().GetRecordDB().Exec(calcDataSQL)

	///更新记录中的日剩余钻石总和
	calcDataSQL = fmt.Sprintf(`update %s as c,
	(select %d as sid, '%s' as date, sum(money) as overyuanbao,sdkname from  %s where 
	maketime between '%s 00:00:00' and '%s 23:59:00' group by sdkname order by maketime desc) as t
	set c.overyuanbao=t.overyuanbao where c.sid=t.sid and c.date=t.date and c.sdkname=t.sdkname;`,
		complexTableName, serverID, todayTime, tableRecordPay, todayTime, todayTime)
	_, rowsAffected = GetServer().GetRecordDB().Exec(calcDataSQL)

	///更新记录中的活跃用户日剩余钻石总和
	//calcDataSQL = fmt.Sprintf(`update %s as c,
	//(select %d as sid, '%s 00:00:00' as date, sum(ticket) as overyuanbao72 from (select id,name,ticket from %s.%s as b,
	//      (select distinct teamid  from %s where type = 1 and maketime between '%s 00:00:00' and '%s 23:59:59') as a
	//      where b.id=a.teamid) as havemoney) as t
	//set c.overyuanbao72=t.overyuanbao72 where c.sid=t.sid and c.date=t.date;`,
	//	complexTableName, serverID, todayTime, dynamicName, tableTeam, tableRecordAction, todayTime, todayTime)
	//_, rowsAffected = GetServer().GetRecordDB().Exec(calcDataSQL)

	loger().Info("DataCalcMgr::calcComplexData end at %s insert rs count:%d",
		time.Now().Format("2006-01-02 15:04:05"), rowsAffected)
}

func (self *DataCalcMgr) calcRegLogin(startTime string, endTime string, serverID int) { ///统计注册登录数据
	loger().Info("DataCalcMgr::calcRegLogin begin at %s", time.Now().Format("2006-01-02 15:04:05"))
	targetTableName := "mmo2d_recordljzm.createdata" ///后台对应的表
	///type: 1.注册 2.创建
	calcDataSQL := fmt.Sprintf(`insert %s select %d as sid,'%s' as date,t.sdkname,
		sum(case type when '1' then count else 0 end) as registernum,
		sum(case type when '2' then count else 0 end) as createnum from 
		(select type,sdkname,count(id) as count from %s where maketime 
		between '%s' and '%s' group by sdkname,type) as t group by sdkname;`,
		targetTableName, serverID, endTime, tableRecordCreate, startTime, endTime)
	_, rowsAffected := GetServer().GetRecordDB().Exec(calcDataSQL) ///每隔一段时间写入注册人数与创建人数
	loger().Info("DataCalcMgr::calcRegLogin end at %s insert rs count:%d",
		time.Now().Format("2006-01-02 15:04:05"), rowsAffected)
}

func (self *DataCalcMgr) calcLevelDistributed(startTime string, endTime string, serverID int) { ///统计等级分布数据
	loger().Info("DataCalcMgr::calcLevelDistributed begin at %s", time.Now().Format("2006-01-02 15:04:05"))
	targetTableName := "mmo2d_recordljzm.leveldata" ///后台对应的表
	//	recordTableName := "record_create"
	todayTime := time.Now().Format("2006/01/02") ///取得今天日期串
	//	recordDBName := GetServer().GetRecordDB().GetDBName()
	dynamicDBName := GetServer().GetDynamicDB().GetDBName()
	loginDBName := GetServer().GetLoginDB().GetDBName()
	todayTimeNew := time.Now().Format("2006-01-02")

	removeOldDataSQL := fmt.Sprintf("delete from %s where calcdate='%s'", targetTableName, todayTimeNew)
	_, rowsAffected := GetServer().GetRecordDB().Exec(removeOldDataSQL) ///清除上次残留数据

	calcDataSQL := fmt.Sprintf(`replace %s(sid,sdkname,calcdate,levels,levelsnum) 
		select %d as sid,sdkname,'%s', level as levels, count(1) as levelsnum from (select sdkname,level from 
		%s.%s as t,%s.%s as a where t.accountid=a.id) as tt group by sdkname,level;`,
		targetTableName, serverID, todayTime, dynamicDBName, tableTeam, loginDBName, tableAccount)

	_, rowsAffected = GetServer().GetRecordDB().Exec(calcDataSQL)
	loger().Info("DataCalcMgr::calcLevelDistributed end at %s insert rs count:%d",
		time.Now().Format("2006-01-02 15:04:05"), rowsAffected)
}

func (self *DataCalcMgr) Run() {
	self.calcTimeDelaySecs = 600 ///每隔指定秒数运行统计逻辑
	//self.calcTimeDelaySecs = 30                          ///每隔指定秒数运行统计逻辑
	serverUpdateTimer := time.NewTicker(time.Second * 1) ///数据统计系统逻辑每1秒执行一次
	for now := range serverUpdateTimer.C {
		self.onTimer(now)
	}
}

func (self *DataCalcMgr) onTimer(now time.Time) {
	defer func() {
		x := recover()
		if x != nil {
			GetServer().GetLoger().Error("%v\r\n%s", x, debug.Stack())
		}
	}()
	nowUTC := now.Unix()
	if nowUTC%self.calcTimeDelaySecs > 0 {
		return ///每过十分钟统计一次
	}
	detalSecs := time.Second * time.Duration(self.calcTimeDelaySecs*-1)
	begin := now.Add(detalSecs)                      ///得到前一个间隔的时间对象
	startTime := begin.Format("2006/01/02 15:04:05") ///得到开始时间串
	endTime := now.Format("2006/01/02 15:04:05")     ///得到结束时间串
	serverID := GetServer().GetConfig().ServerID     ///得到服务器编号
	fmt.Println("DataCalcMgr::onTimer begin at", now.Format("2006-01-02 15:04:05"))
	self.calcRegLogin(startTime, endTime, serverID)         ///统计注册登录数据 ok
	self.calcServerOnline(startTime, endTime, serverID)     ///统计服务器在线数据 ok
	self.calcComplexData(startTime, endTime, serverID)      ///统计服务器综合数据 ok
	self.calcRemainUserData(startTime, endTime, serverID)   ///统计服务器用户留存数据
	self.calcLevelDistributed(startTime, endTime, serverID) ///统计服务器用户等级分布数据 ok
	fmt.Println("DataCalcMgr::onTimer end at", time.Now().Format("2006-01-02 15:04:05"))
}
