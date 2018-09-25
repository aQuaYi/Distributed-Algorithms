package mutual

import (
	"fmt"
)

// request 是 priorityQueue 中的元素
type request struct {
	timestamp int // request 的时间
	process   int // request 的 process
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
	if q[i].timestamp == q[j].timestamp {
		// timestamp 一样时
		// request 的顺序取决于 process 的顺序
		// process 的顺序可以任意指定
		// 我的选择是按照 process 的序号升序排列
		return q[i].process < q[j].process
	}
	// timestamp 不同时
	// 按照 timestamp 排序
	return q[i].timestamp < q[j].timestamp
}

func (q requestPriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

// Push 往 pq 中放 entry
func (q *requestPriorityQueue) Push(x interface{}) {
	temp := x.(*request)
	*q = append(*q, temp)
}

// Pop 从 pq 中取出最优先的 entry
func (q *requestPriorityQueue) Pop() interface{} {
	temp := (*q)[len(*q)-1]
	*q = (*q)[0 : len(*q)-1]
	return temp
}
