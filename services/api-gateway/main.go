package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /trip/preview", handleTripPreview)

	server := &http.Server{
		Addr:         httpAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go (func() {
		log.Printf("Server is listening on PORT: %s", httpAddr)
		if err := server.ListenAndServe(); err != nil {
			serverErrors <- err
		}
	})()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Go server failed to start: %v", err)
	case signal := <-shutdownSignal:
		log.Fatalf("Shutdown signal received: %s", signal)
	}

	shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		server.Close()
		log.Printf("Graceful shutdown failed: %v", err)
	}
}
