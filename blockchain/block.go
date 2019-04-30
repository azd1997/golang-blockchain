package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

//在实现POW之后弃用
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	//block.DeriveHash()	//弃用
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

/*对区块进行序列化，返回字节数组*/
func (b *Block) Serialize() []byte {
	//创建编码器对象
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	//序列化区块
	err := encoder.Encode(b)

	Handle(err)

	return result.Bytes()
}

/*对序列化后的区块进行反序列化*/
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

/*err处理程序*/
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
