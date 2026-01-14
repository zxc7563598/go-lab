package main

import (
	"fmt"
	"sync"
	"time"
)

func demoChannelPhilosophy() {
	fmt.Println("----------")
	fmt.Println("channel 的设计哲学")

	var wg sync.WaitGroup
	wg.Add(1)

	// 创建一个无缓冲 channel
	// 无缓冲 channel 的发送和接收必须“同时发生”
	ch := make(chan string)

	// 启动发送者 goroutine
	go func() {
		defer wg.Done()
		fmt.Println("sender: 准备发送消息（此时会阻塞，直到有人接收）")

		// 这里并不是简单地“把数据放进去”
		// 而是: 等待某个接收者出现，完成一次同步
		ch <- "hello"

		// 只有当接收者真正接收完成后，才会继续往下执行
		fmt.Println("sender: 发送完成（说明接收者已经收到）")
	}()

	// 主 goroutine 故意 sleep
	// 用来证明: 即使发送者已经准备好了，
	// 只要接收者没出现，发送就无法完成
	time.Sleep(time.Second)

	fmt.Println("receiver: 准备接收消息")

	// 接收操作也是一次同步点
	// 它会和发送方“同时完成”
	msg := <-ch

	fmt.Println("receiver: 收到消息 ->", msg)

	wg.Wait()
}
