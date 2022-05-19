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
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/llshrly/K-Times-Redactable-Blockchain/blockchain/chain"
	"github.com/llshrly/K-Times-Redactable-Blockchain/common/comm"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// editMessage 编辑请求
type editMessage struct {
	Index   int    `json:"index"`
	Message string `json:"message"`

	Mpk string `json:"mpk"`
	Spk string `json:"spk"`
}

// handleEditBlock 编辑区块链的请求
func handleEditBlock(w http.ResponseWriter, r *http.Request) {
	log.Info("=======handleGetBlockchain==========")
	var m editMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		fmt.Println(err.Error())
		comm.RespondWithJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	// 替换区块
	if m.Index >= len(chain.Blockchain) {
		comm.RespondWithJSON(w, r, http.StatusInternalServerError, "block index exceed")
		return
	}
	// 验证
	if !verifyEditable(m.Index, m.Mpk, m.Spk) {
		comm.RespondWithJSON(w, r, http.StatusOK, "block edit verify failed, check the mpk & spk")
		return
	}

	//编辑
	newBlock, err := chain.EditBlock(chain.Blockchain[m.Index], m.Message)
	if err != nil {
		comm.RespondWithJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if chain.IsBlockValid(newBlock, chain.Blockchain[len(chain.Blockchain)-1]) {
		newBlockchain := append(chain.Blockchain, newBlock)
		chain.ReplaceChain(newBlockchain)
		spew.Dump(chain.Blockchain)
	}
	chain.ReplaceBlock(newBlock, m.Index)

	comm.RespondWithJSON(w, r, http.StatusCreated, newBlock)
	log.Info("=================================")
}

// verifyEditable 验证
func verifyEditable(index int, mpk, spk string) bool {
	if chain.Blockchain[index].Tx.Mpk != mpk || chain.Blockchain[index].Tx.Spk != spk {
		return false
	}
	return true
}
