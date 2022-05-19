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
	"errors"
	"math/big"
	"time"

	mathrand "math/rand"
)

const (
	Bit128      = 128
	Bit256      = 256
	Probability = 32
)

var zero = big.NewInt(0)
var one = big.NewInt(1)
var two = big.NewInt(2)

var ErrMessageLarge = errors.New("elgamal: message is larger than public key size")
var ErrCipherLarge = errors.New("elgamal: cipher is larger than public key size")

// PublicKey represents a Elgamal public key.
type PublicKey struct {
	G, P, Y *big.Int
}

// PrivateKey represents Elgamal private key.
type PrivateKey struct {
	PublicKey
	X *big.Int
}

// GenerateKey generates elgamal private key according
// to given bit size and probability. Moreover, the given probability
// value is used in choosing prime number P for performing n Miller-Rabin
// tests with 1 - 1/(4^n) probability false rate.
func GenerateKey(bitsize, probability int) (*PrivateKey, error) {
	// p is prime number
	// q is prime group order
	// g is cyclic group generator Zp
	p, q, g, err := GeneratePQZp(bitsize, probability)
	if err != nil {
		return nil, err
	}

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

func GeneratePQZp(bitsize, probability int) (p, q, g *big.Int, err error) {
	return Gen(bitsize, probability)
}

// Encrypt encrypts a plain text represented as a byte array. It returns
// an error if plain text value is larger than modulus P of Public key.
func (pub *PublicKey) Encrypt(message []byte) ([]byte, []byte, error) {
	// choose random integer k from {1...p}
	k, err := rand.Int(rand.Reader, pub.P)
	if err != nil {
		return nil, nil, err
	}

	m := new(big.Int).SetBytes(message)
	if m.Cmp(pub.P) == 1 { //  m < P
		return nil, nil, ErrMessageLarge
	}

	// c1 = g^k mod p
	c1 := new(big.Int).Exp(pub.G, k, pub.P)
	// s = y^k mod p
	s := new(big.Int).Exp(pub.Y, k, pub.P)
	// c2 = m*s mod p
	c2 := new(big.Int).Mod(
		new(big.Int).Mul(m, s),
		pub.P,
	)
	return c1.Bytes(), c2.Bytes(), nil
}

// Decrypt decrypts the passed cipher text. It returns an
// error if cipher text value is larger than modulus P of Public key.
func (priv *PrivateKey) Decrypt(cipher1, cipher2 []byte) ([]byte, error) {
	c1 := new(big.Int).SetBytes(cipher1)
	c2 := new(big.Int).SetBytes(cipher2)
	if c1.Cmp(priv.P) == 1 && c2.Cmp(priv.P) == 1 { //  (c1, c2) < P
		return nil, ErrCipherLarge
	}

	// s = c^x mod p
	s := new(big.Int).Exp(c1, priv.X, priv.P)
	// s = s(inv) = s^(-1) mod p
	if s.ModInverse(s, priv.P) == nil {
		return nil, errors.New("elgamal: invalid private key")
	}

	// m = s(inv) * c2 mod p
	m := new(big.Int).Mod(
		new(big.Int).Mul(s, c2),
		priv.P,
	)
	return m.Bytes(), nil
}

// HomomorphicEncTwo performs homomorphic operation over two passed chiphers.
// Elgamal has multiplicative homomorphic property, so resultant cipher
// contains the product of two numbers.
func (pub *PublicKey) HomomorphicEncTwo(c1, c2, c1dash, c2dash []byte) ([]byte, []byte, error) {
	cipher1 := new(big.Int).SetBytes(c1)
	cipher2 := new(big.Int).SetBytes(c2)
	if cipher1.Cmp(pub.P) == 1 && cipher2.Cmp(pub.P) == 1 { //  (c1, c2) < P
		return nil, nil, ErrCipherLarge
	}

	// In the context of elgamal encryption, (cipher1,cipher2) and
	// (cipher1dash, cipher2dash) both are valid ciphers and represented
	// by different variable names.
	cipher1dash := new(big.Int).SetBytes(c1dash)
	cipher2dash := new(big.Int).SetBytes(c2dash)
	if cipher1dash.Cmp(pub.P) == 1 && cipher2dash.Cmp(pub.P) == 1 { //  (c1dash, c2dash) < P
		return nil, nil, ErrCipherLarge
	}

	// C1 = c1 * c1dash mod p
	C1 := new(big.Int).Mod(
		new(big.Int).Mul(cipher1, cipher1dash),
		pub.P,
	)

	// C2 = c2 * c2dash mod p
	C2 := new(big.Int).Mod(
		new(big.Int).Mul(cipher2, cipher2dash),
		pub.P,
	)
	return C1.Bytes(), C2.Bytes(), nil
}

// HommorphicEncMultiple performs homomorphic operation over multiple passed chiphers.
// Elgamal has multiplicative homomorphic property, so resultant cipher
// contains the product of multiple numbers.
func (pub *PublicKey) HommorphicEncMultiple(ciphertext [][2][]byte) ([]byte, []byte, error) {
	// C1, C2, _ := pub.Encrypt(one.Bytes())
	C1 := one // since, c = 1^e mod n is equal to 1
	C2 := one

	for i := 0; i < len(ciphertext); i++ {
		c1 := new(big.Int).SetBytes(ciphertext[i][0])
		c2 := new(big.Int).SetBytes(ciphertext[i][1])

		if c1.Cmp(pub.P) == 1 && c2.Cmp(pub.P) == 1 { //  (c1, c2) < P
			return nil, nil, ErrCipherLarge
		}

		// C1 = (c1)_1 * (c1)_2 * (c1)_3 ...(c1)_n mod p
		C1 = new(big.Int).Mod(
			new(big.Int).Mul(
				C1,
				c1),
			pub.P,
		)

		// C2 = (c2)_1 * (c2)_2 * (c2)_3 ...(c2)_n mod p
		C2 = new(big.Int).Mod(
			new(big.Int).Mul(
				C2,
				c2),
			pub.P,
		)
	}
	return C1.Bytes(), C2.Bytes(), nil
}

// Signature generates signature over the given message. It returns signature
// value consisting of two parts "r" and "s" as byte arrays.
func (priv *PrivateKey) Signature(message []byte) ([]byte, []byte, error) {
	k := new(big.Int)
	gcd := new(big.Int)
	var err error

	// choosing random integer k from {1...(p-2)}, such that
	// gcd(k,(p-1)) should be equal to 1.
	for {
		k, err = rand.Int(rand.Reader, new(big.Int).Sub(priv.P, two))
		if err != nil {
			return nil, nil, err
		}
		if k.Cmp(one) == 0 {
			continue
		} else {
			gcd = gcd.GCD(nil, nil, k, new(big.Int).Sub(priv.P, one))
			if gcd.Cmp(one) == 0 {
				break
			}
		}
	}

	// taking SHA256 of the message
	hashofm := sha256.Sum256(message)
	// m as H(m)
	m := new(big.Int).SetBytes(hashofm[:])

	// r = g^k mod p
	r := new(big.Int).Exp(priv.G, k, priv.P)
	// xr = x * r
	xr := new(big.Int).Mod(
		new(big.Int).Mul(r, priv.X),
		new(big.Int).Sub(priv.P, one),
	)

	// hmxr = [H(m) -xr]
	hmxr := new(big.Int).Sub(m, xr)
	// k = k^(-1)
	k = k.ModInverse(k, new(big.Int).Sub(priv.P, one))

	// s = [H(m) -xr]k^(-1) mod (p-1)
	s := new(big.Int).Mod(
		new(big.Int).Mul(hmxr, k),
		new(big.Int).Sub(priv.P, one),
	)
	return r.Bytes(), s.Bytes(), nil
}

// SigVerify verifies signature over the given message and signature values (r & s).
// It returns true as a boolean value if signature is verify correctly. Otherwise
// it returns false along with error message.
func (pub *PublicKey) SigVerify(r, s, message []byte) (bool, error) {
	// verify that 0 < r < p
	signr := new(big.Int).SetBytes(r)
	if signr.Cmp(zero) == -1 {
		return false, errors.New("r is smaller than zero")
	} else if signr.Cmp(pub.P) == +1 {
		return false, errors.New("r is larger than public key p")
	}

	signs := new(big.Int).SetBytes(s)
	if signs.Cmp(zero) == -1 {
		return false, errors.New("s is smaller than zero")
	} else if signs.Cmp(new(big.Int).Sub(pub.P, one)) == +1 {
		return false, errors.New("s is larger than public key p")
	}

	// taking SHA256 of the message
	hashofm := sha256.Sum256(message)
	// m as H(m)
	m := new(big.Int).SetBytes(hashofm[:])
	// ghashm = g^[H(m)] mod p
	ghashm := new(big.Int).Exp(pub.G, m, pub.P)

	// y^r * r*s mod p
	YrRs := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(pub.Y, signr, pub.P),
			new(big.Int).Exp(signr, signs, pub.P),
		),
		pub.P,
	)

	// g^H(m) y^r * r*s mod p
	if ghashm.Cmp(YrRs) == 0 {
		return true, nil // signature is verified
	}
	return false, errors.New("signature is not verified")
}

// Note : this section of code is taken from (https://github.com/ldinc/pqg).
// Author of this code is "Drogunov Igor".
// Gen emit <p,q,g>.
// p = 2q + 1, p,q - safe primes
// g - cyclic group generator Zp
// performs n Miller-Rabin tests with 1 - 1/(4^n) probability false rate.
// Gain n - bit width for integer & probability rang for MR.
// It returns p, q, g and write error message.
func Gen(n, probability int) (*big.Int, *big.Int, *big.Int, error) {
	for {
		q, err := rand.Prime(rand.Reader, n-1)
		if err != nil {
			return nil, nil, nil, err
		}
		t := new(big.Int).Mul(q, two)
		p := new(big.Int).Add(t, one)
		if p.ProbablyPrime(probability) {
			for {
				g, err := rand.Int(rand.Reader, p)
				if err != nil {
					return nil, nil, nil, err
				}
				b := new(big.Int).Exp(g, two, p)
				if b.Cmp(one) == 0 {
					continue
				}
				b = new(big.Int).Exp(g, q, p)
				if b.Cmp(one) == 0 {
					return p, q, g, nil
				}
			}
		}
	}
	return nil, nil, nil, errors.New("can't emit <p,q,g>")
}
