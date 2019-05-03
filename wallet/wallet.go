package wallet

import (
	"bytes"
	//"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/azd1997/golang-blockchain/utils"
	//ripemd160 "github.com/azd1997/golang-blockchain/mycrypto/myripemd160"
	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	WPrivateKey ecdsa.PrivateKey
	WPublicKey  []byte
}

//方法列表
//1.func NewKeyPair() (ecdsa.PrivateKey, []byte)
//2.func MakeWallet() *Wallet
//3.func PublicKeyHash(pubKey []byte)
//4.func Checksum(payload []byte) []byte
//5.func (w Wallet) Address() []byte
//6.func ValidateAddress(address string) bool


/*生成ECDSA公私钥对*/
//TODO：弄清ECDSA代码实现
//curve -> ecdsa -> privateKey/publicKey
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	utils.Handle(err)

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return *privateKey, publicKey
}

/*生成钱包对象*/
func MakeWallet() *Wallet {
	privateKey, publicKey := NewKeyPair()
	wallet := Wallet{privateKey, publicKey}

	return &wallet
}

/*对公钥取sha256哈希，再进行RIPEMD160哈希*/
//publicKey -> sha256 -> publicKeyHash -> ripemd160 -> publicKeyHash(ripemd160)
//TODO:理解RIPEMD160
func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	//hasher := crypto.RIPEMD160.New()
	_, err := hasher.Write(pubHash[:])
	utils.Handle(err)

	publicRipMD := hasher.Sum(nil)

	return publicRipMD
}

/*对公钥哈希进行双sha256哈希，再取前若干字节作为校验码*/
//publicKeyHah(ripemd160) -> sha256 -> sha256 -> [:checksum] -> checksum
func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

/*将公钥哈希、校验码、版本号合一进行base58编码得到账户地址*/
//publicKeyHash(ripemd160) + checksum + version -> base58 -> address
func (w Wallet) Address() []byte {
	//取公钥哈希
	pubHash := PublicKeyHash(w.WPublicKey)
	//将公钥哈希和版本号拼接成新slice切片
	versionedHash := append([]byte{version}, pubHash...) //PubHash...表示将字节切片中的内容打散再做操作
	//对包含了version和公钥哈希信息的slice取校验码
	checksum := Checksum(versionedHash)

	//再把校验码也给打散拼接上
	fullHash := append(versionedHash, checksum...)
	//将之转换成地址
	address := utils.Base58Encode(fullHash)

	fmt.Printf("Pub Key: %x\n", w.WPublicKey)
	fmt.Printf("Pub Hash: %x\n", pubHash)
	fmt.Printf("Address: %x\n", address)

	return address
}

/*账户（钱包创建流程）*/
//1.privateKey -> ecdsa -> publicKey -> sha256 -> ripemd160 -> publicKeyHash
//2.publicKeyHash -> sha256 -> sha256 -> 4bytes -> checksum -> 3.
//publicKeyHash -> 3.
//version -> 3.
//3.checksum;publicKeyHash;version} -> base58 -> address

/*数字签名实现*/
//Address
//FullHash
//[Version]
//[Pub Key Hash]
//[CheckSum]

/*验证钱包地址是否是合法的钱包地址*/
func ValidateAddress(address string) bool {
	//由字符串钱包地址解码得到所谓的公钥哈希（加入了校验码和版本号的）
	pubKeyHash := utils.Base58Decode([]byte(address))
	//取出检验码
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	//取出版本号
	version := pubKeyHash[0]
	//取出实际上的公钥哈希（sha256+ripemd160）
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength-1]
	//将得到的实际公钥哈希在本地进行一次计算校验码
	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))

	//返回两个校验码是否相等，相等说明钱包地址有效
	return bytes.Compare(actualChecksum, targetChecksum) == 0
}