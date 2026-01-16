package main

import "fmt"

func DemoStructureDesign() {
	fmt.Println("----------")
	fmt.Println("如何通过结构设计避免并发 bug")
	// 对外暴露的只是“请求”，而不是 map 本身
	type request struct {
		key int
	}
	// 用 channel 作为唯一的交互入口
	reqCh := make(chan request)
	done := make(chan struct{})
	// map 的“唯一拥有者”
	go func() {
		defer close(done)

		m := make(map[int]int)
		fmt.Println("map 由单一 goroutine 持有，开始处理请求")

		for req := range reqCh {
			m[req.key] = req.key
			fmt.Println("处理请求：写入 key =", req.key)
		}

		fmt.Println("请求通道已关闭，map 最终状态为：", m)
	}()
	// 模拟多个请求发送方
	for i := 0; i < 5; i++ {
		fmt.Println("发送请求：key =", i)
		reqCh <- request{key: i}
	}
	// 所有请求发送完毕，关闭通道
	close(reqCh)
	// 等待 map 所属 goroutine 正常结束
	<-done
	fmt.Println("所有请求处理完成，程序安全退出")
}
