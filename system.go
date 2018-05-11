package mutual

type system struct {
	processes []*process
}

// size: process 的数量
func newSystem(size int, r *resource) *system {
	chans := make([]chan *message, size)
	for i := range chans {
		// TODO: chan 可以带缓冲吗？
		chans[i] = make(chan *message, 100)
	}

	ps := make([]*process, size)
	for i := range ps {
		ps[i] = newProcess(i, r, chans)
	}

	return &system{
		processes: ps,
	}
}
