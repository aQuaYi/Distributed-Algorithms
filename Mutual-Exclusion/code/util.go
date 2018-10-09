package mutualexclusion

import (
	"log"
	"math/rand"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	debugPrintf("程序开始运行")
	rand.Seed(time.Now().UnixNano())
}

var needDebug = false

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
