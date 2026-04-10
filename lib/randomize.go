package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

func Randomize(baseURL string, settings []byte, spoilerToken string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Settings file
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="settings"; filename="settings.json"`)
	h.Set("Content-Type", "text/plain")
	part, err := writer.CreatePart(h)
	if err != nil {
		return "", err
	}
	part.Write(settings)

	// Spoiler token
	err = writer.WriteField("spoiler_token", spoilerToken)
	if err != nil {
		return "", err
	}

	writer.Close()

	fmt.Printf("Sending request to: %s/randomize\n", baseURL)
	req, err := http.NewRequest("POST", baseURL+"/randomize", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		return "", err
	}
	fmt.Printf("Response status: %s\n", resp.Status)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var result struct {
		SeedURL string `json:"seed_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response JSON: %v", err)
	}

	seedURL := result.SeedURL
	if strings.HasPrefix(seedURL, "/") {
		seedURL = baseURL + seedURL
	}

	return seedURL, nil
}
