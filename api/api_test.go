package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"maprandoseedroller/lib/models"
	"maprandoseedroller/lib/randomize"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// withFakeMapRando swaps randomize.HTTPClient for a client whose Transport is
// intercepted by fn, so tests never make a real network call. It restores the
// original client via t.Cleanup and returns a counter of how many requests
// were served by the fake transport.
func withFakeMapRando(t *testing.T, fn func(*http.Request) (*http.Response, error)) *int {
	t.Helper()
	called := 0
	orig := randomize.HTTPClient
	randomize.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			called++
			return fn(req)
		}),
	}
	t.Cleanup(func() { randomize.HTTPClient = orig })
	return &called
}

type rollResponseWire struct {
	Status string                  `json:"status"`
	Data   models.RollResponseData `json:"data,omitempty"`
}

func TestRandomizeHandler(t *testing.T) {
	called := withFakeMapRando(t, func(_ *http.Request) (*http.Response, error) {
		body := `{"seed_url":"/seed/abc123","seed_hash":"HASH HASH"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	handler := http.HandlerFunc(RandomizeHandler)

	body, _ := json.Marshal(models.RequestRaw{
		Action: "roll",
		Source: "test",
		Data:   json.RawMessage(`"s5"`),
	})
	req := httptest.NewRequest("POST", "/api/roll", strings.NewReader(string(body)))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rec.Code)
	}

	var res rollResponseWire
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if res.Status != "success" {
		t.Errorf("Expected status %q, got %q", "success", res.Status)
	}
	if want := "https://maprando.com/seed/abc123"; res.Data.SeedURL != want {
		t.Errorf("Expected seed URL %q, got %q", want, res.Data.SeedURL)
	}
	if res.Data.SeedHash != "HASH HASH" {
		t.Errorf("Expected seed hash %q, got %q", "HASH HASH", res.Data.SeedHash)
	}
	if *called != 1 {
		t.Errorf("Expected exactly 1 request to the fake MapRando backend, got %d", *called)
	}
}

func TestRandomizeHandler_UpstreamFailure(t *testing.T) {
	withFakeMapRando(t, func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
		}, nil
	})

	handler := http.HandlerFunc(RandomizeHandler)

	body, _ := json.Marshal(models.RequestRaw{
		Action: "roll",
		Source: "test",
		Data:   json.RawMessage(`"s5"`),
	})
	req := httptest.NewRequest("POST", "/api/roll", strings.NewReader(string(body)))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("Expected status 502, got %d", rec.Code)
	}

	var res models.ResponseOut
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if res.Status != "error" {
		t.Errorf("Expected status %q, got %q", "error", res.Status)
	}
}
