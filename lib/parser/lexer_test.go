package parser

import (
	"reflect"
	"testing"

	"maprandoseedroller/lib"
	"maprandoseedroller/lib/models"
)

func TestLex(t *testing.T) {
	aliasTable := lib.MergeAndSortAliases(
		models.ObjectiveAliases,
		models.ItemAliases,
		models.FlagAliases,
	)
	tests := []struct {
		name    string
		input   string
		wantOut []models.Token
		wantErr bool
	}{
		{
			name:  "Valid input",
			input: "D",
			wantOut: []models.Token{
				{Flag: 'd', ID: "version", Value: models.True},
			},
			wantErr: false,
		},
		{
			name:    "No input",
			input:   "",
			wantOut: []models.Token{},
			wantErr: false,
		},
		{
			name:    "Ignore non-alias characters",
			input:   "758943",
			wantOut: []models.Token{},
			wantErr: false,
		},
		{
			name:  "Array input",
			input: "s:MORPHSPACEchaRgewave",
			wantOut: []models.Token{
				{Flag: 's', ID: "starting_items", Value: models.False},
				{Flag: 's', ID: "Morph", Value: models.True},
				{Flag: 's', ID: "SpaceJump", Value: models.True},
				{Flag: 's', ID: "Charge", Value: models.Maybe},
				{Flag: 's', ID: "Wave", Value: models.False},
			},
			wantErr: false,
		},
		{
			name:  "Multiple flags",
			input: "RDS:hjbgrappleSPEEDVARIA",
			wantOut: []models.Token{
				{Flag: 'r', ID: "race_mode", Value: models.True},
				{Flag: 'd', ID: "version", Value: models.True},
				{Flag: 's', ID: "starting_items", Value: models.True},
				{Flag: 's', ID: "HiJump", Value: models.False},
				{Flag: 's', ID: "Grapple", Value: models.False},
				{Flag: 's', ID: "SpeedBooster", Value: models.True},
				{Flag: 's', ID: "Varia", Value: models.True},
			},
			wantErr: false,
		},
		{
			name:  "Aliases without initial flag",
			input: "hjbgrappleSPEEDVARIA",
			wantOut: []models.Token{
				{Flag: rune(-1), ID: "HiJump", Value: models.False},
				{Flag: rune(-1), ID: "Grapple", Value: models.False},
				{Flag: rune(-1), ID: "SpeedBooster", Value: models.True},
				{Flag: rune(-1), ID: "Varia", Value: models.True},
			},
			wantErr: false,
		},
		{
			name:  "Delineator Flexibility",
			input: "s:morph,pb,missile",
			wantOut: []models.Token{
				{Flag: 's', ID: "starting_items", Value: models.False},
				{Flag: 's', ID: "Morph", Value: models.False},
				{Flag: 's', ID: "PowerBomb", Value: models.False},
				{Flag: 's', ID: "Missile", Value: models.False},
			},
			wantErr: false,
		},
		{
			name:  "Noise Filtering",
			input: "s:!!MORPH??pb!!",
			wantOut: []models.Token{
				{Flag: 's', ID: "starting_items", Value: models.False},
				{Flag: 's', ID: "Morph", Value: models.True},
				{Flag: 's', ID: "PowerBomb", Value: models.False},
			},
			wantErr: false,
		},
		{
			name:  "Objective Batching",
			input: "o:KRAIDphanDRAYridley",
			wantOut: []models.Token{
				{Flag: 'o', ID: "objective_options", Value: models.False},
				{Flag: 'o', ID: "Kraid", Value: models.True},
				{Flag: 'o', ID: "Phantoon", Value: models.False},
				{Flag: 'o', ID: "Draygon", Value: models.True},
				{Flag: 'o', ID: "Ridley", Value: models.False},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Lex(tt.input, aliasTable)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.wantOut) {
				t.Errorf("Lex() got = %v, want %v", got, tt.wantOut)
			}
		})
	}
}

func TestDetermineTristate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantOut models.TriState
	}{
		{
			name:    "Lowercase input",
			input:   "morph",
			wantOut: models.False,
		},
		{
			name:    "Uppercase input",
			input:   "MORPH",
			wantOut: models.True,
		},
		{
			name:    "Mixed case input",
			input:   "MoRpH",
			wantOut: models.Maybe,
		},
		{
			name:    "Empty input",
			input:   "",
			wantOut: models.Maybe,
		},
		{
			name:    "Numeric input",
			input:   "12345",
			wantOut: models.True,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetermineTriState(tt.input)
			if got != tt.wantOut {
				t.Errorf("DetermineTriState() actual output %v does not match expected output %v", got, tt.wantOut)
				return
			}
		})
	}
}
