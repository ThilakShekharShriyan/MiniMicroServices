package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thilakshekharshriyan/api"
	"github.com/thilakshekharshriyan/nodemetrics"
)


func main() {

	
	nodemetrics.InitMetrics()
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

	log.Println("Starting consistent hashing server on :8080")
	server := api.NewServer(5) // 5 virtual nodes
	log.Fatal(http.ListenAndServe(":8080", server.Routes()))


}