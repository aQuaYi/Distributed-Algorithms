package mutual

import (
	"fmt"
)

func (p *process) append(r *request) {
	p.requestQueue = append(p.requestQueue, r)
	debugPrintf("[%d]P%d append %s, MRT=%d, RT%v, RQ%v", p.clock.getTime(), p.me, r, p.minReceiveTime, p.receiveTime, p.requestQueue)
}

func (p *process) delete(r *request) {
	i := 0
	for p.requestQueue[i] != r {
		i++
	}
	last := len(p.requestQueue) - 1

	// 删除的时候，需要保持 requestQueue 的顺序
	copy(p.requestQueue[i:], p.requestQueue[i+1:])

	p.requestQueue = p.requestQueue[:last]

	debugPrintf("[%d]P%d delete %s, MRT=%d, RT%v, RQ%v", p.clock.getTime(), p.me, r, p.minReceiveTime, p.receiveTime, p.requestQueue)

	// p.requestQueue 变化时，都需要检查是否符合了 rule5
	p.toCheckRule5Chan <- struct{}{}
}

// requet 是 priorityQueue 中的元素
type request struct {
	timestamp int // request 的时间
	process   int // request 的 process
}

func (r *request) String() string {
	if r == nil {
		return "<>"
	}
	return fmt.Sprintf("<%d:%d>", r.timestamp, r.process)
}

// rpq implements heap.Interface and holds entries.
type rpq []*request

func (pq rpq) Len() int { return len(pq) }

func (pq rpq) Less(i, j int) bool {
	if pq[i].timestamp == pq[j].timestamp {
		return pq[i].process < pq[j].process
	}
	return pq[i].timestamp < pq[j].timestamp
}

func (pq rpq) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push 往 pq 中放 entry
func (pq *rpq) Push(x interface{}) {
	temp := x.(*request)
	*pq = append(*pq, temp)
}

// Pop 从 pq 中取出最优先的 entry
func (pq *rpq) Pop() interface{} {
	temp := (*pq)[len(*pq)-1]
	*pq = (*pq)[0 : len(*pq)-1]
	return temp
}
