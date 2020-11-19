package main

import (
	"fmt"
	"thread_pool/threadpool"
	"time"
)

type Test struct {
}

func (tt *Test) test(args interface{}, ch chan int) interface{} {
	// time.Sleep(time.Millisecond * 1000)
	// fmt.Println(args)
	// return ""
	for {
		select {
		case <-ch:
			//do something…………
			//跳出循环，退出
			goto loop
		default:
			//do something…………
			time.Sleep(time.Millisecond * 1000)
			fmt.Println(args)

		}
	}
loop:

	return 0
}
func main() {
	tt := &Test{}
	pool := &threadpool.ThreadPool{}
	pool.Run(30, 100)
	for i := 0; i < 100; i++ {
		pool.AddTask(tt.test, i)
	}

	var str string
	fmt.Scan(&str)
	fmt.Println("回收call ")
	pool.BadClose2()

	var str2 string
	fmt.Scan(&str2)
	pool.Get()
	var str3 string
	fmt.Scan(&str3)
}
