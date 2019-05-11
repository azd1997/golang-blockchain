package cli

import (
	"fmt"
	"github.com/azd1997/golang-blockchain/blockchain"
	"github.com/azd1997/golang-blockchain/utils"
	"github.com/azd1997/golang-blockchain/wallet"
	"log"
)

/*获取账户余额*/
func (cli *CommandLine) getBalance(address, nodeID string) {

	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not valid")
	}

	chain := blockchain.ContinueBlockChain(nodeID)

	UTXOSet := blockchain.UTXOSet{chain}
	defer chain.Db.Close()

	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	UTXOs := UTXOSet.FindUnspentTransactions(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}
