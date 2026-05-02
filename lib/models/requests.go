package models

type RequestIn struct {
	Preset string `json:"preset"`
	Source string `json:"source"`
	Flags  string `json:"flags,omitempty"`
	Title  string `json:"title,omitempty"`
}
type RequestIn2 struct {
	Args       string `json:"args"`
	Source     string `json:"source"`
	SourceInfo string `json:"source_info,omitempty"`
}

type ResponseOut struct {
	SeedURL string `json:"seed_url"`
}
type ResponseOut2 struct {
	SeedURL  string `json:"seed_url"`
	SeedHash string `json:"seed_hash"`
	Info     string `json:"info"`
	Message  string `json:"message"`
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
