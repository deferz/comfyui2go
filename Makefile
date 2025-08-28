# ComfyUI2Go Makefile

.PHONY: help test test-unit test-integration test-all clean build examples

# 默认目标
help:
	@echo "ComfyUI2Go 构建和测试工具"
	@echo ""
	@echo "可用目标:"
	@echo "  help            显示此帮助信息"
	@echo "  build           构建项目"
	@echo "  test            运行所有测试"
	@echo "  test-unit       运行单元测试"
	@echo "  test-integration 运行集成测试"
	@echo "  test-verbose    运行详细模式测试"
	@echo "  clean           清理构建文件"
	@echo "  examples        构建示例程序"
	@echo ""
	@echo "示例:"
	@echo "  make test-unit              # 快速单元测试"
	@echo "  make test-integration       # 集成测试"
	@echo "  make test                   # 所有测试"

# 构建项目
build:
	@echo "🔨 构建项目..."
	go build ./...
	@echo "✅ 构建完成"

# 运行所有测试
test:
	@./tests/run_tests.sh all

# 运行单元测试
test-unit:
	@./tests/run_tests.sh unit

# 运行集成测试
test-integration:
	@./tests/run_tests.sh integration

# 详细模式测试
test-verbose:
	@TEST_VERBOSE=true ./tests/run_tests.sh all

# 快速测试（跳过集成测试）
test-fast:
	@echo "🚀 运行快速测试（跳过集成测试）..."
	go test -short -v ./tests/...

# 测试覆盖率
test-coverage:
	@echo "📊 生成测试覆盖率..."
	go test -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ 覆盖率报告已生成: coverage.html"

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	@./tests/run_tests.sh clean
	@find . -name "*.test" -delete 2>/dev/null || true
	@find . -name "*.out" -delete 2>/dev/null || true
	@find . -name "*.prof" -delete 2>/dev/null || true
	@find . -name "coverage.html" -delete 2>/dev/null || true
	@find . -name "*.log" -delete 2>/dev/null || true
	@find . -name "*.tmp" -delete 2>/dev/null || true
	@find . -name "*.temp" -delete 2>/dev/null || true
	@find examples/ -name "basic_usage" -delete 2>/dev/null || true
	@find examples/ -name "concurrent_tasks" -delete 2>/dev/null || true
	@find examples/ -name "callback_demo" -delete 2>/dev/null || true
	@find examples/ -name "new_api_demo" -delete 2>/dev/null || true
	@find examples/ -name "websocket_optional_demo" -delete 2>/dev/null || true
	@find . -name "simple_test_*" -delete 2>/dev/null || true
	@find . -name "concurrent_test_*" -delete 2>/dev/null || true
	@find . -name "callback_test_*" -delete 2>/dev/null || true
	@find . -name "websocket_test_*" -delete 2>/dev/null || true
	@rm -rf build/ dist/ tmp/ 2>/dev/null || true
	@echo "✅ 清理完成"

# 构建示例程序
examples:
	@echo "📋 构建示例程序..."
	@cd examples && for f in *.go; do \
		if [ -f "$$f" ]; then \
			echo "  构建 $$f..."; \
			go build "$$f" || exit 1; \
		fi \
	done
	@echo "✅ 示例程序构建完成"

# 运行示例程序
run-example:
	@echo "可用的示例程序:"
	@cd examples && ls -1 *.go | sed 's/.go$$//' | sed 's/^/  make run-/'

run-basic:
	@cd examples && go run basic_usage.go

run-concurrent:
	@cd examples && go run concurrent_tasks.go

run-callback:
	@cd examples && go run callback_demo.go

run-websocket:
	@cd examples && go run websocket_optional_demo.go

run-api:
	@cd examples && go run new_api_demo.go

# 开发工具
fmt:
	@echo "🎨 格式化代码..."
	go fmt ./...
	@echo "✅ 代码格式化完成"

vet:
	@echo "🔍 静态分析..."
	go vet ./...
	@echo "✅ 静态分析完成"

lint:
	@echo "📝 代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint 未安装，跳过代码检查"; \
		echo "   安装: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# 检查所有内容
check: fmt vet test-unit
	@echo "🎯 所有检查完成"

# 完整验证（包括集成测试）
verify: fmt vet test
	@echo "🎉 完整验证通过"

# 安装开发依赖
dev-deps:
	@echo "📦 安装开发依赖..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ 开发依赖安装完成"
