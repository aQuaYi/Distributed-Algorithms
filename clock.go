package mutual

import (
	"sync"
)

type clock struct {
	time int
	// clock 使用单独的锁
	// 避免与 process 的锁冲突
	rwmu sync.RWMutex
}

// 每个 process 的 clock 的 initial time，都是随机的
func newClock() *clock {
	return &clock{
		// time: 1 + rand.Intn(100),
		time: 0,
	}
}

func (c *clock) getTime() int {
	c.rwmu.RLock()
	t := c.time
	c.rwmu.RUnlock()
	return t
}

func (c *clock) update(t int) {
	c.rwmu.Lock()
	c.time = max(c.time, t) + 1
	c.rwmu.Unlock()
}

func (c *clock) tick() int {
	c.rwmu.Lock()
	c.time++
	ts := c.time
	c.rwmu.Unlock()
	return ts
}
