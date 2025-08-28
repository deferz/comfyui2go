package comfyui2go

import (
	"time"

	resty "resty.dev/v3"
)

// Option 用于在创建客户端时自定义配置。
// 推荐使用 NewClientWithOptions(clientID, baseURL, opts...) 而不是 New(opts...)
type Option func(*Client)

func WithDebug(enable bool) Option {
	return func(c *Client) {
		c.cli.SetDebug(enable)
	}
}

// WithClientID 设置用于在 ComfyUI 端标记任务来源的稳定客户端 ID。
func WithClientID(id string) Option { return func(c *Client) { c.clientID = id } }

// WithHTTP 允许传入已预配置的 Resty 客户端（如代理、TLS、全局头等）。
func WithHTTP(rc *resty.Client) Option { return func(c *Client) { c.cli = rc } }

// WithTimeout 设置底层 Resty 客户端的请求超时时间。
func WithTimeout(d time.Duration) Option { return func(c *Client) { c.cli.SetTimeout(d) } }

// WithBaseURL 设置 ComfyUI 服务器的基础地址。
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.cli.SetBaseURL(url)
		c.baseURL = url
	}
}

// WithBasicAuth 为所有请求设置 HTTP 基本认证。
func WithBasicAuth(user, pass string) Option {
	return func(c *Client) {
		c.cli.SetBasicAuth(user, pass)
		c.username = user
		c.password = pass
	}
}

// WithProgressCallback 设置进度回调函数
func WithProgressCallback(callback ProgressCallback) Option {
	return func(c *Client) {
		c.onProgress = callback
	}
}

// WithStatusCallback 设置状态变化回调函数
func WithStatusCallback(callback StatusCallback) Option {
	return func(c *Client) {
		c.onStatus = callback
	}
}

// WithExecutionCallback 设置执行状态回调函数
func WithExecutionCallback(callback ExecutionCallback) Option {
	return func(c *Client) {
		c.onExecution = callback
	}
}

// WithErrorCallback 设置错误回调函数
func WithErrorCallback(callback ErrorCallback) Option {
	return func(c *Client) {
		c.onError = callback
	}
}

// WithWebSocketCallbacks 一次性设置所有WebSocket回调函数
func WithWebSocketCallbacks(config WSCallbackConfig) Option {
	return func(c *Client) {
		if config.OnProgress != nil {
			c.onProgress = config.OnProgress
		}
		if config.OnStatus != nil {
			c.onStatus = config.OnStatus
		}
		if config.OnExecution != nil {
			c.onExecution = config.OnExecution
		}
		if config.OnError != nil {
			c.onError = config.OnError
		}
	}
}

// WithWebSocketEnabled 设置是否启用WebSocket（默认为true）
func WithWebSocketEnabled(enabled bool) Option {
	return func(c *Client) {
		c.wsEnabled = enabled
	}
}

// WithoutWebSocket 禁用WebSocket（用于不支持WebSocket的开放平台）
func WithoutWebSocket() Option {
	return func(c *Client) {
		c.wsEnabled = false
	}
}

// WSCallbackConfig WebSocket回调配置
type WSCallbackConfig struct {
	OnProgress  ProgressCallback  // 进度回调
	OnStatus    StatusCallback    // 状态回调
	OnExecution ExecutionCallback // 执行回调
	OnError     ErrorCallback     // 错误回调
}
