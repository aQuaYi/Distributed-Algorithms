package mutual

import (
	"sync"
)

type request struct {
	time    int // request 的时间
	process int // request 的 process
}

type process struct {
	rwmu             *sync.RWMutex
	me               int
	clock            *clock
	chans            []chan message
	requestQueue     []*request
	sentTime         []int         // 最近一次给别的 process 发送的消息，所携带的最后时间
	receiveTime      []int         // 最近一次从别的 process 收到的消息，所携带的最后时间
	minReceiveTime   int           // lastReceiveTime 中的最小值
	toCheckRule5Chan chan struct{} // 每次收到 message 后，都靠这个 chan 来通知检查此 process 是否已经满足 rule 5，以便决定是否占有 resource

	isHoldingResource bool // process 正在占用资源

	takenCounterDown int // 每个 preocess 想要占用 resource 的次数
}

func newProcess(me, takenTimes, initTakonProcess int, r *resource, chans []chan message) *process {
	rq := make([]*request, 1, len(chans)*takenTimes*2)
	p := &process{
		me:               me,
		clock:            newClock(),
		chans:            chans,
		requestQueue:     rq,
		sentTime:         make([]int, len(chans)),
		receiveTime:      make([]int, len(chans)),
		minReceiveTime:   0,
		toCheckRule5Chan: make(chan struct{}, 1),
		takenCounterDown: takenTimes,
	}

	if me == initTakonProcess {
		p.isHoldingResource = true
	}

	return p
}

func (p *process) occupy(rsc *resource) {

	return
}

func (p *process) request() {
	r := &request{
		time:    p.clock.getTime(),
		process: p.me,
	}

	p.append(r)

	// TODO: 这算不算一个 event 呢
	p.clock.tick()

	p.messaging(requestResource, r)
}

func (p *process) release() {
	i := 0
	for p.requestQueue[i].process != p.me {
		i++
	}
	r := p.requestQueue[i]

	rsc.release(p.me)

	p.isHoldingResource = false

	// TODO: 这算不算一个 event 呢
	p.clock.tick()

	p.delete()

	// TODO: 这算不算一个 event 呢
	p.clock.tick()

	p.messaging(releaseResource, r)
}

func (p *process) messaging(mt msgType, r *request) {
	for i, ch := range p.chans {
		if i == p.me {
			continue
		}
		ch <- message{
			msgType:  mt,
			time:     p.clock.getTime(),
			senderID: p.me,
			request:  r,
		}
		// sending 是一个 event
		// 所以，发送完成后，需要 clock.tick()
		p.clock.tick()
	}
}

func (p *process) isCounterDown() bool {
	return p.takenCounterDown == 0 && !p.isHoldingResource
}

func (p *process) receiveLoop() {
	msgChan := p.chans[p.me]
	for {
		msg := <-msgChan

		p.rwmu.Lock()

		// 接收到了一个新的消息
		// 根据 IR2
		// process 的 clock 需要根据 msg.time 进行更新
		// 无论 msg 是什么类型的消息
		p.clock.update(msg.time)
		p.receiveTime[msg.senderID] = msg.time
		p.updateMinReceiveTime()

		switch msg.msgType {
		case requestResource:
			p.append(msg.request)
		case releaseResource:
			p.delete()
		}

		p.toCheckRule5Chan <- struct{}{}

		p.rwmu.Unlock()
	}
}

func (p *process) append(r *request) {
	p.requestQueue = append(p.requestQueue, r)
}

func (p *process) delete() {
	last := len(p.requestQueue) - 1
	p.requestQueue[0], p.requestQueue[last] = p.requestQueue[last], p.requestQueue[0]
	p.requestQueue = p.requestQueue[:last]
}

func (p *process) updateMinReceiveTime() {
	idx := (p.me + 1) % len(p.chans)
	minTime := p.receiveTime[idx]
	for i, t := range p.receiveTime {
		if i == p.me {
			continue
		}
		minTime = min(minTime, t)
	}
	p.minReceiveTime = minTime
}

func (p *process) occupyLoop() {
	for {
		<-p.toCheckRule5Chan

		p.rwmu.Lock()

		if len(p.requestQueue) > 0 && // p.requestQueue 中还有元素
			p.requestQueue[0].process == p.me && // 排在首位的 repuest 是 p 自己的
			p.requestQueue[0].time < p.minReceiveTime { // p 在 request 后，收到过所有其他 p 的回复

			p.occupy(rsc)

			p.isHoldingResource = true

			// TODO: 这里需要 tick 一下吗
			p.clock.tick()
		}

		p.rwmu.Unlock()
	}
}
