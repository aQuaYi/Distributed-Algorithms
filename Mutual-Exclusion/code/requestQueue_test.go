package mutualExclusion

import (
	"strings"
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
		rq.Push(ts) // 每次放入到都是新的最小值
		expected := ts
		actual := rq.Min()
		ast.Equal(expected, actual)
	}
	//
	for i := 0; i+1 < size; i++ {
		rq.Remove(tss[i])
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
		rq.Push(ts)
	}
	//
	expected := tss[0]
	for i, j := 1, size-1; i < j; i, j = i+1, j-1 {
		rq.Remove(tss[i])
		actual := rq.Min()
		ast.Equal(expected, actual)
		//
		rq.Remove(tss[j])
		actual = rq.Min()
		ast.Equal(expected, actual)
	}
}

func Test_requestQueue_MinOfEmpty(t *testing.T) {
	ast := assert.New(t)
	rq := newRequestQueue()
	ast.Nil(rq.Min())
}

func Test_requestQueue_String(t *testing.T) {
	ast := assert.New(t)
	size := 100
	// 创建 timestamps
	timestamps := make([]Timestamp, 0, size)
	for i := 1; i < size; i++ {
		timestamps = append(timestamps, newTimestamp(i, i))
	}
	// 创建 requestQueue，并添加 timestamp
	rq := newRequestQueue()
	for i := range timestamps {
		rq.Push(timestamps[i])
	}
	// 获取 rq 的字符输出
	rqs := rq.String()
	// 验证 rqs 中的内容
	for i := range timestamps {
		tss := timestamps[i].String()
		ast.True(strings.Contains(rqs, tss))
	}
}
