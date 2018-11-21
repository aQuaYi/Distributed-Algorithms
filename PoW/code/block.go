package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// Block represents a block in the blockchain
type Block struct {
	Timestamp     int64          // 创建此区块的时间
	Transactions  []*Transaction // 所有新生成的交易
	PrevBlockHash []byte         // 上一个区块的哈希值，即父哈希
	Hash          []byte         // 当前区块的哈希值
	Nonce         int            // 计算目标哈希值所需的 "计数器"
	Height        int            // 区块在区块链中索引值，创世区块为 0
}

// NewBlock creates and returns Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
		Height:        height,
	}
	pow := NewProofOfWork(block)

	nonce, hash := pow.mining()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(coinbase *Transaction) *Block {
	// 创世区块也需要包含 coinbase 交易
	return NewBlock([]*Transaction{coinbase}, []byte{}, 0)
}

// HashTransactions returns a hash of the transactions in the block
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	// 顺序会影响最后的返回结果
	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)

	return mTree.RootNode.Data
}

// Serialize serializes the block
// 序列化的目的是为了把 block 保存到 k-v 数据库中
func (b *Block) Serialize() []byte {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

// DeserializeBlock deserializes a block
// 反向序列化是为了从数据库中取出 block
func DeserializeBlock(date []byte) *Block {
	r := bytes.NewReader(date)
	dec := gob.NewDecoder(r)

	var block Block

	err := dec.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
