package randomize

import (
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"maprandoseedroller/lib/models"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func withFakeTransport(t *testing.T, fn roundTripFunc) {
	t.Helper()
	orig := HTTPClient
	HTTPClient = &http.Client{Transport: fn}
	t.Cleanup(func() { HTTPClient = orig })
}

func TestMakeRequest_Success(t *testing.T) {
	var capturedFields map[string]string
	withFakeTransport(t, func(req *http.Request) (*http.Response, error) {
		_, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
		if err != nil {
			t.Fatalf("failed to parse content type: %v", err)
		}
		mr := multipart.NewReader(req.Body, params["boundary"])
		capturedFields = map[string]string{}
		for {
			part, err := mr.NextPart()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				t.Fatalf("failed to read multipart part: %v", err)
			}
			data, err := io.ReadAll(part)
			if err != nil {
				t.Fatalf("failed to read part data: %v", err)
			}
			capturedFields[part.FormName()] = string(data)
		}

		body := `{"seed_url":"/seed/xyz789","seed_hash":"ONE TWO THREE FOUR"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	result, err := MakeRequest("https://maprando.com", models.RequestMapRando{
		Settings:     []byte(`{"foo":"bar"}`),
		SpoilerToken: "test-token",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if capturedFields["settings"] != `{"foo":"bar"}` {
		t.Errorf("expected settings field %q, got %q", `{"foo":"bar"}`, capturedFields["settings"])
	}
	if capturedFields["spoiler_token"] != "test-token" {
		t.Errorf("expected spoiler_token field %q, got %q", "test-token", capturedFields["spoiler_token"])
	}

	if want := "https://maprando.com/seed/xyz789"; result.SeedURL != want {
		t.Errorf("expected seed URL %q, got %q", want, result.SeedURL)
	}
	if result.SeedHash != "ONE TWO THREE FOUR" {
		t.Errorf("expected seed hash %q, got %q", "ONE TWO THREE FOUR", result.SeedHash)
	}
}

func TestMakeRequest_UpstreamErrors(t *testing.T) {
	tests := []struct {
		name      string
		transport roundTripFunc
	}{
		{
			name: "non-200 status",
			transport: func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader("")),
					Header:     make(http.Header),
				}, nil
			},
		},
		{
			name: "transport failure",
			transport: func(_ *http.Request) (*http.Response, error) {
				return nil, errors.New("connection refused")
			},
		},
		{
			name: "malformed response body",
			transport: func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("not json")),
					Header:     make(http.Header),
				}, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withFakeTransport(t, tt.transport)

			_, err := MakeRequest("https://maprando.com", models.RequestMapRando{
				Settings:     []byte(`{}`),
				SpoilerToken: "test-token",
			})
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			var upstreamErr *UpstreamError
			if !errors.As(err, &upstreamErr) {
				t.Errorf("expected error to be an *UpstreamError, got %T: %v", err, err)
			}
		})
	}
}
