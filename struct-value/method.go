package main

import "fmt"

func testMethod() {
	fmt.Println("----------")
	fmt.Println("【方法接收者】方法并不等于对象行为")
	// 定义一个 Counter 结构体变量
	c := Counter{n: 10}
	// 调用值接收者方法
	// 实际上传入的是 c 的一份拷贝
	c.incByValue()
	// 原始变量 c 不会被修改
	fmt.Println("调用值接收者方法后，c.n 的值是：", c.n) // 10
	// 调用指针接收者方法
	// 这里传入的是 c 的地址
	c.incByPointer()
	// 指针接收者方法可以修改原始变量
	fmt.Println("调用指针接收者方法后，c.n 的值是：", c.n) // 11
}

// 值接收者方法
// 接收的是 Counter 的一个副本
func (c Counter) incByValue() {
	// 修改的只是副本中的字段
	c.n++
}

// 指针接收者方法
// 接收的是 Counter 的指针
func (c *Counter) incByPointer() {
	// 通过指针修改的是原始变量
	c.n++
}
