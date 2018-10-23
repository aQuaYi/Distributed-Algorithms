package raft

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ElectionTimeout(t *testing.T) {
	ast := assert.New(t)
	for i := 0; i < 1000; i++ {
		rate := electionTimeout() / heartBeat
		ast.True(rate >= 10, "electionTimeout 没有比 heartBeatInterval 大 10 倍")
	}
}

func Test_heartBeat_isInRange(t *testing.T) {
	ast := assert.New(t)
	minInterval := 30 * time.Millisecond
	maxInterval := 100 * time.Millisecond
	isInRange := minInterval <= heartBeat && heartBeat <= maxInterval
	ast.True(isInRange, " heartBeat 设置的过大或过小")
}
