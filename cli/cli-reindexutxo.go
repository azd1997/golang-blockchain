package cli

import (
	"fmt"
	"github.com/azd1997/golang-blockchain/blockchain"
)

func (cli *CommandLine) reindexUTXO() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Db.Close()

	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransaction()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}
