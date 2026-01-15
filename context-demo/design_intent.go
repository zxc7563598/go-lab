package main

import (
	"context"
	"fmt"
	"time"
)

func demoDesignIntent() {
	fmt.Println("----------")
	fmt.Println("context 的设计初衷：传递“是否还值得继续”这个信号")

	// 创建一个可取消的 context
	// 这里的关键不在于“可取消”，
	// 而在于：它为下游提供了一种统一的方式，
	// 用来判断“这条执行链是否还有效”
	ctx, cancel := context.WithCancel(context.Background())

	// 启动一个 goroutine，模拟一段可能持续运行的工作
	// 注意：goroutine 的生命周期并不会因为外层函数返回而自动结束，
	// 它是否退出，完全取决于自己是否感知并尊重 context 的状态
	go func() {
		for {
			select {
			case <-ctx.Done():
				// 当 ctx.Done() 可读时，说明上游已经放弃了这件事
				// goroutine 并不会被强制终止，
				// 而是“主动意识到已经不值得继续”，然后选择退出
				fmt.Println("goroutine 感知到：上游已放弃，主动结束执行")
				return
			default:
				// 在没有收到取消信号之前，
				// goroutine 会认为这件事仍然值得继续做
				fmt.Println("goroutine 判断：当前任务仍然值得继续执行")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// 模拟一段正常的业务执行时间
	time.Sleep(2 * time.Second)

	// 上游在某个时刻做出决定：
	// 不再需要这条执行链的结果
	// 这一步不会“杀死”任何 goroutine，
	// 只是向下游发出一个“可以结束了”的信号
	fmt.Println("上游决定：这件事已经没有继续的意义了，发出取消信号")
	cancel()

	// 留出时间观察 goroutine 对取消信号的响应
	time.Sleep(1 * time.Second)

	fmt.Println("演示结束：context 本身不做事，只负责传递状态变化")
}
