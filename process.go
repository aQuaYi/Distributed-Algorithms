package mutual

import (
	"sync"
)

type process struct {
	rwmu             sync.RWMutex
	me               int
	clock            *clock
	chans            []chan *message
	requestQueue     []*request
	sentTime         []int         // 最近一次给别的 process 发送的消息，所携带的最后时间
	receiveTime      []int         // 最近一次从别的 process 收到的消息，所携带的最后时间
	minReceiveTime   int           // lastReceiveTime 中的最小值
	toCheckRule5Chan chan struct{} // 每次收到 message 后，都靠这个 chan 来通知检查此 process 是否已经满足 rule 5，以便决定是否占有 resource

	occupying *request
}

func newProcess(me int, chans []chan *message) *process {
	p := &process{
		me:               me,
		clock:            newClock(),
		chans:            chans,
		requestQueue:     make([]*request, 0, 1024),
		sentTime:         make([]int, len(chans)),
		receiveTime:      make([]int, len(chans)),
		minReceiveTime:   0,
		toCheckRule5Chan: make(chan struct{}, 1),
	}

	go p.receiveLoop()

	go p.occupyLoop()

	return p
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
		p.receiveTime[msg.senderID] = p.clock.getTime()
		p.updateMinReceiveTime()

		switch msg.msgType {
		case requestResource:
			p.append(msg.request)

			// rule 2
			// 收到 request message 后
			// 需要发送一个 acknowledgement message
			if p.sentTime[msg.senderID] <= msg.time {
				t := p.clock.getTime()
				p.sentTime[msg.senderID] = t

				p.send(msg.senderID, &message{
					msgType: acknowledgment,
					time:    t,
				})

			}

		case releaseResource:
			p.delete(msg.request)
		}

		p.rwmu.Unlock()

		p.toCheckRule5Chan <- struct{}{}

	}
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
}

func (p *process) occupyLoop() {
	for {
		<-p.toCheckRule5Chan

		p.rwmu.Lock()

		if len(p.requestQueue) > 0 && // p.requestQueue 中还有元素
			p.requestQueue[0].process == p.me && // 排在首位的 repuest 是 p 自己的
			p.requestQueue[0].time < p.minReceiveTime && // p 在 request 后，收到过所有其他 p 的回复
			p.occupying != p.requestQueue[0] { // 不能是正占用的资源

			p.occupy()
		}

		p.rwmu.Unlock()
	}
}
