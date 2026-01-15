package main

import (
	"context"
	"fmt"
	"time"
)

func demoWebRequestExit() {
	fmt.Println("----------")
	fmt.Println("Web 请求结束后 goroutine 如何退出")

	// 模拟一个“请求级别”的 context
	// 在真实的 Web 框架中，这个 ctx 通常来自：
	//   r.Context()
	// 请求结束、客户端断开、超时等事件，
	// 都会通过这个 context 向下游传播
	ctx, cancel := context.WithCancel(context.Background())

	// 启动一个与请求相关的 goroutine
	// 关键点不在于“启动了 goroutine”，
	// 而在于：这个 goroutine 是否愿意感知请求的生命周期
	go func() {
		for {
			select {
			case <-ctx.Done():
				// 只有当 goroutine 主动监听 ctx.Done()，
				// 它才能感知到“请求已经结束”这件事
				fmt.Println("goroutine 感知到请求已结束，主动退出")
				return
			default:
				// 如果没有 ctx.Done() 这一支，
				// 即使请求已经结束，goroutine 也会继续执行
				fmt.Println("goroutine 正在处理与请求绑定的业务逻辑")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// 模拟请求正在正常处理
	time.Sleep(2 * time.Second)

	// 模拟 Web 请求结束（比如：响应已经返回给客户端）
	// 注意：这一步并不会“强制杀死”任何 goroutine，
	// 它只是通过 context 发出了一个“可以结束了”的信号
	fmt.Println("模拟场景：Web 请求结束，触发 context cancel()")
	cancel()

	// 再等待一小段时间，观察 goroutine 的行为
	time.Sleep(1 * time.Second)

	fmt.Println("演示结束：goroutine 的退出依赖于是否监听了 context")
}
