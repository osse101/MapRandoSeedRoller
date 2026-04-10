package preset

import (
	_ "embed"
)

//go:embed data/season4.json
var season4JSON []byte

//go:embed data/mentor.json
var mentorJSON []byte

//go:embed data/Default.json
var defaultJSON []byte

type Season4 struct{}

func (s *Season4) Settings() ([]byte, error) {
	return season4JSON, nil
}

type Mentor struct{}

func (m *Mentor) Settings() ([]byte, error) {
	return mentorJSON, nil
}

type Default struct{}

func (d *Default) Settings() ([]byte, error) {
	return defaultJSON, nil
}
