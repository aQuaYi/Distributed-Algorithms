package raft

import (
	"fmt"
	"time"

	"github.com/aQuaYi/Distributed-Algorithms/Raft/code/labrpc"
)

/**
 * // create a new Raft server instance:
 * rf := Make(peers, me, persister, applyCh)
 *
 * // start agreement on a new log entry:
 * rf.Start(command interface{}) (index, term, isLeader)
 *
 * // ask a Raft for its current term, and whether it thinks it is leader
 * rf.GetState() (term, isLeader)
 *
 * // each time a new entry is committed to the log, each Raft peer
 * // should send an ApplyMsg to the service (or tester).
 * type ApplyMsg
 *
 */

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
	rf.state = FOLLOWER
	rf.votedFor = NOBODY
	rf.logs = append(rf.logs, LogEntry{LogIndex: 0, LogTerm: 0, Command: 0})
	rf.currentTerm = 0
	rf.chanCommit = make(chan struct{}, 100)
	rf.chanHeartbeat = make(chan struct{}, 100)
	rf.chanGrantVote = make(chan struct{}, 100)
	rf.chanLeader = make(chan struct{}, 100)
	rf.chanApply = applyCh

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	go rf.stateLoop()

	go rf.applyLoop()

	return rf
}

func (rf *Raft) stateLoop() {
	for {
		switch rf.state {
		case FOLLOWER:
			select {
			case <-time.After(electionTimeout()):
				rf.state = CANDIDATE
			case <-rf.chanHeartbeat:
			case <-rf.chanGrantVote:
			}
		case CANDIDATE:
			rf.newElection()
		case LEADER:
			rf.newHeartBeat()
		}
	}
}

func (rf *Raft) newElection() {
	rf.mu.Lock()
	rf.currentTerm++

	DPrintf("%s begin new election\n", rf)

	rf.votedFor = rf.me
	rf.voteCount = 1

	rf.persist()

	rf.mu.Unlock()

	go rf.broadcastRequestVote()

	select {
	case <-rf.chanHeartbeat:
		rf.state = FOLLOWER
		DPrintf("%s receives chanHeartbeat", rf)
	case <-rf.chanLeader:
		rf.mu.Lock()
		rf.state = LEADER
		DPrintf("%s is Leader now", rf)
		rf.nextIndex = make([]int, len(rf.peers))
		rf.matchIndex = make([]int, len(rf.peers))
		for i := range rf.peers {
			rf.nextIndex[i] = rf.getLastIndex() + 1
			rf.matchIndex[i] = 0
		}
		rf.mu.Unlock()
	case <-time.After(electionTimeout()):
	}
}

func (rf *Raft) newHeartBeat() {
	DPrintf("%s broadcastAppendEntries", rf)
	rf.broadcastAppendEntries()
	<-time.After(heartBeat)
}

func (rf *Raft) applyLoop() {
	for {
		select { // TODO: select 只有一个 case ，可以删掉
		case <-rf.chanCommit:
			rf.mu.Lock()
			commitIndex := rf.commitIndex
			baseIndex := rf.getBaseIndex()
			for i := rf.lastApplied + 1; i <= commitIndex; i++ {
				msg := ApplyMsg{
					CommandValid: true,
					CommandIndex: i,
					Command:      rf.logs[i-baseIndex].Command,
				}
				rf.chanApply <- msg
				rf.lastApplied = i
				DPrintf("%s ApplyMSG: %s %s", rf, msg, rf.details())
			}
			rf.mu.Unlock()
		}
	}
}

// Start is
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and
// ** return immediately, without waiting for the log appends to complete. **
// there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	// Your code here (2B).
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if !rf.isLeader() {
		return -1, -1, false
	}

	DPrintf("%s Start %v", rf, command)

	logIndex := rf.getLastIndex() + 1
	term := rf.currentTerm
	isLeader := rf.isLeader()

	rf.logs = append(rf.logs,
		LogEntry{
			LogIndex: logIndex,
			LogTerm:  term,
			Command:  command,
		}) // append new entry from client

	rf.persist()

	// Your code above
	return logIndex, term, isLeader
}

// GetState is
// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isLeader bool
	// Your code here (2A).

	term = rf.currentTerm
	isLeader = rf.isLeader()

	// Your code above (2A)
	return term, isLeader
}

// ApplyMsg is
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in Lab 3 you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh; at that point you can add fields to
// ApplyMsg, but set CommandValid to false for these other uses.
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

func (m ApplyMsg) String() string {
	return fmt.Sprintf("ApplyMsg{Valid:%t,Index:%d,Command:%v}", m.CommandValid, m.CommandIndex, m.Command)
}
