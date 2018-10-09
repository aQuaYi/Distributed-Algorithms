package mutualexclusion

import (
	"fmt"
	"time"
)

func main() {
	beginTime := time.Now()
	count := 0
	amount := 131072 // NOTICE: 为了保证测试结果的可比性，请勿修改此数值
	for all := 2; all <= 128; all *= 2 {
		times := amount / all
		fmt.Printf("~~~ %d Process，每个占用资源 %d 次，共计 %d 次 ~~~\n", all, times, amount)
		newRound(all, times)
		count++
	}

	fmt.Printf("一共测试了 %d 轮，全部通过。共耗时 %s 。\n", count, time.Now().Sub(beginTime))
}
