package cli

import (
	"fmt"
	"github.com/azd1997/golang-blockchain/wallet"
)

func (cli *CommandLine) listAddresses() {

	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddress()

	for _, address := range addresses {
		fmt.Println(address)
	}

}
