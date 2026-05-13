package workflow

import (
	"fmt"
	"slices"
	"strings"

	"maprandoseedroller/lib"
	"maprandoseedroller/lib/models"
	"maprandoseedroller/lib/parser"
	"maprandoseedroller/lib/randomize"
	"maprandoseedroller/preset"
)

func ExecuteRoll(data string) (models.SeedData, error) {
	gameData, isDev, err := PrepareGameData(data)
	if err != nil {
		return models.SeedData{}, err
	}

	//Send to MapRando
	resp, err := randomize.Randomize(gameData, isDev)
	if err != nil {
		return models.SeedData{}, err
	}

	return resp, nil
}

func PrepareGameData(data string) ([]byte, bool, error) {
	//Get Keywords
	flagTable := lib.MergeAndSortAliases(
		models.ObjectiveAliases,
		models.ItemAliases,
		models.FlagAliases,
	)
	validPresets := preset.GetPresetNames()

	// Separate preset and flags
	selectedPreset := "s5"
	flags := ""
	if len(data) > 0{
		words := strings.Split(data, " ")
		switch len(words) {
		case 0:
		case 1:
			selectedPreset = data
		default:
			selectedPreset = words[0]
			flags = words[1]
		}
	}

	if !slices.Contains(validPresets, selectedPreset) {
		return nil, false, fmt.Errorf("invalid preset selected.")
	}

	// Parse Flags
	tokens, err := parser.Lex(flags, flagTable)
	if err != nil {
		return nil, false, err
	}

	//Write json preset
	tmpl, err := preset.LoadTemplate(selectedPreset)
	if err != nil {
		return nil, false, err
	}
	return parser.Hydrate(tmpl, tokens)
}
