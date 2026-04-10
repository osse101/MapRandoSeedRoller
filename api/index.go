package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"maprandoseedroller/api/preset"
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
		fmt.Fprintf(w, "MapRando Seed Roller API is running. Please use POST to / with a preset name.")
		return
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received request: %+v\n", req)
	defer r.Body.Close()

	p, err := route(req.Preset)
	if err != nil {
		fmt.Printf("Routing failed for preset %s: %v\n", req.Preset, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Selected preset: %T\n", p)

	settings, err := p.Settings()
	if err != nil {
		http.Error(w, "failed to load settings", http.StatusInternalServerError)
		return
	}

	// Modify race mode if requested
	if req.Race {
		var s map[string]interface{}
		if err := json.Unmarshal(settings, &s); err != nil {
			http.Error(w, "failed to parse settings", http.StatusInternalServerError)
			return
		}
		s["race_mode"] = true
		settings, err = json.Marshal(s)
		if err != nil {
			http.Error(w, "failed to encode modified settings", http.StatusInternalServerError)
			return
		}
		fmt.Println("Race mode enabled in settings")
	}

	site := "main"
	if req.Dev {
		site = "dev"
	}

	spoilerToken := os.Getenv("SPOILER_TOKEN")
	if spoilerToken == "" {
		spoilerToken = "default_token_please_set_env_var"
	}

	fmt.Printf("Using site: %s\n", site)
	//seedURL, err := Randomize(site, settings, spoilerToken)
	seedURL := "https://maprando.com/seed/1234567890"
	if err != nil {
		fmt.Printf("Randomization failed: %v\n", err)
		http.Error(w, fmt.Sprintf("randomization failed: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Randomization successful: %s\n", seedURL)

	res := Response{
		SeedURL: seedURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func route(name string) (preset.Preset, error) {
	switch strings.ToLower(name) {
	case "s4", "season4", "community race season 4":
		return &preset.Season4{}, nil
	case "mentor", "mentor tournament":
		return &preset.Mentor{}, nil
	case "default":
		return &preset.Default{}, nil
	}
	return nil, fmt.Errorf("unknown preset: %s", name)
}