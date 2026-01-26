.PHONY: run test tidy compose-up compose-down compose-logs compose-shell

PORT ?= 8080
READ_TIMEOUT_SEC ?= 5
WRITE_TIMEOUT_SEC ?= 10
IDLE_TIMEOUT_SEC ?= 60
SHUTDOWN_TIMEOUT_SEC ?= 10


run:
	PORT=$(PORT) READ_TIMEOUT_SEC=$(READ_TIMEOUT_SEC) WRITE_TIMEOUT_SEC=$(WRITE_TIMEOUT_SEC) IDLE_TIMEOUT_SEC=$(IDLE_TIMEOUT_SEC) SHUTDOWN_TIMEOUT_SEC=$(SHUTDOWN_TIMEOUT_SEC) go run ./cmd/api

test:
	go test ./... -count=1

tidy:
	go mod tidy

# 一键构建并后台启动（MySQL + API）
compose-up:
	docker compose -f deploy/docker-compose.yaml up -d --build

# 停止并删除容器；-v 会删除 volume（会清空 MySQL 数据）
compose-down:
	docker compose -f deploy/docker-compose.yaml down -v --remove-orphans

# 跟随日志（默认显示最近 200 行）
compose-logs:
	docker compose -f deploy/docker-compose.yaml logs -f --tail=200

# 进入api容器（Alpine有sh，可直接排障）
compose-shell:
	docker compose -f deploy/docker-compose.yaml exec api sh
