package main

import (
	"fmt"
	"sync"
)

// DemoDefensiveUsage 演示一种最基础、也是最常见的 map 防御性使用方式：
// 在 map 的所有读写操作外层，统一加一把互斥锁。
func DemoDefensiveUsage() {
	fmt.Println("----------")
	fmt.Println("通过互斥锁保护 map 的读写")

	// 使用 make 初始化 map
	// 这里只关心并发安全问题，不讨论 nil map 的情况
	m := make(map[string]int)

	// 一把互斥锁，用来保护对 map 的所有访问
	var mu sync.Mutex

	// 封装写操作：
	// 每次写入前加锁，写完后释放
	// 确保任意时刻，只有一个 goroutine 在修改 map
	write := func(k string, v int) {
		mu.Lock()
		defer mu.Unlock()

		m[k] = v
	}

	// 封装读操作：
	// 即便只是读取，也同样需要加锁
	// 否则在并发场景下，仍然可能触发 runtime panic
	read := func(k string) int {
		mu.Lock()
		defer mu.Unlock()

		return m[k]
	}

	// 在单线程下调用，看起来和普通 map 用法没有区别
	// 但这些封装使得它在并发场景下依然是安全的
	write("a", 1)
	write("b", 2)

	fmt.Println("read a ->", read("a"))
	fmt.Println("read b ->", read("b"))
}
