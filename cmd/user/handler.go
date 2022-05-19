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
	"github.com/go-resty/resty/v2"
	"github.com/llshrly/K-Times-Redactable-Blockchain/blockchain/core"
	"github.com/llshrly/K-Times-Redactable-Blockchain/blockchain/transaction"
	"github.com/llshrly/K-Times-Redactable-Blockchain/common/comm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

const userKey = "user.key"

type txRequest struct {
	Message string `json:"message"`
	Mpk     string `json:"mpk"`
	Spk     string `json:"spk"`
}

// Post请求的 Handler
func handleSendTx(w http.ResponseWriter, r *http.Request) {
	log.Info("=======SendTx==========")
	log.Infof("At %s\n", time.Now().Format(time.RFC3339Nano))
	// 解析请求
	var m txRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		fmt.Println(err.Error())
		comm.RespondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}

	// 生成交易
	txBytes, err := func() ([]byte, error) {
		signer, err := core.Load(userKey)
		if err != nil {
			return nil, err
		}

		got, err := transaction.NewSignedTx(*signer, m.Message, m.Mpk, m.Spk)
		if err != nil {
			return nil, err
		}
		return got.Bytes(), nil
	}()
	if err != nil {
		comm.RespondWithJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// 发送节点上链
	req := struct {
		Tx []byte `json:"tx"`
	}{
		Tx: txBytes,
	}
	data, _ := json.Marshal(req)
	client := resty.New()
	resp, _ := client.R().SetHeader("Content-Type", "application/json").SetBody(data).Post(os.Getenv("CHAIN_ENDPOINT"))
	comm.Respond(w, r, http.StatusCreated, resp.Body())

	log.Infof("At %s\n", time.Now().Format(time.RFC3339Nano))
	log.Info("=======SendTx End==========")
	defer r.Body.Close()
}
