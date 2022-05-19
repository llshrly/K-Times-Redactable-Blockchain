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
	"github.com/llshrly/K-Times-Redactable-Blockchain/common/vc"
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
func handleEdit(w http.ResponseWriter, r *http.Request) {
	log.Info("=======EditTx==========")
	log.Infof("At %s\n", time.Now().Format(time.RFC3339Nano))
	// 解析请求
	var m editRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		fmt.Println(err.Error())
		comm.RespondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	fmt.Printf("req %+v\n", m)

	// 校验
	vcp := vc.Query(m.Editor)
	c := vc.NewZero().SetBytes(m.C)
	proof := vc.NewZero().SetBytes(m.Proof)
	if !vcp.Check(m.C, m.Proof) {
		comm.RespondWithJSON(w, r, http.StatusOK, "proof already been used or invalid or revoked")
		return
	}

	if !vcp.VerifyPP(c, m.I, proof) {
		comm.RespondWithJSON(w, r, http.StatusOK, "not authorized")
		return
	}

	// 修改区块
	{
		reqEditBlock := struct {
			Index   int    `json:"index"`
			Message string `json:"message"`
			Mpk     string `json:"mpk"`
			Spk     string `json:"spk"`
		}{
			Index:   m.Index,
			Message: m.Message,
			Mpk:     os.Getenv("SUPERVISOR_PUBLIC"),
			Spk:     m.Editor,
		}
		data, _ := json.Marshal(reqEditBlock)
		client := resty.New()
		_, err := client.R().SetHeader("Content-Type", "application/json").SetBody(data).Post(fmt.Sprintf("%s/edit", os.Getenv("CHAIN_ENDPOINT")))
		if err != nil {
			comm.RespondWithJSON(w, r, http.StatusOK, "edit block failed")
			return
		}
	}

	// update
	cUpdate := vcp.Update(c, m.I)
	ΛUpdate := vcp.ProofUpdate(m.I+1, m.I, proof)
	fmt.Println(vcp.VerifyPP(cUpdate, m.I+1, ΛUpdate))
	if !vcp.VerifyPP(cUpdate, m.I+1, ΛUpdate) {
		comm.RespondWithJSON(w, r, http.StatusOK, "no edit time rest")
		return
	}
	vc.Update(m.Editor, cUpdate.Bytes(), ΛUpdate.Bytes())

	result := enrollResponse{
		Q:     m.Q,
		C:     cUpdate.Bytes(),
		Proof: ΛUpdate.Bytes(),
		I:     m.I + 1,
	}
	resp, _ := json.Marshal(result)
	comm.Respond(w, r, http.StatusCreated, resp)

	log.Infof("At %s\n", time.Now().Format(time.RFC3339Nano))
	log.Info("=======EditTx END==========")
	defer r.Body.Close()
}

// enrollRequest 请求
type enrollRequest struct {
	Index  int    `json:"index"`
	Count  int    `json:"count"`
	Editor string `json:"editor"`
}

// enrollResponse 返回
type enrollResponse struct {
	Q     int    `json:"q"`
	C     []byte `json:"c"`
	Proof []byte `json:"proof"`
	I     int    `json:"i"`
}

// Post请求的 Handler
func handleEnroll(w http.ResponseWriter, r *http.Request) {
	log.Info("=======EditTimesSet==========")
	log.Infof("At %s\n", time.Now().Format(time.RFC3339Nano))
	var m enrollRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		fmt.Println(err.Error())
		comm.RespondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	fmt.Printf("req %+v\n", m)

	/* vc gen */
	vcp := &vc.VectorCommitParam{}
	vcp.KenGen(m.Count)
	if err := vc.Insert(m.Editor, vcp); err != nil {
		comm.RespondWithJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	c := vcp.Compute()
	Λi := vcp.OpenPP(0)
	vc.Update(m.Editor, c.Bytes(), Λi.Bytes())

	result := enrollResponse{
		Q:     m.Count,
		C:     c.Bytes(),
		Proof: Λi.Bytes(),
		I:     0,
	}
	resp, _ := json.Marshal(result)
	comm.Respond(w, r, http.StatusCreated, resp)
	defer r.Body.Close()
	log.Infof("At %s\n", time.Now().Format(time.RFC3339Nano))
	log.Info("=======EditTimesSet END==========")
}

// revokeRequest 请求
type revokeRequest struct {
	Editor string `json:"editor"`
}

// Post请求的 Handler
func handleRevoke(w http.ResponseWriter, r *http.Request) {
	log.Info("=======REVOKE==========")
	log.Infof("At %s\n", time.Now().Format(time.RFC3339Nano))
	var m revokeRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		fmt.Println(err.Error())
		comm.RespondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	fmt.Printf("req %+v\n", m)

	/* vc gen */
	vc.Delete(m.Editor)

	comm.RespondWithJSON(w, r, http.StatusCreated, "ok")
	defer r.Body.Close()
	log.Infof("At %s\n", time.Now().Format(time.RFC3339Nano))
	log.Info("=======Revoke END==========")
}
