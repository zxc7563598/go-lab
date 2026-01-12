package main

import "fmt"

// DemoNilVsMake 演示 nil map 与通过 make 创建的 map 在行为上的差异
func DemoNilVsMake() {
	fmt.Println("----------")
	fmt.Println("nil map 与 make(map) 的差异")

	// 只声明，不初始化
	// m1 的零值是 nil，此时并没有真正的底层数据结构
	var m1 map[string]int

	// 使用 make 初始化
	// m2 指向一块已经分配好的 map 结构，可以安全读写
	m2 := make(map[string]int)

	// 判断是否为 nil
	fmt.Println("m1 == nil ->", m1 == nil)
	fmt.Println("m2 == nil ->", m2 == nil)

	// 从 nil map 中读取是安全的
	// 如果 key 不存在，返回对应类型的零值
	fmt.Println("从 nil map 中读取 m1[\"x\"] ->", m1["x"])

	// 向 make 创建的 map 中写入是正常的
	m2["x"] = 1
	fmt.Println("向已初始化的 map 写入后 m2 ->", m2)

	// 向 nil map 写入会直接触发 panic
	// 因为此时并没有可供写入的底层结构
	// m1["x"] = 1 // 会 panic
}
