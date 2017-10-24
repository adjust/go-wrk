package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "runtime"
)

type TPSReport struct {
    Urls []string
    TotalCalls []int
    Threads int
    NumConnections []int
}
var (
    numThreads        = flag.Int("t", 1, "the numbers of threads used")
    method            = flag.String("m", "GET", "the http request method")
    requestBody       = flag.String("b", "", "the http requst body")
    numConnections    = flag.Int("c", 100, "the max numbers of connections used")
    totalCalls        = flag.Int("n", 1000, "the total number of calls processed")
    disableKeepAlives = flag.Bool("k", true, "if keep-alives are disabled")
    configFile        = flag.String("f", "", "json config file")
    headers           = flag.String("H", "User-Agent: go-wrk 0.1 bechmark\nContent-Type: text/html;", "the http headers sent separated by '\\n'")
    certFile          = flag.String("cert", "someCertFile", "A PEM eoncoded certificate file.")
    keyFile           = flag.String("key", "someKeyFile", "A PEM encoded private key file.")
    caFile            = flag.String("CA", "someCertCAFile", "A PEM eoncoded CA's certificate file.")
    insecure          = flag.Bool("i", true, "TLS checks are disabled")
    tps TPSReport
)


func init() {
    //flag.Parse()
    //target = os.Args[len(os.Args)-1]
    configFile := os.Args[len(os.Args)-1]
    if configFile != "" {
        readConfig(configFile)
    }
    runtime.GOMAXPROCS(tps.Threads)
}

func readConfig(configFile string) {
    configData, err := ioutil.ReadFile(configFile)
    if err != nil {
        fmt.Println(err)
    }
    err = json.Unmarshal(configData, &tps)
    if err != nil {
        fmt.Println(err)
    }
}

func main() {
    // warmup cache on first route
    // TODO: may want to make this more general in case Urls[0] is not always the first one hit
    fmt.Println("Warming up cache on route " + tps.Urls[0])
    SingleNode(tps.Urls[0], 1, 10, true)
    fmt.Println("Warmup complete")

    for i, url := range tps.Urls {
        go SingleNode(url, tps.NumConnections[i], tps.TotalCalls[i], false)
    }
    for {}
}
