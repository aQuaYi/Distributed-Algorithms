package raft

import (
	"fmt"
	"sync"

	"github.com/aQuaYi/Distributed-Algorithms/Raft/code/labrpc"
)

const (
	// NOBODY used for Raft.votedFor, means vote for none
	NOBODY = -1
)

// Raft is
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	/* ↓ state of raft on Figure 2 ↓ */

	// Persistent state on all servers:
	currentTerm int        // latest term server has seen. Initialized to 0.
	votedFor    int        // candidateID that received vote in current Term
	logs        []LogEntry // NOTICE: first LogEntry.LogIndex is 1

	// Volatile state on all servers: initialized to 0, increase monotonically
	commitIndex int // index of highest log entry known to be committed
	lastApplied int // index of highest log entry known to be applied to state machine

	// Volatile state on leader:
	// nextIndex : for each server, index of the next log entry to send to that server
	// initialized to leader last LogIndex+1
	nextIndex []int
	// matchIndex : for each server, index of highest log entry known to be replicated on server
	// initialized to 0, increases monotonically
	matchIndex []int

	/* ↑ state of raft on Figure 2 ↑ */

	state     state
	voteCount int

	chanApply chan ApplyMsg

	//channel
	chanCommit    chan struct{}
	chanHeartBeat chan struct{}
	chanBeElected chan struct{}
}

func (rf *Raft) String() string {
	return fmt.Sprintf(" <R%d:T%d:%s> ", rf.me, rf.currentTerm, rf.state)
}

func (rf *Raft) details() string {
	postfix := ""
	if rf.state == LEADER {
		postfix = fmt.Sprintf(", nextIndex%v, matchIndex%v", rf.nextIndex, rf.matchIndex)
	}
	return fmt.Sprintf("@@ votedFor:%2d, commitIndex:%d, lastApplied:%d, logs:%v%s @@",
		rf.votedFor, rf.commitIndex, rf.lastApplied, rf.logs, postfix)
}
