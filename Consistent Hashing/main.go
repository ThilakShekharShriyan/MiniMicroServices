package main

import (
	"log"
	"net/http"

	"github.com/thilakshekharshriyan/api"
)

func main() {
	log.Println("Starting consistent hashing server on :8080")
	server := api.NewServer(5) // 5 virtual nodes
	log.Fatal(http.ListenAndServe(":8080", server.Routes()))
}