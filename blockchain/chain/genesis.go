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
	"github.com/llshrly/K-Times-Redactable-Blockchain/blockchain/transaction"
	"time"
)

var GenesisBlock = Block{
	Index:     0,
	Timestamp: time.Now().String(),
	Tx:        transaction.SignedTx{},
	PrevHash:  "",
	R:         nil,
	S:         nil,
	Hash:      "",
}
