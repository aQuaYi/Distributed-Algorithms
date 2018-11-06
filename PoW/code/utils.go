package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

// Int64ToHex converts an int64 to a byte array
// Int64ToHex 把 num 的 16 进制数值，按照大端法转换成了 []bytes
// 例如， Int64ToHex(0x18) = []byte("0000000000000018")
func Int64ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
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
