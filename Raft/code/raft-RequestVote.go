package raft

import "fmt"

// RequestVoteArgs 获取投票参数
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term         int // candidate's term
	CandidateID  int // candidate requesting vote
	LastLogIndex int // index of candidate's last log entry
	LastLogTerm  int // term of candidate's last log entry
}

func (a RequestVoteArgs) String() string {
	return fmt.Sprintf("voteArgs{R%d:T%d;LastLogIndex:%d;LastLogTerm:%d}",
		a.CandidateID, a.Term, a.LastLogIndex, a.LastLogTerm)
}

// RequestVoteReply is
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Term        int
	VoteGranted bool
}

func (reply RequestVoteReply) String() string {
	return fmt.Sprintf("voteReply{T%d,Granted:%t}", reply.Term, reply.VoteGranted)
}

// RequestVote is
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).

	DPrintf("%s 收到投票请求 [%s]", rf, args)

	// rf.rwmu.Lock() // TODO: 这里是否需要锁
	// defer rf.rwmu.Unlock()
	// defer rf.persist()

	// 1. replay false if term < currentTerm
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return
	}

	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.state = FOLLOWER
		rf.votedFor = NOBODY
	}

	reply.Term = rf.currentTerm

	// 2. votedFor is null or candidateId and
	//    candidate's log is at least as up-to-date as receiver's log, then grant vote
	//    If the logs have last entries with different terms, then the log with the later term is more up-to-date
	//    If the logs end with the same term, then whichever log is longer is more up-to-date

	if isValidArgs(rf, args) {
		reply.VoteGranted = true
		rf.chanGrantVote <- struct{}{}
		rf.state = FOLLOWER
		rf.votedFor = args.CandidateID
		DPrintf("%s voted for %s", rf, args)
		return
	}
	DPrintf("%s **NOT** voted for %s", rf, args)
}

func isValidArgs(rf *Raft, args *RequestVoteArgs) bool {
	term := rf.getLastTerm()
	index := rf.getLastIndex()
	return (rf.votedFor == NOBODY || rf.votedFor == args.CandidateID) &&
		isUpToDate(args, term, index)
}

func isUpToDate(args *RequestVoteArgs, term, index int) bool {
	return (args.LastLogTerm > term) ||
		(args.LastLogTerm == term && args.LastLogIndex >= index)
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
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in struts passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

func (rf *Raft) boatcastRequestVote() {
	panic("boatcastRequestVote is empty")
}
