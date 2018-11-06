package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

// 目标哈希值的前 targetBits 位必须是 0
// 代表了挖矿的难度，数值越大越难
const targetBits = 16

// ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	// Lsh 是把 target 左移 256-targetBits 位
	target.Lsh(target, uint(256-targetBits))
	// 目标哈希值，需要比此时的 target 小

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			Int64ToHex(pow.block.Timestamp),
			Int64ToHex(int64(targetBits)),
			Int64ToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// mining performs a proof-of-work
// 这就是俗称的挖矿过程，找到合适的 nonce 使得区块的哈希值 < target
func (pow *ProofOfWork) mining() (int, []byte) {
	// hashInt 是把 hash 按照大端无符号的方式，解释成整数
	var hashInt big.Int
	var hash [32]byte
	nonce := 0 // once number 简称 nonce https://zh.wikipedia.org/zh-cn/Nonce

	fmt.Printf("Mining a new block")

	for nonce < maxNonce {
		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)

		// TODO: 清洁代码时，可以删除此处
		if math.Remainder(float64(nonce), 100000) == 0 {
			fmt.Printf("\r%x", hash)
		}

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			// 当哈希所代表的数值，小于 pow.target 的时候，
			// 说明找到了想要的 nonce 及其哈希值
			break
		} else {
			// 没有找到就顺着继续找
			nonce++
		}
	}

	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate validates block's PoW
// 验证 block，hash，nonce 和 target 是否匹配
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
