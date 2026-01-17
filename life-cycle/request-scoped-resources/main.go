package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// 正确示例：request 级资源（context + timeout）与请求生命周期绑定
	http.HandleFunc("/timeout", handler)

	// 危险示例：request 级资源被 goroutine 带出了 handler
	http.HandleFunc("/timeout-error", handlerError)

	fmt.Println("request 级资源的创建与释放")
	fmt.Println("开始监听端口 :8080")
	http.ListenAndServe(":8080", nil)
}

// handler 演示：
// request 级资源的创建与主动释放
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[/timeout] handler 开始执行")

	// 基于 HTTP 请求的 context，创建一个带超时的子 context
	// 这个 context 只应该活在这一次请求里
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer func() {
		// 主动释放 request 级资源
		fmt.Println("[/timeout] 调用 cancel，释放 request 级资源")
		cancel()
	}()

	result := callService(ctx)

	fmt.Println("[/timeout] service 返回结果：", result)

	w.Write([]byte(result))
	fmt.Println("[/timeout] handler 执行结束")
}

// handlerError 演示一种“语法正确，但生命周期危险”的写法
func handlerError(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[/timeout-error] handler 开始执行")

	ctx := r.Context()

	go func() {
		fmt.Println("[/timeout-error] goroutine 启动（已脱离 handler）")

		// 这里仍然在使用 request 级 context
		// 但此时 handler 很可能已经返回
		result := callService(ctx)

		fmt.Println("[/timeout-error] goroutine 中的 service 返回结果：", result)
	}()

	// handler 很快结束，HTTP 请求生命周期到此为止
	w.Write([]byte("ok"))
	fmt.Println("[/timeout-error] handler 执行结束")
}

// callService 模拟一个耗时操作：
// 它完全依赖 ctx 来判断自己是否应该继续执行
func callService(ctx context.Context) string {
	fmt.Println("[service] 开始执行业务逻辑")

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("[service] 业务逻辑正常完成")
		return "ok"

	case <-ctx.Done():
		fmt.Println("[service] context 被取消，提前结束")
		return "timeout"
	}
}
