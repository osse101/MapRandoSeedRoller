package models

import "encoding/json"

type RequestIn struct {
	Preset string `json:"preset"`
	Source string `json:"source"`
	Flags  string `json:"flags,omitempty"`
	Title  string `json:"title,omitempty"`
}

type RequestRaw struct {
	Event  string          `json:"event"`  // e.g. rt.raceend
	Action string          `json:"action"` // e.g. unlock
	Source string          `json:"source"`
	Data   json.RawMessage `json:"data"`
}

type ResponseOut struct {
	SeedURL  string `json:"seed_url"`
	SeedHash string `json:"seed_hash"`
}
type ResponseOut2 struct {
	SeedURL  string `json:"seed_url"`
	SeedHash string `json:"seed_hash"`
	Info     string `json:"info"`
	Message  string `json:"message"`
}

type ResponseOut3 struct {
	Status  string      `json:"status"` // "success" or "error"
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type RequestMapRando struct {
	// The "settings" field is actually a file (settings.json)
	Settings     []byte `form:"settings" filename:"settings.json" content-type:"text/plain"`
	SpoilerToken string `form:"spoiler_token"`
}

type ResponseMapRando struct {
	SeedURL  string `json:"seed_url"`
	SeedHash string `json:"seed_hash"`
}
