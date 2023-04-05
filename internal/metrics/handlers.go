package metrics

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

func MetricsHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println(DefaultStorage)
	for _, m := range DefaultStorage.Counters {
		fmt.Println(m.Get())
	}
	for _, m := range DefaultStorage.Gauges {
		fmt.Println(m.Get())
	}
	if req.Method != http.MethodPost {
		http.Error(res, "Unexpected method", http.StatusMethodNotAllowed)
	}
	if req.Header.Get("Content-Type") != "text/plain" {
		http.Error(res, "Unexpected content-type should be text/plain", http.StatusUnsupportedMediaType)
	}
	re, err := regexp.Compile("/update/(?P<Type>gauge|counter)/(?P<Name>[a-zA-Z]+)/(?P<Value>[0-9]+)")
	if err != nil {
		http.Error(res, "Internal error", http.StatusInternalServerError)
		return
	}
	matches := re.FindStringSubmatch(req.RequestURI)
	if len(matches)-1 < 3 {
		http.Error(res, "Incorrect URL", http.StatusUnprocessableEntity)
		return
	}
	metricType, metricName, metricValue := matches[1], matches[2], matches[3]
	if err := metricTypeHandle(metricType, metricName, metricValue); err != nil {
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
	}
}

func metricTypeHandle(metricType string, metricName string, metricValue string) error {
	switch metricType {
	case "gauge":
		return gaugeHandle(metricName, metricValue)
	case "counter":
		return counterHandle(metricName, metricValue)
	default:
		return nil
	}
}

func gaugeHandle(metricName string, metricValue string) error {
	var m Metric
	if v, ok := DefaultStorage.Gauges[metricName]; !ok {
		m = new(Gauge)
	} else {
		m = v
	}
	g, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		return errors.New("Unexpected value")
	}
	m.Update(g)
	DefaultStorage.Gauges[metricName] = m
	return nil
}

func counterHandle(metricName string, metricValue string) error {
	var m Metric
	if v, ok := DefaultStorage.Counters[metricName]; !ok {
		m = new(Counter)
	} else {
		m = v
	}
	g, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		return errors.New("Unexpected value")
	}
	m.Update(g)
	DefaultStorage.Counters[metricName] = m
	return nil
}
