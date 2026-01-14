package main

import "fmt"

// 场景一: 没有接收者的发送
// 本质: 无缓冲 channel 的发送必须“发送者 + 接收者同时就绪”
func deadlockNoReceiver() {
	ch := make(chan int)

	fmt.Println("场景一: 没有接收者的发送")
	fmt.Println("准备向 channel 发送数据")

	// 这里会发生阻塞:
	// 因为这是一个无缓冲 channel，
	// 发送操作必须等待某个 goroutine 来接收
	ch <- 1

	// 永远不会执行到这里
	fmt.Println("发送完成（这一行永远不会打印）")
}

// 场景二: 接收者在等待一个永远不会结束的 channel
// 本质: range 会一直等到 channel 被 close
func deadlockWaitNeverClose() {
	fmt.Println("场景二: 接收者在等待一个永远不会结束的 channel")

	ch := make(chan int)

	go func() {
		fmt.Println("发送者: 发送一个值")
		ch <- 1

		// 注意: 这里没有 close(ch)
		// 发送者退出了，但 channel 仍然是“未关闭状态”
		fmt.Println("发送者: 退出（但 channel 未关闭）")
	}()

	// range ch 的语义是:
	// 1. 只要 channel 没关闭，就继续等待
	// 2. 即使当前没有值，也会阻塞
	for v := range ch {
		fmt.Println("接收者: 收到数据:", v)
	}

	// 因为 channel 永远不会被 close，
	// 所以永远无法跳出 range
	fmt.Println("channel 已结束（这一行永远不会执行）")
}

// 场景三: goroutine 之间互相等待，形成循环依赖
// 本质: 多个 goroutine 在 channel 上形成“闭环等待”
func deadlockMutualWaiting() {
	fmt.Println("场景三: goroutine 之间互相等待，形成循环依赖")

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		fmt.Println("goroutine A: 向 ch1 发送数据")
		ch1 <- 1

		fmt.Println("goroutine A: 等待从 ch2 接收数据")
		<-ch2
	}()

	go func() {
		fmt.Println("goroutine B: 向 ch2 发送数据")
		ch2 <- 1

		fmt.Println("goroutine B: 等待从 ch1 接收数据")
		<-ch1
	}()

	// main goroutine 在这里阻塞，防止程序直接退出
	// 这样可以清楚地观察到死锁状态
	select {}
}

func demoDeadlock() {
	fmt.Println("----------")
	fmt.Println("channel 常见死锁场景演示")

	// 三个场景一次只打开一个运行观察
	deadlockNoReceiver()
	// deadlockWaitNeverClose()
	// deadlockMutualWaiting()
}
