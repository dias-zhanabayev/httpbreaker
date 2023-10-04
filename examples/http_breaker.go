package main

import (
	"fmt"
	"github.com/dias-zhanabayev/httpbreaker"
	"io"
	"net/http"
)

var cb *httpbreaker.CircuitBreaker

func initialize() {
	var st httpbreaker.Settings
	st.Name = "HTTP GET"
	st.ReadyToTrip = func(counts httpbreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures()) / float64(counts.Requests())
		return counts.Requests() >= 3 && failureRatio >= 0.6
	}
	st.TracerTransport = http.DefaultTransport

	cb = httpbreaker.NewCircuitBreaker(st)
}

func main() {
	initialize()
	client := &http.Client{
		Transport: cb,
	}

	// nolint
	resp, err := client.Get("https://www.google.com/robots.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}(resp.Body)
	fmt.Println("Response:", resp.Status)
}
