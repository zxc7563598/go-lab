package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// gmpWork 用来模拟在不同 GOMAXPROCS 设置下的执行情况
//
// 这里关注的不是 goroutine 的数量（始终是 100 个），
// 而是：
// - 在不同 P 数量下
// - 这些 goroutine 是如何被调度执行的
func gmpWork(maxProcs int) {
	// 设置当前运行时允许的最大并行度
	// 本质上就是：P 的数量
	runtime.GOMAXPROCS(maxProcs)

	start := time.Now()

	var wg sync.WaitGroup

	// 启动固定数量的 goroutine
	// 无论 GOMAXPROCS 是多少，goroutine 的数量都保持不变
	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// 一个纯 CPU 计算任务
			// 不涉及 I/O、不涉及锁，
			// 目的是把注意力集中在“调度与并行度”本身
			sum := 0
			for i := 0; i < 10_000_000; i++ {
				sum += i
			}
			_ = sum
		}()
	}

	// 等待所有 goroutine 执行完成
	wg.Wait()

	elapsed := time.Since(start)

	fmt.Printf(
		"GOMAXPROCS=%2d | 总耗时=%v\n",
		maxProcs,
		elapsed,
	)
}

func demoGMP() {
	fmt.Println("----------")
	fmt.Println("GMP 调度模型的直觉理解")
	fmt.Println("关注点：P 才是并行执行 Go 代码的上限")
	fmt.Println()

	// 在不同的 P 数量下，执行同样的任务
	// 对比总耗时的变化趋势
	gmpWork(1)
	gmpWork(5)
	gmpWork(10)
	gmpWork(15)
	gmpWork(20)
	gmpWork(30)
	gmpWork(50)
}
