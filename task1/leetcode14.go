package task1

import "fmt"

func main() {

	strs := []string{
		"flower",
		"flow",
		"flight",
	}
	fmt.Println(longestCommonPrefix(strs))
}

// {[()]}
// 遍历字符串：
/**
遍历字符串：
如果当前字符是左括号，入栈，
如果是
	右括号，再看栈顶元素，如果为空 或 不是匹配的左括号，直接return false
继续遍历，最后return true
*/

/*
*
最长公共前缀
*/
func longestCommonPrefix(strs []string) string {

	longestStr := strs[0]
	for i := 1; i < len(strs); i++ {
		currentStr := strs[i]
		maxMatchLength := getMaxMatchLength(longestStr, currentStr)

		if maxMatchLength == 0 {
			return ""
		}

		longestStr = longestStr[:maxMatchLength]
	}

	return longestStr
}

func getMaxMatchLength(str1 string, str2 string) int {
	minLength := min(len(str1), len(str2))
	var res int
	for i := 0; i < minLength; i++ {
		if str1[i] == str2[i] {
			res++
		} else {
			return res
		}
	}
	return res
}
