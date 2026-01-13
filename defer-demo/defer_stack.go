package main

import "fmt"

func demoDeferStack() {
	fmt.Println("----------")
	fmt.Println("defer 的执行时机与栈模型")

	// 第一次 defer：最早注册，最晚执行
	// 此时并不会立刻打印，而是把这次调用压入 defer 栈
	defer fmt.Println("defer 执行：第一个注册的 defer")

	// 第二次 defer：后注册
	// 在栈中位置比“第一个 defer”更靠上
	defer fmt.Println("defer 执行：第二个注册的 defer")

	// 第三次 defer：最后注册
	// 位于栈顶，会在函数返回时最先执行
	defer fmt.Println("defer 执行：第三个注册的 defer")

	// 普通语句会立即执行
	// defer 中注册的函数，直到当前函数返回时才会统一执行
	fmt.Println("函数主体执行结束，即将开始返回")
}
