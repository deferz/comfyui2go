package main // new_api_demo

import (
	"context"
	"fmt"
	"time"

	"github.com/deferz/comfyui2go"
)

func main() {
	// 🚀 新API方式1: 简单客户端创建
	fmt.Println("=== 方式1: 简单客户端创建 ===")
	client1 := comfyui2go.NewClient("my-app-v1", "http://localhost:8188")
	defer client1.CloseWebSocket()

	fmt.Printf("✅ 客户端1创建成功: ClientID=%s, BaseURL=%s\n",
		"my-app-v1", "http://localhost:8188")

	// 🚀 新API方式2: 带选项的客户端创建
	fmt.Println("\n=== 方式2: 带选项的客户端创建 ===")
	client2 := comfyui2go.NewClientWithOptions(
		"my-advanced-app",
		"http://115.238.30.185:8812",

		// 认证配置
		comfyui2go.WithBasicAuth("admin", "ZmtBGthP5TPFqs2U5m68"),

		// 超时配置
		comfyui2go.WithTimeout(30*time.Second),

		// 回调配置
		comfyui2go.WithProgressCallback(func(promptID string, progress comfyui2go.WSProgressMessage) {
			percentage := float64(progress.Value) / float64(progress.Max) * 100
			fmt.Printf("📊 [新API] 任务 %s 进度: %.1f%%\n", promptID, percentage)
		}),
		comfyui2go.WithExecutionCallback(func(promptID string, nodeID *string) {
			if nodeID == nil {
				fmt.Printf("🎉 [新API] 任务 %s 执行完成\n", promptID)
			} else {
				fmt.Printf("⚙️ [新API] 任务 %s 正在执行节点: %s\n", promptID, *nodeID)
			}
		}),

		// 调试配置
		comfyui2go.WithDebug(false),
	)
	defer client2.CloseWebSocket()

	fmt.Printf("✅ 客户端2创建成功: ClientID=%s, BaseURL=%s\n",
		"my-advanced-app", "http://115.238.30.185:8812")

	// 🚀 完整配置示例
	fmt.Println("\n=== 完整配置示例 ===")
	client3 := comfyui2go.NewClientWithOptions(
		"full-config-client", "http://localhost:8188",
		comfyui2go.WithBasicAuth("user", "pass"),
		comfyui2go.WithTimeout(30*time.Second),
	)
	defer client3.CloseWebSocket()

	fmt.Println("✅ 客户端3创建成功 (完整配置)")

	// 🎯 API对比演示
	fmt.Println("\n=== API对比演示 ===")
	demonstrateAPIComparison()

	// 🧪 功能测试
	fmt.Println("\n=== 功能测试 ===")
	testClientFunctionality(client2)
}

func demonstrateAPIComparison() {
	fmt.Println("🔍 API设计对比:")

	fmt.Println("\n💡 客户端创建方式:")

	fmt.Println("\n✅ 新方式1 (简单场景):")
	fmt.Println("   client := comfyui2go.NewClient(\"my-app\", \"http://localhost:8188\")")

	fmt.Println("\n✅ 新方式2 (复杂场景):")
	fmt.Println("   client := comfyui2go.NewClientWithOptions(")
	fmt.Println("       \"my-app\", \"http://localhost:8188\",")
	fmt.Println("       comfyui2go.WithBasicAuth(\"user\", \"pass\"),")
	fmt.Println("       comfyui2go.WithProgressCallback(...),")
	fmt.Println("       // 其他选项...")
	fmt.Println("   )")

	fmt.Println("\n💡 新API的优势:")
	fmt.Println("   1. 必需参数更明确 (clientID, baseURL)")
	fmt.Println("   2. 简单场景更简洁")
	fmt.Println("   3. 复杂场景更有序")
	fmt.Println("   4. 向后兼容，不破坏现有代码")
}

func testClientFunctionality(client *comfyui2go.Client) {
	ctx := context.Background()

	// 测试基本连接
	fmt.Println("🔌 测试WebSocket连接...")
	wsClient, err := client.GetWebSocketClient(ctx)
	if err != nil {
		fmt.Printf("❌ WebSocket连接失败: %v\n", err)
		return
	}

	if wsClient.IsConnected() {
		fmt.Println("✅ WebSocket连接成功")
	}

	// 测试队列查询
	fmt.Println("📋 测试队列查询...")
	queue, err := client.GetQueue(ctx)
	if err != nil {
		fmt.Printf("❌ 队列查询失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 队列查询成功: 执行中=%d, 待处理=%d\n",
		len(queue.QueueRunning), len(queue.QueuePending))

	fmt.Println("🎯 所有功能测试通过！")
}

// 演示不同的创建模式
func demonstrateCreationPatterns() {
	fmt.Println("\n=== 创建模式演示 ===")

	// 模式1: 开发环境 - 简单快速
	devClient := comfyui2go.NewClient("dev-client", "http://localhost:8188")
	fmt.Println("🔧 开发环境客户端创建")
	defer devClient.CloseWebSocket()

	// 模式2: 生产环境 - 完整配置
	prodClient := comfyui2go.NewClientWithOptions(
		"prod-client-v1.0",
		"https://comfyui.production.com",
		comfyui2go.WithBasicAuth("prod-user", "secure-password"),
		comfyui2go.WithTimeout(60*time.Second),
		comfyui2go.WithProgressCallback(func(promptID string, progress comfyui2go.WSProgressMessage) {
			// 生产环境的进度处理
			fmt.Printf("🏭 [生产] 任务进度: %.1f%%\n",
				float64(progress.Value)/float64(progress.Max)*100)
		}),
	)
	fmt.Println("🏭 生产环境客户端创建")
	defer prodClient.CloseWebSocket()

	// 模式3: 测试环境 - 带调试
	testClient := comfyui2go.NewClientWithOptions(
		"test-client",
		"http://test.comfyui.local:8188",
		comfyui2go.WithDebug(true),
		comfyui2go.WithTimeout(120*time.Second),
	)
	fmt.Println("🧪 测试环境客户端创建")
	defer testClient.CloseWebSocket()
}
