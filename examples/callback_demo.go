package main // callback_demo

import (
	"context"
	"fmt"
	"time"

	"github.com/deferz/comfyui2go"
)

func main() {
	// 🚀 方式1: 分别设置回调函数
	client1 := comfyui2go.NewClientWithOptions(
		"callback-demo-1", "http://localhost:8188",
		comfyui2go.WithProgressCallback(func(promptID string, progress comfyui2go.WSProgressMessage) {
			percentage := float64(progress.Value) / float64(progress.Max) * 100
			fmt.Printf("📊 [方式1] 任务 %s 进度: %.1f%% (%d/%d)\n",
				promptID, percentage, progress.Value, progress.Max)
		}),
		comfyui2go.WithExecutionCallback(func(promptID string, nodeID *string) {
			if nodeID == nil {
				fmt.Printf("🎉 [方式1] 任务 %s 执行完成\n", promptID)
			} else {
				fmt.Printf("⚙️ [方式1] 任务 %s 正在执行节点: %s\n", promptID, *nodeID)
			}
		}),
		comfyui2go.WithErrorCallback(func(promptID string, err error) {
			fmt.Printf("❌ [方式1] 任务 %s 出错: %v\n", promptID, err)
		}),
	)
	defer client1.CloseWebSocket()

	// 🚀 方式2: 一次性设置所有回调
	client2 := comfyui2go.NewClientWithOptions(
		"callback-demo-2", "http://localhost:8188",
		comfyui2go.WithWebSocketCallbacks(comfyui2go.WSCallbackConfig{
			OnProgress: func(promptID string, progress comfyui2go.WSProgressMessage) {
				percentage := float64(progress.Value) / float64(progress.Max) * 100
				fmt.Printf("📈 [方式2] 任务 %s 进度: %.1f%% (%d/%d)\n",
					promptID, percentage, progress.Value, progress.Max)
			},
			OnExecution: func(promptID string, nodeID *string) {
				if nodeID == nil {
					fmt.Printf("✅ [方式2] 任务 %s 执行完成\n", promptID)
				} else {
					fmt.Printf("🔄 [方式2] 任务 %s 正在执行节点: %s\n", promptID, *nodeID)
				}
			},
			OnStatus: func(promptID string, status string) {
				fmt.Printf("📋 [方式2] 任务 %s 状态变化: %s\n", promptID, status)
			},
			OnError: func(promptID string, err error) {
				fmt.Printf("🔥 [方式2] 任务 %s 出错: %v\n", promptID, err)
			},
		}),
	)
	defer client2.CloseWebSocket()

	// 🚀 方式3: 使用全局回调的简洁客户端
	client3 := comfyui2go.NewClientWithOptions(
		"callback-demo-3", "http://localhost:8188",
		comfyui2go.WithProgressCallback(createGlobalProgressTracker()),
		comfyui2go.WithExecutionCallback(createGlobalExecutionTracker()),
	)
	defer client3.CloseWebSocket()

	fmt.Println("🎯 演示不同的回调配置方式")
	fmt.Println("💡 提示：这只是演示回调配置，实际使用时选择其中一个客户端即可")

	// 实际使用示例（使用client1）
	demonstrateUsage(client1)
}

func demonstrateUsage(client *comfyui2go.Client) {
	ctx := context.Background()

	// 创建一个简单的工作流
	workflow := comfyui2go.JSON{
		"1": map[string]interface{}{
			"class_type": "CheckpointLoaderSimple",
			"inputs": map[string]interface{}{
				"ckpt_name": "v1-5-pruned-emaonly.ckpt",
			},
		},
		"2": map[string]interface{}{
			"class_type": "CLIPTextEncode",
			"inputs": map[string]interface{}{
				"text": "a cute cat",
				"clip": []interface{}{"1", 1},
			},
		},
		// ... 可以继续添加更多节点
	}

	fmt.Println("\n🚀 开始执行工作流...")

	// 提交工作流
	promptID, err := client.Prompt(ctx, workflow)
	if err != nil {
		fmt.Printf("❌ 提交失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 工作流已提交: %s\n", promptID)
	fmt.Println("📡 WebSocket回调将自动处理进度和状态更新...")

	// 等待完成（回调函数会自动报告进度）
	result, err := client.WaitForCompletionWithWS(ctx, promptID, 2*time.Minute)
	if err != nil {
		fmt.Printf("❌ 执行失败: %v\n", err)
		return
	}

	fmt.Printf("🎊 最终结果: %s\n", result.PromptID)
}

// 创建全局进度追踪器
func createGlobalProgressTracker() comfyui2go.ProgressCallback {
	taskProgress := make(map[string]int)

	return func(promptID string, progress comfyui2go.WSProgressMessage) {
		lastProgress, exists := taskProgress[promptID]
		currentProgress := int(float64(progress.Value) / float64(progress.Max) * 100)

		// 只在进度有显著变化时输出
		if !exists || currentProgress-lastProgress >= 5 {
			taskProgress[promptID] = currentProgress
			fmt.Printf("🌟 [全局追踪] 任务 %s: %d%% 完成\n", promptID, currentProgress)
		}
	}
}

// 创建全局执行追踪器
func createGlobalExecutionTracker() comfyui2go.ExecutionCallback {
	return func(promptID string, nodeID *string) {
		if nodeID == nil {
			fmt.Printf("🏁 [全局追踪] 任务 %s 全部完成！\n", promptID)
		} else {
			fmt.Printf("🎯 [全局追踪] 任务 %s 执行节点 %s\n", promptID, *nodeID)
		}
	}
}
