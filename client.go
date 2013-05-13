package main

import (
	"fmt"
	"net/http"
	"sync"
)

func StartClient(u, m string, ch chan bool, dka bool, bc chan int64, rc chan int, wg *sync.WaitGroup) {
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
