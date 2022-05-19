/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Description:
 * @File:
 * @Version: 1.0.0
 * @Date: 2022.05.29
 */

package chain

import "errors"

var Blockchain []Block

// ReplaceChain 将本地的过期的链切换成最新的链, 最长链法则
func ReplaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func ReplaceBlock(newBlock Block, index int) error {
	if index >= len(Blockchain) {
		return errors.New("block index out of range")
	}
	Blockchain[index] = newBlock
	return nil
}
