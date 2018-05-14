package mutual

type system struct {
	processes []*process
}

func newSystem(size int) (*system, *resource) {
	chans := make([]chan *message, size)
	for i := range chans {
		chans[i] = make(chan *message, 100)
	}

	r := newResource()

	ps := make([]*process, size)
	for i := range ps {
		ps[i] = newProcess(i, r, chans)
	}

	return &system{
		processes: ps,
	}, r
}
