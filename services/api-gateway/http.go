package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody tripPreviewRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	tripPreview, err := tripService.Client.PreviewTrip(r.Context(), reqBody.toProto())
	if err != nil {
		log.Printf("Failed to preview trip: %v", err)
		http.Error(w, "Failed to preview trip", http.StatusInternalServerError)
		return
	}

	// validation
	if reqBody.UserID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	// Call the Trip service

	resp := contracts.APIResponse{
		Data: tripPreview,
	}

	WriteJSON(w, http.StatusCreated, resp)
}
