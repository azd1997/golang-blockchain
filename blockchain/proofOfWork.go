package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

//工作量证明流程
//1.从Block中取数据
//2.创建从0开始的迭代器nonce
//3.创建区块数据加上随机数后的哈希值
//4.检查哈希是否满足target等要求

//要求：
//哈希的前面有若干个0

const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

/*创建新的工作量证明对象*/
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)                  //Target初始化为大数 1
	target.Lsh(target, uint(256-Difficulty)) //Target左移(256-Difficulty)位，得到需要的Target

	pow := &ProofOfWork{b, target}

	return pow
}

/*初始化数据，将区块数据拼接成字节数组*/
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(), //注意此时没有Hash，需要后边计算再赋进来
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{})
	return data
}

/*工作量证明运行程序，返回有效的Nonce和哈希*/
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	//采用  for {	CODE	}进行死循环也是一样的
	for nonce < math.MaxInt64 {
		//计算区块哈希，直至哈希满足条件
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data) //sha256哈希计算

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:]) //哈希值字节数组转为大数型，用以比较

		if intHash.Cmp(pow.Target) == -1 { // intHash < Target
			break
		} else {
			nonce++
		}

	}
	fmt.Println()

	return nonce, hash[:]
}

/*区块的哈希验证*/
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	//满足条件时 Cmp == -1 ，返回True
	return intHash.Cmp(pow.Target) == -1
}


