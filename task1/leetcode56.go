package task1

import "slices"

/*
*
合并区间
*/
func merge(intervals [][]int) (res [][]int) {

	slices.SortFunc(intervals, func(a, b []int) int {
		return a[0] - b[0]
	})

	for _, interval := range intervals {

		if len(res) == 0 || res[len(res)-1][1] < interval[0] { // 如果res为空，或者res中的最后一个区间的右边界小于 当前遍历到的左边界，不需要合并
			res = append(res, interval)
		} else {
			res[len(res)-1][1] = max(interval[1], res[len(res)-1][1]) // 更新最后一个元素的右边界
		}

	}
	return

}
