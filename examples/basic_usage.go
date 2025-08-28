package main // basic_usage

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/deferz/comfyui2go"
)

func main() {
	// 创建客户端
	client := comfyui2go.NewClientWithOptions(
		"basic-example", "http://localhost:8188",
		// comfyui2go.WithBasicAuth("admin", "password"), // 如果需要认证
	)
	defer client.CloseWebSocket()

	ctx := context.Background()

	// 示例工作流 - 简单的文本到图像生成
	workflow := comfyui2go.JSON{
		"1": map[string]interface{}{
			"class_type": "CheckpointLoaderSimple",
			"inputs": map[string]interface{}{
				"ckpt_name": "v1-5-pruned-emaonly.ckpt",
			},
		},
		"2": map[string]interface{}{
			"class_type": "CLIPTextEncode",
			"inputs": map[string]interface{}{
				"text": "a beautiful sunset over mountains",
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
				"width":      512,
				"height":     512,
				"batch_size": 1,
			},
		},
		"5": map[string]interface{}{
			"class_type": "KSampler",
			"inputs": map[string]interface{}{
				"seed":         42,
				"steps":        20,
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
				"filename_prefix": "ComfyUI",
				"images":          []interface{}{"6", 0},
			},
		},
	}

	fmt.Println("🚀 开始生成图像...")

	// 提交工作流
	promptID, err := client.Prompt(ctx, workflow)
	if err != nil {
		fmt.Printf("❌ 提交工作流失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 工作流已提交: %s\n", promptID)

	// 使用WebSocket等待完成
	result, err := client.WaitForCompletionWithWS(ctx, promptID, 5*time.Minute)
	if err != nil {
		fmt.Printf("❌ 等待完成失败: %v\n", err)
		return
	}

	fmt.Printf("🎉 图像生成完成: %s\n", result.PromptID)

	// 下载生成的图像
	downloadImages(client, result)
}

func downloadImages(client *comfyui2go.Client, result *comfyui2go.WaitResult) {
	ctx := context.Background()

	for nodeName, nodeOutput := range result.Item.Outputs {
		if nodeMap, ok := nodeOutput.(map[string]interface{}); ok {
			if images, ok := nodeMap["images"].([]interface{}); ok {
				for i, img := range images {
					if imgData, ok := img.(map[string]interface{}); ok {
						filename := imgData["filename"].(string)

						fmt.Printf("📥 下载图像: %s (节点: %s)\n", filename, nodeName)

						// 下载图像数据
						data, err := client.Download(ctx, filename, "", "output")
						if err != nil {
							fmt.Printf("❌ 下载失败: %v\n", err)
							continue
						}

						// 保存到当前目录
						outputFile := fmt.Sprintf("generated_%d_%s", i+1, filename)
						if err := ioutil.WriteFile(outputFile, data, 0644); err != nil {
							fmt.Printf("❌ 保存失败: %v\n", err)
							continue
						}

						fmt.Printf("💾 图像已保存: %s\n", outputFile)
					}
				}
			}
		}
	}
}
