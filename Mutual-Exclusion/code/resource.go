package main

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
	occupiedBy Timestamp
	timestamps []Timestamp
	times      []time.Time
	wg         sync.WaitGroup
	// NOTICE: 为了验证 mutual exclusion 算法，resource 不能加锁
}

func newResource(times int) Resource {
	r := &resource{}
	r.wg.Add(times)
	return r
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
	if !r.occupiedBy.IsEqual(ts) {
		msg := fmt.Sprintf("%s 想要释放正在被 P%s 占据的资源。", ts, r.occupiedBy)
		panic(msg)
	}
	r.occupiedBy = nil
	r.times = append(r.times, time.Now())
	r.wg.Done()
	debugPrintf("~~~ @resource: %s released ~~~ ", ts)
}

func (r *resource) report() string {
	if !r.isSortedOccupied() {
		panic("resource 不是按照顺序被占用的")
	}
	var b strings.Builder
	// 计算占用率
	occupiedTime := time.Duration(0)
	size := len(r.times)
	busys := make([]float64, 0, size/2)
	for i := 0; i+1 < size; i += 2 {
		dif := r.times[i+1].Sub(r.times[i])
		occupiedTime += dif
		busys = append(busys, float64(dif.Nanoseconds()))
	}
	totalTime := r.times[size-1].Sub(r.times[0])
	rate := occupiedTime.Nanoseconds() * 10000 / totalTime.Nanoseconds()
	format := "resource 按照顺序被占用了 %s，占用比率为 %02d.%02d%%。\n"
	fmt.Fprintf(&b, format, totalTime, rate/100, rate%100)
	// 计算资源占用时间的均值和方差
	format = "资源占用 %8d 次，最短 %8.2fus， 最长 %8.2fus， 均值 %8.2fus， 方差 %8.2f。\n"
	minBusy, _ := stats.Min(busys)
	maxBusy, _ := stats.Max(busys)
	meanBusy, _ := stats.Mean(busys)
	sdBusy, _ := stats.StandardDeviation(busys)
	fmt.Fprintf(&b, format, len(busys), minBusy/1000, maxBusy/1000, meanBusy/1000, sdBusy/1000)
	// 计算资源空闲间隙的均值和方差
	idles := make([]float64, 0, size/2-1)
	for i := 1; i+1 < size; i += 2 {
		idles = append(idles,
			float64(r.times[i+1].Sub(r.times[i]).Nanoseconds()),
		)
	}
	format = "资源空闲 %8d 次，最短 %8.2fus， 最长 %8.2fus， 均值 %8.2fus， 方差 %8.2f。\n"
	minIdle, _ := stats.Min(idles)
	maxIdle, _ := stats.Max(idles)
	meanIdle, _ := stats.Mean(idles)
	sdIdle, _ := stats.StandardDeviation(idles)
	fmt.Fprintf(&b, format, len(idles), minIdle/1000, maxIdle/1000, meanIdle/1000, sdIdle/1000)
	//
	return b.String()
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
