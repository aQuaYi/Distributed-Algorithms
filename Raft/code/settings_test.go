package raft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ElectionTimeout(t *testing.T) {
	ast := assert.New(t)
	for i := 0; i < 1000; i++ {
		rate := electionTimeout() / heartBeatInterval
		ast.True(rate >= 10, "electionTimeout 没有比 heartBeatInterval 大 10 倍")
	}
}
