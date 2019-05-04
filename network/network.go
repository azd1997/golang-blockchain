package network

import (
	"fmt"

	"github.com/azd1997/golang-blockchain/blockchain"
	"gopkg.in/vrecan/death.v3"
)

const (
	protocal      = "tcp"
	version       = 1
	commandLength = 12
)

var (
	nodeAddress    string
	mineAddress    string
	KnownNodes     = []string{"localhost:3000"}
	blockInTransit = [][]byte{}
	memoryPool     = make(map[string]blockchain.Transaction)
)

type Addr struct {
	AddrList []string
}

type Block struct {
	AddrFrom string
	Block    []byte
}

type GetBlocks struct {
	AddrFrom string
}

type Inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

type Tx struct {
	AddrFrom    string
	Transaction []byte
}

type Version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

//流程
//1.create Blockchain
//2.Wallet connects and Download Blockchain
//3.Miner connects and download Blockchain
//4.Wallets creates tx
//5.miner gets tx to memory pool
//6.enough tx -> mine block
//7.block sent to central node	//
//8.wallet syncs and verifies

func CmdToBytes(cmd string) []byte {
	var bytes [commandLength]byte

	for i, c := range cmd {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func BytesToCmd(bytes []byte) string {
	var cmd []byte

	for _, b := range bytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}

	return fmt.Sprintf("%s", cmd)
}

func Ex
