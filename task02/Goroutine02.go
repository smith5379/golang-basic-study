package main

import (
	"fmt"
	"sync"
	"time"
)

type TaskScheduler struct {
}

/*
题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
考察点 ：协程原理、并发任务调度。
*/
func (t *TaskScheduler) executeTasks(funcs []func()) {
	wg := sync.WaitGroup{}
	wg.Add(len(funcs))
	for _, fn := range funcs {
		go func() {
			fmt.Println("开始执行任务: ", fn)
			start := time.Now()
			fn()
			fmt.Println("任务耗时: ", time.Since(start))
			wg.Done()
		}()
	}
	wg.Wait()
}

func main() {

	ts := &TaskScheduler{}

	funcs := []func(){}
	for i := 0; i < 10; i++ {
		funcs = append(funcs, func() {
			fmt.Println("执行任务 ", i)
		})
	}

	ts.executeTasks(funcs)

}
