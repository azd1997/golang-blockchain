package blockchain

import (
	"bytes"
	"github.com/azd1997/golang-blockchain/wallet"
)

//当前交易的输入  是其来源交易的输出
//张三、李四、王五等人曾经给过我一笔钱，现在我是用这笔钱，并对这笔钱签名，转给某某
type TXInput struct {
	ID  []byte //来源交易的ID
	Out int    //来源转账在来源交易中是第几个输出
	//Sig string //转账发起者的签名
	Signature []byte
	PubKey    []byte //公钥
}

/*检查转账交易的私钥签名（检查该次交易是否是这个人发起的）*/
//验证这笔交易时调用。
//比特币中标准版本是：我检查这笔交易时，取这笔交易转账者公钥来对其签名进行验证
//这里先进行简化验证：给定data(转账者的地址信息或其它绑定信息)，看和转账者签名对不对得上
//对的上说明：这笔交易的输入确实这个转账者可以使用（可以解锁）
//比如说：王五曾经给了张三100块，现在张三想用这笔钱转给李四。我们需要检查这笔钱是不是张三在用，防止被别人盗用
//张三提供其签名，而其公钥地址在王五给他的转账中已经公开，所以我们拿来看看对不对得上
//arg data 一般是指转账者的公钥地址

//A向B转账，A用自己的私钥签名，其他人用A的公钥进行验证
//func (in *TxInput) CanUnlockOutputWith(data string) bool {
//	//如果输入的数据能和转账者签名对的上，则表示这笔钱可以解锁（可以使用）
//	return in.Sig == data
//}

/*检查一笔交易的输入，这笔交易将 给入的公钥哈希H1 和利用交易输入中的公钥生成的公钥哈希H2 进行比较*/
//匹配那么就可以取出这笔钱使用
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

