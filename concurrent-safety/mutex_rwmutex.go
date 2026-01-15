package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func DemoMutexAndRWMutex() {
	fmt.Println("----------")
	fmt.Println("mutex 与 RWMutex 的使用边界")

	demoWithMutex()
	demoWithRWMutexStable()
	demoWithRWMutexUnstable()
}

func demoWithMutex() {
	fmt.Println("普通 Mutex: 保守但稳定，不区分读写")

	var (
		mu    sync.Mutex
		value int
		wg    sync.WaitGroup
	)

	// 多个 goroutine 同时写同一个变量
	// 不管是读还是写，全部串行化
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 3; j++ {

				// 进入临界区
				mu.Lock()

				// 对共享数据的修改是确定的、可预期的
				value++
				fmt.Println("write by goroutine", id, "value =", value)

				// 离开临界区
				mu.Unlock()

				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}

	// 等待 goroutine 执行结束
	wg.Wait()
}

func demoWithRWMutexStable() {
	fmt.Println("RWMutex: 读多写少，且边界清晰（成立）")

	var (
		mu    sync.RWMutex
		value int
		wg    sync.WaitGroup
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 写操作集中、数量少
	// 明确知道：这是“会修改状态的路径”
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				mu.Lock()
				value++
				fmt.Println("write value =", value)
				mu.Unlock()
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// 读操作频繁，且保证是“纯读”
	// 多个 goroutine 可以同时持有读锁
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					mu.RLock()
					fmt.Println("read by goroutine", id, "value =", value)
					mu.RUnlock()

					time.Sleep(100 * time.Millisecond)
				}
			}
		}(i)
	}

	wg.Wait()
}

func demoWithRWMutexUnstable() {
	fmt.Println("RWMutex: 读写交错、模式不稳定（边界被打破）")

	var (
		mu    sync.RWMutex
		value int
		wg    sync.WaitGroup
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 同一个 goroutine 内部，读写频繁交替
	// 很难再用“读多写少”来描述这个模型
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var j = 0
			for {
				select {
				case <-ctx.Done():
					return
				default:
					if j%2 == 0 {
						// 表面上是读
						mu.RLock()
						fmt.Println("read by goroutine", id, "value =", value)
						mu.RUnlock()
					} else {
						// 紧接着就是写
						mu.Lock()
						value++
						fmt.Println("write by goroutine", id, "value =", value)
						mu.Unlock()
					}
					time.Sleep(100 * time.Millisecond)
					j++
				}
			}
		}(i)
	}

	wg.Wait()
}
