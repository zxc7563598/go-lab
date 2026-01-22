package main

import (
	"fmt"
	"time"
)

// doWork 模拟一个需要耗时处理的任务
// 注意：这个函数本身并不知道自己是被如何调度的
// 它只负责“开始执行 -> 阻塞一段时间 -> 结束执行”
func doWork(id int) {
	// 打印任务开始执行的时间点，方便观察并发启动情况
	fmt.Printf("[work %d] 开始执行\n", id)

	// 模拟真实业务中的耗时操作（例如 IO、计算、外部请求等）
	time.Sleep(500 * time.Millisecond)

	// 打印任务结束执行
	fmt.Printf("[work %d] 执行结束\n", id)
}

// UnboundedGoroutinesDemo 演示「不受控的 goroutine 并发」
//
// 核心关注点不在于 doWork 做了什么，
// 而在于：每一次循环都会直接创建一个新的 goroutine。
//
// 当任务数量增长时，并发数量会线性增长，
// 这个结构本身没有任何“限流”或“背压”机制。
func UnboundedGoroutinesDemo() {
	// 连续启动多个 goroutine
	// 这里的 10 只是示例，如果是 1000、10000，结构上并没有区别
	for i := 0; i < 10; i++ {
		go doWork(i)
	}

	// 主 goroutine 主动 sleep
	// 目的只是为了防止 main 函数过早退出
	// 并不是一种可控或推荐的并发等待方式
	time.Sleep(3 * time.Second)
}
