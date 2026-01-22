package main

import (
	"context"
	"log"
	"time"
)

// RunContextAwareWorkerPoolJob 表示一个需要被处理的任务
// 这里仍然只保留最小字段，避免干扰对结构本身的观察
type RunContextAwareWorkerPoolJob struct {
	ID int
}

// runContextAwareWorkerPoolWorker 表示一个“感知上下文”的 worker
//
// 相比固定 worker pool，这里的核心变化是：
// worker 不再只是被动地等待任务结束，
// 而是可以被外部信号主动要求退出。
//
// ctx 控制的是 worker 的“生存权”
// jobs 控制的是任务的“供给”
func runContextAwareWorkerPoolWorker(
	ctx context.Context,
	id int,
	jobs <-chan RunContextAwareWorkerPoolJob,
) {
	for {
		select {
		// 上下文被取消：表示程序层面不再需要这个 worker
		case <-ctx.Done():
			log.Printf("[worker %d] 收到取消信号，准备退出\n", id)
			return

		// 从任务队列中取任务
		case job, ok := <-jobs:
			if !ok {
				// jobs channel 被关闭，且任务已消费完毕
				// 表示“没有新任务了”，但不一定是程序要退出
				log.Printf("[worker %d] 任务通道关闭，正常结束\n", id)
				return
			}

			// 正常处理任务
			log.Printf("[worker %d] 开始处理任务 %d\n", id, job.ID)

			// 模拟耗时操作
			time.Sleep(500 * time.Millisecond)

			log.Printf("[worker %d] 完成任务 %d\n", id, job.ID)
		}
	}
}

// RunContextAwareWorkerPoolDemo 演示一个「可感知生命周期的 worker pool」
//
// 在这个示例中，有两条不同的退出路径：
// 1. jobs 被关闭：表示任务已经全部处理完成
// 2. context 被取消：表示程序层面要求 worker 尽快退出
//
// worker 同时监听这两种信号，从而成为系统生命周期的一部分
func RunContextAwareWorkerPoolDemo() {
	// 创建一个可取消的上下文
	// 它代表“程序还是否希望 worker 继续存在”
	ctx, cancel := context.WithCancel(context.Background())

	// 创建任务通道
	jobs := make(chan RunContextAwareWorkerPoolJob)

	// 启动固定数量的 worker
	workerNum := 3
	for i := 0; i < workerNum; i++ {
		go runContextAwareWorkerPoolWorker(ctx, i, jobs)
	}

	// 提交一批任务
	for i := 0; i < 10; i++ {
		jobs <- RunContextAwareWorkerPoolJob{ID: i}
	}

	// 主动关闭任务通道
	// 告诉 worker：不会再有新任务了
	close(jobs)

	// 等待一段时间，观察部分任务被正常处理
	time.Sleep(2 * time.Second)

	// 取消上下文
	// 这一步不关心还有没有任务未完成
	// 表示程序层面要求 worker 尽快退出
	cancel()

	// 再等待一段时间，确保可以观察到 worker 的退出日志
	time.Sleep(5 * time.Second)
}
