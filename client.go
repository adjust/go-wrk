package main

import (
	"net/http"
	"strings"
	"sync"
)

func StartClient(url, heads, meth string, dka bool, responseChan chan *Response, waitGroup *sync.WaitGroup, tc int) {
	defer waitGroup.Done()

	tr := &http.Transport{DisableKeepAlives: dka}
	req, _ := http.NewRequest(meth, url, nil)
	sets := strings.Split(heads, "\n")

	//Split incoming header string by \n and build header pairs
	for i := range sets {
		split := strings.SplitN(sets[i], ":", 2)
		if len(split) == 2 {
			req.Header.Set(split[0], split[1])
		}
	}

	timer := NewTimer()
	for {
		timer.Reset()

		resp, err := tr.RoundTrip(req)

		respObj := &Response{}

		if err != nil {
			respObj.Error = true
		} else {
			respObj.Size = resp.ContentLength
			resp.Body.Close()
		}

		respObj.Duration = timer.Duration()

		if len(responseChan) >= tc {
			break
		}
		responseChan <- respObj
	}
}
