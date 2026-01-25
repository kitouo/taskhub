# TaskHub

[![Go Version](https://img.shields.io/badge/Go-1.25.5+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

一个基于Go语言开发的轻量级任务管理系统，采用简洁的架构设计，提供RESTful API接口进行任务的创建、查询和状态管理。

## ✨ 功能特性

- 📝 **任务管理** - 创建、查询和更新任务状态
- ✅ **状态跟踪** - 支持任务完成状态切换
- � **轻量级** - 内存存储，快速响应
- � **RESTful API** - 标准HTTP接口
- � **配置灵活** - 环境变量配置
- � **健康检查** - 内置健康检查端点
- �️ **错误处理** - 完善的错误响应机制

## 🏗️ 项目架构

```
taskhub/
├── cmd/                    # 应用程序入口点
│   └── api/               # HTTP API服务启动
│       ├── main.go        # 主程序入口
│       └── main_test.go   # 主程序测试
├── internal/              # 内部应用代码（不对外暴露）
│   ├── api/              # HTTP处理器和路由
│   │   ├── router.go     # 路由配置
│   │   └── task_handler.go # 任务处理器
│   ├── app/              # 应用程序初始化
│   │   └── app.go        # 应用程序结构
│   ├── config/           # 配置管理
│   │   └── config.go     # 配置加载器
│   ├── httpx/            # HTTP工具和中间件
│   │   ├── middleware.go # 中间件
│   │   ├── recover.go    # 恢复中间件
│   │   └── response.go   # 响应工具
│   ├── logx/             # 日志处理
│   │   └── logx.go       # 日志工具
│   ├── model/            # 数据模型定义
│   │   └── task.go       # 任务模型
│   ├── repo/             # 数据访问层
│   │   ├── memory/       # 内存存储实现
│   │   └── task_repo.go  # 任务仓库接口
│   └── service/          # 业务服务层
│       └── task_service.go # 任务业务逻辑
├── api/                   # API文档目录（预留）
├── configs/               # 配置文件目录（预留）
├── deploy/                # 部署文件目录（预留）
├── scripts/               # 脚本目录（预留）
├── Makefile              # 构建命令
└── go.mod                # Go模块依赖
```

## 🚀 快速开始

### 环境要求

- **Go**: 1.25.5 或更高版本
- **内存**: 最少 128MB RAM
- **操作系统**: Linux, macOS, Windows

### 安装和运行

1. **克隆项目**
```bash
git clone https://github.com/kitouo/taskhub.git
cd taskhub
```

2. **安装依赖**
```bash
make tidy
# 或者
go mod tidy
```

3. **配置环境变量（可选）**
```bash
export APP_ENV=dev                    # 运行环境: dev/staging/prod
export HTTP_PORT=8080                 # HTTP端口
export LOG_LEVEL=info                 # 日志级别: debug/info/warn/error
export READ_TIMEOUT_SEC=5             # 读取超时
export WRITE_TIMEOUT_SEC=10           # 写入超时
export IDLE_TIMEOUT_SEC=60            # 空闲超时
export SHUTDOWN_TIMEOUT_SEC=10        # 关闭超时
```

4. **启动服务**
```bash
# 使用 Makefile（推荐）
make run

# 或者直接运行
go run ./cmd/api
```

5. **验证安装**
```bash
# 健康检查
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz

# 获取任务列表
curl http://localhost:8080/tasks
```

### 🐳 Docker 部署

```bash
# 构建镜像
docker build -t taskhub:latest .

# 运行容器
docker run -p 8080:8080 -e LOG_LEVEL=info taskhub:latest
```

## 📖 API文档

### 健康检查端点

- **健康检查**: `GET /healthz` - 返回服务健康状态
- **就绪检查**: `GET /readyz` - 返回服务就绪状态

### 任务管理API

#### 获取任务列表
```http
GET /tasks
```

**响应示例:**
```json
[
  {
    "id": "task-123",
    "title": "完成项目文档",
    "done": false,
    "create_at": "2024-01-20T10:30:00Z"
  }
]
```

#### 创建新任务
```http
POST /tasks
Content-Type: application/json

{
  "title": "新任务标题"
}
```

**响应示例:**
```json
{
  "id": "task-456",
  "title": "新任务标题",
  "done": false,
  "create_at": "2024-01-20T11:00:00Z"
}
```

#### 获取单个任务
```http
GET /tasks/{id}
```

#### 更新任务状态
```http
PATCH /tasks/{id}
Content-Type: application/json

{
  "done": true
}
```

### 错误响应格式

```json
{
  "error": {
    "code": "INVALID_ARGUMENT",
    "message": "title is required (<= 200)",
    "request_id": "req-123"
  }
}
```

## �️ 开发指南

### 本地开发

```bash
# 运行测试
make test

# 代码整理
make tidy

# 启动开发服务器
make run
```

### 代码规范

- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 `gofmt` 格式化代码
- 使用 `golint` 进行代码检查
- 编写单元测试，保持测试覆盖率 > 80%
- 遵循清洁架构原则

### 项目结构说明

- **cmd/**: 应用程序入口，包含main函数
- **internal/**: 私有应用代码，不会被其他项目导入
- **api/**: API定义，包括OpenAPI规范（预留）
- **configs/**: 配置文件和模板（预留）
- **deploy/**: 部署相关文件（Dockerfile, k8s manifests等）（预留）

### Git提交规范

遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**类型说明：**
- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式调整（不影响功能）
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动
- `perf`: 性能优化

**示例：**
```
feat(api): add task priority filtering
fix(auth): resolve token expiration issue
docs: update API documentation
```

## 🧪 测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./internal/service/...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📊 性能监控

项目集成了以下监控功能：

- **健康检查**: `/healthz` 和 `/readyz` 端点
- **日志记录**: 结构化日志输出
- **请求ID**: 每个请求都有唯一标识符用于追踪

## 🔧 配置说明

支持通过环境变量进行配置：

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `APP_ENV` | dev | 运行环境（dev/staging/prod） |
| `HTTP_PORT` | 8080 | HTTP服务端口 |
| `LOG_LEVEL` | info | 日志级别（debug/info/warn/error） |
| `READ_TIMEOUT_SEC` | 5 | 读取超时时间（秒） |
| `WRITE_TIMEOUT_SEC` | 10 | 写入超时时间（秒） |
| `IDLE_TIMEOUT_SEC` | 60 | 空闲超时时间（秒） |
| `SHUTDOWN_TIMEOUT_SEC` | 10 | 优雅关闭超时时间（秒） |

## 🤝 贡献指南

我们欢迎所有形式的贡献！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

### 贡献者

感谢所有为这个项目做出贡献的开发者！

## 📄 许可证

本项目采用 [MIT License](LICENSE) 许可证。

## 📞 联系方式

- **项目地址**: [https://github.com/kitouo/taskhub](https://github.com/kitouo/taskhub)
- **问题反馈**: [Issues](https://github.com/kitouo/taskhub/issues)
- **功能建议**: [Discussions](https://github.com/kitouo/taskhub/discussions)

---

⭐ 如果这个项目对你有帮助，请给我们一个星标！