package raft

import "testing"

func Test_state_String(t *testing.T) {
	tests := []struct {
		name string
		s    state
		want string
	}{

		{
			"Follower",
			FOLLOWER,
			"Follower",
		},

		{
			"Candidate",
			CANDIDATE,
			"Candidate",
		},

		{
			"Leader",
			LEADER,
			"Leader",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("state.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
