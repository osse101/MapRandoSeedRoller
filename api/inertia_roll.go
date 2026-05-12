package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	svix "github.com/svix/svix-webhooks/go"

	"maprandoseedroller/lib/models"
)

func InertiaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("SVIX_INERTIA_SECRET")
	wh, err := svix.NewWebhook(secret)
	if err != nil {
		slog.Error("Svix initialization failed", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Read the RAW body bytes
	// We MUST do this before any JSON decoding to preserve the signature integrity
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Verify the signature using headers and raw payload
	// This checks svix-id, svix-timestamp, and svix-signature automatically
	if err := wh.Verify(payload, r.Header); err != nil {
		slog.Warn("Unauthorized webhook attempt", slog.Any("error", err))
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Now unmarshal the payload
	var req models.RequestIn
	if err := json.Unmarshal(payload, &req); err != nil {
		slog.Error("JSON unmarshal failed", slog.Any("error", err))
		http.Error(w, "Invalid payload format", http.StatusBadRequest)
		return
	}

	// Continue to roller logic
	result, err := roller.ExecuteRoll(req)
	if err != nil {
		slog.Error("Execution failed", slog.Any("error", err))
		http.Error(w, fmt.Sprintf("randomization failed: %v", err), http.StatusInternalServerError)
		return
	}

	err = writeResponse(result.SeedURL, w)
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}
