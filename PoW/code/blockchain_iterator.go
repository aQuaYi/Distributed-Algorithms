package main

import (
	"log"

	"github.com/boltdb/bolt"
)

// BlockchainIterator is used to iterate over blockchain blocks
// 迭代器用于从新往旧，依次访问区块链中的每一个区块。
type BlockchainIterator struct {
	nextHash []byte   // 下一次调用 next 时，返回区块的 hash 值
	db       *bolt.DB // 存储区块链的数据库
}

// Next returns next block starting from the tip
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		data := b.Get(i.nextHash)
		block = DeserializeBlock(data)

		return nil
	})
	if err != nil {
		log.Panicf("迭代时，无法取出 %x : %s", i.nextHash, err)
	}

	i.nextHash = block.PrevBlockHash

	return block
}

// HasNext 告诉你能否继续迭代
func (i *BlockchainIterator) HasNext() bool {
	// 生成创世区块时，其 prevBlockHash 参数设置为了 []byte，而非 nil
	return len(i.nextHash) != 0
}

// Iterator returns a BlockchainIterator
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}
