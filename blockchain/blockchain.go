package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"runtime"
)

const genesisData = "First Transaction from Genesis"

type BlockChain struct {
	//Blocks []*Block
	LastHash []byte
	Db       *badger.DB
}

/*创建带有创世区块的区块链，创世区块需指定创世区块coinbase收款人地址*/
func InitBlockChain(address string) *BlockChain {
	var lastHash []byte

	//检查区块链是否存在，不存在才执行下边的初始化区块链流程
	if DbExists() {
		fmt.Println("区块链已存在...")
		runtime.Goexit()
	}

	//打开数据库
	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	//更新数据库，存入创世区块和lasthash
	err = db.Update(func(txn *badger.Txn) error {

		//创世区块的coinbase交易
		cbtx := CoinbaseTx(address, genesisData)
		//创世区块
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created...")
		//存入区块链的第一个区块的键值对
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		//安排一个键值对用来存储链上最新区块的哈希，在工程代码里常称为lasthash、lh
		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err

	})
	Handle(err)

	//创建BlockChain对象并返回
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

/*区块链已存在时，调用此函数，创建并返回此时最新的区块链对象*/
//TODO:ContinueBlockChain调用了不必要的参数
func ContinueBlockChain(address string) *BlockChain {
	if DbExists() == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	//配置并打开数据库
	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	//查取("lh", lastHash)
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()
		return err
	})
	Handle(err)

	//创建并返回BlockChain对象
	chain := BlockChain{lastHash, db}

	return &chain
}

/*向区块链中添加新区块*/
func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var lastHash []byte
	//获取lastHash
	err := chain.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()
		return err
	})
	Handle(err)

	//将Transactions和PrevHash(lastHash)，打包、工作量证明，挖出新区块
	newBlock := CreateBlock(transactions, lastHash)

	//将新区块信息存入数据库；更新数据库中lastHash；更新BlockChain对象中lastHash
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
func (chain *BlockChain) Iterator() *BCIterator {
	iter := &BCIterator{chain.LastHash, chain.Db}

	return iter
}

/*寻找某账户未花费的交易*/
//返回：包含有 该用户的未花费输出 的交易 的集合
//TODO:好好理解这段代码！
//	思考：这段代码为何不分两段走，先循环遍历得到已花费输出，再循环遍历在所有输出中过滤掉已花费输出？
func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	//未花费的交易集合
	var unspentTxs []Transaction

	//已花费的交易输出的集合，map类型
	spentTXOs := make(map[string][]int) // string -> []int

	//遍历区块链中从后往前所有区块
	iter := chain.Iterator()

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
				if out.CanBeUnlockedWith(address) {
					unspentTxs = append(unspentTxs, *tx)
				}

			}

			//判断非Coinbase之后，再检查交易的输入
			//检查交易输入的作用是，通过循环迭代，获取spentTXOs map
			if tx.IsCoinbase() == false {
				for _, in := range tx.TXInputs {
					//检查每笔交易的输入（也就是已花费的输出）是否是这个账户地址花出去的
					if in.CanUnlockOutputWith(address) {
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
func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	//在未花费交易的输出里边找到这个账户地址的未花费输出
	for _, tx := range unspentTransactions {
		for _, out := range tx.TXOutputs {
			if out.CanBeUnlockedWith(address) {
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
func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)                //未花费输出，记录 UTXO_TX_ID -> OutNumOfThisAddress
	unspentTxs := chain.FindUnspentTransactions(address) //未花费交易
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.TXOutputs {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
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
