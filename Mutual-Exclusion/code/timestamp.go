package mutual

import "fmt"

// Timestamp 是用于全局排序的接口
type Timestamp interface {
	// Less 比较两个 Timestamp 的大小
	Less(Timestamp) bool
	// Equal 判断两个 Timestamp 是否相等
	isEqual(Timestamp) bool
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

// TODO: 把返回值改成 接口
func newTimestamp(time, process int) Timestamp {
	return &timestamp{
		time:    time,
		process: process,
	}
}

func (ts *timestamp) String() string {
	return fmt.Sprintf("<T%d:P%d>", ts.time, ts.process)
}

func (ts *timestamp) Less(tsi Timestamp) bool {
	ts2, _ := tsi.(*timestamp)
	if ts.time == ts2.time {
		return ts.process < ts2.process
	}
	return ts.time < ts2.time
}

func (ts *timestamp) isEqual(tsi Timestamp) bool {
	return !ts.Less(tsi) && !tsi.Less(ts)
}

func (ts *timestamp) Time() int {
	return ts.time
}

func (ts *timestamp) Process() int {
	return ts.process
}

// TODO: 删除此处内容
func less(ts, b timestamp) bool {
	if ts.time == b.time {
		return ts.process < b.process
	}
	return ts.time < b.time
}
