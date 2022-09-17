package main

import (
	"fmt"
	"time"
)

var totalOperations int32 = 0

func inc() {
	// не атомарная операция
	totalOperations++
}

func main() {
	// runtime.GOMAXPROCS(1)
	for i := 0; i < 1000; i++ {
		go inc()
	}
	time.Sleep(20 * time.Millisecond)
	// ождается 1000
	fmt.Println("total operation = ", totalOperations)
}
