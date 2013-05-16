package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

func StartClient(u, h, m string, dka bool, work chan struct{}, rc chan response, wg *sync.WaitGroup) {
	defer wg.Done()
	tr := &http.Transport{DisableKeepAlives: dka}
	req, err := http.NewRequest(m, u, nil)
	if err != nil {
		panic(err)
	}
	sets := strings.Split(h, "\n")
	for i := range sets {
		split := strings.SplitN(sets[i], ":", 2)
		if len(split) == 2 {
			req.Header.Set(split[0], split[1])
		}
	}
	for {
		if _, ok := <-work; !ok {
			return
		}
		t1 := time.Now()
		resp, err := tr.RoundTrip(req)
		if err != nil {
			fmt.Println(err)
			continue
		}
		size, err := io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
			fmt.Println("error reading response:", err)
		}
		rc <- response{resp.StatusCode, time.Now().Sub(t1), size}
		resp.Body.Close()
	}
}
