package main // concurrent_tasks

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/deferz/comfyui2go"
)

func main() {
	// 创建客户端
	client := comfyui2go.NewClientWithOptions(
		"concurrent-tasks", "http://localhost:8188",
		// comfyui2go.WithBasicAuth("admin", "password"), // 如果需要认证
	)
	defer client.CloseWebSocket()

	// 准备多个工作流变体
	prompts := []string{
		"a cat sitting on a chair",
		"a dog running in the park",
		"a bird flying in the sky",
	}

	fmt.Printf("🚀 开始并发处理 %d 个任务...\n", len(prompts))

	// 使用 WaitGroup 来等待所有任务完成
	var wg sync.WaitGroup
	results := make(chan string, len(prompts))
	errors := make(chan error, len(prompts))

	for i, prompt := range prompts {
		wg.Add(1)
		go func(taskID int, promptText string) {
			defer wg.Done()

			err := processTask(client, taskID, promptText)
			if err != nil {
				errors <- fmt.Errorf("任务 %d 失败: %v", taskID, err)
			} else {
				results <- fmt.Sprintf("任务 %d 完成", taskID)
			}
		}(i+1, prompt)
	}

	// 等待所有任务完成
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// 收集结果
	var completed, failed int
	for {
		select {
		case result, ok := <-results:
			if !ok && len(errors) == 0 {
				goto done
			}
			if ok {
				fmt.Printf("✅ %s\n", result)
				completed++
			}
		case err, ok := <-errors:
			if !ok {
				continue
			}
			fmt.Printf("❌ %v\n", err)
			failed++
		}
	}

done:
	fmt.Printf("🎯 处理完成! 成功: %d, 失败: %d\n", completed, failed)
	fmt.Println("📊 注意: 所有任务共享同一个WebSocket连接")
}

func processTask(client *comfyui2go.Client, taskID int, promptText string) error {
	ctx := context.Background()

	// 创建工作流（简化版本）
	workflow := createWorkflow(promptText)

	fmt.Printf("📤 任务 %d: 提交工作流 '%s'\n", taskID, promptText)

	// 提交工作流
	promptID, err := client.Prompt(ctx, workflow)
	if err != nil {
		return fmt.Errorf("提交失败: %v", err)
	}

	// 使用WebSocket等待完成（连接会自动复用）
	result, err := client.WaitForCompletionWithWS(ctx, promptID, 3*time.Minute)
	if err != nil {
		return fmt.Errorf("执行失败: %v", err)
	}

	fmt.Printf("🎉 任务 %d: 生成完成 (%s)\n", taskID, result.PromptID)
	return nil
}

func createWorkflow(promptText string) comfyui2go.JSON {
	return comfyui2go.JSON{
		"1": map[string]interface{}{
			"class_type": "CheckpointLoaderSimple",
			"inputs": map[string]interface{}{
				"ckpt_name": "v1-5-pruned-emaonly.ckpt",
			},
		},
		"2": map[string]interface{}{
			"class_type": "CLIPTextEncode",
			"inputs": map[string]interface{}{
				"text": promptText,
				"clip": []interface{}{"1", 1},
			},
		},
		"3": map[string]interface{}{
			"class_type": "CLIPTextEncode",
			"inputs": map[string]interface{}{
				"text": "bad quality, blurry",
				"clip": []interface{}{"1", 1},
			},
		},
		"4": map[string]interface{}{
			"class_type": "EmptyLatentImage",
			"inputs": map[string]interface{}{
				"width":      256,
				"height":     256,
				"batch_size": 1,
			},
		},
		"5": map[string]interface{}{
			"class_type": "KSampler",
			"inputs": map[string]interface{}{
				"seed":         42,
				"steps":        10, // 减少步数以加快速度
				"cfg":          7.0,
				"sampler_name": "euler",
				"scheduler":    "normal",
				"denoise":      1.0,
				"model":        []interface{}{"1", 0},
				"positive":     []interface{}{"2", 0},
				"negative":     []interface{}{"3", 0},
				"latent_image": []interface{}{"4", 0},
			},
		},
		"6": map[string]interface{}{
			"class_type": "VAEDecode",
			"inputs": map[string]interface{}{
				"samples": []interface{}{"5", 0},
				"vae":     []interface{}{"1", 2},
			},
		},
		"7": map[string]interface{}{
			"class_type": "SaveImage",
			"inputs": map[string]interface{}{
				"filename_prefix": "concurrent_test",
				"images":          []interface{}{"6", 0},
			},
		},
	}
}
