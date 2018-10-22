package raft

import "log"

// Debugging is
const Debugging = 0

// DPrintf is
func DPrintf(format string, a ...interface{}) {
	if Debugging > 0 {
		log.Printf(format, a...)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
