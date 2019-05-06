package cli

import (
	"fmt"
	"github.com/azd1997/golang-blockchain/wallet"
)

func (cli *CommandLine) listAddresses(nodeID string) {

	wallets, _ := wallet.CreateWallets(nodeID)
	addresses := wallets.GetAllAddress()

	for _, address := range addresses {
		fmt.Println(address)
	}

}
