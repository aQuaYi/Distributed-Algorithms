package mutual

import (
	"container/heap"
	"fmt"
	"sync"
)

const (
	// NULL 表示没有赋予任何人
	NULL = -1
)

type resource struct {
	grantedTo   int
	occupyOrder []int
	occupied    sync.WaitGroup
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
	r.occupied.Done()
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

func (p *process) handleRequest() {
	r := &request{
		timestamp: p.clock.getTime(),
		process:   p.me,
	}

	// 根据 Rule1
	// 把 r 放入自身的 request queue
	p.push(r)

	// 根据 Rule1
	// 给其他的 process 发消息

	for i := range p.chans {
		if i == p.me {
			continue
		}

		go func(i int) {
			sm := &sendMsg{
				receiveID: i,
				msg: &message{
					msgType: requestResource,
					// timestamp 留在真正发送前更新
					senderID: p.me,
					request:  r,
				},
			}
			p.sendChan <- sm
		}(i)
	}
}

func (p *process) handleOccupy() {
	p.isOccupying = true
	p.resource.occupy(p.me)
	randSleep()
	p.resource.release(p.me)
	p.isOccupying = false
	p.handleRelease()
}

func (p *process) handleRelease() {
	// 根据 Rule3
	// 删除自身的 request
	r := heap.Pop(&p.requestQueue).(*request)

	// 根据 Rule3
	// 给其他的 process 发消息

	for i := range p.chans {
		if i == p.me {
			continue
		}

		go func(i int) {
			sm := &sendMsg{
				receiveID: i,
				msg: &message{
					msgType: releaseResource,
					// timestamp 留在真正发送前更新
					senderID: p.me,
					request:  r,
				},
			}
			p.sendChan <- sm
		}(i)
	}
}
