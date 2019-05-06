package cli

import (
	"fmt"
	"github.com/azd1997/golang-blockchain/blockchain"
	"github.com/azd1997/golang-blockchain/network"
	"github.com/azd1997/golang-blockchain/utils"
	"github.com/azd1997/golang-blockchain/wallet"
	"log"
)

/*转账*/
func (cli *CommandLine) send(from, to, nodeID string, amount int, mineNow bool) {

	if !wallet.ValidateAddress(to) {
		log.Panic("Address is not valid")
	}

	if !wallet.ValidateAddress(from) {
		log.Panic("Address is not valid")
	}

	chain := blockchain.ContinueBlockChain(nodeID)
	UTXOSet := blockchain.UTXOSet{chain}
	defer chain.Db.Close()

	//从钱包文件读取内容，创建钱包对象
	wallets, err := wallet.CreateWallets(nodeID)
	utils.Handle(err)
	fromWallet := wallets.GetWallet(from)

	//创建新交易
	tx := blockchain.NewTransaction(&fromWallet, to, amount, &UTXOSet)

	if mineNow { // mineNow == true
		//挖矿
		cbTx := blockchain.CoinbaseTx(from, "")
		txs := []*blockchain.Transaction{cbTx, tx}
		block := chain.MineBlock(txs)
		UTXOSet.Update(block)
	} else { // mineNow == false
		//向本地节点发送，用以调试
		network.SendTx(network.KnownNodes[0], tx)
		fmt.Println("Send tx")
	}

	fmt.Println("Success!")
}
