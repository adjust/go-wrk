# go-wrk 0.1

this is a small http benchmark utility similar to https://github.com/wg/wrk but written in go.
it has a couple of features absent from wrk
 
  - https support (quite expensive on the client side with disabled keep alives)
  - http POST support
  - more statistics
  - leaner codebase

## status

this tool is in early stage development but stable enough to run larger benchmark sets.
missing features will be added as needed, pull requests are welcome ;)

## building

you need go 1.0+ (1.1 is suggested for performance)

```
git clone git://github.com/adeven/go-wrk.git
cd go-wrk
go build
```

## usage

basic usage is quite simple:
```
go-wrk [flags] url
```

with the flags being
```
    -H="User-Agent: go-wrk 0.1 bechmark\nContent-Type: text/html;": the http headers sent separated by '\n'
    -c=100: the max numbers of connections used
    -k=true: if keep-alives are disabled
    -i=false: if TLS security checks are disabled
    -m="GET": the http request method
    -n=1000: the total number of calls processed
    -t=1: the numbers of threads used
    -b="" the http request body
```
for example
```
go-wrk -c=400 -t=8 -n=100000 http://localhost:8080/index.html
```


## example output

 ```
==========================BENCHMARK==========================
URL:				http://localhost:8509/startup?app_id=479516143&mac=123456789

Used Connections:			100
Used Threads:				1
Total number of calls:		100000

============================TIMES============================
Total time passed:			19.47s
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

## License

This Software is licensed under the MIT License.

Copyright (c) 2013 adeven GmbH,
http://www.adeven.com

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
