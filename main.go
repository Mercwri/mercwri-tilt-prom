package main

import (
	"net/http"
	"time"

	"github.com/alexhowarth/go-tilt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type TiltReading struct {
	Gravity float64
	TempF   float64
	TempC   float64
	Colour  string
}

func TiltMetrics() []TiltReading {
	var readings []TiltReading
	s := tilt.NewScanner()
	s.Scan(20 * time.Second)
	for _, t := range s.Tilts() {
		reading := &TiltReading{
			Gravity: t.Gravity(),
			TempF:   float64(t.Fahrenheit()),
			TempC:   float64(t.Celsius()),
		}
		readings = append(readings, *reading)
	}
	return readings
}

type TiltCollector struct {
	Gravity *prometheus.Desc
	TempF   *prometheus.Desc
	TempC   *prometheus.Desc
}

func NewTiltCollector() *TiltCollector {
	return &TiltCollector{
		Gravity: prometheus.NewDesc("tilt_gravity_r", "latest specfic gravity reading", []string{"colour"}, nil),
		TempF:   prometheus.NewDesc("tilt_temperature_reading_f", "latest temperature reading", []string{"colour"}, nil),
		TempC:   prometheus.NewDesc("tilt_temperature_reading_c", "latest temperature reading", []string{"colour"}, nil),
	}
}

func (collector *TiltCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.Gravity
	ch <- collector.TempC
	ch <- collector.TempF
}

func (collector *TiltCollector) Collect(ch chan<- prometheus.Metric) {
	readings := TiltMetrics()
	for _, r := range readings {
		ch <- prometheus.MustNewConstMetric(collector.Gravity, prometheus.GaugeValue, r.Gravity, r.Colour)
		ch <- prometheus.MustNewConstMetric(collector.Gravity, prometheus.GaugeValue, r.TempF, r.Colour)
		ch <- prometheus.MustNewConstMetric(collector.Gravity, prometheus.GaugeValue, r.TempC, r.Colour)
	}
}

func main() {
	tiltCollect := NewTiltCollector()
	prometheus.MustRegister(tiltCollect)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
