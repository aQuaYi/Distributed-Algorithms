package mutual

import (
	"container/heap"
	"fmt"
	"sync"
)

type process struct {
	rwmu         sync.RWMutex
	me           int
	clock        *clock
	chans        []chan *message
	requestQueue rpq

	resource *resource

	sentTime       []int // 最近一次给别的 process 发送的消息，所携带的最后时间
	receiveTime    []int // 最近一次从别的 process 收到的消息，所携带的最后时间
	minReceiveTime int   // lastReceiveTime 中的最小值

	toCheckRule5Chan chan struct{} // 每次收到 message 后，都靠这个 chan 来通知检查此 process 是否已经满足 rule 5，以便决定是否占有 resource

	requestChan chan struct{}
	sendChan    chan *sendMsg
	occupyChan  chan struct{}

	// TODO: 删除此处内容
	isOccupying bool
}

func newProcess(me int, r *resource, chans []chan *message) *process {
	p := &process{
		me:             me,
		resource:       r,
		clock:          newClock(),
		chans:          chans,
		requestQueue:   make(rpq, 0, 1024),
		sentTime:       make([]int, len(chans)),
		receiveTime:    make([]int, len(chans)),
		minReceiveTime: 0,

		toCheckRule5Chan: make(chan struct{}),
		requestChan:      make(chan struct{}),
		sendChan:         make(chan *sendMsg),
		occupyChan:       make(chan struct{}),
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
}

func (p *process) handleCheckRule5() {
	debugPrintf("[%d]P%d 将要检查是否满足了 Rule5", p.clock.getTime(), p.me)

	if len(p.requestQueue) > 0 && // p.requestQueue 中还有元素
		p.requestQueue[0].process == p.me && // 排在首位的 repuest 是 p 自己的
		p.requestQueue[0].timestamp < p.minReceiveTime && // p 在 request 后，收到过所有其他 p 的回复
		!p.isOccupying { // 不能是正占用的资源

		debugPrintf("[%d]P%d 满足 Rule5 MRT=%d RT%v PQ%v", p.clock.getTime(), p.me, p.minReceiveTime, p.receiveTime, p.requestQueue)
		p.handleOccupy()

	}

	debugPrintf("[%d]P%d 不满足 Rule5 MRT=%d RT%v PQ%v", p.clock.getTime(), p.me, p.minReceiveTime, p.receiveTime, p.requestQueue)
}

func eventLoop(p *process) {
	debugPrintf("[%d]P%d 启动 eventLoop", p.clock.getTime(), p.me)

	go func() {

		for {
			p.clock.tick()
			select {
			case msg := <-p.chans[p.me]:
				p.handleMsg(msg)
			case <-p.requestChan:
				p.handleRequest()
			case sm := <-p.sendChan:
				p.handleSend(sm)
			// case <-p.occupyChan:
			// p.handleOccupy()
			case <-p.toCheckRule5Chan:
				p.handleCheckRule5()
			}
		}

	}()
}

func (p *process) handleMsg(msg *message) {
	debugPrintf("[%d]P%d receive %s", p.clock.getTime(), p.me, msg)

	// 接收到了一个新的消息
	// 根据 IR2
	// process 的 clock 需要根据 msg.time 进行更新
	// 无论 msg 是什么类型的消息
	p.clock.update(msg.timestamp)
	p.receiveTime[msg.senderID] = msg.timestamp
	p.updateMinReceiveTime()
	debugPrintf("[%d]P%d MRT=%d, RT%v, RQ%v ", p.clock.getTime(), p.me, p.minReceiveTime, p.receiveTime, p.requestQueue)

	r := msg.request
	switch msg.msgType {
	case requestResource:
		p.push(r)
		// rule 2
		// 收到 request message 后
		// 需要发送一个 acknowledgement message
		if p.sentTime[msg.senderID] <= r.timestamp {
			go func() {
				p.sendChan <- &sendMsg{
					receiveID: msg.senderID,
					msg: &message{
						msgType: acknowledgment,
						// timestamp 真正发送的时候更新
						senderID: p.me,
						// request: nil,
					},
				}
			}()
		}

	case releaseResource:
		p.pop(r)
	}

	go func() {
		p.toCheckRule5Chan <- struct{}{}
	}()
}

// TODO: finish this
type sendMsg struct {
	receiveID int
	msg       *message
}

func (p *process) handleSend(sm *sendMsg) {
	debugPrintf("[%d]P%d -> P%d，消息内容 %s", p.clock.getTime(), p.me, sm.receiveID, sm.msg)

	sm.msg.timestamp = p.clock.getTime()
	p.sentTime[sm.receiveID] = max(p.sentTime[sm.receiveID], p.clock.getTime())

	go func() {
		p.chans[sm.receiveID] <- sm.msg
	}()

}

func (p *process) push(r *request) {
	heap.Push(&p.requestQueue, r)
}

func (p *process) pop(r *request) {
	pr := heap.Pop(&p.requestQueue)
	if pr != r {
		msg := fmt.Sprintf("P%d 删除的 %s 不是需要删除的 %s", p.me, pr, r)
		panic(msg)
	}
}

func (p *process) request() {
	debugPrintf("[%d]P%d 准备 request", p.clock.getTime(), p.me)
	p.requestChan <- struct{}{}
}
