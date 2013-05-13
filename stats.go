package main

import (
	"fmt"
	"sort"
)

func CalcStats(benchChannel chan int64, responseChannel chan int, url string, connections, threads int, duration int64) {
	var total, sum, i int64

	times := make([]int, len(benchChannel))
	for rt := range benchChannel {
		sum += rt
		times[i] = int(rt)
		i++
		if len(benchChannel) == 0 {
			break
		}
	}
	total = int64(len(times))
	sort.Ints(times)
	var resp200, resp300, resp400, resp500 int64

	for res := range responseChannel {
		switch {
		case res < 300:
			resp200++
		case res < 400:
			resp300++
		case res < 500:
			resp400++
		case res < 600:
			resp500++
		}
		if len(responseChannel) == 0 {
			break
		}
	}

	durationFloat := float64(duration)
	totalFloat := float64(total)
	sumFloat := float64(sum)

	fmt.Println("==========================BENCHMARK==========================")
	fmt.Printf("URL:\t\t\t\t%s\n\n", url)
	fmt.Printf("Used Connections:\t\t%d\n", connections)
	fmt.Printf("Used Threads:\t\t\t%d\n", threads)
	fmt.Printf("Total number of calls:\t\t%d\n\n", total)
	fmt.Println("============================TIMES============================")
	fmt.Printf("Total time passed:\t\t%.2fs\n", durationFloat/1E6)
	fmt.Printf("Avg time per request:\t\t%.2fms\n", sumFloat/totalFloat/1000)
	fmt.Printf("Requests per second:\t\t%.2f\n", totalFloat/(durationFloat/1E6))
	fmt.Printf("Median time per request:\t%.2fms\n", float64(times[(total-1)/2])/1000)
	fmt.Printf("99th percentile time:\t\t%.2fms\n", float64(times[(total/100*99)])/1000)
	fmt.Printf("Slowest time for request:\t%.2fms\n\n", float64(times[total-1]/1000))
	fmt.Println("==========================RESPONSES==========================")
	fmt.Printf("20X responses:\t\t%d\t(%d%%)\n", resp200, resp200/total*100)
	fmt.Printf("30X responses:\t\t%d\t(%d%%)\n", resp300, resp300/total*100)
	fmt.Printf("40X responses:\t\t%d\t(%d%%)\n", resp400, resp400/total*100)
	fmt.Printf("50X responses:\t\t%d\t(%d%%)\n", resp500, resp500/total*100)
}
