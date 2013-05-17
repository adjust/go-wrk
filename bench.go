package main

import (
	"sync"
	"time"
)

// type response represents the result of an HTTP request
type response struct {
	code     int
	duration time.Duration
	size     int64
}

func bench() []byte {
	work := make(chan struct{}, *totalCalls)
	responses := make(chan response)

	//TODO check ulimit
	wg := &sync.WaitGroup{}

	for i := 0; i < *totalCalls; i++ {
		work <- struct{}{}
	}
	close(work)
	t1 := time.Now()
	for i := 0; i < *numConnections; i++ {
		go StartClient(
			toCall,
			*headers,
			*method,
			*disableKeepAlives,
			*disableCompression,
			work,
			responses,
			wg,
		)
		wg.Add(1)
	}
	go func() {
		wg.Wait()
		close(responses)
	}()
	return CalcStats(responses, t1)
}
