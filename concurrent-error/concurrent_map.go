package main

import (
	"fmt"
	"sync"
)

// 演示最典型、也最容易遇到的情况：
// - 多个 goroutine 并发写同一个 map
// - map 不是并发安全的
// - 运行时直接 panic，而不是数据错乱
// 这是 Go 在这里选择的策略：与其让你得到“看起来还能用但已经坏掉的数据”，不如直接让程序崩溃。
func panicByConcurrentWrite() {
	fmt.Println("示例一: 多个 goroutine 并发写 map，直接 panic")

	m := make(map[int]int)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m[i] = i
		}(i)
	}

	wg.Wait()
}

// 演示一种更隐蔽的情况：
// - 一部分 goroutine 只读 map
// - 另一部分 goroutine 在写 map
// - 即使“读看起来是安全的”
// - 只要存在并发写，整体就是不安全的
// 读写混合并不会降低风险，反而更容易让人放松警惕。
func panicByConcurrentReadWrite() {
	fmt.Println("示例二: 并发读 + 并发写，同样会触发 panic")

	m := make(map[int]int)
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = m[i]
		}(i)
	}

	for i := 5; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m[i] = i
		}(i)
	}

	wg.Wait()
}

// 演示最直接、也最常见的修复方式：
// - 使用 Mutex 保护 map
// - 所有读写都必须经过同一把锁
// - map 本身不变，结构发生了变化
// 这里 map 依然不是“并发安全的”，只是我们保证了“不会并发访问它”。
func safeMapWithMutex() {
	fmt.Println("示例三: 使用 Mutex 串行化 map 的访问")

	m := make(map[int]int)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mu.Lock()
			m[i] = i
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	fmt.Println("最终 map 内容：", m)
}

// 演示一种“结构性规避并发”的写法：
// - map 只存在于一个 goroutine 中
// - 其他 goroutine 通过 channel 发送数据
// - 所有写操作在同一个 goroutine 内完成
// 这里甚至不需要 Mutex，因为从结构上就避免了并发访问。
func safeMapBySingleOwner() {
	fmt.Println("示例四: 通过 channel 让 map 只被一个 goroutine 持有")

	ch := make(chan int)
	m := make(map[int]int)

	go func() {
		for v := range ch {
			m[v] = v
		}
		fmt.Println("map 写入完成：", m)
	}()

	for i := 0; i < 10; i++ {
		ch <- i
	}

	close(ch)
}

func DemoConcurrentMap() {
	fmt.Println("----------")
	fmt.Println("并发 map 写导致的 panic")

	panicByConcurrentWrite()
	panicByConcurrentReadWrite()
	safeMapWithMutex()
	safeMapBySingleOwner()
}
