package main

import (
    "sync"
    "fmt"
    "strconv"
)

func SingleNode(toCall string, numConnections, totalCalls int, isWarmup bool) []byte {
    // totalCalls*2 probably so that the channel can hold resquests+responses
    responseChannel := make(chan *Response, totalCalls*2)

    benchTime := NewTimer()
    benchTime.Reset()
    //TODO check ulimit
    wg := &sync.WaitGroup{}

    // Allow reuse of TCP connection after warmup sequence
    dka := *disableKeepAlives
    if isWarmup {
        dka = false
    }

    for i := 0; i < numConnections; i++ {
        fmt.Println("Starting connection " + strconv.Itoa(i) + " to " + toCall)
        
        wg.Add(1)
        go StartClient(
            toCall,
            *headers,
            *requestBody,
            *method,
            dka,
            responseChannel,
            wg,
            totalCalls,
        )
    }

    wg.Wait()

    // initialize empty byte array incase of warmup
    result := make([]byte, 0)

    if !isWarmup {
        result = CalcStats(
            responseChannel,
            benchTime.Duration(),
            toCall,
        )
    }

    return result
}
