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
	me    int
	mutex sync.Mutex
	wg    sync.WaitGroup

	// 为了保证 Process 内部的事件顺序
	// 访问修改以下属性时，需要加锁
	resource         Resource
	clock            Clock
	receivedTime     ReceivedTime
	requestQueue     RequestQueue
	prop             observer.Property
	stream           observer.Stream
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

	debugPrintf("%s 完成创建工作", p)

	return p
}

func (p *process) String() string {
	// 这个方法会在别的加锁方法内部出现
	// 为了避免死锁，此方法不加锁
	return fmt.Sprintf("[%d]P%d", p.clock.Now(), p.me)
}

func (p *process) Listening() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	stream := p.prop.Observe()

	debugPrintf("%s 获取了 stream 开始监听", p)

	go func() {
		for {
			msg := stream.Next().(*message)
			if msg.from == p.me {
				// 忽略自己发出的消息
				continue
			}
			switch msg.msgType {
			case requestResource:
				p.handleRequestMessage(msg)
			case releaseResource:
				p.handleReleaseMessage(msg)
			case acknowledgment:
				if msg.to != p.me {
					continue
				}
				p.handleAcknowledgeMessage(msg)
			}
			p.checkRule5()
		}
	}()
}

func (p *process) handleRequestMessage(msg *message) {
	p.mutex.Lock()

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

	p.mutex.Unlock()
}

func (p *process) handleReleaseMessage(msg *message) {
	p.mutex.Lock()

	// 收到消息，总是先更新自己的时间
	p.updateTime(msg.from, msg.msgTime)
	// rule 4: 收到就从 request queue 中删除相应的申请
	p.requestQueue.Remove(msg.timestamp)

	debugPrintf("%s 删除了 %s 后的 request queue 是 %s", p, msg.timestamp, p.requestQueue)

	p.mutex.Unlock()

	// checkRule5 有自己的锁
	p.checkRule5()
}

func (p *process) handleAcknowledgeMessage(msg *message) {
	p.mutex.Lock()

	// 收到消息，总是先更新自己的时间
	p.updateTime(msg.from, msg.msgTime)

	p.mutex.Unlock()

	// checkRule5 有自己的锁
	p.checkRule5()
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

func (p *process) updateTime(id, time int) {
	// 被带锁的方法引用，所以，不再加锁
	p.clock.Update(time)
	p.receivedTime.Update(id, p.clock.Now())
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
