package main

import "fmt"

func testAssign() {
	fmt.Println("----------")
	fmt.Println("【赋值行为】struct 赋值会发生值拷贝")
	// a 是一个 Counter 结构体变量
	a := Counter{n: 10}
	// b := a 会把 a 的整个值复制一份给 b
	// 此时 a 和 b 是两个完全独立的变量
	b := a
	// 修改 b 的字段，只会影响 b 自己
	b.n = 20
	// a 的值不会发生任何变化
	fmt.Println("a.n 的值是：", a.n) // 10
	fmt.Println("b.n 的值是：", b.n) // 20
}

func testParam() {
	fmt.Println("----------")
	fmt.Println("【参数传递】struct 作为函数参数会再次发生值拷贝")
	// c 是一个 Counter 结构体变量
	c := Counter{n: 10}
	// 调用 inc 时，c 会被复制一份传入函数
	inc(c)
	// 原始变量 c 不会被修改
	fmt.Println("函数调用后，c.n 的值是：", c.n) // 10
}

func inc(c Counter) {
	// 这里的 c 只是外部变量的一个副本
	// 对 c 的修改不会影响调用方
	c.n++
}
