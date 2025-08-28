# ComfyUI2Go

Go语言的ComfyUI客户端库，支持完整的API调用和WebSocket实时通信。

## 安装

```bash
go get github.com/deferz/comfyui2go
```

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "github.com/deferz/comfyui2go"
)

func main() {
    // 创建客户端
    client := comfyui2go.NewClient("my-app", "http://localhost:8188")
    defer client.CloseWebSocket()
    
    ctx := context.Background()
    
    // 提交工作流
    workflow := comfyui2go.JSON{
        "1": map[string]interface{}{
            "class_type": "CheckpointLoaderSimple",
            "inputs": map[string]interface{}{
                "ckpt_name": "v1-5-pruned-emaonly.ckpt",
            },
        },
    }
    
    promptID, err := client.Prompt(ctx, workflow)
    if err != nil {
        panic(err)
    }
    
    // 等待完成
    result, err := client.WaitForCompletion(ctx, promptID, 2*time.Second)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("任务完成: %s\n", result.PromptID)
}
```

## API文档

### 客户端创建

```go
// 简单创建（默认启用WebSocket）
client := comfyui2go.NewClient("clientID", "http://localhost:8188")

// 完整配置
client := comfyui2go.NewClientWithOptions(
    "clientID", "http://localhost:8188",
    comfyui2go.WithBasicAuth("user", "pass"),     // HTTP认证
    comfyui2go.WithTimeout(30*time.Second),       // 请求超时
    comfyui2go.WithoutWebSocket(),                // 禁用WebSocket
)
```

### 基本API

```go
// 提交工作流
promptID, err := client.Prompt(ctx, workflow)

// 查询队列
queue, err := client.GetQueue(ctx)

// 查询历史
history, err := client.GetHistory(ctx, promptID)

// 中断任务
err := client.Interrupt(ctx)

// 上传图片
err := client.UploadImage(ctx, imageData, "image.png", "input", true)

// 下载文件
data, err := client.Download(ctx, "filename.png", "", "output")
```

### 等待任务完成

```go
// 轮询方式（适用于所有环境）
result, err := client.WaitForCompletion(ctx, promptID, 2*time.Second)

// WebSocket方式（需要启用WebSocket）
result, err := client.WaitForCompletionWithWS(ctx, promptID, 30*time.Second)
```

## WebSocket配置

### 启用/禁用WebSocket

```go
// 默认启用WebSocket
client := comfyui2go.NewClient("app", "http://localhost:8188")

// 禁用WebSocket
client := comfyui2go.NewClientWithOptions(
    "app", "http://localhost:8188",
    comfyui2go.WithoutWebSocket(),
)

// 条件启用
client := comfyui2go.NewClientWithOptions(
    "app", "http://localhost:8188",
    comfyui2go.WithWebSocketEnabled(true), // 或 false
)
```

### WebSocket回调函数

```go
client := comfyui2go.NewClientWithOptions(
    "app", "http://localhost:8188",
    comfyui2go.WithProgressCallback(func(promptID string, progress comfyui2go.WSProgressMessage) {
        percentage := float64(progress.Value) / float64(progress.Max) * 100
        fmt.Printf("进度: %.1f%%\n", percentage)
    }),
    comfyui2go.WithExecutionCallback(func(promptID string, nodeID *string) {
        if nodeID == nil {
            fmt.Printf("任务 %s 完成\n", promptID)
        } else {
            fmt.Printf("执行节点: %s\n", *nodeID)
        }
    }),
)
```

### 批量回调配置

```go
callbacks := comfyui2go.WSCallbackConfig{
    OnProgress: func(promptID string, progress comfyui2go.WSProgressMessage) {
        // 处理进度更新
    },
    OnExecution: func(promptID string, nodeID *string) {
        // 处理执行状态
    },
    OnError: func(promptID string, err error) {
        // 处理错误
    },
}

client := comfyui2go.NewClientWithOptions(
    "app", "http://localhost:8188",
    comfyui2go.WithWebSocketCallbacks(callbacks),
)
```

## 配置选项

### 可用选项

```go
comfyui2go.WithBasicAuth("username", "password")    // HTTP基础认证
comfyui2go.WithTimeout(30*time.Second)              // 请求超时时间
comfyui2go.WithWebSocketEnabled(true)               // 启用/禁用WebSocket
comfyui2go.WithoutWebSocket()                       // 禁用WebSocket的快捷方式
comfyui2go.WithProgressCallback(callback)           // 进度回调
comfyui2go.WithExecutionCallback(callback)          // 执行回调
comfyui2go.WithStatusCallback(callback)             // 状态回调
comfyui2go.WithErrorCallback(callback)              // 错误回调
comfyui2go.WithWebSocketCallbacks(config)           // 批量回调配置
```

### WebSocket状态检查

```go
// 检查WebSocket是否启用
if client.IsWebSocketEnabled() {
    // 使用WebSocket方式
} else {
    // 使用轮询方式
}

// 检查WebSocket连接状态（需要启用WebSocket）
connected, err := client.GetWebSocketStatus(ctx)
```

## 数据类型

### JSON工作流

```go
type JSON = map[string]interface{}

workflow := comfyui2go.JSON{
    "1": map[string]interface{}{
        "class_type": "CheckpointLoaderSimple",
        "inputs": map[string]interface{}{
            "ckpt_name": "model.ckpt",
        },
    },
}
```

## 使用示例

### 启用WebSocket的场景

```go
client := comfyui2go.NewClient("app", "http://localhost:8188")
defer client.CloseWebSocket()

// 实时进度监控
result, err := client.WaitForCompletionWithWS(ctx, promptID, 30*time.Second)
```

### 禁用WebSocket的场景

```go
client := comfyui2go.NewClientWithOptions(
    "app", "http://api-server.com",
    comfyui2go.WithoutWebSocket(),
    comfyui2go.WithBasicAuth("key", "secret"),
)

// 轮询等待
result, err := client.WaitForCompletion(ctx, promptID, 5*time.Second)
```

### 根据状态选择等待方式

```go
if client.IsWebSocketEnabled() {
    result, err = client.WaitForCompletionWithWS(ctx, promptID, timeout)
} else {
    result, err = client.WaitForCompletion(ctx, promptID, pollInterval)
}
```

## 错误处理

```go
promptID, err := client.Prompt(ctx, workflow)
if err != nil {
    log.Printf("提交失败: %v", err)
    return
}

result, err := client.WaitForCompletion(ctx, promptID, 2*time.Second)
if err != nil {
    log.Printf("等待失败: %v", err)
    return
}
```
