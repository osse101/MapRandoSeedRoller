package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"maprandoseedroller/lib/models"
)

// Example ResponseOut
// var mockResult = models.ResponseOut{
// 	Status:  "success",
// 	Message: "Your seed: https://maprando.com/seed/tc2pHBSZc/",
// 	Data: map[string]string{
// 		"seedURL":  "https://maprando.com/seed/tc2pHBSZc/",
// 		"seedHash": "YARD YARD YARD YARD",
// 	},
// }
func TestRandomizeHandler(t *testing.T) {
	handler := http.HandlerFunc(RandomizeHandler)

	body, _ := json.Marshal(models.RequestRaw{
		Action: "test",
		Source: "",
	})
	req := httptest.NewRequest("POST", "/api/roll", strings.NewReader(string(body)))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var res models.ResponseOut
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	resData, ok := res.Data.(models.SeedData)
	if !ok {
		t.Fatalf("expected result to be ResponseOut, got %T", res.Data)
	}

	if resData.SeedURL != "http://mock-seed-url.com/123" {
		t.Errorf("Expected mock URL, got %q", resData.SeedURL)
	}
}
