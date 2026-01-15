package main

import (
	"fmt"
	"sync"
)

type Config struct {
	Name string
}

func DemoCopyVsLock() {
	fmt.Println("----------")
	fmt.Println("复制数据 vs 加锁：是在维护共享，还是切断关系")

	// 场景一：通过加锁来维持共享
	// config 是一份被多个 goroutine 共享的数据
	// 为了保证并发安全，所有访问都必须受锁保护
	var (
		mu     sync.RWMutex
		config = Config{Name: "v1"}
	)

	// 读操作需要显式加锁
	// 这里返回的是同一份共享数据的“当前视图”
	readWithLock := func() Config {
		mu.RLock()
		defer mu.RUnlock()
		return config
	}

	// 场景二：通过复制来切断共享
	// 每次更新都生成一份新的配置快照
	// 调用方拿到的数据，从此与原状态无关
	updateAndCopy := func() Config {
		return Config{Name: "v2"} // 新快照
	}

	// 对比两种读取方式的语义差异
	fmt.Println("[加锁读取] 当前共享配置 Name =", readWithLock().Name)
	fmt.Println("[复制读取] 独立配置快照 Name =", updateAndCopy().Name)
}
