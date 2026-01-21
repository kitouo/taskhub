.PHONY: run test tidy

run:
	PORT=8080 go run ./cmd/api
test:
	go test ./... -count=1
tidy:
	go mod tidy