package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Serialize_And_Deserialize(t *testing.T) {
	ast := assert.New(t)
	//
	block := &Block{
		Timestamp: 1024,
		Transactions: []*Transaction{
			&Transaction{
				Hash: []byte("This is Transactions' ID"),
				Vin:  []TXInput{TXInput{}},
				Vout: []TXOutput{TXOutput{}},
			},
		},
		PrevBlockHash: []byte("abcdefghijklmn"),
		Hash:          []byte("ABCDEFGHIJKLMN"),
		Nonce:         1999,
		Height:        8102,
	}
	//
	data := block.Serialize()
	newBlock := DeserializeBlock(data)
	//
	ast.False(block == newBlock, "两个区块不能指向同一个地址")
	ast.Equal(block, newBlock)
}
