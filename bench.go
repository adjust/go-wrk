package main

import (
	"sync"
)

func bench() []byte {
	responseChannel := make(chan int, *totalCalls*2)
	countChannel := make(chan bool, *totalCalls*2)
	benchChannel := make(chan int64, *totalCalls*2)

	benchTime := NewTimer()
	benchTime.Reset()

	//TODO check ulimit
	wg := &sync.WaitGroup{}

	for i := 0; i < *numConnections; i++ {
		go StartClient(
			target,
			*headers,
			*method,
			countChannel,
			*disableKeepAlives,
			benchChannel,
			responseChannel,
			wg,
		)
		wg.Add(1)
	}

	wg.Wait()

	result := CalcStats(
		benchChannel,
		responseChannel,
		benchTime.Duration(),
	)
	return result
}
