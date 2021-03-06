package main

import (
	"fmt"
	"time"
	"regexp"
	"io/ioutil"
	"encoding/json"
	"os"
	"log"
	"io"
)

//the struct of log
type InputMsgOfLog struct{
	NodeName     string
	Time         time.Time
	ProcessLine  string
	Num          uint32
	LogMsg       string
}

//json文件
type AlarmMsg struct{
	Msg string
}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

//当前log日志数
var LogNum uint32

//log日志的最大个数
var MaxLogNum uint32 = 3

//日志
var LogFile *os.File

func main(){


	//确定程序启动时从第几个文件开始写
	InitLogNum()

	//读取json文件
	JsonParse := NewJsonStruct()
	//v := AlarmMsg{}
	//var v []map [string] string
	var v []AlarmMsg
	//下面使用的是相对路径，config.json文件和main.go文件处于同一目录下
	JsonParse.Load("F://编程语言学习/Go语言/code-practice/src/PracticeCode/LogExtraction/main/msgCls.json", &v)
	fmt.Println("get json msg:",v)

	//reg1 := regexp.MustCompile("bigger than length")
	//reg2 := regexp.MustCompile("fail to attach")

	gzsca01 := regexp.MustCompile("gzsca01")
	gzsca02 := regexp.MustCompile("gzsca02")
	msg := GetLogMsg()

	for _,temp := range v{
		
		if temp.Msg == msg.LogMsg && gzsca01.MatchString(msg.NodeName) {
			fmt.Println("show is [gzsca01] msg",msg)
			err := WriteMsgToLog(msg)
			if err != nil{
				return
			}
		}

		if temp.Msg == msg.LogMsg && gzsca02.MatchString(msg.NodeName) {
			fmt.Println("show is [gzsca02] msg",msg)
			err := WriteMsgToLog(msg)
			if err != nil{
				return
			}
		}
	}
}

//获取日志内容
func GetLogMsg() InputMsgOfLog {
	logMsg := InputMsgOfLog{}


	logMsg.NodeName = "gzsca01"
	logMsg.Time = time.Now().In(time.FixedZone("CST", 8*60*60))
	logMsg.ProcessLine = "111:qtdubi"
	logMsg.Num = 1

	//var msg string
	//strMsg, _ := fmt.Scanf("%s",&msg)
	// logMsg.LogMsg = fmt.Sprintln(strMsg)


	strMsg := "fail to attach"
	logMsg.LogMsg = strMsg

	return logMsg
}


func (jst *JsonStruct) Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("get filename failed,err is :",err)
		return
	}

	fmt.Println("data is :",string(data),data)

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		fmt.Println("Unmarshal failed,err is :",err)
		return
	}
}


func InitLogNum(){
	//确定程序启动时从第几个文件开始写日志
	for temp := uint32(0); temp < MaxLogNum+1; temp ++ {
		path := fmt.Sprintf("./AllErrorMsg%d_%s.log",LogNum,time.Now().Format("2006-01-02"))
		//新建一个日志
		logFile, err := os.OpenFile(path , os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			return
		}
		logSize,_ := logFile.Seek(0, io.SeekEnd)
		if logSize > 1000 {
			LogNum ++
			logFile.Close()
		}else{
			break
		}
	}
}


//将内容写入日志
func WriteMsgToLog(inMsg InputMsgOfLog) error {

	//定义日志文件的路径和内容，坑，必须写2006-01-02才能正确的转换
	//FileName := "./"+"AllErrorMsg1_" + time.Now().Format("2006-01-02") + ".log"
	//fmt.Println("Filename is ",FileName)
	//FileName = fmt.Sprintln("aaaaaaa")
	//fmt.Printf("./"+"AllErrorMsg%d_" + time.Now().Format("2006-01-02") + ".log",1)
	//path,_ := fmt.Printf("./"+"AllErrorMsg%d_" + time.Now().Format("2006-01-02") + ".log",1)
	var err error

	//FileName := fmt.Sprintf("./AllErrorMsg%d_%s_%d:%d:%d.%d.log",LogNum,time.Now().Format("2006-01-02"),
	//	time.Now().Hour(),time.Now().Minute(),time.Now().Second(),time.Now().Nanosecond())

	FileName := fmt.Sprintf("./AllErrorMsg%d_%s.log",LogNum,time.Now().Format("2006-01-02"))
	fmt.Println("LogNum",LogNum)
	fmt.Println(FileName)

	//新建一个日志
	LogFile, err = os.OpenFile(FileName , os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		return  err
	}

	logSize,_ := LogFile.Seek(0, io.SeekEnd)
	fmt.Println(logSize)
	if logSize > 1000 {
		LogNum ++

		if LogNum > MaxLogNum{
			LogNum = 0
			FileName := fmt.Sprintf("./AllErrorMsg%d_%s.log",LogNum,time.Now().Format("2006-01-02"))
			_, err = os.Stat(FileName)
			if err == nil {
				os.Remove(FileName)
			}
		}

		fmt.Println("Size is big than 1000",LogNum)
		FileName1 := fmt.Sprintf("./AllErrorMsg%d_%s.log",LogNum,time.Now().Format("2006-01-02"))
		_, err = os.Stat(FileName1)
		if err == nil {
			os.Remove(FileName1)
		}
		LogFile.Close()

		LogFile, err = os.OpenFile(FileName1 , os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			return  err
		}
	}

	//创建一个Logger
	//	//参数1：日志写入目的地
	//	//参数2：每条日志的前缀
	//	//参数3：日志属性
	//loger := log.New(logFile, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	loger := log.New(LogFile, "[INFO] ", log.Ldate|log.Ltime)
	loger.Println("[",inMsg.NodeName,"]",inMsg.Time,"line:",inMsg.ProcessLine,"[",inMsg.LogMsg,"]")
	
	fmt.Println(LogNum)

	LogFile.Close()

	return nil
	
}






















//package main
//
//import (
//	"io/ioutil"
//	"encoding/json"
//	"fmt"
//)
////定义配置文件解析后的结构
//type MongoConfig struct {
//	MongoAddr      string
//	MongoPoolLimit int
//	MongoDb        string
//	MongoCol       string
//}
//
//type Config struct {
//	Addr  string
//	Mongo MongoConfig
//}
//
//func main() {
//	JsonParse := NewJsonStruct()
//	v := Config{}
//	//下面使用的是相对路径，config.json文件和main.go文件处于同一目录下
//	JsonParse.Load("F://编程语言学习/Go语言/code-practice/src/PracticeCode/LogExtraction/main/config.json", &v)
//	fmt.Println(v.Addr)
//	fmt.Println(v.Mongo.MongoDb)
//}
//
//type JsonStruct struct {
//}
//
//func NewJsonStruct() *JsonStruct {
//	return &JsonStruct{}
//}
//
//func (jst *JsonStruct) Load(filename string, v interface{}) {
//	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
//	data, err := ioutil.ReadFile(filename)
//	if err != nil {
//		return
//	}
//	fmt.Println("data is:",data)
//
//	//读取的数据为json格式，需要进行解码
//	err = json.Unmarshal(data, v)
//	if err != nil {
//		return
//	}
//}
