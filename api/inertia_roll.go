package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"maprandoseedroller/lib/httpio"
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
				Data:   json.RawMessage(`"ammo-balance"`),
			}
			resp, err := workflow.Process(request)
			if err != nil {
				httpio.WriteJSONResponse(w, http.StatusBadRequest, models.ResponseOut{
					Status:  "error",
					Message: err.Error(),
				})
				return
			}

			rollData, ok := resp.Data.(models.RollResponseData)
			if !ok {
				httpio.WriteJSONResponse(w, http.StatusInternalServerError, models.ResponseOut{
					Status:  "error",
					Message: "unexpected response shape from roll action",
				})
				return
			}

			msg := fmt.Sprintf("Your seed is ready: %s | %s.", rollData.SeedURL, rollData.SeedHash)
			iResp := models.InertiaResponseOut{
				URL:     rollData.SeedURL,
				Hash:    rollData.SeedHash,
				Message: msg,
				Extra:   rollData.Extra,
			}

			writeInertiaJSONResponse(w, http.StatusOK, iResp)
			return
		}
		fmt.Fprintf(w, "MapRando Seed Roller API is running. Please use POST with a preset name.")
		return
	}
	httpio.WriteJSONResponse(w, http.StatusBadRequest, models.ResponseOut{Status: "error", Message: "Bad User-Agent."})
}

func writeInertiaJSONResponse(w http.ResponseWriter, statusCode int, payload models.InertiaResponseOut) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("Failed to encode JSON response", slog.Any("error", err))
	}
}
