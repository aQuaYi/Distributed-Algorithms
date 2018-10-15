package raft

import (
	"bytes"
	"io"
	"math/rand"
	"time"
)

import "github.com/aQuaYi/Distributed-Algorithms/Raft/code/labgob"

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
// func (*Decoder) Decode(e interface{}) error
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

func (rf *Raft) resetElectionTimer() {
	timeout := time.Duration(150+rand.Int63n(151)) * time.Millisecond
	rf.electionTimer.Reset(timeout)
	debugPrintf("%s election timer 已经重置, 时长： %s", rf, timeout)
}
