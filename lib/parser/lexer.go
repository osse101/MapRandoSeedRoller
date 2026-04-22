package parser

import (
	"strings"
	"unicode"

	"maprandoseedroller/lib/models"
)

func Lex(input string, aliases []models.AliasEntry) ([]models.Token, error) {
	var tokens []models.Token

	for len(input) > 0 {
		matchFound := false

		for _, entry := range aliases {
			// Case-insensitive prefix check
			if len(input) >= len(entry.ShortName) &&
				strings.EqualFold(input[:len(entry.ShortName)], entry.ShortName) {
				// 1. Extract the specific chunk (e.g., "Dray")
				chunk := input[:len(entry.ShortName)]

				// 2. Determine TriState based on the chunk's casing
				state := determineTriState(chunk)

				// 3. Create token
				tokens = append(tokens, models.Token{
					ID:    entry.LongName,
					Value: state,
				})

				// 4. Advance input pointer
				input = input[len(entry.ShortName):]
				matchFound = true
				break
			}
		}

		// If no alias matches the current start of the string, skip 1 char
		if !matchFound {
			// You could log a warning here for the "Warning Level" feedback
			input = input[1:]
		}
	}
	return tokens, nil
}

func determineTriState(s string) models.TriState {
	isAllUpper := true
	isAllLower := true

	for _, r := range s {
		if unicode.IsLower(r) {
			isAllUpper = false
		} else if unicode.IsUpper(r) {
			isAllLower = false
		}
	}

	if isAllUpper {
		return models.True // e.g., "KRAID"
	}
	if isAllLower {
		return models.False // e.g., "kraid"
	}
	return models.Maybe // e.g., "Kraid" or "KraId"
}
