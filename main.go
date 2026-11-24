package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stope := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stope()
	go func() {
		log.Println("gatway listening on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down gateway, waiting for in-flight requests to complete...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shoutdown error: %v", err)
	}
	log.Println("gateway stopped")
}
