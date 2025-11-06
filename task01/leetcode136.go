package main

import "fmt"

/**
 * 只出现一次的数字
 */
func main() {
	res := singleNumber2([]int{1, 3, 1, 3, 4})
	fmt.Println("结果是:", res)
}

// hashmap统计法
func singleNumber(nums []int) int {
	var countMap = make(map[int]int)
	for _, value := range nums {
		countMap[value]++
	}

	for num, count := range countMap {
		if count == 1 {
			return num
		}
	}
	return -1
}

// 异或位运算法
func singleNumber2(nums []int) int {

	res := 0
	for _, value := range nums {
		res = res ^ value
	}
	return res
}
