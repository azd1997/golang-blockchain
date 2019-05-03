package cli

import (
	"fmt"
	"github.com/azd1997/golang-blockchain/blockchain"
	"github.com/azd1997/golang-blockchain/wallet"
	"log"
)

/*转账*/
func (cli *CommandLine) send(from, to string, amount int) {

	if !wallet.ValidateAddress(to) {
		log.Panic("Address is not valid")
	}

	if !wallet.ValidateAddress(from) {
		log.Panic("Address is not valid")
	}

	chain := blockchain.ContinueBlockChain(from)
	defer chain.Db.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}
