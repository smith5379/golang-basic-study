package main

/*
*
加一
*/
func plusOne(digits []int) []int {
	// 判断最后一位是不是9，如果不是，最后一位直接加1，
	// 如果是，当前位置置为0，前面一个加一
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] != 9 {
			digits[i]++
			return digits
		} else {
			digits[i] = 0
		}
	}

	digits = append([]int{1}, digits...)
	return digits
}
