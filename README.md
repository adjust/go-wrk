# go-wrk

    -c=100: the max numbers of connections used
    -k=true: if keep-alives are disabled
    -m="GET": the http request method
    -n=1000: the total number of calls processed
    -t=1: the numbers of threads used


## example output

 ```
 ==========================BENCHMARK==========================
URL:				http://localhost:8509/startup?app_id=479516143&mac=123456789

Used Connections:		100
Used Threads:			1
Total number of calls:		100000

============================TIMES============================
Total time passed:		19.47s
Avg time per request:		19.45ms
Requests per second:		5135.02
Median time per request:	11.30ms
99th percentile time:		65.23ms
Slowest time for request:	1698.00ms

==========================RESPONSES==========================
20X responses:		100000	(100%)
30X responses:		0	(0%)
40X responses:		0	(0%)
50X responses:		0	(0%)
```