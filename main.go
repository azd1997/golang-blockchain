package main

import (
	//"../golang-blockchain/blockchain"  在不采用GO MOD或者GOPATH时可以使用相对路径导包
	"github.com/azd1997/golang-blockchain/blockchain"
	"os"
)



func main() {

	defer os.Exit(0)
	cli := blockchain.CommandLine{}
	cli.Run()

}
