package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/**
题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ：原子操作、并发数据安全。
*/

func main() {
	wg := sync.WaitGroup{}
	wg.Add(10)

	var count int64 = 0

	for i := 0; i < 10; i++ {
		go func() {
			for i := 0; i < 1000; i++ {
				atomic.AddInt64(&count, 1)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(count)

}
