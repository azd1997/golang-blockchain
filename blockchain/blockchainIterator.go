package blockchain

import (
	"github.com/azd1997/golang-blockchain/utils"
	"github.com/dgraph-io/badger"
)

type BCIterator struct {
	CurrentHash []byte
	Db          *badger.DB
}

/*迭代器对象的Next方法，用以返回当前区块，并更新BCIterator对象至对应前一个区块*/
func (iter *BCIterator) Next() *Block {
	var block *Block

	//从数据库取出当前区块的序列化字节，反序列化
	err := iter.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		utils.Handle(err)
		encodedBlock, err := item.Value()
		block = Deserialize(encodedBlock)

		return err
	})
	utils.Handle(err)

	//更新BCIterator对象
	iter.CurrentHash = block.PrevHash

	return block
}
