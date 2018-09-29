package main

import (
	"fmt"
	"sync"

	"github.com/aQuaYi/observer"
)

// OTHERS 表示信息接收方为其他所有 process
const OTHERS = -1

// Process 是进程的接口
type Process interface {
	// 检查 process 是否需要申请资源
	CanRequest() bool
	// 申请资源
	Request()
	// 输出信息
	String() string
}

type process struct {
	me int

	resource     Resource
	clock        Clock
	receivedTime ReceivedTime
	requestQueue RequestQueue

	prop   observer.Property
	stream observer.Stream

	rwm sync.RWMutex
	// 访问修改以下属性时，需要加锁
	requestTimestamp Timestamp
	isOccupying      bool
}

func newProcess(all, me int, r Resource, prop observer.Property) Process {
	p := &process{
		me:           me,
		resource:     r,
		prop:         prop,
		clock:        newClock(),
		requestQueue: newRequestQueue(),
		receivedTime: newReceivedTime(all, me),
	}

	p.Listening()

	debugPrintf("[%d]P%d 完成创建工作", p.clock.Now(), p.me)

	return p
}

func (p *process) String() string {
	return fmt.Sprintf("[%d]P%d", p.clock.Now(), p.me)
}

func (p *process) Listening() {
	stream := p.prop.Observe()

	debugPrintf("[%d]P%d 获取了 stream 开始监听", p.clock.Now(), p.me)

	go func() {
		for {
			msg := stream.Next().(*message)
			switch msg.msgType {
			case requestResource:
				p.handleRequestMessage(msg)
			case releaseResource:
				p.handleReleaseMessage(msg)
			case acknowledgment:
				p.handleAcknowledgeMessage(msg)
			}
		}
	}()
}

func (p *process) handleRequestMessage(msg *message) {
	if msg.from == p.me {
		return
	}
	// 收到消息，总是先更新自己的时间
	p.updateTime(msg.from, msg.msgTime)
	// rule 2: 把 msg.timestamp 放入自己的 requestQueue 当中
	p.requestQueue.Push(msg.timestamp)

	debugPrintf("%s 添加了 %s 后的 request queue 是 %s", p, msg.timestamp, p.requestQueue)

	// rule 2: 给对方发送一条 acknowledge 消息
	p.prop.Update(newMessage(
		acknowledgment,
		p.clock.Tick(),
		p.me,
		msg.from,
		msg.timestamp,
	))
	p.checkRule5()
}

func (p *process) handleReleaseMessage(msg *message) {
	if msg.from == p.me {
		return
	}
	// 收到消息，总是先更新自己的时间
	p.updateTime(msg.from, msg.msgTime)
	// rule 4: 收到就从 request queue 中删除相应的申请
	p.requestQueue.Remove(msg.timestamp)

	debugPrintf("%s 删除了 %s 后的 request queue 是 %s", p, msg.timestamp, p.requestQueue)

	p.checkRule5()
}

func (p *process) handleAcknowledgeMessage(msg *message) {
	if msg.to != p.me {
		return
	}
	// 收到消息，总是先更新自己的时间
	p.updateTime(msg.from, msg.msgTime)
	p.checkRule5()
}

func (p *process) Request() {
	ts := newTimestamp(p.clock.Tick(), p.me)
	p.rwm.Lock()
	p.requestTimestamp = ts
	p.rwm.Unlock()

	msg := newMessage(requestResource, p.clock.Tick(), p.me, OTHERS, ts)
	// Rule 1: 发送申请信息给其他的 process
	p.prop.Update(msg)

	p.requestQueue.Push(ts)
}

func (p *process) CanRequest() bool {
	p.rwm.RLock()
	defer p.rwm.RUnlock()
	return p.requestTimestamp == nil
}

func (p *process) updateTime(id, time int) {
	p.clock.Update(time)
	p.receivedTime.Update(id, time)
}

func (p *process) checkRule5() {
	p.rwm.Lock()
	defer p.rwm.Unlock()
	if !p.isOccupying &&
		p.requestTimestamp != nil &&
		p.requestTimestamp.IsEqual(p.requestQueue.Min()) &&
		p.requestTimestamp.Time() < p.receivedTime.Min() {
		p.occupyResource()
		go func() {
			randSleep()
			p.releaseResource()
		}()
	}
}

func (p *process) occupyResource() {
	debugPrintf("%s 准备占用资源 %s", p, p.requestQueue)
	p.isOccupying = true
	p.resource.Occupy(p.requestTimestamp)
}

func (p *process) releaseResource() {
	p.rwm.RLock()
	ts := p.requestTimestamp
	p.rwm.RUnlock()

	// rule 3: 先释放资源
	p.resource.Release(ts)
	// rule 3: 在 requestQueue 中删除 ts
	p.requestQueue.Remove(ts) // FIXME: 到底是先释放好，还是先删除好呢

	p.rwm.Lock()
	p.isOccupying = false
	p.requestTimestamp = nil
	p.rwm.Unlock()

	// rule 3: 把释放的消息发送给其他 process
	msg := newMessage(releaseResource, p.clock.Tick(), p.me, OTHERS, ts)
	p.prop.Update(msg)

}
