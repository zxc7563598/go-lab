package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// WebhookPayload 定义接收的数据结构
type WebhookPayload struct {
	Event     string         `json:"event"`
	Timestamp time.Time      `json:"timestamp"`
	Data      map[string]any `json:"data"`
	Signature string         `json:"signature"`
}

// WebhookProcessor 表示一个典型的“生产者-消费者”模型：
// HTTP 层作为生产者，worker 池作为消费者。
// 这个结构体本身并不关心 HTTP，只关心如何并发、安全地处理任务。
type WebhookProcessor struct {
	// queue 是任务缓冲队列，用来承接 HTTP 请求与后台处理之间的速率差
	// 使用带缓冲的 channel，可以避免 webhook 高峰期直接把服务压垮
	queue chan WebhookPayload

	// workers 表示 worker goroutine 的数量
	// 在这个模型里，它同时也代表了最大并发处理能力
	workers int

	// wg 用来等待所有 worker 正常退出
	// 这使得 Stop() 可以是一个“阻塞直到完全停止”的操作
	wg sync.WaitGroup

	// stopChan 用来广播“停止信号”
	// 一旦关闭，所有 worker 都应该开始走退出路径
	stopChan chan struct{}

	// processing 作为一个信号量（semaphore），用于限制同时处理的任务数
	// 在本例中，由于 worker 数量已经限制了并发度，这个 channel 实际上是演示用途
	processing chan struct{}
}

// NewWebhookProcessor 创建新的处理器
func NewWebhookProcessor(workers int, queueSize int) *WebhookProcessor {
	return &WebhookProcessor{
		queue:      make(chan WebhookPayload, queueSize),
		workers:    workers,
		stopChan:   make(chan struct{}),
		processing: make(chan struct{}, workers),
	}
}

// Start 启动 worker 池。
// 每个 worker 都是一个独立的 goroutine，
// 它们会持续从 queue 中取任务，直到收到 stop 信号。
func (wp *WebhookProcessor) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	fmt.Printf("启动 %d 个 worker 处理 webhook\n", wp.workers)
}

// Stop 用来通知所有 worker 停止工作，并等待它们退出。
// 这里并不会强制中断正在处理的任务，而是给 worker 一个“该结束了”的信号。
func (wp *WebhookProcessor) Stop() {
	close(wp.stopChan)
	wp.wg.Wait()
	fmt.Printf("所有 worker 已停止\n")
}

// worker 是实际执行 webhook 处理逻辑的 goroutine。
// 每个 worker 会在一个循环中等待任务或停止信号。
func (wp *WebhookProcessor) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		// 一旦 stopChan 被关闭，所有 worker 都会走到这里并退出
		case <-wp.stopChan:
			fmt.Printf("Worker %d 收到停止信号，退出\n", id)
			return

		// 从任务队列中取出一个 payload
		// 如果此时队列为空，worker 会阻塞在这里等待新任务
		case payload := <-wp.queue:
			// 通过向 processing channel 写入一个值，表示“占用一个处理名额”
			wp.processing <- struct{}{}

			fmt.Printf("Worker %d 开始处理事件: %s\n", id, payload.Event)

			// 实际的业务处理逻辑
			wp.processPayload(id, payload)

			// 释放处理名额
			<-wp.processing
		}
	}
}

// AddPayload 尝试将 webhook payload 放入处理队列。
// 这里使用非阻塞写入，是为了避免 HTTP 请求被无限期卡住。
func (wp *WebhookProcessor) AddPayload(payload WebhookPayload) error {
	select {
	case wp.queue <- payload:
		fmt.Printf("事件 %s 已加入队列\n", payload.Event)
		return nil
	default:
		// 当队列已满时，直接返回错误，由 HTTP 层决定如何响应
		return fmt.Errorf("处理队列已满")
	}
}

// processPayload 实际处理webhook的业务逻辑
func (wp *WebhookProcessor) processPayload(workerID int, payload WebhookPayload) {
	// 模拟随机处理时间
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := r.Intn(5) + 1
	time.Sleep(time.Second * time.Duration(randomNum))
	// 打印payload
	dataJSON, _ := json.MarshalIndent(payload.Data, "", "  ")
	fmt.Printf("Worker %d 处理完成 - 事件: %s, 时间: %s\n数据: %s\n", workerID, payload.Event, payload.Timestamp.Format(time.DateTime), string(dataJSON))
}

// healthCheck 检查接口是否正常
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "healthy",
		"time":   time.Now().Format(time.DateTime),
	})
}

// webhookHandler 是 HTTP 层与处理器之间的“边界”。
// 它的职责非常明确：
// 1. 校验请求
// 2. 解析数据
// 3. 尝试入队
// 4. 尽快返回响应
func webhookHandler(processor *WebhookProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "只支持 POST 请求", http.StatusMethodNotAllowed)
			return
		}

		// 这里对 Content-Type 的判断是简化版本
		// 真实环境中可能需要兼容 charset 等参数
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "只支持 application/json", http.StatusUnsupportedMediaType)
			return
		}

		var payload WebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "无效的 JSON 格式", http.StatusBadRequest)
			return
		}

		// 补充一些基础字段，避免后续处理时出现空值
		if payload.Timestamp.IsZero() {
			payload.Timestamp = time.Now()
		}
		if payload.Event == "" {
			payload.Event = "未知"
		}

		// 尝试将任务交给后台处理器
		if err := processor.AddPayload(payload); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		// webhook 的最佳实践通常是“尽快确认已接收”
		w.WriteHeader(http.StatusAccepted)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "accepted",
			"message": "webhook 已接收",
			"event":   payload.Event,
		})
	}
}

// metricsHandler 监控信息：显示队列状态
func metricsHandler(processor *WebhookProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"queue_length":     len(processor.queue),
			"queue_capacity":   cap(processor.queue),
			"processing_count": processor.workers,
			"timestamp":        time.Now().Format(time.DateTime),
		})
	}
}

func main() {
	// 配置参数
	port := ":8080"
	workers := 5
	queueSize := 100

	// 创建webhook处理器
	processor := NewWebhookProcessor(workers, queueSize)

	// 启动worker池
	processor.Start()
	defer processor.Stop()

	// 设置http路由
	http.HandleFunc("/webhook", webhookHandler(processor))
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/metrics", metricsHandler(processor))

	// 启动http服务
	fmt.Printf("http服务启动在: http://localhost%s\n", port)
	fmt.Println("可用端点:")
	fmt.Println("  POST /webhook - 接收webhook")
	fmt.Println("  GET  /health  - 健康检查")
	fmt.Println("  GET  /metrics - 查看队列状态")
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("服务器启动失败", err)
	}

}
