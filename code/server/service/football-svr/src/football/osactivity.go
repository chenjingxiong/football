package football

import (
	"fmt"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////////////
//! 数据库结构
//! 开服活动信息
type OSActivityInfo struct {
	ID     int    //! 记录id
	Teamid int    //! 球队信息
	Time1  string //! 首次充值时间
	Time2  string //! 首次10连抽时间
	Mask   int    //! 充值奖励掩码
}

type OSActivity struct {
	OSActivityInfo //! 结构
	DataUpdater    //! 数据更新接口
}

func NewOSActivity(teamid int) (*OSActivity, bool) {
	p := new(OSActivity)

	str := fmt.Sprintf("select * from `dy_osactivity` where teamid = %d", teamid)
	GetServer().GetDynamicDB().fetchOneRow(str, &p.OSActivityInfo)

	//! 这里为了兼容老号
	if p.ID <= 0 {
		strSql := fmt.Sprintf("insert into `dy_osactivity`(`teamid`, `time1`, `time2`, `mask`) values (%d, '1990-01-01 00:00:00', '1990-01-01 00:00:00', 0)", teamid)
		lastInsertID, rows := GetServer().GetDynamicDB().Exec(strSql)
		if rows <= 0 {
			GetServer().GetLoger().Warn("TeamCreateMsg processAction strSql fail! msg:%s", strSql)
			return p, false
		} else {
			p.OSActivityInfo.ID = lastInsertID
			p.OSActivityInfo.Teamid = teamid
			p.OSActivityInfo.Time1 = "1990-01-01 00:00:00"
			p.OSActivityInfo.Time2 = "1990-01-01 00:00:00"
			p.OSActivityInfo.Mask = 0
		}
	}

	p.InitDataUpdater("dy_osactivity", &p.OSActivityInfo)

	return p, p.ID > 0
}

func (self *OSActivity) Refresh(team *Team, index int) {
	if index == 0 { //! 充值
		//! 判断首次充值
		now := time.Now()
		gettime, _ := time.ParseInLocation(TimeFormat, self.Time1, time.Local)
		if gettime.Day() != now.Day() || gettime.Month() != now.Month() || gettime.Year() != now.Year() {
			//! 发送奖励
			team.GetMailMgr().SendSysAwardMail(SystemMail, LuckAwardMail, IntList{505016, 300307, 300052}, IntList{0, 0, 0}, IntList{5, 1, 1}, "", "")
			self.Time1 = now.Format(TimeFormat)
		}

		//! 判断累计充值
		if team.VipExp >= 1000 && (self.Mask&1) == 0 {
			team.GetMailMgr().SendSysAwardMail(SystemMail, LuckAwardMail, IntList{510003, 300021}, IntList{0, 0}, IntList{30, 5}, "", "")
			self.Mask |= 1
			GetServer().GetLoger().CYDebug("mask:%d", self.Mask)
		}

		if team.VipExp >= 5000 && (self.Mask&2) == 0 {
			team.GetMailMgr().SendSysAwardMail(SystemMail, LuckAwardMail, IntList{300301, 300011}, IntList{0, 0}, IntList{1, 10}, "", "")
			self.Mask |= 2
		}

		if team.VipExp >= 10000 && (self.Mask&4) == 0 {
			team.GetMailMgr().SendSysAwardMail(SystemMail, LuckAwardMail, IntList{300302, 300022}, IntList{0, 0}, IntList{1, 5}, "", "")
			self.Mask |= 4
		}

		if team.VipExp >= 20000 && (self.Mask&8) == 0 {
			team.GetMailMgr().SendSysAwardMail(SystemMail, LuckAwardMail, IntList{300303, 300012}, IntList{0, 0}, IntList{1, 5}, "", "")
			self.Mask |= 8
		}

		if team.VipExp >= 50000 && (self.Mask&16) == 0 {
			team.GetMailMgr().SendSysAwardMail(SystemMail, LuckAwardMail, IntList{300304, 300022}, IntList{0, 0}, IntList{1, 10}, "", "")
			self.Mask |= 16
		}

		if team.VipExp >= 100000 && (self.Mask&32) == 0 {
			team.GetMailMgr().SendSysAwardMail(SystemMail, LuckAwardMail, IntList{300305, 300013}, IntList{0, 0}, IntList{1, 6}, "", "")
			self.Mask |= 32
		}

	} else if index == 1 { //! 10连抽
		now := time.Now()
		gettime, _ := time.ParseInLocation(TimeFormat, self.Time2, time.Local)
		if gettime.Day() != now.Day() || gettime.Month() != now.Month() || gettime.Year() != now.Year() {
			//! 发送奖励
			team.GetMailMgr().SendSysAwardMail(SystemMail, LuckAwardMail, IntList{awardTypeTicket}, IntList{0}, IntList{525}, "", "")
			self.Time2 = now.Format(TimeFormat)
		}
	}
}
