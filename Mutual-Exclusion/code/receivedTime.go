package mutual

import (
	"container/heap"
	"sync"
)

type receivedTime struct {
	trq   *timeRecordQueue
	trs   []*timeRecord
	mutex sync.Mutex
}

func newReceivedTime(all, me int) *receivedTime {
	trq := new(timeRecordQueue)
	trs := make([]*timeRecord, all)
	for i := range trs {
		if i == me {
			continue
		}
		trs[i] = &timeRecord{}
		heap.Push(trq, trs[i])
	}
	return &receivedTime{
		trq: trq,
		trs: trs,
	}
}

// update 以后，返回 rt 中的最小值
func (rt *receivedTime) update(id, time int) int {
	rt.mutex.Lock()
	defer rt.mutex.Unlock()
	rt.trq.update(rt.trs[id], time)
	return (*rt.trq)[0].time
}

// timeRecord 是 priorityQueue 中的元素
type timeRecord struct {
	time  int
	index int
}

type timeRecordQueue []*timeRecord

func (trq timeRecordQueue) Len() int { return len(trq) }

func (trq timeRecordQueue) Less(i, j int) bool {
	return trq[i].time < trq[j].time
}

func (trq timeRecordQueue) Swap(i, j int) {
	trq[i], trq[j] = trq[j], trq[i]
	trq[i].index = i
	trq[j].index = j
}

// Push 往 pq 中放 entry
func (trq *timeRecordQueue) Push(x interface{}) {
	temp := x.(*timeRecord)
	temp.index = len(*trq)
	*trq = append(*trq, temp)
}

// Pop 从 pq 中取出最优先的 entry
func (trq *timeRecordQueue) Pop() interface{} {
	temp := (*trq)[len(*trq)-1]
	temp.index = -1 // for safety
	*trq = (*trq)[0 : len(*trq)-1]
	return temp
}

func (trq *timeRecordQueue) update(entry *timeRecord, priority int) {
	entry.time = priority
	heap.Fix(trq, entry.index)
}
