package lib

import (
	"os"
	"sort"
	"strings"

	"maprandoseedroller/lib/models"
)

func BuildSpoilerToken() string {
	return os.Getenv("SPOILER_TOKEN")
}

func BuildSite(isDev bool) string {
	if isDev {
		return "https://dev.maprando.com"
	}
	return "https://maprando.com"
}

func MergeAndSortAliases(tables ...map[string]string) []models.AliasEntry {
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
