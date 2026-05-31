package models

import "encoding/json"

type RequestRaw struct {
	Event  string          `json:"event"`  // e.g. rt.raceend
	Action string          `json:"action"` // e.g. unlock
	Source string          `json:"source"`
	Data   json.RawMessage `json:"data"`
}

type ResponseOut struct {
	Status  string      `json:"status"` // "success" or "error"
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type InertiaResponseOut struct {
	URL     string `json:"url"`
	Hash    string `json:"hash,omitempty"`
	Message string `json:"message,omitempty"`
}

type SeedData struct {
	SeedURL  string `json:"seed_url"` // Full URL
	SeedHash string `json:"seed_hash"`
}

type RequestMapRando struct {
	// The "settings" field is actually a file (settings.json)
	Settings     []byte `form:"settings" filename:"settings.json" content-type:"text/plain"`
	SpoilerToken string `form:"spoiler_token"`
}

type ResponseMapRando struct {
	SeedURL  string `json:"seed_url"` // seed id, not full URL
	SeedHash string `json:"seed_hash"`
}
