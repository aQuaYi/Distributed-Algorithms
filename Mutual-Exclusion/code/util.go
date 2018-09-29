package main

import (
	"log"
	"math/rand"
	"time"
)

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

// NOTICE: 为了保证测试结果的可比性，请勿修改此函数
func randSleep() {
	timeout := time.Duration(2+rand.Intn(2)) * time.Millisecond
	time.Sleep(timeout)
}
