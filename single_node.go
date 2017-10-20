package main

import (
    "sync"
    "fmt"
)

func SingleNode(toCall string, numConnections, totalCalls int) []byte {
    responseChannel := make(chan *Response, totalCalls*2)

    benchTime := NewTimer()
    benchTime.Reset()
    //TODO check ulimit
    wg := &sync.WaitGroup{}

    for i := 0; i < numConnections; i++ {
        go StartClient(
            toCall,
            *headers,
            *requestBody,
            *method,
            *disableKeepAlives,
            responseChannel,
            wg,
            totalCalls,
        )
        wg.Add(1)
    }

    fmt.Println("WAITING")
    wg.Wait()

    fmt.Println("HELELHELEHLEHLEH")

    result := CalcStats(
        responseChannel,
        benchTime.Duration(),
        toCall,
    )
    return result
}
