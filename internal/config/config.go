package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config 配置结构体
type Config struct {
	AppEnv string // 运行环境 dev/staging/prod

	HTTPPort string // 端口

	/**
	日志级别 debug/info/warn/error
	debug：开发调试细节（变量、分支、请求参数摘要等）
	info：正常业务流程（启动、请求完成、关键状态变化）
	warn：不致命但异常/可疑（重试、降级、慢查询、第三方超时但已兜底）
	error：明确失败（请求失败、依赖不可用、任务处理失败）
	**/
	LogLevel string

	ReadTimeoutSec     int
	WriteTimeoutSec    int
	IdleTimeoutSec     int
	ShutdownTimeoutSec int
}

// Load 加载器
func Load() (Config, error) {
	// 创建配置文件对象
	cfg := Config{
		AppEnv:   getenv("APP_ENV", "dev"),
		HTTPPort: getenv("HTTP_PORT", "8080"),
		LogLevel: strings.ToLower(getenv("LOG_LEVEL", "info")),

		ReadTimeoutSec:     getenvInt("READ_TIMEOUT_SEC", 5),
		WriteTimeoutSec:    getenvInt("WRITE_TIMEOUT_SEC", 10),
		IdleTimeoutSec:     getenvInt("IDLE_TIMEOUT_SEC", 60),
		ShutdownTimeoutSec: getenvInt("SHUTDOWN_TIMEOUT_SEC", 10),
	}

	if cfg.HTTPPort == "" {
		return Config{}, fmt.Errorf("HTTPPort must not be empty")
	}

	// 选择日志等级
	switch cfg.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return Config{}, fmt.Errorf("invalid LOG_LEVEL: %s", cfg.LogLevel)
	}

	return cfg, nil
}

func (c Config) SafeString() string {
	return fmt.Sprintf("app_env: %s, http_port: %s, level: %s", c.AppEnv, c.HTTPPort, c.LogLevel)
}

// 读取环境变量
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}

	return n

}
