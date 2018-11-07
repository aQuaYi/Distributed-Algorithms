package main

import (
	"fmt"
	"log"
)

// 从命令行创建区块链
// TODO: 精简掉这个方法
func (cli *CLI) createBlockchain(address, nodeID string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := CreateBlockchain(address, nodeID)
	defer bc.db.Close()

	UTXOSet := UTXOSet{bc}
	UTXOSet.ReIndex()

	fmt.Println("Done!")
}
