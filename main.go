package main

import (
	"flag"
	"fmt"
	"net/http"
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
	wg                *sync.WaitGroup
)

func init() {
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(*numThreads)
	url := os.Args[len(os.Args)-1]
	countChannel := make(chan bool, *totalCalls*2)
	benchChannel := make(chan int64, *totalCalls*2)
	responseChannel := make(chan int, *totalCalls*2)

	//TODO check ulimit
	wg = &sync.WaitGroup{}
	for i := 0; i < *numConnections; i++ {
		go StartClient(url, *method, countChannel, *disableKeepAlives, benchChannel, responseChannel)
		wg.Add(1)
	}
	wg.Wait()
	CalcStats(benchChannel, responseChannel, url, *numConnections, *numThreads)
}

func StartClient(u, m string, ch chan bool, dka bool, bc chan int64, rc chan int) {
	tr := &http.Transport{DisableKeepAlives: dka}
	req, _ := http.NewRequest(m, u, nil)
	timer := NewTimer()
	for {

		timer.Reset()
		resp, err := tr.RoundTrip(req)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(ch) >= *totalCalls {
			break
		}
		ch <- true
		rc <- resp.StatusCode
		bc <- timer.Duration()
		resp.Body.Close()
	}
	wg.Done()
}

func CalcStats(benchChannel chan int64, responseChannel chan int, url string, connections, threads int) {
	var total, sum int64
	for rt := range benchChannel {
		total++
		sum += rt
		if len(benchChannel) == 0 {
			break
		}
	}
	var resp200, resp300, resp400, resp500 int

	for res := range responseChannel {
		switch {
		case res < 300:
			resp200++
		case res < 400:
			resp300++
		case res < 500:
			resp400++
		case res < 600:
			resp500++
		}
		if len(responseChannel) == 0 {
			break
		}
	}
	fmt.Println("==========================BENCHMARK==========================")
	fmt.Printf("URL:\t\t\t\t%s\n\n", url)
	fmt.Printf("Used Connections:\t\t%d\n", connections)
	fmt.Printf("Used Threads:\t\t\t%d\n", threads)
	fmt.Printf("Total number of calls:\t\t%d\n\n", total)
	fmt.Println("============================TIMES============================")
	fmt.Printf("Total time passed:\t\t%ds\n", sum/1E6)
	fmt.Printf("Avg time per request:\t\t%.2fms\n", float64(sum/total)/1000)
	fmt.Printf("Median time per request:\t%.2fms\n", float64(sum/total)/1000)
	fmt.Printf("99th percentile time:\t\t%.2fms\n", float64(sum/total)/1000)
	fmt.Printf("Slowest time for request:\t%.2fms\n\n", float64(sum/total)/1000)
	fmt.Println("==========================RESPONSES==========================")
	fmt.Printf("20X responses:\t\t%d\n", resp200)
	fmt.Printf("30X responses:\t\t%d\n", resp300)
	fmt.Printf("40X responses:\t\t%d\n", resp400)
	fmt.Printf("50X responses:\t\t%d\n", resp500)
}
