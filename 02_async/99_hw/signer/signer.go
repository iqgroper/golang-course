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

func takeMd5Hash(data string, out chan string) {
	md5Mutex.Lock()
	out <- DataSignerMd5(data)
	md5Mutex.Unlock()
}

func takeCrc32Hash(data string, out chan string) {
	out <- DataSignerCrc32(data)
}

func InnerSingleHash(in, out chan interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()
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

	result := <-outString + "~" + <-outString
	out <- result

	fmt.Println("singlehash - done", result)

}

func SingleHash(in, out chan interface{}) {
	fmt.Println("singlecall", len(in))
	singleWaiter := &sync.WaitGroup{}
	for i := 0; i < len(in); i++ {
		singleWaiter.Add(1)
		go InnerSingleHash(in, out, singleWaiter)
	}

	singleWaiter.Wait()
}

func InnerMultiHash(in, out chan interface{}, waiter *sync.WaitGroup) {
	fmt.Println("innermulticall", len(in))
	defer waiter.Done()
	dataRaw := <-in
	data, ok := dataRaw.(string)
	if !ok {
		fmt.Println("cant convert result data to string in MultiHash", dataRaw)
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
		answer += msg
	}
	out <- answer
	fmt.Println("multihash - done", answer)
}

func MultiHash(in, out chan interface{}) {
	fmt.Println("multicall", len(in))

	multiWaiter := &sync.WaitGroup{}

	for i := 0; i < 7; i++ {
		multiWaiter.Add(1)
		go InnerMultiHash(in, out, multiWaiter)
	}
	multiWaiter.Wait()
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
	fmt.Println(strings.Join(allData, "_"))
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

	waitJobs := &sync.WaitGroup{}

	channels := make([]chan interface{}, len(jobs)+1)

	for i := range channels {
		channels[i] = make(chan interface{}, 7)
	}

	for i := 0; i < len(jobs); i++ {
		waitJobs.Add(1)
		go callingJob(jobs[i], channels[i], channels[i+1], true, waitJobs)
	}
	waitJobs.Wait()

}

func main() {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	testResult := "NOT_SET"
	runtime.GOMAXPROCS(1)

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
