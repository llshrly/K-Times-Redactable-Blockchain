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
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/Nik-U/pbc"
	"github.com/cloudflare/bn256"
	"testing"
)

func TestProcess(t *testing.T) {
	vcp := &VectorCommitParam{}

	vcp.KenGen(2)
	c := vcp.Compute()
	Λi := vcp.OpenPP(1)

	fmt.Println(vcp.VerifyPP(c, 1, Λi))
}

func TestTT(t *testing.T) {
	pairing := testPairing()

	/* key gen */
	g := pairing.NewG1()
	z1, _ := randomK(rand.Reader)
	z2, _ := randomK(rand.Reader)
	h1 := new(pbc.Element).PowBig(g, z1)
	h2 := new(pbc.Element).PowBig(g, z2)

	m1, _ := randomK(rand.Reader)
	m2, _ := randomK(rand.Reader)

	/* Com */
	h1m1 := new(pbc.Element).PowBig(h1, m1)
	h2m2 := new(pbc.Element).PowBig(h2, m2)
	c := new(pbc.Element).Mul(h1m1, h2m2)

	/* Open */
	Λ1 := new(pbc.Element).PowBig(h2m2, z1)

	/* Verify */
	cDivh1m1 := new(pbc.Element).Div(c, h1m1)

	t10 := pairing.NewGT()
	t1 := t10.Pair(cDivh1m1, h1)
	t2 := t10.Pair(Λ1, g)
	fmt.Println(t1.Equals(t2))
}

func TestEncryptDecrypt(t *testing.T) {
	// Each of three parties, a, b and c, generate a private value.
	a, _ := rand.Int(rand.Reader, bn256.Order)
	b, _ := rand.Int(rand.Reader, bn256.Order)
	c, _ := rand.Int(rand.Reader, bn256.Order)

	// Then each party calculates g₁ and g₂ times their private value.
	pa := new(bn256.G1).ScalarBaseMult(a)
	qa := new(bn256.G2).ScalarBaseMult(a)

	pb := new(bn256.G1).ScalarBaseMult(b)
	qb := new(bn256.G2).ScalarBaseMult(b)

	pc := new(bn256.G1).ScalarBaseMult(c)
	qc := new(bn256.G2).ScalarBaseMult(c)

	// Now each party exchanges its public values with the other two and
	// all parties can calculate the shared key.
	k1 := bn256.Pair(pb, qc)
	k1.ScalarMult(k1, a)

	k2 := bn256.Pair(pc, qa)
	k2.ScalarMult(k2, b)

	k3 := bn256.Pair(pa, qb)
	k3.ScalarMult(k3, c)

	// k1, k2 and k3 will all be equal.
	fmt.Println(bytes.Equal(k1.Marshal(), k2.Marshal()))

}
