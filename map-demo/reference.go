package main

import "fmt"

// DemoReferenceButNotSafe 演示 map 在函数间传递时的“引用式行为”
func DemoReferenceButNotSafe() {
	fmt.Println("----------")
	fmt.Println("map 在函数间传递时会共享同一份数据")

	// 使用 make 创建一个可写的 map
	m := make(map[string]int)

	// 将 map 传入函数，在函数内部进行修改
	modify(m)

	// 函数返回后，观察 m 的内容
	// 可以看到，函数内部的修改对这里是可见的
	fmt.Println("函数调用结束后 m 的内容 ->", m)
}

// modify 接收一个 map，并直接对其进行写入
func modify(m map[string]int) {
	// 这里并没有重新分配 map
	// 而是通过同一个底层结构进行写操作
	m["a"] = 1
	m["b"] = 2
}
