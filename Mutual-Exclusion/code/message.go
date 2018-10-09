package mutualexclusion

import "fmt"

type message struct {
	msgType   msgType
	from      int // message 发送方的 ID
	to        int // message 接收方的 ID， 当值为 OTHERS 的时候，表示接收方为除 from 外的所有
	timestamp Timestamp
	msgTime   int
}

func newMessage(mt msgType, msgTime, from, to int, ts Timestamp) *message {
	return &message{
		msgType:   mt,
		msgTime:   msgTime,
		from:      from,
		to:        to,
		timestamp: ts,
	}
}

func (m *message) String() string {
	return fmt.Sprintf("{%s, Time:%d, From:%d, To:%2d, %s}", m.msgType, m.msgTime, m.from, m.to, m.timestamp)
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
	default:
		return "确认"
	}
}
