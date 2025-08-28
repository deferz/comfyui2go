# ComfyUI2Go 测试指南

## 📁 新的测试结构

经过重新组织，所有测试文件现在都位于独立的 `tests/` 目录中，并按功能进行了分类：

```
tests/
├── helpers/               # 🛠️ 测试辅助工具
│   └── test_helpers.go   # 通用客户端创建和测试工具
├── unit/                 # 🧪 单元测试
│   ├── client_test.go    # 客户端创建和配置
│   └── websocket_test.go # WebSocket功能
├── integration/          # 🔗 集成测试
│   ├── basic_api_test.go # 基本API测试
│   ├── websocket_test.go # WebSocket集成测试
│   └── concurrent_test.go # 并发功能测试
├── run_tests.sh         # 🚀 测试运行器脚本
└── README.md           # 📖 详细文档
```

## 🚀 快速开始

### 运行测试的方式

1. **使用Makefile（推荐）**：
   ```bash
   make test-unit        # 单元测试
   make test-integration # 集成测试
   make test            # 所有测试
   ```

2. **使用测试运行器**：
   ```bash
   ./tests/run_tests.sh unit        # 单元测试
   ./tests/run_tests.sh integration # 集成测试
   ./tests/run_tests.sh all         # 所有测试
   ```

3. **使用Go命令**：
   ```bash
   go test ./tests/unit/...         # 单元测试
   go test ./tests/integration/...  # 集成测试
   go test ./tests/...              # 所有测试
   ```

## 🧪 测试类型详解

### 单元测试
- **目的**: 测试独立功能，无外部依赖
- **特点**: 快速执行，不需要真实服务器
- **覆盖**: 客户端创建、WebSocket配置、回调设置
- **运行**: `make test-unit`

### 集成测试
- **目的**: 测试与真实ComfyUI服务器的交互
- **特点**: 需要真实环境，执行时间较长
- **覆盖**: API调用、WebSocket连接、工作流执行
- **运行**: `make test-integration`

## 🛠️ 测试辅助工具

新的测试辅助函数简化了客户端创建：

```go
// 创建标准测试客户端
client := helpers.NewTestClient("test-id")

// 创建带选项的客户端  
client := helpers.NewTestClientWithOptions("test-id", 
    comfyui2go.WithoutWebSocket())

// 创建带回调的客户端
client := helpers.NewTestClientWithCallbacks("test-id", callbacks)

// 清理资源
defer helpers.CleanupClient(client)

// 跳过集成测试
helpers.SkipIntegrationTest(t)
```

## 🎯 好处与改进

### ✅ 解决的问题

1. **代码重复**: 统一了客户端初始化代码
2. **结构混乱**: 按功能组织测试文件
3. **运行困难**: 提供了多种运行方式
4. **配置分散**: 集中管理测试配置

### 🚀 新的优势

1. **清晰的分离**: 单元测试 vs 集成测试
2. **易于维护**: 通用辅助函数减少重复
3. **灵活运行**: 可选择性运行不同类型测试
4. **标准化**: 统一的测试模式和最佳实践

## 📊 测试统计

### 测试覆盖

- **单元测试**: 4个测试套件，覆盖核心功能
- **集成测试**: 3个测试套件，覆盖实际使用场景
- **辅助函数**: 简化测试代码，提高可维护性

### 文件减少

- **之前**: 根目录下5个测试文件
- **现在**: `tests/` 目录下结构化组织
- **减少重复**: 聚合了客户端初始化代码

## 🔧 开发工作流

### 添加新测试

1. **单元测试**: 添加到 `tests/unit/` 对应文件
2. **集成测试**: 添加到 `tests/integration/` 对应文件
3. **使用辅助函数**: 避免重复的客户端创建代码

### 运行工作流

```bash
# 开发时快速检查
make test-unit

# 功能完成后完整测试
make test

# 生成覆盖率报告
make test-coverage

# 代码检查和格式化
make check
```

## 🎉 总结

新的测试结构提供了：

- 📁 **清晰的组织**: 按功能分类的目录结构
- 🛠️ **强大工具**: 统一的辅助函数和运行器
- 🚀 **灵活运行**: 多种运行方式适应不同需求
- 📈 **易于扩展**: 标准化的测试模式
- 🔧 **简化维护**: 减少重复代码，提高可读性

这个新结构让测试更加专业化和易于管理，为项目的长期维护奠定了良好的基础！
