package mutual

type system struct {
	processes []*process
}

// size: process 的数量
func newSystem(size int) *system {
	chans := make([]chan message, size)
	for i := range chans {
		chans[i] = make(chan message)
	}

	ps := make([]*process, size)
	for i := range ps {
		ps[i] = newProcess(i, chans)
	}

	return &system{
		processes: ps,
	}
}
