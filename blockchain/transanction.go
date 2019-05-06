package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/azd1997/golang-blockchain/utils"
	"github.com/azd1997/golang-blockchain/wallet"
	"log"
	"math/big"
	"strings"
)

type Transaction struct {
	ID        []byte //即交易哈西
	TXInputs  []TXInput
	TXOutputs []TXOutput
}

//方法列表
//1.func DeserializeTransaction(data []byte) Transaction
//2.func (tx *Transaction) IsCoinbase() bool
//3.func CoinbaseTx(to, data string) *Transaction
//4.func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction
//5.func (tx Transaction) Serialize() []byte
//6.func (tx *Transaction) Hash() []byte
//7.func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction)
//8.func (tx *Transaction) TrimmedCopy() Transaction
//9.func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool
//10.func (tx Transaction) String() string

func DeserializeTransaction(data []byte) Transaction {

	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	utils.Handle(err)

	return transaction
}

/*判断交易是否是Coinbase交易*/
func (tx *Transaction) IsCoinbase() bool {
	//交易输入只有1个且该交易输入的ID也就是哈希为0（ID为空，所以哈希为0）
	return len(tx.TXInputs) == 1 && len(tx.TXInputs[0].ID) == 0
}

/*出块奖励交易*/
func CoinbaseTx(to, data string) *Transaction {

	//若Coinbase交易未指定Data内容，则默认为下方内容
	//对于挖出区块的矿工而言，可以在Coinbase交易的data域填想填的东西
	if data == "" {
		//长度为24字节的随机数据
		randData := make([]byte, 24)
		_, err := rand.Read(randData)
		utils.Handle(err)

		data = fmt.Sprintf("%x", randData)
	}

	//Coinbase来源交易不存在，所以填空字节，其来源交易占来源输出序号也不存在，这里以-1表示
	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(100, to)

	//Coinbase交易只有一笔输入一笔输出，其交易ID或者说哈希需要进行哈希才能得到
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	//tx.SetID()
	tx.ID = tx.Hash()

	return &tx
}

/*产生一笔新交易*/
func NewTransaction(w *wallet.Wallet, to string, amount int, UTXO *UTXOSet) *Transaction {
	var inputs []TXInput   //当前交易的输入
	var outputs []TXOutput //当前交易输出


	pubKeyHash := wallet.PublicKeyHash(w.WPublicKey)

	//用户进行转账时，需指定转账者A，被转账者B以及转账金额S
	//主要需要考虑以下几点：
	//1.找到A的所有UTXO，计算A的余额，检查A余额是否足够
	//2.检查A余额

	acc, validOutputs := UTXO.FindSpendableOutputs(pubKeyHash, amount)
	//注意返回的acc是有可能小于amount的！！！

	//余额不够，报错
	if acc < amount {
		log.Panic("Error: not enough funds...")
	}

	for txid, outs := range validOutputs { //txid为键，outs为值
		txID, err := hex.DecodeString(txid) //由字符串解码回十六进制字节数组
		utils.Handle(err)

		for _, out := range outs {
			input := TXInput{txID, out, nil, w.WPublicKey}
			inputs = append(inputs, input)
		}
	}

	from := fmt.Sprintf("%s", w.Address())
	outputs = append(outputs, *NewTXOutput(amount, to))

	//增加找零的交易输出
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	//对交易进行签名
	UTXO.UBlockChain.SignTransaction(&tx, w.WPrivateKey)

	return &tx
}

/*将交易序列化成字节切片数组*/
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	utils.Handle(err)

	return encoded.Bytes()
}

/*对交易进行序列化并取哈希*/
//将交易的ID置空，计算哈希。
//和setID有点像，但有区别，setID不返回
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

//TODO:理解
/*使用私钥对交易的来源交易字典进行签名*/
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//Coinbase交易没有签名
	if tx.IsCoinbase() {
		return
	}

	//查提供的来源交易（除了Coinbase）里边, 按照TXInput里面的来源交易ID去查，看有没有
	for _, in := range tx.TXInputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct!")
		}
	}

	//获取不含签名和接收者公钥的交易副本
	txCopy := tx.TrimmedCopy()

	for inId, in := range txCopy.TXInputs {
		prevTX := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.TXInputs[inId].Signature = nil
		txCopy.TXInputs[inId].PubKey = prevTX.TXOutputs[in.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.TXInputs[inId].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		utils.Handle(err)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.TXInputs[inId].Signature = signature

	}

}

/*获取被裁剪的交易对象（没有接收者公钥和转账者签名）*/
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	//获取交易输入集合（但不含公钥和签名）
	for _, in := range tx.TXInputs {
		inputs = append(inputs, TXInput{in.ID, in.Out, nil, nil})
	}

	//获取交易输出集合
	for _, out := range tx.TXOutputs {
		outputs = append(outputs, TXOutput{out.Value, out.PubKeyHash})
	}

	//构建裁剪过的交易副本
	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

/*验证交易是否合法*/
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, in := range tx.TXInputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("Previous transaction does not exist!")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inId, in := range tx.TXInputs {
		prevTx := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.TXInputs[inId].Signature = nil
		txCopy.TXInputs[inId].PubKey = prevTx.TXOutputs[in.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.TXInputs[inId].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(keyLen / 2)])
		y.SetBytes(in.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}

	}

	return true
}

/*输出交易信息*/
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	for i, input := range tx.TXInputs {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:     %x", input.ID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Out))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.TXOutputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

