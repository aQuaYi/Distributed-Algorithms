package main

import "bytes"

// TXInput represents a transaction input
type TXInput struct {
	RefTxHash []byte
	OutIndex  int
	// RefTxHash + Vout 表明 input 所 “引用” 的是 Hash 值为 RefTxHash 的 Vout[OutIndex]
	Signature []byte
	PubKey    []byte
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
