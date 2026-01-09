package main

// Counterer 是一个接口类型，定义了两个方法。
// 接口可以存放任何实现了这两个方法的类型（值类型或指针类型）。
type Counterer interface {
	Inc()       // 增加计数
	Value() int // 返回当前计数值
}

// Inc 和 Value 方法由 *Counter 实现。
// 使用指针接收者意味着方法内部会操作原始数据，而不是副本。
func (c *Counter) Inc() {
	c.n++
}

func (c *Counter) Value() int {
	return c.n
}

// NewCounterInterface 返回一个 Counterer 接口。
// 这里返回的是 *Counter 指针，实现了接口。
// 通过接口调用方法时，内部存储的是指针，可以修改原始 Counter。
func NewCounterInterface() Counterer {
	return &Counter{n: 1}
}
