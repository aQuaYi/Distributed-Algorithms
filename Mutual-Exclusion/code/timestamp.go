package main

import "fmt"

// Timestamp 是用于全局排序的接口
type Timestamp interface {
	// Less 比较两个 Timestamp 的大小
	Less(interface{}) bool
	// Equal 判断两个 Timestamp 是否相等
	isEqual(interface{}) bool
	// Time 输出的是 time 属性的值
	Time() int
	// Process 输出 process 属性的值
	Process() int
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
	ts2, ok := tsi.(*timestamp)
	if !ok {
		panic("ts.Less：无法转换 tsi 到 *timestamp 类型")
	}
	// 这就是将局部顺序推广到全局顺序的关键
	if ts.time == ts2.time {
		return ts.process < ts2.process
	}
	return ts.time < ts2.time
}

func (ts *timestamp) isEqual(tsi interface{}) bool {
	ts2, ok := tsi.(*timestamp)
	if !ok {
		panic("ts.Less：无法转换 tsi 到 *timestamp 类型")
	}
	return ts.time == ts2.time && ts.process == ts2.process
}

func (ts *timestamp) Time() int {
	return ts.time
}

func (ts *timestamp) Process() int {
	return ts.process
}
