package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

	svix "github.com/svix/svix-webhooks/go"

	"maprandoseedroller/lib/httpio"
	"maprandoseedroller/lib/models"
	"maprandoseedroller/lib/workflow"
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
	defer r.Body.Close()
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
	var req models.RequestRaw
	if err := json.Unmarshal(payload, &req); err != nil {
		slog.Error("JSON unmarshal failed", slog.Any("error", err))
		http.Error(w, "Invalid payload format", http.StatusBadRequest)
		return
	}

	// Delegate to manager
	resp, err := workflow.Process(req)
	if err != nil {
		httpio.WriteJSONResponse(w, http.StatusBadRequest, models.ResponseOut{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	httpio.WriteJSONResponse(w, http.StatusOK, resp)
}
