package main

import (
	"fmt"
	"time"
)

func demoDeferInWebLikeContext() {
	fmt.Println("----------")
	fmt.Println("defer 在“请求生命周期”中的作用边界")

	// 可以把这个 resource 想象成：
	// - 一个数据库连接
	// - 一个 HTTP 请求上下文里的资源
	// - 或者一次请求中创建的临时对象
	resource := "connection"

	// defer 注册的清理逻辑，只会在「当前函数返回」时执行
	// 在 Web 场景中，通常等价于：请求处理函数结束
	defer fmt.Println("defer 执行：关闭资源 ->", resource)

	// 启动一个 goroutine，模拟异步逻辑
	// 注意：这里把 resource 作为参数传进去，
	// 是为了避免和 defer + 闭包变量混在一起讨论
	go func(res string) {
		// 模拟异步任务比当前函数“活得更久”
		time.Sleep(100 * time.Millisecond)

		// 这个输出发生在 defer 之后还是之前，
		// 取决于当前函数什么时候返回
		fmt.Println("goroutine 中尝试使用资源 ->", res)
	}(resource)

	// 模拟请求处理逻辑已经完成
	// 在真实 Web 框架中，这一行之后，
	// handler 函数很快就会返回
	fmt.Println("请求处理逻辑结束，handler 即将返回")
}
