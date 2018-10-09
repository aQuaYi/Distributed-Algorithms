package main

import "fmt"

// Timestamp 是用于全局排序的接口
type Timestamp interface {
	// Less 比较两个 Timestamp 的大小
	Less(interface{}) bool
	// Equal 判断两个 Timestamp 是否相等
	IsEqual(interface{}) bool
	// IsBefore 在比较同一个 clock 的时间，所以，不需要 process
	IsBefore(int) bool
	// String 输出 Timestamp 的内容
	String() string
}

type timestamp struct {
	time, process int
}

func newTimestamp(time, process int) Timestamp {
	return &timestamp{
		time:    time,
		process: process,
	}
}

func (ts *timestamp) String() string {
	return fmt.Sprintf("<T%d:P%d>", ts.time, ts.process)
}

func (ts *timestamp) Less(tsi interface{}) bool {
	ts2 := tsi.(*timestamp)
	// 这就是将局部顺序推广到全局顺序的关键
	if ts.time == ts2.time {
		return ts.process < ts2.process
	}
	return ts.time < ts2.time
}

func (ts *timestamp) IsEqual(tsi interface{}) bool {
	if tsi == nil {
		return false
	}
	ts2 := tsi.(*timestamp)
	return ts.time == ts2.time && ts.process == ts2.process
}

func (ts *timestamp) IsBefore(t int) bool {
	return ts.time < t
}
