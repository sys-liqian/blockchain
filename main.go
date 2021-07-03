package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"time"
)

//区块
type Block struct {
	Timestamp    int64  //当前时间戳
	PreBlockHash []byte //当前块前一块hash
	Hash         []byte //当前块hash
	Data         []byte //当前块数据
	Nonce        int    //计数器
}

func (b *Block) setHash() {
	// 将时间戳转为10进制int
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PreBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

//Block构造
func NewBlock(data string, preBlockHash []byte) *Block {
	block := &Block{
		Timestamp:    time.Now().Unix(),
		PreBlockHash: preBlockHash,
		Data:         []byte(data),
		Nonce:        0,
	}
	//block.setHash()
	p := NewProofOfWork(block) // 创建工作量证明
	nonce, hash := p.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

//链
type BlockChain struct {
	Blocks []*Block
}

//添加块到链
func (b *BlockChain) AddBlock(data string) {
	//计算当前链中最后一个Block的Hash
	preBlockHash := b.Blocks[len(b.Blocks)-1]
	block := NewBlock(data, preBlockHash.Hash)
	b.Blocks = append(b.Blocks, block)
}

//创建创世块
func NewGenesisBlock() *Block {
	return NewBlock("hello block chain", []byte{})
}

//初始化有创世块的链
func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}

//挖矿难度值
const targetBits = 32

//工作量证明
type ProofOfWork struct {
	Block  *Block   //指向一个块的指针
	target *big.Int //上边界
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}
	return pow
}

func (p *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			p.Block.PreBlockHash,
			p.Block.Data,
			IntToHex(p.Block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func IntToHex(i int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, i)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func (p *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int //hash的整形表示
	var hash [32]byte
	nonce := 0 //计数器

	for nonce < math.MaxInt64 {
		data := p.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		//hashInt < target 正解
		if hashInt.Cmp(p.target) == -1 {
			fmt.Printf("\r%x\n", hash)
			break
		} else {
			nonce += 1
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func main() {
	chain := NewBlockChain()

	startTime := time.Now()
	chain.AddBlock("Send 1 BTC to Ivan")
	chain.AddBlock("Send 2 BTC to Ivan")
	elapse:=time.Since(startTime)

	for _, v := range chain.Blocks {
		fmt.Printf("Pre.Hash: %x\n", v.PreBlockHash)
		fmt.Printf("Data    : %s\n", v.Data)
		fmt.Printf("Hash    : %x\n", v.Hash)
		fmt.Println("Nonce   : ",v.Nonce)
		fmt.Println()
	}
	fmt.Println("耗时:",elapse)


}
