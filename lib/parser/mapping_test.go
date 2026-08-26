package parser

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/go-openapi/testify/v2/require"

	"maprandoseedroller/lib/models"
)

func TestHydrate(t *testing.T) {
	//Setup
	templateRaw := loadTemplate(t)

	tests := []struct {
		name    string
		tokens  []models.Token
		golden  string
		wantErr bool
	}{
		{
			name:    "Simple Case",
			tokens:  []models.Token{},
			golden:  "simple_case.json",
			wantErr: false,
		},
		{
			name: "Skill Preset",
			tokens: []models.Token{
				{Flag: rune(-1), ID: "Expert", Value: models.False},
			},
			golden:  "skill_preset_case.json",
			wantErr: false,
		},
		{
			name: "Objective Preset",
			tokens: []models.Token{
				{Flag: 'o', ID: "objective_options", Value: models.False},
				{Flag: 'o', ID: "Chozos", Value: models.False},
			},
			golden:  "objective_preset_case.json",
			wantErr: false,
		},
		{
			name: "Skill preset wins over escape multiplier nil-out",
			tokens: []models.Token{
				{Flag: 'x', ID: "escape_timer_multiplier", Value: models.True, RawValue: "1.5"},
				{Flag: rune(-1), ID: "Expert", Value: models.False},
			},
			golden:  "skill_preset_priority_case.json",
			wantErr: false,
		},
		{
			name: "Objective preset wins over objective options nil-out",
			tokens: []models.Token{
				{Flag: 'o', ID: "objective_options", Value: models.False},
				{Flag: 'o', ID: "Kraid", Value: models.True},
				{Flag: 'o', ID: "Metroids", Value: models.False},
			},
			golden:  "objective_preset_priority_case.json",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := Hydrate(freshTemplate(t, templateRaw), tt.tokens)
			require.NoError(t, err)
			want := loadGolden(t, tt.golden)
			require.JSONEq(t, string(want), string(got))
		})
	}
}

func loadTemplate(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile("../../testdata/mapping_test/template.json")
	require.NoError(t, err)
	return b
}

func freshTemplate(t *testing.T, raw []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(raw, &m))
	return m
}

func loadGolden(t *testing.T, goldenName string) []byte {
	t.Helper()
	b, err := os.ReadFile("../../testdata/mapping_test/" + goldenName)
	require.NoError(t, err)
	return b
}
