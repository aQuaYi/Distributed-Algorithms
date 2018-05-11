package mutual

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	// NULL 表示没有赋予任何人
	NULL = -1
)

var (
	rsc         *resource // 全局变量，随时随地都可以访问
	occupyOrder []int     // rsc 被占用的顺序
)

func init() {
	rsc = &resource{
		grantedTo: NULL,
	}
}

type resource struct {
	grantedTo int
}

func (r *resource) occupy(p int) {
	if r.grantedTo != NULL {
		msg := fmt.Sprintf("资源正在被 P%d 占据，P%d 却想获取资源。", r.grantedTo, p)
		panic(msg)
	}
	r.grantedTo = p
	occupyOrder = append(occupyOrder, p)
}

func (r *resource) release(p int) {
	if r.grantedTo != p {
		msg := fmt.Sprintf("P%d 想要释放正在被 P%d 占据的资源。", p, r.grantedTo)
		panic(msg)
	}
	r.grantedTo = NULL
}

func (p *process) request() {
	r := &request{
		time:    p.clock.getTime(),
		process: p.me,
	}

	p.rwmu.Lock()

	p.append(r)

	p.clock.tick()

	p.messaging(requestResource, r)

	p.rwmu.Unlock()
}

func (p *process) occupy() {

	rsc.occupy(p.me)

	p.occupying = p.requestQueue[0]

	p.clock.tick()

	// 经过一段时间，就释放资源
	go func(p *process) {
		occupyPeriod := time.Duration(100+rand.Intn(900)) * time.Millisecond
		time.Sleep(occupyPeriod)

		p.rwmu.Lock()

		p.release()

		p.clock.tick()

		p.rwmu.Unlock()
	}(p)
}

func (p *process) release() {
	r := p.requestQueue[0]

	rsc.release(p.me)
	p.occupying = nil

	p.delete(p.requestQueue[0])

	// TODO: 这算不算一个 event 呢
	p.clock.tick()

	p.messaging(releaseResource, r)
}
