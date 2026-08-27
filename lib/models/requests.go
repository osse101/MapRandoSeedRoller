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
	URL     string      `json:"url"`
	Hash    string      `json:"hash,omitempty"`
	Message string      `json:"message,omitempty"`
	Extra   interface{} `json:"extra,omitempty"`
}

type SeedData struct {
	SeedURL  string `json:"seed_url"` // Full URL
	SeedHash string `json:"seed_hash"`
}

type RollResponseData struct {
	SeedData
	Extra interface{} `json:"extra,omitempty"`
}

type SpriteExtra struct {
	SpriteName string `json:"sprite_name"`
}

// PresetSelectionExtra describes a meta preset action's random pick — which
// preset/variation was selected, any starting items chosen for it, and (when
// the pick chains into that preset's own action, e.g. randompreset landing
// on nis) that action's own Extra payload nested underneath.
type PresetSelectionExtra struct {
	SelectedPreset string      `json:"selected_preset"`
	StartingItems  []string    `json:"starting_items,omitempty"`
	Extra          interface{} `json:"extra,omitempty"`
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
