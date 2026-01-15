package main

import (
	"fmt"
	"sync"
)

func DemoDataRace() {
	fmt.Println("----------")
	fmt.Println("data race 是怎么产生的")

	// 这个示例的目的不是“算错了多少”
	// 而是演示：在没有同步手段的情况下，
	// 多个 goroutine 对同一份数据的读写顺序，
	// 不再受 Go 语言和运行时保证
	fmt.Println("预期逻辑结果应为 10000，但程序行为本身是未定义的")

	var counter int
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			for j := 0; j < 1000; j++ {

				// 下面这三步在“逻辑上”看起来是顺序的：
				// 1. 读取 counter
				// 2. 基于旧值计算新值
				// 3. 写回 counter
				//
				// 但在并发环境下，它们并不是一个原子操作，
				// 多个 goroutine 之间可以任意交错执行
				value := counter
				counter = value + 1
			}
		}(i)
	}

	// 等待所有 goroutine 执行完毕
	// 注意：WaitGroup 只能保证“结束时机”，
	// 并不能保证中间过程是并发安全的
	wg.Wait()

	fmt.Printf(
		"最终 counter 值 = %d（从逻辑上推导应为 10000，但这里没有任何顺序保证）\n",
		counter,
	)
}
