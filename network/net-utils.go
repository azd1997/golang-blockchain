package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/azd1997/golang-blockchain/utils"
)

/*将命令（字符串）转为字节数组*/
func CmdToBytes(cmd string) []byte {
	var cmdBytes [commandLength]byte

	for i, c := range cmd {
		cmdBytes[i] = byte(c)
	}

	return cmdBytes[:]
}

/*将命令（字节数组）转为字符串并输出*/
func BytesToCmd(cmdBytes []byte) string {
	var cmd []byte

	for _, b := range cmdBytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}

	return fmt.Sprintf("%s", cmd)
}

/*将数据进行编码得到字节数组*/
func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	utils.Handle(err)

	return buff.Bytes()
}

/*从request信息中抽取前12字节作为命令*/
func ExtractCmd(request []byte) []byte {
	return request[:commandLength]
}
