package main

import (
	"fmt"
)

func demoPropagateError() {
	fmt.Println("----------")
	fmt.Println("错误如何在调用链中逐层向上传递")

	// 从最外层发起调用，只关心“最终是否成功”
	err := service()
	if err != nil {
		fmt.Println("错误一路向上传递后，在最外层被观察到:", err)
	}
}

func service() error {
	// service 层负责编排业务流程
	// 并不直接处理底层失败的细节
	if err := repository(); err != nil {
		// 在无法修复错误的情况下
		// 这一层做的事情是补充当前业务语境
		return fmt.Errorf("service 层执行失败: %w", err)
	}
	return nil
}

func repository() error {
	// repository 层负责数据访问
	// 如果底层出错，它会说明“在做什么操作时失败了”
	if err := queryDB(); err != nil {
		return fmt.Errorf("repository 查询数据失败: %w", err)
	}
	return nil
}

func queryDB() error {
	// 最底层只返回“事实性错误”
	// 不关心调用者是谁、会如何处理
	return fmt.Errorf("数据库连接被拒绝")
}
