package main

import "fmt"

func testPointer() {
	fmt.Println("----------")
	fmt.Println("【指针语义】只有通过指针，才能共享并修改同一份数据")
	// 定义一个 Counter 结构体变量
	c := Counter{n: 10}
	// 将 c 的地址传入函数
	// 这里传递的不是值的拷贝，而是指向同一块内存的指针
	incPtr(&c)
	// 由于函数内部通过指针修改了数据
	// 原始变量 c 的值发生了变化
	fmt.Println("通过指针调用函数后，c.n 的值是：", c.n) // 11
}

func incPtr(c *Counter) {
	// c 是一个指向 Counter 的指针
	// 通过指针修改字段，等价于修改原始变量
	c.n++
}
