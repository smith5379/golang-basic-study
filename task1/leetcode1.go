package task1

import "fmt"

func main() {
	nums := []int{2, 7, 11, 15}
	fmt.Println(twoSum(nums, 9))
}

/**
 * 两数之和
 */
func twoSum(nums []int, target int) []int {
	mp := make(map[int]int)
	for i, num := range nums {
		value, ok := mp[target-num]
		if ok {
			return []int{value, i}
		} else {
			mp[num] = i
		}
	}
	return nil
}
