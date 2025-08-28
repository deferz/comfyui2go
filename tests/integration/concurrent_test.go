package integration

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/deferz/comfyui2go"
	"github.com/deferz/comfyui2go/tests/helpers"
)

// TestConcurrentTasks 测试并发任务处理
func TestConcurrentTasks(t *testing.T) {
	helpers.SkipIntegrationTest(t)

	client := helpers.NewTestClient("concurrent-test")
	defer helpers.CleanupClient(client)

	ctx := context.Background()

	// 读取测试工作流
	workflowData, err := os.ReadFile("../test.json")
	if err != nil {
		t.Skipf("跳过并发测试，找不到 test.json: %v", err)
	}

	var baseWorkflow comfyui2go.JSON
	if err := json.Unmarshal(workflowData, &baseWorkflow); err != nil {
		t.Fatalf("解析工作流JSON失败: %v", err)
	}

	t.Run("多任务并发提交", func(t *testing.T) {
		taskCount := 3
		var wg sync.WaitGroup
		results := make(chan string, taskCount)
		errors := make(chan error, taskCount)

		for i := 0; i < taskCount; i++ {
			wg.Add(1)
			go func(taskID int) {
				defer wg.Done()

				// 创建任务特定的工作流副本
				workflow := make(comfyui2go.JSON)
				for k, v := range baseWorkflow {
					workflow[k] = v
				}

				// 提交任务
				promptID, err := client.Prompt(ctx, workflow)
				if err != nil {
					errors <- err
					return
				}

				t.Logf("📤 任务 %d 已提交: %s", taskID+1, promptID)
				results <- promptID
			}(i)
		}

		// 等待所有任务提交完成
		go func() {
			wg.Wait()
			close(results)
			close(errors)
		}()

		// 收集结果
		var promptIDs []string
		var submitErrors []error

		for {
			select {
			case promptID, ok := <-results:
				if !ok {
					goto done
				}
				promptIDs = append(promptIDs, promptID)
			case err, ok := <-errors:
				if !ok {
					continue
				}
				submitErrors = append(submitErrors, err)
			case <-time.After(30 * time.Second):
				t.Error("并发提交超时")
				goto done
			}
		}

	done:
		t.Logf("🎯 并发提交完成: 成功 %d 个，失败 %d 个", len(promptIDs), len(submitErrors))

		if len(submitErrors) > 0 {
			for i, err := range submitErrors {
				t.Logf("❌ 错误 %d: %v", i+1, err)
			}
		}

		if len(promptIDs) > 0 {
			t.Log("✅ 并发任务提交成功")
		} else {
			t.Error("没有任务提交成功")
		}
	})
}

// TestWebSocketConnectionSharing 测试WebSocket连接共享
func TestWebSocketConnectionSharing(t *testing.T) {
	helpers.SkipIntegrationTest(t)

	client := helpers.NewTestClient("connection-sharing-test")
	defer helpers.CleanupClient(client)

	if !client.IsWebSocketEnabled() {
		t.Skip("WebSocket未启用，跳过连接共享测试")
	}

	ctx := context.Background()

	t.Run("并发WebSocket操作", func(t *testing.T) {
		var wg sync.WaitGroup
		connectCount := 5
		results := make([]bool, connectCount)

		for i := 0; i < connectCount; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				// 并发获取WebSocket状态
				connected, err := client.GetWebSocketStatus(ctx)
				if err != nil {
					t.Errorf("连接 %d 失败: %v", index+1, err)
					return
				}

				results[index] = connected
				t.Logf("🔗 连接 %d 状态: %v", index+1, connected)
			}(i)
		}

		wg.Wait()

		// 检查所有连接是否成功
		successCount := 0
		for _, connected := range results {
			if connected {
				successCount++
			}
		}

		t.Logf("✅ WebSocket连接共享测试: %d/%d 成功", successCount, connectCount)

		if successCount == connectCount {
			t.Log("🎉 所有并发连接都成功，连接共享正常")
		}
	})
}

// TestMixedOperations 测试混合操作
func TestMixedOperations(t *testing.T) {
	helpers.SkipIntegrationTest(t)

	client := helpers.NewTestClient("mixed-ops-test")
	defer helpers.CleanupClient(client)

	ctx := context.Background()

	t.Run("混合API操作", func(t *testing.T) {
		var wg sync.WaitGroup

		// 同时进行队列查询和WebSocket状态检查
		wg.Add(2)

		go func() {
			defer wg.Done()
			queue, err := client.GetQueue(ctx)
			if err != nil {
				t.Errorf("队列查询失败: %v", err)
			} else {
				t.Logf("📋 队列查询成功: %d 个运行中任务", len(queue.QueueRunning))
			}
		}()

		go func() {
			defer wg.Done()
			if client.IsWebSocketEnabled() {
				connected, err := client.GetWebSocketStatus(ctx)
				if err != nil {
					t.Errorf("WebSocket状态查询失败: %v", err)
				} else {
					t.Logf("🔌 WebSocket状态: %v", connected)
				}
			} else {
				t.Log("WebSocket未启用")
			}
		}()

		wg.Wait()
		t.Log("✅ 混合操作测试完成")
	})
}
