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

const InertiaUserAgent string = "Inertia/1.0 (Action=\"RollSeed\") Contact=https://inertia.run"

func InertiaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		userAgent := r.Header.Get("User-Agent")
		isInertia := userAgent == InertiaUserAgent
		if isInertia {
			slog.Info("Received Inertia GET request")
			request := models.RequestRaw{
				Action: "roll",
				Source: "inertia",
				Data:   json.RawMessage(`"mentor"`),
			}
			resp, err := workflow.Process(request)
			if err != nil {
				writeJSONResponse(w, http.StatusBadRequest, models.ResponseOut{
					Status:  "error",
					Message: err.Error(),
				})
				return
			}
			//Convert to InertiaResponseOut
			seedJSON, err := json.Marshal(resp.Data)
			if err != nil {
				writeJSONResponse(w, http.StatusInternalServerError, models.ResponseOut{
					Status:  "error",
					Message: "failed to marshal seed data",
				})
				return
			}
			var seedData models.SeedData
			if err := json.Unmarshal(seedJSON, &seedData); err != nil {
				writeJSONResponse(w, http.StatusInternalServerError, models.ResponseOut{
					Status:  "error",
					Message: "failed to unmarshal seed data",
				})
				return
			}
			msg := fmt.Sprintf("Your seed is ready: %s | %s.", seedData.SeedURL, seedData.SeedHash)
			iResp := models.InertiaResponseOut{
				URL:     seedData.SeedURL,
				Hash:    seedData.SeedHash,
				Message: msg,
			}

			writeInertiaJSONResponse(w, http.StatusOK, iResp)
			return
		}
		fmt.Fprintf(w, "MapRando Seed Roller API is running. Please use POST with a preset name.")
		return
	}
	writeJSONResponse(w, http.StatusBadRequest, models.ResponseOut{Message: "Bad User-Agent."})
}

func writeInertiaJSONResponse(w http.ResponseWriter, statusCode int, payload models.InertiaResponseOut) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("Failed to encode JSON response", slog.Any("error", err))
	}
}
