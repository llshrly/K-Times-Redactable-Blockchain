/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Description:
 * @File:
 * @Version: 1.0.0
 * @Date: 2022.05.29
 */

package vc

import (
	"crypto/rand"
	"github.com/Nik-U/pbc"
	log "github.com/sirupsen/logrus"
	"io"
	"math/big"
)

func randomK(r io.Reader) (k *big.Int, err error) {
	for {
		limit, _ := new(big.Int).SetString("6518589491078791937", 10)
		k, err = rand.Int(r, limit)
		if err != nil || k.Sign() > 0 {
			return
		}
	}
	return
}

func testPairing() *pbc.Pairing {
	// Generated with pbc_param_init_a_gen(p, 10, 32);
	pairing, err := pbc.NewPairingFromString("type a\nq 4025338979\nh 6279780\nr 641\nexp2 9\nexp1 7\nsign1 1\nsign0 1\n")
	if err != nil {
		log.Fatalf("Cant init pairing, err: %s\n", err.Error())
	}
	return pairing
}
