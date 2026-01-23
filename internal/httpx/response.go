package httpx

import (
	"encoding/json"
	"net/http"
)

/*
统一 json 响应与错误类型
*/

type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
}

func WriteJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, code, msg string, rid string) {
	WriteJson(w, status, ErrorResponse{
		Code:      code,
		Message:   msg,
		RequestId: rid,
	})
}
