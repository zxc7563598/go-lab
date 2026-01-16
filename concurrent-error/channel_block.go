package main

import "fmt"

// 演示最直接的一种阻塞：
// - 向无缓冲 channel 发送数据
// - 但此时根本没有任何接收方
// - 发送操作会一直阻塞
func blockBySendWithoutReceiver() {
	fmt.Println("示例一: 无接收方的发送会永久阻塞")

	ch := make(chan int)

	// 这里没有任何 goroutine 在接收
	ch <- 1

	// 这一行永远不会被执行
	fmt.Println("这一行代码永远到不了")
}

// 演示一种“结构上断裂”的情况：
// - 启动了一个接收 goroutine
// - 但主 goroutine 没有任何等待逻辑
// - main 直接返回，程序整体退出
// 这里不是 channel 卡住，
// 而是你以为会发生的通信，其实根本没来得及发生。
func blockByGoroutineExitBeforeReceive() {
	fmt.Println("示例二: 接收 goroutine 还没等到数据，程序就结束了")

	ch := make(chan int)

	go func() {
		v := <-ch
		fmt.Println("子 goroutine 收到数据：", v)
	}()

	fmt.Println("main goroutine 结束，程序直接退出")
}

// 演示：
// - goroutine 使用 range 读取 channel
// - channel 发送完数据后没有 close
// - 接收方永远卡在 range 上
// - 主 goroutine 用 select{} 强行挂住程序
func blockByRangeWithoutClose() {
	fmt.Println("示例三: range channel，但永远等不到 close")

	ch := make(chan int)

	go func() {
		for v := range ch {
			fmt.Println("收到数据：", v)
		}
		fmt.Println("这一行只有在 channel 关闭时才会打印")
	}()

	ch <- 1
	ch <- 2

	fmt.Println("数据已发送完，但 channel 没有关闭")

	// 主 goroutine 永久阻塞，方便观察子 goroutine 的状态
	select {}
}

// 演示：
// - select 中只有一个发送分支
// - 但此时没有接收方
// - 没有 default
// - select 会整体阻塞
func blockBySelectWithoutReadyCase() {
	fmt.Println("示例四: select 中没有任何可执行的分支")

	ch := make(chan int)

	select {
	case ch <- 1:
		fmt.Println("发送成功（但实际上永远不会发生）")
	}
}

// 演示：
// - 有缓冲 channel 并不等于“永远不会阻塞”
// - 当缓冲区满了之后
// - 再次发送仍然会阻塞
func blockByExceedBufferCapacity() {
	fmt.Println("示例五: 超过缓冲区容量的发送会阻塞")

	ch := make(chan int, 1)

	ch <- 1
	fmt.Println("第一次发送成功（缓冲区容量已满）")

	// 第二次发送没有接收方，会阻塞
	ch <- 2

	fmt.Println("这一行同样不会被执行")
}

func DemoChannelBlock() {
	fmt.Println("----------")
	fmt.Println("channel 永远阻塞的原因")

	blockBySendWithoutReceiver()
	blockByGoroutineExitBeforeReceive()
	blockByRangeWithoutClose()
	blockBySelectWithoutReadyCase()
	blockByExceedBufferCapacity()
}
