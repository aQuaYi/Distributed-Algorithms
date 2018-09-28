package mutual

import (
	"github.com/aQuaYi/observer"
)

type system struct {
	processes []*process
}

func newSystem(size int, r Resource) *system {
	chans := make([]chan *message, size)
	for i := range chans {
		chans[i] = make(chan *message, 100)
	}
	prop := observer.NewProperty(nil)
	ps := make([]*process, size)
	for i := range ps {
		ps[i] = newProcess(size, i, r, prop)
	}

	return &system{
		processes: ps,
	}
}
