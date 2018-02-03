package main

import (
	"encoding/json"
	"fmt"
	"github.com/unrolled/render"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

func BenchmarkHTTPLatency(b *testing.B) {
	GetAverageHTTPLatencyOverInterval(1, "https://gitlab.com")
}

func TestHTTPLatency(t *testing.T) {
	latency, err := GetHTTPLatency("https://gitlab.com")
	if err != nil {
		t.Errorf("Error getting stats for single connection: %v\n", err)
	}
	if latency > 4.0 {
		t.Errorf("Too much latency please check your connection: %d\n", latency)
	}
}

func TestHTTPLatencyOverInterval(t *testing.T) {
	client := &http.Client{}
	server := httptest.NewServer(http.HandlerFunc(PublishHttpLatency(formatter)))
	defer server.Close()
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Errorf("Error in creating GET request: %v\n", err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in GET to PublishHttpLatency: %v\n", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected code: 200 recieved code: %d\n", res.StatusCode)
	}

	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v\n", err)
	}

	var stats LatencyStats
	err = json.Unmarshal(payload, &stats)
	if err != nil {
		t.Errorf("Error in Decoding JSON: %v\n", err)
	}

	if average := func() float64 {
		sum := 0.0
		for _, latency := range stats.Latencies {
			sum += latency
		}
		return sum / float64(len(stats.Latencies))
	}(); average != stats.Average {
		t.Errorf("Stats does not match: %f with recieved stats: %f\n", average, stats.Average)
	}

	fmt.Println("Test passed successfully")
}
