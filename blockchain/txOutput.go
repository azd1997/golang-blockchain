package blockchain

import (
	"bytes"
	"github.com/azd1997/golang-blockchain/utils"
)

//当前交易的输出，是其下个交易的输入，或者未花费
//转给某某多少钱
type TXOutput struct {
	Value  int    //转账金额
	//PubKey string //转账对象地址
	PubKeyHash []byte //公钥哈希
}

//方法列表
//1.func (out *TXOutput) Lock(address []byte)
//2.func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool
//3.func NewTXOutput(value int, address string) *TXOutpu





/*检查转账输出的公钥地址*/
//检查该笔交易的输出是否接受者用接收者给的数据进行解锁
//A向B转账，A填写的目标地址是B的公钥地址，B需要用B的私钥来解锁这笔钱（或者说这笔未花费输出）
//func (out *TxOutput) CanBeUnlockedWith(data string) bool {
//	return out.PubKey == data
//}

/*对交易输出使用接收者的公钥进行上锁，使得只有接收者使用私钥才能解开*/
//本质是转账者将接收者地址转换成公钥哈希放进交易输出；接收者需要去用自己的私钥匹配（对其解锁）
func (out *TXOutput) Lock(address []byte) {

	//对传入的钱包地址进行解码得到解码前的 公钥哈希、版本号、校验码拼接成的字节串
	pubKeyHash := utils.Base58Decode(address)
	//丢掉字节串的最后3个字节（校验码）和最前头1个字节（版本号），
	//得到真正的经过sha256和ripemd160的公钥哈希
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

/*检查交易输出是否上锁*/
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

/*根据转账地址和金额生成新的交易输出，并且上锁*/
//本质就是创建一个交易输出对象
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}