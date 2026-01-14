package main

import (
	"fmt"
	"sync"
)

func demoWhoCloseChannel() {
	fmt.Println("----------")
	fmt.Println("谁负责关闭 channel")

	// ch 用来在多个发送者和接收者之间传递数据
	// 注意: 接收者并不会主动关闭它
	ch := make(chan int)

	var wg sync.WaitGroup

	// 启动多个发送者 goroutine
	// 每个发送者只负责发送数据，
	// 并不知道还有没有其他发送者存在
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			fmt.Println("发送者发送数据:", val)
			ch <- val
		}(i)
	}

	// 启动一个“统一关闭者”
	// 它并不发送具体数据，而是站在“所有发送者整体”的角度
	// 在确认所有发送者都结束后，再关闭 channel
	go func() {
		wg.Wait()
		// 关闭 channel 表示:
		// “不会再有新的数据被发送进来了”
		fmt.Println("所有发送者已完成，准备关闭 channel")
		close(ch)
	}()

	// 接收者通过 range 读取 channel
	// range 会一直读取，直到 channel 被关闭
	for v := range ch {
		fmt.Println("接收者收到数据:", v)
	}

	// 能执行到这里，说明 channel 已经被关闭，
	// 且 channel 中的所有数据都已经被接收完毕
	fmt.Println("接收者确认: channel 已关闭，数据接收完成")
}
