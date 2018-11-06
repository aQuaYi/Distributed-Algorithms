package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ProofOfWork_mining(t *testing.T) {
	ast := assert.New(t)
	//
	ast.Equal(16, targetBits, "挖矿的难度应该等于 16")
	//
	block := &Block{
		Transactions: []*Transaction{&Transaction{}},
	}
	//
	expected := 27765
	//
	pow := NewProofOfWork(block)
	actual, _ := pow.mining()
	ast.Equal(expected, actual)
}
