/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Description:
 * @File:
 * @Version: 1.0.0
 * @Date: 2022.05.29
 */

package chain

import (
	"encoding/hex"
	"github.com/llshrly/K-Times-Redactable-Blockchain/blockchain/transaction"
	"github.com/llshrly/K-Times-Redactable-Blockchain/common/chamelemon"
	"time"
)

type Block struct {
	// 区块数据
	Index     int                  // 这个块在整个链中的位置
	Timestamp string               // 块生成的时间戳
	Tx        transaction.SignedTx // 交易，每个区块只有一笔交易
	PrevHash  string               // 代表前一个块的SHA256散列值
	MetaData  string               // 可编辑区域, 目前就是唯一交易中的message

	// chameleon hash r, s
	R, S []byte
	Hash string // 这个块通过SHA256的散列值
}

//entity 签名实体
func (b *Block) entity() []byte {
	record := string(b.Index) + b.Timestamp + string(b.Tx.Bytes()) + string(b.PrevHash) + string(b.MetaData)
	return []byte(record)
}

// 用来计算给定的数据的 SHA256 散列值
func calculateHash(block Block) (r, s []byte, hash string) {
	r, s, hashBytes := chamelemon.Hash128(block.entity())
	hash = hex.EncodeToString(hashBytes)
	return
}

// 生成碰撞哈希区块
func generateCollision(oldBlock, newBlock Block, r1, s1 []byte) (r2, s2 []byte, hash2 string) {
	r2, s2, hash2Bytes := chamelemon.GenerateCollision128(oldBlock.entity(), newBlock.entity(), r1, s1)
	hash2 = hex.EncodeToString(hash2Bytes)
	return
}

// GenerateBlock 生成一个块
func GenerateBlock(oldBlock Block, tx *transaction.SignedTx) (Block, error) {
	var newBlock Block

	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Tx = *tx
	newBlock.PrevHash = oldBlock.Hash
	newBlock.MetaData = string(tx.Message)
	newBlock.R, newBlock.S, newBlock.Hash = calculateHash(newBlock)
	return newBlock, nil
}

// EditBlock 编辑区块
func EditBlock(oldBlock Block, message string) (Block, error) {
	var newBlock Block
	newBlock.Index = oldBlock.Index
	newBlock.Timestamp = oldBlock.Timestamp
	newBlock.Tx = oldBlock.Tx
	newBlock.PrevHash = oldBlock.PrevHash
	newBlock.MetaData = message
	newBlock.R, newBlock.S, newBlock.Hash = generateCollision(oldBlock, newBlock, oldBlock.R, oldBlock.S)
	return newBlock, nil
}

// IsBlockValid 校验块
func IsBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if string(oldBlock.Hash) != string(newBlock.PrevHash) {
		return false
	}
	r, s, hash := calculateHash(newBlock)
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}
	chamelemon.VerifyHash(r, s, hashBytes, newBlock.entity())
	return true
}
