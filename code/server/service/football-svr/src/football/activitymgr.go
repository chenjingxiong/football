package football

import (
	"fmt"
	"reflect"
)

const (
	ActivityAwardMaxCount = 10 ///活动奖励道具最大种类数
)

type ActivityAwardTypeList []*ActivityAwardType
type ActivitCodeAwardList []*ActivitCodeAward

type ActivityType struct { ///活动类型配置静态表
	ID          int    ///活动类型
	Name        string ///活动名称
	NeedType    int    ///条件类型 1使用道具次数 2日充值球票 3团购道具数 4消费球票数 5刷新球员数 6排行榜 7训练次数 8训练赛积分 9增加员属性点 10消费培养点数 11vip人数
	NeedSort    int    ///条件小分类 1普通个人条件 2普通服务器条件 3可连续个人 4可连续服务器
	TimeType    int    ///活动时间类型 1日活动 2周活动 3月活动 4循环活动 5每日活动
	StarTime    int    ///开始时间 例如1403111200 去掉年前两位精确到分钟
	TimeParam   int    ///时间参数 排行榜时为领奖时间 循环活动时为出现索引1表示第一天
	EndTime     int    ///结束时间
	NeedParam01 int    ///条件参数1 可以为0 道具类型 球员星级 训练赛项目品质颜色 训练赛积分标准 needtype类型专用排行
	NeedParam02 int    ///条件参数2 道具数量 道具星数 道具类型
	NeedParam03 int    ///条件参数2 道具数量 道具星数 排行榜名次
	AwardType0  int    ///奖励类型 对应activitytype中的id
	NeedParam11 int    ///条件参数1
	NeedParam12 int    ///条件参数2
	NeedParam13 int    ///条件参数2 道具数量 道具星数 排行榜名次
	AwardType1  int    ///奖励类型
	NeedParam21 int    ///条件参数1 道具类型
	NeedParam22 int    ///条件参数2 道具数量 道具星数
	NeedParam23 int    ///条件参数2 道具数量 道具星数 排行榜名次
	AwardType2  int    ///奖励类型
	NeedParam31 int    ///条件参数1 道具类型
	NeedParam32 int    ///条件参数2 道具数量 道具星数
	NeedParam33 int    ///条件参数2 道具数量 道具星数 排行榜名次
	AwardType3  int    ///奖励类型
	Needparam41 int    ///条件参数1 道具类型
	Needparam42 int    ///条件参数2 道具数量 道具星数
	Needparam43 int    ///条件参数2 道具数量 道具星数 排行榜名次
	AwardType4  int    ///奖励类型
	NeedParam51 int    ///条件参数1 道具类型
	NeedParam52 int    ///条件参数2 道具数量 道具星数
	NeedParam53 int    ///条件参数2 道具数量 道具星数 排行榜名次
	AwardType5  int    ///奖励类型
	NeedParam61 int    ///条件参数1 道具类型
	NeedParam62 int    ///条件参数2 道具数量 道具星数
	NeedParam63 int    ///条件参数2 道具数量 道具星数 排行榜名次
	AwardType6  int    ///奖励类型
	NeedParam71 int    ///条件参数1 道具类型
	NeedParam72 int    ///条件参数2 道具数量 道具星数
	NeedParam73 int    ///条件参数2 道具数量 道具星数 排行榜名次
	AwardType7  int    ///奖励类型
	NeedParam81 int    ///条件参数1 道具类型
	NeedParam82 int    ///条件参数2 道具数量 道具星数
	NeedParam83 int    ///条件参数2 道具数量 道具星数 排行榜名次
	AwardType8  int    ///奖励类型
	NeedParam91 int    ///条件参数1 道具类型
	NeedParam92 int    ///条件参数2 道具数量 道具星数
	NeedParam93 int    ///条件参数2 道具数量 道具星数 排行榜名次
	AwardType9  int    ///奖励类型
	Desc        string ///描述
}

type ActivityAwardType struct {
	ID           int ///活动奖励类型
	AwardItem01  int ///道具类型0(vip等级)1(奖励的第一个道具)
	AwardCount01 int ///道具数量
	AwardItem02  int ///道具类型
	AwardCount02 int ///道具数量
	AwardItem03  int ///道具类型
	AwardCount03 int ///道具数量
	AwardItem04  int ///道具类型
	AwardCount04 int ///道具数量
	AwardItem05  int ///道具类型
	AwardCount05 int ///道具数量
	AwardItem11  int ///道具类型
	AwardCount11 int ///道具星数
	AwardItem12  int ///道具类型
	AwardCount12 int ///道具数量
	AwardItem13  int ///道具类型
	AwardCount13 int ///道具数量
	AwardItem14  int ///道具类型
	AwardCount14 int ///道具数量
	AwardItem15  int ///道具类型
	AwardCount15 int ///道具数量
	AwardItem21  int ///道具类型
	AwardCount21 int ///道具星数
	AwardItem22  int ///道具类型
	AwardCount22 int ///道具数量
	AwardItem23  int ///道具类型
	AwardCount23 int ///道具数量
	AwardItem24  int ///道具类型
	AwardCount24 int ///道具数量
	AwardItem25  int ///道具类型
	AwardCount25 int ///道具数量
	AwardItem31  int ///道具类型
	AwardCount31 int ///道具星数
	AwardItem32  int ///道具类型
	AwardCount32 int ///道具数量
	AwardItem33  int ///道具类型
	AwardCount33 int ///道具数量
	AwardItem34  int ///道具类型
	AwardCount34 int ///道具数量
	AwardItem35  int ///道具类型
	AwardCount35 int ///道具数量
	AwardItem41  int ///道具类型
	AwardCount41 int ///道具星数
	AwardItem42  int ///道具类型
	AwardCount42 int ///道具数量
	AwardItem43  int ///道具类型
	AwardCount43 int ///道具数量
	AwardItem44  int ///道具类型
	AwardCount44 int ///道具数量
	AwardItem45  int ///道具类型
	AwardCount45 int ///道具数量
	AwardItem51  int ///道具类型
	AwardCount51 int ///道具星数
	AwardItem52  int ///道具类型
	AwardCount52 int ///道具数量
	AwardItem53  int ///道具类型
	AwardCount53 int ///道具数量
	AwardItem54  int ///道具类型
	AwardCount54 int ///道具数量
	AwardItem55  int ///道具类型
	AwardCount55 int ///道具数量
}

func (self *ActivityType) GetAwardType(activityItemIndex int) int { ///通过Vip等级得到奖励道具类型列表与数量列表
	activityAwardList := IntList{self.AwardType0, self.AwardType1,
		self.AwardType2, self.AwardType3,
		self.AwardType4, self.AwardType5, self.AwardType6,
		self.AwardType7, self.AwardType8, self.AwardType9}
	activityItemIndexReal := activityItemIndex - 1 ///转换成从0开始的索引
	activityAwardType := activityAwardList[activityItemIndexReal]
	return activityAwardType
}

func (self *ActivityAwardType) GetAwardItemList(vipLevel int) (IntList, IntList) { ///通过Vip等级得到奖励道具类型列表与数量列表
	awardTypeList, awardCountList := IntList{}, IntList{}
	value := reflect.ValueOf(self).Elem()
	awardTypeFieldName, awardCountFieldName := "", ""
	for i := 1; i < ActivityAwardMaxCount; i++ {
		awardTypeFieldName = fmt.Sprintf("AwardItem%d%d", vipLevel, i)
		awardTypeField := value.FieldByName(awardTypeFieldName)
		if awardTypeField.IsValid() == false {
			break ///没有此字段直接跳出
		}
		awardCountFieldName = fmt.Sprintf("AwardCount%d%d", vipLevel, i)
		awardCountField := value.FieldByName(awardCountFieldName)
		if awardCountField.IsValid() == false {
			break ///没有此字段直接跳出
		}
		awardType := int(awardTypeField.Int())
		awardCount := int(awardCountField.Int())
		if awardType <= 0 || awardCount <= 0 {
			continue
		}
		awardTypeList = append(awardTypeList, awardType)
		awardCountList = append(awardCountList, awardCount)
	}
	return awardTypeList, awardCountList
}

type ActivityInfo struct { ///活动动态表,记录当前每个球队活动完成情况
	ID             int `json:"id"`             ///id编号
	Type           int `json:"type"`           ///活动类型
	TeamID         int `json:"teamid"`         ///球队id 如果id为-1表示此记录为服务器所有
	Progress       int `json:"progress"`       ///当前进度
	AwardCurCount0 int `json:"awardcurcount0"` ///已领奖次数
	AwardMaxCount0 int `json:"awardmaxcount0"` ///可领奖次数
	AwardCurCount1 int `json:"awardCurCount1"` ///已领奖次数
	AwardMaxCount1 int `json:"awardmaxcount1"` ///可领奖次数
	AwardCurCount2 int `json:"awardcurcount2"` ///已领奖次数
	AwardMaxCount2 int `json:"awardmaxcount2"` ///可领奖次数
	AwardCurCount3 int `json:"awardcurcount3"` ///已领奖次数
	AwardMaxCount3 int `json:"awardmaxcount3"` ///可领奖次数
	AwardCurCount4 int `json:"awardcurcount4"` ///已领奖次数
	AwardMaxCount4 int `json:"awardmaxcount4"` ///可领奖次数
	AwardCurCount5 int `json:"awardcurcount5"` ///已领奖次数
	AwardMaxCount5 int `json:"awardmaxcount5"` ///可领奖次数
	AwardCurCount6 int `json:"awardcurcount6"` ///已领奖次数
	AwardMaxCount6 int `json:"awardmaxcount6"` ///可领奖次数
	AwardCurCount7 int `json:"awardcurcount7"` ///已领奖次数
	AwardMaxCount7 int `json:"awardmaxcount7"` ///可领奖次数
	AwardCurCount8 int `json:"awardcurcount8"` ///已领奖次数
	AwardMaxCount8 int `json:"awardmaxcount8"` ///可领奖次数
	AwardCurCount9 int `json:"awardcurcount9"` ///已领奖次数
	AwardMaxCount9 int `json:"awardmaxcount9"` ///可领奖次数
}

type Activity struct { ///活动对象
	ActivityInfo
	DataUpdater ///信息更新组件
}

func (self *Activity) GetReflectValue() reflect.Value { ///得到球队反射对象
	reflectValue := reflect.ValueOf(self).Elem()
	return reflectValue
}

func (self *Activity) GetCurrentAwardInfo(itemIndex int) (int, int) { ///得到活动当前已领奖次数
	itemIndexReal := itemIndex - 1
	reflectValue := reflect.ValueOf(self).Elem()
	curCountFieldName := fmt.Sprintf("AwardCurCount%d", itemIndexReal) ///当前领奖次数
	maxCountFieldName := fmt.Sprintf("AwardMaxCount%d", itemIndexReal) ///最大领奖次数
	curCountField := reflectValue.FieldByName(curCountFieldName)
	maxCountField := reflectValue.FieldByName(maxCountFieldName)
	curCount := int(curCountField.Int())
	maxCount := int(maxCountField.Int())
	return curCount, maxCount
}

func (self *Activity) SetCurrentAwardInfo(itemIndex int, curAwardCount int, maxAwardCount int) { ///得到活动当前已领奖次数
	itemIndexReal := itemIndex - 1
	reflectValue := reflect.ValueOf(self).Elem()
	curCountFieldName := fmt.Sprintf("AwardCurCount%d", itemIndexReal) ///当前领奖次数
	maxCountFieldName := fmt.Sprintf("AwardMaxCount%d", itemIndexReal) ///最大领奖次数
	curCountField := reflectValue.FieldByName(curCountFieldName)
	maxCountField := reflectValue.FieldByName(maxCountFieldName)
	curCountField.SetInt(int64(curAwardCount))
	maxCountField.SetInt(int64(maxAwardCount))
}

type ActivityInfoList []ActivityInfo
type ActivityList map[int]*Activity

type ActivityMgr struct { ///活动管理器
	GameMgr                   ///逻辑系统管理器
	activityList ActivityList ///活动动态信息列表
}

func (self *ActivityMgr) GetType() int { ///得到管理器类型
	return mgrTypeActivityMgr ///任务管理器
}

///创建默认的活动对象并写到数据库中,并返回新建的对象id
func (self *ActivityMgr) createDefaultActivity(activityType int) *Activity {
	activity := new(Activity)
	createActivityQuery := fmt.Sprintf("insert %s (type,teamid) VALUES(%d,%d)",
		tableActivity, activityType, self.team.ID)
	lastInsertActivityID, _ := GetServer().GetDynamicDB().Exec(createActivityQuery)
	if lastInsertActivityID <= 0 {
		return nil ///插入失败返回nil
	}
	activityQuery := fmt.Sprintf("select * from %s where id=%d limit 1", tableActivity, lastInsertActivityID)
	GetServer().GetDynamicDB().fetchOneRow(activityQuery, &activity.ActivityInfo)
	activity.InitDataUpdater(tableActivity, &activity.ActivityInfo)
	self.activityList[activity.Type] = activity ///保存了到列表中
	return activity
}

func (self *ActivityMgr) updateActivity(activity *Activity) bool { ///更新活动对象数据,true表示有数据更新
	return true
}

func (self *ActivityMgr) GetActivity(activityType int) *Activity { ///得到指定类型的活动对象
	activity := self.activityList[activityType]
	if nil == activity {
		///还没有此活动的信息,需要创建一个
		activity = self.createDefaultActivity(activityType)
	}
	return activity
}

///查询指定类型的活动对象,同时对此活动进行更新操作
func (self *ActivityMgr) QueryActivity(activityType int) *Activity {
	activity := self.GetActivity(activityType)
	if activity != nil {
		self.updateActivity(activity)
	}
	return activity
}

func (self *ActivityMgr) SaveInfo() { ///保存数据
	for _, v := range self.activityList {
		v.Save()
	}
}

func NewActivityMgr(teamID int) IGameMgr {
	activityMgr := new(ActivityMgr)
	if activityMgr.Init(teamID) == false {
		return nil
	}
	return activityMgr
}

func NewActivity(activityInfo *ActivityInfo) *Activity {
	activity := new(Activity)
	activity.ActivityInfo = *activityInfo
	activity.InitDataUpdater(tableActivity, &activity.ActivityInfo)
	return activity
}

func (self *ActivityMgr) Init(teamID int) bool {
	self.activityList = make(ActivityList) ///创建活动列表
	taskListQuery := fmt.Sprintf("select * from %s where teamid=%d limit 500", tableActivity, teamID)
	activityInfo := new(ActivityInfo)
	activityInfoList := GetServer().GetDynamicDB().fetchAllRows(taskListQuery, activityInfo)
	for i := range activityInfoList {
		activityInfo = activityInfoList[i].(*ActivityInfo)
		self.activityList[activityInfo.Type] = NewActivity(activityInfo) ///注意此处是用活动类型做id索引
	}
	return true
}

func (self *ActivityMgr) OnUseItem(itemType int, count int) { ///捕捉道具使用情况并进行统计

}

func GetActivityType(activityType int) *ActivityType {
	staticDataMgr := GetServer().GetStaticDataMgr()
	element := staticDataMgr.GetStaticData(tableActivityType, activityType)
	if nil == element {
		return nil
	}
	return element.(*ActivityType)
}

func (self *ActivityType) GetActivityAwardType(activityAwardType int) *ActivityAwardType {
	staticDataMgr := GetServer().GetStaticDataMgr()
	element := staticDataMgr.GetStaticData(tableActivityAward, activityAwardType)
	if nil == element {
		return nil
	}
	return element.(*ActivityAwardType)
}

func (self *ActivitCodeMgr) HasPrefix(s, prefix string) bool { ///前缀判断函数
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

//func (self *ActivitCodeMgr) FindPrefix(activitcodeprefix string) int { ///查找前缀函数
//	staticDataMgr := GetServer().GetStaticDataMgr()
//	//println("传入的前缀：", activitcodeprefix)
//	for _, v := range staticDataMgr.staticDataList[tableAwardType] {
//		staticData := v.(*ActivitCodeAward)
//		if staticData.Prefix == activitcodeprefix {
//			return v.(*ActivitCodeAward).ID
//		}
//	}
//	return 0
//}

///通过激活码前缀查找对应的奖励信息
func (self *ActivitCodeMgr) GetActiveCodeAwardType(activitCodeType int) *ActivitCodeAward {
	//rows := GetServer().GetLoginDB().Query("select type from %s where code='%s' limit 1", tableActivationCode, activitCode)
	//if nil == rows {
	//	return nil
	//}
	//activitCodeType := 0
	//for rows.Next() {
	//	rows.Scan(&activitCodeType)
	//}
	//rows.Close()
	staticDataMgr := GetServer().GetStaticDataMgr()
	element := staticDataMgr.GetStaticData(tableAwardType, activitCodeType)
	if nil == element {
		return nil
	}
	return element.(*ActivitCodeAward)
}

//func (self *ActivitCodeMgr) SubString(str string, begin, length int) string {
//	rs := []rune(str) /// 将字符串的转换成[]rune
//	lth := len(rs)
//	if begin < 0 { /// 越界判断
//		begin = 0
//	}
//	if begin >= lth {
//		begin = lth
//	}
//	end := begin + length
//	if end > lth {
//		end = lth
//	}
//	//println("函数截取的前四个字符：", string(rs[begin:end]))
//	return string(rs[begin:end])

//}
