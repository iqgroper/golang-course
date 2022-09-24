package main

import "fmt"

func ExecutePipeline() {

}

func SingleHash(in, out chan string) {
	data := <-in
	md5 := DataSignerMd5(data)
	crc32FromMd5 := DataSignerCrc32(md5)
	crc32FromData := DataSignerCrc32(data)
	out <- crc32FromData + "~" + crc32FromMd5
}

func MultiHash(in, out chan string) {

}

func CombineResult() {

}

func main() {
	in := make(chan string, 1)
	out := make(chan string, 1)
	in <- "0"
	SingleHash(in, out)

	fmt.Println(<-out)

}
