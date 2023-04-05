package main

import (
	"github.com/Kroch4ka/yandex-go-alert/internal/metrics"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", metrics.MetricsHandler)
	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Panic(err)
	}

}
