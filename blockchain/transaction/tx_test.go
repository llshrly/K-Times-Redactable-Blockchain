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
	"encoding/base64"
	"fmt"
	"github.com/llshrly/K-Times-Redactable-Blockchain/blockchain/core"
	core2 "github.com/llshrly/K-Times-Redactable-Blockchain/common/elgamal/core"
	"testing"
)

func TestNewSignedTx(t *testing.T) {
	priv, err := core2.GenerateKey(core2.Bit128, core2.Probability)
	if err != nil {
		t.Fatalf("GenerateKey failed with error: %s", priv)
	}
	signer := core.NewSigner(*priv)
	signer.Store("test.key")
	signer2, _ := core.Load("test.key")

	got, err := NewSignedTx(*signer2, "Hello BTC!!", "mpk", "spk")
	if err != nil {
		t.Fatalf("NewSignedTx failed with error: %s", priv)
	}
	fmt.Println(got.Verify())
	fmt.Println("tx base64: ", base64.StdEncoding.EncodeToString(got.Bytes()))

}
