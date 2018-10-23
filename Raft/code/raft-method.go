package raft

// 这里的方法都是被内部引用的，所以无需加锁

func (rf *Raft) getLastIndex() int {
	return rf.logs[len(rf.logs)-1].LogIndex
}

func (rf *Raft) getBaseIndex() int {
	return rf.logs[0].LogIndex
}

func (rf *Raft) getLastTerm() int {
	return rf.logs[len(rf.logs)-1].LogTerm
}

func (rf *Raft) isLeader() bool {
	return rf.state == LEADER
}

func (rf *Raft) isCandidate() bool {
	return rf.state == CANDIDATE
}

func (rf *Raft) isFollower() bool {
	return rf.state == FOLLOWER
}
