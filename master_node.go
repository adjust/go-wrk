package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

func MasterNode() {
	distChannel := make(chan string, len(config.Nodes)*2)
	wg := &sync.WaitGroup{}
	for _, node := range config.Nodes {
		go runChild(distChannel, wg, node)
		wg.Add(1)
	}
	wg.Wait()
	CalcDistStats(distChannel)
}

func runChild(distChan chan string, wg *sync.WaitGroup, node string) {
	defer wg.Done()
	toCall := fmt.Sprintf(
		"http://%s/t=%d&m=%s&c=%d&n=%d&k=%t&url=%s",
		node,
		*numThreads,
		*method,
		*numConnections,
		*totalCalls,
		*disableKeepAlives,
		url.QueryEscape(url.QueryEscape(target)),
	)
	fmt.Println(toCall)
	resp, err := http.Get(toCall)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	distChan <- string(body)
}
