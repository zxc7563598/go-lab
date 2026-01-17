package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// 注册三个 handler，用来演示不同的请求生命周期形态
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/work", workHandler)
	http.HandleFunc("/async", asyncHandler)

	fmt.Println("HTTP 请求在 Go 中的完整生命周期")
	fmt.Println("开始监听端口 :8080")
	// ListenAndServe 会阻塞在这里
	// 每一个进入的 HTTP 请求，都会由 net/http
	// 分配一个独立的 goroutine 来执行对应的 handler
	http.ListenAndServe(":8080", nil)
}

// helloHandler 用来演示：
// 一个最“干净”的请求生命周期
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[/hello] handler 开始执行")

	// 在 handler 被调用的这一刻，请求正式进入我们的代码
	w.Write([]byte("hello"))

	// handler 返回，意味着这次 HTTP 请求在我们这里结束
	fmt.Println("[/hello] handler 执行结束")
}

// workHandler 用来演示：
// 请求生命周期是如何通过 context 向下传递的
func workHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[/work] handler 开始执行")

	// r.Context() 代表“这一次 HTTP 请求的生命周期”
	ctx := r.Context()

	// 将请求的 context 继续向下传递
	result := doWork(ctx)

	fmt.Println("[/work] doWork 返回结果：", result)

	w.Write([]byte(result))
	fmt.Println("[/work] handler 执行结束")
}

func doWork(ctx context.Context) string {
	fmt.Println("[/work] 开始执行业务逻辑，等待 5 秒")

	select {
	case <-time.After(5 * time.Second):
		// 这里代表：业务逻辑本身完成
		fmt.Println("[/work] 业务逻辑正常完成")
		return "业务完成"

	case <-ctx.Done():
		// 这里代表：HTTP 请求已经结束（客户端断开、超时等）
		fmt.Println("[/work] context 被取消，请求已结束")
		return "请求取消"
	}
}

// asyncHandler 用来演示：
// goroutine 脱离 HTTP 请求生命周期后的状态
func asyncHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[/async] handler 开始执行")

	go func() {
		fmt.Println("[/async] 异步 goroutine 启动")
		time.Sleep(2 * time.Second)
		// 这里执行时，请求很可能已经结束
		fmt.Println("[/async] 异步执行完成（已脱离请求生命周期）")
	}()

	// handler 很快返回，HTTP 请求在这里就结束了
	w.Write([]byte("请求结束"))
	fmt.Println("[/async] handler 执行结束")
}
