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

func ExecuteRoll(data string) (models.RollResponseData, error) {
	gameData, isDev, extra, err := PrepareGameData(data)
	if err != nil {
		return models.RollResponseData{}, err
	}

	//Send to MapRando
	resp, err := randomize.Randomize(gameData, isDev)
	if err != nil {
		return models.RollResponseData{}, err
	}

	return models.RollResponseData{SeedData: resp, Extra: extra}, nil
}

func PrepareGameData(data string) ([]byte, bool, interface{}, error) {
	//Get Keywords
	flagTable := lib.MergeAndSortAliases(
		models.ObjectiveAliases,
		models.ItemAliases,
		models.FlagAliases,
		models.LayoutAliases,
	)
	validPresets := preset.GetPresetNames()

	// Separate preset and flags
	selectedPreset := "s5"
	flags := ""
	if len(data) > 0 {
		words := strings.SplitN(data, " ", 2)
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
		return nil, false, nil, fmt.Errorf("invalid preset selected")
	}

	// Parse Flags
	tokens, err := parser.Lex(flags, flagTable)
	if err != nil {
		return nil, false, nil, err
	}

	//Write json preset
	tmpl, err := preset.LoadTemplate(selectedPreset)
	if err != nil {
		return nil, false, nil, err
	}

	// Run any preset-specific custom action before hydration, so user flags
	// still apply on top of whatever the action mutates.
	var extra interface{}
	if action, ok := PresetActions[strings.ToLower(selectedPreset)]; ok {
		extra, err = action(tmpl)
		if err != nil {
			return nil, false, nil, err
		}
	}

	gameData, isDev, err := parser.Hydrate(tmpl, tokens)
	if err != nil {
		return nil, false, nil, err
	}
	return gameData, isDev, extra, nil
}
