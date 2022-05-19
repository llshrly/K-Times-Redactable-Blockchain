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
	"errors"
)

var vcMap = make(map[string]VectorCommitParam)

func Insert(id string, param *VectorCommitParam) error {
	if _, ok := vcMap[id]; ok {
		return errors.New("vc id already existed")
	}
	vcMap[id] = *param
	return nil
}

func Update(id string, c, proof []byte) error {
	if _, ok := vcMap[id]; !ok {
		return errors.New("vc id not exist")
	}
	vc := vcMap[id]
	vc.SetCurrent(c, proof)
	vcMap[id] = vc
	return nil
}

func Query(id string) VectorCommitParam {
	return vcMap[id]
}

func Delete(id string) {
	delete(vcMap, id)
}
