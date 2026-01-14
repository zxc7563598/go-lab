package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func demoNotThread() {
	fmt.Println("----------")
	fmt.Println("goroutine 不是线程")

	// 将 GOMAXPROCS 设置为 1
	//
	// 含义不是“只能起一个 goroutine”，
	// 而是：同一时刻，最多只有一个 P 在执行 Go 代码
	// 也就是说，这里刻意把“并行”限制为 1
	runtime.GOMAXPROCS(1)

	start := time.Now()

	log := func(format string, args ...any) {
		elapsed := time.Since(start).Truncate(time.Millisecond)
		prefix := fmt.Sprintf("[%6s] ", elapsed)
		fmt.Printf(prefix+format+"\n", args...)
	}

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			// 这里的输出并不表示“真的同时运行”
			// 在 GOMAXPROCS=1 的情况下，
			// 所有 goroutine 实际上是在一个执行位上被轮流调度
			log("goroutine %d 开始执行", id)

			// Sleep 会让当前 goroutine 主动让出执行权
			// runtime 可以在这段时间内切换去执行其他 goroutine
			time.Sleep(500 * time.Millisecond)

			log("goroutine %d 执行结束", id)
		}(i)
	}

	// 等待所有 goroutine 执行完成
	wg.Wait()
}
