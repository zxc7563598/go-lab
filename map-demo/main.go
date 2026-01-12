package main

func main() {
	DemoReferenceButNotSafe()
	DemoNilVsMake()
	DemoPassBetweenFuncs()
	// DemoConcurrentPanic() // 会 panic，手动打开
	DemoDefensiveUsage()
}
