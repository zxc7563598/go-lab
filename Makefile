# Go 语义实验仓库
#
# 这个仓库中的示例代码并不是业务代码，
# 而是一组用于验证 Go 中「值 / 指针 / 返回语义」的小实验。
#
# 每个目录都可以通过 go run 单独执行，
# Makefile 只是为了把这些实验入口集中管理。

.PHONY: help struct-value return-demo semantics slice-demo map-demo defer-demo error-demo

# 默认命令：列出所有可用实验
help:
	@echo ""
	@echo "Go语义演示"
	@echo ""
	@echo "可用命令:"
	@echo ""
	@echo "  make struct-value   # struct 的传递：值拷贝 vs 指针"
	@echo "  make return-demo    # 值接收者 vs 指针接收者，以及返回值/指针/interface 的区别"
	@echo "  make semantics      # & 和 * 的真实含义"
	@echo "  make slice-demo     # slice 的底层结构：指针、长度、容量"
	@echo "  make map-demo       # map：看起来简单，实则暗雷密布"
	@echo "  make defer-demo     # defer 与资源生命周期"
	@echo "  make error-demo     # error：Go 的“显式异常系统”"
	@echo ""

# struct 的传递：值拷贝 vs 指针
struct-value:
	go run ./struct-value

# 返回语义相关实验：
# - 值接收者 vs 指针接收者
# - 返回 struct
# - 返回指针
# - 返回 interface
return-demo:
	go run ./return-demo

# & 和 * 的真实含义：
# - 接收值 / 接收指针
# - 返回值 / 返回指针
# - 解引用发生的位置
semantics:
	go run ./semantics

# slice 的底层结构：指针、长度、容量
slice-demo:
	go run ./slice-demo

# map：看起来简单，实则暗雷密布
map-demo:
	go run ./map-demo

# defer 与资源生命周期
defer-demo:
	go run ./defer-demo

# error：Go 的“显式异常系统”
error-demo:
	go run ./error-demo