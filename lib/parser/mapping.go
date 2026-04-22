package parser

import (
	"encoding/json"
	"fmt"
	"maprandoseedroller/lib/models"
	"reflect"
	"strings"
)

// Hydrate converts tokens into preset overrides and applies them onto the template.
// Returns the merged preset as JSON bytes ready for the MapRando API.
func Hydrate(template map[string]interface{}, tokens []models.Token) ([]byte, error) {
	// --- Step 1: tokens → PresetFields ---
	fields, err := tokensToPresetFields(tokens)
	if err != nil {
		return nil, fmt.Errorf("tokensToPresetFields: %w", err)
	}

	// --- Step 2: apply PresetFields onto template ---
	// Not 1:1: some fields drive multiple JSON paths, others require
	// lookup tables, and complex fields (slices, enums) need custom logic.
	if err := applyPresetFields(template, fields); err != nil {
		return nil, fmt.Errorf("applyPresetFields: %w", err)
	}

	// --- Step 3: postprocess ---
	// e.g. clamp objectives, resolve preset name collisions, etc.
	postprocess(template, fields)

	return json.Marshal(template)
}

// tokensToPresetFields maps a flat token list onto the PresetFields struct.
func tokensToPresetFields(tokens []models.Token) (models.PresetFields, error) {
	var f models.PresetFields
	for _, tok := range tokens {
		switch tok.ID {
		case "race_mode":
			f.IsRace = tok.Value == models.True
		case "escape_timer_multiplier":
			// TODO: map TriState → float64 scalar
		case "map_layout":
			// TODO: map TriState → layout string
		case "starting_items":
			// TODO: accumulate StartingItem entries
		case "objective_options":
			f.ObjectiveOptions = append(f.ObjectiveOptions, models.ObjectiveOption{
				Objective: tok.ID,
				Setting:   tok.Value,
			})
		}
	}
	return f, nil
}

// applyPresetFields writes the populated PresetFields back into the template
// map using the `path` struct tags. Complex/non-scalar fields are handled
// individually because they are not direct 1:1 assignments.
func applyPresetFields(m map[string]interface{}, f models.PresetFields) error {
	v := reflect.ValueOf(f)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		path := field.Tag.Get("path")
		if path == "" {
			continue
		}
		// TODO: handle slice fields (StartingItems, ObjectiveOptions) specially;
		// for now only write scalar fields that are non-zero.
		fval := v.Field(i)
		if fval.IsZero() {
			continue
		}
		SetNestedValue(m, path, fval.Interface())
	}
	return nil
}

// postprocess applies any cross-field fixups after the main apply step.
// e.g. clamping min/max objectives, resolving preset label conflicts.
func postprocess(m map[string]interface{}, f models.PresetFields) {
	// TODO
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
