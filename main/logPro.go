package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func main(){

	//读取文件内容
	path := "result.log"
	fileMsg, err := ioutil.ReadFile(path)
	if err != nil{
		fmt.Println(err)
		panic("ReadFile failed")
	}
	//fmt.Println(string(fileMsg))

	//将读取的内容按空行处理，且只取2：end作为一个字符串放进切片里面
	var str bytes.Buffer
	strSlice := make([]string,0)
	spMsg := strings.Split(string(fileMsg),"\n")
	for line, newMsg := range spMsg {
		//fmt.Println(newMsg,line)
		if line == 0 && newMsg == ""{
			continue
		}

		if newMsg != "" && spMsg[line+1] != "" && spMsg[line-1] != ""{
			str.WriteString(newMsg)
			//fmt.Println(spMsg[line+1])
		}
		if newMsg != "" && spMsg[line+1] == "" && spMsg[line-1] != ""{
			str.WriteString(newMsg)
		}
		if newMsg == "" && line > 0{
			//fmt.Println(str.String())
			strSlice = append(strSlice,str.String())
			str.Reset()
		}
	}
	sort.Strings(strSlice)
	//以读写方式打开文件，如果不存在，则创建
	file1, error := os.OpenFile("./result_all.log", os.O_RDWR|os.O_CREATE, 0766)
	if error != nil {
		fmt.Println(error)
	}

	for {
		var strLen int
		for i := 0; i < len(strSlice); i ++ {
			if len(strSlice[i]) == 0 {
				strSlice = append(strSlice[:i], strSlice[i+1:]...)
			}else{
				strLen = strLen + 1
			}
		}
		if strLen == len(strSlice){
			break
		}
	}

	for i := 0; i < len(strSlice); i ++ {
		file1.WriteString(strSlice[i])
		file1.WriteString("\n")
	}

	file1.Close()

	result := RemoveRepByLoop(strSlice)
	sort.Strings(result)

	//删除相似的字符串
	//result = DeleteSimilarStrings(result)
	//for num := range result{
	//	fmt.Println(result[num])
	//}

	//循环删除相似的字符串，直到不能再删除为止
	for{
		beforLen := len(result)
		//fmt.Println("beforLen",beforLen)
		result = DeleteSimilarStrings(result,1)
		afterLen := len(result)
		//fmt.Println("afterLen",afterLen)
		if afterLen == beforLen{
			break
		}
	}

	//多次去重之后，字符串中含两个不同的项也删除
	//result = DeleteSimilarStrings(result,3)
	for{
		beforLen := len(result)
		//fmt.Println("beforLen",beforLen)
		result = DeleteSimilarStrings(result,2)
		afterLen := len(result)
		//fmt.Println("afterLen",afterLen)
		if afterLen == beforLen{
			break
		}
	}

	file2, error := os.OpenFile("./result_all.log", os.O_RDWR|os.O_CREATE, 0766)
	if error != nil {
		fmt.Println(error)
	}

	for num := range result{
		file2.WriteString(result[num])
		file2.WriteString("\n")
		fmt.Println(result[num])
	}
	file2.Close()

}


//删除相似的字符串
func DeleteSimilarStrings(result []string,num uint32)[]string{
	//循环处理每一个字符串
	totalLen := len(result)
	//fmt.Println(totalLen)
	totalArr := make([][]string,0)
	for _,str := range result{
		//将每一个字符串按照空格分为数组
		array := SplitStringToArray(str)
		totalArr = append(totalArr,array)
		//for j := i; j < totalLen; j ++{
		//	reqMsg := result[j]
		//	flag := CompareStringWithOneElement(array,reqMsg)
		//	if flag == true{
		//		result = append(result[:j],result[j+1:]...)
		//		totalLen = totalLen-1
		//	}
		//}
	}
	if totalArr[0][0] == " "{
		totalArr = append(totalArr[:0],totalArr[1:]...)
		result = append(result[:0],result[1:]...)
	}

	//fmt.Println(totalArr)
	//fmt.Println(result)
	totalLen = len(result)
	//将当前字符串的每个数组元素与后面的每个字符串的数组元素比较
	//如果后面的字符串与当前的字符串只有一个数组元素不一致，则去除掉
	for i := 0; i < totalLen; i ++{
		inum := num
		msg := totalArr[i]
		//todo:*******************************************************
		if len(totalArr[i]) < 5{
			inum = 2
		}
		if len(totalArr[i]) < 3{
			inum = 1
		}
		//todo:*******************************************************
		//fmt.Println(msg)
		for j := i+1; j < totalLen; j ++{
			reqMsg := totalArr[j]
			//fmt.Println(reqMsg)
			flag := CompareStringWithOneElement(msg,result[i],reqMsg,result[j],inum)
			if flag == true{
				totalArr = append(totalArr[:j],totalArr[j+1:]...)
				result = append(result[:j],result[j+1:]...)
				totalLen = totalLen-1
			}
		}
	}

	sort.Strings(result)
	return result
}

// 通过两重循环过滤重复元素
func RemoveRepByLoop(slc []string) []string {
	result := []string{}  // 存放结果
	for i := range slc{
		flag := true
		for j := range result{
			if slc[i] == result[j] {
				flag = false  // 存在重复元素，标识为false
				break
			}
		}
		if flag {  // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}


// 通过两重循环过滤重复元素,包括过滤掉只变化数字的元素
func RemoveRepByLoopEvol(slc []string) []string {
	result := []string{}  // 存放结果
	for i := range slc{
		flag := true
		for j := range result{
			if slc[i] == result[j] {
				flag = false  // 存在重复元素，标识为false
				break
			}
		}
		if flag {  // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}


//将每一个字符串按照空格分为数组
func SplitStringToArray(str string) []string{
	reMsg := make([]string,0)
	for _,temp := range strings.Split(str," "){
		reMsg = append(reMsg,temp)
	}
	return reMsg
}

//比较两个数组的元素
func CompareStringWithOneElement(arra []string, stra string, arrb []string, strb string, num uint32) bool{
	var indexa uint32
	var indexb uint32
	for _,str := range arra{
		if strings.Contains(strb,str) == false{
			indexa ++
			//fmt.Println(indexa)
		}
	}

	for _,str := range arrb{
		if strings.Contains(stra,str) == false{
			indexb ++
			//fmt.Println(indexb)
		}
	}
	
	if indexa <= num && indexb <= num{
		return true
	}
	return false
}

