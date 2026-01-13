package main

import "fmt"

func demoExplicitResourceRelease() {
	fmt.Println("----------")
	fmt.Println("Go 中资源释放为什么必须显式")

	fmt.Println("进入外层函数 demoExplicitResourceRelease")

	// 这里使用一个匿名函数，刻意制造一个“更小的作用域”
	// 用来观察 defer 是如何严格绑定在函数边界上的
	func() {
		// 可以把 resource 想象成：
		// - 打开的文件
		// - 数据库连接
		// - 网络连接等需要手动释放的资源
		resource := "file"

		// defer 注册的释放逻辑，只属于“当前这个函数作用域”
		// 一旦这个函数返回，这里的 defer 就会立刻执行
		defer fmt.Println("defer 执行：关闭资源 ->", resource)

		// 模拟资源的使用过程
		fmt.Println("函数内部：正在使用资源 ->", resource)

		// 注意：这里并没有显式调用 Close，
		// 资源的释放完全依赖 defer + 函数返回
	}()

	// 当执行到这里时：
	// - 上面的匿名函数已经返回
	// - 该函数内注册的 defer 已经全部执行完毕
	fmt.Println("离开外层函数 demoExplicitResourceRelease")
}
