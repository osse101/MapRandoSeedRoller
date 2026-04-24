package workflow

import (
	"maprandoseedroller/lib"
	"maprandoseedroller/lib/models"
	"maprandoseedroller/lib/parser"
	"maprandoseedroller/lib/randomize"
	"maprandoseedroller/preset"
)

func ExecuteRoll(req models.RequestIn) (models.ResponseOut, error) {
	gameData, isDev, err := PrepareGameData(req)
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

// PrepareGameData handles lexing flags, loading templates, and hydrating game data.
func PrepareGameData(req models.RequestIn) ([]byte, bool, error) {
	//Parse flags
	flagTable := lib.MergeAndSortAliases(
		models.ObjectiveAliases,
		models.ItemAliases,
		models.FlagAliases,
	)

	tokens, err := parser.Lex(req.Flags, flagTable)
	if err != nil {
		return nil, false, err
	}

	//Write json preset
	tmpl, err := preset.LoadTemplate(req.Preset)
	if err != nil {
		return nil, false, err
	}
	return parser.Hydrate(tmpl, tokens)
}
