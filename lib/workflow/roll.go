package workflow

import (
	"maprandoseedroller/lib/models"
	"maprandoseedroller/lib/parser"
	"maprandoseedroller/lib/randomize"
	"maprandoseedroller/preset"
	"sort"
	"strings"
)

func ExecuteRoll(req models.RequestIn)(models.ResponseOut, error){
	//Parse flags
	flagTable, err := mergeAndSortAliases();
	if err != nil{
		return models.ResponseOut{}, err
	}

	tokens, err := parser.Lex(req.Flags, flagTable);
	if err != nil{
		return models.ResponseOut{}, err
	}
	
	//Write json preset
	tmpl, err := preset.LoadTemplate(req.Preset)
	if err != nil {
		return models.ResponseOut{}, err
	}
	gameData, err := parser.Hydrate(tmpl, tokens)
	if err != nil {
		return models.ResponseOut{}, err
	}

	//Send to MapRando
	seed_url, err :=  randomize.Randomize(gameData);
	if err != nil{
		return models.ResponseOut{}, err
	}

	//Determine Discord/Racetime fields
	var resp = models.ResponseOut{
		SeedURL : seed_url,
	}

	return resp, nil
}

func mergeAndSortAliases(tables ...map[string]string)([]models.AliasEntry, error){
	tempMap := make(map[string]string)
	// Merge all maps into one to handle overrides/deduplication
	for _, table := range tables{
		for alias, id := range table{
			tempMap[strings.ToLower(alias)] = id
		}
	}

	// Convert to a slice
	var list []models.AliasEntry
	for alias, id := range tempMap{
		list = append(list, models.AliasEntry{
			ShortName: alias,
			LongName: id,
		})
	}
	
	//Sort by length DESCENDING
	sort.Slice(list, func(i, j int)bool{
		return len(list[i].ShortName) > len(list[j].ShortName)
	})
	return list, nil
}