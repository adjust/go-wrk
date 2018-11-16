package main

type Response struct {
	Size       int64
	Duration   int64
	StatusCode int
	Error      bool
	Body       string
}
