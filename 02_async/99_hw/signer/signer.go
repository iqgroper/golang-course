package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var md5Mutex = &sync.Mutex{}

func takeMd5Hash(data string, out chan string) {
	md5Mutex.Lock()
	result := DataSignerMd5(data)
	md5Mutex.Unlock()
	out <- result
}

func takeCrc32Hash(data string, out chan string) {
	out <- DataSignerCrc32(data)
}

func InnerSingleHash(in int, out chan interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()

	data := strconv.Itoa(in)
	outStringCrc := make(chan string)
	outStringMdCrc := make(chan string, 1)

	go takeCrc32Hash(data, outStringCrc)
	takeMd5Hash(data, outStringMdCrc)
	takeCrc32Hash(<-outStringMdCrc, outStringMdCrc)

	result := <-outStringCrc + "~" + <-outStringMdCrc
	out <- result

	// fmt.Println("singlehash - done", result)

}

func SingleHash(in, out chan interface{}) {
	singleWaiter := &sync.WaitGroup{}
	for dataRaw := range in {
		checkedData, ok := dataRaw.(int)
		if !ok {
			fmt.Println("input value is not int")
			return
		}
		singleWaiter.Add(1)
		go InnerSingleHash(checkedData, out, singleWaiter)
	}

	singleWaiter.Wait()
}

func InnerMultiHash(in string, out chan interface{}, waiter *sync.WaitGroup) {
	// fmt.Println("innermulticall", len(in))
	defer waiter.Done()
	data := in

	outStrings := make([]chan string, 6)

	for i := range outStrings {
		outStrings[i] = make(chan string)
	}

	for i := 0; i < 6; i++ {
		go takeCrc32Hash(strconv.Itoa(i)+data, outStrings[i])
	}

	var answer string
	for i := 0; i < 6; i++ {
		msg := <-outStrings[i]
		answer += msg
	}
	out <- answer
	// fmt.Println("multihash - done", answer)
}

func MultiHash(in, out chan interface{}) {
	// fmt.Println("multicall", len(in))

	multiWaiter := &sync.WaitGroup{}

	for itemRaw := range in {
		item, ok := itemRaw.(string)
		if !ok {
			fmt.Println("cant convert result data to string in MultiHash", itemRaw)
			return
		}
		multiWaiter.Add(1)
		go InnerMultiHash(item, out, multiWaiter)
	}
	multiWaiter.Wait()
}

func CombineResults(in, out chan interface{}) {
	allData := make([]string, 0)
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
	// fmt.Println(strings.Join(allData, "_"))
}

func callingJob(someJob job, in, out chan interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()

	someJob(in, out)

	close(out)
}

func ExecutePipeline(jobs ...job) {

	waitJobs := &sync.WaitGroup{}

	channels := make([]chan interface{}, len(jobs)+1)

	for i := range channels {
		channels[i] = make(chan interface{})
	}

	for i := 0; i < len(jobs); i++ {
		waitJobs.Add(1)
		go callingJob(jobs[i], channels[i], channels[i+1], waitJobs)
	}
	waitJobs.Wait()
}

func main() {

	// inputData := []int{}
	// var inputData [100]int
	// inputData := []int{0, 1, 1, 2, 3, 5, 8}
	inputData := []int{0, 1}

	testResult := "NOT_SET"
	// runtime.GOMAXPROCS(0)

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
