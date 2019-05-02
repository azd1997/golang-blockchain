package cli

import (
	"fmt"
	"github.com/azd1997/golang-blockchain/blockchain"
	"strconv"
)

/*命令行打印区块信息*/
func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Db.Close()

	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("TransactionsHash: %x\n", block.HashTransactions())
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 { //创世区块PrevHash设为0
			break
		}
	}
}
