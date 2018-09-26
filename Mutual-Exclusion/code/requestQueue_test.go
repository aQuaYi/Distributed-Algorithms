package mutual

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeIncreasingTimestamps(half int) []timestamp {
	res := make([]timestamp, 0, half*2)
	for i := 0; i < half; i++ {
		res = append(res,
			timestamp{time: i, process: i * 2},
			timestamp{time: i, process: i*2 + 1},
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

	for i := size - 1; i >= 0; i-- {
		ts := tss[i]
		rq.push(ts) // 每次放入到都是新的最小值
		expected := ts
		actual := rq.first()
		ast.Equal(expected, actual)
	}

	for i := 0; i+1 < size; i++ {
		rq.remove(tss[i])
		expected := tss[i+1] // 删除了最小值后，下个就是新的最小值
		actual := rq.first()
		ast.Equal(expected, actual)
	}
}
