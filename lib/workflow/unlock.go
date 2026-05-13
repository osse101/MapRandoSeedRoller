package workflow

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"maprandoseedroller/lib"
)

type UnlockRequest struct {
	BaseURL    string
	SeedString string
}

func ExecuteUnlock(seedURL string) error {
	req, err := parseSeedURL(seedURL)
	if err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}
	slog.Info("Received unlock request", slog.Any("request", req))

	err = sendUnlockRequest(*req)
	if err != nil {
		slog.Error("Unlock failed", slog.Any("error", err))
		return fmt.Errorf("unlock failed: %w", err)
	}

	return nil
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

func sendUnlockRequest(req UnlockRequest) error {
	// Build Path
	path := req.BaseURL + "/unlock"

	data := url.Values{}
	data.Set("spoiler_token", lib.BuildSpoilerToken())

	// Send Request
	httpReq, err := http.NewRequest("POST", path, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		return fmt.Errorf("backend returned status: %s", resp.Status)
	}

	return nil
}
