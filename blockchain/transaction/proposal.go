/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Description:
 * @File:
 * @Version: 1.0.0
 * @Date: 2022.05.29
 */

package transaction

import (
	"encoding/json"
	"github.com/llshrly/K-Times-Redactable-Blockchain/common/chamelemon"
)

// Proposal 原始交易
type Proposal struct {
	// 上链消息
	Message []byte `json:"message"`

	/* 变色龙哈希 */
	// 监管方变色龙哈希主公钥mpk
	Mpk string `json:"mpk"`
	// 编辑方变色龙哈希子公钥spk
	Spk string `json:"spk"`

	// 哈希信息
	R             []byte `json:"r"`
	S             []byte `json:"s"`
	ChameleonHash []byte `json:"chameleon_hash"`
}

// newProposal 创建proposal
func newProposal(message string, mpk string, spk string) *Proposal {
	r, s, hash := chamelemon.Hash128([]byte(message))
	return &Proposal{
		Message:       []byte(message),
		Mpk:           mpk,
		Spk:           spk,
		R:             r,
		S:             s,
		ChameleonHash: hash,
	}
}

// Bytes proposal序列化
func (p *Proposal) Bytes() []byte {
	buf, _ := json.Marshal(p)
	return buf
}
