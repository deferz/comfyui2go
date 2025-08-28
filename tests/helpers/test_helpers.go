package helpers

import (
	"testing"

	"github.com/deferz/comfyui2go"
)

// 测试配置常量
const (
	TestBaseURL  = "http://127.0.0.1:8812/"
	TestUsername = "admin"
	TestPassword = "admin123456"
)

// NewTestClient 创建标准测试客户端
func NewTestClient(clientID string) *comfyui2go.Client {
	return comfyui2go.NewClientWithOptions(
		clientID, TestBaseURL,
		comfyui2go.WithBasicAuth(TestUsername, TestPassword),
	)
}

// NewTestClientWithCallbacks 创建带回调的测试客户端
func NewTestClientWithCallbacks(clientID string, callbacks comfyui2go.WSCallbackConfig) *comfyui2go.Client {
	return comfyui2go.NewClientWithOptions(
		clientID, TestBaseURL,
		comfyui2go.WithBasicAuth(TestUsername, TestPassword),
		comfyui2go.WithWebSocketCallbacks(callbacks),
	)
}

// NewTestClientWithOptions 创建自定义选项的测试客户端
func NewTestClientWithOptions(clientID string, extraOptions ...comfyui2go.Option) *comfyui2go.Client {
	options := []comfyui2go.Option{
		comfyui2go.WithBasicAuth(TestUsername, TestPassword),
	}
	options = append(options, extraOptions...)

	return comfyui2go.NewClientWithOptions(clientID, TestBaseURL, options...)
}

// SkipIntegrationTest 检查是否跳过集成测试
func SkipIntegrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
	}
}

// CleanupClient 清理客户端资源
func CleanupClient(client *comfyui2go.Client) {
	if client != nil {
		client.CloseWebSocket()
	}
}
