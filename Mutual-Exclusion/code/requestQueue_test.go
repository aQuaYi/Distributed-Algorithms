package mutual

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeIncreasingTimestamps(half int) []Timestamp {
	res := make([]Timestamp, 0, half*2)
	for i := 0; i < half; i++ {
		res = append(res,
			newTimestamp(i, i*2),
			newTimestamp(i, i*2+1),
		)
	}
	return res
}

func Test_requestQueue(t *testing.T) {
	ast := assert.New(t)
	//
	half := 10
	size := half * 2
	tss := makeIncreasingTimestamps(half)
	rq := newRequestQueue()
	//
	for i := size - 1; i >= 0; i-- {
		ts := tss[i]
		rq.push(ts) // 每次放入到都是新的最小值
		expected := ts
		actual := rq.Min()
		ast.Equal(expected, actual)
	}
	//
	for i := 0; i+1 < size; i++ {
		rq.remove(tss[i])
		expected := tss[i+1] // 删除了最小值后，下个就是新的最小值
		actual := rq.Min()
		ast.Equal(expected, actual)
	}
}

func Test_requestQueue_remove(t *testing.T) {
	ast := assert.New(t)
	//
	half := 10
	size := half * 2
	tss := makeIncreasingTimestamps(half)
	rq := newRequestQueue()
	//
	for i := 0; i < size; i++ {
		ts := tss[i]
		rq.push(ts)
	}
	//
	expected := tss[0]
	for i, j := 1, size-1; i < j; i, j = i+1, j-1 {
		rq.remove(tss[i])
		actual := rq.Min()
		ast.Equal(expected, actual)
		//
		rq.remove(tss[j])
		actual = rq.Min()
		ast.Equal(expected, actual)
	}
}

func Test_requestQueue_emptyToFirst(t *testing.T) {
	ast := assert.New(t)
	rq := newRequestQueue()
	expected := others
	actual := rq.Min().Process()
	ast.Equal(expected, actual)
}
