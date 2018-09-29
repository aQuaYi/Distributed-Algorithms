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
	for all := 2; all <= 100; all++ {
		for times := 10; times <= 10000; times *= 2 {
			if all*times > 20480 {
				continue
			}
			fmt.Printf("~~~ %d Process，每个占用资源 %d 次 ~~~\n", all, times)
			oneRound(all, times)
		}
	}
	return
}

func oneRound(all, occupyTimesPerProcess int) {
	rsc := new(resource)
	rsc.wg.Add(all * occupyTimesPerProcess)

	prop := observer.NewProperty(nil)

	ps := make([]Process, all)
	for i := range ps {
		p := newProcess(all, i, rsc, prop)
		ps[i] = p
	}

	debugPrintf("~~~ 已经成功创建了 %d 个 Process ~~~", all)

	for _, p := range ps {
		go func(p Process, times int) {
			i := 0
			debugPrintf("%s 开始随机申请资源", p)
			for i < times {
				if p.CanRequest() {
					p.Request()
					i++
					// debugPrintf("%s 第 %d 次申请资源", p, i)
				}
			}
		}(p, occupyTimesPerProcess)
	}

	rsc.wg.Wait()

	fmt.Println(rsc.report())
}
