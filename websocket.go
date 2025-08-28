package comfyui2go

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
)

// WSClient WebSocket客户端，用于实时接收ComfyUI的进度和状态更新
type WSClient struct {
	conn     *websocket.Conn
	baseURL  string
	clientID string
	username string
	password string

	// 回调函数
	onProgress  ProgressCallback
	onStatus    StatusCallback
	onExecution ExecutionCallback
	onError     ErrorCallback

	// 内部状态
	running bool
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
}

// WSConfig WebSocket客户端配置
type WSConfig struct {
	BaseURL     string            // ComfyUI服务器地址 (如: "http://localhost:8188")
	ClientID    string            // 客户端ID
	Username    string            // 用户名（用于基本认证）
	Password    string            // 密码（用于基本认证）
	OnProgress  ProgressCallback  // 进度回调
	OnStatus    StatusCallback    // 状态回调
	OnExecution ExecutionCallback // 执行状态回调
	OnError     ErrorCallback     // 错误回调
}

// NewWSClient 创建新的WebSocket客户端
func NewWSClient(config WSConfig) *WSClient {
	return &WSClient{
		baseURL:     config.BaseURL,
		clientID:    config.ClientID,
		username:    config.Username,
		password:    config.Password,
		onProgress:  config.OnProgress,
		onStatus:    config.OnStatus,
		onExecution: config.OnExecution,
		onError:     config.OnError,
	}
}

// Connect 连接到ComfyUI WebSocket服务
func (ws *WSClient) Connect(ctx context.Context) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.running {
		return fmt.Errorf("WebSocket客户端已经在运行")
	}

	// 构建WebSocket URL
	wsURL, err := ws.buildWebSocketURL()
	if err != nil {
		return fmt.Errorf("构建WebSocket URL失败: %v", err)
	}

	// 准备连接选项
	opts := &websocket.DialOptions{}

	// 如果有认证信息，添加Basic Auth头
	if ws.username != "" && ws.password != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(ws.username + ":" + ws.password))
		opts.HTTPHeader = map[string][]string{
			"Authorization": {"Basic " + auth},
		}
	}

	// 连接WebSocket
	conn, _, err := websocket.Dial(ctx, wsURL, opts)
	if err != nil {
		return fmt.Errorf("连接WebSocket失败: %v", err)
	}

	ws.conn = conn
	ws.ctx, ws.cancel = context.WithCancel(ctx)
	ws.running = true

	// 启动消息处理循环
	go ws.messageLoop()

	return nil
}

// Close 关闭WebSocket连接
func (ws *WSClient) Close() error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if !ws.running {
		return nil
	}

	ws.running = false

	if ws.cancel != nil {
		ws.cancel()
	}

	if ws.conn != nil {
		return ws.conn.Close(websocket.StatusNormalClosure, "客户端关闭")
	}

	return nil
}

// IsConnected 检查是否已连接
func (ws *WSClient) IsConnected() bool {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.running && ws.conn != nil
}

// buildWebSocketURL 构建WebSocket连接URL
func (ws *WSClient) buildWebSocketURL() (string, error) {
	// 解析基础URL
	baseURL := strings.TrimSuffix(ws.baseURL, "/")

	// 将HTTP(S)协议转换为WS(S)
	if strings.HasPrefix(baseURL, "http://") {
		baseURL = strings.Replace(baseURL, "http://", "ws://", 1)
	} else if strings.HasPrefix(baseURL, "https://") {
		baseURL = strings.Replace(baseURL, "https://", "wss://", 1)
	} else if !strings.HasPrefix(baseURL, "ws://") && !strings.HasPrefix(baseURL, "wss://") {
		// 默认使用ws://
		baseURL = "ws://" + baseURL
	}

	// 构建完整的WebSocket URL
	wsURL := baseURL + "/ws"

	// 添加clientID参数
	if ws.clientID != "" {
		u, err := url.Parse(wsURL)
		if err != nil {
			return "", err
		}
		q := u.Query()
		q.Set("clientId", ws.clientID)
		u.RawQuery = q.Encode()
		wsURL = u.String()
	}

	return wsURL, nil
}

// messageLoop 消息处理循环
func (ws *WSClient) messageLoop() {
	defer func() {
		ws.mu.Lock()
		ws.running = false
		ws.mu.Unlock()
	}()

	for {
		select {
		case <-ws.ctx.Done():
			return
		default:
			// 读取消息
			_, messageData, err := ws.conn.Read(ws.ctx)
			if err != nil {
				if ws.onError != nil {
					ws.onError("", fmt.Errorf("读取WebSocket消息失败: %v", err))
				}
				return
			}

			// 处理消息
			ws.handleMessage(messageData)
		}
	}
}

// handleMessage 处理接收到的消息
func (ws *WSClient) handleMessage(data []byte) {
	var msg WSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		if ws.onError != nil {
			ws.onError("", fmt.Errorf("解析WebSocket消息失败: %v", err))
		}
		return
	}

	switch msg.Type {
	case "status":
		ws.handleStatusMessage(msg.Data)
	case "execution_start":
		ws.handleExecutionStartMessage(msg.Data)
	case "executing":
		ws.handleExecutingMessage(msg.Data)
	case "progress":
		ws.handleProgressMessage(msg.Data)
	case "execution_error":
		ws.handleExecutionErrorMessage(msg.Data)
	case "execution_interrupted":
		ws.handleExecutionInterruptedMessage(msg.Data)
	default:
		// 未知消息类型，可以记录日志但不报错
	}
}

// handleStatusMessage 处理状态消息
func (ws *WSClient) handleStatusMessage(data JSON) {
	if ws.onStatus == nil {
		return
	}

	var statusMsg WSStatusMessage
	if dataBytes, err := json.Marshal(data); err == nil {
		if err := json.Unmarshal(dataBytes, &statusMsg); err == nil {
			// 调用状态回调
			ws.onStatus("", fmt.Sprintf("队列剩余: %d", statusMsg.Status.ExecInfo.QueueRemaining))
		}
	}
}

// handleExecutionStartMessage 处理执行开始消息
func (ws *WSClient) handleExecutionStartMessage(data JSON) {
	if ws.onExecution == nil {
		return
	}

	var execMsg WSExecutionStartMessage
	if dataBytes, err := json.Marshal(data); err == nil {
		if err := json.Unmarshal(dataBytes, &execMsg); err == nil {
			ws.onExecution(execMsg.PromptID, nil) // nil表示开始执行
		}
	}
}

// handleExecutingMessage 处理当前执行节点消息
func (ws *WSClient) handleExecutingMessage(data JSON) {
	if ws.onExecution == nil {
		return
	}

	var execMsg WSExecutingMessage
	if dataBytes, err := json.Marshal(data); err == nil {
		if err := json.Unmarshal(dataBytes, &execMsg); err == nil {
			ws.onExecution(execMsg.PromptID, execMsg.Node)
		}
	}
}

// handleProgressMessage 处理进度消息
func (ws *WSClient) handleProgressMessage(data JSON) {
	if ws.onProgress == nil {
		return
	}

	var progMsg WSProgressMessage
	if dataBytes, err := json.Marshal(data); err == nil {
		if err := json.Unmarshal(dataBytes, &progMsg); err == nil {
			ws.onProgress("", progMsg) // 进度消息通常不包含prompt_id
		}
	}
}

// handleExecutionErrorMessage 处理执行错误消息
func (ws *WSClient) handleExecutionErrorMessage(data JSON) {
	if ws.onError == nil {
		return
	}

	var errMsg WSExecutionErrorMessage
	if dataBytes, err := json.Marshal(data); err == nil {
		if err := json.Unmarshal(dataBytes, &errMsg); err == nil {
			err := fmt.Errorf("节点 %s (%s) 执行错误: %s",
				errMsg.NodeID, errMsg.NodeType, errMsg.Exception)
			ws.onError(errMsg.PromptID, err)
		}
	}
}

// handleExecutionInterruptedMessage 处理执行中断消息
func (ws *WSClient) handleExecutionInterruptedMessage(data JSON) {
	if ws.onError == nil {
		return
	}

	if promptID, ok := data["prompt_id"].(string); ok {
		ws.onError(promptID, fmt.Errorf("任务执行被中断"))
	}
}

// WaitForCompletionWithWS 使用WebSocket等待任务完成（复用Client的WebSocket连接）
func (c *Client) WaitForCompletionWithWS(ctx context.Context, promptID string, timeout time.Duration) (*WaitResult, error) {
	if timeout <= 0 {
		timeout = 5 * time.Minute
	}

	// 获取共享的WebSocket客户端
	wsClient, err := c.GetWebSocketClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取WebSocket连接失败: %v", err)
	}

	return c.waitForCompletionWithExistingWS(ctx, promptID, timeout, wsClient)
}

// waitForCompletionWithExistingWS 使用已有的WebSocket连接等待完成
func (c *Client) waitForCompletionWithExistingWS(ctx context.Context, promptID string, timeout time.Duration, wsClient *WSClient) (*WaitResult, error) {
	resultChan := make(chan *WaitResult, 1)
	errorChan := make(chan error, 1)
	completed := false
	var mu sync.Mutex

	// 临时设置回调函数（注意：这会覆盖之前的回调）
	originalOnExecution := wsClient.onExecution
	originalOnError := wsClient.onError

	// 设置临时回调
	wsClient.onExecution = func(msgPromptID string, nodeID *string) {
		// 先调用原始回调
		if originalOnExecution != nil {
			originalOnExecution(msgPromptID, nodeID)
		}

		// 处理当前任务的完成
		if msgPromptID == promptID && nodeID == nil {
			mu.Lock()
			if !completed {
				completed = true
				// 获取最终结果
				go func() {
					// 稍等一下确保历史记录已更新
					time.Sleep(500 * time.Millisecond)
					history, err := c.GetHistory(ctx, promptID)
					if err != nil {
						errorChan <- err
						return
					}
					if item, ok := history[promptID]; ok {
						resultChan <- &WaitResult{PromptID: promptID, Item: item}
					} else {
						errorChan <- fmt.Errorf("历史记录中未找到 prompt_id: %s", promptID)
					}
				}()
			}
			mu.Unlock()
		}
	}

	wsClient.onError = func(msgPromptID string, err error) {
		// 先调用原始回调
		if originalOnError != nil {
			originalOnError(msgPromptID, err)
		}

		// 处理当前任务的错误
		if msgPromptID == promptID {
			mu.Lock()
			if !completed {
				completed = true
				errorChan <- err
			}
			mu.Unlock()
		}
	}

	// 在函数结束时恢复原始回调
	defer func() {
		wsClient.onExecution = originalOnExecution
		wsClient.onError = originalOnError
	}()

	// 等待结果
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	case <-timeoutCtx.Done():
		return nil, fmt.Errorf("等待任务完成超时")
	}
}

// getBaseURL 获取基础URL
func (c *Client) getBaseURL() string {
	if c.baseURL != "" {
		return c.baseURL
	}
	return "http://localhost:8188" // 默认值
}

// getUsername 获取用户名
func (c *Client) getUsername() string {
	return c.username
}

// getPassword 获取密码
func (c *Client) getPassword() string {
	return c.password
}
