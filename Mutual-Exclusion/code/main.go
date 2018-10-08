package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aQuaYi/observer"
)

var needDebug = false

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	debugPrintf("程序开始运行")
	rand.Seed(time.Now().UnixNano())
}

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

func newRound(all, occupyTimesPerProcess int) {
	rsc := newResource(all * occupyTimesPerProcess)

	prop := observer.NewProperty(nil)

	ps := make([]Process, all)
	// 需要一口气同时生成，保证所有的 stream 都能从同样的位置开始观察
	for i := range ps {
		p := newProcess(all, i, rsc, prop)
		ps[i] = p
	}
	debugPrintf("~~~ 已经成功创建了 %d 个 Process ~~~", all)

	stream := prop.Observe()
	go func() {
		for {
			msg := stream.Next().(*message)
			debugPrintf(" ## %s", msg)
		}
	}()

	for _, p := range ps {
		go func(p Process, times int) {
			i := 0
			debugPrintf("%s 开始申请资源", p)
			for i < times {
				p.Request()
				i++
			}
		}(p, occupyTimesPerProcess)
	}

	rsc.Wait()

	fmt.Println(rsc.Report())
}
