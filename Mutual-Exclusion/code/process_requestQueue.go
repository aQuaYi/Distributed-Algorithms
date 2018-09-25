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

// // entry 是 priorityQueue 中的元素
// type entry struct {
// 	key      string
// 	priority int
// 	// index 是 entry 在 heap 中的索引号
// 	// entry 加入 Priority Queue 后， Priority 会变化时，很有用
// 	// 如果 entry.priority 一直不变的话，可以删除 index
// 	index int
// }

// // PQ implements heap.Interface and holds entries.
// type PQ []*entry

// func (pq PQ) Len() int { return len(pq) }

// func (pq PQ) Less(i, j int) bool {
// 	return pq[i].priority < pq[j].priority
// }

// func (pq PQ) Swap(i, j int) {
// 	pq[i], pq[j] = pq[j], pq[i]
// 	pq[i].index = i
// 	pq[j].index = j
// }

// // Push 往 pq 中放 entry
// func (pq *PQ) Push(x interface{}) {
// 	temp := x.(*entry)
// 	temp.index = len(*pq)
// 	*pq = append(*pq, temp)
// }

// // Pop 从 pq 中取出最优先的 entry
// func (pq *PQ) Pop() interface{} {
// 	temp := (*pq)[len(*pq)-1]
// 	temp.index = -1 // for safety
// 	*pq = (*pq)[0 : len(*pq)-1]
// 	return temp
// }

// // update modifies the priority and value of an entry in the queue.
// func (pq *PQ) update(entry *entry, value string, priority int) {
// 	entry.key = value
// 	entry.priority = priority
// 	heap.Fix(pq, entry.index)
// }
