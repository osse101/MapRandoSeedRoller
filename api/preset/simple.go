package preset

import (
	_ "embed"
)

//go:embed data/min_preset.json
var raceJSON []byte
//go:embed data/min_preset.json
var mentorJSON []byte


type Simple struct{
	Data []byte // The embedded JSON
}

func (s *Simple) Transform(input string)([]byte, error){
	return mentorJSON, nil
}