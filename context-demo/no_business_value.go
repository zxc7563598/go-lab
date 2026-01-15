package main

import (
	"context"
	"fmt"
)

type User struct {
	ID   int
	Name string
}

func demoNoBusinessValue() {
	fmt.Println("----------")
	fmt.Println("为什么 context 不应该传业务参数")

	// 创建一个最基础的 context
	ctx := context.Background()

	// 将业务数据塞进 context
	// 这一行在语法和运行时层面都是完全合法的，
	// 但问题不在“能不能这么写”，而在于“这样写会带来什么后果”
	ctx = context.WithValue(ctx, "user", User{ID: 1, Name: "Tom"})

	// 调用下游函数
	// 从函数签名上，你完全看不出来：
	//   1. 这个函数依赖了 user
	//   2. user 是必须存在的
	doSomething(ctx)
}

func doSomething(ctx context.Context) {
	// 仅从函数签名来看：
	//   func doSomething(ctx context.Context)
	// 你无法判断：
	//   - 这个函数是否依赖业务数据
	//   - 依赖的是哪个字段
	//   - 如果缺失会发生什么
	//
	// 这种依赖是“隐式的”，而且编译器无法帮你兜底
	user := ctx.Value("user").(User)

	// 如果上游没有正确注入 user，
	// 这里会直接 panic，而问题只会在运行时暴露
	fmt.Println("从 context 中读取到业务数据 user：", user.Name)
}
