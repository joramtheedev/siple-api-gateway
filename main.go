package main

import (
	"log"
	"net/http"
	"os"

	"simple-gateway/config"
	"simple-gateway/gateway"
)

func main() {
	apiKey := os.Getenv("GATEWAY_API_KEY")
	if apiKey == "" {
		log.Fatal("GATEWAY_API_KEY is not set")
	}

	routes, err := config.Load("routes.json")
	if err != nil {
		log.Fatal(err)
	}

	handler, err := gateway.New(routes, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
