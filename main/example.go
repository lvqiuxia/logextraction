
package main

import (
	"fmt"
)


func main() {
	fmt.Println("hello world")

	//防止主进程退出
	i := 1
	for{
		if i < 0 {
			return
		}
	}

	////防止主进程退出
	//ch := make(chan int )
	//go func(){
	//	<-ch
	//}()

}