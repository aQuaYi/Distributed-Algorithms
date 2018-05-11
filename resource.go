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

type resource struct {
	grantedTo   int
	occupyOrder []int
}

func newResource() *resource {
	return &resource{
		grantedTo: NULL,
	}
}

func (r *resource) occupy(p int) {
	if r.grantedTo != NULL {
		msg := fmt.Sprintf("资源正在被 P%d 占据，P%d 却想获取资源。", r.grantedTo, p)
		panic(msg)
	}
	r.grantedTo = p
	r.occupyOrder = append(r.occupyOrder, p)
	wg.Done()
	debugPrintf("~~~ @resource: P%d occupy ~~~", p)
}

func (r *resource) release(p int) {
	if r.grantedTo != p {
		msg := fmt.Sprintf("P%d 想要释放正在被 P%d 占据的资源。", p, r.grantedTo)
		panic(msg)
	}
	r.grantedTo = NULL
	debugPrintf("~~~ @resource: P%d release ~~~", p)
}

func (p *process) request() {
	r := &request{
		time:    p.clock.getTime(),
		process: p.me,
	}

	p.rwmu.Lock()

	debugPrintf("[%d]P%d request %s", p.clock.getTime(), p.me, r)

	p.append(r)

	p.clock.tick()

	p.messaging(requestResource, r)

	p.rwmu.Unlock()
}

func (p *process) occupy() {

	p.rsc.occupy(p.me)

	p.isOccupying = true

	p.clock.tick()

	// 经过一段时间，就释放资源
	go func(p *process) {
		occupyTime := time.Duration(5+rand.Intn(20)) * time.Millisecond
		time.Sleep(occupyTime)

		p.rwmu.Lock()

		p.release()

		p.clock.tick()

		p.rwmu.Unlock()
	}(p)
}

func (p *process) release() {
	r := p.requestQueue[0]

	p.rsc.release(p.me)
	p.isOccupying = false

	p.delete(p.requestQueue[0])

	// TODO: 这算不算一个 event 呢
	p.clock.tick()

	p.messaging(releaseResource, r)
}
