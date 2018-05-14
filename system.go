package mutual

type system struct {
	processes []*process
}

// size: process 的数量
func newSystem(size int, r *resource, channel chan int) *system {
	chans := make([]chan *message, size)
	for i := range chans {
		// TODO: chan 可以带缓冲吗？
		chans[i] = make(chan *message, 1000)
	}

	ps := make([]*process, size)
	for i := range ps {
		ps[i] = newProcess(i, r, channel, chans)
	}

	return &system{
		processes: ps,
	}
}
