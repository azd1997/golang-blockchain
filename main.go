package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	//"../golang-blockchain/blockchain"  在不采用GO MOD或者GOPATH时可以使用相对路径导包
	"github.com/azd1997/golang-blockchain/blockchain"
	"strconv"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

/*打印命令行工具用法*/
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to chain")
	fmt.Println(" Print - Prints the blocks in the chain")
}

/*检查命令行输入参数是否至少有两个*/
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
		//os.Exit(1)
	}
}

/*命令行添加区块*/
func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

/*命令行打印区块信息*/
func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 { //创世区块PrevHash设为0
			break
		}
	}
}

/*运行命令行程序*/
func (cli *CommandLine) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	addBlockData := addBlockCmd.String("block", "", "请指定区块数据...")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}

func main() {

	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer chain.Db.Close()

	cli := CommandLine{chain}
	cli.Run()

}
