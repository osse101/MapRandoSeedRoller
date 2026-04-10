package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"maprandoseedroller/lib"
	"maprandoseedroller/preset"
	"net/http"
	"os"
)

type Request struct {
	Preset string `json:"preset"`
	Race   bool   `json:"race"`
	Dev    bool   `json:"dev"`
}

type Response struct {
	SeedURL string `json:"seed_url"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, "MapRando Seed Roller API is running. Please use POST with a preset name.")
		return
	}

	req, err := decode(r)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received request: %+v\n", req)

	p, err := route(req.Preset)
	if err != nil {
		fmt.Printf("Routing failed for preset %s: %v\n", req.Preset, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Selected preset: %T\n", p)

	settings, err := buildSettings(p, req.Race)
	if err != nil {
		http.Error(w, "failed to build settings", http.StatusInternalServerError)
		return
	}
	baseURL := buildSite(req.Dev)
	spoilerToken := buildSpoilerToken()
	fmt.Printf("Using site: %s\n", baseURL)

	seedURL, err := lib.Randomize(baseURL, settings, spoilerToken)
	if err != nil {
		fmt.Printf("Randomization failed: %v\n", err)
		http.Error(w, fmt.Sprintf("randomization failed: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Randomization successful: %s\n", seedURL)

	err = writeResponse(seedURL, w)
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}

func decode(r *http.Request) (*Request, error) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return &req, nil
}

var presets = map[string]preset.Preset{
	"s4": &preset.Season4{},
	"season4": &preset.Season4{},
	"community race season 4": &preset.Season4{},
	"mentor": &preset.Mentor{},
	"mentor tournament": &preset.Mentor{},
	"default": &preset.Default{},
}

func route(name string) (preset.Preset, error) {
	p, ok := presets[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("unknown preset: %s", name)
	}
	return p, nil
}

func buildSettings(p preset.Preset, race bool) ([]byte, error) {
	settings, err := p.Settings()
	if err != nil {
		return nil, err
	}
	if race {
		var s map[string]interface{}
		if err := json.Unmarshal(settings, &s); err != nil {
			return nil, err
		}
		s["race_mode"] = true
		settings, err = json.Marshal(s)
		if err != nil {
			return nil, err
		}
	}
	return settings, nil
}

func buildSpoilerToken() string {
	return os.Getenv("SPOILER_TOKEN")
}

func buildSite(isDev bool) string {
	if isDev {
		return "https://dev.maprando.com"
	}
	return "https://maprando.com"
}

func writeResponse(seedURL string, w http.ResponseWriter) error {
	res := Response{
		SeedURL: seedURL,
	}
	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode(res)
	return nil
}