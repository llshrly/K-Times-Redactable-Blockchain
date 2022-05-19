/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Description:
 * @File:
 * @Version: 1.0.0
 * @Date: 2022.05.29
 */

package core

import "encoding/json"

type signature struct {
	R   []byte    `json:"r"`
	S   []byte    `json:"s"`
	Pub PublicKey `json:"pub"`
}

func NewSignature(r []byte, s []byte, pub PublicKey) *signature {
	return &signature{R: r, S: s, Pub: pub}
}

// DecodeSignatureFromBytes 反序列化
func DecodeSignatureFromBytes(data []byte) *signature {
	var s signature
	json.Unmarshal(data, &s)
	return &s
}

// Bytes 序列化
func (s signature) Bytes() []byte {
	buf, _ := json.Marshal(s)
	return buf
}
