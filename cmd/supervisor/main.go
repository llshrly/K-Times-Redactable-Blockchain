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
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

const port = "9002"

/**
 * 后面是创建HTTP请求的部分了
 */

// Web服务 使用Gorilla/mux 包
func run() error {
	mux := makeMuxRouter()
	httpAddr := port
	log.Println("Listening on ", httpAddr)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// 对于HTTP服务器的 Handler
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	// 颁发编辑次数凭证
	muxRouter.HandleFunc("/enroll", handleEnroll).Methods("POST")
	// 申请编辑
	muxRouter.HandleFunc("/edit", handleEdit).Methods("POST")
	// 申请编辑
	muxRouter.HandleFunc("/revoke", handleRevoke).Methods("POST")
	return muxRouter
}

// main 函数
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(run())
}
