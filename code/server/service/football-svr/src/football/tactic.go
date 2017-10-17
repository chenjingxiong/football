package football

import (
//	"encoding/json"
//"log"
)

type TacticTypePtr *TacticTypeStaticData

type TacticTypeStaticData struct { ///阵型战术类型表
	ID            int    ///类型编号
	Name          string ///名字
	OpenFormType  int    ///开启阵型类型条件
	OpenFormLevel int    ///开启阵型等级条件
	OpenTicketPay int    ///开启所需球票价格
	Attrib1       int    ///加成属性类型1
	Attrib2       int    ///加成属性类型2
	Pos1          int    ///加成位置 1加成 0不加成
	Pos2          int    ///加成位置 1加成 0不加成
	Pos3          int    ///加成位置 1加成 0不加成
	Pos4          int    ///加成位置 1加成 0不加成
	Pos5          int    ///加成位置 1加成 0不加成
	Pos6          int    ///加成位置 1加成 0不加成
	Pos7          int    ///加成位置 1加成 0不加成
	Pos8          int    ///加成位置 1加成 0不加成
	Pos9          int    ///加成位置 1加成 0不加成
	Pos10         int    ///加成位置 1加成 0不加成
	Pos11         int    ///加成位置 1加成 0不加成
	Over          int    ///被克制战术类型
	Desc          string ///描述
}

//type TacticBaseStaticData struct { ///战术基础表
//	ID   int    ///类型编号
//	Name string ///名字
//	Over int    ///受制战术类型
//}

//type Tactic struct {
//	TacticInfo
//}

//func (self *Tactic) Create(tacticInfo *TacticInfo) bool { ///加载阵形战术信息
//	if self.ID != 0 {
//		return false
//	}
//	self.TacticInfo = *tacticInfo
//	return true
//}
