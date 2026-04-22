package randomize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"reflect"
	"strings"

	"maprandoseedroller/lib/models"
)

func RandomizeRequest(baseURL string, settings models.RequestMapRando) (string, error) {
	body, contentType, err := buildMultipartRequest(settings)

	fmt.Printf("Sending request to: %s/randomize\n", baseURL)
	req, err := http.NewRequest("POST", baseURL+"/randomize", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentType)

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

	var result models.ResponseMapRando
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response JSON: %v", err)
	}

	seedURL := result.SeedURL
	if strings.HasPrefix(seedURL, "/") {
		seedURL = baseURL + seedURL
	}

	return seedURL, nil
}


func buildMultipartRequest(data interface{}) (*bytes.Buffer, string, error){
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	val := reflect.ValueOf(data)
	typ := val.Type()

	for i:= 0; i<val.NumField(); i++{
		field:= val.Field(i)
		structField := typ.Field(i)

		fieldName := structField.Tag.Get("form")
		fileName := structField.Tag.Get("filename")

		if fileName != ""{
			// It's a file part
			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileName))
			h.Set("Content-Type", structField.Tag.Get("content-type"))

			part, err := writer.CreatePart(h)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create part: %w", err)
			}
			part.Write(field.Bytes())
		}else{
			// It's a regular field
			writer.WriteField(fieldName, fmt.Sprint(field.Interface()))
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close multipart writer: %w", err)
	}
	return body, writer.FormDataContentType(), nil
}