package unit

import (
	"testing"

	"github.com/deferz/comfyui2go"
)

// TestClientCreation 测试客户端创建
func TestClientCreation(t *testing.T) {
	t.Run("NewClient", func(t *testing.T) {
		client := comfyui2go.NewClient("test-client", "http://localhost:8188")

		if client == nil {
			t.Fatal("客户端创建失败")
		}

		// 默认应该启用WebSocket
		if !client.IsWebSocketEnabled() {
			t.Error("默认应该启用WebSocket")
		}

		// 清理
		client.CloseWebSocket()
		t.Log("✅ NewClient 测试通过")
	})

	t.Run("NewClientWithOptions", func(t *testing.T) {
		client := comfyui2go.NewClientWithOptions(
			"test-client", "http://localhost:8188",
			comfyui2go.WithBasicAuth("user", "pass"),
			comfyui2go.WithoutWebSocket(),
		)

		if client == nil {
			t.Fatal("客户端创建失败")
		}

		// 应该禁用WebSocket
		if client.IsWebSocketEnabled() {
			t.Error("应该禁用WebSocket")
		}

		// 清理
		client.CloseWebSocket()
		t.Log("✅ NewClientWithOptions 测试通过")
	})
}

// TestWebSocketConfig 测试WebSocket配置
func TestWebSocketConfig(t *testing.T) {
	t.Run("默认启用", func(t *testing.T) {
		client := comfyui2go.NewClient("test", "http://localhost:8188")
		defer client.CloseWebSocket()

		if !client.IsWebSocketEnabled() {
			t.Error("默认应该启用WebSocket")
		}
	})

	t.Run("明确禁用", func(t *testing.T) {
		client := comfyui2go.NewClientWithOptions(
			"test", "http://localhost:8188",
			comfyui2go.WithoutWebSocket(),
		)
		defer client.CloseWebSocket()

		if client.IsWebSocketEnabled() {
			t.Error("应该禁用WebSocket")
		}
	})

	t.Run("条件启用", func(t *testing.T) {
		// 测试启用
		client1 := comfyui2go.NewClientWithOptions(
			"test1", "http://localhost:8188",
			comfyui2go.WithWebSocketEnabled(true),
		)
		defer client1.CloseWebSocket()

		if !client1.IsWebSocketEnabled() {
			t.Error("应该启用WebSocket")
		}

		// 测试禁用
		client2 := comfyui2go.NewClientWithOptions(
			"test2", "http://localhost:8188",
			comfyui2go.WithWebSocketEnabled(false),
		)
		defer client2.CloseWebSocket()

		if client2.IsWebSocketEnabled() {
			t.Error("应该禁用WebSocket")
		}
	})
}

// TestCallbackConfiguration 测试回调函数配置
func TestCallbackConfiguration(t *testing.T) {
	client := comfyui2go.NewClientWithOptions(
		"callback-test", "http://localhost:8188",
		comfyui2go.WithProgressCallback(func(promptID string, progress comfyui2go.WSProgressMessage) {
			t.Logf("进度回调被调用: %s", promptID)
		}),
		comfyui2go.WithExecutionCallback(func(promptID string, nodeID *string) {
			t.Logf("执行回调被调用: %s", promptID)
		}),
	)
	defer client.CloseWebSocket()

	// 这里只测试回调是否能正确设置，实际调用需要真实的WebSocket连接
	t.Log("✅ 回调函数配置测试通过")
}

// TestBatchCallbackConfig 测试批量回调配置
func TestBatchCallbackConfig(t *testing.T) {
	callbacks := comfyui2go.WSCallbackConfig{
		OnProgress: func(promptID string, progress comfyui2go.WSProgressMessage) {
			t.Logf("批量配置 - 进度: %s", promptID)
		},
		OnExecution: func(promptID string, nodeID *string) {
			t.Logf("批量配置 - 执行: %s", promptID)
		},
		OnError: func(promptID string, err error) {
			t.Logf("批量配置 - 错误: %s, %v", promptID, err)
		},
	}

	client := comfyui2go.NewClientWithOptions(
		"batch-test", "http://localhost:8188",
		comfyui2go.WithWebSocketCallbacks(callbacks),
	)
	defer client.CloseWebSocket()

	t.Log("✅ 批量回调配置测试通过")
}
