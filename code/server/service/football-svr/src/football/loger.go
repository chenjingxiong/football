package football

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const CLR_G = "\x1b[32m" //32: 绿
const CLR_B = "\x1b[34m" //34: 蓝
const CLR_Y = "\x1b[33m" //34: 黄
const CLR_R = "\x1b[31m" //31: 红
const CLR_P = "\x1b[35m" //35: 紫
const CLR_N = "\033[0m"  //重置

const (
	logDebug = 1 ///调试
	logInfo  = 2 ///信息
	logWarn  = 3 ///警告
	logError = 4 ///错误
	logFatal = 5 ///致命错误
)

type Loger struct {
	loger          *log.Logger ///日志组件
	std            *log.Logger ///终端日志组件
	logMinLevel    int         ///最小日志等级
	terminalOutput bool        ///是否同时终端输出
	currentDay     int         ///当前天
	logPath        string      ///当前日志文件路径
}

//type ILoger interface {
//	Debug(format string, v ...interface{})                                                  ///调试
//	Info(format string, v ...interface{})                                                   ///信息
//	Warn(format string, v ...interface{})                                                   ///警告
//	Error(format string, v ...interface{})                                                  ///错误
//	Fatal(format string, v ...interface{})                                                  ///致命错误
//	CheckFail(formula string, result bool, variantA interface{}, variantB interface{}) bool ///不相等检测
//	Print(format string, v ...interface{})                                                  ///Print
//}

func loger() *Loger {
	return GetServer().GetLoger()
}

func (self *Loger) Output(logType int, format string, v ...interface{}) { ///信息
	if logType < self.logMinLevel {
		return ///过滤此等级
	}
	colorDic := map[string]string{"[Debug]": CLR_G, "[Info]": CLR_B, "[Warn]": CLR_Y, "[Error]": CLR_R, "[Fatal]": CLR_P}
	prefixStr := "unknow"
	switch logType {
	case logDebug:
		prefixStr = "[Debug]"
	case logInfo:
		prefixStr = "[Info]"
	case logWarn:
		prefixStr = "[Warn]"
	case logError:
		prefixStr = "[Error]"
	case logFatal:
		prefixStr = "[Fatal]"
	}
	s := fmt.Sprintf(prefixStr+format, v...)
	self.ChangDay() ///跨天则生成新文件
	self.loger.Output(3, s)
	if true == self.terminalOutput {
		s = strings.Replace(s, prefixStr, colorDic[prefixStr]+prefixStr+CLR_N, 1)
		self.std.Output(3, s)
		//self.std.Printf(s)
	}
}

func (self *Loger) CheckFail(formula string, result bool, variantA interface{}, variantB interface{}) bool {
	if true == result {
		return false ///不相等则直接返回不处理
	}
	strA, strB := "", ""
	sA := reflect.ValueOf(variantA)
	sB := reflect.ValueOf(variantB)
	if sA.Kind() == reflect.Ptr || sA.Kind() == reflect.Slice {
		strA = fmt.Sprintf("Get:0x%x", sA.Pointer())
	} else {
		strA = fmt.Sprintf("Get:%v", variantA)
	}
	if sB.Kind() == reflect.Ptr || sB.Kind() == reflect.Slice {
		strB = fmt.Sprintf("Need:%x", sB.Pointer())
	} else {
		strB = fmt.Sprintf("Need:%v", variantB)
	}
	pc, _, _, _ := runtime.Caller(1)
	funcObject := runtime.FuncForPC(pc)
	self.Output(logWarn, "%s Check %s Fail! %s %s", funcObject.Name(), formula, strA, strB)
	return result == false
}

func (self *Loger) Print(format string, v ...interface{}) { ///Print
	fmt.Printf(format, v...)
}

func (self *Loger) Debug(format string, v ...interface{}) { ///调试
	self.Output(logDebug, format, v...)
}

func (self *Loger) Info(format string, v ...interface{}) { ///信息
	self.Output(logInfo, format, v...)
}

func (self *Loger) Warn(format string, v ...interface{}) { ///警告
	self.Output(logWarn, format, v...)
}

func (self *Loger) Error(format string, v ...interface{}) { ///错误
	self.Output(logError, format, v...)
}

func (self *Loger) Fatal(format string, v ...interface{}) { ///致命错误,使用会造成服务器退出,慎用!!
	self.Output(logFatal, format, v...)
	os.Exit(1)
}

//! cy调试
func (self *Loger) CYDebug(format string, v ...interface{}) { ///调试
	return
	self.Output(logDebug, format, v...)
}

func (self *Loger) ChangDay() { ///跨天则生成新文件
	now := time.Now()
	currentDay := now.Day()
	fileName := self.logPath + "/" + now.Format("server-20060102.log")
	_, err := os.Stat(fileName)
	if self.currentDay == currentDay && nil == err {
		return
	}
	self.currentDay = currentDay
	logfile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		log.Println("open log file fail!", err.Error())
	}
	self.loger = log.New(logfile, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)
}

func NewLoger(logPath string, logMinLevel int, terminalOutput bool) *Loger {
	os.Mkdir(logPath, 7777) ///创建log目录,如果存在则忽略  log.Ldate|
	logerServer := new(Loger)
	logerServer.logPath = logPath
	logerServer.logMinLevel = logMinLevel
	logerServer.terminalOutput = terminalOutput
	logerServer.std = log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)
	return logerServer
}
