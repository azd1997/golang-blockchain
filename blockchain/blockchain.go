package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/azd1997/golang-blockchain/utils"
	"github.com/dgraph-io/badger"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const genesisData = "First Transaction from Genesis"

type BlockChain struct {
	//Blocks []*Block
	LastHash []byte
	Db       *badger.DB
}

//方法列表
//1.func InitBlockChain(address string) *BlockChain
//2.func ContinueBlockChain(address string) *BlockChain
//3.func (bc *BlockChain) AddBlock(transactions []*Transaction)
//4.func (bc *BlockChain) Iterator() *BCIterator
//5.func (bc *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction
//6.func (bc *BlockChain) FindUTXO(pubKeyHash []byte) []TXOutput
//7.func (bc *BlockChain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int)
//8.func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error)
//9.func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey)
//10.func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool
//11.func (bc *BlockChain) FindUTXO2() map[string]TXOutputs

//TODO:参数
/*创建带有创世区块的区块链，创世区块需指定创世区块coinbase收款人地址*/
//创世区块奖励接收者建议为	1111 1111 1111 1111 1111 1111 1111 1111 11
func InitBlockChain(address, nodeId string) *BlockChain {
	var lastHash []byte

	//检查区块链是否存在，不存在才执行下边的初始化区块链流程
	path := fmt.Sprintf(dbPath, nodeId)
	if DbExists(path) {
		fmt.Println("区块链已存在...")
		runtime.Goexit()
	}

	//打开数据库
	opts := badger.DefaultOptions
	opts.Dir = path
	opts.ValueDir = path

	db, err := openDB(path, opts)
	utils.Handle(err)

	//更新数据库，存入创世区块和lastHash
	err = db.Update(func(txn *badger.Txn) error {

		//创世区块的coinbase交易
		cbTx := CoinbaseTx(address, genesisData)
		//创世区块
		genesis := Genesis(cbTx)
		fmt.Println("Genesis created...")
		//存入区块链的第一个区块的键值对
		err = txn.Set(genesis.Hash, genesis.Serialize())
		utils.Handle(err)
		//安排一个键值对用来存储链上最新区块的哈希，在工程代码里常称为lasthash、lh
		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err

	})
	utils.Handle(err)

	//创建BlockChain对象并返回
	blockChain := BlockChain{lastHash, db}
	return &blockChain
}

/*区块链已存在时，调用此函数，创建并返回此时最新的区块链对象*/
//TODO:ContinueBlockChain调用了不必要的参数
func ContinueBlockChain(nodeId string) *BlockChain {

	//检查数据库是否存在
	path := fmt.Sprintf(dbPath, nodeId)
	if DbExists(path) == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	//配置并打开数据库
	opts := badger.DefaultOptions
	opts.Dir = path
	opts.ValueDir = path

	db, err := openDB(path, opts)
	utils.Handle(err)

	//查取("lh", lastHash)
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		lastHash, err = item.Value()
		return err
	})
	utils.Handle(err)

	//创建并返回BlockChain对象
	blockChain := BlockChain{lastHash, db}

	return &blockChain
}

/*向区块链中 挖出 新区块*/
func (bc *BlockChain) MineBlock(transactions []*Transaction) *Block {

	//本地验证交易有效性
	for _, tx := range transactions {
		if bc.VerifyTransaction(tx) != true {
			log.Panic("Invalid Transaction")
		}
	}

	var lastHash []byte
	var lastHeight int
	err := bc.Db.View(func(txn *badger.Txn) error {
		//获取lastHash
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		lastHash, err = item.Value()
		//获取lastHeight
		item, err = txn.Get([]byte(lastHash))
		utils.Handle(err)
		lastBlockData, _ := item.Value()
		lastBlock := Deserialize(lastBlockData)
		lastHeight = lastBlock.Height

		return err
	})
	utils.Handle(err)

	//将Transactions和PrevHash(lastHash)，打包、工作量证明，挖出新区块
	newBlock := CreateBlock(transactions, lastHash, lastHeight+1)

	//将新区块信息存入数据库；更新数据库中lastHash；更新BlockChain对象中lastHash
	err = bc.Db.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		utils.Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		bc.LastHash = newBlock.Hash

		return err
	})
	utils.Handle(err)

	return newBlock
}

/*向区块链中 添加 新区块*/
//这个主要是用于当从别的节点接收最新区块时，将这些区块加入到本地区块链
func (bc *BlockChain) AddBlock(block *Block) {
	err := bc.Db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(block.Hash); err == nil {
			return nil
		}

		//将区块存入
		blockData := block.Serialize()
		err := txn.Set(block.Hash, blockData)
		utils.Handle(err)

		//取出最后区块
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		lastHash, _ := item.Value()
		item, err = txn.Get(lastHash)
		utils.Handle(err)
		lastBlockData, _ := item.Value()
		lastBlock := Deserialize(lastBlockData)

		//比较最后区块和刚加入的区块的区块高度
		if block.Height > lastBlock.Height {
			err := txn.Set([]byte("lh"), block.Hash)
			utils.Handle(err)
			bc.LastHash = block.Hash
		}

		return nil
	})
	utils.Handle(err)
}

/*从区块链中查询区块*/
func (bc *BlockChain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	err := bc.Db.View(func(txn *badger.Txn) error {
		if item, err := txn.Get(blockHash); err != nil {
			return errors.New("block is not found")
		} else {
			blockData, _ := item.Value()
			block = *Deserialize(blockData)
		}
		return nil
	})

	if err != nil {
		return block, err
	}

	return block, nil
}

/*获取区块链所有区块哈希集合，用以快速验证不同节点间区块链的一致性*/
func (bc *BlockChain) GetBlockHashes() [][]byte {
	var blockHashes [][]byte

	iter := bc.Iterator()

	for {
		block := iter.Next()

		blockHashes = append(blockHashes, block.Hash)

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return blockHashes
}

func (bc *BlockChain) GetBestHeight() int {
	var lastBlock Block

	err := bc.Db.View(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		lastHash, _ := item.Value()

		item, err = txn.Get(lastHash)
		utils.Handle(err)
		lastBlockData, _ := item.Value()

		lastBlock = *Deserialize(lastBlockData)

		return nil
	})
	utils.Handle(err)

	return lastBlock.Height
}

/*返回区块链迭代器对象*/
func (bc *BlockChain) Iterator() *BCIterator {
	iter := &BCIterator{bc.LastHash, bc.Db}

	return iter
}

func (bc *BlockChain) FindUTXO2() map[string]TXOutputs {

	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)

	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.TXOutputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.TXOutputs = append(outs.TXOutputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.TXInputs {
					inTxID := hex.EncodeToString(in.ID)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}

	}
	return UTXO

}

/*寻找某账户未花费的交易*/
//返回：包含有 该用户的未花费输出 的交易 的集合
//TODO:好好理解这段代码！
//	思考：这段代码为何不分两段走，先循环遍历得到已花费输出，再循环遍历在所有输出中过滤掉已花费输出？
func (bc *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	//未花费的交易集合
	var unspentTxs []Transaction

	//已花费的交易输出的集合，map类型
	spentTXOs := make(map[string][]int) // string -> []int

	//遍历区块链中从后往前所有区块
	iter := bc.Iterator()

	for {
		block := iter.Next()

		//遍历区块中所有交易
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID) //将十六进制的字节数组转换成字符串

			//continue语句只结束本次循环，而不终止整个循环的执行。而break语句则是结束整个循环过程，不再判断执行循环的条件是否成立。
			//goto语句：用几个字来描述就是简单粗暴，也称为无条件转移语句，其一般格式如下： goto 语句标号。 
			//其中语句标号是按标识符规定书写的符号， 放在某一语句行的前面，标号后加冒号(：)。语句标号起标识语句的作用，
			//与goto 语句配合使用。在结构化程序设计中一般不主张使用goto语句， 以免造成程序流程的混乱，使理解和调试程序
			//都产生困难。很容易造成bug。
			//如果使用 continue + 标识符，可以很好的利用continue与goto语句的优点。
		Outputs: //Outputs代码段包含了检查交易输出和输入的两个循环
			//检查该交易的输出
			//1.是否是未花费输出
			//2.是否是这个账户地址的
			for outIdx, out := range tx.TXOutputs {
				//因为遍历循环，最终spentTXOs map会存有所有未花费交易输出所在交易的txID

				//如果在已花费输出中找到了这笔交易，说明这笔交易的输出被花过
				//再检查这笔交易的所有输出是否与外部for循环的交易的输出序号对得上与否
				//这一步似乎是在考虑哈希碰撞的情况
				//如果说在交易哈西一致，且所有输出都对得上，认为是同一个交易

				//这一段用来过滤掉所有已花费的输出
				if spentTXOs[txID] != nil { //说明该笔交易是该地址的已花费输出的交易
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx { //判断条件：该笔交易中该账户地址所持的输出的编号 spentOut;
							// 是否等于该笔交易中这笔输出的编号
							continue Outputs //跳出此次循环，同时保留当前值，再一次从Outputs处执行
						}
					}
				}

				//检查这笔输出（未花费的）是否能被该账户地址用他的私钥信息解锁
				//也就是检查这笔输出是否是这个人的 未花费输出
				//if out.CanBeUnlockedWith(address) {
				//	unspentTxs = append(unspentTxs, *tx)
				//}
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTxs = append(unspentTxs, *tx)
				}

			}

			//判断非Coinbase之后，再检查交易的输入
			//检查交易输入的作用是，通过循环迭代，获取spentTXOs map
			if tx.IsCoinbase() == false {
				for _, in := range tx.TXInputs {
					//检查每笔交易的输入（也就是已花费的输出）是否是这个账户地址花出去的
					//if in.CanUnlockOutputWith(address) {
					//	inTxID := hex.EncodeToString(in.ID) //inTxID是所有与该账户参与输入的交易
					//	spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					//}

					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.ID) //inTxID是所有与该账户参与输入的交易
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)

					}
				}
			}

			//显然！
			//每笔交易中每个账户地址都至多有一次输出，也就是说，
			//可以将每个账户的所有转出交易的交易ID和输出编号Out一一映射起来，得到unspentTXOs
			//找到了该账户的已花费输出之后就可以过滤掉已花费输出，找出未花费输出

		}

		//遍历至创世区块则停止
		if len(block.PrevHash) == 0 {
			break
		}
	}

	//返回该账户的所有的未花费输出所在的交易的集合
	return unspentTxs
}

/*返回该账户所有的未花费输出的集合，用来查询账户余额*/
func (bc *BlockChain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	//在未花费交易的输出里边找到这个账户地址的未花费输出
	for _, tx := range unspentTransactions {
		for _, out := range tx.TXOutputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out) //这里的out是TxOutput对象
			}
		}
	}

	return UTXOs
}

/*在该账户的未花费输出 里边去找  可以花费的一个或一组输出（这组输出总金额要大于转账金额）*/
//比如说：张三有来自100个交易的未花费输出，他想给李四转100块，他需要从第一个开始未花费输出开始遍历，并累加金额，直到累加金额超过100块。
//例如他在遍历至第17个输出时累加金额是99块，到第18个时是101块，所以，张三的可花费输出是第1~18个未花费输出的集合
//为什么不用上面的FindUTXO呢？因为它不反回交易ID信息，无法进行验证
//注意：！！！如果所有未花费输出总和也不满足amount条件，依然返回
func (bc *BlockChain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)                //未花费输出，记录 UTXO_TX_ID -> OutNumOfThisAddress
	unspentTxs := bc.FindUnspentTransactions(pubKeyHash) //未花费交易
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.TXOutputs {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}

/*UTXO流程*/
//1.遍历区块链中所有区块、所有交易、所有输入输出，找到所有没被花费的交易
//2.

/*找到指定交易ID的交易并返回*/
func (bc *BlockChain) FindTransaction(txID []byte) (Transaction, error) {
	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, txID) == 0 {
				return *tx, nil //找到交易
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("transaction does not exist")
}

/*使用私钥对当前交易所有的来源交易进行签名，用以转账*/
func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.TXInputs {
		prevTX, err := bc.FindTransaction(in.ID)
		utils.Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(privKey, prevTXs)
}

/*验证一笔交易，通过验证这笔交易的所有来源交易来实现*/
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {

	//检查是不是Coinbase交易，是则返回true
	if tx.IsCoinbase() {
		return true
	}

	//获取所有来源交易
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.TXInputs {
		prevTX, err := bc.FindTransaction(in.ID)
		utils.Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	//验证所有来源交易
	return tx.Verify(prevTXs)
}

/*开启数据库失败时调用*/
func retry(dir string, originOpts badger.Options) (*badger.DB, error) {

	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}
	retryOpts := originOpts
	retryOpts.Truncate = true //truncate 截短
	db, err := badger.Open(retryOpts)

	return db, err
}

/*打开数据库，一次不成则retry*/
func openDB(dir string, opts badger.Options) (*badger.DB, error) {
	if db, err := badger.Open(opts); err != nil {
		//报错信息包含“LOCK”，则retry
		if strings.Contains(err.Error(), "LOCK") {
			if db, err := retry(dir, opts); err == nil {
				log.Println("database unlocked, value log truncated")
				return db, nil
			}
			log.Println("could not unlock database:", err)
		}
		return nil, err
	} else {
		return db, nil
	}
}