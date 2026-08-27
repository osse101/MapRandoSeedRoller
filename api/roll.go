package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"maprandoseedroller/lib/httpio"
	// Initialize the global slog logger definition
	_ "maprandoseedroller/lib/logger"
	"maprandoseedroller/lib/models"
	"maprandoseedroller/lib/randomize"
	"maprandoseedroller/lib/workflow"
)

func RandomizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, "MapRando Seed Roller API is running. Please use POST with a preset name.")
		return
	}

	req, err := httpio.DecodeRequest(r)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	slog.Info("Received request", slog.Any("request", req))

	// Delegate to manager
	resp, err := workflow.Process(*req)
	if err != nil {
		status := http.StatusBadRequest
		var upstreamErr *randomize.UpstreamError
		if errors.As(err, &upstreamErr) {
			status = http.StatusBadGateway
		}
		httpio.WriteJSONResponse(w, status, models.ResponseOut{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	httpio.WriteJSONResponse(w, http.StatusOK, resp)
}
