package workflow

import (
	"fmt"
	"math/rand/v2"

	"maprandoseedroller/lib/models"
)

// PresetAction runs custom logic for a specific preset before hydration. It may
// mutate tmpl in place (e.g. to override a settings section) and can return an
// arbitrary payload that gets attached to the roll response as Extra.
type PresetAction func(tmpl map[string]interface{}) (interface{}, error)

var PresetActions = map[string]PresetAction{
	"nis": nisAction,
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
