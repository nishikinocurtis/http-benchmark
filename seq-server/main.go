package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func Seq(n int64) (int64, error) {

	if n <= 1 {
		return n, nil
	}

	if n > 200 {
		return 0, fmt.Errorf("unsupported seq number %d: too large", n)
	}

	var n1 int64 = 0
	for i := int64(1); i <= n; i++ {
		n1 += i
	}

	return n1, nil
}

func calcSeq(w http.ResponseWriter, r *http.Request) {
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

	res, err := Seq(int64(reqN))
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	_, err = io.WriteString(w, strconv.FormatInt(res, 10)+"+"+strconv.FormatInt(now, 10))
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/svc-b/seq", calcSeq)

	_ = http.ListenAndServe(":20730", nil)
}
