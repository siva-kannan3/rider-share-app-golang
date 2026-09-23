package main

import (
	"encoding/json"
	"net/http"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody tripPreviewRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// validation
	if reqBody.UserID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	// Call the Trip service

	resp := contracts.APIResponse{}

	WriteJSON(w, http.StatusCreated, resp)
}
