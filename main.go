package main

import (
	"log"
	"net/http"

	"simple-gateway/config"
	"simple-gateway/gateway"
)

func main() {
	routes, err := config.Load("routes.json")
	if err != nil {
		log.Fatal(err)
	}

	handler, err := gateway.New(routes)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
