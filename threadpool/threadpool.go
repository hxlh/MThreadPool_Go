package threadpool

import (
	"fmt"
	"sync"
	"time"
)

//任务类
type task struct {
	ff   func(args interface{}, ch chan int) interface{}
	args interface{}
}

//执行任务
func (t *task) exec(args interface{}, chFunc chan int) {
	t.ff(args, chFunc)
}

//ThreadPool 线程池
type ThreadPool struct {
	//唤醒线程的通道
	chThread chan int
	chFunc   chan int
	//全局锁
	locker *sync.RWMutex
	//总线程数
	threadNum int
	//关闭标志
	shutdown bool
	//任务列表
	tasks []task
}

/*AddTask ：添加任务
参数一：回到函数
参数二：回调函数的参数
*/
func (pool *ThreadPool) AddTask(f func(args interface{}, ch chan int) interface{}, arg interface{}) {
	tt := task{}
	tt.ff = f
	tt.args = arg
	pool.locker.Lock()
	pool.tasks = append(pool.tasks, tt)
	pool.locker.Unlock()
	fmt.Println("add task func call")
	pool.chThread <- 1

}

//Run 运行线程池
func (pool *ThreadPool) Run(threadNum int, taskNum int) {
	pool.chThread = make(chan int, taskNum)
	pool.chFunc = make(chan int, taskNum)
	pool.shutdown = false
	pool.threadNum = 0
	pool.locker = new(sync.RWMutex)

	pool.tasks = []task{}
	for i := 0; i < threadNum; i++ {
		go pool.threadRunFunc()
		pool.threadNum++
	}
}

//线程内部执行的函数
func (pool *ThreadPool) threadRunFunc() {
	for {
		<-pool.chThread
		pool.locker.Lock()
		if pool.shutdown {
			break
		}
		taskLen := len(pool.tasks)
		if taskLen <= 0 {
			continue
		}
		first := pool.tasks[0]
		pool.tasks = pool.tasks[1:]
		pool.locker.Unlock()
		first.exec(first.args, pool.chFunc)
	}

	pool.threadNum--
	pool.locker.Unlock()
}

//BadClose ：立即线程池
func (pool *ThreadPool) BadClose() {
	go func() {
		pool.locker.Lock()

		num := pool.threadNum
		//清空tasks
		pool.tasks = pool.tasks[:0:0]
		pool.shutdown = true
		pool.locker.Unlock()
		//发送关闭信号
		for i := 0; i < num; i++ {
			pool.chThread <- 1
		}
		close(pool.chThread)
	}()
}

//BadClose2 ：立即线程池
func (pool *ThreadPool) BadClose2() {
	go func() {
		pool.locker.Lock()

		num := pool.threadNum
		//清空tasks
		pool.tasks = pool.tasks[:0:0]
		pool.shutdown = true
		pool.locker.Unlock()
		//发送关闭信号
		for i := 0; i < num; i++ {
			//先退出任务
			pool.chFunc <- 1

		}
		for i := 0; i < num; i++ {
			//先退出任务
			pool.chThread <- 1

		}
		close(pool.chThread)
		close(pool.chFunc)
	}()
}

//FriendClose 轮询等待任务完成再关闭
func (pool *ThreadPool) FriendClose() {
	go func() {
		//循环检测任务是否完成
		for {
			pool.locker.Lock()
			taskNum := len(pool.tasks)
			threadNum := pool.threadNum
			pool.locker.Unlock()
			if taskNum <= 0 {
				//设置结束标识符
				pool.shutdown = true
				//发送关闭channal
				for i := 0; i < threadNum; i++ {
					pool.chThread <- 1
				}
				close(pool.chThread)
				break
			}
			time.Sleep(time.Second * 1)
		}

	}()
}

//测试用
func (pool *ThreadPool) Get() int {
	fmt.Printf("threadNum=%v\n", pool.threadNum)
	fmt.Printf("Tasks=%v\n", pool.tasks)
	return 0
}
