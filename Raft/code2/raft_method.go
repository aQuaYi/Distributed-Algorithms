package raft

func (rf *Raft) getLastIndex() int {
	return rf.logs[len(rf.logs)-1].LogIndex
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
