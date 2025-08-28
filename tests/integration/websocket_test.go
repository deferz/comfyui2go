package integration

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/deferz/comfyui2go"
	"github.com/deferz/comfyui2go/tests/helpers"
)

// TestWebSocketCallbacks 测试WebSocket回调功能
func TestWebSocketCallbacks(t *testing.T) {
	helpers.SkipIntegrationTest(t)

	var mu sync.Mutex

	// 创建带回调的客户端
	client := helpers.NewTestClientWithCallbacks("callback-test", comfyui2go.WSCallbackConfig{
		OnProgress: func(promptID string, progress comfyui2go.WSProgressMessage) {
			mu.Lock()
			defer mu.Unlock()
			t.Logf("📊 进度回调: 任务 %s, 进度 %d/%d", promptID, progress.Value, progress.Max)
		},
		OnExecution: func(promptID string, nodeID *string) {
			mu.Lock()
			defer mu.Unlock()
			if nodeID == nil {
				t.Logf("🎉 执行回调: 任务 %s 完成", promptID)
			} else {
				t.Logf("⚙️ 执行回调: 任务 %s 执行节点 %s", promptID, *nodeID)
			}
		},
		OnError: func(promptID string, err error) {
			mu.Lock()
			defer mu.Unlock()
			t.Logf("❌ 错误回调: 任务 %s 出错: %v", promptID, err)
		},
	})
	defer helpers.CleanupClient(client)

	ctx := context.Background()

	// 确保WebSocket连接
	if client.IsWebSocketEnabled() {
		_, err := client.GetWebSocketStatus(ctx)
		if err != nil {
			t.Skipf("WebSocket连接失败，跳过回调测试: %v", err)
		}
	} else {
		t.Skip("WebSocket未启用，跳过回调测试")
	}

	t.Log("✅ WebSocket回调测试设置完成")
}

// TestSingleWebSocketConnection 测试单WebSocket连接复用
func TestSingleWebSocketConnection(t *testing.T) {
	helpers.SkipIntegrationTest(t)

	client := helpers.NewTestClient("single-conn-test")
	defer helpers.CleanupClient(client)

	ctx := context.Background()

	// 读取测试工作流
	workflowData, err := os.ReadFile("../test.json")
	if err != nil {
		t.Skipf("跳过连接测试，找不到 test.json: %v", err)
	}

	var workflow comfyui2go.JSON
	if err := json.Unmarshal(workflowData, &workflow); err != nil {
		t.Fatalf("解析工作流JSON失败: %v", err)
	}

	if !client.IsWebSocketEnabled() {
		t.Skip("WebSocket未启用，跳过连接复用测试")
	}

	t.Run("连接复用测试", func(t *testing.T) {
		// 多次获取WebSocket状态，应该复用同一个连接
		for i := 0; i < 3; i++ {
			connected, err := client.GetWebSocketStatus(ctx)
			if err != nil {
				t.Errorf("第 %d 次连接失败: %v", i+1, err)
				continue
			}
			t.Logf("🔗 第 %d 次连接状态: %v", i+1, connected)
		}
		t.Log("✅ WebSocket连接复用测试完成")
	})
}

// TestExplicitChoice 测试明确选择等待方式
func TestExplicitChoice(t *testing.T) {
	helpers.SkipIntegrationTest(t)

	t.Run("明确选择WebSocket方式", func(t *testing.T) {
		client := helpers.NewTestClient("explicit-ws")
		defer helpers.CleanupClient(client)

		if !client.IsWebSocketEnabled() {
			t.Error("WebSocket应该启用")
		}

		ctx := context.Background()

		// 这里不执行实际的任务，只测试方法调用
		// 使用假的promptID会超时，但能验证WebSocket路径
		_, err := client.WaitForCompletionWithWS(ctx, "fake-prompt-id", 1*time.Second)
		if err != nil {
			t.Logf("✅ 明确使用WebSocket方式 (预期超时): %v", err)
		}
	})

	t.Run("明确选择轮询方式", func(t *testing.T) {
		client := helpers.NewTestClientWithOptions("explicit-poll", comfyui2go.WithoutWebSocket())
		defer helpers.CleanupClient(client)

		if client.IsWebSocketEnabled() {
			t.Error("WebSocket应该禁用")
		}

		ctx := context.Background()

		// 使用轮询方式，假的promptID会快速失败
		_, err := client.WaitForCompletion(ctx, "fake-prompt-id", 1*time.Second)
		if err != nil {
			t.Logf("✅ 明确使用轮询方式 (预期失败): %v", err)
		}
	})

	t.Run("状态检查后选择", func(t *testing.T) {
		// 启用WebSocket的客户端
		wsClient := helpers.NewTestClient("choice-ws")
		defer helpers.CleanupClient(wsClient)

		// 禁用WebSocket的客户端
		pollClient := helpers.NewTestClientWithOptions("choice-poll", comfyui2go.WithoutWebSocket())
		defer helpers.CleanupClient(pollClient)

		// 测试状态检查
		clients := []*comfyui2go.Client{wsClient, pollClient}
		expectedStates := []bool{true, false}

		for i, client := range clients {
			if client.IsWebSocketEnabled() != expectedStates[i] {
				t.Errorf("客户端%d WebSocket状态不符合预期", i)
			}

			// 模拟根据状态选择等待方式的逻辑
			if client.IsWebSocketEnabled() {
				t.Logf("✅ 客户端%d: 选择WebSocket方式", i)
			} else {
				t.Logf("✅ 客户端%d: 选择轮询方式", i)
			}
		}
	})
}
