package mutual

import (
	"fmt"
	"time"
)

// Resource 是 Process 占用资源的接口
type Resource interface {
	// Occupy 表示占用资源
	Occupy(Timestamp)
	// Release 表示释放资源
	Release(Timestamp)
}
type resource struct {
	occupiedBy Timestamp
	timestamps []Timestamp
	times      []time.Time
}

func newResource() *resource {
	return &resource{
		occupiedBy: nil,
	}
}

func (r *resource) Occupy(ts Timestamp) {
	if r.occupiedBy != nil {
		msg := fmt.Sprintf("资源正在被 %s 占据，%s 却想获取资源。", r.occupiedBy, ts)
		panic(msg)
	}
	r.occupiedBy = ts
	r.timestamps = append(r.timestamps, ts)
	r.times = append(r.times, time.Now())
	debugPrintf("~~~ @resource: %s occupied ~~~ ", ts)
}

func (r *resource) Release(ts Timestamp) {
	if !r.occupiedBy.isEqual(ts) {
		msg := fmt.Sprintf("%s 想要释放正在被 P%s 占据的资源。", ts, r.occupiedBy)
		panic(msg)
	}
	r.occupiedBy = nil
	r.times = append(r.times, time.Now())
	debugPrintf("~~~ @resource: %s released ~~~ ", ts)
}

func (r *resource) report() string {
	occupiedTime := time.Duration(0)
	size := len(r.times)
	for i := 0; i+1 < size; i += 2 {
		occupiedTime += r.times[i+1].Sub(r.times[i])
	}
	totalTime := r.times[size-1].Sub(r.times[0])
	rate := occupiedTime.Nanoseconds() * 10000 / totalTime.Nanoseconds()
	format := "resource 的占用比率为 %02d.%02d%%"
	return fmt.Sprintf(format, rate/100, rate%100)
}

func (r *resource) isSortedOccupied() bool {
	size := len(r.timestamps)
	for i := 1; i < size; i++ {
		if !r.timestamps[i-1].Less(r.timestamps[i]) {
			return false
		}
	}
	return true
}
