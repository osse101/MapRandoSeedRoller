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
		models.SkillPresetAliases,
		models.ObjectivePresetAliases,
		models.ProgressionPresetAliases,
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

	action, hasAction := PresetActions[strings.ToLower(selectedPreset)]
	isRealPreset := slices.Contains(validPresets, selectedPreset)
	if !isRealPreset && !hasAction {
		return nil, false, nil, fmt.Errorf("invalid preset selected")
	}

	// Parse Flags
	tokens, err := parser.Lex(flags, flagTable)
	if err != nil {
		return nil, false, nil, err
	}

	// Write json preset. Meta/action-only presets (e.g. "drockyrandom") have
	// no template of their own — the action below is responsible for
	// selecting and loading one into tmpl.
	var tmpl map[string]interface{}
	if isRealPreset {
		tmpl, err = preset.LoadTemplate(selectedPreset)
		if err != nil {
			return nil, false, nil, err
		}
	} else {
		tmpl = map[string]interface{}{}
	}

	// Run any preset-specific custom action before hydration, so user flags
	// still apply on top of whatever the action mutates.
	var extra interface{}
	if hasAction {
		extra, err = action(tmpl)
		if err != nil {
			return nil, false, nil, err
		}
		if !isRealPreset && len(tmpl) == 0 {
			return nil, false, nil, fmt.Errorf("preset action %q did not select a template", selectedPreset)
		}
	}

	gameData, isDev, err := parser.Hydrate(tmpl, tokens)
	if err != nil {
		return nil, false, nil, err
	}
	return gameData, isDev, extra, nil
}
