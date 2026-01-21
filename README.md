# TaskHub

一个基于Go语言开发的任务管理系统，提供高效的任务创建、分配和跟踪功能。

## 功能特性

- 📝 任务创建和管理
- 👥 团队协作和任务分配
- 📊 任务状态跟踪
- 🔔 通知提醒
- 📈 数据统计和报表

## 项目结构

```
taskhub/
├── api/           # API接口定义
├── cmd/           # 应用程序入口
│   └── api/       # API服务启动
├── configs/       # 配置文件
├── deploy/        # 部署相关文件
├── internal/      # 内部应用代码
│   └── app/       # 应用核心逻辑
└── scripts/       # 构建和部署脚本
```

## 快速开始

### 环境要求

- Go 1.25.5+
- 数据库（MySQL/PostgreSQL）

### 安装和运行

1. 克隆项目
```bash
git clone https://github.com/kitouo/taskhub.git
cd taskhub
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置环境
```bash
# 复制配置文件模板
cp configs/config.example.yaml configs/config.yaml
# 编辑配置文件，设置数据库连接等
```

4. 运行应用
```bash
go run cmd/api/main.go
```

## API文档

启动服务后，访问 `http://localhost:8080/docs` 查看API文档。

## 开发指南

### 代码规范

- 遵循Go官方代码规范
- 使用gofmt格式化代码
- 编写单元测试

### 提交规范

- feat: 新功能
- fix: 修复bug
- docs: 文档更新
- style: 代码格式调整
- refactor: 代码重构
- test: 测试相关
- chore: 构建过程或辅助工具的变动

## 贡献

欢迎提交Issue和Pull Request来帮助改进项目。

## 许可证

本项目采用 [LICENSE](LICENSE) 许可证。