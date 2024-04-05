package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func postConfig() {
	config := "svc-a-110\n10.214.96.108\n10728\n127.0.0.1\n10730"
	resp, err := http.Post("http://127.0.0.1:9903/rr_endpoint",
		"text/plain", bytes.NewBufferString(config))
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(body)
	return
}

func makeReqAndCount() (string, bool) {
	// fmt.Println("making req")
	resp, err := http.Get("http://127.0.0.1:10729/svc-a/fib?n=3")
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	bString := string(body)
	fmt.Println(bString)
	_, cut, found := strings.Cut(bString, "+")

	return cut, found
	// fmt.Println(body)
}

func worker(jobs chan struct{}, ret chan struct{}, finished chan struct{}, start time.Time) {
	for {
		select {
		case <-jobs:
			cut, found := makeReqAndCount()
			if found {
				if strings.Compare(cut, "") != 0 {
					i, err := strconv.ParseInt(cut, 10, 64)
					if err != nil {
						panic(err)
					}
					tm := time.UnixMicro(i)
					elapse := tm.Sub(start).Microseconds()
					fmt.Printf("%d\n", elapse)
				}
				finished <- struct{}{}
				return
			}
		case <-ret:
			finished <- struct{}{}
			return
		}
	}
}

func main() {
	startTime := time.Now()
	go postConfig()

	timeCh := make(chan struct{})
	limit := make(chan struct{}, 128)
	finished := make(chan struct{}, 1)

	for w := 1; w <= 128; w++ {
		go worker(limit, timeCh, finished, startTime)
	}

	go func() {
		for {
			select {
			case limit <- struct{}{}:

				// time.Sleep(5 * time.Millisecond)
			case <-timeCh:
				return
			}
			// time.Sleep(10 * time.Millisecond)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(128)

	go func() {
		<-finished
		close(timeCh)
		return
	}()

	for i := 1; i <= 128; i++ {
		go func() {
			<-finished
			wg.Done()
		}()
	}

	wg.Wait()
	close(limit)

	fmt.Println("all finished")
}
