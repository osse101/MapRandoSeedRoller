package models

import "strings"

type TriState int

const (
	Maybe TriState = iota
	True
	False
)

type Token struct{
	Flag rune
	ID string
	Value TriState
}

func DetermineState(input string) TriState{
	if input == strings.ToUpper(input) {return True}
	if input == strings.ToLower(input) {return False}
	return Maybe
}