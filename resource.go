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

func (r *resource) occupy(req *request) {
	if r.grantedTo != NULL {
		msg := fmt.Sprintf("资源正在被 P%d 占据，P%d 却想获取资源。", r.grantedTo, req.process)
		panic(msg)
	}
	r.grantedTo = req.process
	r.occupyOrder = append(r.occupyOrder, req.process)
	r.occupied.Done()
	debugPrintf("~~~ @resource: %s occupy ~~~", req)
}

func (r *resource) release(req *request) {
	if r.grantedTo != req.process {
		msg := fmt.Sprintf("P%d 想要释放正在被 P%d 占据的资源。", req.process, r.grantedTo)
		panic(msg)
	}
	r.grantedTo = NULL
	debugPrintf("~~~ @resource: %s release ~~~", req)
}

func (p *process) handleRequest() {
	r := &request{
		timestamp: p.clock.getTime(),
		process:   p.me,
	}

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

	debugPrintf("[%d]P%d handleRequest，已分配好了所有发送消息的任务", p.clock.getTime(), p.me)

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

	// 根据 Rule3
	// 删除自身的 request
	debugPrintf("[%d]P%d handleRelease %s request queue %v", p.clock.getTime(), p.me, req, p.requestQueue)

	p.resource.release(req)
	p.isOccupying = false
	r := heap.Pop(&p.requestQueue).(*request)

	// 根据 Rule3
	// 给其他的 process 发消息

	for i := range p.chans {
		if i == p.me {
			continue
		}

		go func(i int, req *request) {
			sm := &sendMsg{
				receiveID: i,
				msg: &message{
					msgType: releaseResource,
					// timestamp 留在真正发送前更新
					senderID: p.me,
					request:  req,
				},
			}
			p.sendChan <- sm
		}(i, r)
	}
}
