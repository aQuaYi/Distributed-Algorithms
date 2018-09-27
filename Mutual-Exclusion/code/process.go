package mutual

import (
	"container/heap"
	"fmt"

	"github.com/aQuaYi/observer"
)

const others = -1

type process struct {
	me           int
	clock        *clock
	chans        []chan *message
	requestQueue requestPriorityQueue

	resource *resource

	receiveTime    []int // 最近一次从别的 process 收到的消息，所携带的最后时间
	minReceiveTime int   // lastReceiveTime 中的最小值

	toCheckRule5Chan chan struct{} // 每次收到 message 后，都靠这个 chan 来通知检查此 process 是否已经满足 rule 5，以便决定是否占有 resource

	requestChan chan struct{}
	releaseChan chan struct{}
	occupyChan  chan struct{}

	isOccupying bool

	// 新的属性
	requestTimestamp timestamp
	prop             observer.Property
	stream           observer.Stream
	receivedTime     *receivedTime
	requestQueue2    *requestQueue
	occupyTimes      int // process 可以占用资源的次数
}

func newProcess(me int, r *resource, chans []chan *message) *process {
	p := &process{
		me:               me,
		resource:         r,
		clock:            newClock(),
		chans:            chans,
		requestQueue:     make(requestPriorityQueue, 0, 1024),
		receiveTime:      make([]int, len(chans)),
		minReceiveTime:   0,
		toCheckRule5Chan: make(chan struct{}),
		requestChan:      make(chan struct{}),
		occupyChan:       make(chan struct{}),
	}
	eventLoop(p)
	debugPrintf("[%d]P%d 完成创建工作", p.clock.getTime(), p.me)
	return p
}

func newProcess2(all, me int, r *resource, prop observer.Property) *process {
	p := &process{
		me:               me,
		resource:         r,
		clock:            newClock(),
		requestQueue2:    newRequestQueue(),
		receivedTime:     newReceivedTime(all, me),
		requestTimestamp: NOBODY2,
	}
	eventLoop(p)
	debugPrintf("[%d]P%d 完成创建工作", p.clock.getTime(), p.me)
	return p
}

func (p *process) updateMinReceiveTime() {
	i := (p.me + 1) % len(p.chans)
	minTime := p.receiveTime[i]
	for i, t := range p.receiveTime {
		if i == p.me {
			continue
		}
		minTime = min(minTime, t)
	}
	p.minReceiveTime = minTime

	debugPrintf("[%d]P%d updated p.minReceiveTime=%d, RT%v, RQ%v ", p.clock.getTime(), p.me, p.minReceiveTime, p.receiveTime, p.requestQueue)
}

func (p *process) push(r *request) {
	heap.Push(&p.requestQueue, r)
	debugPrintf("[%d]P%d push(%s) 后，request queue %v", p.clock.getTime(), p.me, r, p.requestQueue)
}

func (p *process) pop(r *request) {
	req := heap.Pop(&p.requestQueue).(*request)
	if req != r {
		msg := fmt.Sprintf("需要删除的是 %s，实际删除的是 %s，P%d.RQ%s", r, req, p.me, p.requestQueue)
		panic(msg)
	}

	debugPrintf("[%d]P%d pop(%s) 后，request queue %v", p.clock.getTime(), p.me, req, p.requestQueue)
}

func (p *process) request() {
	ts := newTimestamp(p.clock.tick(), p.me)
	msg := newMessage2(requestResource, p.clock.tick(), p.me, others, ts)
	// Rule 1: 发送申请信息给其他的 process
	p.prop.Update(msg)
	p.requestQueue2.push(ts)
}

func (p *process) occupy() {
	p.isOccupying = true
	p.resource.occupy2(p.requestTimestamp)
}

func (p *process) release() {
	ts := p.requestTimestamp
	// rule 3: 先释放资源
	p.resource.release2(ts)
	// rule 3: 在 requestQueue 中删除 ts
	p.requestQueue2.remove(ts)
	// rule 3: 把释放的消息发送给其他 process
	msg := newMessage2(releaseResource, p.clock.tick(), p.me, others, ts)
	p.prop.Update(msg)

	p.requestTimestamp = NOBODY2
	p.isOccupying = false
}
