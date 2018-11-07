package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

// Int64ToHex converts an int64 to a byte array
// Int64ToHex 把 num 的二进制编码，按照大端法的顺序存储在了 []byte
// []byte 中的每个元素会存放 一个字节的内容
// 具体用法，请参考其单元测试
func Int64ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// ReverseBytes reverses a byte array
// 反转 data 的顺序
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// exist 用于判断文件是否存在
func exist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

// FindUnspentTransactions TODO: 删除此处内容
// func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
// 	var unspentTXs []Transaction
// 	spentTXOs := make(map[string][]int)
// 	bci := bc.Iterator()

// 	for {
// 		block := bci.Next()

// 		for _, tx := range block.Transactions {
// 			txID := hex.EncodeToString(tx.ID)

// 		Outputs:
// 			for outIdx, out := range tx.Vout {
// 				// Was the output spent?
// 				if spentTXOs[txID] != nil {
// 					for _, spentOut := range spentTXOs[txID] {
// 						if spentOut == outIdx {
// 							continue Outputs
// 						}
// 					}
// 				}

// 				if out.CanBeUnlockedWith(address) {
// 					unspentTXs = append(unspentTXs, *tx)
// 				}
// 			}

// 			if tx.IsCoinbase() == false {
// 				for _, in := range tx.Vin {
// 					if in.CanUnlockOutputWith(address) {
// 						inTxID := hex.EncodeToString(in.Txid)
// 						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
// 					}
// 				}
// 			}
// 		}

// 		if len(block.PrevBlockHash) == 0 {
// 			break
// 		}
// 	}

// 	return unspentTXs
// }
