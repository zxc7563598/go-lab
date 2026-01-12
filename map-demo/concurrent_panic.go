package main

import (
	"fmt"
	"time"
)

// DemoConcurrentPanic 用来演示：
// map 在没有任何同步手段的情况下，
// 一旦发生并发读 + 并发写，运行时会直接 panic。
func DemoConcurrentPanic() {
	fmt.Println("----------")
	fmt.Println("map 在并发读写场景下的行为")

	// 使用 make 初始化一个可写的 map
	// 注意：这里 map 本身是引用类型，但并不意味着它是并发安全的
	m := make(map[string]int)

	// 启动一个 goroutine，不断对 map 进行写操作
	// 这在单线程下是完全合法的
	go func() {
		for i := 0; i < 10000; i++ {
			// 对同一个 key 反复写入
			m["a"] = i
		}
	}()

	// 启动另一个 goroutine，同时对 map 进行读操作
	go func() {
		for i := 0; i < 10000; i++ {
			// 读取同一个 key
			// 即便只是读，也会参与并发访问
			_ = m["a"]
		}
	}()

	// 这里用 Sleep 只是为了：
	// 1. 防止 main goroutine 过早退出
	// 2. 给上面的两个 goroutine 足够的时间并发执行
	//
	// 在大多数运行环境中，这段代码都会触发（并非一定触发，所以说是坑）：
	// fatal error: concurrent map read and map write
	time.Sleep(time.Second)
}
