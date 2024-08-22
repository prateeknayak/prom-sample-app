package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


var (
	counter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "golang",
			Name:      "my_counter",
			Help:      "This is my counter",
		})

	gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "golang",
			Name:      "my_gauge",
			Help:      "This is my gauge",
		})

	histogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "golang",
			Name:      "my_histogram",
			Help:      "This is my histogram",
		})

	summary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "golang",
			Name:      "my_summary",
			Help:      "This is my summary",
		})
)



// var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
//     Name: "myapp_http_duration_seconds",
//     Help: "Duration of HTTP requests.",
//   }, []string{"path"})


// var randomRequestCounter = prometheus.NewCounter(
//    prometheus.CounterOpts{
//        Name: "random_request_counter",
//        Help: "No of request handled by Random Response handler",
//    },
// )


func main() {

	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "prom_request_time",
		Help: "Time it has taken to retrieve the metrics",
	}, []string{"time"})
	
	prometheus.Register(histogramVec)
	prometheus.MustRegister(counter)
	prometheus.MustRegister(gauge)
	prometheus.MustRegister(histogram)
	prometheus.MustRegister(summary)

	go func() {
		for {
			counter.Add(rand.Float64() * 5)
			gauge.Add(rand.Float64()*15 - 5)
			histogram.Observe(rand.Float64() * 10)
			summary.Observe(rand.Float64() * 10)
			time.Sleep(time.Second)
		}
	}()
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleList)
	mux.HandleFunc("/health", healthCheck)
	mux.HandleFunc("/random", randomResponse)
	mux.Handle("/metrics", newHandlerWithHistogram(promhttp.Handler(), histogramVec))
	log.Fatalf("received error while serving. %v", http.ListenAndServe(":8080", mux))
}

func newHandlerWithHistogram(handler http.Handler, histogram *prometheus.HistogramVec) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		status := http.StatusOK

		defer func() {
			histogram.WithLabelValues(fmt.Sprintf("%d", status)).Observe(time.Since(start).Seconds())
		}()

		if req.Method == http.MethodGet {
			handler.ServeHTTP(w, req)
			return
		}
		status = http.StatusBadRequest

		w.WriteHeader(status)
	})
}

func handleList(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	body, err := getList(ctx)
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	_, err = fmt.Fprint(w, string(body))
	if err != nil {
		fmt.Printf("error responding")
	}
}

func getList(ctx context.Context) ([]byte, error) {

	fmt.Println("calling list from httpbin")
	r, err := http.Get("http://httpbin.org/stream/5")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body);
	fmt.Println("finished calling httpbin")
	return body, nil
}

func healthCheck(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprint(w, "up")
	if err != nil {
		fmt.Printf("error responding to health check")
	}

}

func randomResponse(w http.ResponseWriter, req *http.Request) {

	resp := randomiseResponseCode(&http.Response{})
	resp.Body = io.NopCloser(strings.NewReader("random response from the server"))

	w.WriteHeader(resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	_, err := w.Write(body)
	if err != nil {
		fmt.Printf("error while sending respone %v", err)
	}

	_ = resp.Body.Close()

}

func randomiseResponseCode(resp *http.Response) *http.Response {
	codes := []int{200, 203, 500, 501, 502, 503, 504, 505}
	rn := rand.Intn(len(codes))
	i := rn % len(codes)
	resp.StatusCode = codes[i]

	return resp
}
