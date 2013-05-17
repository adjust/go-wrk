package main

import (
	"time"
)

type Timer struct {
	start time.Time
}

func NewTimer() *Timer {
	timer := &Timer{}
	return timer
}

func (timer *Timer) Reset() {
	timer.start = time.Now().UTC()
}

func (timer *Timer) Duration() int64 {
	now := time.Now().UTC()
	nanos := now.Sub(timer.start).Nanoseconds()
	micros := nanos / 1000
	return micros
}
