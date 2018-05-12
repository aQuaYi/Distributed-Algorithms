package mutual

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func start(size, occupyNumber int, r *resource) []int {
	r.occupied.Add(occupyNumber)

	sys := newSystem(size, r)

	requestOrder := requestLoop(sys.processes, occupyNumber)

	r.occupied.Wait()

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
		waitingTime := time.Duration(50+rand.Intn(50)) * time.Millisecond
		time.Sleep(waitingTime)
		idx++
	}

	debugPrintf("完成全部 request 工作，len(requestOrder)=%d, requestOrder=%v", occupyNumber, requestOrder)

	return
}
