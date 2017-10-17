package football

import (
	"fmt"
	"reflect"
)

const (
	activitCodeNotUse = 0
	activitCodeBeUsed = 1
)

type ActivitCodeAward struct { //st_activitcodeaward表结构
	ID          int    `json:"id"`
	Name        string `json:"name"`
	SDKName     string `json:"sdkname"`
	Awarditem1  int    `json:"item"`
	Awardgrade1 int    `json:"grade"`
	Awardcount1 int    `json:"count"`
	Awarditem2  int    `json:"item"`
	Awardgrade2 int    `json:"grade"`
	Awardcount2 int    `json:"count"`
	Awarditem3  int    `json:"item"`
	Awardgrade3 int    `json:"grade"`
	Awardcount3 int    `json:"count"`
	Awarditem4  int    `json:"item"`
	Awardgrade4 int    `json:"grade"`
	Awardcount4 int    `json:"count"`
	Awarditem5  int    `json:"item"`
	Awardgrade5 int    `json:"grade"`
	Awardcount5 int    `json:"count"`
	Awarditem6  int    `json:"item"`
	Awardgrade6 int    `json:"grade"`
	Awardcount6 int    `json:"count"`
}

type ActivitCode struct {
	ID            int    `json:"id"`            //记录ID
	Code          string `json:"code"`          //激活码
	Type          int    `json:"type"`          //状态
	UserAccountID int    `json:"useraccountid"` //使用者帐号ID
	TeamID        int    `json:"teamid"`        //使用者球队id
	State         int    `json:"state"`         //状态
}

type ActivitCodeInfo struct {
	ActivitCode
	DataUpdater
}

func (self *ActivitCodeInfo) GetReflectValue() reflect.Value {
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

type ActivitCodeMgr struct {
}

func (self *ActivitCodeMgr) GetActiveCode(code string) *ActivitCode { ///判断激活码是否能用
	//	codeState := true
	loginDB := GetServer().GetLoginDB()
	activitCode := new(ActivitCode)
	activitCodeQuery := fmt.Sprintf("select * from dy_activationcode where code = '%s' limit 1", code)
	loginDB.fetchOneRow(activitCodeQuery, activitCode)
	if activitCode.ID <= 0 {
		return nil
	}
	return activitCode
	//if isExist == false {
	//	codeState = false
	//}

	//if activitCode.State == activitCodeBeUsed {
	//	codeState = false
	//}

	//return codeState
}

func (self *ActivitCodeMgr) Use(code string, userCountID int, teamID int) { ///使用该激活码
	loginDB := GetServer().GetLoginDB()
	activitCodeUpdate := fmt.Sprintf("update %s set state = %d,useraccountid = %d,teamid=%d where code = '%s'",
		tableActivationCode, activitCodeBeUsed, userCountID, teamID, code)
	loginDB.Exec(activitCodeUpdate)
}
