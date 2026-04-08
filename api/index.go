package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"main.go/api/preset"
)

func Handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")

	// Parse
	body, err := io.ReadAll(r.Body)
	if err != nil{
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	input := strings.TrimSpace(string(body))

	// Route
	preset, err := route(input);
	if err != nil{
		http.Error(w, "Error in preset routing", http.StatusBadRequest)
		return
	}
	preset.Transform("")
	if err != nil{
		http.Error(w, "Error building preset settings", http.StatusInternalServerError)
		return
	}

	// Wait
	
	seedLink := []byte("/seed")

	// Respond
	w.WriteHeader(http.StatusOK)
	w.Write(seedLink)
	return 
}

func route(name string)(preset.Preset, error){
	switch strings.ToLower(name){
	case "s4":
	case "season4":
		return &preset.Race{}, nil
	case "mentor":
		return &preset.Simple{Data:nil}, nil
	}
	return nil, fmt.Errorf("unkown preset: %s", name)
}