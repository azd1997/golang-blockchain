package main

import (
	"fmt"
	"github.com/azd1997/golang-blockchain/blockchain"
)


func main() {

	chain := blockchain.InitBlockChain()

	chain.AddBlock("第二个区块")
	chain.AddBlock("第三个区块")
	chain.AddBlock("第四个区块")

	for _, block := range chain.Blocks {
		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}

}
