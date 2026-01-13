package main

import "fmt"

func demoDeferLoopWrong() {
	fmt.Println("----------")
	fmt.Println("defer + loop（闭包捕获变量，直觉容易出错）")

	// 显式声明 i，强调：整个循环过程中只有这一个变量实例
	var i int

	for i = 0; i < 3; i++ {
		// defer 的是一个“无参函数”
		// 这个函数内部直接使用了外部变量 i
		//
		// 注意：这里 defer 并不会在此刻打印 i，
		// 只是把“这个函数”压入 defer 栈
		defer func() {
			// 当函数真正执行时，循环早已结束
			// 此时 i 的值已经变成了 3
			fmt.Println("defer 执行时看到的 i =", i)
		}()
	}

	// 循环结束后，i == 3
	// 接下来函数返回，defer 栈开始出栈执行
}

func demoDeferLoopCorrect() {
	fmt.Println("----------")
	fmt.Println("defer + loop（通过参数显式绑定当次循环的值）")

	// 同样只有一个 i 变量，但我们不再直接让闭包引用它
	var i int

	for i = 0; i < 3; i++ {
		// 这里 defer 的是“一次已经绑定好参数的函数调用”
		//
		// 当前轮的 i 会被拷贝一份，作为参数 n 传入
		// defer 入栈时，n 的值就已经确定
		defer func(n int) {
			fmt.Println("defer 执行时绑定的 n =", n)
		}(i)
	}

	// 函数返回时，defer 栈按后进先出顺序执行
	// 输出结果为：
	// 2
	// 1
	// 0
}
