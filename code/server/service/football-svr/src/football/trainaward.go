package football

type TrainAwardStaticData struct {
	ID          int ///id编号
	Level       int ///条件等级段
	GreenAward  int ///绿品质训练奖励点数
	BlueAward   int ///蓝品质训练奖励点数
	PurpleAward int ///紫品质训练奖励点数
	OrangeAward int ///橙品质训练奖励点数
	Score1      int ///分数1
	AwardType1  int ///奖励类型1
	AwardObj1   int ///奖励实体1
	AwardNum1   int ///奖励数量1
	Score2      int ///分数2
	AwardType2  int ///奖励类型2
	AwardObj2   int ///奖励实体2
	AwardNum2   int ///奖励数量2
	Score3      int ///分数3
	AwardType3  int ///奖励类型3
	AwardObj3   int ///奖励实体3
	AwardNum3   int ///奖励数量3
	Score4      int ///分数4
	AwardType4  int ///奖励类型4
	AwardObj4   int ///奖励实体4
	AwardNum4   int ///奖励数量4
	Score5      int ///分数5
	AwardType5  int ///奖励类型5
	AwardObj5   int ///奖励实体5
	AwardNum5   int ///奖励数量5
	Score6      int ///分数6
	AwardType6  int ///奖励类型6
	AwardObj6   int ///奖励实体6
	AwardNum6   int ///奖励数量6
	Score7      int ///分数7
	AwardType7  int ///奖励类型7
	AwardObj7   int ///奖励实体7
	AwardNum7   int ///奖励数量7
	Score8      int ///分数8
	AwardType8  int ///奖励类型8
	AwardObj8   int ///奖励实体8
	AwardNum8   int ///奖励数量8
	Score9      int ///分数9
	AwardType9  int ///奖励类型9
	AwardObj9   int ///奖励实体9
	AwardNum9   int ///奖励数量9
}

func GetTrainAwardType(trainAwardType int) *TrainAwardStaticData {
	staticDataMgr := GetServer().GetStaticDataMgr()
	element := staticDataMgr.GetStaticData(tableTrainAward, trainAwardType)
	if nil == element {
		return nil
	}
	return element.(*TrainAwardStaticData)
}

//func (self *StaticDataMgr) GetTaskTypeList() TaskTypePtrList { ///得到所有任务类型对象指针列表
//	if nil == self.staticDataList[tableTaskType] {
//		return TaskTypePtrList{}
//	}
//	taskTypePtrList := TaskTypePtrList{}
//	for _, v := range self.staticDataList[tableTaskType] {
//		taskTypePtrList = append(taskTypePtrList, v.(*TaskTypeStaticData))
//	}
//	return taskTypePtrList
//}

func FindTrainAwardByLevel(teamLevel int) *TrainAwardStaticData {
	staticDataMgr := GetServer().GetStaticDataMgr()
	staticDataList := staticDataMgr.GetStaticDataList(tableTrainAward)
	resultID := 0
	for _, v := range staticDataList {
		trainAwardStaticData := v.(*TrainAwardStaticData)
		if trainAwardStaticData.Level > teamLevel {
			break
		}
		resultID = trainAwardStaticData.ID
	}
	if resultID <= 0 {
		return nil
	}
	trainAwardStaticDataResult := staticDataList[resultID].(*TrainAwardStaticData)
	return trainAwardStaticDataResult
}
