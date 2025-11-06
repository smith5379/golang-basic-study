package task1

import "fmt"

func main() {
	fmt.Println(isValid("]"))
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

/**
 * 有效的括号
 */
func isValid(s string) bool {
	// fast fail
	if len(s)%2 == 1 {
		return false
	}
	mp := map[rune]rune{
		'}': '{',
		')': '(',
		']': '[',
	}

	var arr []rune

	for _, c := range s {
		if mp[c] == 0 { //左括号入栈
			arr = append(arr, c)
		} else { //右括号 校验
			if len(arr) == 0 || arr[len(arr)-1] != mp[c] {
				return false
			}
			arr = arr[:len(arr)-1]
		}
	}

	return len(arr) == 0
}
