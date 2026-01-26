# ==========================================================
# Stage 1: builder
# - 使用 Go 官方镜像编译二进制
# - 产出一个静态二进制，便于在更小更安全的运行时镜像中执行
# ==========================================================
FROM golang:1.25-alpine AS builder

# CA证书用于 o mod下载依赖（HTTPS）
RUN apk add --no-cache ca-certificates

# 工作目录：后续所有COPY/BUILD都在这里完成
WORKDIR /src

# 先复制 go.mod / go.sum
# 好处：只要依赖不变，这一层可以被Docker缓存复用，构建更快
COPY go.mod go.sum ./
RUN go mod download

# 再复制剩余源代码
COPY . .

# 编译二进制：
# - CGO_ENABLED=0：生成静态二进制，运行时不依赖libc（部署更简单）
# - GOOS=linux：容器运行在Linux
# - GOARCH=amd64：大多数云环境/本机 Intel都是amd64
#   （如果你是Apple Silicon，后面可改成多架构构建）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/taskhub-api ./cmd/api

# ==========================================================
# Stage 2: runtime（运行阶段：Alpine，可进入容器排障）
# - 有/bin/sh，可docker exec 进入
# - 安装curl：用于容器内healthcheck与排障
# ==========================================================
FROM alpine:3.23

# 安装：
# - ca-certificates：TLS 请求需要
# - curl：healthcheck / 排障必备
# - tzdata：时间相关（可选，但建议装，避免容器时区混乱）
RUN apk add --no-cache ca-certificates curl tzdata

WORKDIR /app

# 从 builder 阶段拷贝编译产物
COPY --from=builder /out/taskhub-api /app/taskhub-api

# 使用非root用户运行（更接近生产）
# -S 创建系统用户/组；-G 指定组
RUN addgroup -S app && adduser -S app -G app
USER app

EXPOSE 8080

# 启动入口
ENTRYPOINT ["/app/taskhub-api"]

