package blockchain

import (
	"bytes"
	"encoding/gob"
	"github.com/azd1997/golang-blockchain/utils"
	"time"
)

type Block struct {
	Timestamp    int64
	Height       int
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

//方法列表
//1.func (b *Block) HashTransactions() []byte
//2.func CreateBlock(txs []*Transaction, prevHash []byte) *Block
//3.func Genesis(coinbase *Transaction) *Block
//4.func (b *Block) Serialize() []byte
//5.func Deserialize(data []byte) *Block


/*对区块中要打包的交易取哈希，并以哈希表示所有交易*/
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte //单笔交易的哈希的集合（二维字节数组）

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}
	//注意，引入Merkle后这里的txHashes其实不表示哈西了，而表示序列化交易集合
	tree := NewMerkleTree(txHashes)

	return tree.RootNode.Data
}

/*创建区块*/
func CreateBlock(txs []*Transaction, prevHash []byte, height int) *Block {
	block := &Block{time.Now().Unix(), height, []byte{}, txs, prevHash, 0}

	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

/*创建创世区块（只含有一个Coinbase交易，因为这时候只有这个账户得到钱，其他人没钱，也就不可能有其他交易）*/
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{}, 0)
}

/*对区块进行序列化，返回字节数组*/
func (b *Block) Serialize() []byte {
	//创建编码器对象
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	//序列化区块
	err := encoder.Encode(b)

	utils.Handle(err)

	return result.Bytes()
}

/*对序列化后的区块进行反序列化*/
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)

	utils.Handle(err)

	return &block
}


