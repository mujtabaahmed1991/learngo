package main

import (
	"flag"
	"fmt"
	"github.com/unrolled/render"
	"net/http"
	"time"
)

type LatencyStats struct {
	Average   float64
	Latencies []float64
}

var cmd bool

func PublishHttpLatency(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, GetAverageHTTPLatencyOverInterval(1, "https://gitlab.com"))
	}
}

func GetAverageHTTPLatencyOverInterval(interval int, url string) *LatencyStats {
	timeout := time.After(time.Minute * time.Duration(interval))
	done := false
	start := time.Now()
	ch := make(chan float64)
	go func() {
		for !done {
			latency, err := GetHTTPLatency(url)
			if err != nil {
				fmt.Println(err)
			}
			if cmd {
				fmt.Printf("Time Elapsed: %f\t\r", time.Since(start).Seconds())
			}
			ch <- latency
		}
	}()
	latencies := make([]float64, 0)
	for !done {
		select {
		case latency := <-ch:
			latencies = append(latencies, latency)
		case <-timeout:
			done = true
		}
	}
	return &LatencyStats{
		func() float64 {
			sum := 0.0
			for _, latency := range latencies {
				sum += latency
			}
			return sum / float64(len(latencies))
		}(),
		latencies,
	}
}

func GetHTTPLatency(url string) (float64, error) {
	startTime := time.Now()
	_, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	return time.Since(startTime).Seconds(), nil
}

func main() {
	cmd = *(flag.Bool("cmd", false, "Flag to check program type HTTP-Server/CMD"))
	flag.Parse()
	if !(cmd) {
		http.HandleFunc("/", PublishHttpLatency(render.New(render.Options{
			IndentJSON: true,
		})))
		http.ListenAndServe(":8080", nil)
	} else {
		fmt.Println(*(GetAverageHTTPLatencyOverInterval(1, "https://gitlab.com")))
	}
}
