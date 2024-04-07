package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	// "strings"
	"time"

	"github.com/google/uuid"
)

func generateDummyData(size uint32) []byte {
	data := make([]byte, size)
	for i := uint32(0); i < size; i++ {
		data[i] = 'A' + byte(i%26)
	}
	return data
}

func generateStateObj(size uint32) []byte {
	mtString := "10730\n0\n6000\n/svc-a/recover\nd9175a23-65b0-4e78-9802-1d29d0a019d6\n" +
		"svc-a-107\nsvc-a-110\nfib\n10.214.96.110\n10729\n"
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

func readAllBody(closer io.ReadCloser) ([]byte, error) {
	defer func(closer io.ReadCloser) {
		err := closer.Close()
		if err != nil {
			panic(err)
		}
	}(closer)
	body, err := io.ReadAll(closer)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func UpdateRequestGenerator(timeString string, dummy bool) *http.Request {
	bts := []byte("{\"ExpressionAttributeNames\":{\"#0\":\"info\",\"#1\":\"rating\",\"#2\":\"plot\"},\"ExpressionAttributeValues\":{\":0\":{\"S\":\"5.0\"},\":1\":{\"S\":\"ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEF\"}},\"Key\":{\"title\":{\"S\":\"FooBar\"},\"year\":{\"N\":\"2024\"}},\"ReturnValues\":\"UPDATED_NEW\",\"TableName\":\"movieTable\",\"UpdateExpression\":\"SET #0.#1 = :0, #0.#2 = :1\\n\"}\n")
	oriLen := len(bts)

	if dummy {
		bts = append(bts, generateStateObj(512)...)
	}

	buf := bytes.NewBuffer(bts)
	req, _ := http.NewRequest("POST", "http://127.0.0.1:10729/db", buf)

	req.Header.Add("Amz-Sdk-Invocation-Id", uuid.NewString())
	req.Header.Add("Amz-Sdk-Request", "attempt=1; max=3")
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=AaBbCcDc1234/20240407/us_west_2/dynamodb/aws4_request, SignedHeaders=accept-encoding;amz-sdk-invocation-id;amz-sdk-request;content-length;content-type;host;x-amz-date;x-amz-target, Signature=0bf62c65363bbb400cda94abba6e60194d4a9e356f20581fd9e799184a1d43ac")
	req.Header.Set("Content-Type", "application/x-amz-json-1.0")
	req.Header.Add("X-Amz-Date", "20240407T"+timeString+"Z")
	req.Header.Add("X-Amz-Target", "DynamoDB_20120810.UpdateItem")
	req.Header.Set("Accept-Encoding", "identity")

	if dummy {
		req.Header.Add("x-ftmesh-length", fmt.Sprintf("%d", oriLen))
		req.Header.Add("x-ftmesh-mode", "0")
	}

	return req
}

func main() {
	// svcs := []string{"b", "b"}
	// wg := sync.WaitGroup{}

	// dynamodbTable, e := db.GetTableClient("movieTable")
	//if e != nil {
	//	return
	//}

	//_, e = dynamodbTable.CreateMovieTable()
	//if e != nil {
	//	_ = fmt.Errorf("create table error: %w", e)
	//}
	//
	//rating := "5.0"
	//
	//e = dynamodbTable.AddMovie(db.Movie{
	//	Title: "FooBar",
	//	Year:  2024,
	//	Info: map[string]interface{}{
	//		"rating": rating,
	//		"plot":   string(generateDummyData(1436)),
	//	},
	//})
	//
	//m, e := dynamodbTable.UpdateMovie(db.Movie{
	//	Title: "FooBar",
	//	Year:  2024,
	//	Info: map[string]interface{}{
	//		"rating": rating,
	//		"plot":   string(generateDummyData(1436)),
	//	},
	//})
	//
	//if e == nil {
	//	fmt.Printf("%v\n", m)
	//}

	for i := 0; i <= 9; i++ {
		//reqBody := generateDummyData(1536) // 10240, 102400, 1048576
		//// reqBody = append(reqBody, generateStateObj(512)...)
		//buf := bytes.NewBuffer(reqBody)
		//// buf := strings.NewReader("d9175a23-65b0-4e78-9802-1d29d0a019d6")
		//// req, _ := http.NewRequest("POST", "http://127.0.0.1:10729/recover", buf)
		//
		//req, _ := http.NewRequest("POST", fmt.Sprintf("http://127.0.0.1:10729/svc-%s/seq?n=3", svcs[i&1]), buf)
		//fmt.Printf("Req: http://127.0.0.1:10729/svc-%s/fib?n=3\n", svcs[i&1])
		// req, _ := http.NewRequest("POST", "http://127.0.0.1:9903/failure_single", buf)
		//req.Header.Add("x-ftmesh-resource-id", "d9175a23-65b0-4e78-9802-1d29d0a019d6")
		//req.Header.Add("x-ftmesh-states-position", "946176")
		//req.Header.Add("x-ftmesh-recover-port", "10801")
		//req.Header.Add("x-ftmesh-flags", "0")
		//req.Header.Add("x-ftmesh-ttl", "6000")
		//req.Header.Add("x-ftmesh-recover-uri", "/recover")
		//req.Header.Add("x-ftmesh-svc-id", "curl-test")
		//req.Header.Add("x-ftmesh-pod-id", "curl-test")
		//req.Header.Add("x-ftmesh-method-name", "curl-test")
		//req.Header.Add("x-ftmesh-svc-port", "10801")
		//req.Header.Add("x-ftmesh-svc-ip", "127.0.0.1")
		// req.Header.Add("x-ftmesh-mode", "0")
		// req.Header.Add("x-ftmesh-length", "1536")
		// req.Header.Add("Keep-Alive", "timeout=5, max=100")

		req := UpdateRequestGenerator("20590"+fmt.Sprintf("%d", i), false)
		// req.Header.Add("x-ftmesh-mode", "0")
		// req.Header.Add("x-ftmesh-length", "1536")

		req.Header.Add("x-ftmesh-cluster", "cluster_0")

		// marker := strconv.Itoa(i)
		// req.Header.Add("x-ftmesh-bench-marker", marker)

		start := time.Now()
		// tbs, err := dynamodbTable.ListTables()
		//if err != nil {
		//
		//	panic(err)
		//}
		client := &http.Client{}
		fmt.Printf("Before do req ts: %v\n", time.Now().UnixMicro())
		_, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		duration := time.Since(start)
		// fmt.Println(tbs)
		fmt.Println(duration.Microseconds())

		// _, _ = readAllBody(resp.Body)

		time.Sleep(1000 * time.Millisecond)
	}
}
