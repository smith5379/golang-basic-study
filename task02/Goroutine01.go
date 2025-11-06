package main

import (
	"fmt"
	"sync"
)

/*
题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
考察点 ： go 关键字的使用、协程的并发执行。
*/
func main() {

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		i := 1
		for i <= 10 {
			fmt.Println("1 ~ 10 中的奇数: ", i)
			i = i + 2
		}
		wg.Done()
	}()

	go func() {
		i := 2
		for i <= 10 {
			fmt.Println("1 ~ 10 中的偶数: ", i)
			i = i + 2
		}
		wg.Done()
	}()

	wg.Wait()
}
