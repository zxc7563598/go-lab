package main

import (
	"context"
	"fmt"
	"time"
)

func demoContextLeak() {
	fmt.Println("----------")
	fmt.Println("context 泄漏的隐患")

	// 创建一个可取消的 context
	// 注意：这里我们刻意不保存 ctx，只保留 cancel
	// 用来强调一个事实：
	//   即使上游已经“取消”了 context，
	//   下游 goroutine 如果不监听 ctx.Done()，是完全无感知的
	_, cancel := context.WithCancel(context.Background())

	// 启动一个 goroutine
	// 这个 goroutine 的问题在于：
	//   1. 它是一个无限循环
	//   2. 循环内部完全没有检查 ctx.Done()
	// 换句话说，它的生命周期已经“脱离”了 context
	go func() {
		for {
			fmt.Println("goroutine 正在持续运行（未监听 context，无法被取消）")
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// 模拟一段正常的业务执行时间
	time.Sleep(2 * time.Second)

	// 上游调用 cancel()，表示：
	//   “这件事已经不值得继续了”
	fmt.Println("上游调用 cancel()，试图结束这条 context 链路")
	cancel()

	// 再等待一段时间，观察现象
	// 可以看到：
	//   cancel() 已经被调用，
	//   但 goroutine 仍然在继续运行
	time.Sleep(2 * time.Second)

	fmt.Println("演示结束：goroutine 并不会因为 cancel() 而自动退出")
}
