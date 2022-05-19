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
	"encoding/base64"
	"fmt"
	"github.com/Nik-U/pbc"
	log "github.com/sirupsen/logrus"
	"math/big"
	"time"
)

type Params struct {
	pair    *pbc.Pairing
	g       *pbc.Element
	mUpdate *big.Int
}

// dp 系统性参数
var dp *Params

func NewZero() *pbc.Element {
	return dp.pair.NewG1()
}

func init() {
	log.Info("=======keyGen==========")
	log.Infof("At %s\n", time.Now().Format(time.RFC3339Nano))
	pairing := testPairing()
	g := pairing.NewG1().Rand()
	fmt.Println("g: ", g)

	mUpdate, _ := randomK(rand.Reader)
	dp = &Params{
		pair:    pairing,
		g:       g,
		mUpdate: mUpdate,
	}
	log.Infof("At %s\n", time.Now().Format(time.RFC3339Nano))
	log.Info("=======keyGen End==========")
}

type VectorCommitParam struct {
	q     int
	zList []*big.Int
	hList []*pbc.Element
	aux   []*big.Int

	cCurrent     []byte
	proofCurrent []byte
}

func (p *VectorCommitParam) SetCurrent(c, proof []byte) {
	p.cCurrent = c
	p.proofCurrent = proof
}

// KenGen 初始化
func (p *VectorCommitParam) KenGen(q int) {
	// listing
	zList := make([]*big.Int, 0)
	hList := make([]*pbc.Element, 0)
	aux := make([]*big.Int, 0)
	for i := 0; i < q; i++ {
		// z1,z2 ... zq
		zTmp, _ := randomK(rand.Reader)
		zList = append(zList, zTmp)
		fmt.Println("zTmp: ", zTmp)

		// h1,h2 ... hq
		hTmp := dp.pair.NewG1().PowBig(dp.g, zTmp)
		hList = append(hList, hTmp)
		fmt.Println("hTmp: ", hTmp)

		// m1,m2 ... mq
		mTmp, _ := randomK(rand.Reader)
		aux = append(aux, mTmp)
	}

	// init
	p.q = q
	p.zList = zList
	p.hList = hList
	p.aux = aux
}

// Compute 计算
func (p *VectorCommitParam) Compute() *pbc.Element {
	c := p.hList[0]
	for i := 0; i < len(p.hList); i++ {
		if i != 0 {
			c = dp.pair.NewG1().Mul(c, p.hList[i])
		}
	}
	return c
}

func (p *VectorCommitParam) OpenPP(i int) *pbc.Element {
	var product *pbc.Element
	firstTag := true
	for j := 0; j < p.q; j++ {
		if j != i {
			hjExpmj := dp.pair.NewG1().PowBig(p.hList[j], p.aux[j])
			if firstTag {
				product = hjExpmj
				firstTag = false
			} else {
				product = dp.pair.NewG1().Mul(product, hjExpmj)
			}

		}
	}
	Λi := dp.pair.NewG1().PowBig(product, p.zList[i])
	return Λi
}

func (p *VectorCommitParam) VerifyPP(c *pbc.Element, i int, Λi *pbc.Element) bool {
	if i >= len(p.aux) {
		return false
	}
	hiExpmi := dp.pair.NewG1().PowBig(p.hList[i], p.aux[i])
	cDivhimi := dp.pair.NewG1().Div(c, hiExpmi)

	t10 := dp.pair.NewGT()
	t1 := t10.Pair(cDivhimi, p.hList[i])
	t2 := t10.Pair(Λi, dp.g)
	return t1.Equals(t2)
}

func (p *VectorCommitParam) Update(c *pbc.Element, i int) *pbc.Element {
	mUpdateSubmi := big.NewInt(0).Sub(dp.mUpdate, p.aux[i])
	hiExpmUpdateSubmi := dp.pair.NewG1().PowBig(p.hList[i], mUpdateSubmi)

	cUpdate := dp.pair.NewG1().Mul(c, hiExpmUpdateSubmi)
	return cUpdate
}

func (p *VectorCommitParam) ProofUpdate(j, i int, Λj *pbc.Element) *pbc.Element {
	if j >= len(p.aux) {
		return nil
	}
	mUpdateSubmi := big.NewInt(0).Sub(dp.mUpdate, p.aux[j])
	hiExpmUpdateSubmi := dp.pair.NewG1().PowBig(p.hList[i], mUpdateSubmi)

	hiExpmUpdateSubmiExpZj := dp.pair.NewG1().PowBig(hiExpmUpdateSubmi, p.zList[j])
	return dp.pair.NewG1().Mul(Λj, hiExpmUpdateSubmiExpZj)
}

func (p *VectorCommitParam) Check(c, proof []byte) bool {
	if p == nil {
		return false
	}
	fmt.Println(base64.StdEncoding.EncodeToString(p.cCurrent))
	fmt.Println(base64.StdEncoding.EncodeToString(p.proofCurrent))
	fmt.Println(base64.StdEncoding.EncodeToString(c))
	fmt.Println(base64.StdEncoding.EncodeToString(proof))

	if string(p.cCurrent) == string(c) && string(p.proofCurrent) == string(proof) {
		return true
	}
	return false
}
