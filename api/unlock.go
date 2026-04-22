package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"maprandoseedroller/lib"
	// Initialize the global slog logger definition
	_ "maprandoseedroller/lib/logger"
)

type UnlockURL struct {
	SeedURL string `json:"seed_url"`
}

type UnlockRequest struct {
	BaseURL    string
	SeedString string
}

type UnlockResponse struct {
	UnlockMessage string `json:"unlock_message"`
}

func UnlockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, "MapRando Seed Roller API is running. Please use POST with a seed_url.")
		return
	}

	req, err := decodeAndParseSeedURL(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}
	slog.Info("Received unlock request", slog.Any("request", req))

	msg, err := sendUnlockRequest(*req)
	if err != nil {
		slog.Error("Unlock failed", slog.Any("error", err))
		http.Error(w, fmt.Sprintf("unlock failed: %v", err), http.StatusInternalServerError)
		return
	}

	err = writeUnlockResponse(msg, w)
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}

func decodeAndParseSeedURL(r *http.Request) (*UnlockRequest, error) {
	var unlockURL UnlockURL
	if err := json.NewDecoder(r.Body).Decode(&unlockURL); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if unlockURL.SeedURL == "" {
		return nil, fmt.Errorf("missing seed_url")
	}

	var req UnlockRequest
	req.BaseURL = strings.TrimSuffix(unlockURL.SeedURL, "/")

	// Extract seed ID from the end of the URL
	parts := strings.Split(req.BaseURL, "/")
	req.SeedString = parts[len(parts)-1]

	return &req, nil
}

func sendUnlockRequest(req UnlockRequest) (string, error) {
	// Build Path
	path := req.BaseURL + "/unlock"

	data := url.Values{}
	data.Set("spoiler_token", lib.BuildSpoilerToken())

	// Send Request
	httpReq, err := http.NewRequest("POST", path, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("backend returned status: %s", resp.Status)
	}

	return "Seed unlocked.", nil
}

func writeUnlockResponse(UnlockMessage string, w http.ResponseWriter) error {
	res := UnlockResponse{
		UnlockMessage: UnlockMessage,
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}
