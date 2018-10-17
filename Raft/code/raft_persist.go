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

	// 不需要上锁的原因是，persist 总是在锁定的环境中被调用

	w := new(bytes.Buffer)

	e := labgob.NewEncoder(w)
	e.Encode(rf.currentTerm)
	e.Encode(rf.votedFor)
	e.Encode(rf.logs)

	data := w.Bytes()

	rf.persister.SaveRaftState(data)

	debugPrintf("%s persisted!", rf)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	// Your code here.
	// Example:
	r := bytes.NewBuffer(data)
	d := gob.NewDecoder(r)
	d.Decode(&rf.currentTerm)
	d.Decode(&rf.votedFor)
	d.Decode(&rf.logs)
}
