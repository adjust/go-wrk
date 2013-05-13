package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sort"
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
	benchTime := NewTimer()
	benchTime.Reset()
	//TODO check ulimit
	wg = &sync.WaitGroup{}
	for i := 0; i < *numConnections; i++ {
		go StartClient(url, *method, countChannel, *disableKeepAlives, benchChannel, responseChannel)
		wg.Add(1)
	}
	wg.Wait()
	CalcStats(benchChannel, responseChannel, url, *numConnections, *numThreads, benchTime.Duration())
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

func CalcStats(benchChannel chan int64, responseChannel chan int, url string, connections, threads int, duration int64) {
	var total, sum, i int64

	times := make([]int, len(benchChannel))
	for rt := range benchChannel {
		sum += rt
		times[i] = int(rt)
		i++
		if len(benchChannel) == 0 {
			break
		}
	}
	total = int64(len(times))
	sort.Ints(times)
	var resp200, resp300, resp400, resp500 int64

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
	durationFloat := float64(duration)
	totalFloat := float64(total)
	sumFloat := float64(sum)

	fmt.Println("==========================BENCHMARK==========================")
	fmt.Printf("URL:\t\t\t\t%s\n\n", url)
	fmt.Printf("Used Connections:\t\t%d\n", connections)
	fmt.Printf("Used Threads:\t\t\t%d\n", threads)
	fmt.Printf("Total number of calls:\t\t%d\n\n", total)
	fmt.Println("============================TIMES============================")
	fmt.Printf("Total time passed:\t\t%.2fs\n", durationFloat/1E6)
	fmt.Printf("Avg time per request:\t\t%.2fms\n", sumFloat/totalFloat/1000)
	fmt.Printf("Requests per second:\t\t%.2f\n", totalFloat/(durationFloat/1E6))
	fmt.Printf("Median time per request:\t%.2fms\n", float64(times[(total-1)/2])/1000)
	fmt.Printf("99th percentile time:\t\t%.2fms\n", float64(times[(total/100*99)])/1000)
	fmt.Printf("Slowest time for request:\t%.2fms\n\n", float64(times[total-1]/1000))
	fmt.Println("==========================RESPONSES==========================")
	fmt.Printf("20X responses:\t\t%d\t(%d%%)\n", resp200, resp200/total*100)
	fmt.Printf("30X responses:\t\t%d\t(%d%%)\n", resp300, resp300/total*100)
	fmt.Printf("40X responses:\t\t%d\t(%d%%)\n", resp400, resp400/total*100)
	fmt.Printf("50X responses:\t\t%d\t(%d%%)\n", resp500, resp500/total*100)
}
