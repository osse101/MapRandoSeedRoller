package workflow

import (
	"fmt"
	"math/rand/v2"

	"maprandoseedroller/lib/models"
	"maprandoseedroller/preset"
)

// PresetAction runs custom logic for a specific preset before hydration. It may
// mutate tmpl in place (e.g. to override a settings section) and can return an
// arbitrary payload that gets attached to the roll response as Extra.
type PresetAction func(tmpl map[string]interface{}) (interface{}, error)

// PresetActions is populated in init() rather than via a var initializer
// because randomPresetAction reads PresetActions itself (to chain into the
// action of whatever preset it randomly picks); referencing the map from an
// initializer expression that also builds the map creates a (false-positive)
// initialization cycle as far as the compiler is concerned.
var PresetActions = map[string]PresetAction{}

func init() {
	PresetActions["nis"] = nisAction
	PresetActions["drockyrandom"] = drockyRandomAction
	PresetActions["mentorrandom"] = mentorRandomAction
	PresetActions["randompreset"] = randomPresetAction
}

var nisItemProgressionPresets = []string{"Technical", "Challenge", "Desolate"}

// Placeholder pool — swap for real sprite names later.
var nisSprites = []string{
	"Dread Samus",
	"Captain Novolin",
	"Samus Cannon",
	"Moonclif",
	"Officer Donut",
	"SNES Controller",
	"The King of Pop",
	"Invisible Samus",
	"Left Leg Samus",
	"Upside-Down Samus",
	"180-Degree Samus",
	"Tetromino",
	"Space Pirate",
	"Shaktool",
	"140",
	"Metroid",
	"Master Hand",
	"Mario (NES)",
	"Luigi",
	"Link (Z2)",
	"Kirby (Yarn)",
	"Kirby",
	"Infee and Nitee",
	"Elista",
	"Diddy Kong",
	"Cursor",
	"Crewmate",
	"Cacodemon",
}

func nisAction(tmpl map[string]interface{}) (interface{}, error) {
	ips, ok := tmpl["item_progression_settings"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("nis action: item_progression_settings missing or malformed")
	}
	ips["preset"] = nisItemProgressionPresets[rand.IntN(len(nisItemProgressionPresets))]

	return models.SpriteExtra{
		SpriteName: nisSprites[rand.IntN(len(nisSprites))],
	}, nil
}

// loadTemplateInto replaces tmpl's contents with the named template's
// contents. Since maps are reference types, this mutation is visible to the
// caller (roll.go) after the action returns, effectively swapping the
// template that will be hydrated and rolled — without needing to change the
// PresetAction signature to return a new map.
func loadTemplateInto(tmpl map[string]interface{}, name string) error {
	loaded, err := preset.LoadTemplate(name)
	if err != nil {
		return err
	}
	for k := range tmpl {
		delete(tmpl, k)
	}
	for k, v := range loaded {
		tmpl[k] = v
	}
	return nil
}

// starterItemExclusions are items that aren't valid picks for a random
// starting-item draft (ammo/tank items and the non-suit movement extras).
var starterItemExclusions = map[string]bool{
	"Missile":      true,
	"ETank":        true,
	"ReserveTank":  true,
	"Super":        true,
	"PowerBomb":    true,
	"WallJump":     true,
	"BlueBooster":  true,
	"SparkBooster": true,
}

// randomizeStartingItems sets n distinct, randomly chosen (eligible)
// starting_items entries in tmpl to count 1, matching the default count the
// flag lexer uses for a bare -s flag (lib/parser/mapping.go), and returns
// the names of the items that were selected.
func randomizeStartingItems(tmpl map[string]interface{}, n int) ([]string, error) {
	ips, ok := tmpl["item_progression_settings"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("item_progression_settings missing or malformed")
	}
	arr, ok := ips["starting_items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("starting_items missing or malformed")
	}

	var eligible []int
	for i, raw := range arr {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		item, _ := entry["item"].(string)
		if !starterItemExclusions[item] {
			eligible = append(eligible, i)
		}
	}
	if n > len(eligible) {
		n = len(eligible)
	}
	perm := rand.Perm(len(eligible))
	selected := make([]string, 0, n)
	for _, p := range perm[:n] {
		entry, ok := arr[eligible[p]].(map[string]interface{})
		if !ok {
			continue
		}
		entry["count"] = float64(1)
		if item, ok := entry["item"].(string); ok {
			selected = append(selected, item)
		}
	}
	ips["starting_items_preset"] = nil
	return selected, nil
}

var winterTournamentPresets = []string{"objectives", "suits", "g91", "draft", "metroids", "noobjectives", "vmove"}

func drockyRandomAction(tmpl map[string]interface{}) (interface{}, error) {
	//Select preset from Winter_Tournament series of presets
	//  if Item_Draft is selected, select 2 different starting items at random
	chosen := winterTournamentPresets[rand.IntN(len(winterTournamentPresets))]
	if err := loadTemplateInto(tmpl, chosen); err != nil {
		return nil, fmt.Errorf("drockyrandom action: %w", err)
	}
	extra := models.PresetSelectionExtra{SelectedPreset: chosen}
	if chosen == "draft" {
		items, err := randomizeStartingItems(tmpl, 2)
		if err != nil {
			return nil, fmt.Errorf("drockyrandom action: %w", err)
		}
		extra.StartingItems = items
	}
	return extra, nil
}

var mentorObjectivePresets = []string{"Minibosses", "Pirates", "Metroids", "Chozos"}

func mentorRandomAction(tmpl map[string]interface{}) (interface{}, error) {
	//Select preset from Mentor Tournament variations:
	//  Item draft, miniboss objective, pirate objective, metroid objective, chozo objective
	//  if Item_Draft is selected, select 2 different starting items at random
	if err := loadTemplateInto(tmpl, "mentor"); err != nil {
		return nil, fmt.Errorf("mentorrandom action: %w", err)
	}
	if rand.IntN(5) == 0 { // 1-in-5: Item Draft
		items, err := randomizeStartingItems(tmpl, 2)
		if err != nil {
			return nil, fmt.Errorf("mentorrandom action: %w", err)
		}
		return models.PresetSelectionExtra{SelectedPreset: "ItemDraft", StartingItems: items}, nil
	}
	os, ok := tmpl["objective_settings"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("mentorrandom action: objective_settings missing or malformed")
	}
	chosen := mentorObjectivePresets[rand.IntN(len(mentorObjectivePresets))]
	os["preset"] = chosen
	return models.PresetSelectionExtra{SelectedPreset: chosen}, nil
}

func randomPresetAction(tmpl map[string]interface{}) (interface{}, error) {
	//Select preset from all options at random
	names := preset.GetPresetNames()
	chosen := names[rand.IntN(len(names))]
	if err := loadTemplateInto(tmpl, chosen); err != nil {
		return nil, fmt.Errorf("randompreset action: %w", err)
	}
	extra := models.PresetSelectionExtra{SelectedPreset: chosen}
	// Chain into the chosen preset's own action, if it has one (e.g. "nis"),
	// so randompreset is truly "any option" rather than just any raw
	// template. The nested action's own Extra is preserved underneath,
	// rather than replacing the fact that randompreset picked it.
	if action, ok := PresetActions[chosen]; ok {
		nested, err := action(tmpl)
		if err != nil {
			return nil, fmt.Errorf("randompreset action: %w", err)
		}
		extra.Extra = nested
	}
	return extra, nil
}
