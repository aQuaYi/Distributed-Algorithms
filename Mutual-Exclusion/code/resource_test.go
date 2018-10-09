package mutualexclusion

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_resource_occupyAndRelease(t *testing.T) {
	// 避免 debugprint 输出
	temp := needDebug
	needDebug = false
	//
	ast := assert.New(t)
	//
	p := 0
	ts := newTimestamp(0, p)
	r := newResource(1)
	// 占用
	r.Occupy(ts)
	ast.Equal(ts, r.occupiedBy)
	// 释放
	r.Release(ts)
	r.wait()
	ast.Equal(ts, r.lastOccupiedBy)
	ast.Equal(ts, r.timestamps[0])
	// 还原 needDebug
	needDebug = temp
}

func Test_resource_occupy_occupyInvalidResource(t *testing.T) {
	// 避免 debugprint 输出
	temp := needDebug
	needDebug = false
	//
	ast := assert.New(t)
	//
	p0 := 0
	p1 := 1
	ts0 := newTimestamp(0, p0)
	ts1 := newTimestamp(1, p1)
	r := newResource(1)
	r.Occupy(ts0)
	//
	expected := fmt.Sprintf("资源正在被 %s 占据，%s 却想获取资源。", ts0, ts1)
	ast.PanicsWithValue(expected, func() { r.Occupy(ts1) })
	// 还原 needDebug
	needDebug = temp
}

func Test_resource_occupy_panicOfEarlyTimestampWantToOccupy(t *testing.T) {
	// 避免 debugprint 输出
	temp := needDebug
	needDebug = false
	//
	ast := assert.New(t)
	//
	ts0 := newTimestamp(0, 1)
	ts1 := newTimestamp(1, 1)
	r := newResource(2)
	r.Occupy(ts1)
	r.Release(ts1)
	//
	expected := fmt.Sprintf("资源上次被 %s 占据，这次 %s 却想占据资源。", ts1, ts0)
	ast.PanicsWithValue(expected, func() { r.Occupy(ts0) })
	// 还原 needDebug
	needDebug = temp
}

func Test_resource_report(t *testing.T) {
	// 避免 debugprint 输出
	temp := needDebug
	needDebug = false
	//
	ast := assert.New(t)
	//
	p := 0
	ts0 := newTimestamp(0, p)
	ts1 := newTimestamp(1, p)
	r := newResource(3)
	r.Occupy(ts0)
	r.Release(ts0)
	r.Occupy(ts1)
	r.Release(ts1)
	now := time.Now()
	r.times[0] = now
	r.times[1] = now.Add(100 * time.Second)
	r.times[2] = now.Add(200 * time.Second)
	r.times[3] = now.Add(400 * time.Second)
	//
	report := r.report()
	ast.True(strings.Contains(report, "75.00%"), report)
	//
	ast.Equal(4, len(r.times), "资源被占用了 2 次，但是 r.times 的长度不等于 4")
	// 还原 needDebug
	needDebug = temp
}

func Test_resource_Occupy_lenOfTimes(t *testing.T) {
	// 避免 debugprint 输出
	temp := needDebug
	needDebug = false
	//
	ast := assert.New(t)
	//
	times := 100
	r := newResource(times)
	go func(max int) {
		time, p := 0, 0
		for i := 0; i < max; i++ {
			if i%2 == 0 {
				time++
			} else {
				p++
			}
			ts := newTimestamp(time, p)
			r.Occupy(ts)
			r.Release(ts)
		}
	}(times)
	r.wait()
	expected := times * 2
	actual := len(r.times)
	ast.Equal(expected, actual)
	// 还原 needDebug
	needDebug = temp
}

func Test_resource_Release_panicOfReleaseByOther(t *testing.T) {
	// 避免 debugprint 输出
	temp := needDebug
	needDebug = false
	//
	ast := assert.New(t)
	//
	r := newResource(1)
	ts0 := newTimestamp(0, 1)
	ts1 := newTimestamp(1, 1)
	r.Occupy(ts0)
	expected := fmt.Sprintf("%s 想要释放正在被 P%s 占据的资源。", ts1, ts0)
	ast.PanicsWithValue(expected, func() { r.Release(ts1) })
	// 还原 needDebug
	needDebug = temp
}
