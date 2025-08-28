package integration

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/deferz/comfyui2go"
	"github.com/deferz/comfyui2go/tests/helpers"
)

// TestBasicAPIs 测试基本API功能
func TestBasicAPIs(t *testing.T) {
	helpers.SkipIntegrationTest(t)

	client := helpers.NewTestClient("basic-api-test")
	defer helpers.CleanupClient(client)

	ctx := context.Background()

	t.Run("测试队列查询", func(t *testing.T) {
		queue, err := client.GetQueue(ctx)
		if err != nil {
			t.Fatalf("查询队列失败: %v", err)
		}

		t.Logf("📋 当前队列状态: 执行中 %d 个任务", len(queue.QueueRunning))
		t.Log("✅ 队列查询成功")
	})

	t.Run("测试WebSocket状态", func(t *testing.T) {
		if !client.IsWebSocketEnabled() {
			t.Skip("WebSocket未启用，跳过此测试")
		}

		connected, err := client.GetWebSocketStatus(ctx)
		if err != nil {
			t.Errorf("获取WebSocket状态失败: %v", err)
		} else if connected {
			t.Log("✅ WebSocket连接正常")
		} else {
			t.Log("ℹ️ WebSocket未连接")
		}
	})
}

// TestWorkflowExecution 测试工作流执行
func TestWorkflowExecution(t *testing.T) {
	helpers.SkipIntegrationTest(t)

	client := helpers.NewTestClient("workflow-test")
	defer helpers.CleanupClient(client)

	ctx := context.Background()

	// 读取测试工作流
	workflowData, err := os.ReadFile("../test.json")
	if err != nil {
		t.Skipf("跳过工作流测试，找不到 test.json: %v", err)
	}

	var workflow comfyui2go.JSON
	if err := json.Unmarshal(workflowData, &workflow); err != nil {
		t.Fatalf("解析工作流JSON失败: %v", err)
	}

	t.Run("提交工作流", func(t *testing.T) {
		promptID, err := client.Prompt(ctx, workflow)
		if err != nil {
			t.Fatalf("提交工作流失败: %v", err)
		}

		t.Logf("✅ 工作流提交成功: %s", promptID)

		// 使用轮询方式等待完成（更稳定）
		result, err := client.WaitForCompletion(ctx, promptID, 2*time.Second)
		if err != nil {
			t.Errorf("等待完成失败: %v", err)
		} else {
			t.Logf("🎉 任务完成: %s", result.PromptID)
		}
	})
}

// TestWebSocketConnection 测试WebSocket连接
func TestWebSocketConnection(t *testing.T) {
	helpers.SkipIntegrationTest(t)

	client := helpers.NewTestClient("websocket-test")
	defer helpers.CleanupClient(client)

	ctx := context.Background()

	if !client.IsWebSocketEnabled() {
		t.Skip("WebSocket未启用，跳过此测试")
	}

	t.Run("WebSocket连接测试", func(t *testing.T) {
		// 第一次连接
		connected1, err := client.GetWebSocketStatus(ctx)
		if err != nil {
			t.Fatalf("第一次WebSocket连接失败: %v", err)
		}
		t.Logf("🔌 第一次连接状态: %v", connected1)

		// 第二次应该复用连接
		connected2, err := client.GetWebSocketStatus(ctx)
		if err != nil {
			t.Fatalf("第二次WebSocket连接失败: %v", err)
		}
		t.Logf("🔗 第二次连接状态: %v", connected2)

		if connected1 && connected2 {
			t.Log("✅ WebSocket连接复用正常")
		}
	})
}
