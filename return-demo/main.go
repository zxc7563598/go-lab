package main

import "fmt"

// struct（值）定义
type Counter struct {
	n int
}

func main() {
	fmt.Println("----------")
	fmt.Println("返回 struct（值）")
	c1 := NewCounterVal()
	fmt.Println("初始值:", c1.n)
	c1.IncByValue()
	fmt.Println("调用 IncByValue 后:", c1.n) // 值接收者，不会改变原值
	fmt.Println("----------")
	fmt.Println("返回 *struct（指针）")
	c2 := NewCounterPtr()
	fmt.Println("初始值:", c2.n)
	c2.IncByPointer()
	fmt.Println("调用 IncByPointer 后:", c2.n) // 修改生效
	fmt.Println("----------")
	fmt.Println("返回 interface")
	c3 := NewCounterInterface()
	fmt.Println("初始值:", c3.Value())
	c3.Inc()
	fmt.Println("调用 Inc() 后:", c3.Value())
	fmt.Println("----------")
	c4 := c3
	c4.Inc()
	fmt.Println("赋值给 c4 并 Inc() 后:")
	fmt.Println("c3.Value():", c3.Value())
	fmt.Println("c4.Value():", c4.Value())
}
