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
	"github.com/llshrly/K-Times-Redactable-Blockchain/common/chamelemon/core"
)

// 变色龙hash系统参数，随机生成
var (
	p128  = []byte("f89d0d3e104c11a81d07065b12e67381")
	q128  = []byte("7c4e869f082608d40e83832d897339c0")
	g128  = []byte("d80a2fcb71db7853db62d160df4a8a66")
	hk128 = []byte("880aa9e911c8a70f8fba5651fca7bee5")
	tk128 = []byte("5b1186069dc48ec5b94be26962263255")
)

//func init() {
//	core.Keygen(128, &p128, &q128, &g128, &hk128, &tk128)
//
//	fmt.Printf("CHAMELEON HASH PARAMETERS:"+
//		"\np: %s1"+
//		"\nq: %s1"+
//		"\ng: %s1"+
//		"\nhk: %s1"+
//		"\ntk: %s1"+
//		"\nDONE!", p128, q128, g128, hk128, tk128)
//}

// Hash128 哈希
func Hash128(msg []byte) (r, s, hash []byte) {
	r = core.Randgen(&q128)
	s = core.Randgen(&q128)

	core.ChameleonHash(&hk128, &p128, &q128, &g128, &msg, &r, &s, &hash)
	return
}

// VerifyHash 验证哈希
func VerifyHash(msg []byte, r, s, hash1 []byte) bool {
	var hash2 []byte
	core.ChameleonHash(&hk128, &p128, &q128, &g128, &msg, &r, &s, &hash2)
	if string(hash1) == string(hash2) {
		return true
	}
	return false
}

// GenerateCollision128 计算hash碰撞
func GenerateCollision128(msg1, msg2, r1, s1 []byte) (r2, s2, hash2 []byte) {
	core.GenerateCollision(&hk128, &tk128, &p128, &q128, &g128, &msg1, &msg2, &r1, &s1, &r2, &s2)
	core.ChameleonHash(&hk128, &p128, &q128, &g128, &msg2, &r2, &s2, &hash2)
	return
}
