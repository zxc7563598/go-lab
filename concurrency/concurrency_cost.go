package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// 模拟 CPU 密集型任务
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func demoConcurrencyCost() {
	fmt.Println("----------")
	fmt.Println("并发 ≠ 更快")
	fmt.Println("逻辑CPU核数:", runtime.NumCPU())

	// 不同的并行度测试值
	// 这里刻意包含了：
	// - 很小的值（1, 2）
	// - 接近或超过常见 CPU 核心数的值
	testCases := []int{1, 2, 4, 8, 16, 32, 64, 100}

	for _, gmp := range testCases {
		// 设置当前运行时允许的最大并行度（P 的数量）
		runtime.GOMAXPROCS(gmp)

		start := time.Now()
		var wg sync.WaitGroup

		// 启动与 GOMAXPROCS 相同数量的 goroutine
		// 每个 goroutine 执行同样的计算任务
		for i := 0; i < gmp; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// 一个“足够重”的计算，
				// 让调度成本在总耗时中显现出来
				fibonacci(35)
			}()
		}

		// 等待所有 goroutine 完成
		wg.Wait()

		elapsed := time.Since(start)

		fmt.Printf(
			"GOMAXPROCS=%3d | goroutine=%3d | 总耗时=%v\n",
			gmp,
			gmp,
			elapsed,
		)
	}
}
