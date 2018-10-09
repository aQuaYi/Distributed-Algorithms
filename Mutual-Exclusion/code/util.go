package mutualexclusion

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	debugPrintf("程序开始运行")
	rand.Seed(time.Now().UnixNano())
}

var needDebug = false

// 读取和修改 needebug 前需要上锁
var rwm sync.RWMutex // TODO: 删除此处的锁

// debugPrintf 根据设置打印输出
func debugPrintf(format string, a ...interface{}) {
	rwm.RLock() // TODO: 删除此处的锁
	if needDebug {
		log.Printf(format, a...)
	}
	rwm.RUnlock()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
