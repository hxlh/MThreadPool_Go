package main

import (
	"fmt"
	"thread_pool/threadpool"
	"time"
)

func test(args interface{}) interface{} {
	time.Sleep(time.Millisecond * 1000)
	fmt.Println(args)
	return ""
}
func main() {
	pool := &threadpool.ThreadPool{}
	pool.Run(30,100)
	for i := 0; i < 300; i++ {
		pool.AddTask(test, i)
	}
	
	pool.FriendClose()
	fmt.Println("回收call ")
	var str string
	fmt.Scan(&str)
	
	pool.Get()
	var str2 string
	fmt.Scan(&str2)
}
