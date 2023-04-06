package metrics

import (
	"errors"
	"github.com/Kroch4ka/yandex-go-alert/internal/util"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrNoMetricName = util.StatusErr{
		Status:  http.StatusNotFound,
		Message: "No metric name",
	}
	ErrInvalidValue = util.StatusErr{
		Status:  http.StatusBadRequest,
		Message: "Unexpected metric value",
	}
	ErrUnknownMetricType = util.StatusErr{
		Status:  http.StatusNotImplemented,
		Message: "Incorrect metric type",
	}
)

func UpdateHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "unexpected method", http.StatusMethodNotAllowed)
	}
	URLParts := make([]string, 3)
	k := 0
	for _, el := range strings.Split(req.RequestURI, "/") {
		if el != "" && el != "update" {
			URLParts[k] = el
			k += 1
		}
	}
	metricType, metricName, metricValue := URLParts[0], URLParts[1], URLParts[2]
	var statErr util.StatusErr
	if err := metricTypeHandle(metricType, metricName, metricValue); err != nil && errors.As(err, &statErr) {
		http.Error(res, statErr.Message, statErr.Status)
		return
	}
}

func metricTypeHandle(metricType string, metricName string, metricValue string) error {
	switch metricType {
	case "gauge":
		return gaugeHandle(metricName, metricValue)
	case "counter":
		return counterHandle(metricName, metricValue)
	default:
		return ErrUnknownMetricType
	}
}

func gaugeHandle(metricName string, metricValue string) error {
	if metricName == "" {
		return ErrNoMetricName
	}
	var m Metric
	if v, ok := DefaultStorage.Gauges[metricName]; !ok {
		m = new(Gauge)
	} else {
		m = v
	}
	g, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		return ErrInvalidValue
	}
	m.Update(g)
	DefaultStorage.Gauges[metricName] = m
	return nil
}

func counterHandle(metricName string, metricValue string) error {
	if metricName == "" {
		return ErrNoMetricName
	}
	var m Metric
	if v, ok := DefaultStorage.Counters[metricName]; !ok {
		m = new(Counter)
	} else {
		m = v
	}
	g, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		return ErrInvalidValue
	}
	m.Update(g)
	DefaultStorage.Counters[metricName] = m
	return nil
}
