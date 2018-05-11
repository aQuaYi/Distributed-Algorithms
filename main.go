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

func resourceLoop(ps []*process, occupyNumber int) (requestOrder []int) {
	requestOrder = make([]int, occupyNumber)

	for occupyNumber > 0 {
		occupyNumber--

		//
		timeout := time.Duration(100+rand.Intn(900)) * time.Millisecond
		time.Sleep(timeout)

		i := rand.Intn(len(ps))
		p := ps[i]
		p.request()
	}

	return
}
