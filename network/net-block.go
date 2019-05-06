package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/azd1997/golang-blockchain/blockchain"
	"github.com/azd1997/golang-blockchain/utils"
)

type Block struct {
	AddrFrom string
	Block    []byte
}

//处理接收到区块时
func HandleBlock(request []byte, chain *blockchain.BlockChain) {
	//获取request内容
	var buff bytes.Buffer
	var payload Block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	utils.Handle(err)

	//将接收到的区块添加到区块链中
	blockData := payload.Block
	block := blockchain.Deserialize(blockData)

	fmt.Println("Received a new block!")
	chain.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	//若blockInTransit中还有内容，那么据此继续向对方请求区块数据
	//这表示只要blockInTransit非空，就会不断请求，对方不断返回区块，自己不断处理区块
	if len(blockInTransit) > 0 {
		blockHash := blockInTransit[0]
		SendGetData(payload.AddrFrom, "block", blockHash)

		blockInTransit = blockInTransit[1:]
	} else {
		//更新未花费输出集
		UTXOSet := blockchain.UTXOSet{chain}
		UTXOSet.Reindex()
	}
}

/*向某地址发送区块*/
func SendBlock(addr string, b *blockchain.Block) {
	data := Block{nodeAddress, b.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("block"), payload...)

	SendData(addr, request)
}
