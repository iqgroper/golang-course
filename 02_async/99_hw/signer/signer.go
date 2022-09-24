package main

import (
	"fmt"
	"strconv"
	"time"
)

func ExecutePipeline() {

}

func takeMd5Hash(data string, out chan string) {
	out <- DataSignerMd5(data)
}

func takeCrc32Hash(data string, out chan string) {
	out <- DataSignerCrc32(data)
}

func SingleHash(in, out chan string) {
	data := <-in
	go takeMd5Hash(data, out)
	go takeCrc32Hash(data, out)
	go takeCrc32Hash(<-out, out)
	out <- <-out + "~" + <-out
}

func MultiHash(in, out chan string) {
	data := <-in
	// for i := 0; i < 6; i++ {
	// 	go takeCrc32Hash(strconv.Itoa(i)+data, out)
	// }
	ticker := time.NewTicker(10 * time.Millisecond)
	i := 0
	for tickTime := range ticker.C {
		go takeCrc32Hash(strconv.Itoa(i)+data, out)
		if i == 5 {
			ticker.Stop()
			break
			fmt.Println(tickTime)
		}
		i++
	}

	var answer string
	for i := 0; i < 6; i++ {
		msg := <-out
		fmt.Println(i, msg)
		answer += msg
	}
	fmt.Println("MultiHash result", answer)
	out <- answer
}

func CombineResult(in, out chan string) {

}

func main() {
	start := time.Now()
	in := make(chan string, 6)
	out := make(chan string, 6)
	in <- "0"
	SingleHash(in, out)
	// SingleHashOut := <-out
	// in <- SingleHashOut
	MultiHash(out, in)
	// CombineResult(in, out)
	end := time.Since(start)
	fmt.Println(end)

	// fmt.Println(SingleHashOut)

	// for msg := range in {
	// 	fmt.Println(msg)
	// }
}
