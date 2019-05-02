package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"github.com/azd1997/golang-blockchain/utils"
	"io/ioutil"
	"os"
)

const walletFile = "./tmp/wallets/wallets.data"

//注意wallets是要一直维护的，所以所有调用其的操作需要改变其内容时，一定要用指针
type Wallets struct {
	Wallets map[string]*Wallet
}

/*将wallets字典维护的内容编码之后写进文本*/
func (ws *Wallets) SaveFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	utils.Handle(err)

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	utils.Handle(err)
}

/*从文本文件加载钱包文件，解码后还原出钱包字典*/
func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	fileContent, err := ioutil.ReadFile(walletFile)
	utils.Handle(err)

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	utils.Handle(err)

	ws.Wallets = wallets.Wallets

	return nil
}

/*创造钱包字典对象，从钱包文件中读内容，赋给钱包字典*/
func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile()

	return &wallets, err
}

/*将钱包地址作为键，从钱包字典中查找对应钱包*/
func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

/*从钱包字典获取所有钱包地址，并存入钱包地址的切片数组中*/
func (ws *Wallets) GetAllAddress() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

/*生成新钱包并加入钱包字典，返回钱包地址*/
func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet

	return address
}
