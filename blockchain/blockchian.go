package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	//Blocks []*Block
	LastHash []byte
	Db       *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Db          *badger.DB
}

/*创建带有创世区块的区块链*/
func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		//如果数据库中存在区块链...
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("区块链不存在...")
			genesis := Genesis()
			fmt.Println("创世区块已被开发")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)

			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else { //若数据库中已有区块链...
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.Value()
			return err
		}
	})
	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

/*向区块链中添加区块*/
func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte
	//获取最新区块哈希及其区块信息
	err := chain.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()

		return err
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Db.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
}

/*返回区块链迭代器对象*/
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Db}

	return iter
}

/*迭代，用以遍历区块信息*/
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err := item.Value()
		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}
