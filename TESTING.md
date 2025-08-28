# 测试指南

## 目录结构

```
tests/
├── helpers/               # 测试工具
├── unit/                 # 单元测试  
├── integration/          # 集成测试
├── test.json            # 测试工作流
└── run_tests.sh         # 测试脚本
```

## 运行测试

```bash
# Makefile方式（推荐）
make test-unit        # 单元测试
make test-integration # 集成测试
make test            # 所有测试

# 脚本方式
./tests/run_tests.sh unit
./tests/run_tests.sh integration
./tests/run_tests.sh all

# Go命令
go test ./tests/unit/...
go test ./tests/integration/...
go test ./tests/...
```

## 测试类型

### 单元测试
- 独立功能测试
- 无外部依赖
- 快速执行

### 集成测试  
- 真实服务器交互
- 完整流程测试
- 需要ComfyUI服务器

## 配置

集成测试配置（`helpers/test_helpers.go`）：
- 服务器: `http://127.0.0.1:8812/`
- 用户名: `admin`  
- 密码: `admin123456`

跳过集成测试: `go test -short ./tests/...`

## 编写测试

### 单元测试示例

```go
package unit

import (
    "testing"
    "github.com/deferz/comfyui2go"
)

func TestClient(t *testing.T) {
    client := comfyui2go.NewClient("test", "http://localhost:8188")
    defer client.CloseWebSocket()
    
    // 测试逻辑...
}
```

### 集成测试示例

```go
package integration

import (
    "testing"
    "github.com/deferz/comfyui2go/tests/helpers"
)

func TestAPI(t *testing.T) {
    helpers.SkipIntegrationTest(t)
    
    client := helpers.NewTestClient("test")
    defer helpers.CleanupClient(client)
    
    // 测试逻辑...
}
```

## 工具函数

```go
// 创建测试客户端
client := helpers.NewTestClient("id")
client := helpers.NewTestClientWithOptions("id", options...)

// 清理资源
defer helpers.CleanupClient(client)

// 跳过集成测试
helpers.SkipIntegrationTest(t)
```

## 覆盖率

```bash
make test-coverage
# 或
go test -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out
```
