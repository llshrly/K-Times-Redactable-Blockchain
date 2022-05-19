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
	"fmt"
	"github.com/llshrly/K-Times-Redactable-Blockchain/common/vc"
)

func main() {
	vcp := &vc.VectorCommitParam{}

	vcp.KenGen(3)
	c := vcp.Compute()
	Λi := vcp.OpenPP(0)

	fmt.Println("string c: ", c.String())

	fmt.Println(vcp.VerifyPP(c, 0, Λi))

	cUpdate := vcp.Update(c, 0)
	ΛUpdate := vcp.ProofUpdate(1, 0, Λi)
	fmt.Println(vcp.VerifyPP(cUpdate, 1, ΛUpdate))

	cUpdate2 := vcp.Update(cUpdate, 1)
	ΛUpdate2 := vcp.ProofUpdate(2, 1, ΛUpdate)
	fmt.Println(vcp.VerifyPP(cUpdate2, 2, ΛUpdate2))
}
