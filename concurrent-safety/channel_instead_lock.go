package main

import (
	"fmt"
	"time"
)

func DemoChannelInsteadOfLock() {
	fmt.Println("----------")
	fmt.Println("用 channel 替代锁：让状态只属于一个 goroutine")

	// command 用来描述「我想对状态做什么」
	// 而不是「我想怎么改这个变量」
	type command int

	const (
		inc command = iota // 表示一次递增请求
		get                // 表示一次读取请求
	)

	// cmdCh 用来发送“操作意图”
	cmdCh := make(chan command)

	// respCh 用来返回读取结果
	// 这里单独拆出来，是为了强调：
	// 状态不共享，结果通过通信返回
	respCh := make(chan int)

	// 这个 goroutine 是 counter 的唯一拥有者
	// 所有对 counter 的读写，都被串行化在这里
	go func() {
		var counter int

		for cmd := range cmdCh {
			switch cmd {

			case inc:
				// 只有这个 goroutine 能修改 counter
				counter++

			case get:
				// 读取结果通过 channel 返回
				respCh <- counter
			}
		}
	}()

	// 模拟多个并发请求
	// 它们并不直接接触 counter
	// 只能通过发送命令来“间接影响状态”
	for i := 0; i < 3; i++ {
		go func(id int) {
			fmt.Println("goroutine", id, "请求递增 counter")
			cmdCh <- inc
		}(i)
	}

	// 等待递增请求被处理
	time.Sleep(100 * time.Millisecond)

	// 发起一次读取请求
	cmdCh <- get

	// 接收读取结果
	fmt.Println("最终 counter 值 =", <-respCh)
}
