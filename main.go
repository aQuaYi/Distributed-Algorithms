package mutual

import (
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var wg sync.WaitGroup

func start(size, occupyNumber int, r *resource) []int {
	wg.Add(size)

	sys := newSystem(size, r)

	requestOrder := requestLoop(sys.processes, occupyNumber)

	sys.kill()

	wg.Wait()

	return requestOrder
}

func requestLoop(ps []*process, occupyNumber int) (requestOrder []int) {
	requestOrder = make([]int, occupyNumber)
	idx := 0

	for idx < occupyNumber {

		i := rand.Intn(len(ps))

		requestOrder[idx] = i

		ps[i].request()

		// 等待一段时间，再进行下一个 request
		waitingTime := time.Duration(10+rand.Intn(90)) * time.Millisecond
		time.Sleep(waitingTime)
		idx++
	}

	debugPrintf("完成全部 request 工作，len(requestOrder)=%d, requestOrder=%v", occupyNumber, requestOrder)

	return
}
