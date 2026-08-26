package workflow

import (
	"encoding/json"
	"slices"
	"testing"

	"maprandoseedroller/lib/models"
	"maprandoseedroller/preset"
)

func TestPrepareGameData(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "Valid s4 preset",
			data:    "s4",
			wantErr: false,
		},
		{
			name:    "Default values",
			data:    "",
			wantErr: false,
		},
		{
			name:    "Invalid preset",
			data:    "invalid-preset",
			wantErr: true,
		},
		{
			name:    "nis preset (in-place action)",
			data:    "nis",
			wantErr: false,
		},
		{
			name:    "drockyrandom meta preset",
			data:    "drockyrandom",
			wantErr: false,
		},
		{
			name:    "mentorrandom meta preset",
			data:    "mentorrandom",
			wantErr: false,
		},
		{
			name:    "randompreset meta preset",
			data:    "randompreset",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, _, err := PrepareGameData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrepareGameData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) == 0 {
				t.Errorf("PrepareGameData() returned empty data for valid request")
			}
		})
	}
}

// TestDrockyRandomAction runs the action many times to exercise every branch
// (including the low-probability Item Draft + starting-item selection),
// verifying Extra reports the selection accurately and matches the resulting
// template.
func TestDrockyRandomAction(t *testing.T) {
	for i := 0; i < 50; i++ {
		gameData, _, extra, err := PrepareGameData("drockyrandom")
		if err != nil {
			t.Fatalf("PrepareGameData(drockyrandom) error = %v", err)
		}
		sel, ok := extra.(models.PresetSelectionExtra)
		if !ok {
			t.Fatalf("Extra is %T, want models.PresetSelectionExtra", extra)
		}
		if !slices.Contains(winterTournamentPresets, sel.SelectedPreset) {
			t.Fatalf("unexpected selected preset %q, not a Winter Tournament preset", sel.SelectedPreset)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal(gameData, &parsed); err != nil {
			t.Fatalf("failed to unmarshal game data: %v", err)
		}
		wantName, err := templateDisplayName(sel.SelectedPreset)
		if err != nil {
			t.Fatal(err)
		}
		if name, _ := parsed["name"].(string); name != wantName {
			t.Fatalf("Extra.SelectedPreset %q but template name is %q, want %q", sel.SelectedPreset, name, wantName)
		}

		if sel.SelectedPreset == "draft" {
			assertValidStartingItemDraft(t, parsed, sel.StartingItems)
		} else if len(sel.StartingItems) != 0 {
			t.Fatalf("StartingItems set for non-draft selection %q: %v", sel.SelectedPreset, sel.StartingItems)
		}
	}
}

// TestMentorRandomAction runs the action many times to exercise all 5
// variations, verifying Extra reports the selection accurately and matches
// the resulting template.
func TestMentorRandomAction(t *testing.T) {
	for i := 0; i < 50; i++ {
		gameData, _, extra, err := PrepareGameData("mentorrandom")
		if err != nil {
			t.Fatalf("PrepareGameData(mentorrandom) error = %v", err)
		}
		sel, ok := extra.(models.PresetSelectionExtra)
		if !ok {
			t.Fatalf("Extra is %T, want models.PresetSelectionExtra", extra)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal(gameData, &parsed); err != nil {
			t.Fatalf("failed to unmarshal game data: %v", err)
		}
		if name, _ := parsed["name"].(string); name != "Mentor Tournament" {
			t.Fatalf("unexpected template name %q, expected Mentor Tournament", name)
		}

		if sel.SelectedPreset == "ItemDraft" {
			assertValidStartingItemDraft(t, parsed, sel.StartingItems)
			continue
		}
		if !slices.Contains(mentorObjectivePresets, sel.SelectedPreset) {
			t.Fatalf("unexpected selected preset %q", sel.SelectedPreset)
		}
		if len(sel.StartingItems) != 0 {
			t.Fatalf("StartingItems set for non-draft selection %q: %v", sel.SelectedPreset, sel.StartingItems)
		}
		os, ok := parsed["objective_settings"].(map[string]interface{})
		if !ok {
			t.Fatalf("objective_settings missing or malformed")
		}
		if objPreset, _ := os["preset"].(string); objPreset != sel.SelectedPreset {
			t.Fatalf("Extra.SelectedPreset %q but template objective preset is %q", sel.SelectedPreset, objPreset)
		}
	}
}

// TestRandomPresetAction confirms the meta preset always resolves to a real,
// registered template, that Extra reports which one was picked, and that a
// chained action's own Extra (e.g. nis's sprite) is preserved underneath
// rather than discarded.
func TestRandomPresetAction(t *testing.T) {
	validNames := preset.GetPresetNames()
	for i := 0; i < 40; i++ {
		gameData, _, extra, err := PrepareGameData("randompreset")
		if err != nil {
			t.Fatalf("PrepareGameData(randompreset) error = %v", err)
		}
		if len(gameData) == 0 {
			t.Fatalf("PrepareGameData(randompreset) returned empty data")
		}
		sel, ok := extra.(models.PresetSelectionExtra)
		if !ok {
			t.Fatalf("Extra is %T, want models.PresetSelectionExtra", extra)
		}
		if !slices.Contains(validNames, sel.SelectedPreset) {
			t.Fatalf("unexpected selected preset %q", sel.SelectedPreset)
		}
		if sel.SelectedPreset == "nis" {
			nested, ok := sel.Extra.(models.SpriteExtra)
			if !ok || nested.SpriteName == "" {
				t.Fatalf("expected nested nis SpriteExtra, got %#v", sel.Extra)
			}
		}
	}
}

// templateDisplayName loads the real preset's JSON to read its display
// "name" field, so tests can cross-check Extra against the actual template.
func templateDisplayName(shortName string) (string, error) {
	tmpl, err := preset.LoadTemplate(shortName)
	if err != nil {
		return "", err
	}
	name, _ := tmpl["name"].(string)
	return name, nil
}

// assertValidStartingItemDraft checks that Extra reported exactly 2 distinct,
// non-excluded starting items, and that the template's starting_items array
// has exactly those items (and no others) set to a nonzero count.
func assertValidStartingItemDraft(t *testing.T, parsed map[string]interface{}, extraItems []string) {
	t.Helper()
	if len(extraItems) != 2 {
		t.Fatalf("expected Extra to report 2 starting items, got %v", extraItems)
	}
	if extraItems[0] == extraItems[1] {
		t.Fatalf("expected 2 distinct starting items, got duplicate %q", extraItems[0])
	}
	for _, item := range extraItems {
		if starterItemExclusions[item] {
			t.Fatalf("excluded item %q reported in Extra.StartingItems", item)
		}
	}

	ips, ok := parsed["item_progression_settings"].(map[string]interface{})
	if !ok {
		t.Fatalf("item_progression_settings missing or malformed")
	}
	arr, ok := ips["starting_items"].([]interface{})
	if !ok {
		t.Fatalf("starting_items missing or malformed")
	}
	var selected []string
	for _, raw := range arr {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if count, _ := entry["count"].(float64); count != 0 {
			item, _ := entry["item"].(string)
			selected = append(selected, item)
		}
	}
	if !slices.Equal(slices.Sorted(slices.Values(selected)), slices.Sorted(slices.Values(extraItems))) {
		t.Fatalf("template starting_items %v does not match Extra.StartingItems %v", selected, extraItems)
	}
}
