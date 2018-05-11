package mutual

import (
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

func start(size, occupyNumber int) []int {
	occupyOrder = nil

	sys := newSystem(size)

	requestOrder := requestLoop(sys.processes, occupyNumber)

	sys.kill()

	return requestOrder
}

func requestLoop(ps []*process, occupyNumber int) (requestOrder []int) {
	requestOrder = make([]int, occupyNumber)
	idx := 0

	for idx < occupyNumber {
		idx++

		// 等待一段时间，再进行下一个 request
		waitingTime := time.Duration(100+rand.Intn(900)) * time.Millisecond
		time.Sleep(waitingTime)

		i := rand.Intn(len(ps))

		requestOrder[idx] = i

		ps[i].request()
	}

	return
}
