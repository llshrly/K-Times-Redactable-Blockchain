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
	"github.com/llshrly/K-Times-Redactable-Blockchain/common/comm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type editRequest struct {
	Index   int    `json:"index"`
	Message string `json:"message"`
	Editor  string `json:"editor"`

	Q     int    `json:"q"`
	C     []byte `json:"c"`
	Proof []byte `json:"proof"`
	I     int    `json:"i"`
}

// Post请求的 Handler
func handleEditBlock(w http.ResponseWriter, r *http.Request) {
	log.Info("=======EditBlock==========")
	log.Infof("At %s\n", time.Now().Format(time.RFC3339Nano))
	var m editRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		fmt.Println(err.Error())
		comm.RespondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	data, _ := json.Marshal(m)
	client := resty.New()
	resp, _ := client.R().SetHeader("Content-Type", "application/json").SetBody(data).Post(fmt.Sprintf("%s/edit", os.Getenv("SUPERVISOR_ENDPOINT")))
	comm.Respond(w, r, http.StatusCreated, resp.Body())
	defer r.Body.Close()
	log.Infof("At %s\n", time.Now().Format(time.RFC3339Nano))
	log.Info("=======EditBlock END==========")
}

type enrollRequest struct {
	Index  int    `json:"index"`
	Count  int    `json:"count"`
	Editor string `json:"editor"`
}

// Post请求的 Handler
func handleEnrollBlock(w http.ResponseWriter, r *http.Request) {
	log.Info("=======EditTimesSet==========")
	log.Infof("At %s\n", time.Now().Format(time.RFC3339Nano))
	var m enrollRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		fmt.Println(err.Error())
		comm.RespondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	data, _ := json.Marshal(m)
	client := resty.New()
	resp, _ := client.R().SetHeader("Content-Type", "application/json").SetBody(data).Post(fmt.Sprintf("%s/enroll", os.Getenv("SUPERVISOR_ENDPOINT")))
	comm.Respond(w, r, http.StatusCreated, resp.Body())
	defer r.Body.Close()

	log.Infof("At %s\n", time.Now().Format(time.RFC3339Nano))
	log.Info("=======EditTimesSet END==========")
}
