package main

import (
	"math/rand"
	"time"

	"github.com/aQuaYi/observer"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	all := 10
	occupyTimes := 10
	rsc := newResource(all * occupyTimes)
	ps := make([]Process, all)
	prop := observer.NewProperty(nil)
	for i := range ps {
		p := newProcess(all, i, rsc, prop)
		p.AddOccupyTimes(occupyTimes)
		ps[i] = p
	}
}
