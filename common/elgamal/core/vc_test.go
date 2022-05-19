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
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/cloudflare/bn256"
	"math/big"
	"testing"
)

func TestTT(t *testing.T) {
	p, q, g, err := GeneratePQZp(Bit128, Probability)
	if err != nil {
		t.Fatalf(err.Error())
	}

	priv1, err := GenerateKeyVC(p, q, g)
	priv2, err := GenerateKeyVC(p, q, g)

	/* keygen */
	//===============================================
	z1 := priv1.X
	z2 := priv2.X

	h1 := priv1.Y //g^z1
	h2 := priv2.Y //g^z2

	// h12 = g^z1z2
	h12 := big.NewInt(0).Exp(h1, z2, p)
	h21 := big.NewInt(0).Exp(h2, z1, p)
	fmt.Println(h12, h21)
	// pp
	// pp = {g, h1,h2, h12}

	/* Com */
	//===============================================
	m1, _ := rand.Int(rand.Reader, p)
	m2, _ := rand.Int(rand.Reader, p)
	h1Expm1 := big.NewInt(0).Exp(h1, m1, p)
	h2Expm2 := big.NewInt(0).Exp(h2, m2, p)
	c := big.NewInt(0).Mul(h1Expm1, h2Expm2)
	// aux
	// aux = {c, m1, m2}

	/* Open(m1, 1, aux) */
	//===============================================
	Λ1 := big.NewInt(0).Exp(h12, m1, p)

	/* Verpp(C, m1, 1, Λ1) */
	h1Expm1version := h1Expm1.ModInverse(h1Expm1, p)
	cDivh1m1 := big.NewInt(0).Mul(c, h1Expm1version)
	//
	//cDivh1m1Mulh1 := big.NewInt(0).Mul(cDivh1m1, h1)
	//cDivh1m1Mulh1Mod := cDivh1m1Mulh1.Mod(cDivh1m1Mulh1, p)
	//Λ1g := big.NewInt(0).Mul(Λ1, g)
	//Λ1gMod := Λ1g.Mod(Λ1g, p)

	//fmt.Println(big.NewInt(0).Exp(g, cDivh1m1Mulh1Mod, p))
	//fmt.Println(big.NewInt(0).Exp(g, Λ1gMod, p))
	c11 := new(bn256.G1).ScalarBaseMult(cDivh1m1)
	c12 := new(bn256.G2).ScalarBaseMult(h1)
	t1 := bn256.Pair(c11, c12)

	c21 := new(bn256.G1).ScalarBaseMult(Λ1)
	c22 := new(bn256.G2).ScalarBaseMult(g)
	t2 := bn256.Pair(c21, c22)

	fmt.Println(bytes.Equal(t1.Marshal(), t2.Marshal()))

}
