package mutualExclusion

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_receivedTime_update(t *testing.T) {
	ast := assert.New(t)
	all, me := 10, 0
	rt := newReceivedTime(all, me)
	// 把所有的接受值调整到较大的值
	for i := 1; i < all; i++ {
		rt.Update(i, all+1)
	}
	// 依次按照以最小值更新第 i 个时间值
	for i := all - 1; i > me; i-- {
		expected := i
		rt.Update(i, i)
		actual := rt.Min()
		ast.Equal(expected, actual)
	}
}

func Test_receivedTime_updateItselfWillPanic(t *testing.T) {
	ast := assert.New(t)
	all, me := 10, 0
	rt := newReceivedTime(all, me)
	ast.Panics(func() { rt.Update(me, 1) })
}

func Test_timeRecordQueue_Pop(t *testing.T) {
	ast := assert.New(t)
	trq := new(timeRecordQueue)
	expected := &timeRecord{
		time: 1,
	}
	heap.Push(trq, expected)
	actual := heap.Pop(trq).(*timeRecord)
	ast.Equal(expected.time, actual.time)
}
