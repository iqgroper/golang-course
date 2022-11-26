package main

import (
	context "context"
	"fmt"
	"time"
)

func main() {
	ctx, finish := context.WithCancel(context.Background())
	err := StartMyMicroservice(ctx, "127.0.0.1:8082", "")
	if err != nil {
		fmt.Println("error in microservices")
	}
	time.Sleep(1 * time.Second)
	finish()
	fmt.Println("called finish")
}
