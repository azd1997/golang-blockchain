package network

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/azd1997/golang-blockchain/blockchain"
	"github.com/azd1997/golang-blockchain/utils"
	"io"
	"net"
)

type GetData struct {
	AddrFrom string
	Type     string
	ID       []byte
}

//处理获取数据请求
func HandleGetData(request []byte, chain *blockchain.BlockChain) {

	//获取request中的内容
	var buff bytes.Buffer
	var payload GetData

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	utils.Handle(err)

	//getdata有获取区块和获取交易两种情况

	if payload.Type == "block" {
		block, err := chain.GetBlock([]byte(payload.ID))
		utils.Handle(err)

		//向给自己发请求的节点发送单个区块
		SendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := memoryPool[txID]

		SendTx(payload.AddrFrom, &tx)
	}

}

/*发送获取数据的请求*/
func SendGetData(address, kind string, id []byte) {
	payload := GobEncode(GetData{nodeAddress, kind, id})
	request := append(CmdToBytes("getdata"), payload...)

	SendData(address, request)
}

func SendData(addr string, data []byte) {
	//向addr发起tcp连接
	conn, err := net.Dial(protocol, addr)

	//连接不可用，则更新已知节点集
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range KnownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		KnownNodes = updatedNodes

		return
	}

	defer conn.Close()

	//将data []byte复制一份通过conn发给对方
	_, err = io.Copy(conn, bytes.NewReader(data))
	utils.Handle(err)
}
