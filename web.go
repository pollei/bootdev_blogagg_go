package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func getBytesfromUrl(ctx context.Context, url string) ([]byte, error) {
	var retBytes []byte
	// https://pkg.go.dev/net/http#NewRequestWithContext
	//req, err := http.NewRequest("GET", url, nil)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("bad get\n")
		return retBytes, err
	}
	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("bad do\n")
		return retBytes, err
	}
	defer res.Body.Close()

	if res.StatusCode > 270 {
		//fmt.Printf("status %d\n", res.StatusCode)
		//return retBytes, errors.New("bad status")
		return retBytes, fmt.Errorf("bad status: %d", res.StatusCode)
	}
	data, err := io.ReadAll(res.Body)
	if nil != err {
		fmt.Printf("bad readall\n")
		return retBytes, err
	}

	// fmt.Println(string(data))
	return data, nil
}
