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
	"github.com/llshrly/K-Times-Redactable-Blockchain/blockchain/chain"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

/**
 * 后面是创建HTTP请求的部分了
 */

// Web服务 使用Gorilla/mux 包
func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", os.Getenv("ADDR"))
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
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	muxRouter.HandleFunc("/edit", handleEditBlock).Methods("POST")
	return muxRouter
}

// main 函数
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		genesisBlock := chain.GenesisBlock // 创世区块
		spew.Dump(genesisBlock)
		chain.Blockchain = append(chain.Blockchain, genesisBlock)
	}()
	log.Fatal(run())

}
