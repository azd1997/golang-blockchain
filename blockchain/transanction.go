package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/azd1997/golang-blockchain/utils"
	"log"
)

type Transaction struct {
	ID        []byte //即交易哈西
	TXInputs  []TxInput
	TXOutputs []TxOutput
}

/*对交易进行哈希，设置其交易ID*/
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	utils.Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

/*判断交易是否是Coinbase交易*/
func (tx *Transaction) IsCoinbase() bool {
	//交易输入只有1个且该交易输入的ID也就是哈希为0（ID为空，所以哈希为0）
	return len(tx.TXInputs) == 1 && len(tx.TXInputs[0].ID) == 0
}

/*出块奖励交易*/
func CoinbaseTx(to, data string) *Transaction {

	//若Coinbase交易未指定Data内容，则默认为下方内容
	//对于挖出区块的矿工而言，可以在Coinbase交易的data域填想填的东西
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	//Coinbase来源交易不存在，所以填空字节，其来源交易占来源输出序号也不存在，这里以-1表示
	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{100, to}

	//Coinbase交易只有一笔输入一笔输出，其交易ID或者说哈希需要进行哈希才能得到
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return &tx
}

/*产生一笔新交易*/
func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput   //当前交易的输入
	var outputs []TxOutput //当前交易输出

	//用户进行转账时，需指定转账者A，被转账者B以及转账金额S
	//主要需要考虑以下几点：
	//1.找到A的所有UTXO，计算A的余额，检查A余额是否足够
	//2.检查A余额

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)
	//注意返回的acc是有可能小于amount的！！！

	//余额不够，报错
	if acc < amount {
		log.Panic("Error: not enough funds...")
	}

	for txid, outs := range validOutputs { //txid为键，outs为值
		txID, err := hex.DecodeString(txid) //由字符串解码回十六进制字节数组
		utils.Handle(err)

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	//增加找零的交易输出
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
