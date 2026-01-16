package main

import (
	"fmt"
	"time"
)

// 演示一种最基础、也最隐蔽的 goroutine 泄漏形式：
// - goroutine 中存在一个 for + channel 接收
// - channel 永远不会再有发送方
// - goroutine 就会永久阻塞在 `<-ch` 这一行
func worker(ch chan int) {
	for {
		v := <-ch // 如果再也没有人向 ch 发送数据，这里会永久阻塞
		fmt.Println("worker 收到数据：", v)
	}
}

// 使用 range 读取 channel，看起来比 worker 安全一些
// 但它隐含了一个前提：channel 必须被 close，如果 channel 一直不 close，range 永远不会结束，goroutine 也就永远不会退出
func consumer(ch <-chan int) {
	for v := range ch {
		fmt.Println("consumer 消费数据：", v)
	}
	fmt.Println("consumer 正常退出（只有在 channel 被关闭时才会发生）")
}

// 演示一种“更危险但不容易被察觉”的情况：
// - 使用 select + default
// - 看起来既不阻塞，也没有死锁
// - 但 goroutine 会一直空转，占用 CPU，并且永远不会退出
// 这是“活着”的另一种形式：不阻塞，但也不结束
func busyLoopWithSelect(ch <-chan int) {
	for {
		select {
		case v := <-ch:
			fmt.Println("select 分支收到数据：", v)
		default:
			// 什么都不做，但会立刻进入下一次循环
			// 这是一个典型的忙等（busy loop）
		}
	}
}

// 演示：
// - goroutine 启动了
// - 但 main goroutine 从未向 channel 发送数据
// - worker 永久阻塞在接收操作上
func leakByUnreceivedChannel() {
	fmt.Println("示例一: goroutine 等待 channel 数据，但永远等不到")

	ch := make(chan int)

	go worker(ch)

	// 主 goroutine 什么都不做，只是等待一会儿
	time.Sleep(2 * time.Second)
	fmt.Println("main 退出，但 worker goroutine 仍然被阻塞着")
}

// 演示：
// - goroutine 使用 range 读取 channel
// - 发送完成后却忘记 close
// - consumer goroutine 永远卡在 range 上
func leakByNeverClosedChannel() {
	fmt.Println("示例二: range channel，但忘记 close")

	ch := make(chan int)

	go consumer(ch)

	ch <- 1
	ch <- 2

	// 这里如果不 close(ch)，consumer 永远不会退出
	// close(ch)

	time.Sleep(2 * time.Second)
	fmt.Println("main 退出，但 consumer goroutine 仍在等待 channel 关闭")
}

// 演示：
// - goroutine 使用 select + default
// - 表面上没有阻塞
// - 实际上进入了无限循环
// - goroutine 永远不会自然结束
func leakByBusyLoop() {
	fmt.Println("示例三: select + default 导致 goroutine 空转")

	ch := make(chan int)

	go busyLoopWithSelect(ch)

	time.Sleep(2 * time.Second)
	fmt.Println("main 退出，但 goroutine 仍在无限循环中")
}

func DemoGoroutineLeak() {
	fmt.Println("----------")
	fmt.Println("goroutine 泄漏的几种常见写法演示")

	leakByUnreceivedChannel()
	leakByNeverClosedChannel()
	leakByBusyLoop()
}
