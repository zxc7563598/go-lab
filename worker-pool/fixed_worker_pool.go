package main

import (
	"log"
	"time"
)

// FixedWorkerPoolJob 表示一个需要被 worker 处理的任务
// 在这个示例中，任务本身非常简单，只携带一个 ID
// 重点不在任务内容，而在任务是如何被“分发”和“处理”的
type FixedWorkerPoolJob struct {
	ID int
}

// fixedWorkerPoolWorker 表示一个固定存在的 worker
//
// worker 的生命周期由外部控制：
// - 它在 goroutine 中启动
// - 持续从 jobs channel 中取任务
// - 当 channel 被关闭且任务取尽后，自动退出
//
// 注意：
// worker 并不知道一共有多少个任务，也不知道是否还有其他 worker 存在
// 它只负责一件事：顺序地处理自己拿到的任务
func fixedWorkerPoolWorker(id int, jobs <-chan FixedWorkerPoolJob) {
	for job := range jobs {
		// 打印 worker 开始处理任务
		// 通过输出可以观察到：同一时间最多只有固定数量的 worker 在工作
		log.Printf("[worker %d] 开始处理任务 %d\n", id, job.ID)

		// 模拟真实业务中的耗时操作
		time.Sleep(500 * time.Millisecond)

		// 打印 worker 完成任务
		log.Printf("[worker %d] 完成任务 %d\n", id, job.ID)
	}

	// 当 jobs channel 被关闭且任务消费完毕后，会走到这里
	// 表示这个 worker 的生命周期自然结束
	log.Printf("[worker %d] 退出\n", id)
}

// FixedWorkerPoolDemo 演示一个「固定大小的 worker pool」
//
// 核心变化在于：
// - worker 的数量是提前确定的
// - 任务通过 channel 排队进入 pool
// - 并发度由 worker 数量决定，而不是由任务数量决定
func FixedWorkerPoolDemo() {
	// 创建任务 channel
	// 这里使用无缓冲 channel，任务的发送和接收会自然形成节奏
	jobs := make(chan FixedWorkerPoolJob)

	// 定义 worker 数量
	// 无论任务有多少，同时工作的 goroutine 数量最多为 workerNum
	workerNum := 3

	// 启动固定数量的 worker
	for i := 0; i < workerNum; i++ {
		go fixedWorkerPoolWorker(i, jobs)
	}

	// 连续提交多个任务
	// 即使这里快速提交 10 个任务，也不会产生 10 个并发 worker
	for i := 0; i < 10; i++ {
		jobs <- FixedWorkerPoolJob{ID: i}
	}

	// 关闭任务 channel
	// 告诉所有 worker：不会再有新的任务了
	close(jobs)

	// 主 goroutine 等待一段时间
	// 仅用于演示，保证可以看到所有 worker 的输出
	time.Sleep(3 * time.Second)
}
