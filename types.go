package comfyui2go

import "time"

// 通用 JSON 类型别名，保持灵活性。
// 许多字段使用 map[string]interface{}，因为 ComfyUI 的节点/工作流结构可能随版本和节点实现而变化。

type JSON = map[string]interface{}

// PromptRequest 表示 POST /prompt 的请求体。
// - Prompt: 工作流图（节点、连接等）的 JSON。
// - ClientID: 可选的客户端标识（与 ComfyUI 约定一致）。
// - ExtraData: 可选的额外字段，某些 ComfyUI 或自定义节点可能支持（如提供则合并入请求体）。
type PromptRequest struct {
	Prompt    JSON   `json:"prompt"`
	ClientID  string `json:"client_id,omitempty"`
	ExtraData JSON   `json:"-"`
}

// PromptResponse 为 POST /prompt 的主要响应。
// 包含 prompt_id、队列编号和节点错误信息。
type PromptResponse struct {
	PromptID   string                 `json:"prompt_id"`
	Number     int                    `json:"number"`
	NodeErrors map[string]interface{} `json:"node_errors"`
}

// QueueResponse 对应 GET /queue 的最小结构。
// 通常包含运行中和等待中的队列列表。
type QueueResponse struct {
	QueueRunning []interface{} `json:"queue_running"`
	QueuePending []interface{} `json:"queue_pending"`
}

// HistoryResponse 对应 GET /history/{prompt_id}。
// 实际返回形如：{ "<id>": HistoryItem }。
// 这里用映射包装，便于后续按需取值。
type HistoryResponse map[string]HistoryItem

// HistoryItem 是一个尽量通用的历史条目视图。
// 关键字段：Status.Completed 或 Outputs 存在时通常表示已完成。
// Outputs 通常包含 node_id -> []OutputAsset 的映射。
// 保持松散结构以兼容不同节点输出。
type HistoryItem struct {
	Status  *HistoryStatus `json:"status,omitempty"`
	Outputs JSON           `json:"outputs,omitempty"`
	// Raw 保留剩余字段，避免信息丢失。
	Raw JSON `json:"-"`
}

type HistoryStatus struct {
	StatusStr string `json:"status_str,omitempty"`
	Completed bool   `json:"completed"`
	// Messages 是复杂结构的数组，包含执行事件信息
	Messages []interface{} `json:"messages,omitempty"`
	// 其他时间字段（实际可能不存在）
	Started  *time.Time `json:"started_at,omitempty"`
	Finished *time.Time `json:"finished_at,omitempty"`
}

// UploadResponse 用于上传相关响应的占位（若后续扩展上传 API）。
type UploadResponse struct {
	Name      string `json:"name,omitempty"`
	Subfolder string `json:"subfolder,omitempty"`
	Type      string `json:"type,omitempty"`
	// extra
	Unknown JSON `json:"-"`
}

// WaitResult 为 WaitForCompletion 的便捷结果类型。
// 当任务完成时包含对应的 HistoryItem。
type WaitResult struct {
	PromptID string      `json:"prompt_id"`
	Item     HistoryItem `json:"item"`
}

// WebSocket 消息类型定义

// WSMessage 表示通过WebSocket接收的消息
type WSMessage struct {
	Type string `json:"type"`
	Data JSON   `json:"data"`
}

// WSStatusMessage 队列状态更新消息
type WSStatusMessage struct {
	Status WSStatusData `json:"status"`
}

type WSStatusData struct {
	ExecInfo WSExecInfo `json:"exec_info"`
}

type WSExecInfo struct {
	QueueRemaining int `json:"queue_remaining"`
}

// WSExecutionStartMessage 任务开始执行消息
type WSExecutionStartMessage struct {
	PromptID string `json:"prompt_id"`
}

// WSExecutingMessage 当前执行节点消息
type WSExecutingMessage struct {
	Node     *string `json:"node"` // null表示执行完成
	PromptID string  `json:"prompt_id"`
}

// WSProgressMessage 进度更新消息
type WSProgressMessage struct {
	Value int `json:"value"`
	Max   int `json:"max"`
}

// WSExecutionErrorMessage 执行错误消息
type WSExecutionErrorMessage struct {
	PromptID  string `json:"prompt_id"`
	NodeID    string `json:"node_id"`
	NodeType  string `json:"node_type"`
	Exception string `json:"exception_message"`
}

// ProgressCallback 进度回调函数类型
type ProgressCallback func(promptID string, progress WSProgressMessage)

// StatusCallback 状态变化回调函数类型
type StatusCallback func(promptID string, status string)

// ExecutionCallback 执行状态回调函数类型
type ExecutionCallback func(promptID string, nodeID *string)

// ErrorCallback 错误回调函数类型
type ErrorCallback func(promptID string, err error)
