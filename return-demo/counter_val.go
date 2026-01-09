package main

// IncByValue 是 Counter 的值接收者方法。
// 使用值接收者意味着方法内部操作的是副本，不会修改原始 Counter。
// 调用此方法只会改变副本的 n 字段，外部的 Counter 不受影响。
func (c Counter) IncByValue() {
	c.n++
}

// NewCounterVal 返回一个 Counter 值类型。
// 返回值类型会生成副本，调用方拿到的是独立的数据。
// 修改返回的 Counter 不会影响其他 Counter 实例。
func NewCounterVal() Counter {
	return Counter{n: 1}
}
