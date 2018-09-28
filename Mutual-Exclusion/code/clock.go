package mutual

import (
	"math/rand"
	"sync"
)

type clock struct {
	time int
	rwmu sync.RWMutex
}

// 每个 process 的 clock 的 initial time，都是随机的
func newClock() *clock {
	return &clock{
		time: 1 + rand.Intn(100),
	}
}

func (c *clock) getTime() int {
	c.rwmu.RLock()
	t := c.time
	c.rwmu.RUnlock()
	return t
}

func (c *clock) update(t int) int {
	c.rwmu.Lock()
	defer c.rwmu.Unlock()
	c.time = max(c.time, t+1)
	return c.time
}

func (c *clock) tick() int {
	c.rwmu.Lock()
	defer c.rwmu.Unlock()
	c.time++
	return c.time
}
