package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	// 演示：请求内部的并发模型
	http.HandleFunc("/multi", handlerMulti)

	fmt.Println("Web 中的并发模型与 goroutine 数量控制")
	fmt.Println("开始监听端口 :8080")
	http.ListenAndServe(":8080", nil)
}

// handlerMulti 演示：
// 并发发生在请求生命周期之内，而不是之外
func handlerMulti(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[/multi] handler 开始执行")

	// 每一个 HTTP 请求，net/http 都已经为我们提供了并发环境
	ctx := r.Context()

	var wg sync.WaitGroup
	wg.Add(2)

	var a, b string

	// 并发调用 A 服务
	go func() {
		fmt.Println("[/multi] goroutine A 启动")
		defer wg.Done()

		a = callA(ctx)

		fmt.Println("[/multi] goroutine A 结束")
	}()

	// 并发调用 B 服务
	go func() {
		fmt.Println("[/multi] goroutine B 启动")
		defer wg.Done()

		b = callB(ctx)

		fmt.Println("[/multi] goroutine B 结束")
	}()

	// handler 不返回，就意味着这次 HTTP 请求仍然存在
	wg.Wait()

	w.Write([]byte(a + " & " + b))
	fmt.Println("[/multi] handler 执行结束")
}

// callA 模拟一个下游调用
func callA(ctx context.Context) string {
	fmt.Println("[callA] 开始执行")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("[callA] 执行完成")
	return "A"
}

// callB 模拟另一个下游调用
func callB(ctx context.Context) string {
	fmt.Println("[callB] 开始执行")
	time.Sleep(700 * time.Millisecond)
	fmt.Println("[callB] 执行完成")
	return "B"
}
