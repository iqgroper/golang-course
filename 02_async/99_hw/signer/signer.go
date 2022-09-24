package main

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var md5Mutex = &sync.Mutex{}
var wg = &sync.WaitGroup{}

func takeMd5Hash(data string, out chan string) {
	md5Mutex.Lock()
	out <- DataSignerMd5(data)
	md5Mutex.Unlock()
}

func takeCrc32Hash(data string, out chan string) {
	out <- DataSignerCrc32(data)
}

func SingleHash(in, out chan interface{}) {
	dataRaw := <-in
	data, ok := dataRaw.(string)
	if !ok {
		fmt.Println("cant convert result data to string")
	}
	outString := make(chan string, 1)
	go takeMd5Hash(data, outString)
	go takeCrc32Hash(data, outString)
	go takeCrc32Hash(<-outString, outString)

	out <- <-outString + "~" + <-outString
}

func MultiHash(in, out chan interface{}) {
	dataRaw := <-in
	data, ok := dataRaw.(string)
	if !ok {
		fmt.Println("cant convert result data to string")
	}
	outString := make(chan string, 6)

	ticker := time.NewTicker(5 * time.Millisecond)
	i := 0
	for tickTime := range ticker.C {
		go takeCrc32Hash(strconv.Itoa(i)+data, outString)
		if i == 5 {
			ticker.Stop()
			break
			fmt.Println(tickTime)
		}
		i++
	}

	var answer string
	for i := 0; i < 6; i++ {
		msg := <-outString
		// fmt.Println(i, msg)
		answer += msg
	}
	fmt.Println("MultiHash result", answer)
	out <- answer
	wg.Done()
}

func CombineResults(in, out chan interface{}) {
	var allData []string
	for rawData := range in {
		data, ok := rawData.(string)
		if !ok {
			fmt.Println("cant convert result data to string")
		}
		allData = append(allData, data)
	}

	sort.Strings(allData)
	out <- strings.Join(allData, "_")

}

// func ExecutePipeline(funcs ...job) {
// }

func main() {
	runtime.GOMAXPROCS(0)
	start := time.Now()
	// ExecutePipeline()
	in := make(chan interface{}, 7)
	out := make(chan interface{}, 7)
	in <- "0"
	in <- "1"
	in <- "1"
	in <- "1"
	in <- "1"
	in <- "1"
	in <- "1"

	go SingleHash(in, out)
	go SingleHash(in, out)
	go SingleHash(in, out)
	go SingleHash(in, out)
	go SingleHash(in, out)
	go SingleHash(in, out)
	go SingleHash(in, out)

	// go MultiHash(out, in)
	wg.Add(7)
	go MultiHash(out, in)
	go MultiHash(out, in)
	go MultiHash(out, in)
	go MultiHash(out, in)
	go MultiHash(out, in)
	go MultiHash(out, in)
	go MultiHash(out, in)
	// time.Sleep(5 * time.Second)
	wg.Wait()
	close(in)
	CombineResults(in, out)
	fmt.Println(<-out)

	end := time.Since(start)
	fmt.Println(end)
	// time.Sleep(2 * time.Second)
	// fmt.Println(<-out)

	// for msg := range in {
	// 	fmt.Println(msg)
	// }
}
