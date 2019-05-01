package blockchain

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct {
	//blockchain *blockchain.BlockChain
}

/*打印命令行工具用法*/
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS - get the balance for ADDRESS")
	fmt.Println(" createblockchain -address ADDRESS - creates a blockchain and sends genesis reward to ADDRESS")
	//fmt.Println(" add -block BLOCK_DATA - add a block to chain")
	fmt.Println(" printchain - Prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send amount of coins")
}

/*检查命令行输入参数是否至少有两个*/
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
		//os.Exit(1)
	}
}

/*命令行打印区块信息*/
func (cli *CommandLine) printChain() {
	chain := ContinueBlockChain("")
	defer chain.Db.Close()

	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("TransactionsHash: %x\n", block.HashTransactions())
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := NewProof(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 { //创世区块PrevHash设为0
			break
		}
	}
}

/*创建区块链，其创世区块coinbase交易地址给定*/
func (cli *CommandLine) createBlockChain(address string) {
	chain := InitBlockChain(address)
	chain.Db.Close()
	fmt.Println("Finished!")
}

/*获取账户余额*/
func (cli *CommandLine) getBalance(address string) {
	chain := ContinueBlockChain(address)
	defer chain.Db.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

/*转账*/
func (cli *CommandLine) send(from, to string, amount int) {
	chain := ContinueBlockChain(from)
	defer chain.Db.Close()

	tx := NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*Transaction{tx})
	fmt.Println("Success!")
}

/*运行命令行程序*/
func (cli *CommandLine) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for.")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to.")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		Handle(err)
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		Handle(err)
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		Handle(err)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		Handle(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}
