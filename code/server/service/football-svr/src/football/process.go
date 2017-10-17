package football

import (
	"reflect"
)

const (
	ProcessTypeStarTrain = 1 ///球员训练处理
	ProcessTypeTeamCD    = 2 ///球队CD
)

const (
	getProcessListAll     = 0  ///得到所有
	trainAwardExpInterval = 60 ///训练每60秒给一次经验
)

type IProcess interface {
	//	IObject
	IDataUpdater
	GetInfo() *ProcessInfo
	Reset() ///重置处理进程状态
	CalcTrainAwardExp(teamLevel int, isTotal bool) int
	CanProcess() bool               ///判断是否可以处理
	GetID() int                     ///得到ID
	GetReflectValue() reflect.Value ///得到反射对象
}

type ProcessInfo struct { ///处理对象,用于记录持续处理事务
	ID              int `json:"id"`              ///处理id
	TeamID          int `json:"teamid"`          ///所属球队id
	Type            int `json:"type"`            ///类型
	Pos             int `json:"pos"`             ///格子号
	Level           int `json:"level"`           ///位置等级
	ObjID           int `json:"objid"`           ///处理对象id
	ExpireTime      int `json:"expiretime"`      ///失效时间
	NextProcessTime int `json:"nextprocesstime"` ///下次处理时间
	Param1          int `json:"param1"`          ///参数1
	Param2          int `json:"param2"`          ///参数2
}

type Process struct { ///被球队拥有处理管理器,用于时间推移给效果的系统
	//	Object
	ProcessInfo
	DataUpdater
}

func (self *Process) GetInfo() *ProcessInfo {
	return &self.ProcessInfo
}

func (self *Process) GetReflectValue() reflect.Value { ///得到反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

func (self *Process) CanProcess() bool {
	canProcess := Now() >= self.NextProcessTime
	return canProcess
}

func (self *Process) Reset() { ///重置处理进程状态
	self.ObjID = 0
	self.ExpireTime = 0
	self.NextProcessTime = 0
	self.Param1 = 0
	self.Param2 = 0
}

func (self *Process) GetID() int { ///得到ID
	return self.ID
}

///计算训练奖励经验,isTotal为true是计算总剩余经验
func (self *Process) CalcTrainAwardExp(teamLevel int, isTotal bool) int {
	awardExp := teamLevel * 10
	if true == isTotal { ///期望得到总剩余经验
		now := Now()
		awardExp = Max(0, awardExp*(self.ExpireTime-now)/60) ///每分钟给一次
	}
	return awardExp
}

func NewProcess(processInfo *ProcessInfo) IProcess {
	process := new(Process)
	process.ProcessInfo = *processInfo
	process.InitDataUpdater(tableProcessCenter, &process.ProcessInfo)
	return process
}
