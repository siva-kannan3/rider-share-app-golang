package main

import (
	"bytes"
	"encoding/json"
	"log"
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

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Call the Trip service
	response, err := http.Post("http://trip-service:8083/preview", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		http.Error(w, "Downstream trip service failed", http.StatusBadGateway)
		return
	}

	defer response.Body.Close()

	var tripServiceResponse any
	resp := contracts.APIResponse{}
	if err := json.NewDecoder(response.Body).Decode(&resp.Data); err != nil {
		http.Error(w, "Failed to parse trip service response", http.StatusBadGateway)
		return
	}

	log.Println(tripServiceResponse)

	WriteJSON(w, http.StatusCreated, resp)
}
