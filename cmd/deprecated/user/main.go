/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Description:
 * @File:
 * @Version: 1.0.0
 * @Date: 2022.05.29
 */

package main

import (
	"crypto/rand"
	"fmt"
	"github.com/Nik-U/pbc"
	"io"
	"math/big"
	"os"
)

func randomK(r io.Reader) (k *big.Int, err error) {
	for {
		limit, _ := new(big.Int).SetString("78791937", 10)
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
		os.Exit(0)
	}
	return pairing
}

func main() {
	pairing := testPairing()

	/* key gen */
	g := pairing.NewG1()
	z1, _ := randomK(rand.Reader)
	z2, _ := randomK(rand.Reader)

	fmt.Println("key gen")
	fmt.Println(g)
	h1 := pairing.NewG1().PowBig(g, z1)
	h2 := pairing.NewG1().PowBig(g, z2)

	m1, _ := randomK(rand.Reader)
	m2, _ := randomK(rand.Reader)

	/* Com */
	fmt.Println("Com")
	h1m1 := pairing.NewG1().PowBig(h1, m1)
	h2m2 := pairing.NewG1().PowBig(h2, m2)
	c := pairing.NewG1().Mul(h1m1, h2m2)

	/* Open */
	fmt.Println("Open")
	Λ1 := pairing.NewG1().PowBig(h2m2, z1)

	/* Verify */
	fmt.Println("Verify")
	cDivh1m1 := pairing.NewG1().Div(c, h1m1)

	t10 := pairing.NewGT()
	t1 := t10.Pair(cDivh1m1, h1)
	t2 := t10.Pair(Λ1, g)
	fmt.Println(t1.Equals(t2))

	/* Update */
	mSub, _ := randomK(rand.Reader)
	h1Sub := pairing.NewG1().PowBig(h1, big.NewInt(0).Sub(mSub, m1))
	cSub := pairing.NewG1().Mul(c, h1Sub)
	cSubDivh2m2 := pairing.NewG1().Div(cSub, h2m2)

	/* Proof Update */
	Λ2 := pairing.NewG1().PowBig(h1m1, z2)
	h1Subz2 := pairing.NewG1().PowBig(h1Sub, z2)
	Λ2Sub := pairing.NewG1().Mul(Λ2, h1Subz2)

	/* Verify */
	t101 := pairing.NewGT()
	t11 := t101.Pair(cSubDivh2m2, h2)
	t21 := t101.Pair(Λ2Sub, g)
	fmt.Println(t11.Equals(t21))
}
