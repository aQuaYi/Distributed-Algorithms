package mutualexclusion

import (
	"math/rand"
	"sync"
)

// Clock 是系统的逻辑时钟接口
type Clock interface {
	// Update 根据输入参数更新自身的值
	Update(int)
	// Tick 时钟跳动一次，并返回最新的时间值
	Tick() int
	// Now 返回当前的时间值
	Now() int
}

type clock struct {
	time int
	rwmu sync.RWMutex
}

// 每个 process 的 clock 的 initial time，都是随机的
func newClock() Clock {
	return &clock{
		time: 1 + rand.Intn(100),
	}
}

func (c *clock) Update(t int) {
	c.rwmu.Lock()
	c.time = max(c.time, t+1)
	c.rwmu.Unlock()
}

func (c *clock) Tick() int {
	c.rwmu.Lock()
	c.time++
	t := c.time
	c.rwmu.Unlock()
	return t
}

func (c *clock) Now() int {
	c.rwmu.RLock()
	t := c.time
	c.rwmu.RUnlock()
	return t
}
