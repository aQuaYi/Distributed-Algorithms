package mutual

import (
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

func start(processNumber, takenTimes int) {
	wg.Add(processNumber)

	wg.Wait()
}

func requestLoop(ps []*process, occupyNumber int) (requestOrder []int) {
	requestOrder = make([]int, occupyNumber)
	idx := 0

	for idx < occupyNumber {
		idx++

		//
		timeout := time.Duration(100+rand.Intn(900)) * time.Millisecond
		time.Sleep(timeout)

		i := rand.Intn(len(ps))

		requestOrder[idx] = i

		p := ps[i]
		p.request()
	}

	return
}
