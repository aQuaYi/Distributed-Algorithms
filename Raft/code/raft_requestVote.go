package raft

import (
	"fmt"
)

// RequestVoteArgs 获取投票参数
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	Term         int // candidate's term
	CandidateID  int // candidate requesting vote
	LastLogIndex int // index of candidate's last log entry
	LastLogTerm  int // term of candidate's last log entry
}

func (a RequestVoteArgs) String() string {
	return fmt.Sprintf("voteArgs{R%d;Term:%d;LastLogIndex:%d;LastLogTerm:%d}",
		a.CandidateID, a.Term, a.LastLogIndex, a.LastLogTerm)
}

// RequestVoteReply 投票回复
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Term          int  // 投票人的 currentTerm
	IsVoteGranted bool // 返回 true，表示获得投票
}

func (reply RequestVoteReply) String() string {
	return fmt.Sprintf("voteReply{Term:%d,isGranted:%t}", reply.Term, reply.IsVoteGranted)
}

// RequestVote 投票工作
// example RequestVote RPC handler.
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	debugPrintf("%s  收到投票请求 [%s]", rf, args)

	// 1. replay false if term < currentTerm
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.IsVoteGranted = false
		return
	}

	// 如果 args.Term > rf.currentTerm 的话
	// 更新选民的 currentTerm
	if args.Term > rf.currentTerm {
		rf.call(discoverNewTermEvent,
			toFollowerArgs{
				term:     args.Term,
				votedFor: NOBODY,
			})
	}

	// 2. votedFor is null or candidateId and
	//    candidate's log is at least as up-to-date as receiver's log, then grant vote
	//    If the logs have last entries with different terms, then the log with the later term is more up-to-date
	//    If the logs end with the same term, then whichever log is longer is more up-to-date
	if isValidArgs(rf, args) {
		debugPrintf("%s   投票给了 < %s >", rf, args)
		reply.Term = rf.currentTerm
		reply.IsVoteGranted = true
		rf.votedFor = args.CandidateID

		// 运行到这里，可以认为接收到了合格的 rpc 信号，可以重置 election timer 了
		debugPrintf("%s  准备发送重置 election timer 信号", rf)
		rf.heartbeatChan <- struct{}{}
	} else {
		debugPrintf("%s  拒绝投票给 < %s >", rf, args)
	}

}

func isValidArgs(rf *Raft, args *RequestVoteArgs) bool {
	return (rf.votedFor == NOBODY || rf.votedFor == args.CandidateID) &&
		((args.LastLogTerm > rf.logs[len(rf.logs)-1].LogTerm) ||
			((args.LastLogTerm == rf.logs[len(rf.logs)-1].LogTerm) && args.LastLogIndex >= len(rf.logs)-1))
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ./labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in struct passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

func (rf *Raft) newRequestVoteArgs() *RequestVoteArgs {
	args := &RequestVoteArgs{
		Term:         rf.currentTerm,
		CandidateID:  rf.me,
		LastLogIndex: len(rf.logs) - 1,
		LastLogTerm:  rf.logs[len(rf.logs)-1].LogTerm,
	}
	return args
}
