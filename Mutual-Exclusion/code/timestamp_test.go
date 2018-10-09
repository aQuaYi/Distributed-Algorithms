package main

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

func Test_timestamp_Less(t *testing.T) {
	ast := assert.New(t)

	// a < b < c
	a := newTimestamp(1, 1)
	b := newTimestamp(1, 2)
	c := newTimestamp(2, 3)

	ast.True(a.Less(b))
	ast.True(a.Less(c))
	ast.True(b.Less(c))

	ast.False(b.Less(a))
	ast.False(c.Less(a))
	ast.False(c.Less(b))
}

func Test_timestamp_IsEqual_nil_false(t *testing.T) {
	ast := assert.New(t)
	ts := newTimestamp(0, 0)
	ast.False(ts.IsEqual(nil))
}

func Test_timestamp_IsEqual_same_true(t *testing.T) {
	ast := assert.New(t)
	time, process := 0, 0
	ts := newTimestamp(time, process)
	tsi := newTimestamp(time, process)
	ast.True(ts.IsEqual(tsi))
}

func Test_timestamp_IsBefore(t *testing.T) {
	ast := assert.New(t)
	time, process := 1, 0
	ts := newTimestamp(time, process)
	ast.False(ts.IsBefore(0))
	ast.True(ts.IsBefore(1))
}
