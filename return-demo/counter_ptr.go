package main

// IncByPointer 是 Counter 的指针接收者方法。
// 使用指针接收者意味着方法内部操作的是原始数据，而不是副本。
// 调用此方法会直接修改 Counter 的 n 字段。
func (c *Counter) IncByPointer() {
	c.n++
}

// NewCounterPtr 返回一个 *Counter 指针。
// 返回指针类型的好处是，调用方可以直接修改原始 Counter 数据，
// 并且避免复制较大的 struct。
func NewCounterPtr() *Counter {
	return &Counter{n: 1}
}
