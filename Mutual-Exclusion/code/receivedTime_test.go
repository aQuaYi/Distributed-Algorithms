package mutual

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_receivedTime_update(t *testing.T) {
	ast := assert.New(t)
	all, me := 10, 0
	rt := newReceivedTime(all, me)
	// 把所有的接受值调整到较大的值
	for i := 1; i < all; i++ {
		rt.update(i, all+1)
	}
	// 依次按照以最小值更新第 i 个时间值
	for i := all - 1; i > me; i-- {
		expected := i
		rt.update(i, i)
		actual := rt.min()
		ast.Equal(expected, actual)
	}
}

func Test_receivedTime_updateItselfWillPanic(t *testing.T) {
	ast := assert.New(t)
	all, me := 10, 0
	rt := newReceivedTime(all, me)
	ast.Panics(func() { rt.update(me, 1) })
}
