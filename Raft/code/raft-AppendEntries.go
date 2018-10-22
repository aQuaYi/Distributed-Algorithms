package raft

import "fmt"

// AppendEntriesArgs 是添加 log 的参数
type AppendEntriesArgs struct {
	Term         int        // leader.currentTerm
	LeaderID     int        // leader.me
	PrevLogIndex int        // index of log entry immediately preceding new ones
	PrevLogTerm  int        // term of prevLogIndex entry
	LeaderCommit int        // leader.commitIndex
	Entries      []LogEntry // 需要添加的 log 单元，为空时，表示此条消息是 heartBeat
}

func (a AppendEntriesArgs) String() string {
	return fmt.Sprintf("appendEntriesArgs{R%d:T%d, PrevLogIndex:%d, PrevLogTerm:%d, LeaderCommit:%d, entries:%v}",
		a.LeaderID, a.Term, a.PrevLogIndex, a.PrevLogTerm, a.LeaderCommit, a.Entries)
}

func newAppendEntriesArgs(leader *Raft, server int) AppendEntriesArgs {
	prevLogIndex := leader.nextIndex[server] - 1
	baseIndex := leader.getBaseIndex()
	return AppendEntriesArgs{
		Term:         leader.currentTerm,
		LeaderID:     leader.me,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  leader.logs[prevLogIndex-baseIndex].LogTerm,
		Entries:      leader.logs[prevLogIndex+1-baseIndex:],
		LeaderCommit: leader.commitIndex,
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

func (rf *Raft) broadcastAppendEntries() {

	rf.mu.Lock()
	defer rf.mu.Unlock()
	N := rf.commitIndex
	last := rf.getLastIndex()
	baseIndex := rf.getBaseIndex()

	// 统计 leader 的此 term 的已复制 log 数量，超过半数，就可以 commit 了
	for i := rf.commitIndex + 1; i <= last; i++ {
		num := 1
		for j := range rf.peers {
			if j != rf.me && rf.matchIndex[j] >= i && rf.logs[i-baseIndex].LogTerm == rf.currentTerm {
				num++
			}
		}
		if 2*num > len(rf.peers) {
			N = i
		}
	}
	if N != rf.commitIndex {
		rf.commitIndex = N
		rf.chanCommit <- struct{}{}
		DPrintf("%s COMMITTED %s", rf, rf.details())
	}

	for i := range rf.peers {
		if i != rf.me && rf.isLeader() {
			args := newAppendEntriesArgs(rf, i)
			go rf.sendAppendEntriesAndDealReply(i, args)
		}
	}
}

func (rf *Raft) sendAppendEntriesAndDealReply(i int, args AppendEntriesArgs) {
	var reply AppendEntriesReply

	DPrintf("%s AppendEntries to R%d with %s", rf, i, args)

	ok := rf.sendAppendEntries(i, args, &reply)
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
		rf.nextIndex[i] = reply.NextIndex
		return
	}

	if len(args.Entries) == 0 {
		return
	}

	lastArgsLogIndex := args.Entries[len(args.Entries)-1].LogIndex
	rf.matchIndex[i] = lastArgsLogIndex
	rf.nextIndex[i] = lastArgsLogIndex + 1
}

// AppendEntries 会处理收到 AppendEntries RPC
func (rf *Raft) AppendEntries(args AppendEntriesArgs, reply *AppendEntriesReply) {

	// Your code here.

	rf.mu.Lock()
	defer rf.mu.Unlock()
	defer rf.persist()

	reply.Success = false

	// 1. Replay false at once if term < currentTerm
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.NextIndex = rf.getLastIndex() + 1
		DPrintf("%s rejected %s", rf, args)
		return
	}

	rf.chanHeartbeat <- struct{}{}

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
