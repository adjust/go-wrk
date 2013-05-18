package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func SlaveNode() {
	http.HandleFunc("/", rootHandler)
	err := http.ListenAndServe(
		fmt.Sprintf(":%s", config.Port),
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	select {}
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	values, err := url.ParseQuery(req.URL.String()[1:])
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	*numThreads, _ = strconv.Atoi(values.Get("t"))
	*method = values.Get("m")
	*numConnections, _ = strconv.Atoi(values.Get("c"))
	*totalCalls, _ = strconv.Atoi(values.Get("n"))
	*disableKeepAlives, _ = strconv.ParseBool(values.Get("k"))
	toCall, _ := url.QueryUnescape(values.Get("url"))
	fmt.Fprintf(w, string(SingleNode(toCall)))
}
