package mutual

import "fmt"

type message struct {
	msgType msgType
	// TODO: 删除此处内容
	timestamp  int // 发送 message 时， process.clock 的时间
	from       int // message 发送方的 ID
	to         int // message 接收方的 ID， 当值为 others 的时候，表示接收方为除 from 外的所有
	request    *request
	timestamp2 timestamp
	msgTime    int
}

// TODO: 删除此处内容
func newMessage(mt msgType, timestamp int, from int, request *request) *message {
	return &message{
		msgType:   mt,
		timestamp: timestamp,
		from:      from,
		request:   request,
	}
}

func newMessage2(mt msgType, msgTime, from, to int, ts timestamp) *message {
	return &message{
		msgType:    mt,
		msgTime:    msgTime,
		from:       from,
		to:         to,
		timestamp2: ts,
	}
}

func (m *message) String() string {
	return fmt.Sprintf("{%s, Time:%d, From:%d, To:%d, %s}", m.msgType, m.msgTime, m.from, m.to, m.timestamp2)
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
		return "申请资源"
	case releaseResource:
		return "释放资源"
	default:
		return "确认收到"
	}
}
