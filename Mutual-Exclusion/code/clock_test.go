package mutual

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_clock_update(t *testing.T) {
	ast := assert.New(t)
	//
	c := newClock()
	newTime := 1000
	ast.True(newTime+1 >= c.Now())
	//
	c.Update(newTime)
	//
	expected := newTime + 1
	actual := c.Now()
	ast.Equal(expected, actual)
}

func Test_clock_tick(t *testing.T) {
	ast := assert.New(t)
	//
	c := newClock()
	expected := c.Now() + 1
	actual := c.Tick()
	ast.Equal(expected, actual)
}
