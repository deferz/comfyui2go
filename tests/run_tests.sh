#!/bin/bash

# ComfyUI2Go 测试运行器

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的信息
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 帮助信息
show_help() {
    echo "ComfyUI2Go 测试运行器"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  unit        运行单元测试"
    echo "  integration 运行集成测试"
    echo "  all         运行所有测试"
    echo "  clean       清理测试输出"
    echo "  help        显示此帮助信息"
    echo ""
    echo "环境变量:"
    echo "  TEST_TIMEOUT    测试超时时间 (默认: 60s)"
    echo "  TEST_VERBOSE    详细输出 (默认: false)"
    echo ""
    echo "示例:"
    echo "  $0 unit                    # 只运行单元测试"
    echo "  $0 integration             # 只运行集成测试"
    echo "  TEST_VERBOSE=true $0 all   # 运行所有测试（详细输出）"
}

# 默认配置
TEST_TIMEOUT=${TEST_TIMEOUT:-60s}
TEST_VERBOSE=${TEST_VERBOSE:-false}

# 构建测试参数
build_test_args() {
    local args="-timeout ${TEST_TIMEOUT}"
    
    if [ "$TEST_VERBOSE" = "true" ]; then
        args="$args -v"
    fi
    
    echo "$args"
}

# 运行单元测试
run_unit_tests() {
    print_info "运行单元测试..."
    
    cd "$(dirname "$0")"
    
    local args=$(build_test_args)
    
    if go test $args ./unit/...; then
        print_success "单元测试通过"
        return 0
    else
        print_error "单元测试失败"
        return 1
    fi
}

# 运行集成测试
run_integration_tests() {
    print_info "运行集成测试..."
    print_warning "注意: 集成测试需要连接到真实的ComfyUI服务器"
    
    cd "$(dirname "$0")"
    
    local args=$(build_test_args)
    
    if go test $args ./integration/...; then
        print_success "集成测试通过"
        return 0
    else
        print_error "集成测试失败"
        return 1
    fi
}

# 运行所有测试
run_all_tests() {
    print_info "运行所有测试..."
    
    local unit_result=0
    local integration_result=0
    
    # 运行单元测试
    if ! run_unit_tests; then
        unit_result=1
    fi
    
    echo ""
    
    # 运行集成测试
    if ! run_integration_tests; then
        integration_result=1
    fi
    
    echo ""
    
    # 总结
    if [ $unit_result -eq 0 ] && [ $integration_result -eq 0 ]; then
        print_success "所有测试通过! 🎉"
        return 0
    else
        if [ $unit_result -ne 0 ]; then
            print_error "单元测试失败"
        fi
        if [ $integration_result -ne 0 ]; then
            print_error "集成测试失败"
        fi
        return 1
    fi
}

# 清理测试输出
clean_tests() {
    print_info "清理测试输出..."
    
    cd "$(dirname "$0")/.."
    
    # 删除测试生成的文件（与.gitignore保持一致）
    find . -name "*.test" -delete 2>/dev/null || true
    find . -name "*.out" -delete 2>/dev/null || true  
    find . -name "*.prof" -delete 2>/dev/null || true
    find . -name "coverage.html" -delete 2>/dev/null || true
    find . -name "*.log" -delete 2>/dev/null || true
    find . -name "*.tmp" -delete 2>/dev/null || true
    find . -name "*.temp" -delete 2>/dev/null || true
    find . -name "simple_test_*" -delete 2>/dev/null || true
    find . -name "concurrent_test_*" -delete 2>/dev/null || true
    find . -name "callback_test_*" -delete 2>/dev/null || true
    find . -name "websocket_test_*" -delete 2>/dev/null || true
    
    print_success "清理完成"
}

# 检查Go环境
check_go_env() {
    if ! command -v go &> /dev/null; then
        print_error "Go未安装或不在PATH中"
        exit 1
    fi
    
    local go_version=$(go version | cut -d' ' -f3)
    print_info "Go版本: $go_version"
}

# 主函数
main() {
    check_go_env
    
    case "${1:-all}" in
        "unit")
            run_unit_tests
            ;;
        "integration")
            run_integration_tests
            ;;
        "all")
            run_all_tests
            ;;
        "clean")
            clean_tests
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "未知选项: $1"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@"
