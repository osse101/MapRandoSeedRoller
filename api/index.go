package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"maprandoseedroller/lib/models"
	"maprandoseedroller/lib/workflow"
	"maprandoseedroller/preset"
)

type Roller interface {
	ExecuteRoll(req models.RequestIn) (models.ResponseOut, error)
}

type defaultRoller struct{}

func (d defaultRoller) ExecuteRoll(req models.RequestIn) (models.ResponseOut, error) {
	return workflow.ExecuteRoll(req)
}

var roller Roller = defaultRoller{}

func RandomizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, "MapRando Seed Roller API is running. Please use POST with a preset name.")
		return
	}

	req, err := decode(r)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Received request: %+v\n", req)

	result, err := roller.ExecuteRoll(*req)
	if err != nil {
		log.Printf("Randomization failed: %v\n", err)
		http.Error(w, fmt.Sprintf("randomization failed: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("Randomization successful: %s\n", result)

	err = writeResponse(result.SeedURL, w)
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}

func decode(r *http.Request) (*models.RequestIn, error) {
	var req models.RequestIn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return &req, nil
}

func writeResponse(seedURL string, w http.ResponseWriter) error {
	res := models.ResponseOut{
		SeedURL: seedURL,
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func GetHelpText(input string) string {
	switch input {
	case "preset", "presets":
		presets := preset.GetPresetNames()
		return "Available presets: " + strings.Join(presets, ", ")
	case "flag", "flags":
		return "These are your flags:"
	}
	return "Usage: !roll <preset> <flags>.  For more help use !help presets or !help flags"
}
