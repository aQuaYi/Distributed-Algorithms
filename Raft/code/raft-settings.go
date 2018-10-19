package raft

import (
	"log"
	"math/rand"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	DPrintf("程序开始运行")
}

const (
	// heartBeat 发送心跳的时间间隔，ms
	heartBeat = 50 * time.Millisecond
	// minElection 选举过期的最小时间间隔，ms
	minElection = heartBeat * 10
	// minElectionInterval 选举过期的最大时间间隔，ms
	maxElection = minElection * 8 / 5

	// 按照论文 5.6 Timing and availability 的要求
	// heartBeat 和 minElection 需要相差了一个数量级
)

func electionTimeout() time.Duration {
	interval := int(minElection) +
		rand.Intn(int(maxElection-minElection))
	return time.Duration(interval)
}
