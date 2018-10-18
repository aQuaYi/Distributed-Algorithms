package raft

import (
	"log"
	"math/rand"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	debugPrintf("程序开始运行")
}

const (
	// heartBeatInterval 发送心跳的时间间隔，ms
	heartBeatInterval = 50 * time.Millisecond
	// minElectionInterval 选举过期的最小时间间隔，ms
	minElectionInterval = heartBeatInterval * 10
	// minElectionInterval 选举过期的最大时间间隔，ms
	maxElectionInterval = minElectionInterval * 8 / 5

	// 按照论文 5.6 Timing and availability 的要求
	// HBInterval 和 MinElectionInterval 相差了一个数量级
)

func electionTimeout() time.Duration {
	interval := int(minElectionInterval) +
		rand.Intn(int(maxElectionInterval-minElectionInterval))
	return time.Duration(interval)
}

// needDebug for Debugging
const needDebug = false

// debugPrintf 根据设置打印输出
func debugPrintf(format string, a ...interface{}) {
	if needDebug {
		log.Printf(format, a...)
	}
}
