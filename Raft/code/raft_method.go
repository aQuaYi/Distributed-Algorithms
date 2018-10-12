package raft

import (
	"bytes"
	"io"
	"math/rand"
	"time"
)

import "github.com/aQuaYi/Distributed-Algorithms/Raft/code/labgob"

// GetState 可以获取 raft 对象的状态
// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	var term int
	var isLeader bool
	// Your code here (2A).

	// 添加 RLock 是为了避免在 Lock 期间读取到数据
	rf.rwmu.RLock()
	term = rf.currentTerm
	isLeader = rf.state == LEADER
	rf.rwmu.RUnlock()

	return term, isLeader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).

	// 不需要上锁的原因是，persist 总是在锁定的环境中被调用

	buffer := new(bytes.Buffer)

	e := labgob.NewEncoder(buffer)

	e.Encode(rf.currentTerm)
	e.Encode(rf.votedFor)
	for i := 1; i < len(rf.logs); i++ {
		log := rf.logs[i]
		e.Encode(log.LogTerm)
		e.Encode(&log.Command)
	}

	data := buffer.Bytes()
	rf.persister.SaveRaftState(data)

	debugPrintf("%s persisted!", rf)
}

//
// restore previously persisted state.
// func (*Deocder) Decode(e interface{}) error
//     Decode reads the next value from the input stream and stores it in
//     the data represented by the empty interface value. If e is nil, the
//     value will be discarded.
//     Otherwise, the value underlying e must be a pointer to the correct
//     type for the next data item received. If the input is at EOF,
//     Decode returns io.EOF and does not modify e
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).

	rf.rwmu.Lock()
	defer rf.rwmu.Unlock()

	buffer := bytes.NewBuffer(data)
	d := labgob.NewDecoder(buffer)

	var currentTerm int
	var votedFor int

	if d.Decode(&currentTerm) != nil ||
		d.Decode(&votedFor) != nil {
		debugPrintf("error in decode currentTerm and votedFor, err: %v\n", d.Decode(&currentTerm))
	} else {
		rf.currentTerm = currentTerm
		rf.votedFor = votedFor
	}

	for {
		var log LogEntry
		if err := d.Decode(&log.LogTerm); err != nil {
			if err == io.EOF {
				break
			} else {
				debugPrintf("error when decode log, err: %v\n", err)
			}
		}

		if err := d.Decode(&log.Command); err != nil {
			panic(err)
		}
		rf.logs = append(rf.logs, log)
	}

	debugPrintf("%s readPersisted!", rf)
}

// Start 给 server 发送命令
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
//
func (rf *Raft) Start(command interface{}) (index, term int, isLeader bool) {

	index = -1
	term = -1
	isLeader = false

	// Your code here (2B).
	// if command received from client:
	// append entry to local log, respond after entry applied to state machine
	rf.rwmu.Lock()
	defer rf.rwmu.Unlock()

	if rf.state != LEADER {
		return
	}

	// 修改结果值
	index = len(rf.logs)
	term = rf.currentTerm
	isLeader = true

	// 生成新的 entry
	entry := &LogEntry{
		LogIndex: index,
		LogTerm:  term,
		Command:  command,
	}

	// 修改 rf 的属性
	rf.logs = append(rf.logs, *entry)
	rf.nextIndex[rf.me] = len(rf.logs)
	rf.matchIndex[rf.me] = len(rf.logs) - 1

	debugPrintf("%s 添加了新的 entry:%v", rf, *entry)

	return
}

// Kill is
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
	debugPrintf("S#%d Killing", rf.me)

	// 关闭前，先去检查一遍 apply
	rf.toCheckApplyChan <- struct{}{}

	close(rf.shutdownChan)

	rf.shutdownWG.Wait()
}

func (rf *Raft) resetElectionTimer() {
	timeout := time.Duration(150+rand.Int63n(151)) * time.Millisecond
	rf.electionTimer.Reset(timeout)
	debugPrintf("%s election timer 已经重置, 时长： %s", rf, timeout)
}
