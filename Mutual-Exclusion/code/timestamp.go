package mutual

import "fmt"

type timestamp struct {
	time, process int
}

func (ts timestamp) String() string {
	return fmt.Sprintf("<T%d:P%d>", ts.time, ts.process)
}
