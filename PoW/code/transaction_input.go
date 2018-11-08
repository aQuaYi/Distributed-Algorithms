package main

import "bytes"

// TXInput represents a transaction input
type TXInput struct {
	RefTxID  []byte
	OutIndex int
	// RefTxID + Vout 表明 input 所 “引用” 的是名为 RefTxID.Vout[OutIndex]
	Signature []byte
	PubKey    []byte
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
