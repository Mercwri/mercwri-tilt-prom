package main

import (
	"net/http"
	"time"

	"github.com/alexhowarth/go-tilt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetics() {
	go func() {
		for {
			s := tilt.NewScanner()
			s.Scan(20 * time.Second)
			for _, t := range s.Tilts() {
				beerReading.WithLabelValues(string(t.Colour())).Inc()
				beerGravity.WithLabelValues(string(t.Colour())).Set(t.Gravity())
				beerTemperature.WithLabelValues(string(t.Colour())).Set(float64(t.Fahrenheit()))
			}
			time.Sleep(40 * time.Second)
		}
	}()
}

var (
	beerReading = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tilt_readings_taken_total",
		Help: "total number of beer readings taken",
	},
		[]string{"colour"},
	)
)

var (
	beerGravity = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tilt_gravity_reading",
		Help: "latest specfic gravity reading",
	},
		[]string{"colour"},
	)
)

var (
	beerTemperature = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tilt_temperature_reading",
		Help: "latest temperature reading",
	},
		[]string{"colour"},
	)
)

func main() {
	recordMetics()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
