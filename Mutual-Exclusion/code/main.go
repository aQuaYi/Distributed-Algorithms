package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/aQuaYi/observer"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// TODO: 解决启动后，立刻死锁的问题
// TODO: 解决 request queue 会删除不存在的元素的问题

func main() {
	all := 10
	occupyTimesPerProcess := 10

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
			debugPrintf("%s 准备开始随机申请资源", p)
			for i < times {
				if p.CanRequest() {
					p.Request()
					i++
					debugPrintf("%s 第 %d 次申请资源", p, i)
				}
				randSleep()
			}
		}(p, occupyTimesPerProcess)
	}

	rsc.wg.Wait()

	log.Println(rsc.report())
}
