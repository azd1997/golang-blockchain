package blockchain

import (
	"bytes"
	"encoding/gob"
	"github.com/azd1997/golang-blockchain/utils"
)

type TXOutputs struct {
	TXOutputs []TXOutput
}


//方法列表
//1.func (outs TXOutputs) Serialize() []byte
//2.
//3.

/*对当前交易的交易输出集合进行序列化*/
func (outs TXOutputs) Serialize() []byte {
	//将outs编码存入buffer缓冲器
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(outs)
	utils.Handle(err)

	return buffer.Bytes()
}

/*将序列化后的交易输出数据反序列化*/
func DeserializeOutputs(data []byte) TXOutputs {

	var outputs TXOutputs
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&outputs)
	utils.Handle(err)

	return outputs
}