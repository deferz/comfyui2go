package unit

import (
	"context"
	"testing"
	"time"

	"github.com/deferz/comfyui2go"
)

// TestWebSocketOptional 测试WebSocket可选功能
func TestWebSocketOptional(t *testing.T) {
	t.Run("WebSocket禁用时的错误处理", func(t *testing.T) {
		client := comfyui2go.NewClientWithOptions(
			"test-disabled", "http://localhost:8188",
			comfyui2go.WithoutWebSocket(),
		)
		defer client.CloseWebSocket()

		if client.IsWebSocketEnabled() {
			t.Error("WebSocket应该被禁用")
		}

		ctx := context.Background()

		// 尝试使用WebSocket方法应该返回错误
		_, err := client.WaitForCompletionWithWS(ctx, "test-prompt", 30*time.Second)
		if err == nil {
			t.Error("禁用WebSocket时应该返回错误")
		}

		if err != nil && err.Error() != "" {
			t.Logf("✅ 正确的错误信息: %v", err)
		}
	})

	t.Run("WebSocket状态检查", func(t *testing.T) {
		// 启用WebSocket
		client1 := comfyui2go.NewClient("test-enabled", "http://localhost:8188")
		defer client1.CloseWebSocket()

		if !client1.IsWebSocketEnabled() {
			t.Error("应该启用WebSocket")
		}

		// 禁用WebSocket
		client2 := comfyui2go.NewClientWithOptions(
			"test-disabled", "http://localhost:8188",
			comfyui2go.WithoutWebSocket(),
		)
		defer client2.CloseWebSocket()

		if client2.IsWebSocketEnabled() {
			t.Error("应该禁用WebSocket")
		}

		t.Log("✅ WebSocket状态检查正常")
	})
}
