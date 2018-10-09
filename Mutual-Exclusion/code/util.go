package mutualexclusion

import (
	"log"
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
