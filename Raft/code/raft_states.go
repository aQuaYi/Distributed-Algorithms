package raft

type fsmState int

// 规定了 server 所需的 3 种状态
const (
	LEADER fsmState = iota
	CANDIDATE
	FOLLOWER
)

func (s fsmState) String() string {
	switch s {
	case LEADER:
		return "Leader"
	case CANDIDATE:
		return "Candidate"
	case FOLLOWER:
		return "Follower"
	default:
		panic("出现了第4种 server state")
	}
}
