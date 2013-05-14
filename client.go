package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

func StartClient(u, h, m string, ch chan bool, dka bool, bc chan int64, rc chan int, wg *sync.WaitGroup) {
	tr := &http.Transport{DisableKeepAlives: dka}
	req, _ := http.NewRequest(m, u, nil)
	sets := strings.Split(h, "\n")
	for i := range sets {
		split := strings.SplitN(sets[i], ":", 2)
		if len(split) == 2 {
			req.Header.Set(split[0], split[1])
		}
	}
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
