package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"maprandoseedroller/lib/models"
)

type mockRoller struct{}

func (m mockRoller) ExecuteRoll(req models.RequestIn) (models.ResponseOut, error) {
	return models.ResponseOut{SeedURL: "http://mock-seed-url.com/123"}, nil
}

func TestRandomizeHandler(t *testing.T) {
	// Mock the roller to avoid network calls
	originalRoller := roller
	roller = mockRoller{}
	defer func() { roller = originalRoller }()

	handler := http.HandlerFunc(RandomizeHandler)

	body, _ := json.Marshal(models.RequestIn{
		Preset: "s4",
		Flags:  "",
	})
	req := httptest.NewRequest("POST", "/api/index", strings.NewReader(string(body)))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var res models.ResponseOut
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if res.SeedURL != "http://mock-seed-url.com/123" {
		t.Errorf("Expected mock URL, got %q", res.SeedURL)
	}
}

func TestGetHelpText(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		{"preset", "Available presets:"},
		{"presets", "Available presets:"},
		{"flag", "These are your flags:"},
		{"unknown", "Usage: !roll"},
	}

	for _, tt := range tests {
		got := GetHelpText(tt.input)
		if !strings.Contains(got, tt.contains) {
			t.Errorf("GetHelpText(%q) = %q, want it to contain %q", tt.input, got, tt.contains)
		}
	}
}