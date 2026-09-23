package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8083")
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /preview", func(w http.ResponseWriter, r *http.Request) {
		log.Println("received call from api gateway service")

		response := struct {
			Message string
		}{
			Message: "trip preview created",
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseBytes)
	})

	log.Println("Starting API Gateway")

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println("Trip service failed to start")
	}

}
