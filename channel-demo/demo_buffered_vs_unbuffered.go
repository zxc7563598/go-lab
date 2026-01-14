package main

import (
	"fmt"
	"sync"
	"time"
)

func demoBufferedVsUnbuffered() {
	fmt.Println("----------")
	fmt.Println("无缓冲 channel vs 有缓冲 channel")

	var wg sync.WaitGroup
	wg.Add(1)

	// 无缓冲 channel 的 send 和 receive 必须同时发生，
	// 否则其中一方就会阻塞等待。
	fmt.Println("「无缓冲 channel」")

	chUnbuffered := make(chan int)

	go func() {
		defer wg.Done()
		fmt.Println("无缓冲 channel: 准备发送 1")
		chUnbuffered <- 1
		// 只有当主 goroutine 完成接收后，
		// 这一行才会被执行
		fmt.Println("无缓冲 channel: 发送完成")
		close(chUnbuffered)
	}()

	// 故意 sleep 一下，确保发送方先执行到 send 位置
	time.Sleep(time.Millisecond * 100)

	fmt.Println("无缓冲 channel: 主 goroutine 准备接收")
	fmt.Println("无缓冲 channel: 接收到的值 =", <-chUnbuffered)
	wg.Wait()

	// 有缓冲 channel 允许在缓冲区未满的情况下，
	// 发送方先继续执行，而不必等待接收方。
	fmt.Println("「有缓冲 channel」")

	chBuffered := make(chan int, 2)

	fmt.Println("有缓冲 channel: 准备发送 1")
	chBuffered <- 1

	fmt.Println("有缓冲 channel: 准备发送 2")
	chBuffered <- 2

	// 因为缓冲区容量是 2，此时并不会阻塞
	fmt.Println("有缓冲 channel: 发送完成（未发生阻塞）")

	close(chBuffered)

	fmt.Println("有缓冲 channel: 接收到的值 =", <-chBuffered)
	fmt.Println("有缓冲 channel: 接收到的值 =", <-chBuffered)
}
