package mutual

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_resource_occupy_occupyValidResource(t *testing.T) {
	ast := assert.New(t)
	//
	p := 0
	r := newResource()
	r.occupy2(timestamp{
		time:    0,
		process: p,
	})
	//
	expected := p
	actual := r.occupiedBy
	ast.Equal(expected, actual)
}

func Test_resource_occupy_occupyInvalidResource(t *testing.T) {
	ast := assert.New(t)
	//
	p0 := 0
	p1 := 1
	ts0 := timestamp{time: 0, process: p0}
	ts1 := timestamp{time: 1, process: p1}
	r := newResource()
	r.occupy2(ts0)
	//
	expected := fmt.Sprintf("资源正在被 P%d 占据，P%d 却想获取资源。", p0, p1)
	ast.PanicsWithValue(expected, func() { r.occupy2(ts1) })
}

func Test_resource_report(t *testing.T) {
	ast := assert.New(t)
	now := time.Now()
	r := &resource{}

	r.times = append(r.times, now, now.Add(100*time.Second))
	r.times = append(r.times, now.Add(200*time.Second), now.Add(400*time.Second))
	report := r.report()
	ast.True(strings.Contains(report, "75.00%"), report)
}
