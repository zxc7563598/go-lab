package main

import "time"

func main() {
	demoDeferStack()
	demoDeferLoopWrong()
	demoDeferLoopCorrect()
	demoDeferInWebLikeContext()
	time.Sleep(200 * time.Millisecond)
	demoExplicitResourceRelease()
}
