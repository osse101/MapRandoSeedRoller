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
