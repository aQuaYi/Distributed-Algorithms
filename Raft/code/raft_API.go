package raft

import "github.com/aQuaYi/Distributed-Algorithms/Raft/code/labrpc"

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

// Make is
func Make(peers []*labrpc.ClientEnd, me int, persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := newRaft(peers, me, persister, applyCh)

	go rf.checkApplyLoop(applyCh)

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	return rf
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

// Kill is
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
	debugPrintf("R%d Killing", rf.me)

	// 关闭前，先去检查一遍 apply
	rf.toCheckApplyChan <- struct{}{}

	close(rf.shutdownChan)

	rf.shutdownWG.Wait()
}
