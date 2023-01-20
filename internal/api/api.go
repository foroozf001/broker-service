package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const httpPort = "8080"

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of HTTP requests processed, based on path",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP responses",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests",
}, []string{"path"})

// Config describes all application members
type Config struct{}

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}

func Run() {
	app := Config{}

	log.Printf("listening on http port %s\n", httpPort)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", httpPort),
		Handler:           app.Routes(),
		ReadHeaderTimeout: 2 * time.Second,
	}

	err := srv.ListenAndServe()
	log.Println(err)
}
