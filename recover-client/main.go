package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"strings"

	// "strings"
	"time"
)

func generateDummyData(size uint32) []byte {
	data := make([]byte, size)
	for i := uint32(0); i < size; i++ {
		data[i] = 'A' + byte(i%26)
	}
	return data
}

func generateStateObj(size uint32) []byte {
	mtString := "10730\n6000\n0\n/svc-a/recover\nd9175a23-65b0-4e78-9802-1d29d0a019d6\n" +
		"svc-a-110\nsvc-a-110\nfib\n10.214.96.108\n10729\n"
	mtLen := uint32(len(mtString))

	var buf [4]byte

	binary.LittleEndian.PutUint32(buf[:], mtLen)
	slice := append(buf[:], []byte(mtString)...)

	if mtLen+4 < size {
		dummy := generateDummyData(size - mtLen - 4)
		slice = append(slice, dummy...)
	}

	return slice
}

func main() {
	// svcs := []string{"b", "a", "b"}
	for i := 0; i <= 0; i++ {
		//reqBody := generateDummyData(1536) // 10240, 102400, 1048576
		//reqBody = append(reqBody, generateStateObj(512)...)
		//buf := bytes.NewBuffer(reqBody)
		buf := strings.NewReader("d9175a23-65b0-4e78-9802-1d29d0a019d6")
		req, _ := http.NewRequest("POST", "http://127.0.0.1:9903/failure_single", buf)
		// req, _ := http.NewRequest("POST", fmt.Sprintf("http://127.0.0.1:10728/svc-%s/seq?n=3", svcs[i&1]), buf)
		// fmt.Printf("Req: http://127.0.0.1:10728/svc-%s/fib?n=3\n", svcs[i&1])

		//req.Header.Add("x-ftmesh-mode", "0")
		//req.Header.Add("x-ftmesh-length", "1536")
		//
		//req.Header.Add("x-ftmesh-cluster", "cluster_0")

		// marker := strconv.Itoa(i)
		// req.Header.Add("x-ftmesh-bench-marker", marker)
		client := &http.Client{}
		start := time.Now()
		resp, err := client.Do(req)
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
		duration := time.Since(start)
		fmt.Println(duration.Microseconds())

		time.Sleep(50 * time.Millisecond)
	}

}
