package comfyui2go

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	resty "resty.dev/v3"
)

type Client struct {
	cli      *resty.Client
	clientID string
	baseURL  string
	username string
	password string

	// WebSocket连接管理
	wsClient  *WSClient
	wsMu      sync.RWMutex
	wsEnabled bool // WebSocket是否启用

	// WebSocket回调函数
	onProgress  ProgressCallback
	onStatus    StatusCallback
	onExecution ExecutionCallback
	onError     ErrorCallback
}

// NewClient 创建客户端的简单方式，只需要基本参数（默认启用WebSocket）
func NewClient(clientID, baseURL string) *Client {
	r := resty.New()
	c := &Client{
		cli:       r,
		clientID:  clientID,
		baseURL:   baseURL,
		wsEnabled: true, // 默认启用WebSocket
	}
	r.SetBaseURL(baseURL)
	return c
}

// NewClientWithOptions 创建客户端的完整方式，支持所有配置选项
func NewClientWithOptions(clientID, baseURL string, opts ...Option) *Client {
	c := NewClient(clientID, baseURL)
	for _, o := range opts {
		o(c)
	}
	return c
}

// Prompt 调用 POST /prompt 提交工作流（图形 JSON）。
// 返回可用于后续查询历史记录的 prompt_id。
func (c *Client) Prompt(ctx context.Context, prompt JSON) (string, error) {
	body := JSON{
		"prompt":    prompt,
		"client_id": c.clientID,
	}

	var resp PromptResponse
	r, err := c.cli.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&resp).
		Post("/prompt")
	if err != nil {
		return "", err
	}
	if !r.IsSuccess() {
		return "", fmt.Errorf("/prompt failed: %s", r.String())
	}
	if resp.PromptID == "" {
		return "", errors.New("missing prompt_id in response")
	}
	return resp.PromptID, nil
}

// GetQueue 获取 /queue 队列状态（运行中与等待中）。
func (c *Client) GetQueue(ctx context.Context) (QueueResponse, error) {
	var out QueueResponse
	r, err := c.cli.R().SetContext(ctx).SetResult(&out).Get("/queue")
	if err != nil {
		return out, err
	}
	if !r.IsSuccess() {
		return out, fmt.Errorf("/queue failed: %s", r.String())
	}
	return out, nil
}

// GetHistory 返回指定 promptID 对应的完整历史对象。
func (c *Client) GetHistory(ctx context.Context, promptID string) (HistoryResponse, error) {
	var out HistoryResponse
	url := fmt.Sprintf("/history/%s", promptID)
	r, err := c.cli.R().SetContext(ctx).SetResult(&out).Get(url)
	if err != nil {
		return nil, err
	}
	if !r.IsSuccess() {
		return nil, fmt.Errorf("GET %s failed: %s", url, r.String())
	}
	return out, nil
}

// Interrupt 调用 /interrupt 以中断当前任务。
func (c *Client) Interrupt(ctx context.Context) error {
	r, err := c.cli.R().SetContext(ctx).Post("/interrupt")
	if err != nil {
		return err
	}
	if !r.IsSuccess() {
		return fmt.Errorf("/interrupt failed: %s", r.String())
	}
	return nil
}

// WaitForCompletion 轮询 /history/{promptID}，直到完成或上下文取消。
// pollEvery 为轮询间隔；若 <=0 则默认 1 秒。
func (c *Client) WaitForCompletion(ctx context.Context, promptID string, pollEvery time.Duration) (*WaitResult, error) {
	if pollEvery <= 0 {
		pollEvery = time.Second
	}
	ticker := time.NewTicker(pollEvery)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			h, err := c.GetHistory(ctx, promptID)
			if err != nil {
				return nil, err
			}
			if item, ok := h[promptID]; ok {
				if item.Status != nil && item.Status.Completed {
					return &WaitResult{PromptID: promptID, Item: item}, nil
				}
				if item.Outputs != nil {
					return &WaitResult{PromptID: promptID, Item: item}, nil
				}
			}
		}
	}
}

// UploadImage 上传图片文件到 ComfyUI 服务器。
// 返回文件名和子文件夹等信息，可用于后续工作流中引用。
func (c *Client) UploadImage(ctx context.Context, filename string, data io.Reader) (*UploadResponse, error) {
	var resp UploadResponse
	r, err := c.cli.R().
		SetContext(ctx).
		SetFileReader("image", filename, data).
		SetFormData(map[string]string{
			"type": "input",
		}).
		SetResult(&resp).
		Post("/upload/image")
	if err != nil {
		return nil, err
	}
	if !r.IsSuccess() {
		return nil, fmt.Errorf("/upload/image failed: %s", r.String())
	}
	return &resp, nil
}

// Download 下载 ComfyUI 生成的文件
// filename: 文件名（如 "test_00002_.png"）
// subfolder: 子文件夹（通常为空字符串）
// filetype: 文件类型（"output" 或 "input"）
func (c *Client) Download(ctx context.Context, filename, subfolder, filetype string) ([]byte, error) {
	url := "/view"
	r, err := c.cli.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"filename":  filename,
			"subfolder": subfolder,
			"type":      filetype,
		}).
		Get(url)
	if err != nil {
		return nil, err
	}
	if !r.IsSuccess() {
		return nil, fmt.Errorf("download failed: %s", r.String())
	}
	return r.Bytes(), nil
}

// ensureWebSocketConnected 确保WebSocket连接已建立
func (c *Client) ensureWebSocketConnected(ctx context.Context) error {
	// 如果WebSocket未启用，返回错误
	if !c.wsEnabled {
		return fmt.Errorf("WebSocket未启用，请使用 WithWebSocketEnabled(true) 或移除 WithoutWebSocket() 选项")
	}

	c.wsMu.Lock()
	defer c.wsMu.Unlock()

	// 如果连接存在且有效，直接返回
	if c.wsClient != nil && c.wsClient.IsConnected() {
		return nil
	}

	// 创建新的WebSocket连接，使用Client配置的回调函数
	c.wsClient = NewWSClient(WSConfig{
		BaseURL:     c.getBaseURL(),
		ClientID:    c.clientID,
		Username:    c.getUsername(),
		Password:    c.getPassword(),
		OnProgress:  c.onProgress,
		OnStatus:    c.onStatus,
		OnExecution: c.onExecution,
		OnError:     c.onError,
	})

	return c.wsClient.Connect(ctx)
}

// GetWebSocketClient 获取WebSocket客户端（确保已连接）- 内部使用
func (c *Client) GetWebSocketClient(ctx context.Context) (*WSClient, error) {
	if err := c.ensureWebSocketConnected(ctx); err != nil {
		return nil, err
	}

	c.wsMu.RLock()
	defer c.wsMu.RUnlock()
	return c.wsClient, nil
}

// IsWebSocketConnected 检查WebSocket是否已连接（用于测试和状态检查）
func (c *Client) IsWebSocketConnected() bool {
	c.wsMu.RLock()
	defer c.wsMu.RUnlock()
	return c.wsClient != nil && c.wsClient.IsConnected()
}

// GetWebSocketStatus 获取WebSocket连接状态信息（主要用于测试）
func (c *Client) GetWebSocketStatus(ctx context.Context) (bool, error) {
	// 尝试确保连接建立
	if err := c.ensureWebSocketConnected(ctx); err != nil {
		return false, err
	}

	return c.IsWebSocketConnected(), nil
}

// CloseWebSocket 关闭WebSocket连接
func (c *Client) CloseWebSocket() error {
	c.wsMu.Lock()
	defer c.wsMu.Unlock()

	if c.wsClient != nil {
		err := c.wsClient.Close()
		c.wsClient = nil
		return err
	}
	return nil
}

// IsWebSocketEnabled 检查是否启用了WebSocket
func (c *Client) IsWebSocketEnabled() bool {
	return c.wsEnabled
}
