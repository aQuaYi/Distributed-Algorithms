package mutual

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_timestamp_String(t *testing.T) {
	ast := assert.New(t)
	ts := newTimestamp(0, 1)
	actual := ts.String()
	expected := "<T0:P1>"
	ast.Equal(expected, actual)

}

func Test_less(t *testing.T) {
	ast := assert.New(t)

	// a < b < c
	a := timestamp{time: 1, process: 1}
	b := timestamp{time: 1, process: 2}
	c := timestamp{time: 2, process: 3}

	ast.True(less(a, b))
	ast.True(less(b, c))
	ast.True(less(a, c))

	ast.False(less(b, a))
	ast.False(less(c, a))
	ast.False(less(c, b))
}
