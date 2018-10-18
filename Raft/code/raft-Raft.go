package raft

import (
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

	// TODO: 继续添加代码
}
