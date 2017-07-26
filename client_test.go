package main

import (
	"net/http"
	"reflect"
	"sort"
	"testing"
)

func TestBuildHeaders(t *testing.T) {

	var testSet = []struct {
		H   string // headers passed via flag
		Res map[string][]string
	}{
		{
			"User-Agent:go-wrk 0.1 bechmark\nContent-Type:text/html;",
			map[string][]string{
				"User-Agent":   []string{"go-wrk 0.1 bechmark"},
				"Content-Type": []string{"text/html;"},
			},
		},
		{
			"Key:Value",
			map[string][]string{
				"Key": []string{"Value"},
			},
		},
		{
			"Key1:Value1\nKey2:Value2",
			map[string][]string{
				"Key1": []string{"Value1"},
				"Key2": []string{"Value2"},
			},
		},
		{
			// the headers are set (not added) thus same key values
			// are replaced.
			"Key1:Value1A\nKey1:Value1B",
			map[string][]string{
				"Key1": []string{"Value1B"},
			},
		},
		{
			// a key with no value gets removed by design of the package.
			"Key1",
			map[string][]string{},
		},
	}

	for _, set := range testSet {

		tmpHeaders := http.Header{}
		for k, v := range set.Res {
			tmpHeaders[k] = append(tmpHeaders[k], v...)
			sort.Strings(tmpHeaders[k])
		}

		headers, _ := buildHeaders(set.H)
		for _, v := range headers {
			sort.Strings(v)
		}

		// comparison; using the not very efficient reflect.DeepEqual
		// because its a small test suite.
		if !reflect.DeepEqual(tmpHeaders, headers) {
			t.Errorf("Different results")
		}
	}
}
