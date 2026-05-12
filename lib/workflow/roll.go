package workflow

import (
	"slices"
	"strings"

	"maprandoseedroller/lib"
	"maprandoseedroller/lib/models"
	"maprandoseedroller/lib/parser"
	"maprandoseedroller/lib/randomize"
	"maprandoseedroller/preset"
)

func ExecuteRoll(data string) (models.ResponseOut, error) {
	gameData, isDev, err := PrepareGameData(data)
	if err != nil {
		return models.ResponseOut{}, err
	}

	//Send to MapRando
	seedURL, err := randomize.Randomize(gameData, isDev)
	if err != nil {
		return models.ResponseOut{}, err
	}

	//Determine Discord/Racetime fields
	var resp = models.ResponseOut{
		SeedURL: seedURL,
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
	words := strings.Split(data, " ")
	selectedPreset := "s5"
	flags := ""

	switch len(words) {
	case 0:
	case 1:
		if slices.Contains(validPresets, data) {
			selectedPreset = data
		} else {
			flags = data
		}
	default:
		selectedPreset = words[0]
		flags = words[1]
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
