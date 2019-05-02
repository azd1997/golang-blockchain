package cli

import (
	"fmt"
	"github.com/azd1997/golang-blockchain/wallet"
)

func (cli *CommandLine) createWallet() {

	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Printf("New address is: %s\n", address)

}
