package blockchain

import (
	"bytes"
	"encoding/hex"
	"github.com/azd1997/golang-blockchain/utils"
	"github.com/dgraph-io/badger"
	"log"
)

var (
	utxoPrefix = []byte("utxo-")
	prefixLength = len(utxoPrefix)
)

type UTXOSet struct {
	UBlockChain *BlockChain
}

//方法列表
//1.func (u UTXOSet) DeleteByPrefix(prefix []byte)
//2.func (u UTXOSet) Reindex()
//3.func (u *UTXOSet) Update(block *Block)
//4.func (u UTXOSet) CountTransaction() int
//5.func (u UTXOSet) FindUnspentTransactions(pubKeyHash []byte) []TXOutput
//6.func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int)

func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	accumulated := 0
	db := u.UBlockChain.Db

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.Value()
			utils.Handle(err)
			k = bytes.TrimPrefix(k, utxoPrefix)
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)

			for outIdx, out := range outs.TXOutputs {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOuts[txID] = append(unspentOuts[txID], outIdx)
				}
			}
		}
		return nil
	})
	utils.Handle(err)
	return accumulated, unspentOuts
}



func (u UTXOSet) FindUnspentTransactions(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput

	db := u.UBlockChain.Db

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			item := it.Item()
			v, err := item.Value()
			utils.Handle(err)
			outs := DeserializeOutputs(v)

			for _, out := range outs.TXOutputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}
		return nil
	})
	utils.Handle(err)

	return UTXOs
}

func (u UTXOSet) CountTransaction() int {
	db := u.UBlockChain.Db
	counter := 0

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			counter++
		}
		return nil
	})
	utils.Handle(err)

	return counter
}

/**/
func (u UTXOSet) Reindex() {
	db := u.UBlockChain.Db

	u.DeleteByPrefix(utxoPrefix)

	UTXO := u.UBlockChain.FindUTXO2()

	err := db.Update(func(txn *badger.Txn) error {
		for txId, outs := range UTXO {
			key, err := hex.DecodeString(txId)
			utils.Handle(err)
			key = append(utxoPrefix, key...)

			err = txn.Set(key, outs.Serialize())
			utils.Handle(err)
		}

		return nil
	})
	utils.Handle(err)
}

func (u *UTXOSet) Update(block *Block) {
	db := u.UBlockChain.Db

	err := db.Update(func(txn *badger.Txn) error {
		for _, tx := range block.Transactions {
			if tx.IsCoinbase() == false {
				for _, in := range tx.TXInputs {
					updatedOuts := TXOutputs{}
					inID := append(utxoPrefix, in.ID...)
					item, err := txn.Get(inID)
					utils.Handle(err)
					v, err := item.Value()

					outs := DeserializeOutputs(v)

					for outIdx, out := range outs.TXOutputs {
						if outIdx != in.Out {
							updatedOuts.TXOutputs = append(updatedOuts.TXOutputs, out)
						}
					}

					if len(updatedOuts.TXOutputs) == 0 {
						if err := txn.Delete(inID); err != nil {
							log.Panic(err)
						}
					} else {
						if err := txn.Set(inID, updatedOuts.Serialize()); err != nil {
							log.Panic(err)
						}
					}
				}
			}

			newOutputs := TXOutputs{}
			for _, out := range tx.TXOutputs {
				newOutputs.TXOutputs = append(newOutputs.TXOutputs, out)
			}

			txID := append(utxoPrefix, tx.ID...)
			if err := txn.Set(txID, newOutputs.Serialize()); err != nil {
				log.Panic(err)
			}

		}
		return nil
	})
	utils.Handle(err)
}

func (u UTXOSet) DeleteByPrefix(prefix []byte) {

	deleteKeys := func(keysForDelete [][]byte) error {
		if err := u.UBlockChain.Db.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}

	collectSize := 100000
	_ = u.UBlockChain.Db.View(func(txn *badger.Txn) error {
		//数据库内置迭代器
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		keysForDelete := make([][]byte, 0, collectSize)
		keysCollected := 0
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysForDelete = append(keysForDelete, key)
			keysCollected++
			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					log.Panic(err)
				}
				keysForDelete = make([][]byte, 0, collectSize)
				keysCollected = 0
			}
		}

		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				log.Panic(err)
			}
		}
		return nil
	})

}



