.PHONY: run test tidy

PORT ?= 8080

run:
	PORT=$(PORT) go run ./cmd/api

test:
	go test ./... -count=1

tidy:
	go mod tidy
