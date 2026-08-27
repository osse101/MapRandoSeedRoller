package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"maprandoseedroller/lib/models"
	"maprandoseedroller/preset"
)

var skillPresetValues = valueSet(models.SkillPresetAliases)
var objectivePresetValues = valueSet(models.ObjectivePresetAliases)
var progressionPresetValues = valueSet(models.ProgressionPresetAliases)

// valueSet builds a membership set from an alias table's long-form values,
// so a matched token ID can be tested without a linear scan.
func valueSet(m map[string]string) map[string]bool {
	s := make(map[string]bool, len(m))
	for _, v := range m {
		s[v] = true
	}
	return s
}

// Hydrate converts tokens into preset overrides and applies them onto the template.
// Returns the merged preset as JSON bytes ready for the MapRando API.
func Hydrate(template map[string]interface{}, tokens []models.Token) ([]byte, bool, error) {
	// --- Step 1: tokens → PresetFields ---
	fields, err := tokensToPresetFields(tokens)
	if err != nil {
		return nil, false, fmt.Errorf("tokensToPresetFields: %w", err)
	}

	// --- Step 2: apply PresetFields onto template ---
	if len(tokens) > 0 {
		SetNestedValue(template, "name", "Custom")
	}
	if err := applyPresetFields(template, fields); err != nil {
		return nil, false, fmt.Errorf("applyPresetFields: %w", err)
	}

	// --- Step 3: postprocess ---
	postprocess(template, fields)
	isDev := fields.Version == preset.DevVersion
	data, err := json.Marshal(template)

	return data, isDev, err
}

// tokensToPresetFields maps a flat token list onto the PresetFields struct.
// Routing is based on tok.Flag for multi-value flags (o, s, l) and tok.ID
// for single-value flags (r, d, x).
func tokensToPresetFields(tokens []models.Token) (models.PresetFields, error) {
	var f models.PresetFields
	for _, tok := range tokens {
		// Skip flag-indicator tokens (their ID is the flag's long name)
		switch tok.ID {
		case "race_mode":
			f.IsRace = tok.Value != models.False
			continue
		case "version":
			if tok.Value != models.False {
				f.Version = preset.DevVersion
			}
			continue
		case "escape_timer_multiplier":
			if tok.RawValue == "" {
				continue
			}
			v, err := strconv.ParseFloat(tok.RawValue, 64)
			if err != nil {
				return f, fmt.Errorf("invalid escape_timer_multiplier %q: %w", tok.RawValue, err)
			}
			f.EscapeMultiplier = v
			continue
		case "objective_options", "starting_items", "map_layout":
			// These are flag-indicator tokens; the actual values come
			// from subsequent tokens under the same flag.
			continue
		}

		// Skill presets are unique across every alias table, so they're
		// recognized standalone by value, regardless of any sticky flag.
		if skillPresetValues[tok.ID] {
			f.SkillPreset = tok.ID
			continue
		}

		// Progression presets are likewise unique across every alias table.
		if progressionPresetValues[tok.ID] {
			f.ProgressionPreset = tok.ID
			continue
		}

		// Route value tokens by their flag
		switch tok.Flag {
		case 'o':
			if objectivePresetValues[tok.ID] {
				f.ObjectivePreset = tok.ID
			} else {
				f.ObjectiveOptions = append(f.ObjectiveOptions, models.ObjectiveOption{
					Objective: tok.ID,
					Setting:   tok.Value,
				})
			}
		case 's':
			count := 1
			if tok.RawValue != "" {
				n, err := strconv.Atoi(tok.RawValue)
				if err != nil {
					return f, fmt.Errorf("invalid starting item count %q: %w", tok.RawValue, err)
				}
				count = n
			}
			if tok.Value != models.False {
				f.StartingItems = append(f.StartingItems, models.StartingItem{
					Item:  tok.ID,
					Count: count,
				})
			}
		case 'l':
			f.MapLayout = tok.ID
		}
	}
	return f, nil
}

// applyPresetFields writes the populated PresetFields back into the template map.
// Scalar fields use SetNestedValue. Slice fields are merged into existing
// template arrays. TriState values are converted to strings.
func applyPresetFields(m map[string]interface{}, f models.PresetFields) error {
	// --- Scalar fields ---
	if f.IsRace {
		SetNestedValue(m, "other_settings.race_mode", true)
	}
	if f.Version != 0 {
		SetNestedValue(m, "version", f.Version)
	}
	if f.EscapeMultiplier != 0 {
		SetNestedValue(m, "skill_assumption_settings.escape_timer_multiplier", f.EscapeMultiplier)
		SetNestedValue(m, "skill_assumption_settings.preset", nil)
	}
	if f.MapLayout != "" {
		SetNestedValue(m, "map_layout", f.MapLayout)
	}
	if f.SaveAnimals != 0 {
		SetNestedValue(m, "save_animals", f.SaveAnimals.String())
	}
	if f.FreeShinesparks {
		SetNestedValue(m, "other_settings.energy_free_shinesparks", true)
	}
	if f.WallJump != "" {
		SetNestedValue(m, "other_settings.wall_jump", f.WallJump)
	}
	if f.SplitSpeed != "" {
		SetNestedValue(m, "other_settings.speed_booster", f.SplitSpeed)
	}

	// --- Starting Items (merge into template array) ---
	if len(f.StartingItems) > 0 {
		mergeStartingItems(m, f.StartingItems)
		SetNestedValue(m, "item_progression_settings.starting_items_preset", nil)
		SetNestedValue(m, "item_progression_settings.preset", nil)
		applyStartingItemSideEffects(m, f)
	}
	if f.StartingPreset != "" {
		SetNestedValue(m, "item_progression_settings.starting_items_preset", f.StartingPreset)
		SetNestedValue(m, "item_progression_settings.preset", nil)
	}

	// --- Objective Options (merge into template array) ---
	if len(f.ObjectiveOptions) > 0 {
		mergeObjectiveOptions(m, f.ObjectiveOptions)
		SetNestedValue(m, "objective_settings.preset", nil)
	}

	// --- Named presets (win over the nil-outs above: an explicit preset
	// selection takes priority over sibling overrides in the same command) ---
	if f.SkillPreset != "" {
		SetNestedValue(m, "skill_assumption_settings.preset", f.SkillPreset)
	}
	if f.ObjectivePreset != "" {
		SetNestedValue(m, "objective_settings.preset", f.ObjectivePreset)
	}
	if f.ProgressionPreset != "" {
		SetNestedValue(m, "item_progression_settings.preset", f.ProgressionPreset)
	}

	return nil
}

// applyStartingItemSideEffects sets the settings that starting with certain
// movement items implies: Wall Jump must be collectible if it isn't already
// on the player, and either speed booster color requires the ability to
// split the beam from the run, so its setting must be split as well.
// Explicit flags for these settings take priority if present.
func applyStartingItemSideEffects(m map[string]interface{}, f models.PresetFields) {
	for _, si := range f.StartingItems {
		switch si.Item {
		case "WallJump":
			if f.WallJump == "" {
				SetNestedValue(m, "other_settings.wall_jump", "Collectible")
			}
		case "BlueBooster", "SparkBooster":
			if f.SplitSpeed == "" {
				SetNestedValue(m, "other_settings.speed_booster", "Split")
			}
		}
	}
}

// mergeStartingItems updates matching entries in the template's starting_items array.
func mergeStartingItems(m map[string]interface{}, items []models.StartingItem) {
	ips, ok := m["item_progression_settings"].(map[string]interface{})
	if !ok {
		return
	}
	arr, ok := ips["starting_items"].([]interface{})
	if !ok {
		return
	}
	for _, si := range items {
		for _, entry := range arr {
			e, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}
			if e["item"] == si.Item {
				e["count"] = float64(si.Count) // JSON numbers are float64
				break
			}
		}
	}
}

// mergeObjectiveOptions updates matching entries in the template's objective_options array.
func mergeObjectiveOptions(m map[string]interface{}, opts []models.ObjectiveOption) {
	os, ok := m["objective_settings"].(map[string]interface{})
	if !ok {
		return
	}
	arr, ok := os["objective_options"].([]interface{})
	if !ok {
		return
	}
	for _, oo := range opts {
		for _, entry := range arr {
			e, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}
			if e["objective"] == oo.Objective {
				e["setting"] = oo.Setting.String()
				break
			}
		}
	}
}

// postprocess applies cross-field fixups after the main apply step.
func postprocess(m map[string]interface{}, f models.PresetFields) {
	// Clamp objective counts based on "Yes" entries
	if len(f.ObjectiveOptions) > 0 {
		clampObjectiveCounts(m)
	}
}

// clampObjectiveCounts counts objectives set to "Yes" and updates min/max.
func clampObjectiveCounts(m map[string]interface{}) {
	os, ok := m["objective_settings"].(map[string]interface{})
	if !ok {
		return
	}
	arr, ok := os["objective_options"].([]interface{})
	if !ok {
		return
	}
	count := 0
	for _, entry := range arr {
		e, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		if e["setting"] == "Yes" {
			count++
		}
	}
	if count < 1 {
		count = 1
	}
	if count > 19 {
		count = 19
	}
	os["min_objectives"] = float64(count)
	os["max_objectives"] = float64(count)
}

func SetNestedValue(m map[string]interface{}, path string, value interface{}) {
	parts := strings.Split(path, ".")
	for i := 0; i < len(parts)-1; i++ {
		key := parts[i]
		if _, ok := m[key]; !ok {
			m[key] = make(map[string]interface{})
		}
		m = m[key].(map[string]interface{})
	}
	m[parts[len(parts)-1]] = value
}
