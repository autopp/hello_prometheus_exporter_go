package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "SampleMetric"
)

type myCollector struct {
	counter int
}

var (
	exampleCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "example_count",
		Help:      "example counter help",
	})
	exampleGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "example_gauge",
		Help:        "example gauge help",
		ConstLabels: prometheus.Labels{"answer": "42"},
	})
	exampleLabelCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "example_label_count",
		Help:      "example counter help",
	}, []string{"message", "status"})
)

func (c *myCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- exampleCount.Desc()
	ch <- exampleGauge.Desc()
}

func (c *myCollector) Collect(ch chan<- prometheus.Metric) {
	exampleValue := float64(12345)

	ch <- prometheus.MustNewConstMetric(
		exampleCount.Desc(),
		prometheus.CounterValue,
		float64(c.counter),
	)
	ch <- prometheus.MustNewConstMetric(
		exampleGauge.Desc(),
		prometheus.GaugeValue,
		float64(exampleValue),
	)

	ch <- prometheus.MustNewConstMetric(
		exampleLabelCount.WithLabelValues("goodbye", "200").Desc(),
		prometheus.CounterValue,
		float64(c.counter),
		"hello", "200")
	c.counter++
}

func main() {
	var c myCollector
	prometheus.MustRegister(&c)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("127.0.0.1:5000", nil))
}
