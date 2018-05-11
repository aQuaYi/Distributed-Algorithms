package mutual

import "fmt"

type message struct {
	time     int // 发送 message 时， process.clock 的时间
	msgType  msgType
	senderID int // message 发送方的 ID
	request  *request
}

func (m *message) String() string {
	return fmt.Sprintf("msg{T%d,P%d,%s,request[%s]}", m.time, m.senderID, m.msgType, m.request)
}

type msgType int

// 枚举了 message 的所有类型
const (
	// REQUEST_RESOURCE 请求资源
	requestResource msgType = iota
	releaseResource
	acknowledgment
)

func (mt msgType) String() string {
	switch mt {
	case requestResource:
		return "申请"
	case releaseResource:
		return "释放"
	case acknowledgment:
		return "确认"
	default:
		panic("出现了多余的 msgType")
	}
}

type request struct {
	time    int // request 的时间
	process int // request 的 process
}

func (r *request) String() string {
	return fmt.Sprintf("[T%d:P%d]", r.time, r.process)
}

func (p *process) messaging(mt msgType, r *request) {
	for i := range p.chans {
		if i == p.me {
			continue
		}

		p.send(i, &message{
			msgType:  mt,
			time:     p.clock.getTime(),
			senderID: p.me,
			request:  r,
		})

	}
}

func (p *process) send(id int, msg *message) {
	p.chans[id] <- msg
	// send 是一个 event
	// 所以，发送完成后，需要 clock.tick()
	p.clock.tick()
}
