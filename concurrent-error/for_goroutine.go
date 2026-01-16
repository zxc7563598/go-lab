package main

import (
	"fmt"
	"os"
	"sync"
)

// 演示一种非常常见、但不容易被意识到的问题：
// - for 循环里启动了多个 goroutine
// - 但 main goroutine 没有任何等待逻辑
// - main 很快结束，程序整体退出
// - 很多 goroutine 甚至来不及执行
// 这里没有变量捕获问题，问题出在：主流程比并发任务更早结束。
func earlyExitBeforeGoroutineRun() {
	fmt.Println("示例一: goroutine 还没来得及执行，main 就退出了")

	for i := 0; i < 5; i++ {
		go func(i int) {
			fmt.Println("子 goroutine 打印：", i)
		}(i)
	}

	fmt.Println("main goroutine 结束，程序直接退出")
}

// 演示：
// - 多个 goroutine 并发修改同一个 slice
// - append 不是并发安全的
// - 即使使用了 WaitGroup 等待结束
// - 数据竞争仍然已经发生
// WaitGroup 只能保证“等到结束”，但不能保证“过程是安全的”。
func concurrentAppendWithoutProtection() {
	fmt.Println("示例二: 多个 goroutine 并发 append 同一个 slice")

	var result []int
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			result = append(result, i)
		}(i)
	}

	wg.Wait()
	fmt.Println("最终结果（顺序和内容都不可靠）：", result)
}

// 演示：
// - 多个 goroutine 共享同一个 *os.File
// - 并发调用 Read
// - 文件内部的偏移量是共享状态
// - 每个 goroutine 实际读到的内容不可预测
// 问题不在 for，也不在 goroutine，而在于：共享了一个“带内部状态的对象”
func concurrentReadFromSameFile() {
	fmt.Println("示例三: 多个 goroutine 并发读取同一个文件")

	file, _ := os.Open("concurrent-error.txt")
	defer file.Close()

	for i := 0; i < 3; i++ {
		go func(i int) {
			buf := make([]byte, 10)
			n, _ := file.Read(buf)
			fmt.Println("goroutine", i, "读取到：", string(buf[:n]))
		}(i)
	}

	// 没有任何等待，main 可能提前退出
	fmt.Println("main goroutine 继续执行，读取结果不可预测")
}

func DemoForGoroutine() {
	fmt.Println("----------")
	fmt.Println("for + goroutine 的经典陷阱")

	earlyExitBeforeGoroutineRun()
	concurrentAppendWithoutProtection()
	concurrentReadFromSameFile()
}
