package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthz 健康监测/**
func TestHealthz(t *testing.T) {

	// 创建路由
	mux := http.NewServeMux()

	// 路由注册
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// 请求
	req := httptest.NewRequest("GET", "/healthz", nil)

	// 响应
	resp := httptest.NewRecorder()

	mux.ServeHTTP(resp, req)

	// 是否成功
	if resp.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", resp.Code, http.StatusOK)
	}
}
