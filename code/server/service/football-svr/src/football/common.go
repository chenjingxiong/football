package football

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	//	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
)

const skipStackNumDefault = 2 ///默认跳过堆栈次数
const TimeFormat = "2006-01-02 15:04:05"

//const TimeFormatZeroTail = "2006-01-02" ///扔弃小时分钟秒
///通用函数与通用数据结构算法等

type IntList []int               ///整型数值列表
type IntPtrList []*int           ///整型指针列表
func (self *IntList) Len() int { ///取得列表长度
	return len(*self)
}
func (self *IntPtrList) Len() int { ///取得列表长度
	return len(*self)
}

func (self *IntList) Search(value int) int { ///搜索
	for i := range *self {
		v := (*self)[i]
		if v == value {
			return i
		}
	}
	return -1
}

func (self *IntList) AppendIfMissing(slice []int, i int) []int {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func (self *IntList) IsUnique() bool { ///判断数组是否唯一
	uniqueList := self.Unique()
	uniqueListLen := uniqueList.Len()
	selfLen := self.Len()
	return uniqueListLen == selfLen
}

func (self *IntList) Copy() IntList {
	newList := IntList{}
	newList = append(newList, *self...)
	return newList
}

func (self *IntList) Unique() IntList { ///返回唯一列表
	intList := self.Copy()
	sort.Ints(intList)
	newIntList := IntList{}
	oldValue := 0
	for i := range intList {
		if i <= 0 {
			i = oldValue
			newIntList = append(newIntList, intList[i])
			oldValue = intList[i]
			continue
		}

		if intList[i] == oldValue {
			continue
		}

		newIntList = append(newIntList, intList[i])
		oldValue = intList[i]
	}

	return newIntList
}

//hmac ,use sha1
func CalcHmac(data string, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

func CalcHmacRaw(data string, key string) []byte {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return mac.Sum(nil)
}

func CalcMD5(data string) string {
	t := md5.New()
	t.Write([]byte(data))
	return hex.EncodeToString(t.Sum(nil))
}

func GetMacAddress() string {
	//inter, _ := net.InterfaceByName("eth0")
	//macAddress := inter.HardwareAddr.String()
	//return macAddress
	macAddress := ""
	netInterfaces, _ := net.Interfaces()
	for i := range netInterfaces {
		netInterface := netInterfaces[i]
		hardwareAddr := netInterface.HardwareAddr.String()
		if hardwareAddr != "" {
			macAddress = hardwareAddr
			break
		}
	}
	return macAddress
}

func GetExpireTime(seconds int) time.Time { ///生成指定秒之后的到期时间
	now := time.Now()
	expireTime := now.Add(time.Duration(seconds) * time.Second)
	return expireTime
}

func IsExpireTime(expireTime int) bool { ///判断当前时间是否已经过期
	now := int(time.Now().Unix()) ///得到当前时间
	isExpireTime := now >= expireTime
	return isExpireTime
}

func GetHourＭinTime(hour int, minute int) time.Time { ///得到指定小时和分钟的utc时间秒
	now := time.Now()
	nowHour := now.Hour()
	nowMinute := now.Minute()
	if nowHour > hour {
		now = now.AddDate(0, 0, 1) ///下一天
	} else if nowHour == hour && nowMinute >= minute {
		now = now.AddDate(0, 0, 1) ///下一天
	}
	newTime := time.Date(now.Year(), now.Month(), now.Day(),
		hour, minute, 0, 0, now.Location())
	return newTime
}

func GetHourＭinUTC(hour int, minute int) int { ///得到指定小时和分钟的utc时间秒
	now := time.Now()
	nowHour := now.Hour()
	nowMinute := now.Minute()
	if nowHour >= hour && nowMinute >= minute {
		now = now.AddDate(0, 0, 1) ///下一天
	}
	newTime := time.Date(now.Year(), now.Month(), now.Day(),
		hour, minute, 0, 0, now.Location())
	return int(newTime.Unix())
}

func GetHourUTC(hour int) int { ///得到指定小时的utc时间秒
	now := time.Now()
	nowHour := now.Hour()
	if nowHour >= hour {
		now = now.AddDate(0, 0, 1) ///下一天
	}
	newTime := time.Date(now.Year(), now.Month(), now.Day(),
		hour, 0, 0, 0, now.Location())
	return int(newTime.Unix())
}

func CloneType(obj interface{}) interface{} { ///克隆对象,只有类型,没有值,值是空
	newObj := reflect.New(reflect.TypeOf(obj).Elem()).Elem()
	return newObj.Addr().Interface()
}

func Now() int { ///得到当前utc时间秒
	return int(time.Now().Unix())
}

func NowMS() int64 {
	nowMS := time.Now().UnixNano() / 1000000
	return nowMS
}

func Min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func Min_float(a float32, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
}

func Max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.Trim(lines[n], " \t")
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func SwapInt(value1, value2 *int) {
	tmp := *value1
	*value1 = *value2
	*value2 = tmp
}

func Random(min, max int) int { ///[min,max]
	randRate := max - min + 1
	if randRate <= 0 {
		return 0
	}
	rand.Seed(time.Now().UnixNano())
	result := rand.Intn(randRate) + min
	return result
}

func myStack(keepStackNum int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	stackNum := 0
	for i := 4; ; i++ { // Caller we care about is the user, 2 frames up
		if stackNum >= keepStackNum {
			break ///到达所取指定帧数
		}
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		line-- // in stack trace, lines are 1-indexed but our array is 0-indexed
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
		stackNum++ ///得到一帧
	}
	return buf.Bytes()
}

func SeparateIntList(valueList int) (IntList, IntList) {
	firstList, secList := IntList{}, IntList{}
	trainTargetString := fmt.Sprintf("%d", valueList) ///转成串处理
	for i, v := range trainTargetString {
		n := int(v) - 48 ///减0的ascii值
		if (i+1)%2 == 0 {
			secList = append(secList, n)
		} else {
			firstList = append(firstList, n)
		}
	}
	return firstList, secList
}

///overNum为克制系统数
func CalcGoalCount(attackTurns int, attackScore float32, defenseScore float32, overNum float32) (int, int) { ///计算进球数
	GetServer().GetLoger().CYDebug("overNum = %f", overNum)

	goalCount := 0
	goalRate_1 := int(attackScore * 100 / (attackScore + defenseScore))
	goalRate := int(attackScore*100/(attackScore+defenseScore) + (overNum * 100)) ///由数值策划林之冠要求更改

	GetServer().GetLoger().CYDebug("goalRate_1 = %d, goalRate = %d", goalRate_1, goalRate)

	goalRate = Min(goalRate, 100)       ///不得超出100
	for i := 1; i <= attackTurns; i++ { ///计算玩家进球数
		randRate := rand.Intn(100)
		if randRate <= goalRate {
			goalCount++ ///进一球
		}
	}
	return goalCount, goalRate
}

func TestMask(mask int, bit int) bool { ///测试掩码

	bit = Max(0, bit-1)

	bitValue := 1 << uint(bit)
	result := mask & bitValue

	return result > 0
}

func SetMask(mask int, bit int, value int) int { ///设置掩码
	bit = Max(0, bit-1)
	result := mask
	if value >= 1 {
		result |= 1 << uint(bit)
	} else {
		result |= bit << uint(0)
	}
	return result
}

func SetMask64(mask int64, bit int, value int) int64 { ///设置掩码
	bit64 := int64(Max(0, bit-1))
	result := mask
	testcode := int64(1)
	if value >= 1 {
		result |= testcode << uint(bit64)
	} else {
		result |= bit64 << uint(0)
	}
	return result
}

func TestMask64(mask int64, bit int) bool { ///测试掩码
	bit = Max(0, bit-1)
	testcode := int64(1)
	bitValue := testcode << uint(bit)
	result := mask & bitValue
	return result > 0
}

func round(val float64, prec int) int {

	var rounder float64
	intermed := val * math.Pow(10, float64(prec))

	if val >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return int(rounder / math.Pow(10, float64(prec)))
}

func GetOpenServerSamsaraDay(samsara int) int { //得到从开服时间起为轮回第几天
	if samsara <= 0 {
		return 0
	}
	currentTime := time.Now()

	openServerTime := GetServer().config.OpenServerTime

	opentime, _ := time.ParseInLocation(TimeFormat, openServerTime, time.Local)

	opentime = time.Date(opentime.Year(), opentime.Month(), opentime.Day(),
		0, 0, 0, 0, currentTime.Location())

	duration := currentTime.Sub(opentime)

	value1 := (int(duration.Seconds())) / (24 * 60 * 60)
	value2 := int(value1) % samsara
	//currentDay := () % samsara
	return value2
}

func GetOpenServerSamsaraDay_Plus(samsara int) int { //得到从开服时间起为轮回第几天(以4点钟计算)
	if samsara <= 0 {
		return 0
	}
	currentTime := time.Now()

	openServerTime := GetServer().config.OpenServerTime

	opentime, _ := time.ParseInLocation(TimeFormat, openServerTime, time.Local)

	opentime = time.Date(opentime.Year(), opentime.Month(), opentime.Day(),
		0, 0, 0, 0, currentTime.Location())

	duration := currentTime.Sub(opentime)

	value1 := (int(duration.Seconds()) - 4*60*60) / (24 * 60 * 60)
	value2 := int(value1) % samsara

	return value2
}

func TakeAccuracy(value float32, prec int) float32 { //取精度 prec为取几位小数
	multiple := float32(math.Pow(10, float64(prec+1)))
	rounder := value * float32(multiple)
	if int(rounder)%10 >= 5 {
		rounder = (rounder + 10) / multiple
	} else {
		rounder = rounder / multiple
	}

	return rounder
}

func SubString(inputStr string, start int, end int) string {
	strLen := len(inputStr)
	if start < 0 || start >= strLen {
		return ""
	}
	if end < 0 || end > strLen {
		return ""
	}
	resultStr := fmt.Sprint(inputStr[start:end])
	return resultStr
}

///抽一张卡,isDraw为true时抽选,false为展示
func discoverDrawOne(drawGroupList *IntList, totalTakeWeight *int, totalShowWeight *int, forceStarType int, isDraw bool) (int, int) {
	staticDataMgr := GetServer().GetStaticDataMgr()
	resultDrawOne := 0
	resultDrawOneNum := 0
	isHitDraw := false                          ///随机命中
	diceDraw := rand.Intn(*totalShowWeight + 1) ///生成展示抽取随机值
	if true == isDraw {
		diceDraw = rand.Intn(*totalTakeWeight + 1) //生成抽取随机值
	}

	for i := range *drawGroupList {
		drawGroupIndex := (*drawGroupList)[i]
		drawGroupStaticData := staticDataMgr.GetDrawGroupStaticData(drawGroupIndex)
		if nil == drawGroupStaticData {
			GetServer().GetLoger().Warn("StarSpyDiscoverMsg discoverDrawOne fail! drawGroupIndex %d is Invalid!", drawGroupIndex)
			break
		}
		if true == isDraw {
			isHitDraw = diceDraw <= drawGroupStaticData.TakeWeight ///命中判断
			diceDraw -= drawGroupStaticData.TakeWeight             ///扣除抽取权重
		} else {
			isHitDraw = diceDraw <= drawGroupStaticData.ShowWeight ///命中判断
			diceDraw -= drawGroupStaticData.ShowWeight             ///扣除展示权重
		}
		if forceStarType == drawGroupStaticData.AwardType {
			isHitDraw = true ///强制命中
		}
		if true == isHitDraw { ///抽中处理 (s[:i], s[i+1:]...)
			resultDrawOne = drawGroupStaticData.AwardType
			resultDrawOneNum = drawGroupStaticData.AwardCount
			*drawGroupList = append((*drawGroupList)[:i], (*drawGroupList)[i+1:]...) ///去掉被抽中的项
			*totalTakeWeight -= drawGroupStaticData.TakeWeight                       ///去掉已抽中抽取权重
			*totalShowWeight -= drawGroupStaticData.ShowWeight                       ///去掉已抽中展示权重
			break
		}
	}
	return resultDrawOne, resultDrawOneNum
}

func discoverGetDrawWeightTotal(drawGroupList []int) (int, int) { ///计算抽卡权重总和,抽取权重总和,展示权重总和
	staticDataMgr := GetServer().GetStaticDataMgr()
	totalTakeWeight := 0 ///抽取权重总和
	totalShowWeight := 0 ///显示权重总和
	for i := range drawGroupList {
		drawGroupIndex := drawGroupList[i]
		drawGroupStaticData := staticDataMgr.GetDrawGroupStaticData(drawGroupIndex)
		if nil == drawGroupStaticData {
			GetServer().GetLoger().Warn("StarSpyDiscoverMsg discoverGetDrawWeightTotal fail! drawGroupIndex %d is Invalid!", drawGroupIndex)
			break
		}
		///计算权重总和
		totalTakeWeight += drawGroupStaticData.TakeWeight ///抽取权重
		totalShowWeight += drawGroupStaticData.ShowWeight ///展示权重
	}
	return totalTakeWeight, totalShowWeight
}

///比较两个版本号,0表示完成相等,>0表示目标版本号大于源版本号,<0表示目标版本号小于源版本号
func CompareVersion(srcVersion string, dstVersion string) int {
	result := 0
	if srcVersion == dstVersion {
		return 0 ///完全相同
	}
	srcVersionNumber := strings.Replace(srcVersion, ".", "", -1)
	dstVersionNumber := strings.Replace(dstVersion, ".", "", -1)
	srcNum, _ := strconv.Atoi(srcVersionNumber)
	dstNum, _ := strconv.Atoi(dstVersionNumber)
	result = dstNum - srcNum
	return result
}

func RandMatchResult(resultList [][2]int) (int, int) { ///随机生成一个输了的比分
	resultLen := len(resultList)
	if resultLen <= 0 {
		return 0, 0
	}
	randIndex := Random(0, resultLen-1)
	a := resultList[randIndex][0]
	b := resultList[randIndex][1]
	return a, b
}

func EncodeValue(inputStr string) string {
	rst := ""
	for _, v := range inputStr {
		c := string(v)
		isMatch, _ := regexp.Match("[a-zA-Z0-9!()*]{1,1}", []byte(c))
		if isMatch {
			rst += c
		} else {
			rst += fmt.Sprintf("%%%02X", c)
		}
	}
	return rst
}
