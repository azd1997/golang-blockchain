package cli

import (
	"fmt"
	"github.com/azd1997/golang-blockchain/blockchain"
	"github.com/azd1997/golang-blockchain/wallet"
	"log"
)

/*创建区块链，其创世区块coinbase交易地址给定*/
func (cli *CommandLine) createBlockChain(address, nodeID string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not valid")
	}

	chain := blockchain.InitBlockChain(address, nodeID)
	defer chain.Db.Close()

	//UTXO
	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()

	fmt.Println("Finished!")
}
