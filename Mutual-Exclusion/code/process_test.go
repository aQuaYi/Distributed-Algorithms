package mutualexclusion

import (
	"fmt"

	"github.com/aQuaYi/observer"
)

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

	rsc.wait()

	fmt.Println(rsc.report())
}
