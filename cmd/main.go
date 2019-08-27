package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/capitalone/go-future-context"
)

//BatchSize ...
const BatchSize = 2

func main() {
	// one bad url to  print error
	urls := []string{"http://www.yahoo.com", "http://www.google.com", "http://www.slashdot.com", "http://www.sdfasasdfasdfasfasdff.com", "http://www.youtube.com", "http://www.gmail.com", "http://github.com"}
	// call in batches
	for i := 0; i < len(urls); i += BatchSize {
		batch := urls[i:min(i+BatchSize, len(urls))]
		requests := processBatch(batch)
		processRequests(requests)
	}
}

//processBatch ...
func processBatch(batch []string) []future.Interface {
	requests := make([]future.Interface, 0)
	// make concurrent requests
	for j := range batch {
		count := j
		req := future.New(func() (interface{}, error) {
			return makeRequest(batch[count])
		})
		requests = append(requests, req)
	}
	return requests
}

//processRequests ...
func processRequests(requests []future.Interface) {
	for k := range requests {
		ct := k
		response, timeout, err := requests[ct].GetUntil(2 * time.Second)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if timeout {
			fmt.Println("timed out:")
			fmt.Println(requests[ct])
			continue
		}
		fmt.Println(len(response.(string)))
	}
}

//makeRequest ...
func makeRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// min ...
func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
