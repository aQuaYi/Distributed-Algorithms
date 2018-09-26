package mutual

import (
	"container/heap"
	"fmt"
	"sync"
)

type requestQueue struct {
	rpq       *requestPriorityQueue
	requestOf map[timestamp]*request
	// FIXME: 试着删除 mutex 或者删除此条注释
	mutex sync.Mutex
}

func newRequestQueue() *requestQueue {
	return &requestQueue{
		rpq:       new(requestPriorityQueue),
		requestOf: make(map[timestamp]*request, 1024),
	}
}

func (rq *requestQueue) first() timestamp {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	if len(*rq.rpq) == 0 {
		return timestamp{process: others}
	}
	return (*rq.rpq)[0].timestamp2
}

func (rq *requestQueue) push(ts timestamp) {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	r := &request{
		timestamp2: ts,
	}

	rq.requestOf[ts] = r
	heap.Push(rq.rpq, r)
}

func (rq *requestQueue) remove(ts timestamp) {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	rq.rpq.remove(rq.requestOf[ts])
}

// func newRequestQueue() *requestQueue {
// 	return &requestQueue{
// 		rpq:,
// 		rqs:,
// 	}
// }

// request 是 priorityQueue 中的元素
type request struct {
	// TODO: 更名 timestamp2 到 timestamp
	timestamp2 timestamp
	// TODO: 删除 timestamp 和 process
	timestamp int
	process   int
	index     int
}

func (r *request) String() string {
	if r == nil {
		return "<:>"
	}
	return fmt.Sprintf("<%d:%d>", r.timestamp, r.process)
}

// rpq implements heap.Interface and holds entries.
type requestPriorityQueue []*request

func (q requestPriorityQueue) Len() int { return len(q) }

// NOTICE: 这就是将局部顺序推广到全局顺序的关键
func (q requestPriorityQueue) Less(i, j int) bool {
	return less(q[i].timestamp2, q[j].timestamp2)
}

func (q requestPriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

// Push 往 pq 中放 entry
func (q *requestPriorityQueue) Push(x interface{}) {
	temp := x.(*request)
	temp.index = len(*q)
	*q = append(*q, temp)
}

// Pop 从 pq 中取出最优先的 entry
func (q *requestPriorityQueue) Pop() interface{} {
	temp := (*q)[len(*q)-1]
	temp.index = -1 // for safety
	*q = (*q)[0 : len(*q)-1]
	return temp
}

func (q *requestPriorityQueue) remove(r *request) {
	heap.Remove(q, r.index)
}
