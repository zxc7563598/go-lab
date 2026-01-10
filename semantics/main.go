package main

import "fmt"

// Box 是一个普通的 struct，没有任何引用语义。
// 是否产生共享，完全取决于它是以“值”还是“指针”的形式被传递。
type Box struct {
	n int
}

// takeValue 接收的是 Box 的值。
// 调用时会发生一次值拷贝，函数内部操作的是这份拷贝。
func takeValue(b Box) {
	b.n++ // 只修改函数内部的副本
	fmt.Println("takeValue 内部:", b.n)
}

// takePointer 接收的是 *Box，即指向 Box 的位置。
// 函数通过这个位置，直接修改调用方持有的那份数据。
func takePointer(b *Box) {
	b.n++ // 等价于 (*b).n++
	fmt.Println("takePointer 内部:", b.n)
}

// returnValue 返回的是一个 Box 值。
// 调用方拿到的是一份新的值，不与函数内部的数据共享。
func returnValue() Box {
	b := Box{n: 1}
	return b // 返回值本身，而不是位置
}

// returnPointer 返回的是 *Box，即某个 Box 所在的位置。
// 虽然 b 是函数内的局部变量，但 Go 会保证这个地址是安全的。
func returnPointer() *Box {
	b := Box{n: 1}
	return &b // 返回的是位置，而不是值
}

// returnDereference 接收的是 *Box（位置），返回的是 Box（值）。
// 在函数内部先通过指针修改数据，再把当前值拷贝一份返回。
func returnDereference(b *Box) Box {
	b.n++     // 修改的是指针指向的那份数据
	return *b // 对当前位置解引用，返回一份值拷贝
}

// returnAddressOfValue 接收的是 Box 的值，返回的是 *Box。
// 返回的指针指向的是函数内部那份拷贝，而不是调用方的原始值。
func returnAddressOfValue(b Box) *Box {
	b.n++     // 只影响函数内部的那份值
	return &b // 返回的是“新值”的位置
}

func main() {
	fmt.Println("----------")
	fmt.Println("接收 x（值）")
	box1 := Box{n: 1}
	takeValue(box1)
	// box1 没有发生变化，因为传入的是值的拷贝
	fmt.Println("main 中:", box1.n)

	fmt.Println("----------")
	fmt.Println("接收 &x（位置）")
	box2 := Box{n: 1}
	takePointer(&box2)
	// box2 被修改，因为函数操作的是它的实际位置
	fmt.Println("main 中:", box2.n)

	fmt.Println("----------")
	fmt.Println("返回 x（值）")
	box3 := returnValue()
	box3.n++
	// 对 box3 的修改只作用于 main 中的这份值
	fmt.Println("main 中:", box3.n)

	fmt.Println("----------")
	fmt.Println("返回 &x（位置）")
	box4 := returnPointer()
	box4.n++
	// 通过返回的指针，直接修改指针指向的数据
	fmt.Println("main 中:", box4.n)

	fmt.Println("----------")
	fmt.Println("接收 &x，返回 *x")
	box5 := Box{n: 1}
	box6 := returnDereference(&box5)
	box6.n++
	// box5 在函数中被修改过一次
	// box6 是从 box5 解引用后得到的一份新值
	fmt.Println("原 box:", box5.n)
	fmt.Println("新 box:", box6.n)

	fmt.Println("----------")
	fmt.Println("接收 x，返回 &x")
	box7 := Box{n: 1}
	box8 := returnAddressOfValue(box7)
	box8.n++
	// box7 始终未变
	// box8 指向的是函数内部创建的那份新值
	fmt.Println("原 box:", box7.n)
	fmt.Println("指针 box:", box8.n)
}
