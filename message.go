package mutual

type message struct {
	msgType  msgType
	time     int // 发送 message 时， process.clock 的时间
	senderID int // message 发送方的 ID
	request  *request
}

type msgType int

// 枚举了 message 的所有类型
const (
	// REQUEST_RESOURCE 请求资源
	requestResource msgType = iota
	releaseResource
	acknowledgment
)

type request struct {
	time    int // request 的时间
	process int // request 的 process
}
