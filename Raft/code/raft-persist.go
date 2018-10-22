package raft

import (
	"bytes"
	"encoding/gob"

	"github.com/aQuaYi/Distributed-Algorithms/Raft/code/labgob"
)

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	w := new(bytes.Buffer)
	e := labgob.NewEncoder(w)
	e.Encode(rf.currentTerm)
	e.Encode(rf.votedFor)
	e.Encode(rf.logs)
	data := w.Bytes()
	rf.persister.SaveRaftState(data)

	DPrintf("%s PERSISTED", rf)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var currentTerm int
	// var votedFor int
	// var logs []LogEntry
	// if d.Decode(&currentTerm) != nil ||
	// d.Decode(&votedFor) != nil ||
	// d.Decode(&logs) != nil {
	// panic("readPersist 无法 Decode")
	// } else {
	// rf.currentTerm = currentTerm
	// rf.votedFor = votedFor
	// rf.logs = logs
	// }

	r := bytes.NewBuffer(data)
	d := gob.NewDecoder(r)
	d.Decode(&rf.currentTerm)
	d.Decode(&rf.votedFor)
	d.Decode(&rf.logs)
}
