package mutualexclusion

import (
	"fmt"
	"log"
	"testing"

	"github.com/aQuaYi/observer"
	"github.com/stretchr/testify/assert"
)

func run(all, occupyTimesPerProcess int) {
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

	log.Println(rsc.report())
}

func Test_process(t *testing.T) {
	ast := assert.New(t)
	//
	amount := 131072
	for all := 2; all <= 64; all *= 2 {
		times := amount / all
		name := fmt.Sprintf("%d Process × %d 次 = 共计 %d 次", all, times, amount)
		t.Run(name, func(t *testing.T) {
			ast.NotPanics(func() {
				run(all, times)
			})
		})
	}
}

func Test_process_String(t *testing.T) {
	ast := assert.New(t)
	//
	me := 1
	clock := newClock()
	p := &process{
		me:    me,
		clock: clock,
	}
	time := 999
	p.clock.Update(time)
	expected := fmt.Sprintf("[%d]P%d", time+1, me)
	actual := p.String()
	ast.Equal(expected, actual)
}
