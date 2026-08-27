// Package httpio holds the request/response helpers shared by the api/*.go
// handlers. It exists as its own importable package (rather than living as
// unexported functions in one of the handler files) because Vercel's local
// dev Go builder compiles each handler file in isolation and does not pull
// in unexported sibling declarations from other files in package api.
package httpio

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"maprandoseedroller/lib/models"
)

// DecodeRequest reads and JSON-decodes the body of an inbound API request.
func DecodeRequest(r *http.Request) (*models.RequestRaw, error) {
	defer r.Body.Close()
	var req models.RequestRaw
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

// WriteJSONResponse writes payload as the JSON body of the response with the
// given status code.
func WriteJSONResponse(w http.ResponseWriter, statusCode int, payload models.ResponseOut) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("Failed to encode JSON response", slog.Any("error", err))
	}
}
