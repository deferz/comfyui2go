package main // websocket_optional_demo

import (
	"fmt"

	"github.com/deferz/comfyui2go"
)

func main() {
	fmt.Println("=== WebSocket可选功能演示 ===\n")

	// 🚀 场景1: 默认启用WebSocket（推荐）
	fmt.Println("📡 场景1: 默认启用WebSocket")
	client1 := comfyui2go.NewClientWithOptions(
		"websocket-demo", "http://localhost:8188",
		// 默认启用WebSocket，无需额外配置
		comfyui2go.WithProgressCallback(func(promptID string, progress comfyui2go.WSProgressMessage) {
			fmt.Printf("   📊 WebSocket进度: %.1f%%\n",
				float64(progress.Value)/float64(progress.Max)*100)
		}),
	)
	defer client1.CloseWebSocket()

	fmt.Printf("   WebSocket启用状态: %v\n", client1.IsWebSocketEnabled())
	fmt.Println()

	// 🚀 场景2: 明确禁用WebSocket（开放平台场景）
	fmt.Println("🔌 场景2: 禁用WebSocket（开放平台）")
	client2 := comfyui2go.NewClientWithOptions(
		"openapi-demo", "https://api.third-party.com",
		comfyui2go.WithoutWebSocket(), // 禁用WebSocket
		comfyui2go.WithBasicAuth("api-key", "secret"),
	)
	defer client2.CloseWebSocket()

	fmt.Printf("   WebSocket启用状态: %v\n", client2.IsWebSocketEnabled())
	fmt.Println()

	// 🚀 场景3: 条件性启用WebSocket
	fmt.Println("⚙️ 场景3: 条件性启用WebSocket")
	isProduction := false // 假设这是从环境变量读取
	client3 := comfyui2go.NewClientWithOptions(
		"conditional-demo", "http://localhost:8188",
		comfyui2go.WithWebSocketEnabled(!isProduction), // 开发环境启用，生产环境禁用
	)
	defer client3.CloseWebSocket()

	fmt.Printf("   WebSocket启用状态: %v (生产环境: %v)\n",
		client3.IsWebSocketEnabled(), isProduction)
	fmt.Println()

	// 🚀 场景4: 明确选择等待方式
	fmt.Println("🔧 场景4: 明确选择等待方式")
	demonstrateExplicitChoice()

	// 🚀 场景5: 不同等待方式对比
	fmt.Println("\n📊 场景5: 不同等待方式对比")
	demonstrateWaitingMethods()
}

func demonstrateExplicitChoice() {
	fmt.Println("   💡 根据环境明确选择等待方式:")
	fmt.Println()

	// 示例1: 检查WebSocket状态并选择
	fmt.Println("   📡 方式1: 检查状态后选择")
	fmt.Println("   ```go")
	fmt.Println("   if client.IsWebSocketEnabled() {")
	fmt.Println("       result, err = client.WaitForCompletionWithWS(ctx, promptID, timeout)")
	fmt.Println("   } else {")
	fmt.Println("       result, err = client.WaitForCompletion(ctx, promptID, pollInterval)")
	fmt.Println("   }")
	fmt.Println("   ```")
	fmt.Println()

	// 示例2: 明确配置不同环境
	fmt.Println("   🏠 方式2: 不同环境明确配置")
	fmt.Println("   ```go")
	fmt.Println("   // 本地开发 - 启用WebSocket")
	fmt.Println("   localClient := comfyui2go.NewClient(\"app\", \"http://localhost:8188\")")
	fmt.Println("   result, err := localClient.WaitForCompletionWithWS(ctx, promptID, timeout)")
	fmt.Println()
	fmt.Println("   // 开放平台 - 禁用WebSocket")
	fmt.Println("   apiClient := comfyui2go.NewClientWithOptions(")
	fmt.Println("       \"app\", \"https://api.platform.com\",")
	fmt.Println("       comfyui2go.WithoutWebSocket(),")
	fmt.Println("   )")
	fmt.Println("   result, err := apiClient.WaitForCompletion(ctx, promptID, pollInterval)")
	fmt.Println("   ```")
}

func demonstrateWaitingMethods() {
	fmt.Println("   💡 不同等待方式的使用场景:")
	fmt.Println()

	fmt.Println("   📡 WaitForCompletionWithWS:")
	fmt.Println("      - 实时进度更新")
	fmt.Println("      - 资源利用率低")
	fmt.Println("      - 需要WebSocket支持")
	fmt.Println("      - 适合：本地ComfyUI、支持WS的平台")
	fmt.Println()

	fmt.Println("   🔄 WaitForCompletion:")
	fmt.Println("      - 轮询获取状态")
	fmt.Println("      - 兼容性好")
	fmt.Println("      - 可控制轮询间隔")
	fmt.Println("      - 适合：开放平台、不支持WS的环境")
	fmt.Println()

	fmt.Println("   🔧 明确选择原则:")
	fmt.Println("      - 本地ComfyUI：使用WebSocket")
	fmt.Println("      - 开放平台：使用轮询")
	fmt.Println("      - 混合环境：检查状态后选择")
	fmt.Println("      - 避免隐式降级，保持明确性")
	fmt.Println()
}

// 演示不同配置的使用代码
func demonstrateUsagePatterns() {
	fmt.Println("=== 使用模式演示 ===\n")

	// 1. 本地ComfyUI - 全功能
	fmt.Println("🏠 本地ComfyUI:")
	fmt.Println("client := comfyui2go.NewClientWithOptions(")
	fmt.Println("    \"my-app\", \"http://localhost:8188\",")
	fmt.Println("    // 默认启用WebSocket")
	fmt.Println("    comfyui2go.WithProgressCallback(...),")
	fmt.Println(")")
	fmt.Println()

	// 2. 开放平台 - 禁用WebSocket
	fmt.Println("🌐 开放平台:")
	fmt.Println("client := comfyui2go.NewClientWithOptions(")
	fmt.Println("    \"my-app\", \"https://api.platform.com\",")
	fmt.Println("    comfyui2go.WithoutWebSocket(), // 禁用WebSocket")
	fmt.Println("    comfyui2go.WithBasicAuth(\"key\", \"secret\"),")
	fmt.Println(")")
	fmt.Println()

	// 3. 云服务 - 条件启用
	fmt.Println("☁️ 云服务:")
	fmt.Println("hasWebSocket := checkWebSocketSupport(apiEndpoint)")
	fmt.Println("client := comfyui2go.NewClientWithOptions(")
	fmt.Println("    \"my-app\", apiEndpoint,")
	fmt.Println("    comfyui2go.WithWebSocketEnabled(hasWebSocket),")
	fmt.Println(")")
	fmt.Println()

	// 4. 通用方案 - 明确选择
	fmt.Println("🔧 通用方案:")
	fmt.Println("// 根据WebSocket状态明确选择")
	fmt.Println("if client.IsWebSocketEnabled() {")
	fmt.Println("    result, err = client.WaitForCompletionWithWS(ctx, promptID, timeout)")
	fmt.Println("} else {")
	fmt.Println("    result, err = client.WaitForCompletion(ctx, promptID, pollInterval)")
	fmt.Println("}")
}

func checkWebSocketSupport(endpoint string) bool {
	// 实际实现中可以通过测试连接来检查
	// 这里只是示例
	return true
}
