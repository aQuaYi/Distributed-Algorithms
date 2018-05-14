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
	grantedTo    int
	processOrder []int
	timeOrder    []int
	occupied     sync.WaitGroup
}

func newResource() *resource {
	return &resource{
		grantedTo: NULL,
	}
}

func (r *resource) occupy(req *request) {
	if r.grantedTo != NULL {
		msg := fmt.Sprintf("资源正在被 P%d 占据，P%d 却想获取资源。", r.grantedTo, req.process)
		panic(msg)
	}
	r.grantedTo = req.process
	r.processOrder = append(r.processOrder, req.process)
	r.timeOrder = append(r.timeOrder, req.timestamp)
	debugPrintf("~~~ @resource: %s occupy ~~~ %v", req, r.processOrder[max(0, len(r.processOrder)-6):])
}

func (r *resource) release(req *request) {
	if r.grantedTo != req.process {
		msg := fmt.Sprintf("P%d 想要释放正在被 P%d 占据的资源。", req.process, r.grantedTo)
		panic(msg)
	}
	r.grantedTo = NULL
	debugPrintf("~~~ @resource: %s release ~~~ %v", req, r.processOrder[max(0, len(r.processOrder)-6):])
	r.occupied.Done()
}

func (p *process) request() {
	r := &request{
		timestamp: p.clock.getTime(),
		process:   p.me,
	}

	p.orderChan <- p.me

	debugPrintf("[%d]P%d handleRequest，生成 r=%s", p.clock.getTime(), p.me, r)

	// 根据 Rule1
	// 把 r 放入自身的 request queue
	p.push(r)

	debugPrintf("[%d]P%d handleRequest，已加入 request queue %v", p.clock.getTime(), p.me, p.requestQueue)

	// 根据 Rule1
	// 给其他的 process 发消息

	for i := range p.chans {
		if i == p.me {
			continue
		}

		p.chans[i] <- &message{
			msgType:   requestResource,
			timestamp: p.clock.tick(),
			senderID:  p.me,
			request:   r,
		}
	}

	debugPrintf("[%d]P%d 已经完成所有的 request 通知任务", p.clock.getTime(), p.me)

}

func (p *process) handleOccupy() {
	req := p.requestQueue[0]
	debugPrintf("[%d]P%d handleOccupy %s request queue %v", p.clock.getTime(), p.me, req, p.requestQueue)

	p.isOccupying = true

	p.resource.occupy(req)

	randSleep()

	p.handleRelease()
}

func (p *process) handleRelease() {
	req := p.requestQueue[0]

	p.resource.release(req)
	p.isOccupying = false

	// 根据 Rule3
	// 删除自身的 request
	heap.Pop(&p.requestQueue)

	debugPrintf("[%d]P%d handleRelease %s request queue %v", p.clock.getTime(), p.me, req, p.requestQueue)

	// 根据 Rule3
	// 给其他的 process 发消息

	for i := range p.chans {
		if i == p.me {
			continue
		}
		p.chans[i] <- &message{
			msgType:   releaseResource,
			timestamp: p.clock.tick(),
			senderID:  p.me,
			request:   req,
		}
	}

}
