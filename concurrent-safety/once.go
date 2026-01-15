package main

import (
	"fmt"
	"sync"
)

var (
	// initOnce 用来保证初始化逻辑只会执行一次
	// 它本身并不关心是谁来执行、什么时候执行
	initOnce sync.Once

	// resource 模拟一个需要初始化后才能安全使用的共享资源
	resource int
)

func DemoOnce() {
	fmt.Println("----------")
	fmt.Println("sync.Once：并发环境下的初始化边界")

	var wg sync.WaitGroup

	// 启动多个 goroutine
	// 它们都会“尝试”使用 resource
	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			// Do 会保证传入的函数只执行一次
			// 但并不保证是哪一个 goroutine 来执行
			initOnce.Do(func() {
				fmt.Printf("resource 初始化发生在 goroutine %d 中\n", id)

				// 初始化逻辑：只会有一次
				resource = 42
			})

			// 这里所有 goroutine 都可以安全地使用 resource
			// Once 不仅保证“只初始化一次”，
			// 也保证初始化完成后的结果对其他 goroutine 可见
			fmt.Println("goroutine", id, "使用 resource，当前值 =", resource)
		}(i)
	}

	// 等待所有 goroutine 执行完成
	wg.Wait()
}
