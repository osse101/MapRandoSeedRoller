package randomize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"reflect"
	"strings"

	"maprandoseedroller/lib/models"
)

// HTTPClient is used for all requests to the MapRando backend. Tests may
// swap it for a client with a fake Transport to avoid real network calls.
var HTTPClient = &http.Client{}

// UpstreamError indicates the failure occurred talking to the MapRando
// backend (transport failure, non-200 response, or an unparsable response
// body) rather than being caused by invalid caller input. api/roll.go uses
// errors.As to map this to 502 instead of 400.
type UpstreamError struct {
	Err error
}

func (e *UpstreamError) Error() string { return e.Err.Error() }
func (e *UpstreamError) Unwrap() error { return e.Err }

func MakeRequest(baseURL string, settings models.RequestMapRando) (models.SeedData, error) {
	body, contentType, err := buildMultipartRequest(settings)
	if err != nil {
		return models.SeedData{}, err
	}

	slog.Info("Sending request", slog.String("endpoint", baseURL+"/randomize"))
	req, err := http.NewRequest("POST", baseURL+"/randomize", body)
	if err != nil {
		return models.SeedData{}, err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := HTTPClient.Do(req)
	if err != nil {
		slog.Error("HTTP request failed", slog.Any("error", err))
		return models.SeedData{}, &UpstreamError{Err: err}
	}
	slog.Info("Response received", slog.String("status", resp.Status))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.SeedData{}, &UpstreamError{Err: fmt.Errorf("unexpected status: %s", resp.Status)}
	}

	var result models.ResponseMapRando
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.SeedData{}, &UpstreamError{Err: fmt.Errorf("failed to decode response JSON: %w", err)}
	}

	seedURL := result.SeedURL
	if strings.HasPrefix(seedURL, "/") {
		seedURL = baseURL + seedURL
	}

	return models.SeedData{
		SeedURL:  seedURL,
		SeedHash: result.SeedHash,
	}, nil
}

func buildMultipartRequest(data interface{}) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	val := reflect.ValueOf(data)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		fieldName := structField.Tag.Get("form")
		fileName := structField.Tag.Get("filename")

		if fileName != "" {
			// It's a file part
			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileName))
			h.Set("Content-Type", structField.Tag.Get("content-type"))

			part, err := writer.CreatePart(h)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create part: %w", err)
			}
			if _, err := part.Write(field.Bytes()); err != nil {
				return nil, "", fmt.Errorf("failed to write part: %w", err)
			}
		} else {
			// It's a regular field
			if err := writer.WriteField(fieldName, fmt.Sprint(field.Interface())); err != nil {
				return nil, "", fmt.Errorf("failed to write field: %w", err)
			}
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close multipart writer: %w", err)
	}
	return body, writer.FormDataContentType(), nil
}
