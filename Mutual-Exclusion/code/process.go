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
	// WaitRequest 等待上一次资源占用完毕，然后立马申请资源
	WaitRequest()
}

type process struct {
	me int
	wg sync.WaitGroup

	clock            Clock
	resource         Resource
	receivedTime     ReceivedTime
	requestQueue     RequestQueue
	stream           observer.Stream
	isOccupying      bool
	requestTimestamp Timestamp

	mutex sync.Mutex
	// 操作以下属性，需要加锁
	prop observer.Property
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

	debugPrintf("%s 完成创建工作", p)

	return p
}

func (p *process) String() string {
	return fmt.Sprintf("[%d]P%d", p.clock.Now(), p.me)
}

func (p *process) Listening() {
	// 删除了这个地方的锁
	// stream 的观察起点位置，由上层调用 newProcess 的方式决定
	stream := p.prop.Observe()

	debugPrintf("%s 获取了 stream 开始监听", p)

	go func() {
		for {
			msg := stream.Next().(*message)
			if msg.from == p.me ||
				(msg.msgType == acknowledgment && msg.to != p.me) {
				// 忽略不该看见的消息
				continue
			}

			// 收到消息的第一件，更新自己的 clock
			p.clock.Update(msg.msgTime)
			// 然后为了 Rule5(ii) 记录收到消息的时间
			p.receivedTime.Update(msg.from, p.clock.Now())

			switch msg.msgType {
			// acknowledgment: 收到此类消息只用更新时钟，前面已经做了
			case requestResource:
				p.handleRequestMessage(msg)
			case releaseResource:
				p.handleReleaseMessage(msg)
			}
			p.checkRule5()
		}
	}()
}

func (p *process) handleRequestMessage(msg *message) {

	// rule 2.1: 把 msg.timestamp 放入自己的 requestQueue 当中
	p.requestQueue.Push(msg.timestamp)

	debugPrintf("%s 添加了 %s 后的 request queue 是 %s", p, msg.timestamp, p.requestQueue)

	p.mutex.Lock()

	// rule 2.2: 给对方发送一条 acknowledge 消息
	p.prop.Update(newMessage(
		acknowledgment,
		p.clock.Tick(),
		p.me,
		msg.from,
		msg.timestamp,
	))

	p.mutex.Unlock()
}

func (p *process) handleReleaseMessage(msg *message) {
	// rule 4: 从 request queue 中删除相应的申请
	p.requestQueue.Remove(msg.timestamp)
	debugPrintf("%s 删除了 %s 后的 request queue 是 %s", p, msg.timestamp, p.requestQueue)
}

func (p *process) checkRule5() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.isOccupying &&
		p.requestTimestamp != nil &&
		p.requestTimestamp.IsEqual(p.requestQueue.Min()) &&
		p.requestTimestamp.Time() < p.receivedTime.Min() {
		p.occupyResource()
		go func() {
			p.releaseResource()
		}()
	}
}

func (p *process) occupyResource() {
	// 利用 checkRule5 的锁进行锁定
	debugPrintf("%s 准备占用资源 %s", p, p.requestQueue)
	p.isOccupying = true
	p.resource.Occupy(p.requestTimestamp)
}

func (p *process) releaseResource() {
	p.mutex.Lock()

	ts := p.requestTimestamp
	// rule 3: 先释放资源
	p.resource.Release(ts)
	// rule 3: 在 requestQueue 中删除 ts
	p.requestQueue.Remove(ts) // FIXME: 到底是先释放好，还是先删除好呢?
	p.isOccupying = false
	p.requestTimestamp = nil
	// rule 3: 把释放的消息发送给其他 process
	msg := newMessage(releaseResource, p.clock.Tick(), p.me, OTHERS, ts)
	p.prop.Update(msg)

	p.mutex.Unlock()

	p.wg.Done()
}

func (p *process) WaitRequest() {
	p.wg.Wait()
	p.wg.Add(1)

	p.mutex.Lock()

	ts := newTimestamp(p.clock.Tick(), p.me)
	p.requestTimestamp = ts

	msg := newMessage(requestResource, p.clock.Tick(), p.me, OTHERS, ts)
	// Rule 1: 发送申请信息给其他的 process
	p.prop.Update(msg)

	p.requestQueue.Push(ts)

	p.mutex.Unlock()
}
