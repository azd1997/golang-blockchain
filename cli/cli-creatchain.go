package cli

import (
	"fmt"
	"github.com/azd1997/golang-blockchain/blockchain"
)

/*创建区块链，其创世区块coinbase交易地址给定*/
func (cli *CommandLine) createBlockChain(address string) {
	chain := blockchain.InitBlockChain(address)
	chain.Db.Close()
	fmt.Println("Finished!")
}
