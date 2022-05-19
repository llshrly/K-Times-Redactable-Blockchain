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
	"fmt"
	"testing"
	"time"
)

func Test_ChamelemonHash(t *testing.T) {
	// Generate the parameters.
	var p, q, g, hk, tk, hash1, hash2, r1, s1, r2, s2, msg1, msg2 []byte

	Keygen(128, &p, &q, &g, &hk, &tk)

	msg1 = []byte("YES")
	msg2 = []byte("NO")

	r1 = Randgen(&q)
	s1 = Randgen(&q)

	fmt.Printf("CHAMELEON HASH PARAMETERS:"+
		"\np: %s1"+
		"\nq: %s1"+
		"\ng: %s1"+
		"\nhk: %s1"+
		"\ntk: %s1"+
		"\nDONE!\n", p, q, g, hk, tk)

	fmt.Println("====Time: ", time.Now().Format(time.RFC3339Nano))
	// First we generate a chameleon chamelemon.
	ChameleonHash(&hk, &p, &q, &g, &msg1, &r1, &s1, &hash1)
	fmt.Println("====Time: ", time.Now().Format(time.RFC3339Nano))

	fmt.Printf("\n\nROUND 1:"+
		"\nmsg1: %s"+
		"\nr1: %s1"+
		"\ns1: %s1"+
		"\nhash1: %x\n",
		msg1, r1, s1, hash1)

	fmt.Printf("\n\nGENERATING COLLISION...\n\n")

	// Now we need to generate a collision.
	GenerateCollision(&hk, &tk, &p, &q, &g, &msg1, &msg2, &r1, &s1, &r2, &s2)

	ChameleonHash(&hk, &p, &q, &g, &msg2, &r2, &s2, &hash2)

	fmt.Printf("\nROUND 2:"+
		"\nmsg2: %s"+
		"\nr2: %s"+
		"\ns2: %s"+
		"\nhash2: %x\n",
		msg2, r2, s2, hash2)
}
