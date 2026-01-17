package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// 纯 handler 示例：所有事情都堆在 handler 里
	http.HandleFunc("/user", userHandler)

	// middleware + handler 示例：middleware 只包裹请求生命周期
	http.Handle("/hello", timingMiddleware(http.HandlerFunc(helloWordHandler)))

	// handler + service 示例：handler 负责 HTTP，service 负责业务
	http.HandleFunc("/greet", greetHandler)

	fmt.Println("handler、middleware、service 的职责边界")
	fmt.Println("开始监听端口 :8080")
	http.ListenAndServe(":8080", nil)
}

// userHandler 演示一种“最原始”的 handler 写法：
// 解析参数 + 拼响应，所有责任都在这里
func userHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[/user] handler 开始执行")

	id := r.URL.Query().Get("id")
	w.Write([]byte("user id: " + id))

	fmt.Println("[/user] handler 执行结束")
}

// timingMiddleware 演示 middleware 的职责：
// 不关心业务，只关心“一次请求从开始到结束发生了什么”
func timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[middleware] 请求进入")

		start := time.Now()

		// 调用下一个 handler
		next.ServeHTTP(w, r)

		// handler 返回，意味着请求在业务层面已经结束
		fmt.Println("[middleware] 请求结束，耗时:", time.Since(start))
	})
}

// helloWordHandler 是一个非常纯粹的 handler：
// 它不知道 middleware 的存在，也不关心请求耗时
func helloWordHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[/hello] handler 开始执行")

	time.Sleep(100 * time.Millisecond)
	w.Write([]byte("hello word"))

	fmt.Println("[/hello] handler 执行结束")
}

// greetHandler 演示 handler 与 service 的职责边界：
// handler 负责 HTTP 世界，service 负责业务世界
func greetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[/greet] handler 开始执行")

	// 在 handler 中完成 HTTP 相关的解析
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}

	// 将请求的 context 向业务层传递
	result := greetService(r.Context(), name)

	w.Write([]byte(result))
	fmt.Println("[/greet] handler 执行结束")
}

// greetService 不知道 HTTP，也不应该知道：
// 它只关心业务输入，以及是否仍然处在请求生命周期中
func greetService(ctx context.Context, name string) string {
	fmt.Println("[service] 进入业务逻辑")

	// 这里暂时不使用 ctx，只是明确：
	// 这个 service 运行在一次请求的生命周期之内
	_ = ctx

	return "hello " + name
}
