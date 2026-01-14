package main

import (
	"fmt"
	"sync"
)

func demoAvoidSharedState() {
	fmt.Println("----------")
	fmt.Println("使用 channel 避免共享状态")

	// ch 用来传递数据，而不是共享状态
	// 所有发送者只负责把值送进 ch，不直接操作任何共享变量
	ch := make(chan int)

	var wg sync.WaitGroup

	// 启动多个发送者 goroutine
	// 这些 goroutine 的职责非常单一：只产生数据，不持有状态
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			// 将计算结果发送到 channel
			// 这里并不会修改任何共享变量
			ch <- val
		}(i)
	}

	// resultChan 用来返回最终结果
	// 它的存在是为了让 main goroutine
	// 明确地等待“结果计算完成”这一事件
	resultChan := make(chan int)

	// 启动一个专门的接收者 goroutine
	// 这个 goroutine 独占 total 这个状态
	go func() {
		total := 0

		// 通过 range 顺序接收 channel 中的数据
		// 只要 channel 没有被关闭，这个循环就会一直等待
		for v := range ch {
			total += v
		}

		// 当 channel 被关闭后，说明不会再有新数据
		// 此时 total 的值是确定的，可以安全地发送结果
		resultChan <- total
	}()

	// 启动一个“协调者”goroutine
	// 它并不参与发送具体数据，只负责：
	// 在所有发送者结束后，关闭 channel
	go func() {
		wg.Wait()
		// 关闭 channel 表示：
		// “不会再有新的数据发送进来了”
		close(ch)
	}()

	// 等待接收最终计算结果
	result := <-resultChan

	// 输出最终结果
	// 这里的重点不是“算出来多少”，
	// 而是：整个过程中没有任何共享状态被并发修改
	fmt.Println("最终汇总结果:", result)
}
