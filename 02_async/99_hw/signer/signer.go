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
	checkedData, ok := dataRaw.(int)
	if !ok {
		fmt.Println("input value is not int")
		return
	}
	data := strconv.Itoa(checkedData)
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
		fmt.Println("cant convert result data to string in MultiHash")
		return
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
	// fmt.Println("MultiHash result", answer)
	out <- answer
	wg.Done()
}

func CombineResults(in, out chan interface{}) {
	var allData []string
	for rawData := range in {
		data, ok := rawData.(string)
		if !ok {
			fmt.Println("cant convert result data to string in CombineResults")
			return
		}
		allData = append(allData, data)
	}

	sort.Strings(allData)
	out <- strings.Join(allData, "_")

}

func callingJob(someJob job, in, out chan interface{}, shouldClose bool, waiter *sync.WaitGroup) {
	if waiter != nil {
		defer waiter.Done()
	}

	someJob(in, out)

	if shouldClose {
		close(out)
	}
}

func ExecutePipeline(jobs ...job) {
	// in := make(chan interface{}, 100)
	// out := make(chan interface{}, 100)

	// inChannels := make([]chan interface{}, len(jobs))
	// outChannels := make([]chan interface{}, len(jobs))
	waitJobs := &sync.WaitGroup{}

	channels := make([]chan interface{}, len(jobs)+1)

	for i := range channels {
		channels[i] = make(chan interface{}, 3)
	}

	for i := 0; i < len(jobs); i++ {
		waitJobs.Add(1)
		go callingJob(jobs[i], channels[i], channels[i+1], true, waitJobs)
	}
	waitJobs.Wait()

	// jobs[0](channels[0], channels[1])
	// // close(channels[1])
	// go jobs[1](channels[1], channels[2])
	// // close(channels[2])
	// jobs[2](channels[2], nil)
	// time.Sleep(time.Millisecond)

	// go jobs[1](out, in)
	// jobs[0](in, out)

	// jobs[0](in, out)

	// inputNumber := len(out)
	// for i := 0; i < inputNumber; i++ {
	// 	go jobs[1](out, in)
	// }

	// for i := 0; i < inputNumber; i++ {
	// 	wg.Add(1)
	// 	go jobs[2](in, out)
	// }
	// wg.Wait()

	// close(out)
	// jobs[3](out, in)
	// jobs[4](in, out)
}

func main() {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	testResult := "NOT_SET"
	runtime.GOMAXPROCS(0)

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, _ := dataRaw.(string)
			testResult = data
		}),
	}
	start := time.Now()

	ExecutePipeline(hashSignJobs...)

	end := time.Since(start)
	fmt.Println(end)
	fmt.Println(testResult)
}
