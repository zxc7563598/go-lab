package main

import "fmt"

func demoIfErrChoice() {
	fmt.Println("----------")
	fmt.Println("为什么 Go 选择显式地写 if err != nil")

	// 调用一个可能失败的函数
	data, err := readConfig()

	// 错误不会自动中断执行
	// 是否检查、如何处理，完全由调用方决定
	if err != nil {
		fmt.Println("函数返回了 error，当前这一层选择在这里处理：", err)
		fmt.Println("如果这里不判断 err，后面的代码仍然会继续执行")
		// return
	}

	// 只有在 err == nil 时，返回的数据才有意义
	fmt.Println("配置加载成功，返回的数据是:", data)
}

func readConfig() (string, error) {
	// 在 Go 中，函数只负责“如实返回结果”
	// 是否严重到需要终止流程，并不在这一层做判断
	return "", fmt.Errorf("配置文件未找到")
}
