package parser

import (
	"strings"
	"unicode"

	"maprandoseedroller/lib/models"
)

func Lex(input string, aliases []models.AliasEntry) ([]models.Token, error) {
	tokens := []models.Token{}
	lastFlag := rune(-1)

	for len(input) > 0 {
		matchFound := false

		for _, entry := range aliases {
			// Case-insensitive prefix check
			if len(input) >= len(entry.ShortName) &&
				strings.EqualFold(input[:len(entry.ShortName)], entry.ShortName) {
				// Extract the specific chunk (e.g., "Dray")
				chunk := input[:len(entry.ShortName)]

				// Check if flag or alias
				if len(chunk) == 1 {
					lastFlag = unicode.ToLower(rune(chunk[0]))
				}

				// Determine TriState based on the chunk's casing
				state := DetermineTriState(chunk)

				// Create token
				tokens = append(tokens, models.Token{
					Flag:  lastFlag,
					ID:    entry.LongName,
					Value: state,
				})

				// Advance input pointer
				input = input[len(entry.ShortName):]
				matchFound = true
				break
			}
		}

		// If no alias matches the current start of the string
		if !matchFound {
			input = input[1:]
			continue
		}
	}
	return tokens, nil
}

func DetermineTriState(s string) models.TriState {
	if len(s) == 0 {
		return models.Maybe
	}

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
