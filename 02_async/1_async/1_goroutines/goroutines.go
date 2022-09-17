package main

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	iterationsNum = 7
	goroutinesNum = 7
)

func doWork(th int) {
	for j := 0; j < iterationsNum; j++ {
		fmt.Println(formatWork(th, j))
		// time.Sleep(time.Millisecond)
		// runtime.Gosched()
	}
}

func main() {
	fmt.Println(runtime.NumCPU())
	runtime.GOMAXPROCS(1)
	for i := 0; i < goroutinesNum; i++ {
		go doWork(i)
	}
	fmt.Scanln()
}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}
