package main

import "fmt"

func main() {
	fmt.Println("expected= 4 output=", eatenApples([]int{2, 1, 10}, []int{2, 10, 1}))
	fmt.Println("expected= 7 output=", eatenApples([]int{1, 2, 3, 5, 2}, []int{3, 2, 1, 4, 2}))
	fmt.Println("expected= 5 output=", eatenApples([]int{3, 0, 0, 0, 0, 2}, []int{3, 0, 0, 0, 0, 2}))
}
func eatenApples(apples []int, days []int) int {
	//循环apples数组
	//第i天是0的元素，跳过
	//第i天不是0的元素，轮询当天所有苹果
	//判断当天的苹果是否腐烂
	//不腐烂当天苹果数-1，天+1
	//腐烂天+1，跳过第i天苹果
	rst := 0
	//当前第n天
	day_idx := 0
	apples_len := len(apples)
	//第i天出产的苹果
	for i := 0; i < apples_len; i++ {
		//第i天有多少个苹果
		day_app_len := apples[i]
		//第i天没有苹果跳过
		if 0 == day_app_len {
			//当天没有苹果吃
			if i >= day_idx {
				day_idx++
			}
			continue
		}
		max := (i + days[i])
		//轮询第i天的所有苹果
		for j := 0; j < day_app_len; j++ {
			//当天的苹果腐烂
			if max < day_idx {
				day_idx++
				break
			}
			if max == day_idx {
				break
			}
			day_idx++
			rst++
		}
	}
	return rst
}
