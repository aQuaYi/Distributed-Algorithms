package mutual

import (
	"math/rand"
)

func init() {
	// TODO: 添加随机种子
	// rand.Seed(time.Now().UnixNano())
}

func start(size, occupyNumber int, r *resource) []int {
	r.occupied.Add(occupyNumber)

	recorder := newRecorder()

	sys := newSystem(size, r)

	requestLoop(sys.processes, occupyNumber)

	r.occupied.Wait()

	return *recorder
}

func requestLoop(ps []*process, occupyNumber int) {
	idx := 0

	for idx < occupyNumber {

		i := rand.Intn(len(ps))

		ps[i].request()

		// 等待一段时间，再进行下一个 request
		randSleep()

		idx++
	}

	debugPrintf("完成全部 request 工作", occupyNumber)

	return
}
