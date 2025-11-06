package main

import "fmt"

func main() {
	nums := []int{1, 1, 2}
	fmt.Println(removeDuplicates(nums))
}

/*
*
删除有序数组中的重复项
*/
func removeDuplicates(nums []int) int {

	k := 1
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1] { // nums[i] 不是重复项
			nums[k] = nums[i] // 保留nums[i]
			k++
		}
	}

	return k
}
