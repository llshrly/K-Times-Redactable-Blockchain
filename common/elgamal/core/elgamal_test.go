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
	"fmt"
	"github.com/llshrly/K-Times-Redactable-Blockchain/common/num"
	"testing"
	"time"
)

func Test_verify(t *testing.T) {
	priv, err := GenerateKey(Bit128, Probability)
	if err != nil {
		t.Fatalf("GenerateKey failed with error: %s", priv)
	}

	msg := []byte("Hello BTC!")
	r, s, err := priv.Signature(msg)
	if err != nil {
		t.Fatalf("Signature failed with error: %s", priv)
	}

	match, err := priv.PublicKey.SigVerify(r, s, msg)
	if err != nil {
		t.Fatalf("SigVerify failed with error: %s", priv)
	}
	if match != true {
		t.Fatalf("Signature not match!!!")
	}
	t.Log("Done!")
}

func Test_GenerateKey(t *testing.T) {
	var list []int64
	N := 100
	for i := 0; i < N; i++ {
		start := time.Now()
		GenerateKey(Bit256, Probability)
		list = append(list, time.Now().Sub(start).Milliseconds())

	}
	fmt.Println(list)
	fmt.Println("Max: ", num.Max(list), "Min: ", num.Min(list), "Variance: ", num.Variance(list))
}
