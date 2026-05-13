package models

type TriState int

const (
	Maybe TriState = iota
	True
	False
)

func (t TriState) String() string {
	switch t {
	case True:
		return "Yes"
	case False:
		return "No"
	default:
		return "Maybe"
	}
}

type Token struct {
	Flag     rune
	ID       string
	Value    TriState
	RawValue string // trailing numeric/text argument, e.g. "0.95" or "3"
}
