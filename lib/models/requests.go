package models

type RequestIn struct {
	Preset string `json:"preset"`
	Source string `json:"source"`
	Flags string `json:"flags,omitempty"`
	Title string `json:"title,omitempty"`
}

type ResponseOut struct {
	SeedURL string `json:"seed_url"`
}

type RequestMapRando struct{
	// The "settings" field is actually a file (settings.json)
    Settings     []byte `form:"settings" filename:"settings.json" content-type:"text/plain"`
    SpoilerToken string `form:"spoiler_token"`
}

type ResponseMapRando struct{
	SeedURL string `json:"seed_url"`
}