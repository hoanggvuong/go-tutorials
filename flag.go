package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

// go-tutorials/flag.go
// create a simple HTTP server with command-line flags get status code 200
// and a timeout of 5 seconds.
func main() {
	url := flag.String("url", "http://www.google.com", "URL to fetch")
	timeout := flag.Duration("timeout", 5*time.Second, "Timeout for the request")
	flag.Parse()
	client := http.Client{
		Timeout: *timeout,
	}

	resp, err := client.Get(*url)
	if err != nil {
		fmt.Println("Lỗi khi gửi yêu cầu:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Yêu cầu thành công:", resp.Status)
	} else {
		fmt.Println("Yêu cầu thất bại:", resp.Status)
	}
}
