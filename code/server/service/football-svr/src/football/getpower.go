package football

import (
	"fmt"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////////////
//! 数据库结构
//! 领取体力信息
type GetPowerInfo struct {
	ID     int    //! 记录id
	Teamid int    //! 球队信息
	Time1  string //! 第一次领取时间
	Time2  string //! 第二次领取时间
}

type GetPower struct {
	GetPowerInfo //! 结构
	DataUpdater  //! 数据更新接口
}

func NewGetPower(teamid int) (*GetPower, bool) {
	p := new(GetPower)

	str := fmt.Sprintf("select * from `dy_getpower` where teamid = %d", teamid)
	GetServer().GetDynamicDB().fetchOneRow(str, &p.GetPowerInfo)

	//! 这里为了兼容老号
	if p.ID <= 0 {
		strSql := fmt.Sprintf("insert into `dy_getpower`(`teamid`, `time1`, `time2`) values (%d, '1990-01-01 00:00:00', '1990-01-01 00:00:00')", teamid)
		lastInsertID, rows := GetServer().GetDynamicDB().Exec(strSql)
		if rows <= 0 {
			GetServer().GetLoger().Warn("TeamCreateMsg processAction strSql fail! msg:%s", strSql)
			return p, false
		} else {
			p.GetPowerInfo.ID = lastInsertID
			p.GetPowerInfo.Teamid = teamid
			p.GetPowerInfo.Time1 = "1990-01-01 00:00:00"
			p.GetPowerInfo.Time2 = "1990-01-01 00:00:00"
		}
	}

	p.InitDataUpdater("dy_getpower", &p.GetPowerInfo)

	return p, p.ID > 0
}

func (self *GetPower) GetState() (int, int) {
	state1, state2 := 0, 0

	now := time.Now()

	GetServer().GetLoger().CYDebug("hours = %d", now.Hour())

	if now.Hour() >= 12 && now.Hour() < 14 {
		gettime, _ := time.ParseInLocation(TimeFormat, self.Time1, time.Local)
		GetServer().GetLoger().CYDebug("gettime = %v", gettime)
		if gettime.Day() != now.Day() || gettime.Month() != now.Month() || gettime.Year() != now.Year() {
			state1 = 1
		}
	}

	if now.Hour() >= 20 && now.Hour() < 22 {
		gettime, _ := time.ParseInLocation(TimeFormat, self.Time2, time.Local)
		if gettime.Day() != now.Day() || gettime.Month() != now.Month() || gettime.Year() != now.Year() {
			state2 = 1
		}
	}

	return state1, state2
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//! 消息
//! client2server
type GetPowerMsg struct {
	MsgHead `json:"head"` //! "getpower", "get"
	Index   int           `json:"index"` //! 0,请求;1,领取1;2,领取2
}

func (self *GetPowerMsg) GetTypeAndAction() (string, string) {
	return "getpower", "get"
}

func (self *GetPowerMsg) processAction(client IClient) bool {
	if self.doAction(client) == false {
		return false
	}
	return true
}

func (self *GetPowerMsg) doAction(client IClient) bool {
	team := client.GetTeam()

	getPowerResultMsg := new(GetPowerResultMsg)

	//! 得到当前时间
	now := time.Now()
	switch self.Index {
	case 0:
		getPowerResultMsg.Type = 0
	case 1:
		getPowerResultMsg.State1, getPowerResultMsg.State2 = team.PowerValue.GetState()
		if getPowerResultMsg.State1 == 1 {
			getPowerResultMsg.Type = 1
			team.AwardObject(awardTypeActionPoint, 60, 0, 0)
			team.PowerValue.Time1 = now.Format(TimeFormat)
		} else {
			getPowerResultMsg.Type = 2
		}
	case 2:
		getPowerResultMsg.State1, getPowerResultMsg.State2 = team.PowerValue.GetState()
		if getPowerResultMsg.State2 == 1 {
			getPowerResultMsg.Type = 1
			team.AwardObject(awardTypeActionPoint, 60, 0, 0)
			team.PowerValue.Time2 = now.Format(TimeFormat)
		} else {
			getPowerResultMsg.Type = 2
		}
	}

	getPowerResultMsg.State1, getPowerResultMsg.State2 = team.PowerValue.GetState()
	client.SendMsg(getPowerResultMsg)

	return true
}

//! server2client
type GetPowerResultMsg struct {
	MsgHead `json:"head"` //! "getpower", "getresult"
	State1  int           `json:"state1"` //! 1可领取 0不可领取
	State2  int           `json:"state2"` //! 1可领取 0不可领取
	Type    int           `json:"type"`   //! 0更新 1领取成功 2领取失败
}

func (self *GetPowerResultMsg) GetTypeAndAction() (string, string) {
	return "getpower", "getresult"
}
