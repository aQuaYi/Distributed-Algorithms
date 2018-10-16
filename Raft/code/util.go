package raft

import (
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	debugPrintf("程序开始运行")
}

const (
	// HBInterval 发送心跳的时间间隔，ms
	HBInterval = 50
	// MinElectionInterval 选举过期的最小时间间隔，ms
	MinElectionInterval = 500
	// MaxElectionInterval = MinElectionInterval * 2

	// 按照论文 5.6 Timing and availability 的要求
	// HBInterval 和 MinElectionInterval 相差了一个数量级
)

// needDebug for Debugging
const needDebug = true

// debugPrintf 根据设置打印输出
func debugPrintf(format string, a ...interface{}) {
	if needDebug {
		log.Printf(format, a...)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
