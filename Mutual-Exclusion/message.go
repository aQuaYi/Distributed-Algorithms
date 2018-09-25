package mutual

import "fmt"

type message struct {
	msgType   msgType
	timestamp int // 发送 message 时， process.clock 的时间
	from      int // message 发送方的 ID
	request   *request
}

func newMessage(mt msgType, timestamp int, from int, request *request) *message {
	return &message{
		msgType:   mt,
		timestamp: timestamp,
		from:      from,
		request:   request,
	}
}

func (m *message) String() string {
	return fmt.Sprintf("msg{%s,T%d,P%d,%s}", m.msgType, m.timestamp, m.from, m.request)
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
