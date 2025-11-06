package main

import (
	"fmt"
	"time"
)

/**
题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
考察点 ：通道的缓冲机制。
*/

func main() {

	ch := make(chan int, 10)

	go func() {
		for i := 1; i <= 100; i++ {
			ch <- i
			fmt.Println("发送整数", i, "到channel中 ")
		}
		close(ch)
	}()

	go func() {
		for {
			i, ok := <-ch
			if !ok {
				fmt.Println("通道已关闭")
				break
			} else {
				fmt.Println("从channel中接收到 : ", i)
			}
		}

	}()

	time.Sleep(10 * time.Second)

}
