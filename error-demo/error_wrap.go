package main

import (
	"errors"
	"fmt"
	"io/fs"
)

func demoErrorWrap() {
	fmt.Println("----------")
	fmt.Println("错误包装（wrap）与错误定位")

	// 从最外层调用开始，只拿到一个 error
	err := loadFile()
	if err != nil {
		fmt.Println("最终拿到的错误信息:", err)

		// 即使错误已经被多层包装，仍然可以判断它的根本类型
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("通过 errors.Is 判断：错误根源是 fs.ErrNotExist")
		} else {
			fmt.Println("错误来源无法定位到具体的标准库错误")
		}
	}
}

func loadFile() error {
	// 这一层并不知道如何修复错误
	// 它所做的只是补充“当前正在做什么”的上下文信息
	if err := openFile(); err != nil {
		// 使用 %w 包装原始错误，而不是覆盖它
		// 这样可以在上层继续追溯错误来源
		return fmt.Errorf("加载文件失败: %w", err)
	}
	return nil
}

func openFile() error {
	// 模拟来自标准库的底层错误
	// 这一层只关心“事实本身”，不添加额外语义
	return fs.ErrNotExist
}
