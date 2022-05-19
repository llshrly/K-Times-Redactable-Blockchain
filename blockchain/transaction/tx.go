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
	"github.com/llshrly/K-Times-Redactable-Blockchain/blockchain/core"
)

// SignedTx 签过名的交易
type SignedTx struct {
	*Proposal

	// 签名
	Signature []byte `json:"signature"`
}

// NewSignedTx 新建交易
func NewSignedTx(signer core.Signer, message string, mpk string, spk string) (*SignedTx, error) {
	proposal := newProposal(message, mpk, spk)
	sign, err := signer.Sign(proposal.Bytes())
	if err != nil {
		return nil, err
	}
	return &SignedTx{
		Proposal:  proposal,
		Signature: sign,
	}, nil
}

// DecodeTransactionFromBytes 反序列化
func DecodeTransactionFromBytes(data []byte) (*SignedTx, error) {
	var tx SignedTx
	if err := json.Unmarshal(data, &tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

// Verify 验证交易
func (s SignedTx) Verify() bool {
	return core.Verify(s.Proposal.Bytes(), s.Signature)
}

// Bytes SignedTx
func (s *SignedTx) Bytes() []byte {
	buf, _ := json.Marshal(s)
	return buf
}
