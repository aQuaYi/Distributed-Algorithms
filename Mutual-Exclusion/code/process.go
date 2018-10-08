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
	// Request 会申请占用资源
	// 如果上次 Request 后，还没有占用并释放资源，会发生阻塞
	// 非线程安全
	Request()
}

type process struct {
	me int
	wg sync.WaitGroup

	clock        Clock
	resource     Resource
	receivedTime ReceivedTime
	requestQueue RequestQueue

	mutex sync.Mutex
	// 为了保证发送消息的原子性，
	// 从生成 timestamp 开始到 prop.update 完成，这个过程需要上锁
	prop observer.Property
	// 操作以下属性，需要加锁
	isOccupying      bool
	requestTimestamp Timestamp
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
	// stream 的观察起点位置，由上层调用 newProcess 的方式决定
	// 在生成完所有的 process 后，再更新 prop，
	// 才能保证所有的 process 都能收到全部消息
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

			// TODO: 这里需要加锁吗？

			// 收到消息的第一件，更新自己的 clock
			p.clock.Update(msg.msgTime)
			// 然后为了 Rule5(ii) 记录收到消息的时间
			p.receivedTime.Update(msg.from, p.clock.Now())

			switch msg.msgType {
			// case acknowledgment: 收到此类消息只用更新时钟，前面已经做了
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
	if !p.isOccupying &&
		p.requestTimestamp != nil &&
		p.requestTimestamp.IsEqual(p.requestQueue.Min()) &&
		p.requestTimestamp.Time() < p.receivedTime.Min() {
		p.occupyResource()
		go func() {
			// TODO: 把 releaseResource 从 go func 中拿出来
			p.releaseResource()
		}()
	}
	p.mutex.Unlock()
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
	// rule 3: 把释放的消息发送给其他 process
	msg := newMessage(releaseResource, p.clock.Tick(), p.me, OTHERS, ts)
	p.prop.Update(msg)
	p.isOccupying = false
	p.requestTimestamp = nil

	p.mutex.Unlock()

	p.wg.Done()
}

func (p *process) Request() {
	p.wg.Wait()
	p.wg.Add(1)

	p.mutex.Lock()
	p.clock.Tick()

	ts := newTimestamp(p.clock.Now(), p.me)
	msg := newMessage(requestResource, p.clock.Now(), p.me, OTHERS, ts)
	// Rule 1.1: 发送申请信息给其他的 process
	p.prop.Update(msg)
	// Rule 1.2: 把申请消息放入自己的 request queue
	p.requestQueue.Push(ts)
	// 修改辅助属性，便于后续检查
	p.requestTimestamp = ts

	p.mutex.Unlock()
}
