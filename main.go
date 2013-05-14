package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
)

var (
	numThreads        = flag.Int("t", 1, "the numbers of threads used")
	method            = flag.String("m", "GET", "the http request method")
	numConnections    = flag.Int("c", 100, "the max numbers of connections used")
	totalCalls        = flag.Int("n", 1000, "the total number of calls processed")
	disableKeepAlives = flag.Bool("k", true, "if keep-alives are disabled")
	userAgent         = flag.String("u", "go-wrk 0.1 bechmark", "the user agent sent with each request")
)

func init() {
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(*numThreads)
	if len(os.Args) == 1 {
		fmt.Println("please enter an url")
		os.Exit(0)
	}
	url := os.Args[len(os.Args)-1]
	countChannel := make(chan bool, *totalCalls*2)
	benchChannel := make(chan int64, *totalCalls*2)
	responseChannel := make(chan int, *totalCalls*2)
	benchTime := NewTimer()
	benchTime.Reset()
	//TODO check ulimit
	wg := &sync.WaitGroup{}
	for i := 0; i < *numConnections; i++ {
		go StartClient(url, *userAgent, *method, countChannel, *disableKeepAlives, benchChannel, responseChannel, wg)
		wg.Add(1)
	}
	wg.Wait()
	CalcStats(benchChannel, responseChannel, url, *numConnections, *numThreads, benchTime.Duration())
}
