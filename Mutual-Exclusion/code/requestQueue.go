package mutual

import (
	"container/heap"
	"sync"
)

// RequestQueue 提供了操作 request queue 的接口
type RequestQueue interface {
	// Min 返回最小的 Timestamp 值
	Min() // TODO: 继续写
}

type requestQueue struct {
	rpq       *requestPriorityQueue
	requestOf map[Timestamp]*request
	mutex     sync.Mutex
}

// TODO: 修改返回值为 RequestQueue 接口
func newRequestQueue() *requestQueue {
	return &requestQueue{
		rpq:       new(requestPriorityQueue),
		requestOf: make(map[Timestamp]*request, 1024),
	}
}

func (rq *requestQueue) Min() Timestamp {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	if len(*rq.rpq) == 0 {
		return nil
	}
	return (*rq.rpq)[0].timestamp
}

func (rq *requestQueue) push(ts Timestamp) {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	r := &request{
		timestamp: ts,
	}

	rq.requestOf[ts] = r
	heap.Push(rq.rpq, r)
}

func (rq *requestQueue) remove(ts Timestamp) {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	rq.rpq.remove(rq.requestOf[ts])
	delete(rq.requestOf, ts)
}

// request 是 priorityQueue 中的元素
type request struct {
	timestamp Timestamp
	index     int
}

// rpq implements heap.Interface and holds entries.
type requestPriorityQueue []*request

func (q requestPriorityQueue) Len() int { return len(q) }

// NOTICE: 这就是将局部顺序推广到全局顺序的关键
func (q requestPriorityQueue) Less(i, j int) bool {
	return q[i].timestamp.Less(q[j].timestamp)
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
