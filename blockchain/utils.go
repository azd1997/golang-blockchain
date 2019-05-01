package blockchain

import (
	"bytes"
	"encoding/binary"
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
