/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Description:
 * @File:
 * @Version: 1.0.0
 * @Date: 2022.05.29
 */

package chamelemon

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	msg1 := []byte("Hello BTC!")
	r1, s1, hash1 := Hash128(msg1)
	fmt.Printf("verify hash: %t\n", VerifyHash(msg1, r1, s1, hash1))

	msg2 := []byte("Destroy Luna!")
	r2, s2, hash2 := GenerateCollision128(msg1, msg2, r1, s1)

	fmt.Printf("\n\nROUND 1:"+
		"\nmsg1: %s"+
		"\nr1: %s"+
		"\ns1: %s"+
		"\nhash: %x\n",
		msg1, r1, s1, hash1)

	fmt.Printf("\nROUND 2:"+
		"\nmsg2: %s"+
		"\nr2: %s"+
		"\ns2: %s"+
		"\nhash2: %x\n",
		msg2, r2, s2, hash2)
}
