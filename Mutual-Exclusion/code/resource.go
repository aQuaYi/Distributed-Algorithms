package mutualexclusion

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/montanaflynn/stats"
)

// Resource 是 Process 占用资源的接口
type Resource interface {
	// Occupy 表示占用资源
	Occupy(Timestamp)
	// Release 表示释放资源
	Release(Timestamp)
}

type resource struct {
	lastOccupiedBy Timestamp      // 记录上次占用资源的 timestamp
	occupiedBy     Timestamp      // 记录当前占用资源的 timestamp, nil 表示资源未被占用
	timestamps     []Timestamp    // 按顺序保存占用资源的 timestamp
	times          []time.Time    // 记录每次占用资源的起止时间，用于分析算法的效率
	wg             sync.WaitGroup // 完成全部占用前，阻塞主 goroutine
}

func newResource(times int) *resource {
	r := &resource{
		lastOccupiedBy: newTimestamp(-1, -1),
	}
	r.wg.Add(times)
	return r
}

func (r *resource) wait() {
	r.wg.Wait()
}

func (r *resource) Occupy(ts Timestamp) {
	r.times = append(r.times, time.Now())

	if r.occupiedBy != nil {
		msg := fmt.Sprintf("资源正在被 %s 占据，%s 却想获取资源。", r.occupiedBy, ts)
		panic(msg)
	}

	if !r.lastOccupiedBy.Less(ts) {
		msg := fmt.Sprintf("资源上次被 %s 占据，这次 %s 却想占据资源。", r.lastOccupiedBy, ts)
		panic(msg)
	}

	r.occupiedBy = ts
	r.timestamps = append(r.timestamps, ts)
	debugPrintf("~~~ @resource: %s occupied ~~~ ", ts)
}

func (r *resource) Release(ts Timestamp) {
	if !r.occupiedBy.IsEqual(ts) {
		msg := fmt.Sprintf("%s 想要释放正在被 P%s 占据的资源。", ts, r.occupiedBy)
		panic(msg)
	}

	r.lastOccupiedBy, r.occupiedBy = ts, nil

	r.times = append(r.times, time.Now())

	debugPrintf("~~~ @resource: %s released ~~~ ", ts)

	r.wg.Done() // 完成一次占用

}

func (r *resource) report() string {
	var b strings.Builder

	// 统计资源被占用的时间
	size := len(r.times)
	totalTime := r.times[size-1].Sub(r.times[0])
	format := "resource 被占用了 %s， "
	fmt.Fprintf(&b, format, totalTime)

	// 计算占用率
	busys := make([]float64, 0, size/2)
	idles := make([]float64, 0, size/2)

	var i int
	for i = 0; i+2 < size; i += 2 {
		busys = append(busys, float64(r.times[i+1].Sub(r.times[i]).Nanoseconds()))
		idles = append(idles, float64(r.times[i+2].Sub(r.times[i+1]).Nanoseconds()))
	}
	busys = append(busys, float64(r.times[i+1].Sub(r.times[i]).Nanoseconds()))

	busy, _ := stats.Sum(busys)
	idle, _ := stats.Sum(idles)
	total := busy + idle
	rate := busy * 100 / total

	format = "占用比率为 %4.2f%%。\n"
	fmt.Fprintf(&b, format, rate)

	// 计算资源占用时间的均值和方差
	format = "资源占用: %s\n"
	fmt.Fprintf(&b, format, statisticAnalyze(busys))

	// 计算资源空闲间隙的均值和方差
	format = "资源空闲: %s\n"
	fmt.Fprintf(&b, format, statisticAnalyze(idles))

	return b.String()
}

func statisticAnalyze(floats []float64) string {
	format := "min %8.2fus, max %8.2fus, mean %8.2fus, sd %8.2f"
	min, _ := stats.Min(floats)
	max, _ := stats.Max(floats)
	mean, _ := stats.Mean(floats)
	sd, _ := stats.StandardDeviation(floats)
	return fmt.Sprintf(format, min/1000, max/1000, mean/1000, sd/1000)
}
