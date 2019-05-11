package cli

import (
	"fmt"
	"github.com/azd1997/golang-blockchain/blockchain"
)

func (cli *CommandLine) reindexUTXO(nodeID string) {

	//从数据库中获取最新区块，返回当前区块链对象
	chain := blockchain.ContinueBlockChain(nodeID)
	defer chain.Db.Close()

	//构建UTXOSet对象，调用其reindex方法
	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()

	//统计UTXOSet中交易数
	count := UTXOSet.CountTransaction()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}
