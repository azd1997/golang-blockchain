package cli

import (
	"flag"
	"fmt"
	"github.com/azd1997/golang-blockchain/utils"
	"os"
	"runtime"
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
	fmt.Println(" createwallet - Create a new Wallet")
	fmt.Println(" listaddresses - Lists the addresses in wallet file")
	fmt.Println(" reindexutxo - Rebuild the UTXO set")
}

/*检查命令行输入参数是否至少有两个*/
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
		//os.Exit(1)
	}
}

/*运行命令行程序*/
func (cli *CommandLine) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for.")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to.")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		utils.Handle(err)
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

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO()
	}

}

//调试流程
//0.删除数据库
//钱包文件中已有两个钱包：19mTU5jXmfTuE2uUx3XhRU7jZ2HoeofcCW	18v5TVC3LHsnEUC7VzvqKFVuZb4kTW1R6N
//1.先创建一个钱包
//main createwallet
//2.将第一个地址作为创世区块奖励地址,创建区块链
//main createblockchain -address "address1"
//3.打印地址1余额
//main getbalance -address "address1"
//4.打印区块链
//main printchain
//5.再创建一个钱包
//main createwallet
//6.产生两个地址之间的交易
//main send -from "address1" -to "address2" -amount 45
//7.打印地址2余额
//main getbalance -address "address2"
//8.打印区块链
//main printchain
//9.打印地址列表
//main listaddresses
