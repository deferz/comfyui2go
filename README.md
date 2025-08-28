# ComfyUI2Go

## 📖 简介

ComfyUI2Go是一个用于与ComfyUI服务器进行交互的Go语言客户端库。它提供了简洁的API来提交工作流、监控执行进度、下载结果等功能。

## 🛠 安装

```bash
go get github.com/deferz/comfyui2go
```

## 🚀 快速开始

### 基本使用

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/deferz/comfyui2go"
)

func main() {
    // 1. 创建客户端
    client := comfyui2go.New(
        comfyui2go.WithBaseURL("http://localhost:8188"),
        comfyui2go.WithBasicAuth("admin", "password"), // 可选
    )
    defer client.CloseWebSocket() // 可选：手动关闭WebSocket连接
    
    ctx := context.Background()
    
    // 2. 准备工作流（从文件或直接定义）
    workflow := comfyui2go.JSON{
        "1": map[string]interface{}{
            "class_type": "CheckpointLoaderSimple",
            "inputs": map[string]interface{}{
                "ckpt_name": "v1-5-pruned-emaonly.ckpt",
            },
        },
        // ... 更多节点
    }
    
    // 3. 核心工作流程（仅需2行）
    promptID, err := client.Prompt(ctx, workflow)
    if err != nil {
        panic(err)
    }
    fmt.Printf("✅ 工作流已提交: %s\n", promptID)
    
    // 使用WebSocket等待完成（自动管理连接）
    result, err := client.WaitForCompletionWithWS(ctx, promptID, 5*time.Minute)
    if err != nil {
        panic(err)
    }
    fmt.Printf("🎉 生成完成: %s\n", result.PromptID)
    
    // 4. 下载结果（可选）
    // 遍历输出节点，下载生成的图像
    for nodeName, nodeOutput := range result.Item.Outputs {
        if nodeMap, ok := nodeOutput.(map[string]interface{}); ok {
            if images, ok := nodeMap["images"].([]interface{}); ok {
                for _, img := range images {
                    if imgData, ok := img.(map[string]interface{}); ok {
                        filename := imgData["filename"].(string)
                        
                        // 下载图像
                        data, err := client.Download(ctx, filename, "", "output")
                        if err != nil {
                            fmt.Printf("❌ 下载失败: %v\n", err)
                            continue
                        }
                        
                        // 保存到本地
                        if err := os.WriteFile(filename, data, 0644); err != nil {
                            fmt.Printf("❌ 保存失败: %v\n", err)
                            continue
                        }
                        
                        fmt.Printf("💾 图像已保存: %s (节点: %s)\n", filename, nodeName)
                    }
                }
            }
        }
    }
}
```

## 📚 API参考

### 创建客户端

```go
// 基本客户端
client := comfyui2go.New()

// 带配置选项和回调函数
client := comfyui2go.New(
    comfyui2go.WithBaseURL("http://your-server:8188"),
    comfyui2go.WithBasicAuth("username", "password"),
    comfyui2go.WithClientID("my-client-id"),
    
    // 配置WebSocket回调函数
    comfyui2go.WithProgressCallback(func(promptID string, progress comfyui2go.WSProgressMessage) {
        percentage := float64(progress.Value) / float64(progress.Max) * 100
        fmt.Printf("📊 任务 %s 进度: %.1f%%\n", promptID, percentage)
    }),
    comfyui2go.WithExecutionCallback(func(promptID string, nodeID *string) {
        if nodeID == nil {
            fmt.Printf("🎉 任务 %s 执行完成\n", promptID)
        } else {
            fmt.Printf("⚙️ 任务 %s 正在执行节点: %s\n", promptID, *nodeID)
        }
    }),
)
```

### 核心API

#### 1. 提交工作流

```go
promptID, err := client.Prompt(ctx, workflow)
```

#### 2. 等待完成

```go
// 使用WebSocket实时监控（推荐）
result, err := client.WaitForCompletionWithWS(ctx, promptID, timeout)

// 使用轮询方式
result, err := client.WaitForCompletion(ctx, promptID, pollInterval)
```

#### 3. 查询状态

```go
// 查询队列
queue, err := client.GetQueue(ctx)
fmt.Printf("执行中: %d, 待处理: %d\n", 
    len(queue.QueueRunning), len(queue.QueuePending))

// 查询历史
history, err := client.GetHistory(ctx, promptID)
```

#### 4. 管理任务

```go
// 中断所有任务
err := client.Interrupt(ctx)

// 上传图像
err := client.UploadImage(ctx, imageData, "image.png", "", "input")

// 下载结果
data, err := client.Download(ctx, filename, subfolder, filetype)
```

### WebSocket配置

ComfyUI2Go支持可选的WebSocket功能，适应不同的部署环境：

#### WebSocket启用/禁用

```go
// 默认启用WebSocket（推荐）
client := comfyui2go.NewClient("my-app", "http://localhost:8188")

// 明确禁用WebSocket（开放平台场景）
client := comfyui2go.NewClientWithOptions(
    "my-app", "https://api.platform.com",
    comfyui2go.WithoutWebSocket(), // 禁用WebSocket
)

// 条件性启用WebSocket
client := comfyui2go.NewClientWithOptions(
    "my-app", "http://localhost:8188",
    comfyui2go.WithWebSocketEnabled(hasWebSocketSupport),
)
```

#### 等待方式选择

```go
// WebSocket方式（实时，推荐）- 需要启用WebSocket
result, err := client.WaitForCompletionWithWS(ctx, promptID, timeout)

// 轮询方式（兼容性好）- 适用于所有环境
result, err := client.WaitForCompletion(ctx, promptID, pollInterval)
```

### WebSocket回调配置

当启用WebSocket时，您可以配置回调函数来实时监控任务执行状态：

#### 方式1：分别配置回调函数

```go
client := comfyui2go.New(
    comfyui2go.WithBaseURL("http://localhost:8188"),
    
    // 进度回调
    comfyui2go.WithProgressCallback(func(promptID string, progress comfyui2go.WSProgressMessage) {
        percentage := float64(progress.Value) / float64(progress.Max) * 100
        fmt.Printf("📊 任务 %s 进度: %.1f%% (%d/%d)\n", 
            promptID, percentage, progress.Value, progress.Max)
    }),
    
    // 执行状态回调
    comfyui2go.WithExecutionCallback(func(promptID string, nodeID *string) {
        if nodeID == nil {
            fmt.Printf("🎉 任务 %s 执行完成\n", promptID)
        } else {
            fmt.Printf("⚙️ 任务 %s 正在执行节点: %s\n", promptID, *nodeID)
        }
    }),
    
    // 错误回调
    comfyui2go.WithErrorCallback(func(promptID string, err error) {
        fmt.Printf("❌ 任务 %s 出错: %v\n", promptID, err)
    }),
)
```

#### 方式2：批量配置回调函数

```go
client := comfyui2go.New(
    comfyui2go.WithBaseURL("http://localhost:8188"),
    comfyui2go.WithWebSocketCallbacks(comfyui2go.WSCallbackConfig{
        OnProgress: func(promptID string, progress comfyui2go.WSProgressMessage) {
            percentage := float64(progress.Value) / float64(progress.Max) * 100
            fmt.Printf("📈 进度: %.1f%%\n", percentage)
        },
        OnExecution: func(promptID string, nodeID *string) {
            if nodeID == nil {
                fmt.Println("✅ 任务完成")
            }
        },
        OnStatus: func(promptID string, status string) {
            fmt.Printf("📋 状态: %s\n", status)
        },
        OnError: func(promptID string, err error) {
            fmt.Printf("🔥 错误: %v\n", err)
        },
    }),
)
```

#### 回调函数类型说明

- **`ProgressCallback`**: 任务执行进度更新
- **`ExecutionCallback`**: 节点执行状态变化（开始/完成）
- **`StatusCallback`**: 队列状态变化
- **`ErrorCallback`**: 执行错误

### WebSocket管理

```go
// 获取WebSocket客户端（自动连接，包含配置的回调）
wsClient, err := client.GetWebSocketClient(ctx)

// 检查连接状态
if wsClient.IsConnected() {
    fmt.Println("✅ WebSocket已连接")
}

// 手动关闭连接
err := client.CloseWebSocket()
```

## 🎯 使用场景

### 场景1：本地ComfyUI（推荐WebSocket）

```go
client := comfyui2go.NewClientWithOptions(
    "my-app", "http://localhost:8188",
    // 默认启用WebSocket，获得最佳性能
    comfyui2go.WithProgressCallback(func(promptID string, progress comfyui2go.WSProgressMessage) {
        fmt.Printf("进度: %.1f%%\n", float64(progress.Value)/float64(progress.Max)*100)
    }),
)
```

### 场景2：禁用WebSocket

```go
client := comfyui2go.NewClientWithOptions(
    "openapi-client", "https://api.comfyui-platform.com",
    comfyui2go.WithoutWebSocket(),                    
    comfyui2go.WithBasicAuth("api-key", "api-secret"), 
    comfyui2go.WithTimeout(60*time.Second),          
)

// 使用轮询方式等待完成
result, err := client.WaitForCompletion(ctx, promptID, 5*time.Second)
```

### 场景3：云服务（明确选择）

```go
client := comfyui2go.NewClientWithOptions(
    "cloud-app", cloudEndpoint,
    comfyui2go.WithBasicAuth(apiKey, apiSecret),
    comfyui2go.WithWebSocketEnabled(cloudSupportsWebSocket),
)

// 根据WebSocket启用状态选择等待方式
if client.IsWebSocketEnabled() {
    result, err = client.WaitForCompletionWithWS(ctx, promptID, timeout)
} else {
    result, err = client.WaitForCompletion(ctx, promptID, 3*time.Second)
}
```

### 场景4：条件性启用

```go
// 根据环境变量决定是否启用WebSocket
wsEnabled := os.Getenv("ENABLE_WEBSOCKET") == "true"

client := comfyui2go.NewClientWithOptions(
    "conditional-app", endpoint,
    comfyui2go.WithWebSocketEnabled(wsEnabled),
)
```

## 🎯 高级用法

### 并发处理

```go
func processMultipleTasks(client *comfyui2go.Client, workflows []comfyui2go.JSON) {
    var wg sync.WaitGroup
    results := make(chan string, len(workflows))
    
    for i, workflow := range workflows {
        wg.Add(1)
        go func(taskID int, wf comfyui2go.JSON) {
            defer wg.Done()
            
            ctx := context.Background()
            
            // 提交工作流
            promptID, err := client.Prompt(ctx, wf)
            if err != nil {
                fmt.Printf("❌ 任务 %d 提交失败: %v\n", taskID, err)
                return
            }
            
            // WebSocket连接会自动复用
            result, err := client.WaitForCompletionWithWS(ctx, promptID, 5*time.Minute)
            if err != nil {
                fmt.Printf("❌ 任务 %d 执行失败: %v\n", taskID, err)
                return
            }
            
            results <- result.PromptID
            fmt.Printf("✅ 任务 %d 完成: %s\n", taskID, result.PromptID)
        }(i+1, workflow)
    }
    
    wg.Wait()
    close(results)
    
    fmt.Printf("🎯 所有任务完成，共处理 %d 个工作流\n", len(workflows))
}
```

### 自定义WebSocket回调

```go
// 获取WebSocket客户端
wsClient, err := client.GetWebSocketClient(ctx)
if err != nil {
    panic(err)
}

// 设置自定义回调（注意：这会覆盖默认行为）
wsClient.OnProgress = func(promptID string, progress comfyui2go.WSProgressMessage) {
    percentage := float64(progress.Value) / float64(progress.Max) * 100
    fmt.Printf("📊 任务 %s 进度: %.1f%% (%d/%d)\n", 
        promptID, percentage, progress.Value, progress.Max)
}

wsClient.OnExecution = func(promptID string, nodeID *string) {
    if nodeID == nil {
        fmt.Printf("🎉 任务 %s 执行完成\n", promptID)
    } else {
        fmt.Printf("⚙️ 任务 %s 正在执行节点: %s\n", promptID, *nodeID)
    }
}
```

### 错误处理

```go
result, err := client.WaitForCompletionWithWS(ctx, promptID, timeout)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "timeout"):
        fmt.Println("⏰ 任务执行超时")
    case strings.Contains(err.Error(), "WebSocket"):
        fmt.Println("🔌 WebSocket连接问题")
    case strings.Contains(err.Error(), "中断"):
        fmt.Println("🛑 任务被中断")
    default:
        fmt.Printf("❌ 未知错误: %v\n", err)
    }
}
```
