package main

import "fmt"

// DemoPassBetweenFuncs 演示 map 在函数间传递时的行为
func DemoPassBetweenFuncs() {
	fmt.Println("----------")
	fmt.Println("map 在函数间传递时的行为")

	// 初始化一个 map，包含一个初始值
	m := map[string]int{"x": 10}

	// 将 map 传入函数
	resetMap(m)

	// 观察函数调用结束后，原 map 的内容
	fmt.Println("函数调用结束后 m 的内容 ->", m)
}

// resetMap 接收一个 map 作为参数
func resetMap(m map[string]int) {
	// 在函数内部重新 make 一个新的 map
	// 这里只是修改了参数 m 自己的指向
	m = make(map[string]int)

	// 对新 map 的写入，只影响函数内部的这个新 map
	m["a"] = 1
	m["x"] = 2
}
