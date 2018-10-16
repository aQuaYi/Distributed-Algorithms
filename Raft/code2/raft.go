package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isLeader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//
import (
	"sync"
	"time"

	"github.com/aQuaYi/Distributed-Algorithms/Raft/code/labrpc"
)

// import "bytes"
// import "labgob"

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

	// from Figure 2

	// Persistent state on call servers
	currentTerm int        // 此 server 当前所处的 term 编号
	votedFor    int        // 此 server 在此 term 中投票给了谁，是 peers 中的索引号
	logs        []LogEntry // 此 server 中保存的 logs

	// Volatile state on all servers:
	commitIndex int // logs 中已经 committed 的 log 的最大索引号
	lastApplied int // logs 中已经执行的最后的 log 的索引号

	// Volatile state on leaders:
	nextIndex  []int // 下一个要发送给 follower 的 log 的索引号
	matchIndex []int // leader 与 follower 共有的 log 的最大的索引号

	// Raft 作为 FSM 管理自身状态所需的属性
	state state
	// handlers map[fsmState]map[fsmEvent]fsmHandler

	// 超时，就由 FOLLOWER 变 CANDIDATE
	electionTimer *time.Timer

	// 用于通知 raft 已经关闭的信息
	shutdownChan chan struct{}
	shutdownWG   sync.WaitGroup

	// 当 rf 接收到合格的 rpc 信号时，会通过 resetElectionChan 发送信号
	resetElectionChan chan struct{}

	// candidate 或 leader 中途转变为 follower 的话，就关闭这个 channel 来发送信号
	// 因为，同一个 rf 不可能既是 candidate 又是 leader
	// 所以，用来通知的 channel 只要有一个就好了
	convertToFollowerChan chan struct{}

	// logs 中添加了新的 entries 以后，会通过这个发送信号
	toCheckApplyChan chan struct{}

	// 关闭，则表示需要终结此次 election
	endElectionChan chan struct{}

	// 2018-10-15 新添加的属性
	// closeElectionLoopChan 成为 Leader 时，关闭 electionLoop
	closeElectionLoopChan chan struct{} // TODO: 在 newRaft 中添加
}

// Make is
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	return rf
}
