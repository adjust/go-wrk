package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
)

type Stats struct {
	Url         string
	Connections int
	Threads     int
	AvgDuration float64
	Duration    float64
	Sum         float64
	Times       []int
	Transferred int64
	Resp200     int64
	Resp300     int64
	Resp400     int64
	Resp500     int64
	Errors      int64
	Contains    int64
}

func CalcStats(responseChannel chan *Response, duration int64) []byte {

	stats := &Stats{
		Url:         target,
		Connections: *numConnections,
		Threads:     *numThreads,
		Times:       make([]int, len(responseChannel)),
		Duration:    float64(duration),
		AvgDuration: float64(duration),
	}

	if *respContains != "" {
		log.Printf("search in response for: %v", *respContains)
	}

	i := 0
	for res := range responseChannel {
		switch {
		case res.StatusCode < 200:
			// error
		case res.StatusCode < 300:
			stats.Resp200++
		case res.StatusCode < 400:
			stats.Resp300++
		case res.StatusCode < 500:
			stats.Resp400++
		case res.StatusCode < 600:
			stats.Resp500++
		}

		if *respContains != "" && strings.Contains(res.Body, *respContains) {
			stats.Contains++
		}

		stats.Sum += float64(res.Duration)
		stats.Times[i] = int(res.Duration)
		i++

		stats.Transferred += res.Size

		if res.Error {
			stats.Errors++
		}

		if len(responseChannel) == 0 {
			break
		}
	}

	sort.Ints(stats.Times)

	PrintStats(stats)
	b, err := json.Marshal(&stats)
	if err != nil {
		fmt.Println(err)
	}
	return b
}

func CalcDistStats(distChan chan string) {
	if len(distChan) == 0 {
		return
	}
	allStats := &Stats{
		Url:         target,
		Connections: *numConnections,
		Threads:     *numThreads,
	}
	statCount := len(distChan)
	for res := range distChan {
		var stats Stats
		err := json.Unmarshal([]byte(res), &stats)
		if err != nil {
			fmt.Println(err)
		}
		allStats.Duration += stats.Duration
		allStats.Sum += stats.Sum
		allStats.Times = append(allStats.Times, stats.Times...)
		allStats.Resp200 += stats.Resp200
		allStats.Resp300 += stats.Resp300
		allStats.Resp400 += stats.Resp400
		allStats.Resp500 += stats.Resp500
		allStats.Errors += stats.Errors
		allStats.Contains += stats.Contains

		if len(distChan) == 0 {
			break
		}
	}
	allStats.AvgDuration = allStats.Duration / float64(statCount)
	PrintStats(allStats)
}

func PrintStats(allStats *Stats) {
	sort.Ints(allStats.Times)
	total := float64(len(allStats.Times))
	totalInt := int64(total)
	fmt.Println("==========================BENCHMARK==========================")
	fmt.Printf("URL:\t\t\t\t%s\n\n", allStats.Url)
	fmt.Printf("Used Connections:\t\t%d\n", allStats.Connections)
	fmt.Printf("Used Threads:\t\t\t%d\n", allStats.Threads)
	fmt.Printf("Total number of calls:\t\t%d\n\n", totalInt)
	fmt.Println("===========================TIMINGS===========================")
	fmt.Printf("Total time passed:\t\t%.2fs\n", allStats.AvgDuration/1E6)
	fmt.Printf("Avg time per request:\t\t%.2fms\n", allStats.Sum/total/1000)
	fmt.Printf("Requests per second:\t\t%.2f\n", total/(allStats.AvgDuration/1E6))
	fmt.Printf("Median time per request:\t%.2fms\n", float64(allStats.Times[(totalInt-1)/2])/1000)
	fmt.Printf("99th percentile time:\t\t%.2fms\n", float64(allStats.Times[(totalInt/100*99)])/1000)
	fmt.Printf("Slowest time for request:\t%.2fms\n\n", float64(allStats.Times[totalInt-1]/1000))
	fmt.Println("=============================DATA=============================")
	fmt.Printf("Total response body sizes:\t\t%d\n", allStats.Transferred)
	fmt.Printf("Avg response body per request:\t\t%.2f Byte\n", float64(allStats.Transferred)/total)
	tr := float64(allStats.Transferred) / (allStats.AvgDuration / 1E6)
	fmt.Printf("Transfer rate per second:\t\t%.2f Byte/s (%.2f MByte/s)\n", tr, tr/1E6)
	fmt.Println("==========================RESPONSES==========================")
	fmt.Printf("20X Responses:\t\t%d\t(%.2f%%)\n", allStats.Resp200, float64(allStats.Resp200)/total*1e2)
	fmt.Printf("30X Responses:\t\t%d\t(%.2f%%)\n", allStats.Resp300, float64(allStats.Resp300)/total*1e2)
	fmt.Printf("40X Responses:\t\t%d\t(%.2f%%)\n", allStats.Resp400, float64(allStats.Resp400)/total*1e2)
	fmt.Printf("50X Responses:\t\t%d\t(%.2f%%)\n", allStats.Resp500, float64(allStats.Resp500)/total*1e2)
	if *respContains != "" {
		fmt.Printf("matchResponses:\t\t%d\t(%.2f%%)\n", allStats.Contains, float64(allStats.Contains)/total*1e2)
	}
	fmt.Printf("Errors:\t\t\t%d\t(%.2f%%)\n", allStats.Errors, float64(allStats.Errors)/total*1e2)
}
