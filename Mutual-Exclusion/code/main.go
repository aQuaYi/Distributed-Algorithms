package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aQuaYi/observer"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	count := 0
	amount := 256
	for all := 2; all <= 128; all *= 2 {
		times := amount / all
		fmt.Printf("~~~ %d Process，每个占用资源 %d 次 ~~~\n", all, times)
		newRound(all, times)
		count++
	}

	fmt.Printf("一共测试了 %d 轮，全部通过\n", count)
}

func newRound(all, occupyTimesPerProcess int) {
	rsc := new(resource)
	rsc.wg.Add(all * occupyTimesPerProcess)

	prop := observer.NewProperty(nil)

	ps := make([]Process, all)
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
			debugPrintf("%s 开始随机申请资源", p)
			for i < times {
				if p.CanRequest() {
					p.Request()
					i++
				}
				time.Sleep(time.Millisecond)
			}
		}(p, occupyTimesPerProcess)
	}

	rsc.wg.Wait()

	fmt.Println(rsc.report())
}
