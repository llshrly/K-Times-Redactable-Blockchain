/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Description:
 * @File:
 * @Version: 1.0.0
 * @Data:
 */

package num

func Max(list []int64) int64 {
	maxVal := int64(0)
	for i := 1; i < len(list); i++ {
		//从第二个 元素开始循环比较，如果发现有更大的，则交换
		if maxVal < list[i] {
			maxVal = list[i]
		}
	}
	return maxVal
}

func Min(list []int64) int64 {
	minVal := list[0]
	for i := 1; i < len(list); i++ {
		//从第二个 元素开始循环比较，如果发现有更大的，则交换
		if minVal > list[i] {
			minVal = list[i]
		}
	}
	return minVal
}

func mean(v []int64) float64 {
	var res float64 = 0
	var n int = len(v)
	for i := 0; i < n; i++ {
		res += float64(v[i])
	}
	return res / float64(n)
}

func Variance(v []int64) float64 {
	var res float64 = 0
	var m = mean(v)
	var n int = len(v)
	for i := 0; i < n; i++ {
		res += (float64(v[i]) - m) * (float64(v[i]) - m)
	}
	return res / float64(n-1)
}
