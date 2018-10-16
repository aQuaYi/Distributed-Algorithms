package raft

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/aQuaYi/Distributed-Algorithms/Raft/code/labrpc"
)

// Raft implements a single Raft peer.
type Raft struct {
	rwmu  sync.RWMutex        // Lock to protect shared access to this peer's state
	peers []*labrpc.ClientEnd // RPC end points of all peers

	// Persistent state on all servers
	// (Updated on stable storage before responding to RPCs)
	// This implementation doesn't use disk; ti will save and restore
	// persistent state from a Persister object
	// Raft should initialize its state from Persister,
	// and should use it to save its persistent state each item the state changes
	// Use ReadRaftState() and SaveRaftState
	persister *Persister // Object to hold this peer's persisted state
	me        int        // this peer's index into peers[]

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
	state    fsmState
	handlers map[fsmState]map[fsmEvent]fsmHandler

	// 超时，就由 FOLLOWER 变 CANDIDATE
	electionTimer *time.Timer

	// 用于通知 raft 已经关闭的信息
	shutdownChan chan struct{}
	shutdownWG   sync.WaitGroup

	// 当 rf 接收到合格的 rpc 信号时，会通过 heartbeatChan 发送信号
	heartbeatChan chan struct{}

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

func (rf *Raft) String() string {
	return fmt.Sprintf(" <R%d:T%d> ", rf.me, rf.currentTerm)
}

func (rf *Raft) details() string {
	postfix := ""
	if rf.state == LEADER {
		postfix = fmt.Sprintf(", nextIndex%v, matchIndex%v", rf.nextIndex, rf.matchIndex)
	}
	return fmt.Sprintf("@@ R%d:T%d:L%d:%s:%2d, commitIndex:%d, lastApplied:%d, logs:%v%s @@",
		rf.me, rf.currentTerm, len(rf.logs), rf.state, rf.votedFor,
		rf.commitIndex, rf.lastApplied, rf.logs, postfix)
}

func newRaft(peers []*labrpc.ClientEnd, me int, persister *Persister) *Raft {
	rf := &Raft{
		peers:       peers,
		persister:   persister,
		me:          me,
		currentTerm: 0,
		votedFor:    NOBODY,

		// logs 的序列号从 1 开始
		logs:        make([]LogEntry, 1),
		commitIndex: 0,
		lastApplied: 0,

		// 初始状态都是 FOLLOWER
		state: FOLLOWER,

		handlers: make(map[fsmState]map[fsmEvent]fsmHandler, 3),

		// 并不会等 1 秒，很快就会被重置
		electionTimer: time.NewTimer(time.Second),

		// 靠关闭来传递信号，所以，不设置缓冲
		shutdownChan: make(chan struct{}),
		// endElectionChan 需要用到的时候，再赋值

		// 靠数据来传递信号，所以,  设置缓冲
		heartbeatChan:         make(chan struct{}, 3),
		toCheckApplyChan:      make(chan struct{}, 3),
		closeElectionLoopChan: make(chan struct{}, 2),
	}

	rf.addHandlers()

	electionLoop(rf)

	return rf
}

// 触发 election timer 超时，就开始新的选举
func electionLoop(rf *Raft) {
	rf.heartbeatChan <- struct{}{}

	go func() {
		for {
			select {
			case <-rf.electionTimer.C:
				debugPrintf("%s election timeout", rf)
				rf.call(electionTimeOutEvent, nil)
			case <-rf.heartbeatChan:
				debugPrintf("%s 收到 heartbeat 准备重置 election timer", rf)
				rf.resetElectionTimer()
			case <-rf.closeElectionLoopChan:
				debugPrintf(" R%d 在 electionLoop 的 case <- rf.shutdownChan，收到信号。关闭 electionLoop", rf.me)
				return
			}
		}
	}()
}

func (rf *Raft) resetElectionTimer() {
	interval := MinElectionInterval + rand.Intn(MinElectionInterval)
	d := time.Duration(interval) * time.Millisecond
	rf.electionTimer.Reset(d)
	debugPrintf("%s election timer 已经重置, 时长： %s", rf, d)
}
