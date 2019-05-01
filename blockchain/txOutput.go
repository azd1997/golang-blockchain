package blockchain

//当前交易的输出，是其下个交易的输入，或者未花费
//转给某某多少钱
type TxOutput struct {
	Value  int    //转账金额
	PubKey string //转账对象地址
}

/*检查转账输出的公钥地址*/
//检查该笔交易的输出是否接受者用接收者给的数据进行解锁
//A向B转账，A填写的目标地址是B的公钥地址，B需要用B的私钥来解锁这笔钱（或者说这笔未花费输出）
func (out *TxOutput) CanBeUnlockedWith(data string) bool {
	return out.PubKey == data
}
