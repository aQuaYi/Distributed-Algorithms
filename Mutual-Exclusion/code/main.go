package mutual

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func start(processes, occupieds int) *resource {
	resource := newResource()
	return resource
}

func requestLoop(ps []*process, occupieds int) {
	count := 0

	for count < occupieds {
		count++
		i := rand.Intn(len(ps))
		ps[i].request()
		// 等待一段时间，再进行下一个 request
		randSleep()

	}

	debugPrintf("完成全部 request 工作")
}
