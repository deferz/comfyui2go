# ComfyUI2Go 测试套件

本目录包含 ComfyUI2Go 库的所有测试文件，按功能和类型进行了组织。

## 📁 目录结构

```
tests/
├── helpers/               # 测试辅助函数
│   └── test_helpers.go   # 通用测试工具和客户端创建函数
├── unit/                 # 单元测试
│   ├── client_test.go    # 客户端创建和配置测试
│   └── websocket_test.go # WebSocket功能单元测试
├── integration/          # 集成测试
│   ├── basic_api_test.go # 基本API集成测试
│   ├── websocket_test.go # WebSocket集成测试
│   └── concurrent_test.go # 并发功能测试
├── test.json           # 测试工作流文件
├── run_tests.sh         # 测试运行器脚本
└── README.md           # 本文件
```

## 🧪 测试类型

### 单元测试 (Unit Tests)
- **位置**: `unit/` 目录
- **目的**: 测试独立的功能模块，不需要外部依赖
- **特点**: 快速执行，不需要真实的ComfyUI服务器
- **包含**:
  - 客户端创建和配置
  - WebSocket启用/禁用逻辑
  - 回调函数配置
  - 错误处理

### 集成测试 (Integration Tests)  
- **位置**: `integration/` 目录
- **目的**: 测试与真实ComfyUI服务器的交互
- **特点**: 需要真实的服务器环境，执行时间较长
- **包含**:
  - API调用测试
  - WebSocket连接测试
  - 工作流执行测试
  - 并发操作测试

### 测试辅助 (Test Helpers)
- **位置**: `helpers/` 目录
- **目的**: 提供通用的测试工具和客户端创建函数
- **包含**:
  - `NewTestClient()` - 创建标准测试客户端
  - `NewTestClientWithOptions()` - 创建自定义配置客户端
  - `SkipIntegrationTest()` - 跳过集成测试的条件检查
  - `CleanupClient()` - 清理客户端资源

### 测试数据 (Test Data)
- **文件**: `test.json`
- **目的**: 集成测试使用的标准ComfyUI工作流
- **用途**: 
  - 工作流执行测试
  - WebSocket连接测试
  - 并发任务测试
- **格式**: 标准的ComfyUI工作流JSON格式

## 🚀 运行测试

### 使用测试运行器（推荐）

```bash
# 运行所有测试
./tests/run_tests.sh all

# 只运行单元测试
./tests/run_tests.sh unit

# 只运行集成测试
./tests/run_tests.sh integration

# 详细输出模式
TEST_VERBOSE=true ./tests/run_tests.sh all

# 自定义超时时间
TEST_TIMEOUT=120s ./tests/run_tests.sh integration
```

### 使用 Go 命令

```bash
# 运行单元测试
go test -v ./tests/unit/...

# 运行集成测试
go test -v ./tests/integration/...

# 跳过集成测试（使用-short标志）
go test -short -v ./tests/...

# 运行所有测试
go test -v ./tests/...
```

## ⚙️ 测试配置

### 环境变量

- `TEST_TIMEOUT`: 测试超时时间（默认: 60s）
- `TEST_VERBOSE`: 详细输出模式（默认: false）

### 集成测试配置

集成测试需要本地ComfyUI服务器：
- 地址: `http://127.0.0.1:8812/`  
- 用户名: `admin`
- 密码: `admin123456`

配置位置: `helpers/test_helpers.go`

跳过集成测试: `go test -short ./tests/...`

## 📝 编写新测试

### 单元测试示例

```go
package unit

import (
    "testing"
    "github.com/deferz/comfyui2go"
)

func TestNewFeature(t *testing.T) {
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

func TestNewIntegration(t *testing.T) {
    helpers.SkipIntegrationTest(t)
    
    client := helpers.NewTestClient("test")
    defer helpers.CleanupClient(client)
    
    // 测试逻辑...
}
```

## 🔧 编写规则

- 单元测试：独立功能，无外部依赖
- 集成测试：真实服务器，完整流程
- 清理资源：`defer helpers.CleanupClient(client)`
- 跳过集成：`helpers.SkipIntegrationTest(t)`

## 📊 覆盖率

```bash
go test -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out
```
