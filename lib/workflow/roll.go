package workflow

import (
	"sort"
	"strings"

	"maprandoseedroller/lib/models"
	"maprandoseedroller/lib/parser"
	"maprandoseedroller/lib/randomize"
	"maprandoseedroller/preset"
)

func ExecuteRoll(req models.RequestIn) (models.ResponseOut, error) {
	gameData, err := PrepareGameData(req)
	if err != nil {
		return models.ResponseOut{}, err
	}

	//Send to MapRando
	seedURL, err := randomize.Randomize(gameData)
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
func PrepareGameData(req models.RequestIn) ([]byte, error) {
	//Parse flags
	flagTable := mergeAndSortAliases()

	tokens, err := parser.Lex(req.Flags, flagTable)
	if err != nil {
		return nil, err
	}

	//Write json preset
	tmpl, err := preset.LoadTemplate(req.Preset)
	if err != nil {
		return nil, err
	}
	return parser.Hydrate(tmpl, tokens)
}

func mergeAndSortAliases(tables ...map[string]string) []models.AliasEntry {
	tempMap := make(map[string]string)
	// Merge all maps into one to handle overrides/deduplication
	for _, table := range tables {
		for alias, id := range table {
			tempMap[strings.ToLower(alias)] = id
		}
	}

	// Convert to a slice
	list := make([]models.AliasEntry, 0, len(tempMap))
	for alias, id := range tempMap {
		list = append(list, models.AliasEntry{
			ShortName: alias,
			LongName:  id,
		})
	}

	//Sort by length DESCENDING
	sort.Slice(list, func(i, j int) bool {
		return len(list[i].ShortName) > len(list[j].ShortName)
	})
	return list
}
