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
	"math/big"
	mathrand "math/rand"
	"time"
)

func GenerateKeyVC(p, q, g *big.Int) (*PrivateKey, error) {
	randSource := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	// choose random integer x from {1...(q-1)}
	priv := new(big.Int).Rand(randSource, new(big.Int).Sub(q, one))
	// y = g^p mod p
	y := new(big.Int).Exp(g, priv, p)

	return &PrivateKey{
		PublicKey: PublicKey{
			G: g, // cyclic group generator Zp
			P: p, // prime number
			Y: y, // y = g^p mod p
		},
		X: priv, // secret key x
	}, nil
}
