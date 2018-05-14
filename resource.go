package mutual

import (
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
	p.requestChan <- struct{}{}
}
