package cli

import (
	"fmt"
	"github.com/azd1997/golang-blockchain/network"
	"github.com/azd1997/golang-blockchain/wallet"
	"log"
)

func (cli *CommandLine) StartNode(nodeID, minerAddress string) {
	fmt.Printf("Starting Node %s\n", nodeID)

	if len(minerAddress) > 0 {
		if wallet.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive reward: ", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}
	network.StartServer(nodeID, minerAddress)
}
