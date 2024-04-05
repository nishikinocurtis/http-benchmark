package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func Fibonacci(n int64) (int64, error) {

	if n <= 1 {
		return n, nil
	}

	if n > 93 {
		return 0, fmt.Errorf("unsupported fibonacci number %d: too large", n)
	}

	var n2, n1 int64 = 0, 1
	for i := int64(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n2 + n1, nil
}

func calcFib(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Accepting new request")
	now := time.Now().UnixMicro()
	reqN, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}
	}

	res, err := Fibonacci(int64(reqN))
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	_, err = io.WriteString(w, strconv.FormatInt(res, 10)+"+"+strconv.FormatInt(now, 10))
	if err != nil {
		return
	}
}

func statesAcceptor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Accepting new request")

	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}

	}

	//xFtmeshBenchMarker := r.Header.Get("x-ftmesh-bench-marker")
	rid := r.Header.Get("x-ftmesh-resource-id")
	w.Header().Set("x-ftmesh-resource-id", rid)

	defer r.Body.Close()
	_, _ = io.ReadAll(r.Body)

	// fmt.Println(string(body))
	_, err := io.WriteString(w, "10.214.96.108:10730")
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/svc-a/fib", calcFib)
	http.HandleFunc("/svc-a/recover", statesAcceptor)

	_ = http.ListenAndServe(":10730", nil)
}
