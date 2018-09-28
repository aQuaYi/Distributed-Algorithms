package main

import (
	"log"
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

	rsc := new(resource)
	rsc.wg.Add(all * occupyTimes)

	prop := observer.NewProperty(nil)

	ps := make([]Process, all)
	for i := range ps {
		p := newProcess(all, i, rsc, prop)
		p.AddOccupyTimes(occupyTimes)
		ps[i] = p
	}

	for i := range ps {
		go func(i int) {
			p := ps[i]
			for {
				if p.NeedResource() {
					p.Request()
				}
				randSleep()
			}
		}(i)
	}

	rsc.Wait()

	log.Println(rsc.report())
}
