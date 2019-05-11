package network

import (
	"bytes"
	"encoding/gob"
	"github.com/azd1997/golang-blockchain/blockchain"
	"github.com/azd1997/golang-blockchain/utils"
)

type GetBlocks struct {
	AddrFrom string
}

//处理获取全部区块（哈希）存证请求
func HandleGetBlocks(request []byte, chain *blockchain.BlockChain) {
	//获取request内容
	var buff bytes.Buffer
	var payload GetBlocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	utils.Handle(err)

	//向发请求的节点发送存证，说自己存了所有的区块
	blocks := chain.GetBlockHashes()
	SendInv(payload.AddrFrom, "block", blocks)
}

func SendGetBlocks(address string) {
	payload := GobEncode(GetBlocks{nodeAddress})
	request := append(CmdToBytes("getblocks"), payload...)

	SendData(address, request)
}

/*向已知节点集和中的所有节点发送GetBlocks的请求*/
func RequestBlocks() {
	for _, node := range KnownNodes {
		SendGetBlocks(node)
	}
}
