package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("修改前的nums = ", nums)
	multiTwo(nums)
	fmt.Println("修改后的nums = ", nums)
}

/*
题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
考察点 ：指针运算、切片操作。
*/
func multiTwo(nums []int) {
	for i, num := range nums {
		nums[i] = num * num
	}
}
