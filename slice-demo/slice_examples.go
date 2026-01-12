package main

import "fmt"

// slice 的几种生成方式
func demoSliceCreate() {
	fmt.Println("----------")
	fmt.Println("slice 的几种生成方式")

	// 1. 声明一个 nil slice
	// 此时 slice 本身存在，但没有指向任何底层数组
	var s1 []int
	fmt.Println("1. nil slice")
	fmt.Println("s1:", s1, "len:", len(s1), "cap:", cap(s1))

	// 2. 使用字面量声明
	// 会同时分配底层数组，len 和 cap 初始相等
	s2 := []int{1, 2, 3}
	fmt.Println("2. 字面量 slice")
	fmt.Println("s2:", s2, "len:", len(s2), "cap:", cap(s2))

	// 3. 使用 make 预分配容量
	// len 表示当前可用元素数量
	// cap 表示底层数组的总容量
	s3 := make([]int, 0, 5)
	fmt.Println("3. make 预设容量")
	fmt.Println("s3:", s3, "len:", len(s3), "cap:", cap(s3))
}

// 从 array 切出 slice
func demoSliceFromArray() {
	fmt.Println("----------")
	fmt.Println("从 array 切出 slice")

	// array 的长度是类型的一部分，内存一次性分配
	arr := [5]int{1, 2, 3, 4, 5}

	// 从下标 1（包含）到 4（不包含）
	// slice 会与 array 共享同一块底层内存
	s := arr[1:4]

	fmt.Println("arr:", arr)
	fmt.Println("s  :", s)
	fmt.Println("s len:", len(s), "cap:", cap(s))
}

// 对 slice 再进行切片
func demoSliceSlice() {
	fmt.Println("----------")
	fmt.Println("slice 再切 slice")

	s := []int{10, 20, 30, 40, 50}

	// 所有新的 slice 都仍然指向同一个底层数组
	s1 := s[:2]  // 前两个元素
	s2 := s[2:]  // 从下标 2 到结尾
	s3 := s[1:4] // 中间一段

	fmt.Println("s  :", s)
	fmt.Println("s1(前两个):", s1, "len:", len(s1), "cap:", cap(s1))
	fmt.Println("s2(从下标2到结尾):", s2, "len:", len(s2), "cap:", cap(s2))
	fmt.Println("s3(中间一段):", s3, "len:", len(s3), "cap:", cap(s3))
}

// append 过程中 len / cap 的变化
func demoAppendGrowth() {
	fmt.Println("----------")
	fmt.Println("append 过程中 len / cap 的变化")

	// 从一个空 slice 开始
	s := []int{}

	for i := 0; i < 10; i++ {
		s = append(s, i)

		// len 表示当前元素个数
		// cap 表示当前底层数组还能容纳的最大元素数量
		fmt.Printf("append %d -> len=%d cap=%d\n", i, len(s), cap(s))
	}
}

// append 是否影响原 slice（是否共享底层数组）
func demoAppendShare() {
	fmt.Println("----------")
	fmt.Println("slice 是否共享底层数组")

	// 预分配 cap=5 的 slice
	s := make([]int, 0, 5)
	s = append(s, 1, 2, 3)

	// a 与 s 共享同一底层数组
	a := s[:2]

	fmt.Println("追加前")
	fmt.Println("s:", s, "len:", len(s), "cap:", cap(s))
	fmt.Println("a:", a, "len:", len(a), "cap:", cap(a))

	// 由于没有超过 cap，这次 append 会修改共享的底层数组
	a = append(a, 4)

	fmt.Println("不超过 cap 的追加后（仍共享底层数组）")
	fmt.Println("s:", s, "len:", len(s), "cap:", cap(s))
	fmt.Println("a:", a, "len:", len(a), "cap:", cap(a))

	// 连续 append，最终会超过 cap
	// 此时 Go 会为 a 分配新的底层数组
	a = append(a, 5)
	a = append(a, 6)
	a = append(a, 7)

	fmt.Println("超过 cap 的追加后（a 发生了内存重新分配）")
	fmt.Println("s:", s, "len:", len(s), "cap:", cap(s))
	fmt.Println("a:", a, "len:", len(a), "cap:", cap(a))
}
