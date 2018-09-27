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
	ast.True(newTime+1 >= c.time)
	//
	c.update(newTime)
	//
	expected := newTime + 1
	actual := c.getTime()
	ast.Equal(expected, actual)
}

func Test_clock_tick(t *testing.T) {
	ast := assert.New(t)
	//
	c := newClock()
	expected := c.getTime() + 1
	actual := c.tick()
	ast.Equal(expected, actual)
}
