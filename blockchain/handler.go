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
	"github.com/llshrly/K-Times-Redactable-Blockchain/blockchain/transaction"
	"github.com/llshrly/K-Times-Redactable-Blockchain/common/comm"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// Get请求的 Handler
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	log.Info("=======handleGetBlockchain==========")
	bytes, err := json.MarshalIndent(chain.Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
	log.Info("=================================")
}

// POST请求的 payload
type Message struct {
	Tx []byte `json:"tx"`
}

// Post请求的 Handler
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	log.Info("=======handleWriteBlock==========")
	var m Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		fmt.Println(err.Error())
		comm.RespondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	// 解析交易
	tx, err := transaction.DecodeTransactionFromBytes(m.Tx)
	if err != nil {
		comm.RespondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}

	// 生成交易
	newBlock, err := chain.GenerateBlock(chain.Blockchain[len(chain.Blockchain)-1], tx)
	if err != nil {
		comm.RespondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if chain.IsBlockValid(newBlock, chain.Blockchain[len(chain.Blockchain)-1]) {
		newBlockchain := append(chain.Blockchain, newBlock)
		chain.ReplaceChain(newBlockchain)
		spew.Dump(chain.Blockchain)
	}

	comm.RespondWithJSON(w, r, http.StatusCreated, newBlock)
	log.Info("=================================")
}
