package main

import (
	"sync"
)

func SingleNode(toCall string) []byte {
	responseChannel := make(chan *Response, *totalCalls*2)

	benchTime := NewTimer()
	benchTime.Reset()
	//TODO check ulimit
	wg := &sync.WaitGroup{}

	for i := 0; i < *numConnections; i++ {
		go StartClient(
			toCall,
			*headers,
			*requestBody,
			*method,
			*disableKeepAlives,
			responseChannel,
			wg,
			*totalCalls,
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
