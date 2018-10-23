package raft

import "fmt"

// AppendEntriesArgs 是添加 log 的参数
type AppendEntriesArgs struct {
	Term         int        // leader.currentTerm
	LeaderID     int        // leader.me
	PrevLogIndex int        // index of log entry immediately preceding new ones
	PrevLogTerm  int        // term of prevLogIndex entry
	LeaderCommit int        // leader.commitIndex
	Entries      []LogEntry // 需要添加的 log 单元
}

func (a AppendEntriesArgs) String() string {
	return fmt.Sprintf("appendEntriesArgs{R%d:T%d, PrevLogIndex:%d, PrevLogTerm:%d, LeaderCommit:%d, entries:%v}",
		a.LeaderID, a.Term, a.PrevLogIndex, a.PrevLogTerm, a.LeaderCommit, a.Entries)
}

func (rf *Raft) newAppendEntriesArgs(server int) AppendEntriesArgs {
	prevLogIndex := rf.nextIndex[server] - 1
	baseIndex := rf.getBaseIndex()
	return AppendEntriesArgs{
		Term:         rf.currentTerm,
		LeaderID:     rf.me,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  rf.logs[prevLogIndex-baseIndex].LogTerm,
		Entries:      rf.logs[prevLogIndex+1-baseIndex:],
		LeaderCommit: rf.commitIndex,
	}
}

// AppendEntriesReply 是 flower 回复 leader 的内容
type AppendEntriesReply struct {
	Term      int  // 回复者的 term
	Success   bool // 返回 true，如果被调用的 rf.logs 真的 append 了 entries
	NextIndex int  // 下一次发送的 AppendEntriesArgs.Entries[0] 在 Leader.logs 中的索引号
}

func (r AppendEntriesReply) String() string {
	return fmt.Sprintf("appendEntriesReply{T%d, Success:%t, NextIndex:%d}",
		r.Term, r.Success, r.NextIndex)
}

func (rf *Raft) sendAppendEntries(server int, args AppendEntriesArgs, reply *AppendEntriesReply) bool {
	return rf.peers[server].Call("Raft.AppendEntries", args, reply)
}

// 广播 AppendEntries 有两个作用
// 1. heart beat: 阻止其他 server 发起选举
// 2. 同步 log 到其他 server
func (rf *Raft) broadcastAppendEntries() {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	lastIndex := rf.getLastIndex()
	baseIndex := rf.getBaseIndex()

	newCommitIndex := 0
	// 统计 leader 的此 term 的已复制 log 数量，超过半数，就可以 commit 了
	for idx := rf.commitIndex + 1; idx <= lastIndex; idx++ {
		count := 1 // 1 是 rf 自己的一票
		for id := range rf.peers {
			if id != rf.me &&
				rf.matchIndex[id] >= idx &&
				rf.logs[idx-baseIndex].LogTerm == rf.currentTerm {
				count++
			}
		}
		if 2*count > len(rf.peers) {
			newCommitIndex = idx
		}
	}
	if newCommitIndex > rf.commitIndex {
		rf.commitIndex = newCommitIndex
		rf.chanCommit <- struct{}{}
		DPrintf("%s COMMITTED %s", rf, rf.details())
	}

	for id := range rf.peers {
		if id != rf.me && rf.isLeader() {
			args := rf.newAppendEntriesArgs(id)
			go rf.endAppendEntriesAndDealReply(id, args)
		}
	}
}

func (rf *Raft) endAppendEntriesAndDealReply(id int, args AppendEntriesArgs) {
	var reply AppendEntriesReply

	DPrintf("%s AppendEntries to R%d with %s", rf, id, args)

	ok := rf.sendAppendEntries(id, args, &reply)
	if !ok {
		return
	}

	rf.mu.Lock()
	defer rf.mu.Unlock()

	if reply.Term > rf.currentTerm {
		rf.currentTerm = reply.Term
		rf.state = FOLLOWER
		rf.votedFor = NOBODY
		rf.persist()
		return
	}

	if rf.currentTerm != args.Term {
		// term 已经改变
		return
	}

	if !reply.Success {
		rf.nextIndex[id] = reply.NextIndex
		return
	}

	if len(args.Entries) == 0 {
		// 纯 heartBeat 就无需进一步处理了
		return
	}

	lastArgsLogIndex := args.Entries[len(args.Entries)-1].LogIndex
	rf.matchIndex[id] = lastArgsLogIndex
	rf.nextIndex[id] = lastArgsLogIndex + 1
}

// AppendEntries 会处理收到 AppendEntries RPC
func (rf *Raft) AppendEntries(args AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	reply.Success = false

	// 1. Replay false at once if term < currentTerm
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		DPrintf("%s rejected %s", rf, args)
		return
	}

	defer rf.persist()

	rf.chanHeartBeat <- struct{}{}

	DPrintf("%s 收到了真实有效的信号 %s", rf, args)

	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.state = FOLLOWER
		rf.votedFor = NOBODY
	}

	reply.Term = args.Term

	if args.PrevLogIndex > rf.getLastIndex() {
		reply.NextIndex = rf.getLastIndex() + 1
		return
	}

	baseIndex := rf.getBaseIndex()

	if args.PrevLogIndex > baseIndex {
		term := rf.logs[args.PrevLogIndex-baseIndex].LogTerm
		if args.PrevLogTerm != term {
			for i := args.PrevLogIndex - 1; i >= baseIndex; i-- {
				if rf.logs[i-baseIndex].LogTerm != term {
					reply.NextIndex = i + 1
					break
				}
			}
			return
		}
	}

	if args.PrevLogIndex >= baseIndex {
		rf.logs = rf.logs[:args.PrevLogIndex+1-baseIndex]
		rf.logs = append(rf.logs, args.Entries...)
		reply.Success = true
		reply.NextIndex = rf.getLastIndex() + 1
	}

	// 5. if leadercommit > commitIndex, set commitIndex = min(leaderCommit, index of last new entry)
	if args.LeaderCommit > rf.commitIndex {
		rf.commitIndex = min(args.LeaderCommit, rf.getLastIndex())
		rf.chanCommit <- struct{}{}
		DPrintf("%s COMMITTED %s", rf, rf.details())
	}

}
