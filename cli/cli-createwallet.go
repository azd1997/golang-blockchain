package cli

import (
	"fmt"
	"github.com/azd1997/golang-blockchain/wallet"
)

func (cli *CommandLine) createWallet() {

	//创造钱包集对象
	wallets, _ := wallet.CreateWallets()
	//向钱包集新增一个钱包并保存到文件去
	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Printf("New address is: %s\n", address)

}
