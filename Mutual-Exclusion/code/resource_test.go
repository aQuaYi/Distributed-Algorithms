package mutual

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
	defer func() { needDebug = temp }()
	//
	ast := assert.New(t)
	//
	p := 0
	ts := timestamp{time: 0, process: p}
	r := newResource()
	r.occupy2(ts)
	//
	ast.Equal(ts, r.occupiedBy2)
	ast.Equal(ts, r.timestamps[0])

}

func Test_resource_occupy_occupyInvalidResource(t *testing.T) {
	// 避免 debugprint 输出
	temp := needDebug
	needDebug = false
	defer func() { needDebug = temp }()
	//
	ast := assert.New(t)
	//
	p0 := 0
	p1 := 1
	ts0 := timestamp{time: 0, process: p0}
	ts1 := timestamp{time: 1, process: p1}
	r := newResource()
	r.occupy2(ts0)
	//
	expected := fmt.Sprintf("资源正在被 %s 占据，%s 却想获取资源。", ts0, ts1)
	ast.PanicsWithValue(expected, func() { r.occupy2(ts1) })
}

func Test_resource_report(t *testing.T) {
	// 避免 debugprint 输出
	temp := needDebug
	needDebug = false
	defer func() { needDebug = temp }()
	//
	ast := assert.New(t)
	//
	p := 0
	ts := timestamp{time: 0, process: p}
	r := newResource()
	r.occupy2(ts)
	r.release2(ts)
	r.occupy2(ts)
	r.release2(ts)
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
}

func Test_resource_timestamps(t *testing.T) {
	// 避免 debugprint 输出
	temp := needDebug
	needDebug = false
	defer func() { needDebug = temp }()
	//
	ast := assert.New(t)
	//
	p := 0
	ts := timestamp{time: 0, process: p}
	r := newResource()
	times := 100
	//
	for i := 0; i < times; i++ {
		if i%2 == 0 {
			ts.time++
		} else {
			ts.process++
		}
		r.occupy2(ts)
		r.release2(ts)
	}
	//
	expected := times * 2
	actual := len(r.times)
	ast.Equal(expected, actual)
	//
	ast.True(r.isSortedOccupied())
}
