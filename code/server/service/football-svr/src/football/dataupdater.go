package football

import (
	"fmt"
	"reflect"
	"strings"
)

type DataValueChangItem struct {
	FieldName string ///字段名
	///以下数据二选一
	//Value  int    ///值数据
	//String string ///串数据
}

type DataValueChangList []string ///数值变更列表 key是字段名 value是变更项

type IDataUpdater interface { ///数据更新器接口
	InitDataUpdater(tableName string, dataUser interface{}) ///初始化数据更新组件
	Save()                                                  ///马上保存数据
	//	SavePlus()                                               ///加强版存储
	Sync() DataValueChangList ///得到属性变更列表
	InsertSql() string        ///得到插入SQL语句

}

type DataUpdater struct { ///数据更新器
	baseData  reflect.Value ///基础数据
	syncData  reflect.Value ///同步数据,需要实时刷新
	dirtyData reflect.Value ///保存用户数据指针
	tableName string        ///数据表名
	dataID    int           ///记录主键
}

func NewDataUpdater(tableName string, dataUser interface{}) *DataUpdater {
	dataUpdater := new(DataUpdater)
	dataUpdater.InitDataUpdater(tableName, dataUser)
	return dataUpdater
}

//func (self *DataUpdater) buildSQLString()string { ///根据对象字段名与表名生成创建记录的SQL语句

//}

func (self *DataUpdater) InitDataUpdater(tableName string, dataUser interface{}) { ///初始化数据更新组件
	self.dirtyData = reflect.ValueOf(dataUser).Elem() ///得到用户数据反射数据
	self.baseData = reflect.New(self.dirtyData.Type()).Elem()
	self.baseData.Set(self.dirtyData)
	self.syncData = reflect.New(self.dirtyData.Type()).Elem()
	self.syncData.Set(self.dirtyData)
	self.dataID = int(self.dirtyData.Field(0).Int()) ///保存了数据id
	self.tableName = tableName                       ///保存表名
}

func (self *DataUpdater) InsertSql() string {
	fieldLen := self.dirtyData.NumField() ///得到结构体的字段数
	fieldList := ""
	valueList := ""
	for i := 0; i < fieldLen; i++ {
		dirtyInt := int64(0)
		dirtyStr := ""
		if 0 == i {
			if self.dirtyData.Field(i).Int() != 0 {
				return ""
			}

			continue
		}

		dirtyFieldName := self.baseData.Type().Field(i).Name
		fieldList += fmt.Sprintf("%s,", dirtyFieldName)
		ok := self.dirtyData.Field(i).Type() == self.baseData.Field(i).Type()
		if ok != true {
			continue ///类型不相同则跳过
		}
		switch self.dirtyData.Field(i).Kind() {
		case reflect.Int64:
			dirtyInt = self.dirtyData.Field(i).Int()
			valueList += fmt.Sprintf("%d,", int(dirtyInt))
		case reflect.Int:
			dirtyInt = self.dirtyData.Field(i).Int()
			valueList += fmt.Sprintf("%d,", int(dirtyInt))
		case reflect.String:
			dirtyStr = self.dirtyData.Field(i).String()
			valueList += fmt.Sprintf("'%s',", dirtyStr)
		default:
			continue ///无匹配类型跳过这个字段
		}
	}

	insertQuery := fmt.Sprintf("insert into %s (%s) values(%s)", self.tableName, fieldList, valueList)
	insertQuery = strings.Replace(insertQuery, ",)", ")", 2) ///去掉多余的逗号
	insertQuery = strings.ToLower(insertQuery)
	return insertQuery
}

func (self *DataUpdater) Sync() DataValueChangList { ///同步属性并得到属性变更列表
	dataValueChangList := DataValueChangList{}
	fieldLen := self.dirtyData.NumField() ///得到结构体的字段数
	for i := 1; i < fieldLen; i++ {       ///从1开始,跳过id字段
		syncInt, dirtyInt := int64(0), int64(0)
		syncStr, dirtyStr := "", ""
		dirtyFieldName := self.syncData.Type().Field(i).Name
		ok := self.dirtyData.Field(i).Type() == self.syncData.Field(i).Type()
		if ok != true {
			continue ///类型不相同则跳过
		}
		switch self.dirtyData.Field(i).Kind() {
		case reflect.Int:
			syncInt = self.syncData.Field(i).Int()
			dirtyInt = self.dirtyData.Field(i).Int()
		case reflect.String:
			syncStr = self.syncData.Field(i).String()
			dirtyStr = self.dirtyData.Field(i).String()
		default:
			continue ///无匹配类型跳过这个字段
		}
		if syncInt != dirtyInt {
			dataValueChangList = append(dataValueChangList, dirtyFieldName)
			self.syncData.Field(i).SetInt(dirtyInt) ///更新sync数据
		} else if syncStr != dirtyStr {
			dataValueChangList = append(dataValueChangList, dirtyFieldName)
			self.syncData.Field(i).SetString(dirtyStr) ///更新sync数据
		}
	}
	return dataValueChangList
}

func (self *DataUpdater) Save() { ///马上保存数据
	valueList := ""
	fieldLen := self.dirtyData.NumField() ///得到结构体的字段数
	if self.dataID <= 0 {
		self.dataID = int(self.dirtyData.Field(0).Int())
	}
	for i := 1; i < fieldLen; i++ { ///从1开始,跳过id字段
		baseInt, dirtyInt := int64(0), int64(0)
		baseStr, dirtyStr := "", ""
		//baseFieldName := self.dirtyData.Type().Field(i).Name
		dirtyFieldName := self.baseData.Type().Field(i).Name
		ok := self.dirtyData.Field(i).Type() == self.baseData.Field(i).Type()
		if ok != true {
			continue ///类型不相同则跳过
		}

		switch self.dirtyData.Field(i).Kind() {
		case reflect.Int64:
			baseInt = self.baseData.Field(i).Int()
			dirtyInt = self.dirtyData.Field(i).Int()
		case reflect.Int:
			baseInt = self.baseData.Field(i).Int()
			dirtyInt = self.dirtyData.Field(i).Int()
		case reflect.String:
			baseStr = self.baseData.Field(i).String()
			dirtyStr = self.dirtyData.Field(i).String()
		default:
			continue ///无匹配类型跳过这个字段
		}
		if dirtyFieldName == "Desc" || dirtyFieldName == "desc" {
			continue //! 跳过描述字段
		}

		if baseInt != dirtyInt {
			valueList += fmt.Sprintf("%s=%d,", dirtyFieldName, int(dirtyInt))
			self.baseData.Field(i).SetInt(dirtyInt) ///更新base数据
		} else if baseStr != dirtyStr {
			valueList += fmt.Sprintf("%s='%s',", dirtyFieldName, dirtyStr)
			self.baseData.Field(i).SetString(dirtyStr) ///更新base数据
		}
	}
	if valueList != "" {
		updateQuery := fmt.Sprintf("update %s set %s where id=%d limit 1", self.tableName, valueList, self.dataID)
		updateQuery = strings.Replace(updateQuery, ", where", " where", 1) ///去掉多余的逗号
		updateQuery = strings.ToLower(updateQuery)                         ///转成小写,避免数据库大小写敏感
		GetServer().GetDynamicDB().Exec(updateQuery)                       ///更新到数据库
	}
}

func (self *DataUpdater) SavePlus() {
	valueList := ""
	fieldList := ""
	fieldLen := self.dirtyData.NumField() ///得到结构体的字段数

	///先查询数据是否存在
	selectQuery := fmt.Sprintf("select * from %s where id=%d limit 1", self.tableName, self.dataID)
	selectQuery = strings.Replace(selectQuery, ", where", " where", 1)
	selectQuery = strings.ToLower(selectQuery)
	row := GetServer().GetDynamicDB().Query(selectQuery)
	isHasData := row.Next()

	for i := 1; i < fieldLen; i++ { ///从1开始,跳过id字段
		baseInt, dirtyInt := int64(0), int64(0)
		baseStr, dirtyStr := "", ""
		//baseFieldName := self.dirtyData.Type().Field(i).Name
		dirtyFieldName := self.baseData.Type().Field(i).Name
		ok := self.dirtyData.Field(i).Type() == self.baseData.Field(i).Type()
		if ok != true {
			continue ///类型不相同则跳过
		}
		switch self.dirtyData.Field(i).Kind() {
		case reflect.Int64:
			baseInt = self.baseData.Field(i).Int()
			dirtyInt = self.dirtyData.Field(i).Int()
		case reflect.Int:
			baseInt = self.baseData.Field(i).Int()
			dirtyInt = self.dirtyData.Field(i).Int()
		case reflect.String:
			baseStr = self.baseData.Field(i).String()
			dirtyStr = self.dirtyData.Field(i).String()
		default:
			continue ///无匹配类型跳过这个字段
		}

		if false == isHasData {
			if baseInt == dirtyInt {
				fieldList += fmt.Sprintf("%s,", dirtyFieldName)
				valueList += fmt.Sprintf("%d,", int(dirtyInt))
				self.baseData.Field(i).SetInt(dirtyInt) ///更新base数据
			} else if baseStr == dirtyStr {
				fieldList += fmt.Sprintf("%s,", dirtyFieldName)
				valueList += fmt.Sprintf("'%s',", dirtyStr)
				self.baseData.Field(i).SetString(dirtyStr) ///更新base数据
			}
			continue
		}

		if baseInt != dirtyInt {
			fieldList += fmt.Sprintf("%s=%d,", dirtyFieldName, int(dirtyInt))
			self.baseData.Field(i).SetInt(dirtyInt) ///更新base数据
		} else if baseStr != dirtyStr {
			fieldList += fmt.Sprintf("%s='%s',", dirtyFieldName, dirtyStr)
			self.baseData.Field(i).SetString(dirtyStr) ///更新base数据
		}
	}
	if fieldList != "" {
		if true == isHasData {
			///若存在数据,则直接使用Update
			updateQuery := fmt.Sprintf("update %s set %s where id=%d limit 1", self.tableName, fieldList, self.dataID)
			updateQuery = strings.Replace(updateQuery, ", where", " where", 1) ///去掉多余的逗号
			updateQuery = strings.ToLower(updateQuery)                         ///转成小写,避免数据库大小写敏感
			GetServer().GetDynamicDB().Exec(updateQuery)                       ///更新到数据库
		} else {
			if valueList != "" {
				///若不存在数据,则使用Insert
				insertQuery := fmt.Sprintf("insert into %s (%s) values(%s)", self.tableName, fieldList, valueList)
				insertQuery = strings.Replace(insertQuery, ",)", ")", 2) ///去掉多余的逗号
				insertQuery = strings.ToLower(insertQuery)
				GetServer().GetDynamicDB().Exec(insertQuery)
			}
		}

	}
}
