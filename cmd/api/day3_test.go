package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kitouo/taskhub/internal/httpx"
	"github.com/kitouo/taskhub/internal/logx"
)

func TestErrorResponseIncludesRequestID(t *testing.T) {
	logger := logx.New("service=api", logx.ParseLevel("error"))

	h := newMux()
	h = httpx.AccessLogger(logger, h)
	h = httpx.Recover(logger, h)
	h = httpx.WithRequestID(h)

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	req.Header.Set("X-Request-Id", "abc123")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var body httpx.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v, body=%s", err, rec.Body.String())
	}
	if body.RequestId != "abc123" {
		t.Fatalf("expected request_id %q, got %q", "abc123", body.RequestId)
	}
	if body.Code != "INVALID_ARGUMENT" {
		t.Fatalf("expected code %q, got %q", "INVALID_ARGUMENT", body.Code)
	}
}

func TestPanicIsRecovered(t *testing.T) {
	logger := logx.New("service=api", logx.ParseLevel("error"))

	h := newMux()
	h = httpx.AccessLogger(logger, h)
	h = httpx.Recover(logger, h)
	h = httpx.WithRequestID(h)

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	req.Header.Set("X-Request-Id", "rid-panic")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	var body httpx.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v, body=%s", err, rec.Body.String())
	}
	if body.RequestId != "rid-panic" {
		t.Fatalf("expected request_id %q, got %q", "rid-panic", body.RequestId)
	}
	if body.Code != "INTERNAL" {
		t.Fatalf("expected code %q, got %q", "INTERNAL", body.Code)
	}
}
