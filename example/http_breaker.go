package main

import (
	"fmt"
	"github.com/dias-zhanabayev/httpbreaker"
	"net/http"
)

var cb *httpbreaker.CircuitBreaker

func init() {
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
	client := &http.Client{
		Transport: cb,
	}
	resp, err := client.Get("https://www.google.com/robots.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Response:", resp.Status)
}
