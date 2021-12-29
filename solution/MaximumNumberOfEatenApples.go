package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println("expected= 4 output=", eatenApples([]int{2, 1, 10}, []int{2, 10, 1}))
	fmt.Println("expected= 7 output=", eatenApples([]int{1, 2, 3, 5, 2}, []int{3, 2, 1, 4, 2}))
	fmt.Println("expected= 5 output=", eatenApples([]int{3, 0, 0, 0, 0, 2}, []int{3, 0, 0, 0, 0, 2}))
}

//贪心
func eatenApples(apples []int, days []int) int {
	//吃苹果数
	eaten := 0
	//当前第n天
	day_idx := 0
	//苹果篮子
	appleMap := make(map[int]int)
	appleList := []int{}
	//
	app_len := len(apples)
	for {
		day_idx++
		fmt.Println()
		day_app_len := 0
		//标识当天苹果的腐烂日期，将苹果push到篮子里
		if app_len >= day_idx && 0 != apples[day_idx-1] {
			day_app_len = apples[day_idx-1]
			day_at := day_idx + days[day_idx-1] - 1
			_, ok := appleMap[day_at]
			if ok {
				appleMap[day_at] += day_app_len
			} else {
				appleMap[day_at] = day_app_len
			}

			appleList = append(appleList, day_at)

			sort.Slice(appleList, func(i, j int) bool {
				return appleList[i] < appleList[j]
			})
			fmt.Println("第", day_idx, "天长了", day_app_len, "个可以保存到第", day_at, "天的苹果")
		}

		fmt.Println("第", day_idx, "天吃苹果前篮子还有苹果数：", appleMap, "; 序列：", appleList)

		//从篮子取一个腐烂日期最小的苹果
		count := len(appleList)
		del_at := -1
		for i := 0; i < count; i++ {
			day_at := appleList[i]
			num, ok := appleMap[day_at]
			//该天是否有苹果，List存在重复天就会出现没有苹果的情况
			if !ok {
				del_at = i
				continue
			}
			//是否腐烂了
			if day_idx > day_at {
				fmt.Println("第", day_idx, "天，可以保存到第", day_at, "天的苹果腐烂了", num, "个")
				delete(appleMap, day_at)
				del_at = i
				continue
			}
			//吃一个苹果
			num--
			eaten++
			fmt.Println("第", day_idx, "天吃可以保存到第", day_at, "天的苹果，可以保存到第", day_at, "天的苹果剩余", num, "个")
			//如果该天没有苹果了则删除该天
			if num <= 0 {
				delete(appleMap, day_at)
				del_at = i
				break
			}
			appleMap[day_at] = num
			break
		}
		if del_at+1 >= count {
			appleList = []int{}
		} else if del_at > -1 {
			appleList = appleList[del_at+1:]
		}
		fmt.Println("第", day_idx, "天吃苹果后篮子还有苹果数：", appleMap, "; 序列：", appleList)

		if len(appleMap) < 1 && day_idx >= app_len {
			break
		}
	}
	return eaten
}
