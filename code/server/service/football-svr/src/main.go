package main

import (
	"code.google.com/p/go.net/websocket"
	//"encoding/json"
	"fmt"
	f "football"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func KillOldApp() {
	pidFileName := "./server.pid"                  ///服务器配置文件路径及文件名
	pidServer, err := ioutil.ReadFile(pidFileName) ///尝试打开配置文件
	pidString := fmt.Sprintf("%d", os.Getpid())
	if err != nil {
		err = ioutil.WriteFile(pidFileName, []byte(pidString), os.ModePerm)
		return
	}
	pidOldString := string(pidServer)
	pidOldString = strings.Replace(pidOldString, "\n", "", 1)
	pid, _ := strconv.Atoi(pidOldString)
	process, _ := os.FindProcess(pid)
	if process != nil {
		process.Kill()
	}
	err = ioutil.WriteFile(pidFileName, []byte(pidString), os.ModePerm)
}

func PreventDoubleStart(serverIndex int) { ///防止重复启动
	pidFileName := fmt.Sprintf("./server%d.pid", serverIndex) ///服务器配置文件路径及文件名
	pidServer, err := ioutil.ReadFile(pidFileName)            ///尝试打开配置文件
	pidString := fmt.Sprintf("%d", os.Getpid())
	if err != nil {
		err = ioutil.WriteFile(pidFileName, []byte(pidString), os.ModePerm)
		return
	}
	pidOldString := string(pidServer)
	pidOldString = strings.Replace(pidOldString, "\n", "", 1)
	pid, _ := strconv.Atoi(pidOldString)
	process, _ := os.FindProcess(pid)
	err = process.Signal(syscall.Signal(0))
	if err == nil {
		log.Fatalf("PreventDoubleStart serverIndex:%d has startup!", serverIndex)
	}
	err = ioutil.WriteFile(pidFileName, []byte(pidString), os.ModePerm)
}

///测试客户端，自连自测试
var origin = "http://192.168.20.125/"
var url = "ws://192.168.20.125:8080/"

func clientRobotSend(clientID int, sendMsgChannel chan f.IMsgHead, ws *websocket.Conn) {
	seqID := 1
	for {
		select {
		case msg := <-sendMsgChannel:
			msgType, msgAction := msg.GetTypeAndAction()
			msg.FillMsgHead(seqID, &msgType, &msgAction)
			//log.Println(msg)
			err := websocket.JSON.Send(ws, msg)
			seqID++
			if err != nil {
				log.Println("clientRobotSend get a err!", err)
				os.Exit(1)
				//break
			}
			log.Printf("client %d start send a msg...ok", clientID)
			//time.Sleep(time.Millisecond * 5000)
			//os.Exit(1)
		}
	}
}

func clientRobotReceive(clientID int, sendMsgChannel chan f.IMsgHead, ws *websocket.Conn) {
	msg := ""
	for {
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			log.Println("clientRobotReceive get a err!", err)
			break
		}
		log.Printf("client %d get a msg:%v\n", clientID, msg)
	}
}

func TestClientGo() {
	time.Sleep(time.Second * 1)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Println(err)
		return
	}
	clientID := 1
	seqID := 1 ///消息顺序号
	//	randNum := rand.Intn(1) + 1
	//	randTime := time.Duration(randNum)
	checkLifeTimer := time.NewTicker(time.Second * 1)
	//clientName := fmt.Sprintf("client%d", clientID)

	//msg := ClientLoginMsg{}
	//msg.UserName = "go03"
	//msg.PassWord = "1234"
	//msg.ClientVersion = "1.2.3.4"
	//msg.MsgType = "login" ///聊天消息
	//msg.Action = "clientlogin"
	//msg.SeqID = seqID
	//msg.CreateTime = time.Now().Unix()//TestClient()

	sendMsgChannel := make(chan f.IMsgHead)
	go clientRobotSend(clientID, sendMsgChannel, ws)    ///发送协程
	go clientRobotReceive(clientID, sendMsgChannel, ws) ///收接协程
	for {
		msg := new(f.ClientLoginMsg)
		msg.UserName = "cs88"
		msg.PassWord = "123456"
		msg.ClientVersion = "1.2.3.4"
		msg.MsgType = "login" ///聊天消息
		msg.Action = "clientlogin"

		msg2 := new(f.SDKLoginMsg)
		msg2.UserName = "98153038195b3ae27a33d9c0fdb6ee76894a092d5cb1c2c50"
		msg2.SDKName = "Qihoo360"
		//msg2 := new(f.StarSpyDiscoverMsg)
		//msg2.StarSpyType = 1

		//msg2 := new(f.QueryLevelListMsg)
		//msg2.LeagueType = 1

		//msg2 := new(f.PassLevelMsg)
		//msg2.LevelSort = 3
		//msg2.LevelType = 1301
		//msg2.LevelID = 0
		//msg2.Param1 = 101011

		//msg2 := new(f.PassLevelMsg)
		//msg2.LevelSort = 1
		//msg2.LevelType = 1101
		//msg2.LevelID = 0
		//msg2.Param1 = 101011

		//msg2 := new(f.StarSpyDiscoverMsg)
		//msg2.StarSpyType = 1

		//msg2 := new(QueryItemListMsg)
		//msg2.TeamID = 4568

		//msg2 := new(MergeItemMsg)
		//msg2.MasterItemID = 1Done()
		//msg2.MergeItemIDList = IntList{8}

		//msg2 := new(EquipItemMsg)
		//msg2.StarID = 762
		//msg2.EquipItemID = 1
		//msg2.UnequipItemID = 0
		//msg2.EquipType = equipItemOPTypeWield
		//msg2.ItemSort = ItemSortCloth

		//msg2 := new(f.StarCenterTransferMsg)
		//msg2.StarCenterType = 1
		//msg2.MemberID = 407
		//msg2.ExchangeStarList = []int{}

		//msg3 := new(EquipItemMsg)
		//msg3.StarID = 762
		//msg3.EquipItemID = 2
		//msg3.UnequipItemID = 1
		//msg3.EquipType = equipItemOPTypeReplace
		//msg3.ItemSort = ItemSortCloth

		//msg2 := new(f.TeamCreateMsg)
		//msg2.MsgType = "team"
		//msg2.Action = "createteam"
		//msg2.AccountID = 115
		//msg2.TeamName = "bobo7"
		//msg2.Icon = 1
		//msg2.TeamShirts = 2
		//msg2.StarTypeList = f.IntList{3188, 1091, 1092, 1048, 1070, 1073, 1096, 1081, 1097, 1056, 1098}

		//msg2 := new(QueryVolunteerInfoMsg)
		//msg2.MsgType = "starcenter"
		//msg2.Action = "queryvolunteerinfo"
		//msg2.IsRefreshStarList = false

		//msg3 := new(VolunteerSignMsg)
		//msg3.MsgType = "starcenter" ///聊天消息
		//msg3.Action = "volunteersign"
		//msg3.StarType = 3039

		//msg2 := new(QueryStarCenterMemberListMsg)
		//msg2.MsgType = "starcenter" ///聊天消息
		//msg2.Action = "querystarcentermemberlist"
		//msg2.StarCenterType = 1

		//msg2 := new(f.StarCenterTransferMsg)
		//msg2.StarCenterType = 1
		//msg2.MemberID = 132
		//msg2.ExchangeStarList = []int{399}

		<-checkLifeTimer.C
		//jmsg, err := json.Marshal(msg)
		//log.Println(string(jmsg), err)
		if seqID <= 1 {
			sendMsgChannel <- msg
			sendMsgChannel <- msg2
			//sendMsgChannel <- msg3
		}
		seqID++
		ws.Close()
	}
}

func TestClient() {
	go TestClientGo()
}

//得到随机星级
//func CalcStarTenDrawStarCount() int {

//	///简化公式3^(7-m)/40  m = 星级 星级上限为7 m >= 4
//	///概率 0.675 0.225 0.075 0.025
//	iProbabilityList := []int{669, 225, 74, 24, 8}
//	nTemp := 0
//	randNum := f.Random(0, 1000)
//	randEvolove := 3
//	for i := 3; i <= 7; i++ {
//		if randNum >= nTemp && randNum < nTemp+iProbabilityList[i-3] {
//			randEvolove = i
//			break
//		}
//		nTemp += iProbabilityList[i-3]
//	}
//	return randEvolove
//}

func main() {

	//s1 := f.GetServer()
	//s1.GetLoger()
	//item := foo1()
	//s := foo2()
	//item.GetInfo()
	//s.GetLoger()
	//fmt.Printf("%s\r\n", f.ServerVersion)
	//runtime.GOMAXPROCS(1)
	//f.SeparateIntList(43422145)
	//for i := 0; i < 10; i++ {
	//	n := CalcStarTenDrawStarCount()
	//	fmt.Println(n)
	//}
	//err := f.Mkfifo("./a1", 0666)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//err = f.Mkfifo("./a1", 0666)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//if len(os.Args) < 2 {
	//	log.Fatal("server stop!cause need give serverindex!")
	//}
	serverIndex, _ := strconv.Atoi(os.Args[1])
	PreventDoubleStart(serverIndex)
	f.RunGame(serverIndex)
	//	KillOldApp()
	//	TestClient()
	//
}
