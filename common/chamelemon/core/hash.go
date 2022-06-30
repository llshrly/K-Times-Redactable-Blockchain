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
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// Returns a random hex number within the bounds of 0 and upperBoundHex.
func Randgen(upperBoundHex *[]byte) []byte {
	upperBoundBig := new(big.Int)
	upperBoundBig, success := upperBoundBig.SetString(string(*upperBoundHex), 16)
	if success != true {
		fmt.Printf("Conversion from hex: %s to bigInt failed.", upperBoundHex)
	}

	randomBig, err := rand.Int(rand.Reader, upperBoundBig)
	if err != nil {
		fmt.Printf("Generation of random bigInt in bounds [0...%v] failed.", upperBoundBig)
	}

	return []byte(fmt.Sprintf("%x", randomBig))
}

func Keygen(bits int, p *[]byte, q *[]byte, g *[]byte, hk *[]byte, tk *[]byte) {
	gBig := new(big.Int)
	qBig := new(big.Int)
	hkBig := new(big.Int)
	tkBig := new(big.Int)
	oneBig := new(big.Int)
	twoBig := new(big.Int)

	oneBig.SetInt64(1) // oneBig = 1
	twoBig.SetInt64(2) // twoBig = 2

	pBig, err := rand.Prime(rand.Reader, bits) // pBig is a random prime of length bits
	if err != nil {
		fmt.Printf("Generation of random prime number failed.")
	}
	qBig.Sub(pBig, oneBig) // qBig = pBig - 1
	qBig.Div(qBig, twoBig) // qBig = (pBig - 1) / 2

	gBig, err = rand.Int(rand.Reader, pBig)
	if err != nil {
		fmt.Printf("Generation of random bigInt in bounds [0...%v] failed.", pBig)
	}

	gBig.Exp(gBig, twoBig, pBig) // gBig = gBig ^ 2 % pBig

	// Choosing hk and tk
	tkBig, err = rand.Int(rand.Reader, qBig)
	if err != nil {
		fmt.Printf("Generation of random bigInt in bounds [0...%v] failed.", qBig)
	}

	hkBig.Exp(gBig, tkBig, pBig) // hkBig = gBig ^ tkBig % pBig

	*p = []byte(fmt.Sprintf("%x", pBig))
	*q = []byte(fmt.Sprintf("%x", qBig))
	*g = []byte(fmt.Sprintf("%x", gBig))
	*hk = []byte(fmt.Sprintf("%x", hkBig))
	*tk = []byte(fmt.Sprintf("%x", tkBig))
}

func ChameleonHash(
	hk *[]byte,
	p *[]byte,
	q *[]byte,
	g *[]byte,
	message *[]byte,
	r *[]byte,
	s *[]byte,
	hashOut *[]byte,
) {
	hkeBig := new(big.Int)
	gsBig := new(big.Int)
	tmpBig := new(big.Int)
	eBig := new(big.Int)
	pBig := new(big.Int)
	qBig := new(big.Int)
	gBig := new(big.Int)
	rBig := new(big.Int)
	sBig := new(big.Int)
	hkBig := new(big.Int)
	hBig := new(big.Int)

	// Converting from hex to bigInt
	pBig.SetString(string(*p), 16)
	qBig.SetString(string(*q), 16)
	gBig.SetString(string(*g), 16)
	hkBig.SetString(string(*hk), 16)
	rBig.SetString(string(*r), 16)
	sBig.SetString(string(*s), 16)

	// Generate the hashOut with message || rBig
	hash := sha256.New()
	hash.Write([]byte(*message))
	hash.Write([]byte(fmt.Sprintf("%x", rBig)))

	eBig.SetBytes(hash.Sum(nil))

	hkeBig.Exp(hkBig, eBig, pBig)
	gsBig.Exp(gBig, sBig, pBig)
	tmpBig.Mul(hkeBig, gsBig)
	tmpBig.Mod(tmpBig, pBig)
	hBig.Sub(rBig, tmpBig)
	hBig.Mod(hBig, qBig)

	*hashOut = hBig.Bytes() // Return hBig in big endian encoding as string
}

func GenerateCollision(
	hk *[]byte,
	tk *[]byte,
	p *[]byte,
	q *[]byte,
	g *[]byte,
	msg1 *[]byte,
	msg2 *[]byte,
	r1 *[]byte,
	s1 *[]byte,
	r2 *[]byte,
	s2 *[]byte,
) {
	hkBig := new(big.Int)
	tkBig := new(big.Int)
	pBig := new(big.Int)
	qBig := new(big.Int)
	gBig := new(big.Int)
	r1Big := new(big.Int)
	s1Big := new(big.Int)
	kBig := new(big.Int)
	hBig := new(big.Int)
	eBig := new(big.Int)
	tmpBig := new(big.Int)
	r2Big := new(big.Int)
	s2Big := new(big.Int)

	pBig.SetString(string(*p), 16)
	qBig.SetString(string(*q), 16)
	gBig.SetString(string(*g), 16)
	r1Big.SetString(string(*r1), 16)
	s1Big.SetString(string(*s1), 16)
	hkBig.SetString(string(*hk), 16)
	tkBig.SetString(string(*tk), 16)

	// Generate random k
	kBig, err := rand.Int(rand.Reader, qBig)
	if err != nil {
		fmt.Printf("Generation of random bigInt in bounds [0...%v] failed.", qBig)
	}

	// Get chameleon chamelemon of (msg1, r1, s1)
	var hash []byte
	ChameleonHash(hk, p, q, g, msg1, r1, s1, &hash)
	hBig.SetBytes(hash) // Convert the big endian encoded chamelemon into bigInt.

	// Compute the new r1
	tmpBig.Exp(gBig, kBig, pBig)
	r2Big.Add(hBig, tmpBig)
	r2Big.Mod(r2Big, qBig)

	// Compute e'
	newHash := sha256.New()
	newHash.Write([]byte(*msg2))
	newHash.Write([]byte(fmt.Sprintf("%x", r2Big)))
	eBig.SetBytes(newHash.Sum(nil))

	// Compute s2
	tmpBig.Mul(eBig, tkBig)
	tmpBig.Mod(tmpBig, qBig)
	s2Big.Sub(kBig, tmpBig)
	s2Big.Mod(s2Big, qBig)

	*r2 = []byte(fmt.Sprintf("%x", r2Big))
	*s2 = []byte(fmt.Sprintf("%x", s2Big))
}
