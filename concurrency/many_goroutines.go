package main

import (
	"fmt"
	"sync"
	"time"
)

func demoManyGoroutines() {
	fmt.Println("----------")
	fmt.Println("为什么 Go 可以“随便起协程”")

	var wg sync.WaitGroup
	start := time.Now()

	// 启动大量 goroutine
	//
	// 这里的重点不在于每个 goroutine 做了什么，
	// 而在于：
	// - 创建这么多 goroutine 是否可行
	// - 程序是否还能稳定运行
	// - 是否会立刻出现明显的性能或资源问题
	for i := 0; i < 100000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// 用 Sleep 模拟一个“会被阻塞”的任务
			//
			// 注意：这里不是 CPU 计算，
			// goroutine 在 sleep 期间并不会占用 CPU，
			// runtime 可以把执行机会让给其他 goroutine
			time.Sleep(100 * time.Millisecond)
		}()
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	elapsed := time.Since(start)

	fmt.Printf(
		"启动 100000 个 goroutine 并全部完成，总耗时：%v\n",
		elapsed,
	)
}
