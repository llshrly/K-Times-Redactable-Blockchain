/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Description:
 * @File:
 * @Version: 1.0.0
 * @Date: 2022.05.29
 */

package core

import (
	"encoding/base64"
	"encoding/json"
	"github.com/llshrly/K-Times-Redactable-Blockchain/common/elgamal/core"
	"io/ioutil"
	"math/big"
)

type Signer struct {
	priv core.PrivateKey
}

func NewSigner(priv core.PrivateKey) *Signer {
	return &Signer{priv: priv}
}

// Sign 签名
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	r1, s1, err := s.priv.Signature(msg)
	if err != nil {
		return nil, err
	}
	return core.NewSignature(r1, s1, s.priv.PublicKey).Bytes(), nil
}

// Store 存储
func (s *Signer) Store(path string) error {
	buf, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, buf, 0666)
}

// Load 存储
func Load(path string) (*Signer, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := Signer{}
	if err := json.Unmarshal(buf, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Verify 验证签名
func Verify(msg []byte, signBytes []byte) bool {
	sign := core.DecodeSignatureFromBytes(signBytes)
	bool, err := sign.Pub.SigVerify(sign.R, sign.S, msg)
	if err != nil {
		return false
	}
	return bool
}

//MarshalJSON 序列化
func (s *Signer) MarshalJSON() ([]byte, error) {
	buf, _ := json.Marshal(s.priv.PublicKey)
	pubStr := base64.StdEncoding.EncodeToString(buf)

	t := struct {
		Public  string   `json:"public"`
		Private *big.Int `json:"private"`
	}{
		Public:  pubStr,
		Private: s.priv.X,
	}
	return json.Marshal(t)
}

// UnmarshalJSON 反序列化
func (s *Signer) UnmarshalJSON(data []byte) error {
	t := &struct {
		Public  string   `json:"public"`
		Private *big.Int `json:"private"`
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	pubBytes, err := base64.StdEncoding.DecodeString(t.Public)
	if err != nil {
		return err
	}
	pub := core.PublicKey{}
	if err := json.Unmarshal(pubBytes, &pub); err != nil {
		return err
	}
	s.priv = core.PrivateKey{
		PublicKey: pub,
		X:         t.Private,
	}

	return nil
}
