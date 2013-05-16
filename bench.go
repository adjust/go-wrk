package main

import (
	"sync"
)

// type response represents the result of an HTTP request
type response struct {
	code           int
	duration, size int64
}

func bench() []byte {
	responseChannel := make(chan response, *totalCalls*2)

	benchTime := NewTimer()
	benchTime.Reset()

	//TODO check ulimit
	wg := &sync.WaitGroup{}

	for i := 0; i < *numConnections; i++ {
		go StartClient(
			target,
			*headers,
			*method,
			*disableKeepAlives,
			responseChannel,
			wg,
		)
		wg.Add(1)
	}
	wg.Wait()

	result := CalcStats(
		responseChannel,
		benchTime.Duration(),
	)
	return result
}
