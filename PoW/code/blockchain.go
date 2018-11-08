package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const (
	dbFormat            = "blockchain_%s.db"
	blocksBucket        = "blocks"
	lastBlockHash       = "l"
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

// Blockchain implements interactions with a DB
type Blockchain struct {
	db  *bolt.DB // 存储全部区块数据的 key/value 数据库
	tip []byte   // 最新的区块的哈希值
}

// CreateBlockchain creates a new blockchain DB
// TODO: 把 nodeID 改成 int 类型
func CreateBlockchain(address, nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFormat, nodeID)
	if exist(dbFile) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte

	// cbtx = coinbase transaction
	cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
	genesis := NewGenesisBlock(cbtx)

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panicf("创建 %s 是出错：%s", dbFile, err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panicf("在 %s 中创建 %s 时，失败：%s", dbFile, blocksBucket, err)
		}
		// 把创世区块放入 dbFile 的 blockBucket 中
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		// 把创世区块的 hash 值，放入 dbFile 的 blockBucket 的 lastBlockHash key 中
		err = b.Put([]byte(lastBlockHash), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panicf("创建仅包含创世区块的区块链数据库文件 %s 时，出错：%s", dbFile, err)
	}

	// 区块链就是一串珠子，和，把珠子串起来的线头
	bc := Blockchain{tip: tip, db: db}

	return &bc
}

// NewBlockchain creates a new Blockchain with genesis Block
// CreateBlockchain 用于从无到有地创建区块链数据库
// NewBlockChain    用于打开现有的区块链数据库
func NewBlockchain(nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFormat, nodeID)
	if exist(dbFile) == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panicf("打开 %s 时，失败：%s", dbFile, err)
	}

	var tip []byte // 取 tip 尖端的含义，意为指向最新的区块

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte(lastBlockHash))
		return nil
	})
	if err != nil {
		log.Panicf("在 %s 中获取 lastBlockHash 的值时，失败：%s", dbFile, err)
	}

	bc := Blockchain{tip: tip, db: db}

	return &bc
}

// AddBlock saves the block into the blockchain
func (bc *Blockchain) AddBlock(block *Block) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		blockInDb := b.Get(block.Hash)
		if blockInDb != nil {
			return nil
		}

		blockData := block.Serialize()
		err := b.Put(block.Hash, blockData)
		if err != nil {
			log.Panicf("往 %v 的区块链数据库中存入 %v 时，失败：%s", bc, block, err) // TODO: 为 bc 和 block 添加输出格式，把 %v 改成 %s
		}

		lastHash := b.Get([]byte(lastBlockHash))
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)

		if block.Height > lastBlock.Height {
			err = b.Put([]byte(lastBlockHash), block.Hash)
			if err != nil {
				log.Panic(err)
			}
			bc.tip = block.Hash
		}

		return nil
	})

	if err != nil {
		log.Panicf("把 %v 存入区块链时，失败：%s", block, err)
	}

}

// FindTransaction finds a transaction by its ID
// 具体过程就是每个区块挨个去查验
// TODO: 返回 Transaction 的指针
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()
	//
	for bci.HasNext() {
		block := bci.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.Hash, ID) == 0 {
				return *tx, nil
			}
		}
	}
	//
	return Transaction{}, errors.New("Transaction is not found")
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
// FIXME: 弄清楚这个方法的内容
func (bc *Blockchain) FindUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for bci.HasNext() {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.Hash)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.RefTxHash)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.OutIndex)
				}
			}
		}

	}

	return UTXO
}

// GetBestHeight returns the height of the latest block
// 因为重新运行 node 后，生成新的区块，需要用到 best height
func (bc *Blockchain) GetBestHeight() int {
	var lastBlock Block
	//
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte(lastBlockHash))
		data := b.Get(lastHash)
		lastBlock = *DeserializeBlock(data)
		return nil
	})
	//
	if err != nil {
		log.Panicf("获取最新的区块链失败：%s", err)
	}
	//
	return lastBlock.Height
}

// GetBlock finds a block by its hash and returns it
// 按照区块的 hash 值，获取区块的内容
func (bc *Blockchain) GetBlock(blockHash []byte) (Block, error) {
	var block Block
	//
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		data := b.Get(blockHash)
		if data == nil {
			return errors.New("block is not found")
		}
		block = *DeserializeBlock(data)
		return nil
	})
	//
	if err != nil {
		return block, err
	}
	//
	return block, nil
}

// GetBlockHashes returns a list of hashes of all the blocks in the chain
// 越老的区块，索引值越大
func (bc *Blockchain) GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Iterator()
	for bci.HasNext() {
		block := bci.Next()
		blocks = append(blocks, block.Hash)
	}
	return blocks
}

// MineBlock mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block {
	// 挖矿前，先验证每个交易是否可行
	for _, tx := range transactions {
		if bc.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}
	//
	var lastHash []byte
	var lastHeight int
	//
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte(lastBlockHash))
		data := b.Get(lastHash)
		block := DeserializeBlock(data)
		lastHeight = block.Height
		return nil
	})
	if err != nil {
		log.Panicf("获取 lastHash 和 lastHeight 时，出错：%s", err)
	}
	// NewBlock 中包含了挖矿的过程
	block := NewBlock(transactions, lastHash, lastHeight+1)
	// 把新挖出来的区块，放入数据库
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(block.Hash, block.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte(lastBlockHash), block.Hash)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	// 更新 bc.tip
	bc.tip = block.Hash
	return block
}

// SignTransaction signs inputs of a Transaction
// FIXME: 弄清楚这个方法的内容
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)
	//
	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.RefTxHash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Hash)] = prevTX
	}
	//
	tx.Sign(privKey, prevTXs)
}

// VerifyTransaction verifies transaction input signatures
// FIXME: 弄清楚这个方法的内容
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	//
	prevTXs := make(map[string]Transaction)
	//
	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.RefTxHash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Hash)] = prevTX
	}
	//
	return tx.Verify(prevTXs)
}
