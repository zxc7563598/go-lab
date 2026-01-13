package main

import (
	"errors"
	"fmt"
)

func demoErrorIsValue() {
	fmt.Println("----------")
	fmt.Println("error 在 Go 中只是一个普通返回值")

	// 调用一个可能失败的函数
	// 注意：函数返回后，程序并不会自动中断
	v, err := divide(10, 0)

	// err 是否为 nil，完全需要调用方自己判断
	if err != nil {
		fmt.Println("函数正常返回，但同时返回了一个 error:", err)
		fmt.Println("这里如果不主动判断 err，程序依然会继续向下执行")
		// return
	}

	// 只有在 err == nil 时，返回值 v 才是可信的
	fmt.Println("计算结果:", v)
}

func divide(a, b int) (int, error) {
	// 在 Go 中，错误通常被当作一种“失败结果”返回
	// 而不是通过异常机制中断控制流
	if b == 0 {
		// 返回值 + error
		// 这里只是构造并返回一个 error，不会自动影响调用方的执行流程
		return 0, errors.New("division by zero")
	}

	// 成功时，error 返回 nil
	return a / b, nil
}
