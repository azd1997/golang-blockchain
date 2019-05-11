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
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	//将得到的实际公钥哈希在本地进行一次计算校验码
	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))

	//返回两个校验码是否相等，相等说明钱包地址有效
	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

//神秘的比特币地址详解
//当你看到像这样的一串字符的时候你是什么感想：1M8DPUBQXsVUNnNiXw5oFdRciguXctWpUD，如果在你接触比特币之前，你一定会说这不就是一堆乱码吗？没错这是在你认识比特币之前的时候，而在认识了比特币之后，你所谓的乱码就是你的比特币地址，这个地址就好像你的银行卡账户那样，可以方便快捷的查询和交易你的比特币。
//那么为什么会用这样的一种格式来作为比特币的地址呢？我们还是慢慢的来的了解吧。
//首先，我们常用的比特币地址格式一般有四种：
//1、BASE58格式
//就是人们常说的比特币地址，由1开头的，例如：1M8DPUBQXsVUNnNiXw5oFdRciguXctWpUD
//2、HASH160格式
//Tab content 由RIPEMD160算法对130位公钥的SHA256签名进行计算的结果，
//如：fbfb58defc272942fc31d00c007b59aa4cb5087a
//3、WIF压缩格式
//即钱包输入格式，是将BASE58格式进行压缩后的结果130位公钥格式 这是最原始的由ECDSA算法计算出来的比特币公钥，
//如：0469B0E479C9A358908DB9CF4628BDD643C3F8
//1C4F0096AAD442DA6CA8BCC4FD86A8D47D7A865E178B6D062CC9B70290
//8973952062A1D767DA9B2BD2095D5CCF6E
//4、60位公钥格式
//130位公钥进行压缩后的结果，如：0269B0E479C9A358908DB9CF4628BDD643
//C3F81C4F0096AAD442DA6CA8BCC4FD86

//那么，这些复杂的数字和字符是怎么产生的呢？
//首先，让我们先简单的说说比特币地址是怎么算出来的。比特币是建立在数学加密学基础上的，中本聪大神用了椭圆加密算法（ECDSA）来产生比特币的私钥和公钥。由私钥是可以计算出公钥的，公钥的值经过一系列数字签名运算会得到比特币地址。
//需要说明的是：因为由公钥可以算出比特币地址，所以我们经常把公钥和比特币地址的说法相混淆，但是他们都是指的一个概念。比特币地址只是另一种格式的公钥。
//从比特币私钥得到我们所用的比特币地址需要九个步骤。

//第一步，随机选取一个32字节的数、大小介于1 ~ 0xFFFF FFFF FFFF FFFF FFFF FFFF FFFF FFFE BAAE DCE6 AF48 A03B BFD2 5E8C D036 4141之间，作为私钥。
//18E14A7B6A307F426A94F8114701E7C8E774E7F9A47E2C2035DB29A206321725
//第二步，使用椭圆曲线加密算法（ECDSA-secp256k1）计算私钥所对应的非压缩公钥。 (共65字节， 1字节 0x04, 32字节为x坐标，32字节为y坐标）关于公钥压缩、非压缩的问题另文说明。
//
//0450863AD64A87AE8A2FE83C1AF1A8403CB53F53E486D8511DAD8A04887E5B
//23522CD470243453A299FA9E77237716103ABC11A1DF38855ED6F2EE187E9C582BA6
//第三步，计算公钥的 SHA-256 哈希值
//600FFE422B4E00731A59557A5CCA46CC183944191006324A447BDB2D98D4B408
//第四步，取上一步结果，计算 RIPEMD-160 哈希值
//010966776006953D5567439E5E39F86A0D273BEE
//第五步，取上一步结果，前面加入地址版本号（比特币主网版本号“0x00”）
//00010966776006953D5567439E5E39F86A0D273BEE
//第六步，取上一步结果，计算 SHA-256 哈希值
//445C7A8007A93D8733188288BB320A8FE2DEBD2AE1B47F0F50BC10BAE845C094
//第七步，取上一步结果，再计算一下 SHA-256 哈希值（哈哈）
//D61967F63C7DD183914A4AE452C9F6AD5D462CE3D277798075B107615C1A8A30
//第八步，取上一步结果的前4个字节（8位十六进制）
//D61967F6
//第九步，把这4个字节加在第五步的结果后面，作为校验（这就是比特币地址的16进制形态）。
//00010966776006953D5567439E5E39F86A0D273BEED61967F6
//第十步，用base58表示法变换一下地址（这就是最常见的比特币地址形态）。
//1M8DPUBQXsVUNnNiXw5oFdRciguXctWpUD
//比特币地址生成过车个就是这样，那么会有人问道，既然都是随机生成的，那么比特币的地址会不会重复呢？关于这个问题，想必就更不用担心。因为比特币的私钥长度是256位的二进制串，那么随机生成的两个私钥正好重复的的概率是2^256≈10^77之一，这个数字大到你根本无法想象，比中彩票的概率还要小好多。