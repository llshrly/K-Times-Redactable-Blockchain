/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Description:
 * @File:
 * @Version: 1.0.0
 * @Date: 2022.05.29
 */

package comm

import (
	"encoding/json"
	"log"
	"net/http"
)

// POST处理完成后，返回给客户端一个响应
func RespondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	log.Println("response: ", string(response))
	w.WriteHeader(code)
	w.Write(response)
}

// POST处理完成后，返回给客户端一个响应
func Respond(w http.ResponseWriter, r *http.Request, code int, response []byte) {
	log.Println("response: ", string(response))
	w.WriteHeader(code)
	w.Write(response)
}
