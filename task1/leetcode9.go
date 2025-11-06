package task1

import (
	"fmt"
	"strconv"
)

func main() {
	res := isPalindrome(121232121)
	fmt.Println("结果是:", res)
}

/**
 * 回文数
 */
func isPalindrome(x int) bool {
	str := strconv.Itoa(x)

	start := 0
	end := len(str) - 1
	for start < end {
		if str[start] != str[end] {
			return false
		}
		start++
		end--
	}

	return true
}
