# ComfyUI2Go 示例代码

这个目录包含了一些实用的示例代码，展示如何使用ComfyUI2Go库。

## 📂 示例列表

### 1. basic_usage.go
**基本使用示例**
- 演示最基本的工作流提交和等待完成
- 包含图像下载和保存
- 适合初学者了解基本API

运行方式：
```bash
cd examples
go run basic_usage.go
```

### 2. concurrent_tasks.go
**并发任务处理示例**
- 演示多个任务并发执行
- 展示WebSocket连接复用的优势
- 适合了解高性能处理模式

运行方式：
```bash
cd examples
go run concurrent_tasks.go
```

## 🔧 运行前准备

1. **确保ComfyUI服务器运行**
   ```bash
   # 确保ComfyUI在 http://localhost:8188 运行
   # 或修改示例代码中的URL
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **检查模型文件**
   - 确保有 `v1-5-pruned-emaonly.ckpt` 模型
   - 或修改示例中的模型名称

## 📝 自定义示例

您可以基于这些示例创建自己的工作流：

1. **修改prompt文本**
   ```go
   "text": "your custom prompt here",
   ```

2. **调整图像尺寸**
   ```go
   "width":  1024,
   "height": 1024,
   ```

3. **修改采样参数**
   ```go
   "steps": 30,
   "cfg":   8.0,
   ```

## 🚨 注意事项

- 确保ComfyUI服务器有足够的GPU内存
- 并发任务数量不要超过GPU处理能力
- 网络超时时间根据任务复杂度调整
- 生成的图像会保存在当前目录

## 📖 更多资源

- [主要README](../README.md) - 完整的API文档
- [ComfyUI官方文档](https://github.com/comfyanonymous/ComfyUI)
- [工作流格式说明](https://github.com/comfyanonymous/ComfyUI/wiki)
