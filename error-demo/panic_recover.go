package main

import "fmt"

func demoPanicRecover() {
	fmt.Println("----------")
	fmt.Println("panic / recover 的使用边界")

	// 在函数入口处放置一层 defer + recover
	// 作用不是“忽略 panic”，而是防止程序直接崩溃
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("捕获到 panic，当前层选择记录并继续运行:", r)
		}
	}()

	// 启动应用的核心流程
	startApp()

	// 如果 panic 被 recover，调用 demoPanicRecover 方法的代码后续的方法依然可以执行
	// 代码会在 defer 执行完成后推出
}

func startApp() {
	cfg := loadCriticalConfig()

	// 这里判断的不是“操作是否失败”
	// 而是“程序是否还处在一个合理的状态”
	if cfg == "" {
		// 当关键假设被打破时，直接触发 panic
		// 表达的是：程序不应该在这种状态下继续运行
		panic("关键配置缺失，程序无法启动")
	}
}

func loadCriticalConfig() string {
	// 模拟初始化阶段的关键配置缺失
	// 这类问题通常无法通过普通 error 恢复
	return ""
}
