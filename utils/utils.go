package utils

import (
	"bytes"
	"encoding/binary"
	"github.com/mr-tron/base58"
	"log"
)

/*err处理程序*/
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

/*整型数据转换成十六进制的字节数组，用以让Nonce、Difficulty参与区块哈希计算*/
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	Handle(err)
	return buff.Bytes()
}

//TODO:Base58编码解码过程理解
/*基于第三方base58库实现的Base58编码*/
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

/*基于第三方base58库实现的Base58解码*/
func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	Handle(err)

	return decode
}

//注意点，经常性的，注释中提及[]byte为字节数组，在go中，[]byte称为切片slice，[32]byte（指定长度）称为数组
