package main

import (
	"github.com/azd1997/golang-blockchain/cli"
	"os"
)



func main() {

	defer os.Exit(0)
	Cli := cli.CommandLine{}
	Cli.Run()

}
