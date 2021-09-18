package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

type Config struct {
	Port  string
	Nodes []string
}

var (
	numThreads        = flag.Int("t", 1, "the numbers of threads used")
	method            = flag.String("m", "GET", "the http request method")
	requestBody       = flag.String("b", "", "the http request body")
	requestBodyFile   = flag.String("p", "", "the http request body data file")
	numConnections    = flag.Int("c", 100, "the max numbers of connections used")
	totalCalls        = flag.Int("n", 1000, "the total number of calls processed")
	disableKeepAlives = flag.Bool("k", true, "if keep-alives are disabled")
	dist              = flag.String("d", "", "dist mode")
	configFile        = flag.String("f", "", "json config file")
	config            Config
	target            string
	headers           = flag.String("H", "User-Agent: go-wrk 0.1 benchmark\nContent-Type: text/html;", "the http headers sent separated by '\\n'")
	certFile          = flag.String("cert", "someCertFile", "A PEM encoded certificate file.")
	keyFile           = flag.String("key", "someKeyFile", "A PEM encoded private key file.")
	caFile            = flag.String("CA", "someCertCAFile", "A PEM encode CA's certificate file.")
	insecure          = flag.Bool("i", false, "TLS checks are disabled")
	respContains      = flag.String("s", "", "if specified, it counts how often the searched string s is contained in the responses")
	readAll           = flag.Bool("r", false, "in the case of having stream or file in the response,\n it reads all response body to calculate the response size")
)

func init() {
	flag.Parse()
	target = os.Args[len(os.Args)-1]
	if *configFile != "" {
		readConfig()
	}
	runtime.GOMAXPROCS(*numThreads)
}

func readConfig() {
	configData, err := ioutil.ReadFile(*configFile)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	err = json.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func setRequestBody() {
	// requestBody has been setup and it has highest priority
	if *requestBody != "" {
		return
	}

	if *requestBodyFile == "" {
		return
	}

	// requestBodyFile has been setup
	data, err := ioutil.ReadFile(*requestBodyFile)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	body := string(data)
	requestBody = &body
}

func main() {
	if len(os.Args) == 1 { //If no argument specified
		flag.Usage() //Print usage
		os.Exit(1)
	}
	setRequestBody()
	switch *dist {
	case "m":
		MasterNode()
	case "s":
		SlaveNode()
	default:
		SingleNode(target)
	}
}
