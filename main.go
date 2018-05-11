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

	sys := newSystem(size, r)

	requestOrder := requestLoop(sys.processes, occupyNumber)

	sys.kill()

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
		waitingTime := time.Duration(100+rand.Intn(900)) * time.Millisecond
		time.Sleep(waitingTime)
		idx++
	}

	return
}
