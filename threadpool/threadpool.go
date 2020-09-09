package threadpool

import (
	"sync"
)
//任务类
type task struct {
	ff   func(args interface{}) interface{}
	args interface{}
}
//执行任务
func (t *task) exec(args interface{}) {
	t.ff(args)
}
//ThreadPool 线程池
type ThreadPool struct {
	//唤醒线程的通道
	ch chan int
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
func (pool *ThreadPool) AddTask(f func(args interface{}) interface{}, arg interface{}) {
	tt := task{}
	tt.ff = f
	tt.args = arg
	pool.locker.Lock()
	pool.tasks = append(pool.tasks, tt)
	pool.locker.Unlock()
	pool.ch <- 1
}
//Run 运行线程池
func (pool *ThreadPool) Run(threadNum int) {
	pool.ch = make(chan int)
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
		<-pool.ch
		pool.locker.Lock()
		if pool.shutdown {
			break
		}
		first := pool.tasks[0]
		pool.tasks = pool.tasks[1:]
		pool.locker.Unlock()
		first.exec(first.args)
	}

	pool.threadNum--
	pool.locker.Unlock()
}
//Close ：关闭线程池
func (pool *ThreadPool) Close() {
	pool.locker.Lock()
	num := pool.threadNum
	pool.shutdown = true
	pool.locker.Unlock()
	//发送关闭信号
	for i := 0; i < num; i++ {
		pool.ch <- 1
	}
}

//测试用
// func (pool*ThreadPool)Get() int {
// 	return pool.threadNum
// }
