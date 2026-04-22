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

func MakeRequest(baseURL string, settings models.RequestMapRando) (string, error) {
	body, contentType, err := buildMultipartRequest(settings)
	if err != nil {
		return "", err
	}

	slog.Info("Sending request", slog.String("endpoint", baseURL+"/randomize"))
	req, err := http.NewRequest("POST", baseURL+"/randomize", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("HTTP request failed", slog.Any("error", err))
		return "", err
	}
	slog.Info("Response received", slog.String("status", resp.Status))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var result models.ResponseMapRando
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response JSON: %w", err)
	}

	seedURL := result.SeedURL
	if strings.HasPrefix(seedURL, "/") {
		seedURL = baseURL + seedURL
	}

	return seedURL, nil
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
