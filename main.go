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
	pool.Run(100)
	for i := 0; i < 100; i++ {
		pool.AddTask(test, i)
	}
	// time.Sleep(time.Second*3)
	go pool.Close()
	fmt.Println("回收成功")
	var str string
	fmt.Scan(&str)
	// fmt.Print("threadsNum=")
	// fmt.Println(pool.Get())
	// var str2 string
	// fmt.Scan(&str2)
}
