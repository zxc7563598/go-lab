package main

import (
	"context"
	"fmt"
	"time"
)

func demoWithCancelTimeoutDeadline() {
	fmt.Println("----------")
	fmt.Println("context.WithCancel / Timeout / Deadline")

	// 1. 由上游显式决定“什么时候放弃”
	contextWithCancel()

	// 2. 由时间自动决定“最多等多久”
	contextTimeout()

	// 3. 由一个明确的时间点决定“到此为止”
	contextDeadline()

}

func contextWithCancel() {
	fmt.Println("WithCancel: 取消的决定权在上游")

	// 创建一个可手动取消的 context
	// ctx：向下游传播“是否还值得继续”的状态
	// cancel：由上游调用，用来声明“可以结束了”
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				// goroutine 并不是被“杀死”的，
				// 而是主动感知到上游已经放弃，于是选择退出
				fmt.Println("goroutine 感知到 cancel() 被调用，主动退出")
				return
			default:
				fmt.Println("goroutine 判断：事情目前仍值得继续")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// 模拟一段正常的业务执行时间
	time.Sleep(2 * time.Second)

	// 上游显式声明：这件事不需要再继续了
	fmt.Println("上游调用 cancel()，发出取消信号")
	cancel()

	// 留出时间观察 goroutine 的退出行为
	time.Sleep(1 * time.Second)
}

func contextTimeout() {
	fmt.Println("WithTimeout: 由时间决定是否继续等待")

	// 创建一个带超时机制的 context
	// 超过 2 秒，如果还没有结束，context 会自动进入 Done 状态
	// cancel 依然存在，是为了提前结束或释放资源
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				// 对 goroutine 来说，它并不关心：是超时触发的，还是 cancel() 被调用
				// 它只看到一件事：ctx.Done() 已经关闭
				fmt.Println("goroutine 感知到超时，认为已经不值得继续等待")
				return
			default:
				fmt.Println("goroutine 正在等待任务完成")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// 主线程等待时间超过 timeout，
	// 用来验证：即使没人手动 cancel，goroutine 也会结束
	time.Sleep(3 * time.Second)
}

func contextDeadline() {
	fmt.Println("WithDeadline: 在某个明确时间点之后直接放弃")

	// 指定一个绝对的截止时间
	// 到达这个时间点后，context 会自动进入 Done 状态
	deadline := time.Now().Add(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("goroutine 发现已超过截止时间，选择退出")
				return
			default:
				fmt.Println("goroutine 正在截止时间之前尝试完成任务")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// 等待超过 deadline，用来验证自动取消行为
	time.Sleep(4 * time.Second)
}
