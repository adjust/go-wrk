package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func StartClient(url_, heads, requestBody string, meth string, dka bool, responseChan chan *Response, waitGroup *sync.WaitGroup, tc int) {
	defer waitGroup.Done()

	var tr *http.Transport

	u, err := url.Parse(url_)

	if err == nil && u.Scheme == "https" {
		var tlsConfig *tls.Config
		if *insecure {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		} else {
			// Load client cert
			cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
			if err != nil {
				log.Fatal(err)
			}

			// Load CA cert
			caCert, err := ioutil.ReadFile(*caFile)
			if err != nil {
				log.Fatal(err)
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			// Setup HTTPS client
			tlsConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      caCertPool,
			}
			tlsConfig.BuildNameToCertificate()
		}

		tr = &http.Transport{TLSClientConfig: tlsConfig, DisableKeepAlives: dka}
	} else {
		tr = &http.Transport{DisableKeepAlives: dka}
	}

	timer := NewTimer()
	for {
		requestBodyReader := strings.NewReader(requestBody)
		req, _ := http.NewRequest(meth, url_, requestBodyReader)
		sets := strings.Split(heads, "\n")

		//Split incoming header string by \n and build header pairs
		for i := range sets {
			split := strings.SplitN(sets[i], ":", 2)
			if len(split) == 2 {
				req.Header.Set(split[0], split[1])
			}
		}

		timer.Reset()

		resp, err := tr.RoundTrip(req)

		respObj := &Response{}

		if err != nil {
			respObj.Error = true
		} else {
			if resp.ContentLength < 0 { // -1 if the length is unknown
				data, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					respObj.Size = int64(len(data))
				}
			} else {
				respObj.Size = resp.ContentLength
				if *respContains != "" {
					data, err := ioutil.ReadAll(resp.Body)
					if err == nil {
						respObj.Body = string(data)
					}
				}
			}
			respObj.StatusCode = resp.StatusCode

			resp.Body.Close()
		}

		respObj.Duration = timer.Duration()

		if len(responseChan) >= tc {
			break
		}
		responseChan <- respObj
	}
}
