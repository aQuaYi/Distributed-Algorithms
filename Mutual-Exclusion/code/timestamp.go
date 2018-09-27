package mutual

import "fmt"

type timestamp struct {
	time, process int
}

func newTimestamp(time, process int) *timestamp {
	return &timestamp{
		time:    time,
		process: process,
	}
}

func (ts timestamp) String() string {
	return fmt.Sprintf("<T%d:P%d>", ts.time, ts.process)
}

func less(a, b timestamp) bool {
	if a.time == b.time {
		return a.process < b.process
	}
	return a.time < b.time
}
