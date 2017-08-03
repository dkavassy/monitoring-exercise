package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// TODO: tests

var targetsFlag = flag.String("targets", "localhost:5000,localhost:4000,localhost:3000,localhost:2000", "A comma-separeted list of monitoring targets")
var metricsFlag = flag.String("metrics", "requests,errors", "A comma-separeted list of metrics to collect")
var intervalFlag = flag.Int64("interval", 60, "Scraping interval in seconds")

type Monitoring struct {
	values      map[int64]map[string]int64
	rates       map[string]float64
	logicalTime int64 // Handle overflow
	targets     []string
	metrics     []string
	interval    int64
	mux         sync.Mutex
}

func (m *Monitoring) scrapeTarget(target, metric string, ch chan int64) {
	url := fmt.Sprintf("http://%s/metrics/%s", target, metric)
	resp, err := http.Get(url)
	if err != nil {
		// TODO: handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// TODO: handle error
	}
	value, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		// TODO: handle error
	}
	log.Printf("Scraped: %s %s %d\n", target, metric, value)

	ch <- value
}

// Start starts the scraping loop
func (m *Monitoring) Start() {
	for {
		m.values[m.logicalTime] = make(map[string]int64)
		for _, metric := range m.metrics {
			currentVal := int64(0)
			ch := make(chan int64)
			for _, target := range m.targets {
				go m.scrapeTarget(target, metric, ch)
			}
			for i := 0; i < len(m.targets); i++ {
				currentVal += <-ch
			}
			m.values[m.logicalTime][metric] = currentVal

			if m.logicalTime > 0 {
				m.updateRate(metric, currentVal)
			}

		}
		time.Sleep(time.Duration(m.interval) * time.Second)
		atomic.AddInt64(&m.logicalTime, 1)
	}
}

func (m *Monitoring) updateRate(metric string, currentVal int64) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.rates[metric] = float64(currentVal-m.values[m.logicalTime-1][metric]) / float64(m.interval)
}

// GetRate gets per second rate for specified metrics
func (m *Monitoring) GetRate(metric string) float64 {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.rates[metric]
}

// NewMonitoring creates a new monitoring object from a list of targets, metrics and scraping interval
func NewMonitoring(targets, metrics []string, scrapingInterval int64) *Monitoring {
	m := new(Monitoring)
	m.logicalTime = 0
	m.values = make(map[int64]map[string]int64)
	m.rates = make(map[string]float64)
	m.targets = targets
	m.metrics = metrics
	m.interval = scrapingInterval
	return m
}

func startServer(m *Monitoring) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, metric := range m.metrics {
			fmt.Fprintf(w, "%s/s: %.1f\n", metric, m.GetRate(metric))
		}
	})

	log.Println("Server started on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	flag.Parse()
	m := NewMonitoring(strings.Split(*targetsFlag, ","), strings.Split(*metricsFlag, ","), *intervalFlag)
	go startServer(m)
	m.Start()
}
