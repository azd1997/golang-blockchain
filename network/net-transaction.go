package network

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/azd1997/golang-blockchain/blockchain"
	"github.com/azd1997/golang-blockchain/utils"
)

type Tx struct {
	AddrFrom    string
	Transaction []byte
}

//处理收到一笔交易信息
func HandleTx(request []byte, chain *blockchain.BlockChain) {
	//将request的command后内容解码写入payload(Tx)
	var buff bytes.Buffer
	var payload Tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	utils.Handle(err)

	//将该交易存入内存池
	txData := payload.Transaction
	tx := blockchain.DeserializeTransaction(txData)
	memoryPool[hex.EncodeToString(tx.ID)] = tx

	fmt.Printf("%s, %d\n", nodeAddress, len(memoryPool))

	//若本地节点是已知节点集第一个（已知节点集第一个初始化为本地节点）
	//则向已知节点集中除自己和发交易给自己的节点外的所有节点发送存证，告诉大家我收到了这个交易
	if nodeAddress == KnownNodes[0] {
		for _, node := range KnownNodes {
			if node != nodeAddress && node != payload.AddrFrom {
				SendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else { //内存池有至少两个交易且挖矿节点地址被设定，进行MineTx
		if len(memoryPool) >= 2 && len(mineAddress) > 0 {
			MineTx(chain)
		}
	}
}

/*挖矿，将本地交易打包发布*/
func MineTx(chain *blockchain.BlockChain) {
	var txs []*blockchain.Transaction

	//从内存池（记忆池）中遍历交易，交易符合规则的加入待出块交易集合
	for id := range memoryPool {
		fmt.Printf("tx: %s\n", memoryPool[id].ID)
		tx := memoryPool[id]
		if chain.VerifyTransaction(&tx) {
			txs = append(txs, &tx)
		}
	}

	//若待出块交易集合长度为0，说明内存池所有交易均无效
	if len(txs) == 0 {
		fmt.Println("All Transaction are invalid")
		return
	}

	//挖矿者在出块时自行创建Coinbase交易，数据域可以自行指定，若为空则随机字符串
	cbTx := blockchain.CoinbaseTx(mineAddress, "")
	txs = append(txs, cbTx)

	//新区块
	newBlock := chain.MineBlock(txs)
	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()

	fmt.Println("New Block Mined")

	//在出块之后，删除内存池中已放入待出块交易集合的交易
	for _, tx := range txs {
		txID := hex.EncodeToString(tx.ID)
		delete(memoryPool, txID)
	}

	//向已知节点集中除了本机外的节点发送出块存证（告诉别人我挖出矿了）
	for _, node := range KnownNodes {
		if node != nodeAddress {
			SendInv(node, "block", [][]byte{newBlock.Hash})
		}
	}

	//递归调用，如果内存池不为空，继续MineTx
	if len(memoryPool) > 0 {
		MineTx(chain)
	}

	//注意：比如说比特币，设置出块时间约15分钟，则区块内交易量是挖矿者打包区块之前收了多少算多少
}

/*向某一地址发送交易数据*/
func SendTx(addr string, tx *blockchain.Transaction) {
	data := Tx{nodeAddress, tx.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("tx"), payload...)

	SendData(addr, request)
}
