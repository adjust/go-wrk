package main

import (
	"flag"
	"log"
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
	disableKeepAlives = flag.Bool("k", false, "if keep-alives are disabled")
	url               string
	countChannel      chan bool
	wg                *sync.WaitGroup
)

func init() {
	flag.Parse()
	countChannel = make(chan bool, *totalCalls*2)
	url = os.Args[len(os.Args)-1]
}

func main() {
	runtime.GOMAXPROCS(*numThreads)
	//TODO check ulimit
	wg = &sync.WaitGroup{}
	for i := 0; i < *numConnections; i++ {
		go StartClient(url, *method, countChannel, *disableKeepAlives)
		wg.Add(1)
	}
	wg.Wait()
}

func StartClient(u, m string, ch chan bool, dka bool) {
	tr := &http.Transport{DisableKeepAlives: dka}
	req, _ := http.NewRequest(m, u, nil)
	for {
		resp, err := tr.RoundTrip(req)
		if err != nil {
			log.Println(err)
			continue
		}
		ch <- true
		if resp.StatusCode != 200 {
			//	log.Println(resp.StatusCode)
		}
		resp.Body.Close()
		if len(ch) >= *totalCalls {
			break
		}
	}
	wg.Done()
}
