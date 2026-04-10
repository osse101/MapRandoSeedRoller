package lib

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

func Randomize(baseURL string, settings []byte, spoilerToken string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Settings file
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="settings"; filename="settings.json"`)
	h.Set("Content-Type", "application/json")
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

	// Client that doesn't follow redirects so we can grab the Location header
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		return "", err
	}
	fmt.Printf("Response status: %s\n", resp.Status)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	location := resp.Header.Get("Location")
	fmt.Printf("Location header: %s\n", location)
	if location == "" {
		return "", fmt.Errorf("no location header in response")
	}

	// If it's a relative URL, prepend baseURL
	if location[0] == '/' {
		location = baseURL + location
	}

	return location, nil
}
