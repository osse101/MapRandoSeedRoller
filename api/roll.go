package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	// Initialize the global slog logger definition
	_ "maprandoseedroller/lib/logger"
	"maprandoseedroller/lib/models"
	"maprandoseedroller/lib/workflow"
)

func RandomizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, "MapRando Seed Roller API is running. Please use POST with a preset name.")
		return
	}

	req, err := decode(r)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	slog.Info("Received request", slog.Any("request", req))

	// Delegate to manager
	resp, err := workflow.Process(*req)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, models.ResponseOut3{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	writeJSONResponse(w, http.StatusOK, resp)
}

func decode(r *http.Request) (*models.RequestRaw, error) {
	var req models.RequestRaw
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return &req, nil
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, payload models.ResponseOut3) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("Failed to encode JSON response", slog.Any("error", err))
	}
}
