package workflow

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"maprandoseedroller/lib"
	"maprandoseedroller/lib/models"
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

func ExecuteUnlock(seedURL string) (models.ResponseOut, error) {
	req, err := parseSeedURL(seedURL)
	if err != nil {
		return models.ResponseOut{}, fmt.Errorf("invalid request: %w", err)
	}
	slog.Info("Received unlock request", slog.Any("request", req))

	_, err = sendUnlockRequest(*req)
	if err != nil {
		slog.Error("Unlock failed", slog.Any("error", err))
		return models.ResponseOut{}, fmt.Errorf("unlock failed: %w", err)
	}

	return models.ResponseOut{
		SeedURL: seedURL,
	}, nil
}

func parseSeedURL(seedURL string) (*UnlockRequest, error) {
	if seedURL == "" {
		return nil, fmt.Errorf("missing seed_url")
	}

	var req UnlockRequest
	req.BaseURL = strings.TrimSuffix(seedURL, "/")

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
